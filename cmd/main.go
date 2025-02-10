package main

import (
	"log"
	"os"

	"github.com/gorilla/mux"

	"proximityService/api"
	"proximityService/base"
	"proximityService/database"
)

func main() {
	//connect to database
	svc := base.ProximityService{}

	dB := database.InitDataBase()
	qT := database.InitQuadTree()
	defer dB.Close()

	dbSvc := database.NewDBService(qT, dB)

	svc.ApiService = api.NewService(dbSvc)
	svc.Router = mux.NewRouter()
	svc.Logger = log.New(os.Stdout, "", log.LstdFlags)

	svc.InitRoutesAndStartServer(":8080")
}
