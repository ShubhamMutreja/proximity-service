package base

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"proximityService/api"
	"proximityService/models"
)

// http handler funcitons
func ListAllBusiness(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	resp := api.ListAllBusiness()
	err := json.NewEncoder(writer).Encode(&resp)
	if err != nil {
		log.Println("There was an error encoding the initialized struct")
	}
}
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

func GetNearbyBusinesses(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	var req models.NearbySearchRequest
	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		log.Println("There was an error decoding the request body into the struct")
	}

	resp := api.GetNearbyBusinesses(req)

	err = json.NewEncoder(writer).Encode(&resp)
	if err != nil {
		log.Println("There was an error encoding the initialized struct")
	}
}
