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
	// to be expanded
}

// interface for json needed
func sendError(w http.ResponseWriter, error Error) {
	if err := json.NewEncoder(w).Encode(error); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
		log.Println(err)
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"})
		return
	}

	if erro := app.checkAdmin(user.Token); erro != nil {
		sendError(w, *erro)
		return
	}

	_, err = bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		log.Println(err)
		sendError(w, Error{500, "Could not generate hash from password", "Internal Server Error"})
		return
	}

	err = database.AddUser(app.DB, []string{user.Name, user.Password, user.Email, user.Phone}, user.Role)
	if err != nil {
		log.Println(err)
		sendError(w, Error{500, "Could not add user", "Internal Server Error"})
		return
	}

	user.Password = strings.Repeat("*", len(user.Password)) // should be changed
	if err := json.NewEncoder(w).Encode(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (app *app) login(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	user := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Token    string `json:"token"`
		Role     int    `json:"role"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"})
		return
	}

	hashedPassword, err := database.GetUser(app.DB, user.Email)
	if err != nil {
		log.Println(err)
		sendError(w, Error{400, "Database", "Internal Server Error"})
		return
	}

	user.Role, err = database.GetRole(app.DB, user.Email)
	if err != nil {
		log.Println(err)
		sendError(w, Error{400, "Database", "Internal Server Error"})
		return
	}

	//if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password)); err != nil {
	if hashedPassword != user.Password || hashedPassword == "" {
		log.Println(err)
		sendError(w, Error{401, "Wrong password or login", "Unauthorized"})
		return
	} else {
		token, err := auth.CreateSession(app.CACHE, user.Email)
		if err != nil {
			log.Println(err)
			sendError(w, Error{500, "Could not generate a new token", "Internal Server Error"})
			return
		}

		log.Printf("User: %s - Logged in with token: %s", user.Email, token)

		user.Token = token
		user.Password = strings.Repeat("*", len(user.Password)) // should be changed
		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// func (app *app) issue(w http.ResponseWriter, r *http.Request) {
// 	prepareResponse(w)
//
// 	user := struct {
// 		Token       string `json:"token"`
// 		Description string `json:"description"`
// 	}{}
//
// 	err := json.NewDecoder(r.Body).Decode(&user)
// 	if err != nil {
// 		log.Println(err)
// 		sendError(w, Error{400, "Could not acquire json data", "Bad Request"})
// 		return
// 	}
//
// 	_, err = auth.ValidateSession(app.CACHE, user.Token)
// 	if err != nil {
// 		log.Println(err)
// 		sendError(w, Error{401, "Incorrect Token", "Unauthorized"})
// 		return
// 	}
//
// 	if err = database.AddIssue(app.DB, user.Description); err != nil {
// 		log.Println(err)
// 		sendError(w, Error{500, "Could not add issue", "Internal Server Error"})
// 		return
// 	}
//
// 	w.WriteHeader(http.StatusOK)
// }

func (app *app) info(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	r.ParseMultipartForm(32 << 20)
	token := r.FormValue("token")

	email, err := auth.ValidateSession(app.CACHE, token)
	if err != nil {
		log.Println(err)
		sendError(w, Error{401, "Incorrect Token", "Unauthorized"})
		return
	}

	name, phone, role_id, err := database.GetInfo(app.DB, email)
	info := struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
		Role  int    `json:"role"`
	}{
		Name:  name,
		Email: email,
		Phone: phone,
		Role:  role_id,
	}

	if err := json.NewEncoder(w).Encode(info); err != nil {
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
		log.Println(err)
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"})
		return
	}

	database.DeleteToken(app.CACHE, token.Token)

	token.Token = "" // should be changed
	if err := json.NewEncoder(w).Encode(token); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (app *app) getApartamentList(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	tenants := struct {
		Token       string   `json:"token"`
		Apartaments []string `json:"apartaments"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&tenants)
	if err != nil {
		log.Println(err)
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"})
		return
	}

	if erro := app.checkAdmin(tenants.Token); erro != nil {
		sendError(w, *erro)
		return
	}

	tenants.Apartaments, err = database.GetApartaments(app.DB)
	if err != nil {
		log.Println(err)
		sendError(w, Error{400, "Database", "Internal Server Error"})
	}

	if err := json.NewEncoder(w).Encode(tenants); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *app) addApartament(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	apartament := struct {
		Token          string `json:"token"`
		OwnerEmail     string `json:"owner_email"`
		Name           string `json:"name"`
		Street         string `json:"street"`
		BuildingNumber string `json:"building_number"`
		BuildingName   string `json:"building_name"`
		FlatNumber     string `json:"flat_number"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&apartament)
	if err != nil {
		log.Println(err)
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"})
		return
	}

	if erro := app.checkAdmin(apartament.Token); erro != nil {
		sendError(w, *erro)
		return
	}

	err = database.AddApartament(app.DB, []string{apartament.OwnerEmail, apartament.Name, apartament.Street, apartament.BuildingNumber, apartament.BuildingName, apartament.FlatNumber})
	if err != nil {
		log.Println(err)
		sendError(w, Error{400, "Database", "Internal Server Error"})
	}

	w.WriteHeader(http.StatusOK)
}

//	func (app *app) addApartament(w http.ResponseWriter, r *http.Request) {
//		prepareResponse(w)
//
//		user := struct {
//			Token          string `json:"token"`
//			Name           string `json:"name"`
//			Street         string `json:"street"`
//			BuildingNumber string `json:"building_number"`
//			BuildingName   string `json:"building_name"`
//			FlatNumber     string `json:"flat_number"`
//			OwnerName      string `json:"owner_name"`
//		}{}
//
//		err := json.NewDecoder(r.Body).Decode(&user)
//		if err != nil {
//			log.Println(err)
//			sendError(w, Error{400, "Could not acquire json data", "Bad Request"})
//			return
//		}
//
//		email, err := auth.ValidateSession(app.CACHE, user.Token)
//		if err != nil {
//			log.Println(err)
//			sendError(w, Error{401, "Incorrect Token", "Unauthorized"})
//			return
//		}
//
//		role, err := database.GetRole(app.DB, email)
//		if err != nil {
//			log.Println(err)
//			sendError(w, Error{401, "Database", "Internal Server Error"})
//			return
//		}
//		log.Println(role)
//		if role != 1 {
//			log.Println("Wrong role")
//			sendError(w, Error{401, "Wrong role", "Unauthorized"})
//			return
//		}
//
//
//		err = database.AddApartament(app.DB, []string{user.Name, string(hashedPassword), user.Email, user.Phone}, user.Role)
//		if err != nil {
//			log.Println(err)
//			sendError(w, Error{500, "Could not add user", "Internal Server Error"})
//			return
//		}
//
//		user.Password = strings.Repeat("*", len(user.Password)) // should be changed
//		if err := json.NewEncoder(w).Encode(user); err != nil {
//			w.WriteHeader(http.StatusInternalServerError)
//			return
//		}
//	}
func (app *app) getEmailList(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	tenants := struct {
		Token  string   `json:"token"`
		Emails []string `json:"emails"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&tenants)
	if err != nil {
		log.Println(err)
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"})
		return
	}

	if erro := app.checkAdmin(tenants.Token); erro != nil {
		sendError(w, *erro)
		return
	}

	tenants.Emails, err = database.GetEmails(app.DB)
	if err != nil {
		log.Println(err)
		sendError(w, Error{400, "Database", "Internal Server Error"})
	}

	if err := json.NewEncoder(w).Encode(tenants); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *app) getTenantList(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	tenants := struct {
		Token string   `json:"token"`
		Names []string `json:"names"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&tenants)
	if err != nil {
		log.Println(err)
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"})
		return
	}

	if erro := app.checkAdmin(tenants.Token); erro != nil {
		sendError(w, *erro)
		return
	}

	tenants.Names, err = database.GetTenent(app.DB)
	if err != nil {
		log.Println(err)
		sendError(w, Error{400, "Database", "Internal Server Error"})
	}

	if err := json.NewEncoder(w).Encode(tenants); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
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
		log.Println(err)
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"})
		return
	}

	if erro := app.checkAdmin(data.Token); erro != nil {
		sendError(w, *erro)
		return
	}

	data.Spec, err = database.GetSubcontractorSpec(app.DB)
	if err != nil {
		log.Println(err)
		sendError(w, Error{400, "Database", "Internal Server Error"})
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
		log.Println(err)
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"})
		return
	}

	if erro := app.checkAdmin(data.Token); erro != nil {
		sendError(w, *erro)
		return
	}

	err = database.AddSpec(app.DB, data.Name)
	if err != nil {
		log.Println(err)
		sendError(w, Error{400, "Database", "Internal Server Error"})
	}

	w.WriteHeader(http.StatusOK)
}

func (app *app) checkAdmin(token string) *Error {
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
	if role != 1 {
		log.Println("Wrong role")
		return &Error{401, "Wrong role", "Unauthorized"}
	}
	return nil
}
