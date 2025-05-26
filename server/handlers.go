package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"server/auth"
	"server/database"

	_ "github.com/glebarez/go-sqlite"
)

func prepareResponse(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json")
}

type Error struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	ErrorType  string `json:"error_type"`
}

func sendError(w http.ResponseWriter, error Error, err error) {
	log.Println(err)
	if err := json.NewEncoder(w).Encode(error); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *app) login(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	input := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	output := struct {
		Token string `json:"token"`
		Role  int    `json:"role"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	password, err := database.GetPassword(app.DB, input.Email)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	if password != input.Password || password == "" {
		sendError(w, Error{401, "Wrong password or login", "Unauthorized"}, err)
		return
	}

	output.Token, err = auth.CreateSession(app.CACHE, input.Email)
	if err != nil {
		sendError(w, Error{500, "Could not generate a new token", "Internal Server Error"}, err)
		return
	}

	output.Role, err = database.GetRole(app.DB, input.Email)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	log.Printf("Login -- User: %s - Token: %s", input.Email, output.Token)

	if err := json.NewEncoder(w).Encode(output); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *app) logout(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	token := struct {
		Token string `json:"token"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&token)
	if err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	database.DeleteToken(app.CACHE, token.Token)

	w.WriteHeader(http.StatusOK)
}

func (app *app) info(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	var err error
	output := struct {
		Id    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
		Role  int    `json:"role"`
	}{}

	r.ParseMultipartForm(32 << 20)
	token := r.FormValue("token")

	output.Email, err = auth.ValidateSession(app.CACHE, token)
	if err != nil {
		sendError(w, Error{401, "Incorrect Token", "Unauthorized"}, err)
		return
	}

	output.Id, output.Name, output.Phone, output.Role, err = database.GetInfo(app.DB, output.Email)

	if err := json.NewEncoder(w).Encode(output); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *app) tenantInfo(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	output := struct {
		ApartamentId int     `json:"apartament_id"`
		Rent         float32 `json:"rent"`
	}{}

	r.ParseMultipartForm(32 << 20)
	token := r.FormValue("token")

	if erro := app.checkRole(token, 2); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	email, _ := auth.ValidateSession(app.CACHE, token)
	// change

	id, err := database.GetId(app.DB, email)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	output.ApartamentId, err = database.GetApartamentId(app.DB, id)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	output.Rent, err = database.GetRent(app.DB, output.ApartamentId)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	err = json.NewEncoder(w).Encode(&output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *app) subInfo(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	output := struct {
		Address      string `json:"address"`
		NIP          string `json:"NIP"`
		SpecialityId int    `json:"speciality_id"`
	}{}

	r.ParseMultipartForm(32 << 20)
	token := r.FormValue("token")

	erro := app.checkRole(token, 3)
	if erro != nil {
		sendError(w, *erro, nil)
		return
	}

	email, _ := auth.ValidateSession(app.CACHE, token)
	// change

	id, err := database.GetId(app.DB, email)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	output.Address, output.NIP, output.SpecialityId, err = database.GetSubconInfo(app.DB, id)

	if err = json.NewEncoder(w).Encode(&output); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (app *app) getTenants(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	type tenent struct {
		Id     int    `json:"id"`
		Name   string `json:"name"`
		Email  string `json:"email"`
		Phone  string `json:"phone"`
		RoleId int    `json:"role_id"`
	}

	output := struct {
		Tenants []tenent `json:"apartaments"`
	}{}

	r.ParseMultipartForm(32 << 20)
	token := r.FormValue("token")

	if erro := app.checkRole(token, 1); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	id, name, email, phone, roleId, err := database.GetTenents(app.DB)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	var data []tenent
	for n := range id {
		data = append(data,
			tenent{
				Id:     id[n],
				Name:   name[n],
				Email:  email[n],
				Phone:  phone[n],
				RoleId: roleId[n]})
	}
	output.Tenants = data

	if err = json.NewEncoder(w).Encode(&output); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (app *app) getApartaments(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	type apartament struct {
		Id             int    `json:"id"`
		Name           string `json:"name"`
		Street         string `json:"street"`
		BuildingNumber string `json:"building_number"`
		Buildingname   string `json:"building_name"`
		FlatNumber     string `json:"flat_number"`
		OwnerId        int    `json:"owner_id"`
	}

	output := struct {
		Apartaments []apartament `json:"apartaments"`
	}{}

	r.ParseMultipartForm(32 << 20)
	token := r.FormValue("token")

	if erro := app.checkRole(token, 1); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	id, name, street, buildingNumber, buildingName, flatNumber, ownerId, err := database.GetApartamentsData(app.DB)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	var data []apartament
	for n := range id {
		data = append(data,
			apartament{
				Id:             id[n],
				Name:           name[n],
				Street:         street[n],
				BuildingNumber: buildingNumber[n],
				Buildingname:   buildingName[n],
				FlatNumber:     flatNumber[n],
				OwnerId:        ownerId[n]})
	}
	output.Apartaments = data

	if err = json.NewEncoder(w).Encode(&output); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (app *app) addApartament(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	input := struct {
		Token          string `json:"token"`
		Name           string `json:"name"`
		Street         string `json:"street"`
		BuildingNumber string `json:"building_number"`
		BuildingName   string `json:"building_name"`
		FlatNumber     string `json:"flat_number"`
		OwnerId        int    `json:"owner_id"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	if erro := app.checkRole(input.Token, 1); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	err = database.AddApartament(app.DB, []string{input.Name, input.Street, input.BuildingNumber, input.BuildingName, input.FlatNumber}, input.OwnerId)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *app) getOwners(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	type owner struct {
		Id    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}

	output := struct {
		Owners []owner `json:"owners"`
	}{}

	r.ParseMultipartForm(32 << 20)
	token := r.FormValue("token")
	if erro := app.checkRole(token, 1); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	id, name, email, phone, err := database.GetOwners(app.DB)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	var data []owner
	for n := range id {
		data = append(data,
			owner{
				Id:    id[n],
				Name:  name[n],
				Email: email[n],
				Phone: phone[n]})
	}
	output.Owners = data

	if err = json.NewEncoder(w).Encode(&output); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (app *app) addOwner(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	input := struct {
		Token string `json:"token"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	if erro := app.checkRole(input.Token, 1); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	err = database.AddOwner(app.DB, []string{input.Name, input.Email, input.Phone})
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////

func (app *app) changeRent(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Token        string  `json:"token"`
		ApartamentId int     `json:"apartament_id"`
		Rent         float32 `json:"rent"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	if erro := app.checkRole(data.Token, 1); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	if err = database.ChangeRent(app.DB, data.ApartamentId, data.Rent); err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *app) addUser(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	user := struct {
		Token    string `json:"token"`
		Name     string `json:"name"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Role     int    `json:"role"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	if erro := app.checkRole(user.Token, 1); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	_, err = bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		sendError(w, Error{500, "Could not generate hash from password", "Internal Server Error"}, err)
		return
	}

	err = database.AddUser(app.DB, []string{user.Name, user.Password, user.Email, user.Phone}, user.Role)
	if err != nil {
		sendError(w, Error{500, "Could not add user", "Internal Server Error"}, err)
		return
	}

	user.Password = strings.Repeat("*", len(user.Password)) // should be changed
	if err := json.NewEncoder(w).Encode(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (app *app) getSubContractorSpec(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	data := struct {
		Token string   `json:"token"`
		Email string   `json:"email"`
		Spec  []string `json:"spec"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	if erro := app.checkRole(data.Token, 1); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	data.Spec, err = database.GetSubcontractorSpec(app.DB)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *app) addSubContractorSpec(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	data := struct {
		Token string `json:"token"`
		Name  string `json:"name"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	if erro := app.checkRole(data.Token, 1); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	err = database.AddSpec(app.DB, data.Name)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
	}

	w.WriteHeader(http.StatusOK)
}

func (app *app) checkRole(token string, role_id int) *Error {
	email, err := auth.ValidateSession(app.CACHE, token)
	if err != nil {
		log.Println(err)
		return &Error{401, "Incorrect Token", "Unauthorized"}
	}

	role, err := database.GetRole(app.DB, email)
	if err != nil {
		log.Println(err)
		return &Error{401, "Database", "Internal Server Error"}
	}
	if role != role_id {
		log.Println("Wrong role")
		return &Error{401, "Wrong role", "Unauthorized"}
	}
	return nil
}
