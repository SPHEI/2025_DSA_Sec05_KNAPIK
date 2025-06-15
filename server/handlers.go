package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

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
		ID    int64  `json:"id"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	result, err := app.Query.GetUserPasswordEmail(app.Ctx, input.Email)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	if result.Password != input.Password || result.Password == "" {
		sendError(w, Error{401, "Wrong password or login", "Unauthorized"}, err)
		return
	}

	output.Token, err = auth.CreateSession(app.CACHE, int(result.ID))
	if err != nil {
		sendError(w, Error{500, "Could not generate a new token", "Internal Server Error"}, err)
		return
	}

	output.Role, err = app.Query.GetUserRole(app.Ctx, result.ID)
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
		ApartamentId int64   `json:"apartament_id"`
		Rent         float64 `json:"rent"`
	}{}

	r.ParseMultipartForm(32 << 20)
	token := r.FormValue("token")

	if _, erro := app.checkRole(token, 2); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	userId, _ := auth.ValidateSession(app.CACHE, token)

	log.Println(sql.NullInt64{Int64: int64(userId)})

	apartmentId, err := app.Query.GetApartmentID(app.Ctx, userId)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	output.Rent, err = app.Query.GetRent(app.Ctx, apartmentId)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	output.ApartamentId = apartmentId

	err = json.NewEncoder(w).Encode(&output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *app) subInfo(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	r.ParseMultipartForm(32 << 20)
	token := r.FormValue("token")

	if _, erro := app.checkRole(token, 3); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	userID, _ := auth.ValidateSession(app.CACHE, token)

	output, err := app.Query.GetSubconInfo(app.Ctx, int64(userID))
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

	r.ParseMultipartForm(32 << 20)
	token := r.FormValue("token")

	if _, erro := app.checkRole(token, 1); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	output, err := app.Query.GetSubcontractors(app.Ctx)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	if err := json.NewEncoder(w).Encode(&output); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (app *app) addContractor(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	input := struct {
		Token         string                      `json:"token"`
		Subcontractor sqlc.AddSubcontractorParams `json:"subcontractor"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	if _, erro := app.checkRole(input.Token, 1); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	err = app.Query.AddSubcontractor(app.Ctx, input.Subcontractor)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *app) getTenants(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	r.ParseMultipartForm(32 << 20)
	token := r.FormValue("token")

	if _, erro := app.checkRole(token, 1); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	output, err := app.Query.GetTenets(app.Ctx)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	if err = json.NewEncoder(w).Encode(&output); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (app *app) getApartaments(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	r.ParseMultipartForm(32 << 20)
	token := r.FormValue("token")

	if _, erro := app.checkRole(token, 1); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	output, err := app.Query.GetApartments(app.Ctx)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	if err = json.NewEncoder(w).Encode(&output); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (app *app) addApartament(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	input := struct {
		Token     string                  `json:"token"`
		Apartment sqlc.AddApartmentParams `json:"aparment"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	if _, erro := app.checkRole(input.Token, 1); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	err = app.Query.AddApartment(app.Ctx, input.Apartment)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *app) getOwners(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	r.ParseMultipartForm(32 << 20)
	token := r.FormValue("token")
	if _, erro := app.checkRole(token, 1); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	output, err := app.Query.GetOwners(app.Ctx)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	if err = json.NewEncoder(w).Encode(&output); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (app *app) addOwner(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	input := struct {
		Token string              `json:"token"`
		Owner sqlc.AddOwnerParams `json:"owner"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	if _, erro := app.checkRole(input.Token, 1); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	err = app.Query.AddOwner(app.Ctx, input.Owner)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *app) changeRent(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)
	
	data := struct {
		Token string                 `json:"token"`
		Rent  sqlc.ChangeRent2Params `json:"rent"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	if _, erro := app.checkRole(data.Token, 1); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	if err := app.Query.ChangeRent1(app.Ctx, data.Rent.ApartmentID); err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	if err := app.Query.ChangeRent2(app.Ctx, data.Rent); err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *app) getCurrentRenting(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	r.ParseMultipartForm(32 << 20)
	token := r.FormValue("token")
	if _, erro := app.checkRole(token, 1); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	output, err := app.Query.GetActiveRenting(app.Ctx)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	if err = json.NewEncoder(w).Encode(&output); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (app *app) addNewRenting(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	input := struct {
		Token   string                   `json:"token"`
		Renting sqlc.AddNewRentingParams `josn:"renting"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	if _, erro := app.checkRole(input.Token, 1); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	err = app.Query.AddNewRenting(app.Ctx, input.Renting)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *app) setEndOfRenting(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	input := struct {
		Token string                `json:"token"`
		End   sqlc.SetEndDateParams `json:"end"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	if _, erro := app.checkRole(input.Token, 1); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	err = app.Query.SetEndDate(app.Ctx, input.End)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *app) getReports(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

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

	if role == 2 {
		apart, err := app.Query.GetApartmentID(app.Ctx, userId)
		if err != nil {
			sendError(w, Error{400, "Database", "Internal Server Error"}, err)
			return
		}
		output, err := app.Query.GetFaultReportsUser(app.Ctx, apart)
		if err != nil {
			sendError(w, Error{400, "Database", "Internal Server Error"}, err)
			return
		}

		if err = json.NewEncoder(w).Encode(&output); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		output, err := app.Query.GetFaultReports(app.Ctx)
		if err != nil {
			sendError(w, Error{400, "Database", "Internal Server Error"}, err)
			return
		}

		if err = json.NewEncoder(w).Encode(&output); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (app *app) addFault(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	input := struct {
		Token string              `json:"token"`
		Fault sqlc.AddFaultParams `json:"fault"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	log.Println(input.Fault)

	_, err = auth.ValidateSession(app.CACHE, input.Token)
	if err != nil {
		sendError(w, Error{401, "Incorrect Token", "Unauthorized"}, err)
		return
	}

	err = app.Query.AddFault(app.Ctx, input.Fault)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *app) updateFault(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	input := struct {
		Token string                       `json:"token"`
		Fault sqlc.UpdateFaultStatusParams `json:"fault"`
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

	output, err := app.Query.UpdateFaultStatus(app.Ctx, input.Fault)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	if err = json.NewEncoder(w).Encode(&output); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////

func (app *app) addUser(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	input := struct {
		Token string             `json:"token"`
		User  sqlc.AddUserParams `json:"user"`
	}{}

	output := struct {
		Id int `json:"id"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	if _, erro := app.checkRole(input.Token, 1); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	_, err = bcrypt.GenerateFromPassword([]byte(input.User.Password), 12)
	if err != nil {
		sendError(w, Error{500, "Could not generate hash from password", "Internal Server Error"}, err)
		return
	}

	err = app.Query.AddUser(app.Ctx, input.User)
	if err != nil {
		sendError(w, Error{500, "Could not add user", "Internal Server Error"}, err)
		return
	}

	output.Id, err = database.GetId(app.DB, input.User.Email)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	if err := json.NewEncoder(w).Encode(output); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (app *app) getSubContractorSpec(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	data := struct {
		Token string `json:"token"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	if _, erro := app.checkRole(data.Token, 1); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	output, err := app.Query.GetSubcontractorSpec(app.Ctx)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	if err := json.NewEncoder(w).Encode(output); err != nil {
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

	if _, erro := app.checkRole(data.Token, 1); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	err = app.Query.AddSpec(app.Ctx, data.Name)
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

	role, erro := app.checkRole(token, 1, 2, 3)
	if erro != nil {
		sendError(w, *erro, nil)
		return
	}

	if role == 1 {
		output, err := app.Query.GetRepair(app.Ctx)
		if err != nil {
			sendError(w, Error{400, "Database", "Internal Server Error"}, err)
			return
		}

		if err := json.NewEncoder(w).Encode(&output); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if role == 3 {
		userID, _ := auth.ValidateSession(app.CACHE, token)
		output, err := app.Query.GetRepairSub(app.Ctx, int64(userID))
		if err != nil {
			sendError(w, Error{400, "Database", "Internal Server Error"}, err)
			return
		}

		if err := json.NewEncoder(w).Encode(&output); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		userID, _ := auth.ValidateSession(app.CACHE, token)
		apartId, err := app.Query.GetApartmentID(app.Ctx, userID)
		log.Printf("Apart ID: %d", apartId)
		if err != nil {
			sendError(w, Error{400, "Database", "Internal Server Error"}, err)
			return
		}

		output, err := app.Query.GetRepairApart(app.Ctx, apartId)
		if err != nil {
			sendError(w, Error{400, "Database", "Internal Server Error"}, err)
			return
		}

		if err := json.NewEncoder(w).Encode(&output); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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

	if _, erro := app.checkRole(input.Token, 1); erro != nil {
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

	if _, erro := app.checkRole(input.Token, 1); erro != nil {
		sendError(w, *erro, nil)
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

func (app *app) updateRepairData(w http.ResponseWriter, r *http.Request) {

	input := struct {
		Token  string                      `json:"token"`
		Repair sqlc.UpdateRepairDataParams `json:"repair"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	if _, erro := app.checkRole(input.Token, 1, 3); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	output, err := app.Query.UpdateRepairData(app.Ctx, input.Repair)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	if err := json.NewEncoder(w).Encode(&output); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *app) checkRole(token string, role_id ...int64) (int64, *Error) {
	userId, err := auth.ValidateSession(app.CACHE, token)
	if err != nil {
		log.Println(err)
		return -1, &Error{401, "Incorrect Token", "Unauthorized"}
	}

	role, err := app.Query.GetUserRole(app.Ctx, userId)
	if err != nil {
		log.Println(err)
		return -1, &Error{401, "Database", "Internal Server Error"}
	}

	valid := false

	for i := range role_id {
		if role_id[i] == role {
			valid = true
		}
	}

	if !valid {
		log.Println("Wrong role")
		return role, &Error{401, "Wrong role", "Unauthorized"}
	}
	return role, nil
}
