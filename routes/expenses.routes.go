package routes

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Max23strm/pitz-backend/db"
	"github.com/Max23strm/pitz-backend/helpers"
	"github.com/Max23strm/pitz-backend/models"
	"github.com/Max23strm/pitz-backend/queries"
	"github.com/Max23strm/pitz-backend/validations"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func PostNewExpenseHandler(w http.ResponseWriter, r *http.Request) {
	userUID, ok := r.Context().Value("userID").(string)
	if !ok || userUID == "" {
		fmt.Println(ok, userUID)
		helpers.UnauthorizedResponse(w, "Usuario no autenticado")
		return
	}

	var expense models.PostExpenses

	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		helpers.BadRequestResponse(w, "Invalid form", err, "505")
		return
	}

	validationErrors := validations.ExpenesePostValidations(expense)

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
	expenseSql := queries.PostNewExpense
	_, err := db.DB.Exec(
		expenseSql,
		new_uuid.String(),
		expense.Assigned_uid,
		expense.Reason,
		expense.Amount,
		userUID,
		expense.Date,
		expense.Expense_type_uid,
		expense.Schedule_expense_uid,
	)

	if err != nil {
		helpers.BadRequestResponse(w, "Error posting", err, "25")
		return
	}

	helpers.CreatedResponse(w, "ok", new_uuid.String())
}

func GetMonthExpensesHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}
	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		helpers.BadRequestResponse(w, "Invalid date format. Use YYYY-MM-DD", err, "505")
		return
	}

	// Get the first day of the month
	startOfMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)

	// Get the last day of the month by going to the first day of the next month and subtracting a day
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	startOfMonthFormated := startOfMonth.Format("2006-01-02 15:04:05")
	endOfMonthFormated := endOfMonth.Format("2006-01-02 15:04:05")

	expensesSQL := queries.SelectExpensesListByDatesSQL
	expenses := models.MonthlyExpensesGroup{}

	datos, err := db.DB.Query(expensesSQL, startOfMonthFormated, endOfMonthFormated)
	if err != nil || datos.Err() != nil {
		fmt.Println("-------------", err)
		helpers.InternalServerErrorResponse(w, "Error fetching rows")
		return
	}
	for datos.Next() {
		dato := models.MonthlyExpenses{}
		var scheduledExpenseUID, scheduledExpenseName sql.NullString

		err := datos.Scan(
			&dato.Expense_uid,
			&dato.Reason,
			&dato.Assigned_to,
			&dato.Expense_type_uid,
			&dato.Expense_type,
			&scheduledExpenseUID,
			&scheduledExpenseName,
			&dato.Amount,
			&dato.Date,
		)
		if err != nil {
			fmt.Println(err)
			helpers.InternalServerErrorResponse(w, "Error scanning rows")
			return
		}

		expenses = append(expenses, dato)
	}

	helpers.SuccessResponse(w, "ok", expenses)
}
func GetMonthlyExpensesSummaryHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}

	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		helpers.BadRequestResponse(w, "Invalid date format. Use YYYY-MM-DD", err, "INVALID_DATE")
		return
	}

	monthStart := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
	nextMonth := monthStart.AddDate(0, 1, 0)

	summary := models.MonthlyExpensesSummary{
		Pending_scheduled_expenses: make([]models.ScheduledExpenseSummary, 0),
	}
	if err := db.DB.QueryRowContext(
		r.Context(),
		queries.GetMonthlyExpensesTotal,
		monthStart,
		nextMonth,
	).Scan(&summary.Total_spent); err != nil {
		helpers.InternalServerErrorResponse(w, "Error obteniendo el total de gastos")
		return
	}

	rows, err := db.DB.QueryContext(
		r.Context(),
		queries.GetPendingScheduledExpensesForMonth,
		monthStart,
		nextMonth,
	)
	if err != nil {
		helpers.InternalServerErrorResponse(w, "Error obteniendo los gastos programados pendientes")
		return
	}
	for rows.Next() {
		var expense models.ScheduledExpenseSummary
		if err := rows.Scan(
			&expense.Schedule_expense_uid,
			&expense.Name,
			&expense.Amount,
			&expense.Currency_code,
			&expense.Next_due_date,
		); err != nil {
			rows.Close()
			helpers.InternalServerErrorResponse(w, "Error leyendo los gastos programados pendientes")
			return
		}
		summary.Pending_scheduled_expenses = append(summary.Pending_scheduled_expenses, expense)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		helpers.InternalServerErrorResponse(w, "Error leyendo los gastos programados pendientes")
		return
	}
	if err := rows.Close(); err != nil {
		helpers.InternalServerErrorResponse(w, "Error cerrando la consulta de gastos programados")
		return
	}

	var nextPayment models.NextScheduledExpense
	err = db.DB.QueryRowContext(
		r.Context(),
		queries.GetNextScheduledExpenseForMonth,
		monthStart,
		nextMonth,
	).Scan(&nextPayment.Name, &nextPayment.Date)
	if err != nil && err != sql.ErrNoRows {
		helpers.InternalServerErrorResponse(w, "Error obteniendo el próximo pago programado")
		return
	}
	if err == nil {
		summary.Next_scheduled_payment = &nextPayment
	}

	helpers.SuccessResponse(w, "ok", summary)
}

func GetExpensesTypesHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}

	expenseTypeSql := queries.GetExpensesTypes
	expensesTypes := models.ExpenseTypes{}

	datos, err := db.DB.Query(expenseTypeSql)
	if err != nil {
		fmt.Println("Query expense types:", err)
		helpers.InternalServerErrorResponse(w, "Error querying expense types")
		return
	}
	defer datos.Close()
	if err := datos.Err(); err != nil {
		helpers.InternalServerErrorResponse(w, "Error reading expense types")
		return
	}
	for datos.Next() {
		dato := models.ExpenseType{}
		err := datos.Scan(&dato.Expense_type_uid, &dato.Description)
		if err != nil {
			helpers.InternalServerErrorResponse(w, "Error scanning rows")
			return
		}

		expensesTypes = append(expensesTypes, dato)
	}

	helpers.SuccessResponse(w, "ok", expensesTypes)
}
func GetScheduleExpensesHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}

	schExpenseTypeSql := queries.GetScheduledExpenses
	scheduleExpenses := models.ScheduleExpensesGet{}

	datos, err := db.DB.Query(schExpenseTypeSql)
	if err != nil {
		fmt.Println("Query expense types:", err)
		helpers.InternalServerErrorResponse(w, "Error querying schedule expenses")
		return
	}
	defer datos.Close()
	if err := datos.Err(); err != nil {
		helpers.InternalServerErrorResponse(w, "Error reading expense types")
		return
	}
	for datos.Next() {
		dato := models.ScheduleExpenseGet{}
		err := datos.Scan(&dato.Schedule_expense_uid,
			&dato.Expense_type_uid,
			&dato.Expense_type_name,
			&dato.Name,
			&dato.Amount,
			&dato.Currency_code)
		if err != nil {
			helpers.InternalServerErrorResponse(w, "Error scanning rows")
			return
		}

		scheduleExpenses = append(scheduleExpenses, dato)
	}

	helpers.SuccessResponse(w, "ok", scheduleExpenses)
}

func GetMonthExpensesByIdHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}
	expenseSql := queries.ExpenseByID
	vars := mux.Vars(r)
	fmt.Println(vars["id"])
	expenseRow := db.DB.QueryRow(expenseSql, vars["id"])
	expense := models.GetExpenseId{}

	err := expenseRow.Scan(
		&expense.Assigned_to_uid,
		&expense.Assigned_to,
		&expense.Reason,
		&expense.Amount,
		&expense.Registered_by_uid,
		&expense.Registered_by,
		&expense.Date,
		&expense.Expense_type_uid,
		&expense.Expense_type,
		&expense.Schedule_name,
		&expense.Schedule_expense_uid,
	)
	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows {
			helpers.NotFoundResponse(w, "No encontrado")
			return
		}
		helpers.BadRequestResponse(w, "Error obteniendo pago", err, "30")
		return
	}

	helpers.SuccessResponse(w, "ok", expense)
}

func DeleteExpenseByIdHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}
	expenseSql := "UPDATE \"expenses\" SET \"delete_flag\" = 1 WHERE \"expense_uid\" = $1"

	vars := mux.Vars(r)

	if vars["id"] == "undefined" || len(vars["id"]) == 0 {
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Ok",
			"mensaje":   "EL id es necesario",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)
		return
	}

	expenseRow, err := db.DB.Exec(expenseSql, vars["id"])
	if err != nil {
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   "Error obteniendo pago: " + err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)

		return
	}
	rowsAffected, err := expenseRow.RowsAffected()
	if err != nil {
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   "Error obteniendo pago: " + err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)

		return
	}

	if rowsAffected == 0 {
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   "No encontrado",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(respuesta)

		return
	}

	respuesta := map[string]interface{}{
		"isSuccess":   true,
		"estado":      "Ok",
		"payment_uid": vars["id"],
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respuesta)
}

func GetExpensesReport(w http.ResponseWriter, r *http.Request) {

	var paymentRequested models.PaymentFile

	if err := json.NewDecoder(r.Body).Decode(&paymentRequested); err != nil {
		respuesta := map[string]string{
			"estado":  "Error",
			"mensaje": "Error al crear",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)
		return
	}
	validationErrors := validations.PaymentsFileValidation(paymentRequested)

	if len(validationErrors) > 0 {

		var errors []string
		errors = append(errors, validationErrors...)

		respuesta := map[string]interface{}{
			"isSuccese": false,
			"estado":    "Error",
			"mensaje":   errors,
		}

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)

		return
	}

	// startDate, err := time.Parse("2006-01-02", paymentRequested.Start_date)

	expenseDataSql := "SELECT expenses.expense_uid, expenses.reason, assigned.first_name as assigned_first_name,assigned.last_name as assigned_last_name,expenses.amount, expenses.date FROM \"expenses\" INNER JOIN users AS assigned ON expenses.assigned_uid = assigned.user_uid INNER JOIN users AS registrar ON expenses.registered_by_uid = registrar.user_uid WHERE expenses.delete_flag = 0 AND expenses.date BETWEEN $1 AND $2 ORDER by expenses.date DESC"
	monthlySQL := "SELECT TO_CHAR(date, 'MM-YYYY') AS month, SUM(amount) AS total FROM \"expenses\" WHERE expenses.delete_flag = 0 AND expenses.date BETWEEN $1 AND $2 GROUP BY month ORDER BY month;"

	expenseRows, err := db.DB.Query(expenseDataSql, paymentRequested.Start_date, paymentRequested.End_date)
	expenses := models.ExpensesFileRows{}

	if err != nil {
		if err == sql.ErrNoRows {
			respuesta := map[string]interface{}{
				"isSuccess": false,
				"estado":    "Error",
				"mensaje":   "No encontrado",
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(respuesta)
			return
		}
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   "Error obteniendo pago: " + err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)

		return
	}
	for expenseRows.Next() {
		expense := models.ExpensesFileRow{}
		err := expenseRows.Scan(&expense.Expense_uid, &expense.Reason, &expense.Assigned_first_name, &expense.Assigned_last_name, &expense.Amount, &expense.Date)
		if err != nil {
			log.Println("Error obteniendo datos: ", err)
			http.Error(w, "Error scanning data", http.StatusInternalServerError)
			return
		}

		expenses = append(expenses, expense)
	}
	monthlyRows, err := db.DB.Query(monthlySQL, paymentRequested.Start_date, paymentRequested.End_date)
	monthlyPayments := models.MonthlyFileRows{}
	if err != nil {
		if err == sql.ErrNoRows {
			respuesta := map[string]interface{}{
				"isSuccess": false,
				"estado":    "Error",
				"mensaje":   "No encontrado",
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(respuesta)
			return
		}
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   "Error obteniendo pago: " + err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)

		return
	}

	for monthlyRows.Next() {
		var monthStr string
		var total float64

		err := monthlyRows.Scan(&monthStr, &total)
		if err != nil {
			log.Fatal(err)
		}

		// Parse "08-2025" into time.Time
		parsedTime, err := time.Parse("01-2006", monthStr)
		if err != nil {
			log.Fatal("Failed to parse month:", err)
		}

		monthlyPayments = append(monthlyPayments, models.MonthlyFileRow{
			Month:  parsedTime,
			Amount: total,
		})

	}
	file, err := helpers.CreateExpensesExcel(expenses, monthlyPayments)
	if err != nil {
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   "Error generando archivo: " + err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)

		return
	}
	var buf bytes.Buffer
	if err := file.Write(&buf); err != nil {
		http.Error(w, "Error writing Excel file", http.StatusInternalServerError)
		return
	}

	// Set headers to prompt a download
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=report.xlsx")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))
	w.WriteHeader(http.StatusOK)

	// Write the file to the response
	_, err = w.Write(buf.Bytes())
	if err != nil {
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   "Error writing response: " + err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)

		return
	}
}
