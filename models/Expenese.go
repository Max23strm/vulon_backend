package models

import "time"

type PostExpenses struct {
	Reason               string    `json:"reason"`
	Amount               float64   `json:"amount"`
	Date                 time.Time `json:"date"`
	Registered_by        string    `json:"registered_by"`
	Assigned_uid         string    `json:"assigned_uid"`
	Expense_type_uid     string    `json:"expense_type_uid"`
	Schedule_expense_uid *string   `json:"schedule_expense_uid"`
}
type GetExpenseId struct {
	Assigned_to_uid      string    `json:"assigned_to_uid"`
	Assigned_to          string    `json:"assigned_to"`
	Reason               string    `json:"reason"`
	Amount               string    `json:"amount"`
	Registered_by_uid    string    `json:"registered_by_uid"`
	Registered_by        string    `json:"registered_by"`
	Date                 time.Time `json:"date"`
	Expense_type_uid     string    `json:"expense_type_uid"`
	Expense_type         string    `json:"expense_type"`
	Schedule_name        *string   `json:"schedule_name"`
	Schedule_expense_uid *string   `json:"schedule_expense_uid"`
}

type MonthlyExpenses struct {
	Assigned_to            string    `json:"assigned_to"`
	Amount                 float64   `json:"amount"`
	Date                   time.Time `json:"date"`
	Reason                 string    `json:"reason"`
	Expense_type_uid       string    `json:"expense_type_uid"`
	Scheduled_expense_uid  *string   `json:"scheduled_expense_uid"`
	Scheduled_expense_name *string   `json:"scheduled_expense_name"`
	Expense_type           string    `json:"expense_type"`
	Expense_uid            string    `json:"expense_uid"`
}

type MonthlyExpensesGroup []MonthlyExpenses

type ExpenseType struct {
	Expense_type_uid string `json:"expense_type_uid"`
	Description      string `json:"description"`
}

type ExpenseTypes []ExpenseType

type ScheduleExpenseGet struct {
	Schedule_expense_uid string  `json:"schedule_expense_uid"`
	Expense_type_uid     string  `json:"expense_type_uid"`
	Expense_type_name    string  `json:"expense_type_name"`
	Name                 string  `json:"name"`
	Amount               float32 `json:"amount"`
	Currency_code        string  `json:"currency_code"`
}

type ScheduleExpensesGet []ScheduleExpenseGet

type ScheduledExpenseSummary struct {
	Schedule_expense_uid string    `json:"schedule_expense_uid"`
	Name                 string    `json:"name"`
	Amount               float64   `json:"amount"`
	Currency_code        string    `json:"currency_code"`
	Next_due_date        time.Time `json:"next_due_date"`
}

type NextScheduledExpense struct {
	Name string    `json:"name"`
	Date time.Time `json:"date"`
}

type MonthlyExpensesSummary struct {
	Total_spent                float64                   `json:"total_spent"`
	Pending_scheduled_expenses []ScheduledExpenseSummary `json:"pending_scheduled_expenses"`
	Next_scheduled_payment     *NextScheduledExpense     `json:"next_scheduled_payment"`
}

type ExpensesFileRow struct {
	Expense_uid         string    `json:"expense_uid"`
	Reason              string    `json:"reason"`
	Amount              float64   `json:"amount"`
	Date                time.Time `json:"date"`
	Assigned_first_name string    `json:"assigned_first_name"`
	Assigned_last_name  string    `json:"assigned_last_name"`
	Assigned_uid        string    `json:"assigned_uid"`
}

type ExpensesFileRows []ExpensesFileRow
