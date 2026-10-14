package models

type Team struct {
	Team_uid         string   `json:"team_uid"`
	Description      string   `json:"description"`
	Entity_uid       string   `json:"entity_uid"`
	Entity_name      string   `json:"entity_name"`
	Players_assigned int      `json:"players_assigned"`
	Categories       []TeamCategoryInfo `json:"categories"`
}

type TeamsDetails struct {
	Team_uid         string           `json:"team_uid"`
	Description      string           `json:"description"`
	Entity_uid       string           `json:"entity_uid"`
	Country_code     string           `json:"country_code"`
	Currency_code    string           `json:"currency_code"`
	Entity_name      string           `json:"entity_name"`
	Players_assigned []PlayerAssigned `json:"players_assigned"`
	Categories       []string         `json:"categories"`
}

type TeamPost struct {
	Description string `json:"description"`
	Entity_uid  string `json:"entity_uid"`
}

type Teams []Team

type PlayerAsignation struct {
	Player_assignation_uid string `json:"player_assignation_uid"`
	Player_uid             string `json:"player_uid"`
	Team_uid               string `json:"team_uid"`
	Assigned_by_uid        string `json:"assigned_by_uid"`
}

type PlayerAsignations []PlayerAsignation

type PlayerAsign struct {
	Player_uid string `json:"player_uid"`
	Team_uid   string `json:"team_uid"`
}
type PlayerUnasign struct {
	Assignation_uid string `json:"assignation_uid"`
}

type TeamCategory struct {
	Category_uid string `json:"category_uid"`
	Description  string `json:"description"`
}

type TeamCategoryInfo struct {
	Category_uid string `json:"category_uid"`
	Name         string `json:"name"`
}

type TeamPostCategory struct {
	Entity_uid  string `json:"entity_uid"`
	Description string `json:"description"`
}

type TeamCategories []TeamCategory

type TeamAssignCategory struct {
	Team_uid   string   `json:"team_uid"`
	Categories []string `json:"categories"`
}
