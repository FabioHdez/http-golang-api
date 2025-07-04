package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// package db

type Rateable struct {
	ID       int64  `json:"ID"`
	Name     string `json:"name"`
	Img      string `json:"img"`
	Rating   int    `json:"rating"`
	Review   string `json:"review"`
	AuthorID int64  `json:"authorID"`
	Active   int64  `json:"Active"`
}

func CreateReview(db *sql.DB, review Rateable) error {
	//prep query
	query,err := db.Prepare(`
		INSERT INTO rateables (name, img, rating, review, author_id)
		VALUES (?,?,?,?,?);
	`)
	if err != nil {
		return err
	}
	defer query.Close()

 	res, err := query.Exec(review.Name, review.Img, review.Rating, review.Review, review.AuthorID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}


	if rowsAffected == 0 {
		return errors.New("no rows were affected")
	}

	// print rows affected
	fmt.Printf("Rows affected: %v\n", rowsAffected)
	return nil
}

// refactor this
func SearchReview(db *sql.DB, searchParams Rateable) ([]Rateable,error){
	// where clause wildcards
	// wrap values in %{value}% for searching
	name := "%"+searchParams.Name+"%"
	var ratingStr string
	if searchParams.Rating == -1 { //no rating was specified
		ratingStr = "%"
	}else {
		ratingStr = "%"+strconv.Itoa(searchParams.Rating)+"%"
	}

	var authorIDStr string
	if searchParams.AuthorID == -1 { //no author was specified
		authorIDStr = "%"
	}else {
		authorIDStr = "%"+strconv.FormatInt(searchParams.AuthorID,10)+"%"
	}
	
	query:=`
		SELECT * FROM rateables WHERE
		name like ? AND
		rating like ? AND
		author_id like ? AND 
		active = 1;
	`
	
	rows, err := db.Query(query, name, ratingStr, authorIDStr)
	if err != nil {
		return nil,err
	}
	// db.Query() returns a (*sql.Rows) which mantains an open connection to the db
	defer rows.Close()

	var results []Rateable

	for rows.Next() {
		var r Rateable
		err := rows.Scan(&r.ID,&r.Name,&r.Img,&r.Rating,&r.Review,&r.AuthorID,&r.Active)
		if err != nil{
			return nil,err
		}
		results = append(results, r)
	}

	return results,nil
}

func GetReview(db *sql.DB, id int64) (Rateable, error) {
	query:=`
		SELECT * FROM rateables WHERE
		id = ?
		AND active = 1;
	`
	row := db.QueryRow(query, id)

	var review Rateable
	row.Scan(
		&review.ID,
		&review.Name,
		&review.Img,
		&review.Rating,
		&review.Review,
		&review.AuthorID,
		&review.Active)

	err := row.Err()

	return review, err
}

func DeleteReview(db *sql.DB, id int64) error {
	query,err := db.Prepare(`
		UPDATE rateables
		SET active = 0
		WHERE id = ? 
		AND active = 1;
	`)
	if err != nil {
		return err
	}
	defer query.Close()

	res, err := query.Exec(id)	
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows were affected")
	}

	// print rows affected
	fmt.Printf("Rows affected: %v\n", rowsAffected)
	return nil
}

func UpdateReview(db *sql.DB, reviewChanges Rateable)  error {

	// Build a map based on the Rateables Struct to have some null values
	// Null values are required in order to take advantage of the COALESCE function in mySQL
	changes := make(map[string]interface{})

	if reviewChanges.Name != ""{
		changes["name"] = reviewChanges.Name
	}

	if reviewChanges.Img != ""{
		changes["img"] = reviewChanges.Img
	}
	
	if reviewChanges.Rating != -1 {
		changes["rating"] = reviewChanges.Rating
	}

	if reviewChanges.Review != ""{
		changes["review"] = reviewChanges.Review
	}

	if reviewChanges.AuthorID != -1 {
		changes["authorID"] = reviewChanges.AuthorID
	}

	query,err := db.Prepare(`
		UPDATE rateables
		SET 
		name = COALESCE(?,name),
		img = COALESCE(?,img),
		rating = COALESCE(?,rating),
		review = COALESCE(?,review),
		author_id = COALESCE(?,author_id)
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer query.Close()

	res, err := query.Exec(
		changes["name"],
		changes["img"],
		changes["rating"],
		changes["review"],
		changes["authorID"],
		reviewChanges.ID)
	
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows were affected")
	}

	// print rows affected
	fmt.Printf("Rows affected: %v\n", rowsAffected)
	return nil

}