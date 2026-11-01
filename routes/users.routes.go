package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Max23strm/pitz-backend/db"
	"github.com/Max23strm/pitz-backend/helpers"
	"github.com/Max23strm/pitz-backend/models"
	"github.com/Max23strm/pitz-backend/queries"
	"github.com/Max23strm/pitz-backend/validations"
	"github.com/gorilla/mux"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}
	userSql := queries.UserSql

	allUsers := models.AllUserArr{}
	usersData, err := db.DB.Query(userSql)
	if err != nil {
		helpers.BadRequestResponse(w, "Error fetching users", err, "8894")
		return
	}
	for usersData.Next() {
		currentUser := models.AllUsers{}
		err := usersData.Scan(&currentUser.User_uid, &currentUser.Username, &currentUser.First_name, &currentUser.Last_name)
		if err != nil {
			helpers.BadRequestResponse(w, "Error scanning users", err, "8895")
			return
		}

		allUsers = append(allUsers, currentUser)
	}
	if usersData.Err() != nil {
		helpers.BadRequestResponse(w, "Error obtaining users", err, "8894")
		return
	}

	helpers.SuccessResponse(w, "ok", allUsers)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}
	userSql := "SELECT user_uid, email, username, first_name, last_name from users WHERE user_uid = $1"
	vars := mux.Vars(r)

	userBasicData := db.DB.QueryRow(userSql, vars["id"])
	currentUser := models.FullUser{}

	err := userBasicData.Scan(&currentUser.User_uid, &currentUser.Email, &currentUser.Username, &currentUser.First_name, &currentUser.Last_name)
	if err != nil {
		if err == sql.ErrNoRows {
			// No user found with that ID
			respuesta := map[string]interface{}{
				"isSuccess": false,
				"estado":    "No encontrado",
				"mensaje":   "Usuario no encontrado",
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(respuesta)
			return
		}

		// Other error (e.g., database or scan error)
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   "Error obteniendo usuario: " + err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)
		return
	}

	respuesta := map[string]interface{}{
		"isSuccess": true,
		"estado":    "Error",
		"data":      currentUser,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respuesta)

}
func GetBasicUserHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}

	userSql := "SELECT user_uid, email, username from users WHERE user_uid = $1"
	vars := mux.Vars(r)

	userBasicData := db.DB.QueryRow(userSql, vars["id"])
	currentUser := models.BasicUser{}

	err := userBasicData.Scan(&currentUser.User_uid, &currentUser.Email, &currentUser.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			// No user found with that ID
			respuesta := map[string]interface{}{
				"isSuccess": false,
				"estado":    "No encontrado",
				"mensaje":   "Usuario no encontrado",
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(respuesta)
			return
		}

		// Other error (e.g., database or scan error)
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   "Error obteniendo usuario: " + err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)
		return
	}
	respuesta := map[string]interface{}{
		"isSuccess": true,
		"estado":    "OK",
		"data":      currentUser,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respuesta)

}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}

	vars := mux.Vars(r)

	var NewPass models.PasswordUser
	player_uid := vars["id"]

	if err := json.NewDecoder(r.Body).Decode(&NewPass); err != nil {
		respuesta := map[string]string{
			"estado":  "Error",
			"mensaje": "Error al obtener datos",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)
		return
	}
	hashed, err := validations.HashPassword(NewPass.New_password)

	if err != nil {
		respuesta := map[string]string{
			"estado":  "Error",
			"mensaje": "Error al obtener datos",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)
		return
	}
	newPassSql := "UPDATE \"users\" SET \"hashed_password\" = $1 WHERE \"users\".\"user_uid\" = $2; "

	_, err = db.DB.Exec(newPassSql, string(hashed), player_uid)
	if err != nil {
		respuesta := map[string]string{
			"estado":  "Error",
			"mensaje": "Error al obtener datos",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)
		return
	}

	respuesta := map[string]interface{}{
		"isSuccess":  true,
		"estado":     "OK",
		"player_uid": player_uid,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respuesta)
}
