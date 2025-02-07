package models

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type Business struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Location Location `json:"location"`
	Phone    string   `json:"phone"`
	City     string   `json:"city"`
	State    string   `json:"state"`
	ZipCode  string   `json:"zip_code"`
}

type NearbySearchRequest struct {
	UserLocation Location `json:"location"`
	Radius       float64  `json:"radius"`
}

type BusinessSearch struct {
	BusinessID int       `json:"business_id"`
	Location   *Location `json:"location"`
}

type BusinessResponse struct {
	Action     string     `json:"action"`
	Businesses []Business `json:"businesses"`
}
