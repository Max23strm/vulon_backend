package models

import (
	"time"
)

type PaymentGeneral struct {
	Payment_uid  string    `json:"payment_uid"`
	Player_name  string    `json:"player_name"`
	Creator_name string    `json:"creator_name"`
	Player_uid   string    `json:"player_uid"`
	Amount       float64   `json:"amount"`
	Date         time.Time `json:"date"`
	Payment_name string    `json:"payment_name"`
}

type PostPayments struct {
	Player_uid        string    `json:"player_uid"`
	Payment_reference *string   `json:"payment_reference"`
	Amount            float64   `json:"amount"`
	Comment           *string   `json:"comment"`
	Date              time.Time `json:"date"`
	Payment_type_uid  string    `json:"payment_type_uid"`
	// User_uid          string    `json:"user_uid"`
}
type PaymentsCategory struct {
	Payment_cat_uid   string  `json:"payment_cat_uid"`
	Name              string  `json:"name"`
	Amount            float64 `json:"amount"`
	Payment_type_uid  string  `json:"payment_type_uid"`
	Payment_type_name string  `json:"payment_type_name"`
}
type PaymentsCategories []PaymentsCategory
type PostPaymentsCategory struct {
	Entity_uid       string  `json:"entity_uid"`
	Name             string  `json:"description"`
	Amount           float64 `json:"amount"`
	Payment_type_uid string  `json:"paymetType"`
}

type Payments []PaymentGeneral

type PaymentType struct {
	Payment_type_uid string `json:"payment_type_uid"`
	Payment_name     string `json:"payment_name"`
	System           int8   `json:"is_system"`
}
type PaymentTypes []PaymentType

type PaymentById struct {
	Payment_uid      string    `json:"payment_uid"`
	Player_reference *string   `json:"payment_reference"`
	Amount           string    `json:"amount"`
	Comment          *string   `json:"comment"`
	Date             time.Time `json:"date"`
	Player_name      string    `json:"player_name"`
	Player_uid       string    `json:"player_uid"`
	Payment_name     string    `json:"payment_name"`
	Creator_name     string    `json:"registered_by"`
}
type PaymentUid struct {
	Payment_uid string `json:"payment_uid"`
}

type PaymentFile struct {
	Start_date time.Time `json:"start_date"`
	End_date   time.Time `json:"end_date"`
	File_type  string    `json:"file_type"`
}

type PaymentFileRow struct {
	Payment_uid        string    `json:"payment_uid"`
	Payment_date       time.Time `json:"payment_date"`
	Player_name        string    `json:"player_name"`
	Player_last_name   string    `json:"player_last_name"`
	Registered_by_name string    `json:"registered_by_name"`
	Player_uid         string    `json:"player_uid"`
	Amount             float64   `json:"amount"`
	Comment            string    `json:"comment"`
	Payment_name       string    `json:"payment_name"`
}

type PaymentFileRows []PaymentFileRow

type MonthlyFileRow struct {
	Month  time.Time `json:"month"`
	Amount float64   `json:"amount"`
}

type MonthlyFileRows []MonthlyFileRow
