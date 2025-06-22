package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"example.com/YonkiRating/db"
)
func CreateReviewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//allow only POST
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	img := r.URL.Query().Get("img")
	ratingStr := r.URL.Query().Get("rating") // gets rating as a str
	review := r.URL.Query().Get("review")
	authorIDStr := r.URL.Query().Get("authorID") // gets authorID as a str
	
	// make sure everything is populated
	fields := []string{name,img,ratingStr,review,authorIDStr}
	for _, field := range fields {
		if field == "" {
			http.Error(w, "All fields need to be populated", http.StatusBadRequest)
			return
		}
	}

	// Convert rating to int
	rating, err := strconv.Atoi(ratingStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid rating: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// Convert authorID to int
	authorID, err := strconv.ParseInt(authorIDStr, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid authorID: %s", err.Error()), http.StatusBadRequest)
		return
	}

	rateable := db.Rateable{
		Name: name,
		Img: img,
		Rating: rating, 
		Review: review,
		AuthorID: authorID,
	}
	
	err = db.CreateReview(DB, rateable)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not create review: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	// TODO: OUTPUT JSON
	fmt.Fprintln(w, "Review created")
}


func SearchReviewHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	//allow only GET
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//get only relevant fields
	name := r.URL.Query().Get("name")
	ratingStr := r.URL.Query().Get("rating") // gets rating as a str
	authorIDStr := r.URL.Query().Get("authorID") // gets authorID as a str

	// make sure at least one field is requested
	fields := []string{name,ratingStr,authorIDStr}
	for i, field := range fields {
		// if at least one field is not empty then we can proceed
		if field != "" {
			break
		}
		// if we reach the end of the list and found nothing raise err
		if i == len(fields) - 1{
			http.Error(w, "no fields were specified", http.StatusBadRequest)
			return
		}
	}

	// Convert rating to int or sentinel
	rating, err := strconv.Atoi(ratingStr)
	if err != nil && ratingStr != "" {
		http.Error(w, fmt.Sprintf("invalid rating: %s", err.Error()), http.StatusBadRequest)
		return
	} else if ratingStr == "" {
		// sentinel value for empty param (Golang defaults empty ints it to 0)
		rating = -1
	}
	

	// Convert authorID to int
	authorID, err := strconv.ParseInt(authorIDStr, 10, 64)
	if err != nil && authorIDStr != "" {
		http.Error(w, fmt.Sprintf("invalid authorID: %s", err.Error()), http.StatusBadRequest)
		return
	} else if authorIDStr == "" {
		// sentinel value for empty param (Golang defaults empty ints it to 0)
		authorID = -1
	}

	searchParams := db.Rateable{
		Name: name,
		Rating: rating,
		AuthorID: authorID,
	}

	results, err := db.SearchReview(DB, searchParams)

	if err != nil {
		http.Error(w, fmt.Sprintf("could not get review: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// encode data to w and catch errors
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Internal Server Error",http.StatusInternalServerError)
		return
	}

}

func GetReviewHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	//allow only GET
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// split the path into segments to get the ID
	parts := strings.Split(r.URL.Path, "/")
	// url pattern must be 4 parts /review/id/24
	if(len(parts) != 4){
		http.Error(w, "pattern does not exist", http.StatusNotFound)
		return
	}
	id,err := strconv.ParseInt(parts[3],10,64)
	if err != nil {
		http.Error(w, "'ID' must be an int", http.StatusBadRequest)
		return
	}
	
	result, err := db.GetReview(DB,id)
	if err != nil {
		http.Error(w, "could not execute query", http.StatusInternalServerError)
		return
	}
	// nothing was found and the id value defaulted to 0
	if result.ID == 0 {
		http.Error(w, "could not find any review", http.StatusNotFound)
		return
	}

	// encode data to w and catch errors
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Internal Server Error",http.StatusInternalServerError)
		return
	}
}