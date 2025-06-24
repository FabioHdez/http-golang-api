package main

import (
	"fmt"
	"net/http"
	"os"

	"example.com/YonkiRating/db"
	"example.com/YonkiRating/handlers"
)

func main() {

	fmt.Println("API running...")
	db, err := db.Connect()
	if err != nil {
		fmt.Printf("Error: %s",err.Error())
		os.Exit(1)
	}
	defer db.Close()

	// inject the "db" connection pointer to the "db package"
	handlers.SetDB(db)
	
	// register status
	http.HandleFunc("/status",handlers.StatusHandler)
	http.HandleFunc("/review/create", handlers.CreateReviewHandler)
	http.HandleFunc("/review/search", handlers.SearchReviewHandler)
	http.HandleFunc("/review/id/", handlers.GetReviewHandler)
	http.HandleFunc("/review/update/", func(w http.ResponseWriter, r *http.Request) {fmt.Println("TBA")})
	http.HandleFunc("/review/delete/", handlers.DeleteReviewHandler)

	http.ListenAndServe(":8000",nil)
}