package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nczz/phishguard/config"
	"github.com/nczz/phishguard/internal/db"
	"github.com/nczz/phishguard/internal/tracker"
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
