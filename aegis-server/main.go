package main

import (
	"log"
	"time"
	migrations "nfcunha/aegis/database"
	api "nfcunha/aegis/api"
	"nfcunha/aegis/domain/token"
)

func main() {
	// Initialize database and run migrations
	migrations.Migrate()
	
	// Initialize the token blacklist system
	blacklist := token.NewMemoryBlacklist()
	token.InitializeBlacklist(blacklist)
	log.Println("Token blacklist system initialized")
	
	// Start background cleanup job for expired blacklist entries
	// Runs every hour to remove tokens that have naturally expired
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		
		for range ticker.C {
			log.Println("Running blacklist cleanup job")
			blacklist.Cleanup()
			log.Printf("Blacklist cleanup complete. Current size: %d entries", blacklist.Size())
		}
	}()
	
	// Start the API server
	api.RegisterApis()
}