package validations

import (
	"strings"

	"github.com/Max23strm/pitz-backend/models"
)

func PaymentsPostValidations(payment models.PostPayments) []string {

	var validationErrors []string

	if strings.TrimSpace(payment.Player_uid) == "" {
		validationErrors = append(validationErrors, "El nombre de el jugador es requerido.")
	}

	if strings.TrimSpace(payment.Payment_type_uid) == "" {
		validationErrors = append(validationErrors, "El tipo de pago es requerido.")
	}

	if payment.Date.IsZero() {
		validationErrors = append(validationErrors, "La fecha es requerida.")
	}

	if payment.Amount == 0 || payment.Amount < 0 {
		validationErrors = append(validationErrors, "El monto debe ser mayor a 0.")
	}

	return validationErrors
}
func PaymentsCatPostValidations(paymentCat models.PostPaymentsCategory) []string {

	var validationErrors []string

	if strings.TrimSpace(paymentCat.Entity_uid) == "" {
		validationErrors = append(validationErrors, "Entity_uid is required.")
	}
	if strings.TrimSpace(paymentCat.Name) == "" {
		validationErrors = append(validationErrors, "Name is required.")
	}

	if strings.TrimSpace(paymentCat.Payment_type_uid) == "" {
		validationErrors = append(validationErrors, "Payment type is required.")
	}

	if paymentCat.Amount == 0 || paymentCat.Amount < 0 {
		validationErrors = append(validationErrors, "Amount must be a number bigger than 0.")
	}

	return validationErrors
}

func PaymentsFileValidation(event models.PaymentFile) []string {

	var validationErrors []string

	if strings.TrimSpace(event.File_type) == "" {
		validationErrors = append(validationErrors, "El tipo de archivo es requerido.")
	}
	if event.File_type != "excel" && event.File_type != "pdf" {
		validationErrors = append(validationErrors, "El tipo de archivo es incorrecto: "+event.File_type)
	}

	if event.Start_date.IsZero() {
		validationErrors = append(validationErrors, "La fecha de incio es requerida.")
	}
	if event.End_date.IsZero() {
		validationErrors = append(validationErrors, "La fecha final es requerida.")
	}

	if event.Start_date.After(event.End_date) {
		validationErrors = append(validationErrors, "La fecha final debe ser despues de la inicial")
	}

	return validationErrors
}
