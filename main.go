package main

import (
	"log"
	"net/http"

	"github.com/Max23strm/pitz-backend/db"
	"github.com/Max23strm/pitz-backend/middleware"
	"github.com/Max23strm/pitz-backend/routes"
	"github.com/gorilla/mux"
)

const baseUrl = "/api"

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Or restrict to specific domain
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent) // 204 is more appropriate for preflight
			return
		}

		next.ServeHTTP(w, r)
	})
}

// docker run --name some-postgres -e ¨PSTGRES_USER=max -e POSTGRES_PASSWORD=secretpassword -d postgres
func main() {

	if err := db.DBconnection(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	defer db.CerrarConexion() // ✅ Close on shutdown

	r := mux.NewRouter()
	r.Use(corsMiddleware)
	api := r.PathPrefix(baseUrl).Subrouter()
	api.Use(middleware.AuthMiddleware)

	r.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	r.HandleFunc(baseUrl+"/loginSession", routes.LoginSession).Methods("POST")
	r.HandleFunc(baseUrl+"/passwordRestoration", routes.RestorePassword)

	r.HandleFunc(baseUrl+"/home", routes.HomeHanlder)

	//USERS - PLAYERS
	r.HandleFunc(baseUrl+"/players", routes.GetPlayersHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/players/newPlayer", routes.PostPlayerHandler).Methods("POST")
	r.HandleFunc(baseUrl+"/players/editPlayer/{id}", routes.EditPlayerHandler).Methods("PUT")
	r.HandleFunc(baseUrl+"/players", routes.DeletePlayerHandler).Methods("DELETE")
	r.HandleFunc(baseUrl+"/players/{id}", routes.GetPlayerByIdHandler).Methods("GET")

	//ENTITIES
	r.HandleFunc(baseUrl+"/entities", routes.GetAllEntities).Methods("GET")
	r.Handle(baseUrl+"/entitiesByUser", middleware.AuthMiddleware(http.HandlerFunc(routes.GetEntitiesByUser))).Methods("GET")
	r.HandleFunc(baseUrl+"/newEntity", routes.InsertEntity).Methods("POST")
	r.Handle(baseUrl+"/assignEntity", middleware.AuthMiddleware(http.HandlerFunc(routes.AssignEntity))).Methods("POST")

	//TEAMS
	r.HandleFunc(baseUrl+"/teams", routes.GetTeamsByEntityHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/teams/{id}", routes.GetTemsByIdHandler).Methods("GET")
	r.Handle(baseUrl+"/teams/new", middleware.AuthMiddleware(http.HandlerFunc(routes.NewTeamHandler))).Methods("POST")
	// r.HandleFunc(baseUrl+"/teams/edit/{id}", routes.EditTeamHandler).Methods("PUT")
	r.Handle(baseUrl+"/team/assign", middleware.AuthMiddleware(http.HandlerFunc(routes.AssignTeamHandler))).Methods("POST")
	r.Handle(baseUrl+"/team/unassign", middleware.AuthMiddleware(http.HandlerFunc(routes.UnassignTeamHandler))).Methods("DELETE")

	//TEAMS CATEGORIES
	r.HandleFunc(baseUrl+"/teams-categories", routes.GetTeamsCategoriesByEntityHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/teams-categories/new", routes.NewTeamCatHandler).Methods("POST")
	r.HandleFunc(baseUrl+"/teams-categories/assign", routes.AssingCategoryHandler).Methods("POST")
	r.HandleFunc(baseUrl+"/teams-categories/unassign", routes.UnssingCategoryHandler).Methods("POST")

	//EVENTS
	r.HandleFunc(baseUrl+"/events", routes.GetEventsHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/events/eventById/{id}", routes.GetEventByIdHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/events/newEvent", routes.NewEventHandler).Methods("POST")
	r.HandleFunc(baseUrl+"/events/editEvent/{id}", routes.EditEventHandler).Methods("PUT")
	r.HandleFunc(baseUrl+"/events/getEventsByMonth", routes.GetEventsByMonthHandler).Methods("GET")

	//EVENT TYPES
	r.HandleFunc(baseUrl+"/eventsTypes", routes.GetEventsTypesHandler).Methods("GET")
	//EVENT States
	r.HandleFunc(baseUrl+"/eventsStates", routes.GetEventsStatesHandler).Methods("GET")

	//ASISTANCE TYPES
	r.HandleFunc(baseUrl+"/asistanceTypes", routes.GetAsistanceTypesHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/asistanceByPlayerId/{id}", routes.GetAsistancePlayerbyIdHandler).Methods("GET")

	//PAYMENTS
	r.HandleFunc(baseUrl+"/payments", routes.GetMonthPaymentsHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/payments/new-payment", routes.PostMonthPaymentsHandler).Methods("POST")
	r.HandleFunc(baseUrl+"/paymentsTypes", routes.GetPaymentTypesHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/payments/paymentById/{id}", routes.GetPaymentByIdHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/payments/deletePayment/{id}", routes.DeletePaymentByIdHandler).Methods("DELETE")
	r.HandleFunc(baseUrl+"/payments/getPaymentsReport", routes.GetPaymentsReport).Methods("POST")

	//Expenses
	r.HandleFunc(baseUrl+"/expenses", routes.GetMonthExpensesHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/expenses/new-expense", routes.PostNewExpenseHandler).Methods("POST")
	r.HandleFunc(baseUrl+"/expenses/{id}", routes.GetMonthExpensesByIdHandler).Methods("Get")
	r.HandleFunc(baseUrl+"/expenses/deleteExpense/{id}", routes.DeleteExpenseByIdHandler).Methods("DELETE")
	r.HandleFunc(baseUrl+"/expenses/getExpensesReport", routes.GetExpensesReport).Methods("POST")

	//Users
	r.HandleFunc(baseUrl+"/users/", routes.GetUsersHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/users/basics/{id}", routes.GetBasicUserHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/users/{id}", routes.GetUserHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/updateUserPassword/{id}", routes.UpdateUserHandler).Methods("PUT")

	http.ListenAndServe(":3050", r)
}
