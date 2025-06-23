package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

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

	role, err := app.Query.GetUserRole(app.Ctx, result.ID)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}
	switch role {
	case "admin":
		output.Role = 1
	case "tenant":
		output.Role = 2
	case "subcontractor":
		output.Role = 3
	}

	log.Printf("Login -- User: %s - Token: %s - Role: %d", input.Email, output.Token, output.Role)

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
		Apartment sqlc.GetApartmentAllRow `json:"apartment"`
		Rent      float64                 `json:"rent"`
		Status    string                  `json:"status"`
	}{}

	r.ParseMultipartForm(32 << 20)
	token := r.FormValue("token")

	role, erro := app.checkRole(token, "tenant", "admin")
	if erro != nil {
		sendError(w, *erro, nil)
		return
	}

	userId, _ := auth.ValidateSession(app.CACHE, token)

	if role == "admin" {
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			sendError(w, Error{400, "Wrong id input", "Bad Request"}, err)
			return
		}
		userId = int64(id)
	}

	output.Apartment, err = app.Query.GetApartmentAll(app.Ctx, userId)
	if err != nil && err != sql.ErrNoRows {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	output.Rent, err = app.Query.GetRent(app.Ctx, output.Apartment.ID.Int64)
	if err != nil && err != sql.ErrNoRows {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	_, err = app.Query.GetPendingPaymants(app.Ctx)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Println("GetOverduePayments")
			sendError(w, Error{400, "Database", "Internal Server Error"}, err)
			return
		}
	} else {
		output.Status = "Pending"
	}

	_, err = app.Query.GetOverduePayments(app.Ctx)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Println("GetOverduePayments")
			sendError(w, Error{400, "Database", "Internal Server Error"}, err)
			return
		}
	} else {
		output.Status = "Overdue"
	}

	if output.Status != "Overdue" && output.Status != "Pending" {
		output.Status = "Paid"
	}

	err = json.NewEncoder(w).Encode(&output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *app) subInfo(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	r.ParseMultipartForm(32 << 20)
	token := r.FormValue("token")

	if _, erro := app.checkRole(token, "subcontractor"); erro != nil {
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

	if _, erro := app.checkRole(token, "admin"); erro != nil {
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

	if _, erro := app.checkRole(input.Token, "admin"); erro != nil {
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

	if _, erro := app.checkRole(token, "admin"); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	output, err := app.Query.GetTenetsWithRent(app.Ctx)
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

	if _, erro := app.checkRole(token, "admin"); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	output, err := app.Query.GetApartmentsAndRent(app.Ctx)
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

	if _, erro := app.checkRole(input.Token, "admin"); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	output, err := app.Query.AddApartment(app.Ctx, input.Apartment)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	if err = json.NewEncoder(w).Encode(&output); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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

	if _, erro := app.checkRole(data.Token, "admin"); erro != nil {
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
	if _, erro := app.checkRole(token, "admin"); erro != nil {
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

	if _, erro := app.checkRole(input.Token, "admin"); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	rent, err := app.Query.GetActiveRentingID(app.Ctx, input.Renting.ApartmentID)
	if err != nil && err != sql.ErrNoRows {
		log.Println("GetActiveRentingID:")
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	if err := app.Query.MakeAsEnd(app.Ctx, rent.ApartmentID); err != nil {
		log.Println("MakeAsEnd:")
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
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

	if _, erro := app.checkRole(input.Token, "admin"); erro != nil {
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

func (app *app) setStatusOfRenting(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	input := struct {
		Token     string `json:"token"`
		RentingID int64  `json:"renting_id"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	if _, erro := app.checkRole(input.Token, "admin"); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	err = app.Query.MakeAsEnd(app.Ctx, input.RentingID)
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

	userID, err := auth.ValidateSession(app.CACHE, token)
	if err != nil {
		sendError(w, Error{401, "Incorrect Token", "Unauthorized"}, err)
		return
	}

	role, err := app.Query.GetUserRole(app.Ctx, int64(userID))
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	if string(role) == "tenant" {
		apart, err := app.Query.GetApartmentID(app.Ctx, userID)
		if err != nil {
			sendError(w, Error{400, "Database", "Internal Server Error"}, err)
			return
		}
		log.Println(apart)
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

	input.Fault.UserID, err = auth.ValidateSession(app.CACHE, input.Token)
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

	if _, erro := app.checkRole(input.Token, "admin"); erro != nil {
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

	if _, erro := app.checkRole(data.Token, "admin"); erro != nil {
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

	if _, erro := app.checkRole(data.Token, "admin"); erro != nil {
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

	role, erro := app.checkRole(token, "admin", "tenant", "subcontractor")
	if erro != nil {
		sendError(w, *erro, nil)
		return
	}

	if role == "admin" {
		output, err := app.Query.GetRepair(app.Ctx)
		if err != nil {
			sendError(w, Error{400, "Database", "Internal Server Error"}, err)
			return
		}

		if err := json.NewEncoder(w).Encode(&output); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if role == "subcontractor" {
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

	if _, erro := app.checkRole(input.Token, "admin"); erro != nil {
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

	if _, erro := app.checkRole(input.Token, "admin"); erro != nil {
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
	prepareResponse(w)

	input := struct {
		Token  string                      `json:"token"`
		Repair sqlc.UpdateRepairDataParams `json:"repair"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	if _, erro := app.checkRole(input.Token, "admin", "subcontractor"); erro != nil {
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

func (app *app) checkRole(token string, role_in ...string) (string, *Error) {
	userId, err := auth.ValidateSession(app.CACHE, token)
	if err != nil {
		log.Println(err)
		return "", &Error{401, "Incorrect Token", "Unauthorized"}
	}

	role, err := app.Query.GetUserRole(app.Ctx, userId)
	if err != nil {
		log.Println(err)
		return "", &Error{401, "Database", "Internal Server Error"}
	}

	valid := false

	for i := range role_in {
		if role_in[i] == role {
			valid = true
		}
	}

	if !valid {
		log.Println("Wrong role")
		return role, &Error{401, "Wrong role", "Unauthorized"}
	}
	return role, nil
}

func (app *app) getPayments(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	r.ParseMultipartForm(32 << 20)
	token := r.FormValue("token")

	role, erro := app.checkRole(token, "tenant", "admin")
	if erro != nil {
		sendError(w, *erro, nil)
		return
	}

	app.updateAllPayments(w, r)
	app.updateOverdueRent(w, r)

	if role == "admin" {
		output, err := app.Query.GetAllPayment(app.Ctx)
		if err != nil {
			sendError(w, Error{400, "Database", "Internal Server Error"}, err)
			return
		}

		if err = json.NewEncoder(w).Encode(&output); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {

		id, _ := auth.ValidateSession(app.CACHE, token)

		output, err := app.Query.GetPaymentsId(app.Ctx, id)
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

//func (app *app) updatePayment(w http.ResponseWriter, r *http.Request) {
//	id := r.Context().Value("id").(int64)
//	apartmentID, err := app.Query.GetApartmentID(app.Ctx, id)
//	if err != nil {
//		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
//		return
//	}
//
//	output, err := app.Query.GetActiveRenting(app.Ctx)
//	if err != nil {
//		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
//		return
//	}
//
//	for _, rent := range output {
//		if rent.ApartmentID == apartmentID && rent.UserID == id {
//			startDate := rent.StartDate
//			payments, err := app.Query.GetPayment(app.Ctx, id)
//			if err != nil {
//				sendError(w, Error{400, "Database", "Internal Server Error"}, err)
//				return
//			}
//
//			price, err := app.Query.GetRent(app.Ctx, rent.ApartmentID)
//			if err != nil {
//				sendError(w, Error{400, "Database", "Internal Server Error"}, err)
//				return
//			}
//
//			amount := 0
//			for x := range payments {
//				if startDate.Before(payments[x].DueDate) {
//					amount++
//				}
//			}
//
//			startDate = time.Date(
//				startDate.Year(),
//				startDate.Month(),
//				1, // Set the day to 1
//				0, // Set hour to 0
//				0, // Set minute to 0
//				0, // Set second to 0
//				0, // Set nanosecond to 0
//				startDate.Location(),
//			)
//
//			endDate := time.Date(
//				time.Now().Year(),
//				time.Now().Month(),
//				1, // Set the day to 1
//				0, // Set hour to 0
//				0, // Set minute to 0
//				0, // Set second to 0
//				0, // Set nanosecond to 0
//				time.Now().Location(),
//			)
//			endDate = endDate.AddDate(0, 1, 0)
//
//			for startDate.Before(endDate) {
//				if amount <= 0 {
//					app.Query.AddPayment(app.Ctx, sqlc.AddPaymentParams{UserID: id, UserID_2: id, Amount: price, DueDate: startDate.AddDate(0, 1, 0)})
//				}
//				startDate = startDate.AddDate(0, 1, 0)
//				amount--
//			}
//		}
//	}
//}

func (app *app) pay(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	input := struct {
		Token   string                   `json:"token"`
		Payment sqlc.UpdatePaymentParams `json:"payment"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	if _, erro := app.checkRole(input.Token, "tenant", "subcontractor"); erro != nil {
		sendError(w, *erro, nil)
		return
	}

	output, err := app.Query.UpdatePayment(app.Ctx, input.Payment)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	if err := json.NewEncoder(w).Encode(&output); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *app) updateAllPayments(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)
	activeRentings, err := app.Query.GetActiveRenting(app.Ctx)
	if err != nil {
		log.Println("GetActiveRenting:")
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	for _, activeRenting := range activeRentings {
		price, err := app.Query.GetRent(app.Ctx, activeRenting.ApartmentID)
		if err != nil {
			log.Println("GetRent:")
			sendError(w, Error{400, "Database", "Internal Server Error"}, err)
			return
		}

		payments, err := app.Query.GetPayments(app.Ctx, activeRenting.ID)
		if err != nil {
			log.Println("GetPayments:")
			sendError(w, Error{400, "Database", "Internal Server Error"}, err)
			return
		}

		amount := 0
		startDate := time.Date(
			activeRenting.StartDate.Year(),
			activeRenting.StartDate.Month(),
			2, // Set the day to 1
			0, // Set hour to 0
			0, // Set minute to 0
			0, // Set second to 0
			0, // Set nanosecond to 0
			activeRenting.StartDate.Location(),
		)
		loopStartDate := startDate
		startDate = startDate.AddDate(0, 0, -1)

		if activeRenting.EndDate.Valid {
			endDate := activeRenting.EndDate.Time
			for loopStartDate.Before(endDate) {
				amount++
				loopStartDate = loopStartDate.AddDate(0, 1, 0)
			}

			if amount > len(payments) {
				jump := amount - len(payments)
				startDate = startDate.AddDate(0, len(payments)+1, 0)
				for range jump + 1 {
					err := app.Query.AddPayment(app.Ctx, sqlc.AddPaymentParams{Amount: price, DueDate: startDate, RentingID: activeRenting.ID})
					if err != nil {
						log.Println("AddPayment:")
						sendError(w, Error{400, "Database", "Internal Server Error"}, err)
						return
					}

					startDate = startDate.AddDate(0, 1, 0)
				}
			}
		} else {
			endDate := time.Date(
				time.Now().Year(),
				time.Now().Month(),
				1, // Set the day to 1
				0, // Set hour to 0
				0, // Set minute to 0
				0, // Set second to 0
				0, // Set nanosecond to 0
				time.Now().Location(),
			)

			for loopStartDate.Before(endDate) {
				amount++
				loopStartDate = loopStartDate.AddDate(0, 1, 0)
			}

			if amount > len(payments) {
				jump := amount - len(payments)
				startDate = startDate.AddDate(0, len(payments)+1, 0)
				for range jump + 1 {
					if err := app.Query.AddPayment(app.Ctx, sqlc.AddPaymentParams{Amount: price, DueDate: startDate, RentingID: activeRenting.ID}); err != nil {
						log.Println("AddPayment:")
						sendError(w, Error{400, "Database", "Internal Server Error"}, err)
						return
					}
					startDate = startDate.AddDate(0, 1, 0)
				}
			}
		}
	}
}

func (app *app) updateOverdueRent(w http.ResponseWriter, r *http.Request) {
	pendingPaymants, err := app.Query.GetPendingPaymants(app.Ctx)
	if err != nil {
		log.Println("GetPendingPaymants:")
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	currentDate := time.Now()
	for _, payment := range pendingPaymants {
		if !payment.DueDate.After(currentDate) {
			if _, err := app.Query.SetPaymanyOverdue(app.Ctx, payment.ID); err != nil {
				log.Println("SetPaymentOverdue:")
				sendError(w, Error{400, "Database", "Internal Server Error"}, err)
				return
			}
		}
	}
}
