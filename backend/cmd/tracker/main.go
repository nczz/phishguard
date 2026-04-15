package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/phishguard/phishguard/config"
	"github.com/phishguard/phishguard/internal/db"
	"github.com/phishguard/phishguard/internal/tracker"
)

func main() {
	cfg := config.Load()
	database := db.Init(cfg.DBDSN)

	r := gin.Default()
	h := tracker.NewHandler(database)
	h.RegisterRoutes(r)

	log.Printf("Tracker server starting on %s", cfg.TrackerAddr)
	if err := r.Run(cfg.TrackerAddr); err != nil {
		log.Fatalf("failed to start tracker server: %v", err)
	}
}
