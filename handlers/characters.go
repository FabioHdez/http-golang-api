package handlers

import (
	"encoding/json"
	"net/http"
)

// type Class struct {
// 	Name string
// 	Health int
// 	Mana int
// 	Agility int
// 	Intelligence int
// 	Strength int
// }

type Class struct {
	Name string
	Score float32
	Status string
	Year int
	Genre string
	Img string
}

func AllClassesHandler(w http.ResponseWriter, req *http.Request) {
	// Set headers to Json
	w.Header().Set("Content-Type", "application/json")
	
	data := Class{
		Name: "Naruto",
		Score: 10,
		Status: "Finished",
		Year: 2007,
		Genre: "Accion, Aventura",
		Img: "https://cdn.myanimelist.net/images/anime/1565/111305.jpg",
	}

	// encode data and catch errors
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Internal Server Error",http.StatusInternalServerError)
		return
	}


}