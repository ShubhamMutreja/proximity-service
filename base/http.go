package base

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"proximityService/api"
	"proximityService/models"
)

type ProximityService struct {
	Router     *mux.Router
	Logger     *log.Logger
	ApiService api.Service
}

func (ps *ProximityService) InitRoutesAndStartServer(addr string) {
	//For business
	ps.Router.HandleFunc("/businesses", ps.ListAllBusiness).Methods("GET")
	ps.Router.HandleFunc("/business/", ps.GetBusiness).Methods("GET")
	ps.Router.HandleFunc("/business/create", ps.CreateBusiness).Methods("POST")
	ps.Router.HandleFunc("/business/bulkcreate", ps.BulkCreateBusiness).Methods("POST")
	ps.Router.HandleFunc("/business/update", ps.UpdateBusiness).Methods("PUT")
	ps.Router.HandleFunc("/business/delete", ps.DeleteBusiness).Methods("DELETE")

	//For users
	ps.Router.HandleFunc("/search/nearby", ps.GetNearbyBusinesses).Methods("GET")

	loggedRouter := handlers.LoggingHandler(ps.Logger.Writer(), ps.Router)
	ps.Logger.Fatal(http.ListenAndServe(addr, loggedRouter))
}

// http handler funcitons
// List all business entites present in DB
func (ps *ProximityService) ListAllBusiness(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	params, _ := url.ParseQuery(request.URL.RawQuery)
	var req models.NearbySearchRequest

	req.UserLocation.Latitude, _ = strconv.ParseFloat(params.Get("latitude"), 64)
	req.UserLocation.Longitude, _ = strconv.ParseFloat(params.Get("longitude"), 64)
	req.Radius, _ = strconv.ParseFloat(params.Get("radius"), 64)

	resp := ps.ApiService.ListAllBusiness(req)
	err := json.NewEncoder(writer).Encode(&resp)
	if err != nil {
		log.Println("There was an error encoding the initialized struct")
	}
}

// Publishes new business entity inside DB
func (ps *ProximityService) CreateBusiness(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	var req models.Business
	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		log.Println("There was an error decoding the request body into the struct")
	}

	resp := ps.ApiService.CreateBusiness(req)

	err = json.NewEncoder(writer).Encode(&resp)
	if err != nil {
		log.Println("There was an error encoding the initialized struct")
	}
}

// Publishes businesses in bulk inside DB
func (ps *ProximityService) BulkCreateBusiness(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	var req []models.Business
	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		log.Println("There was an error decoding the request body into the struct")
	}

	resp := ps.ApiService.BulkCreateBusiness(req)

	err = json.NewEncoder(writer).Encode(&resp)
	if err != nil {
		log.Println("There was an error encoding the initialized struct")
	}
}

// Updates esisting business entity inside DB
func (ps *ProximityService) UpdateBusiness(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	var req models.Business
	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		log.Println("There was an error decoding the request body into the struct")
	}
	resp := ps.ApiService.UpdateBusiness(req)
	err = json.NewEncoder(writer).Encode(&resp)
	if err != nil {
		log.Println("There was an error encoding the initialized struct")
	}
}

// Deletes existing business entity inside DB
func (ps *ProximityService) DeleteBusiness(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	params, _ := url.ParseQuery(request.URL.RawQuery)
	var req models.Business
	req.ID = params.Get("ID")

	resp := ps.ApiService.DeleteBusiness(req)

	err := json.NewEncoder(writer).Encode(&resp)
	if err != nil {
		log.Println("There was an error encoding the initialized struct")
	}
}

// Returns a business by using ID
func (ps *ProximityService) GetBusiness(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	params, _ := url.ParseQuery(request.URL.RawQuery)
	var req models.Business
	req.ID = params.Get("ID")

	resp := ps.ApiService.GetBusiness(req)

	err := json.NewEncoder(writer).Encode(&resp)
	if err != nil {
		log.Println("There was an error encoding the initialized struct")
	}
}

// Returns nearby businesses by using user location and radius
func (ps *ProximityService) GetNearbyBusinesses(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	params, _ := url.ParseQuery(request.URL.RawQuery)
	var req models.NearbySearchRequest

	req.UserLocation.Latitude, _ = strconv.ParseFloat(params.Get("latitude"), 64)
	req.UserLocation.Longitude, _ = strconv.ParseFloat(params.Get("longitude"), 64)
	req.Radius, _ = strconv.ParseFloat(params.Get("radius"), 64)

	resp := ps.ApiService.GetNearbyBusinesses(req)

	err := json.NewEncoder(writer).Encode(&resp)
	if err != nil {
		log.Println("There was an error encoding the initialized struct")
	}
}
