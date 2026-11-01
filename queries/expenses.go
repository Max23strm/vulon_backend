package queries

var SelectExpensesListByDatesSQL = `SELECT
    expenses.expense_uid,
    expenses.reason,
    CONCAT(assigned.first_name, ' ', assigned.last_name) AS assigned_to,
    expenses.expense_type_uid,
    eT.description AS expense_type,
    expenses.schedule_expense_uid,
    se.name AS scheduled_expense_name,
    expenses.amount,
    expenses.date
	FROM expenses
	INNER JOIN users AS assigned
		ON expenses.assigned_uid = assigned.user_uid
	INNER JOIN users AS registrar
		ON expenses.registered_by_uid = registrar.user_uid
	INNER JOIN expense_type AS eT
		ON expenses.expense_type_uid = eT.expense_type_uid
	LEFT JOIN schedule_expense AS se
		ON se.schedule_expense_uid = expenses.schedule_expense_uid
		AND expenses.schedule_expense_uid IS NOT NULL
	WHERE expenses.delete_flag = 0
	AND expenses.date BETWEEN $1 AND $2
	ORDER BY expenses.date DESC;`

var GetExpensesTypes = `SELECT
	expense_type_uid,
	description
FROM expense_type WHERE delete_flag = 0
`
var GetScheduledExpenses = `SELECT
	schedule_expense_uid,
	schedule_expense.expense_type_uid,
	et.description as expense_type_name,
	name,
	amount,
	currency_code
FROM schedule_expense 
LEFT JOIN expense_type AS et ON schedule_expense.expense_type_uid = et.expense_type_uid
WHERE schedule_expense.delete_flag = 0
`

var GetMonthlyExpensesTotal = `SELECT COALESCE(SUM(amount), 0)
FROM expenses
WHERE delete_flag = 0
	AND date >= $1
	AND date < $2`

var GetPendingScheduledExpensesForMonth = `SELECT
	schedule_expense_uid,
	name,
	amount,
	currency_code,
	next_due_date
FROM schedule_expense
WHERE delete_flag = 0
	AND is_active = 1
	AND next_due_date >= $1
	AND next_due_date < $2
	AND NOT EXISTS (
		SELECT 1
		FROM expenses
		WHERE expenses.schedule_expense_uid = schedule_expense.schedule_expense_uid
			AND expenses.delete_flag = 0
			AND expenses.date >= $1
			AND expenses.date < $2
	)
ORDER BY next_due_date ASC, name ASC`

var GetNextScheduledExpenseForMonth = `SELECT name, next_due_date
FROM schedule_expense
WHERE delete_flag = 0
	AND is_active = 1
	AND next_due_date >= $1
	AND next_due_date < $2
ORDER BY next_due_date ASC, name ASC
LIMIT 1`

var PostNewExpense = `INSERT INTO expenses (
	expense_uid, 
	assigned_uid,
	reason,
	amount,
	created_at_dttm,
	updated_at_dttm,
	registered_by_uid,
	date,
	expense_type_uid,
	schedule_expense_uid
) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $5, $6, $7, $8)`

var ExpenseByID = `SELECT 
	expenses.assigned_uid as assigned_to_uid, 
	CONCAT(assigned.first_name, ' ', assigned.last_name) as assigned_to, 
	expenses.reason, 
	expenses.amount, 
	expenses.registered_by_uid,
	CONCAT(users.first_name, ' ', users.last_name) as registered_by, 
	expenses.date,
	expenses.expense_type_uid,
	types.description as expense_type,
	sc.name as schedule_name,
	expenses.schedule_expense_uid
	FROM expenses 
	INNER JOIN users on users.user_uid = expenses.registered_by_uid 
	INNER JOIN users as assigned on assigned.user_uid = expenses.assigned_uid 
	INNER JOIN expense_type as types on types.expense_type_uid = expenses.expense_type_uid 
	LEFT JOIN schedule_expense as sc on sc.schedule_expense_uid = expenses.schedule_expense_uid 
	WHERE expense_uid = $1`
