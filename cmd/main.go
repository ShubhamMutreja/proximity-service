package main

import (
	"log"
	"net/http"
	"os"
	"proximityService/base"
	"proximityService/database"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type ProximityService struct {
	Router *mux.Router
	Logger *log.Logger
}

func (ps *ProximityService) InitRoutes() {
	//For business
	ps.Router.HandleFunc("/businesses", base.ListAllBusiness).Methods("GET")
	ps.Router.HandleFunc("/business/", base.GetBusiness).Methods("GET")
	ps.Router.HandleFunc("/business/create", base.CreateBusiness).Methods("POST")
	ps.Router.HandleFunc("/business/bulkcreate", base.BulkCreateBusiness).Methods("POST")
	ps.Router.HandleFunc("/business/update", base.UpdateBusiness).Methods("PUT")
	ps.Router.HandleFunc("/business/delete", base.DeleteBusiness).Methods("DELETE")

	//For users
	ps.Router.HandleFunc("/search/nearby", base.GetNearbyBusinesses).Methods("GET")
}

func (ps *ProximityService) Run(addr string) {
	loggedRouter := handlers.LoggingHandler(ps.Logger.Writer(), ps.Router)
	ps.Logger.Fatal(http.ListenAndServe(addr, loggedRouter))
}

func main() {
	//connect to database
	var svc ProximityService

	dB := database.InitDataBase()
	qT := database.InitQuadTree()
	//defer dB.Close()

	_ = database.NewDBService(qT, dB)
	svc.Router = mux.NewRouter()
	svc.Logger = log.New(os.Stdout, "", log.LstdFlags)

	svc.InitRoutes()
	svc.Run(":8080")
}
