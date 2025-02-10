package base

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"proximityService/api"
	"proximityService/models"
	"strconv"
)

// http handler funcitons
// List all business entites present in DB
func ListAllBusiness(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	params, _ := url.ParseQuery(request.URL.RawQuery)
	var req models.NearbySearchRequest

	req.UserLocation.Latitude, _ = strconv.ParseFloat(params.Get("latitude"), 64)
	req.UserLocation.Longitude, _ = strconv.ParseFloat(params.Get("longitude"), 64)
	req.Radius, _ = strconv.ParseFloat(params.Get("radius"), 64)

	resp := api.ListAllBusiness(req)
	err := json.NewEncoder(writer).Encode(&resp)
	if err != nil {
		log.Println("There was an error encoding the initialized struct")
	}
}

// Publishes new business entity inside DB
func CreateBusiness(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	var req models.Business
	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		log.Println("There was an error decoding the request body into the struct")
	}

	resp := api.CreateBusiness(req)

	err = json.NewEncoder(writer).Encode(&resp)
	if err != nil {
		log.Println("There was an error encoding the initialized struct")
	}
}

// Publishes businesses in bulk inside DB
func BulkCreateBusiness(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	var req []models.Business
	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		log.Println("There was an error decoding the request body into the struct")
	}

	resp := api.BulkCreateBusiness(req)

	err = json.NewEncoder(writer).Encode(&resp)
	if err != nil {
		log.Println("There was an error encoding the initialized struct")
	}
}

// Updates esisting business entity inside DB
func UpdateBusiness(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	var req models.Business
	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		log.Println("There was an error decoding the request body into the struct")
	}
	resp := api.UpdateBusiness(req)
	err = json.NewEncoder(writer).Encode(&resp)
	if err != nil {
		log.Println("There was an error encoding the initialized struct")
	}
}

// Deletes existing business entity inside DB
func DeleteBusiness(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	params, _ := url.ParseQuery(request.URL.RawQuery)
	var req models.Business
	req.ID = params.Get("ID")

	resp := api.DeleteBusiness(req)

	err := json.NewEncoder(writer).Encode(&resp)
	if err != nil {
		log.Println("There was an error encoding the initialized struct")
	}
}

// Returns a business by using ID
func GetBusiness(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	params, _ := url.ParseQuery(request.URL.RawQuery)
	var req models.Business
	req.ID = params.Get("ID")

	resp := api.GetBusiness(req)

	err := json.NewEncoder(writer).Encode(&resp)
	if err != nil {
		log.Println("There was an error encoding the initialized struct")
	}
}

// Returns nearby businesses by using user location and radius
func GetNearbyBusinesses(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	params, _ := url.ParseQuery(request.URL.RawQuery)
	var req models.NearbySearchRequest

	req.UserLocation.Latitude, _ = strconv.ParseFloat(params.Get("latitude"), 64)
	req.UserLocation.Longitude, _ = strconv.ParseFloat(params.Get("longitude"), 64)
	req.Radius, _ = strconv.ParseFloat(params.Get("radius"), 64)

	resp := api.GetNearbyBusinesses(req)

	err := json.NewEncoder(writer).Encode(&resp)
	if err != nil {
		log.Println("There was an error encoding the initialized struct")
	}
}
