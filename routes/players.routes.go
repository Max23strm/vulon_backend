package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Max23strm/pitz-backend/db"
	"github.com/Max23strm/pitz-backend/helpers"
	"github.com/Max23strm/pitz-backend/models"
	"github.com/Max23strm/pitz-backend/validations"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetPlayersHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}
	vars := mux.Vars(r)
	entity_uid, err := uuid.Parse(vars["entity_uid"])

	if len(entity_uid) == 0 || err != nil {
		helpers.BadRequestResponse(w, "Entity is required", errors.New("No entity"), "ENTITY_REQUIRED")
		return
	}

	playersSql := "SELECT players.player_uid, first_name, last_name, status, email FROM players WHERE entity_uid = $1 ORDER BY players.first_name ASC"
	players := models.Players{}

	datos, err := db.DB.Query(playersSql, entity_uid)
	if err != nil {
		fmt.Println(err)
		helpers.InternalServerErrorResponse(w, "Petition error")
		return
	}
	defer datos.Close()

	for datos.Next() {
		dato := models.Player{}
		err := datos.Scan(&dato.Player_uid, &dato.FirstName, &dato.LastName, &dato.Status, &dato.Email)
		if err != nil {
			helpers.ErrorResponse(w, http.StatusBadRequest, "Error scanning players", err)
			return
		}

		players = append(players, dato)
	}

	if err := datos.Err(); err != nil {
		helpers.InternalServerErrorResponse(w, "Error reading players")
		return
	}

	helpers.SuccessResponse(w, "", players)
}

func GetPlayerByIdHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}
	playerSql := "SELECT players.player_uid, players.first_name, players.last_name, players.email, players.status, players.address, players.birth_dt, players.comments, players.blood_type, players.afiliation, players.sex, players.curp, players.enfermedad, players.phone_number, players.emergency_phone, players.insurance, players.insurance_name FROM players WHERE players.player_uid = $1"

	vars := mux.Vars(r)

	parsedUID, err := uuid.Parse(vars["id"])
	if err != nil {
		helpers.BadRequestResponse(
			w,
			"Player identifier is not valid",
			err,
			"NO_VALID_ID",
		)
		return
	}

	playerRow := db.DB.QueryRow(playerSql, parsedUID)

	player := models.PlayerDetails{}

	err = playerRow.Scan(&player.Player_uid, &player.FirstName, &player.LastName, &player.Email, &player.Status, &player.Address, &player.Birth_dt, &player.Comments, &player.BloodType, &player.Afiliation, &player.Sex, &player.Curp, &player.Enfermedad, &player.Phone_number, &player.Emergency_number, &player.Insurance, &player.Insurance_name)
	if err != nil {
		helpers.NotFoundResponse(
			w,
			"Player not found",
		)
		return
	}

	helpers.SuccessResponse(
		w,
		"Fetched successfully",
		player,
	)
}

func PostPlayerHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}

	var player models.PostPlayerDetails
	new_uuid := uuid.New()

	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		helpers.BadRequestResponse(w, "No information sent", err, "EMPTY_DATA")
		return
	}

	//Validacion de campos requeridos
	validationErrors := validations.PlayersPostValidations(player)

	if len(validationErrors) > 0 {
		message := strings.Join(validationErrors, ", ")

		helpers.BadRequestResponse(
			w,
			message,
			errors.New(message),
			"VALIDATION_ERRORS",
		)

		return
	}

	sqlString := "INSERT INTO \"players\" (\"player_uid\", \"first_name\", \"last_name\", \"phone_number\", \"emergency_phone\", \"email\", \"status\", \"positions\", \"birth_dt\", \"blood_type\", \"comments\", \"credential\", \"address\", \"afiliation\", \"sex\", \"curp\", \"enfermedad\", \"insurance\", \"insurance_name\") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)"

	_, err := db.DB.Exec(sqlString, new_uuid.String(), player.FirstName, player.LastName, player.Phone_number, player.Emergency_number, player.Email, player.Status, nil, player.Birth_dt, player.BloodType, player.Comments, player.Credential, player.Address, player.Afiliation, player.Sex, player.Curp, player.Enfermedad, player.Insurance, player.Insurance_name)

	if err != nil {
		helpers.BadRequestResponse(
			w,
			"Unable to store information",
			err,
			"STORING_ERROR",
		)
		return
	}

	helpers.CreatedResponse(
		w,
		"Stored succesfully",
		new_uuid.String(),
	)
}

func EditPlayerHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}
	vars := mux.Vars(r)

	var event models.PutPlayerDetails
	player_uid := vars["id"]

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		helpers.BadRequestResponse(w, "No information sent", err, "EMPTY_DATA")
		return
	}

	fields := []string{}
	values := []interface{}{}

	idx := 1
	addField := func(col string, val interface{}) {
		fields = append(fields, fmt.Sprintf("%s = $%d", col, idx))
		values = append(values, val)
		idx++
	}

	if event.FirstName != nil {
		addField("first_name", *event.FirstName)
	}
	if event.LastName != nil {
		addField("last_name", *event.LastName)
	}
	if event.Email != nil {
		addField("email", *event.Email)
	}
	if event.Status != nil {
		addField("status", *event.Status)
	}
	if event.Birth_dt != nil {
		addField("birth_dt", *event.Birth_dt)
	}
	if event.Address != nil {
		addField("address", *event.Address)
	}
	if event.Sex != nil {
		addField("sex", *event.Sex)
	}
	if event.BloodType != nil {
		addField("blood_type", *event.BloodType)
	}
	if event.Comments != nil {
		addField("comments", *event.Comments)
	}
	if event.Credential != nil {
		addField("credential", *event.Credential)
	}
	if event.Afiliation != nil {
		addField("afiliation", *event.Afiliation)
	}
	if event.Curp != nil {
		addField("curp", *event.Curp)
	}
	if event.Enfermedad != nil {
		addField("enfermedad", *event.Enfermedad)
	}
	if event.Phone_number != nil {
		addField("phone_number", *event.Phone_number)
	}
	if event.Emergency_number != nil {
		addField("emergency_phone", *event.Emergency_number)
	}
	if event.Insurance != nil {
		addField("insurance", *event.Insurance)
	}
	if event.Insurance_name != nil {
		addField("insurance_name", *event.Insurance_name)
	}

	if len(fields) == 0 {
		helpers.BadRequestResponse(w, "No fields to update", errors.New(""), "EMPTY_DATA")
		return
	}

	values = append(values, player_uid)

	query := fmt.Sprintf("UPDATE players SET %s WHERE player_uid = $%d", strings.Join(fields, ", "), idx)

	result, err := db.DB.Exec(query, values...)

	if err != nil {
		helpers.BadRequestResponse(w, "Error updating", err, "ERROR_UPDATING")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		helpers.BadRequestResponse(w, "Error updating", err, "NOT_FOUND")
		return
	}

	helpers.SuccessResponse(w, "Edited succesfully", vars["id"])
}

func DeletePlayerHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete user"))
}
