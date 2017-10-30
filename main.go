package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// db is an interface to interact with data on multiple type of data storage
var db Storage
var router *httprouter.Router

func init() {
	var err error

	// TODO: Add configuration to select type of storage, file location or
	// database connection.
	db, err = newStorage(StorageTypeMemory)
	if err != nil {
		log.Fatal(err)
	}

	PopulateBeers()
	PopulateReviews()

	router = httprouter.New()

	router.GET("/beers", GetBeers)
	router.GET("/beers/:id", GetBeer)
	router.GET("/beers/:id/reviews", GetBeerReviews)

	router.POST("/beers", AddBeer)
	router.POST("/beers/:id/reviews", AddBeerReview)
}

func main() {
	fmt.Println("The beer server is on tap now.")
	log.Fatal(http.ListenAndServe(":8080", router))
}
