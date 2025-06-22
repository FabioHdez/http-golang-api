package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
)

func buildConnStr() (string,error) {
	// Read json config. Return a connection string
	file, err := os.Open("db/config.json")
	if err != nil {
		return "",err
	}
	
	defer file.Close()

	var cfg map[string]string

	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return "",err
	}

	user := cfg["user"]
	pass := cfg["pass"]
	host := cfg["host"]
	name := cfg["name"]

	connStr := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s",user,pass,host,name)

	return connStr, nil
}

func Connect() (*sql.DB, error){
	//read json config
	connStr, err := buildConnStr()
	if err != nil { return nil, err }

	//open connection with the connStr
	db, err := sql.Open("mysql",connStr)
	if err != nil { return nil, err }

	fmt.Println("Connected to DB")
	return db, nil
}
