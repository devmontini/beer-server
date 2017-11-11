package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	model "github.com/golang-bristol/beer-model"
)

func TestGetBeers(t *testing.T) {
	var cellarFromRequest []model.Beer
	var cellarFromStorage []model.Beer

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/beers", nil)

	router.ServeHTTP(w, r)

	cellarFromStorage = db.FindBeers()
	json.Unmarshal(w.Body.Bytes(), &cellarFromRequest)

	if w.Code != http.StatusOK {
		t.Errorf("Expected route GET /beers to be valid.")
		t.FailNow()
	}

	if len(cellarFromRequest) != len(cellarFromStorage) {
		t.Error("Expected number of beers from request to be the same as beers in the storage")
		t.FailNow()
	}

	var mapCellar = make(map[model.Beer]int, len(cellarFromStorage))
	for _, beer := range cellarFromStorage {
		mapCellar[beer] = 1
	}

	for _, beerResp := range cellarFromRequest {
		if _, ok := mapCellar[beerResp]; !ok {
			t.Errorf("Expected all results to match existing records")
			t.FailNow()
			break
		}
	}
}

func TestAddBeer(t *testing.T) {
	newBeer := model.Beer{
		Name:    "Testing beer",
		Abv:     333,
		Brewery: "Testing Beer Inc",
	}

	newBeerJSON, err := json.Marshal(newBeer)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/beers", bytes.NewBuffer(newBeerJSON))

	router.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected route POST /beers to be valid.")
		t.FailNow()
	}

	newBeerMissing := true
	for _, b := range db.FindBeers() {
		if b.Name == newBeer.Name &&
			b.Abv == newBeer.Abv &&
			b.Brewery == newBeer.Brewery {
			newBeerMissing = false
		}
	}

	if newBeerMissing {
		t.Errorf("Expected to find new entry in storage`")
		t.FailNow()
	}

}

func TestGetBeer(t *testing.T) {
	cellar := db.FindBeers()
	choice := rand.Intn(len(cellar) - 1)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", fmt.Sprintf("/beers/%d", cellar[choice].ID), nil)

	router.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected route GET /beers/%d to be valid.", cellar[choice].ID)
		t.FailNow()
	}

	var selectedBeer model.Beer
	json.Unmarshal(w.Body.Bytes(), &selectedBeer)

	if cellar[choice] != selectedBeer {
		t.Errorf("Expected to match results with selected beer")
		t.FailNow()
	}

}
