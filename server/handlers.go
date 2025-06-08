package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"server/auth"
	"server/database"
	"server/sqlc"

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
		Role  int64  `json:"role"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	userId, err := app.Query.GetUserId(app.Ctx, input.Email)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	password, err := app.Query.GetUserPassword(app.Ctx, userId)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	if password != input.Password || password == "" {
		sendError(w, Error{401, "Wrong password or login", "Unauthorized"}, err)
		return
	}

	output.Token, err = auth.CreateSession(app.CACHE, int(userId))
	if err != nil {
		sendError(w, Error{500, "Could not generate a new token", "Internal Server Error"}, err)
		return
	}

	output.Role, err = app.Query.GetUserRole(app.Ctx, userId)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	log.Printf("Login -- User: %s - Token: %s - Admin: %d", input.Email, output.Token, output.Role)

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

	r.ParseMultipartForm(32 << 20)
	token := r.FormValue("token")

	userId, err := auth.ValidateSession(app.CACHE, token)
	if err != nil {
		sendError(w, Error{401, "Incorrect Token", "Unauthorized"}, err)
		return
	}

	output, err := app.Query.GetUserInfo(app.Ctx, int64(userId))

	if err := json.NewEncoder(w).Encode(output); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *app) tenantInfo(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	var err error

	output := struct {
		ApartamentId sql.NullInt64 `json:"apartament_id"`
		Rent         float32       `json:"rent"`
	}{}

	r.ParseMultipartForm(32 << 20)
	token := r.FormValue("token")

	if erro := app.checkRole(token, 2); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	userId, _ := auth.ValidateSession(app.CACHE, token)

	//out, err := app.Query.GetApartmentID(app.Ctx, sql.NullInt64{Int64: int64(id)})
	output.ApartamentId, err = app.Query.GetApartmentID(app.Ctx, sql.NullInt64{Int64: int64(userId)})
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	output.Rent, err = database.GetRent(app.DB, int(output.ApartamentId.Int64))
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

	var err error

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

	userId, _ := auth.ValidateSession(app.CACHE, token)

	output.Address, output.NIP, output.SpecialityId, err = database.GetSubconInfo(app.DB, userId)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	if err = json.NewEncoder(w).Encode(&output); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (app *app) getContractors(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	type contractor struct {
		Id           int    `json:"id"`
		UserId       int    `json:"user_id"`
		Address      string `json:"address"`
		NIP          string `json:"NIP"`
		SpecialityId int    `json:"speciality_id"`
	}

	output := struct {
		Contractors []contractor `json:"apartaments"`
	}{}

	r.ParseMultipartForm(32 << 20)
	token := r.FormValue("token")

	if erro := app.checkRole(token, 1); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	id, useId, address, NIP, specialityId, err := database.GetContractorsData(app.DB)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	var data []contractor
	for n := range id {
		data = append(data,
			contractor{
				Id:           id[n],
				UserId:       useId[n],
				Address:      address[n],
				NIP:          NIP[n],
				SpecialityId: specialityId[n]})
	}
	output.Contractors = data

	if err := json.NewEncoder(w).Encode(&output); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (app *app) addContractor(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	input := struct {
		Token        string `json:"token"`
		UserId       int    `json:"user_id"`
		Address      string `json:"address"`
		NIP          string `json:"NIP"`
		SpecialityId int    `json:"speciality_id"`
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

	err = database.AddContractor(app.DB, []string{input.Address, input.NIP}, input.UserId, input.SpecialityId)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	w.WriteHeader(http.StatusOK)
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
		Id             int     `json:"id"`
		Name           string  `json:"name"`
		Street         string  `json:"street"`
		BuildingNumber string  `json:"building_number"`
		Buildingname   string  `json:"building_name"`
		FlatNumber     string  `json:"flat_number"`
		OwnerId        int     `json:"owner_id"`
		Rent           float32 `json:"rent"`
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
		rent, err := database.GetRent(app.DB, id[n])
		if err != nil {
			sendError(w, Error{400, "Database", "Internal Server Error"}, err)
			return
		}
		data = append(data,
			apartament{
				Id:             id[n],
				Name:           name[n],
				Street:         street[n],
				BuildingNumber: buildingNumber[n],
				Buildingname:   buildingName[n],
				FlatNumber:     flatNumber[n],
				OwnerId:        ownerId[n],
				Rent:           rent})
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

func (app *app) getCurrentRenting(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	type renting struct {
		Id           int    `json:"id"`
		ApartamentId int    `json:"apartament_id"`
		UserId       int    `json:"user_id"`
		StartDate    string `json:"start_date"`
	}

	output := struct {
		Rentings []renting `json:"owners"`
	}{}

	r.ParseMultipartForm(32 << 20)
	token := r.FormValue("token")
	if erro := app.checkRole(token, 1); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	id, apartamentId, userId, startDate, err := database.GetActiveRentings(app.DB)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	var data []renting
	for n := range id {
		data = append(data,
			renting{ApartamentId: apartamentId[n],
				UserId:    userId[n],
				StartDate: startDate[n]})
	}
	output.Rentings = data

	if err = json.NewEncoder(w).Encode(&output); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (app *app) addNewRenting(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	input := struct {
		Token        string `json:"token"`
		ApartamentId int    `json:"apartament_id"`
		UserId       int    `json:"user_id"`
		StartDate    string `json:"start_date"` // DateOnly   = "2006-01-02"
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

	err = database.AddNewRenting(app.DB, input.ApartamentId, input.UserId, input.StartDate)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *app) setEndOfRenting(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	input := struct {
		Token   string `json:"token"`
		Id      int    `json:"id"`
		EndDate string `json:"start_date"` // DateOnly   = "2006-01-02"
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

	err = database.SetEndDate(app.DB, input.Id, input.EndDate)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *app) getReports(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	type fault struct {
		Id           int    `json:"id"`
		Description  string `json:"description"`
		DateReported string `json:"date_reported"`
		StatusId     int    `json:"status_id"`
		ApartamentId int    `json:"apartament_id"`
	}

	output := struct {
		Faults []fault `json:"faults"`
	}{}

	r.ParseMultipartForm(32 << 20)
	token := r.FormValue("token")

	userId, err := auth.ValidateSession(app.CACHE, token)
	if err != nil {
		sendError(w, Error{401, "Incorrect Token", "Unauthorized"}, err)
		return
	}

	role, err := app.Query.GetUserRole(app.Ctx, int64(userId))
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	apartId := -1

	if role == 2 {
		apartId, err = database.GetApartamentId(app.DB, userId)
		if err != nil {
			sendError(w, Error{400, "Database", "Internal Server Error"}, err)
			return
		}
	}

	id, description, dateReported, statusId, apartamentId, err := database.GetFaultReports(app.DB)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	var data []fault
	for n := range id {
		if apartamentId[n] == apartId {
			data = append(data,
				fault{Id: id[n],
					Description:  description[n],
					DateReported: dateReported[n],
					StatusId:     statusId[n],
					ApartamentId: apartamentId[n]})
		}
	}
	output.Faults = data

	if err = json.NewEncoder(w).Encode(&output); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (app *app) addFault(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	input := struct {
		Token        string `json:"token"`
		Description  string `json:"description"`
		DateReported string `json:"date_reported"` // DateOnly   = "2006-01-02"
		StatusId     int    `json:"status_id"`
		ApartamentId int    `json:"apartament_id"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	_, err = auth.ValidateSession(app.CACHE, input.Token)
	if err != nil {
		sendError(w, Error{401, "Incorrect Token", "Unauthorized"}, err)
		return
	}

	err = database.AddFault(app.DB, input.Description, input.DateReported, input.StatusId, input.ApartamentId)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////

func (app *app) addUser(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	user := struct {
		Token    string `json:"token"`
		Name     string `json:"name"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Role     int    `json:"role"`
		Id       int    `json:"id"`
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

	user.Id, err = database.GetId(app.DB, user.Email)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
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
		Token  string   `json:"token"`
		Email  string   `json:"email"`
		Spec   []string `json:"spec"`
		SpecId []int    `json:"spec_id"`
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

	data.SpecId, data.Spec, err = database.GetSubcontractorSpec(app.DB)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
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
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *app) getRepairs(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)
	r.ParseMultipartForm(32 << 20)
	token := r.FormValue("token")

	if erro := app.checkRole(token, 1); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	output, err := app.Query.GetRepair(app.Ctx)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	if err := json.NewEncoder(w).Encode(&output); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *app) addRepair(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	input := struct {
		Token  string               `json:"token"`
		Repair sqlc.AddRepairParams `json:"repair"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	if erro := app.checkRole(input.Token, 1); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	if err := app.Query.AddRepair(app.Ctx, input.Repair); err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *app) assignSubContractor(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	input := struct {
		Token      string                       `json:"token"`
		Contractor sqlc.UpdateSubToRepairParams `json:"contractor"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	output, err := app.Query.UpdateSubToRepair(app.Ctx, input.Contractor)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	if err := json.NewEncoder(w).Encode(&output); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *app) checkRole(token string, role_id ...int) *Error {
	userId, err := auth.ValidateSession(app.CACHE, token)
	if err != nil {
		log.Println(err)
		return &Error{401, "Incorrect Token", "Unauthorized"}
	}

	role, err := app.Query.GetUserRole(app.Ctx, int64(userId))
	if err != nil {
		log.Println(err)
		return &Error{401, "Database", "Internal Server Error"}
	}

	valid := false

	for i := range role_id {
		if role_id[i] == int(role) {
			valid = true
		}
	}

	if !valid {
		log.Println("Wrong role")
		return &Error{401, "Wrong role", "Unauthorized"}
	}
	return nil
}
