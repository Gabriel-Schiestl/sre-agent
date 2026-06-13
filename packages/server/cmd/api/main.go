package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Gabriel-Schiestl/sre-agent/packages/server/config"
	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/data"
	registryhttp "github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/network/http"
	"github.com/Gabriel-Schiestl/sre-agent/packages/server/internal/registry/services"
)

func main() {
	dbCfg, err := config.LoadDB()
	if err != nil {
		log.Fatalf("failed to load db config: %v", err)
	}

	appCfg, err := config.LoadApp()
	if err != nil {
		log.Fatalf("failed to load app config: %v", err)
	}

	db, err := data.Open(dbCfg)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	if err := db.Migrate(); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	if err := os.MkdirAll(appCfg.UploadsDir, 0755); err != nil {
		log.Fatalf("failed to create uploads dir: %v", err)
	}

	suiteDB := data.NewSuiteDB(db)
	microserviceDB := data.NewMicroserviceDB(db)
	runDB := data.NewRunDB(db)
	diagnosisDB := data.NewDiagnosisDB(db)

	suiteSvc := services.NewSuiteService(suiteDB)
	microserviceSvc := services.NewMicroserviceService(microserviceDB)
	// runner and analyst are nil until those modules are implemented
	runSvc := services.NewRunService(runDB, diagnosisDB, nil, nil, appCfg.UploadsDir)

	registryhttp.SetupCORS(appCfg.FrontendURL)
	registryhttp.RegisterRoutes(suiteSvc, microserviceSvc, runSvc)

	port := fmt.Sprintf("%d", appCfg.Port)
	log.Printf("server starting on port %s", port)
	if err := registryhttp.Start(port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
