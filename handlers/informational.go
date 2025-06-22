package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)
var DB *sql.DB

// inject DB pointer from "main" package into this package
func SetDB(db *sql.DB) {
    DB = db
}


// Status:
func StatusHandler (w http.ResponseWriter, req *http.Request) {	
	// Set headers to Json
	w.Header().Set("Content-Type", "application/json")

	// Data
	data := map[string]string{
		"Status":"Active",
	}
	
	// encode data to w and catch errors
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Internal Server Error",http.StatusInternalServerError)
		return
	}
}