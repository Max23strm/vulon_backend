package models

import (
	"time"
)

type Player struct {
	Player_uid string   `json:"player_uid"`
	FirstName  string   `json:"firstName"`
	LastName   string   `json:"last_name"`
	Email      string   `json:"email"`
	Status     int16    `json:"status"`
	Assigned   []string `json:"assigned"`
}
type PlayerAssigned struct {
	Player_uid      string `json:"player_uid"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Status          int16  `json:"status"`
	Assignation_uid string `json:"assignation_uid"`
	Positions       []int8 `json:"positions"`
}

type PlayerDetails struct {
	Player_uid       string    `json:"player_uid"`
	FirstName        string    `json:"firstName"`
	LastName         string    `json:"last_name"`
	Email            string    `json:"email"`
	Status           string    `json:"status"`
	Birth_dt         time.Time `json:"birth_dt"`
	Positions        []*int8   `json:"positions"`
	Address          *string   `json:"address"`
	Sex              string    `json:"sex"`
	BloodType        *string   `json:"blood_type"`
	Comments         *string   `json:"comments"`
	Credential       *string   `json:"credential"`
	Afiliation       *string   `json:"afiliation"`
	Curp             *string   `json:"curp"`
	Enfermedad       *string   `json:"enfermedad"`
	Phone_number     *string   `json:"phone_number"`
	Emergency_number *string   `json:"emergency_number"`
	Insurance        bool      `json:"insurance"`
	Insurance_name   *string   `json:"insurance_name"`
}

type PostPlayerDetails struct {
	Player_uid       string    `json:"player_uid"`
	FirstName        string    `json:"firstName"`
	LastName         string    `json:"last_name"`
	Email            string    `json:"email"`
	Status           int16     `json:"status"`
	Birth_dt         time.Time `json:"birth_dt"`
	Positions        []*int8   `json:"positions,omitempty"`
	Address          *string   `json:"address,omitempty"`
	Sex              string    `json:"sex"`
	BloodType        *string   `json:"blood_type,omitempty"`
	Comments         *string   `json:"comments,omitempty"`
	Credential       *string   `json:"credential,omitempty"`
	Afiliation       *string   `json:"afiliation,omitempty"`
	Curp             *string   `json:"curp,omitempty"`
	Enfermedad       *string   `json:"enfermedad,omitempty"`
	Phone_number     *string   `json:"phone_number,omitempty"`
	Emergency_number *string   `json:"emergency_number,omitempty"`
	Insurance        bool      `json:"insurance"`
	Insurance_name   *string   `json:"insurance_name,omitempty"`
}
type PutPlayerDetails struct {
	FirstName        *string    `json:"firstName,omitempty"`
	LastName         *string    `json:"last_name,omitempty"`
	Email            *string    `json:"email,omitempty"`
	Status           *int16     `json:"status,omitempty"`
	Birth_dt         *time.Time `json:"birth_dt,omitempty"`
	Address          *string    `json:"address,omitempty"`
	Sex              *string    `json:"sex,omitempty"`
	BloodType        *string    `json:"blood_type,omitempty"`
	Comments         *string    `json:"comments,omitempty"`
	Credential       *string    `json:"credential,omitempty"`
	Afiliation       *string    `json:"afiliation,omitempty"`
	Curp             *string    `json:"curp,omitempty"`
	Enfermedad       *string    `json:"enfermedad,omitempty"`
	Phone_number     *string    `json:"phone_number,omitempty"`
	Emergency_number *string    `json:"emergency_number,omitempty"`
	Insurance        *bool      `json:"insurance,omitempty"`
	Insurance_name   *string    `json:"insurance_name,omitempty"`
}
type PlayerDetailsWithAsistance struct {
	GeneralInfo Player     `json:"player_data"`
	Address     *string    `json:"address"`
	Asistance   Asistances `json:"asistance"`
}

type Positions struct {
	Position_uid string `json:"positions_uid"`
	Positions    string `json:"positions"`
}

type Players []Player

type NewPlayer struct {
	Player_uid string
	FirstName  string
	LastName   string
	Email      string
	Status     int16
	Positions  string
}
