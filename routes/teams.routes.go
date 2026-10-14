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
	"github.com/lib/pq"
)

func GetTeamsByEntityHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}
	entity_uid := r.URL.Query().Get("entity_uid")

	if len(entity_uid) == 0 {
		helpers.BadRequestResponse(w, "Entity is required", errors.New("No entity"), "ENTITY_REQUIRED")
		return
	}

	teamsSQL := `
		SELECT
			t.team_uid,
			t.description,
			e.entity_uid,
			e.name AS entity_name,
			(
				SELECT COUNT(DISTINCT pa.player_uid)
				FROM players_assignation AS pa
				WHERE pa.team_uid = t.team_uid
				AND pa.delete_flag = 0
			) AS players_assigned,
			COALESCE(
				json_agg(
					json_build_object(
						'category_uid', tc.category_uid,
						'name', tc.description
					)
					ORDER BY tc.description
				)
				FILTER (WHERE tc.category_uid IS NOT NULL),
				'[]'::JSON
			) AS categories
		FROM entity_teams AS t
		INNER JOIN entities AS e
			ON e.entity_uid = t.entity_uid
		LEFT JOIN entity_team_categories AS etc
			ON etc.team_uid = t.team_uid
		LEFT JOIN team_categories AS tc
			ON tc.category_uid = etc.category_uid
			AND tc.delete_flag = 0
		WHERE t.entity_uid = $1
		AND t.delete_flag = 0
		GROUP BY
			t.team_uid,
			t.description,
			e.entity_uid,
			e.country_code,
			e.currency_code,
			e.name
		ORDER BY t.description;
	`
	teams := models.Teams{}

	datos, err := db.DB.Query(teamsSQL, entity_uid)
	if err != nil {
		fmt.Println(err)
		helpers.InternalServerErrorResponse(w, "Petition error")
		return
	}
	defer datos.Close()

	for datos.Next() {
		dato := models.Team{}
		var categoriesJSON []byte

		err := datos.Scan(
			&dato.Team_uid,
			&dato.Description,
			&dato.Entity_uid,
			&dato.Entity_name,
			&dato.Players_assigned,
			&categoriesJSON,
		)
		if err != nil {
			helpers.BadRequestResponse(
				w,
				"Error scanning teams",
				err,
				"DATA_ERROR",
			)
			return
		}
		if err := json.Unmarshal(categoriesJSON, &dato.Categories); err != nil {
			helpers.BadRequestResponse(
				w,
				"Error decoding team categories",
				err,
				"DATA_ERROR",
			)
			return
		}

		teams = append(teams, dato)
	}

	if err := datos.Err(); err != nil {
		helpers.InternalServerErrorResponse(w, "Error reading events")
		return
	}

	helpers.SuccessResponse(w, "succes", teams)
}

func GetTemsByIdHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}
	vars := mux.Vars(r)

	teamUID, err := uuid.Parse(vars["id"])
	if err != nil {
		helpers.BadRequestResponse(
			w,
			"Team_uid is required or invalid",
			err,
			"INVALID_TEAM_ID",
		)
		return
	}

	teamSQL := `
		SELECT
			t.team_uid,
			t.description,
			e.entity_uid,
			e.country_code,
			e.currency_code,
			e.name AS entity_name,
			COALESCE(p.player_uid::text, '') AS player_uid,
			COALESCE(pa.player_assignation_uid::text, '') AS assignation_uid,
			COALESCE(p.first_name, '') AS first_name,
			COALESCE(p.last_name, '') AS last_name,
			COALESCE(p.email, '') AS email,
			COALESCE(p.status, 0) AS status
		FROM entity_teams AS t
		INNER JOIN entities AS e
			ON e.entity_uid = t.entity_uid
		LEFT JOIN players_assignation AS pa
			ON pa.team_uid = t.team_uid
			AND pa.delete_flag = 0
		LEFT JOIN players AS p
			ON p.player_uid = pa.player_uid
			AND p.delete_flag = 0
		WHERE t.team_uid = $1;
	`

	rows, err := db.DB.Query(teamSQL, teamUID)
	if err != nil {
		fmt.Println(err)
		helpers.InternalServerErrorResponse(w, "Error obtaining team")
		return
	}
	defer rows.Close()

	team := models.TeamsDetails{
		Players_assigned: []models.PlayerAssigned{},
	}

	found := false

	for rows.Next() {
		var player models.PlayerAssigned

		err := rows.Scan(
			&team.Team_uid,
			&team.Description,
			&team.Entity_uid,
			&team.Country_code,
			&team.Currency_code,
			&team.Entity_name,
			&player.Player_uid,
			&player.Assignation_uid,
			&player.FirstName,
			&player.LastName,
			&player.Email,
			&player.Status,
		)
		if err != nil {
			helpers.BadRequestResponse(
				w,
				"Error scanning team",
				err,
				"DATA_ERROR",
			)
			return
		}

		found = true

		if player.Player_uid != "" {
			team.Players_assigned = append(
				team.Players_assigned,
				player,
			)
		}
	}

	if err := rows.Err(); err != nil {
		helpers.InternalServerErrorResponse(w, "Error reading team")
		return
	}

	if !found {
		helpers.NotFoundResponse(w, "Team not found")
		return
	}

	helpers.SuccessResponse(w, "success", team)
}

func NewTeamHandler(w http.ResponseWriter, r *http.Request) {
	userUID, ok := r.Context().Value("userID").(string)
	if !ok || userUID == "" {
		helpers.UnauthorizedResponse(w, "Usuario no autenticado")
		return
	}
	var teamCreate models.TeamPost
	//FALTA IMPLEMENTAR EL RESTO
	if err := json.NewDecoder(r.Body).Decode(&teamCreate); err != nil {
		helpers.BadRequestResponse(w, "Error obtaining data", err, "ERROR_DATA")
		return
	}
	validationErrors := validations.TeamsPostValidations(teamCreate)

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

	new_uuid := uuid.New()
	teamSQL := "INSERT INTO \"entity_teams\"( team_uid, \"description\", \"entity_uid\") VALUES ($1, $2, $3)"

	_, err := db.DB.Exec(teamSQL, new_uuid.String(), teamCreate.Description, teamCreate.Entity_uid)

	if err != nil {
		helpers.BadRequestResponse(
			w,
			"Error creating",
			err,
			"STORING_ERROR",
		)
		return
	}

	uuidResponse := map[string]string{
		"team_uid": new_uuid.String(),
	}
	helpers.CreatedResponse(w, "success", uuidResponse)

	// if !validations.ValidateContext(w, r) {
	// 	return
	// }
	// var team models.TeamPost
	// //FALTA IMPLEMENTAR EL RESTO
	// if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
	// 	helpers.BadRequestResponse(w, "Error obtaining data", err, "ERROR_DATA")
	// 	return
	// }
	// validationErrors := validations.TeamsPostValidations(team)

	// if len(validationErrors) > 0 {
	// 	message := strings.Join(validationErrors, ", ")

	// 	helpers.BadRequestResponse(
	// 		w,
	// 		message,
	// 		errors.New(message),
	// 		"VALIDATION_ERRORS",
	// 	)
	// 	return
	// }

}
func AssignTeamHandler(w http.ResponseWriter, r *http.Request) {
	userUID, ok := r.Context().Value("userID").(string)
	if !ok || userUID == "" {
		helpers.UnauthorizedResponse(w, "Usuario no autenticado")
		return
	}
	var playerAsign models.PlayerAsign
	//FALTA IMPLEMENTAR EL RESTO
	if err := json.NewDecoder(r.Body).Decode(&playerAsign); err != nil {
		helpers.BadRequestResponse(w, "Error obtaining data", err, "ERROR_DATA")
		return
	}
	validationErrors := validations.TeamsAsignValidations(playerAsign)

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

	teamSQL := "INSERT INTO \"players_assignation\"( player_uid, \"team_uid\", \"assigned_by_uid\") VALUES ($1, $2, $3)"

	_, err := db.DB.Exec(teamSQL, playerAsign.Player_uid, playerAsign.Team_uid, userUID)

	if err != nil {
		helpers.BadRequestResponse(
			w,
			"Error creating",
			err,
			"STORING_ERROR",
		)
		return
	}

	helpers.CreatedResponse(w, "success", nil)
}

func UnassignTeamHandler(w http.ResponseWriter, r *http.Request) {
	userUID, ok := r.Context().Value("userID").(string)
	if !ok || userUID == "" {
		helpers.UnauthorizedResponse(w, "Usuario no autenticado")
		return
	}
	//FALTA IMPLEMENTAR EL RESTO
	assignation_uid := r.URL.Query().Get("assignation_uid")

	teamSQL := "UPDATE \"players_assignation\" SET \"delete_flag\" = $1, \"assigned_by_uid\" = $2 WHERE player_assignation_uid =$3"

	_, err := db.DB.Exec(teamSQL, 1, userUID, assignation_uid)

	if err != nil {
		helpers.BadRequestResponse(
			w,
			"Error creating",
			err,
			"STORING_ERROR",
		)
		return
	}

	helpers.CreatedResponse(w, "success", nil)
}

func GetTeamsCategoriesByEntityHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}
	entity_uid := r.URL.Query().Get("entity_uid")

	if len(entity_uid) == 0 {
		helpers.BadRequestResponse(w, "Entity is required", errors.New("No entity"), "ENTITY_REQUIRED")
		return
	}

	catSql := "SELECT category_uid, description FROM team_categories WHERE entity_uid = $1 AND delete_flag = 0"
	categories := models.TeamCategories{}

	datos, err := db.DB.Query(catSql, entity_uid)
	if err != nil {
		fmt.Println(err)
		helpers.InternalServerErrorResponse(w, "Petition error")
		return
	}
	defer datos.Close()

	for datos.Next() {
		dato := models.TeamCategory{}
		err := datos.Scan(&dato.Category_uid, &dato.Description)
		if err != nil {
			helpers.BadRequestResponse(w, "Error scanning", err, "DATA_ERROR")
			return
		}
		categories = append(categories, dato)
	}

	if err := datos.Err(); err != nil {
		helpers.InternalServerErrorResponse(w, "Error reading events")
		return
	}

	helpers.SuccessResponse(w, "succes", categories)
}

func NewTeamCatHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}
	var cat models.TeamPostCategory
	//FALTA IMPLEMENTAR EL RESTO
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil {
		helpers.BadRequestResponse(w, "Error obtaining data", err, "ERROR_DATA")
		return
	}
	validationErrors := validations.TeamsCatPostValidations(cat)

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

	teamSQL := "INSERT INTO \"team_categories\"( \"description\", \"entity_uid\") VALUES ($1, $2)"

	_, err := db.DB.Exec(teamSQL, cat.Description, cat.Entity_uid)

	if err != nil {
		helpers.BadRequestResponse(
			w,
			"Error creating",
			err,
			"STORING_ERROR",
		)
		return
	}

	helpers.CreatedResponse(w, "success", nil)
}

func AssingCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}
	var teamAssign models.TeamAssignCategory

	if err := json.NewDecoder(r.Body).Decode(&teamAssign); err != nil {
		helpers.BadRequestResponse(w, "Error obtaining data", err, "ERROR_DATA")
		return
	}
	validationErrors := validations.TeamsCatAssignValidations(teamAssign)

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

	tx, err := db.DB.Begin()
	if err != nil {
		helpers.BadRequestResponse(
			w,
			"Error starting category assignment",
			err,
			"ASSIGNMENT_ERROR",
		)
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		`DELETE FROM "entity_team_categories" WHERE "team_uid" = $1::UUID;`,
		teamAssign.Team_uid,
	)
	if err != nil {
		helpers.BadRequestResponse(
			w,
			"Error replacing categories",
			err,
			"ASSIGNMENT_ERROR",
		)
		return
	}

	teamSQL := `
		INSERT INTO "entity_team_categories" ("team_uid", "category_uid")
		SELECT $1::UUID, category_uid
		FROM unnest($2::UUID[]) AS category_uid
		ON CONFLICT ("team_uid", "category_uid") DO NOTHING;
	`
	_, err = tx.Exec(
		teamSQL,
		teamAssign.Team_uid,
		pq.Array(teamAssign.Categories),
	)
	if err != nil {
		helpers.BadRequestResponse(
			w,
			"Error assigning categories",
			err,
			"ASSIGNMENT_ERROR",
		)
		return
	}

	if err := tx.Commit(); err != nil {
		helpers.BadRequestResponse(
			w,
			"Error saving categories",
			err,
			"ASSIGNMENT_ERROR",
		)
		return
	}

	helpers.SuccessResponse(w, "success", nil)
}

func UnssingCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}
	var teamAssign models.TeamAssignCategory

	if err := json.NewDecoder(r.Body).Decode(&teamAssign); err != nil {
		helpers.BadRequestResponse(w, "Error obtaining data", err, "ERROR_DATA")
		return
	}
	validationErrors := validations.TeamsCatAssignValidations(teamAssign)

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

	deleteSQL := `
		DELETE FROM "entity_team_categories"
		WHERE "team_uid" = $1::UUID
		AND "category_uid" = ANY($2::UUID[]);
	`

	_, err := db.DB.Exec(
		deleteSQL,
		teamAssign.Team_uid,
		pq.Array(teamAssign.Categories),
	)
	if err != nil {
		helpers.BadRequestResponse(
			w,
			"Error assigning categories",
			err,
			"ASSIGNMENT_ERROR",
		)
		return
	}

	helpers.SuccessResponse(w, "success", nil)
}
