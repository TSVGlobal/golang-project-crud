package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (app *App) Initialize() error {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"))

	var err error
	app.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		return fmt.Errorf("could not open DB: %w", err)
	}

	if err := app.DB.Ping(); err != nil {
		return fmt.Errorf("could not ping DB: %w", err)
	}

	app.Router = mux.NewRouter().StrictSlash(true)
	return nil
}

func (app *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, app.Router))
}

func sendResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func sendError(w http.ResponseWriter, statusCode int, err string) {
	sendResponse(w, statusCode, map[string]string{"error": err})
}

func (app *App) getProducts(w http.ResponseWriter, r *http.Request) {
	products, err := getProducts(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusOK, products)
}

func (app *App) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idstr := vars["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	product, err := getProduct(app.DB, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			sendError(w, http.StatusNotFound, err.Error())

		default:
			sendError(w, http.StatusInternalServerError, err.Error())

		}
		return
	}
	sendResponse(w, http.StatusOK, product)
}

func (app *App) createProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&product); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	result, err := product.createProduct(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusCreated, result)
}

func (app *App) updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idstr := vars["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var product Product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&product); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	product.ID = id
	if err := product.updateProduct(app.DB); err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusOK, product)
}

func (app *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idstr := vars["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var product Product
	if err := product.deleteProduct(app.DB, id); err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusOK, map[string]string{"message": "Product deleted successfully"})
}

func (app *App) handleRoutes() {
	app.Router.HandleFunc("/products", app.getProducts).Methods("GET")
	app.Router.HandleFunc("/product/{id}", app.getProduct).Methods("GET")
	app.Router.HandleFunc("/product", app.createProduct).Methods("POST")
	//update product
	app.Router.HandleFunc("/product/{id}", app.updateProduct).Methods("PUT")

	//delete handler
	app.Router.HandleFunc("/product/{id}", app.deleteProduct).Methods("DELETE")

}
