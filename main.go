package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Product struct {
	Id       string
	Name     string
	Quantity int
	Price    float64
}

var Products []Product

type Data_res struct {
	id   int
	name string
}

func welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my website!")
	log.Println("Endpoint Hit: welcome")

}

func returnAllProducts(w http.ResponseWriter, r *http.Request) {
	log.Infoln("returnAllProducts Endpoint Hit")
	json.NewEncoder(w).Encode(Products)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	log.Info("getProduct Endpoint Hit")
	params := mux.Vars(r)
	id := params["id"]
	for _, product := range Products {
		if product.Id == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(product)

		}
	}

}

func checkError(err error) {
	if err != nil {
		log.Error(err)
	}

}

func handleRequests() {
	r := mux.NewRouter()
	r.HandleFunc("/products", returnAllProducts).Methods("GET")
	r.HandleFunc("/product/{id}", getProduct).Methods("GET")
	r.HandleFunc("/hello", welcome).Methods("GET")
	http.ListenAndServe("127.0.0.1:9989", r)
}

// main
// func main() {
// 	log.SetLevel(log.InfoLevel)
// 	Products = append(Products, Product{Id: "1", Name: "Laptop", Quantity: 10, Price: 1000.0})
// 	Products = append(Products, Product{Id: "2", Name: "Mobile", Quantity: 5, Price: 500.0})
// 	handleRequests()
// }

func main() {
	connectionString := "user1:user1@tcp(127.0.0.1:3306)/learning"

	db, err := sql.Open("mysql", connectionString)

	checkError(err)

	defer db.Close()

	log.Info("Connected to database")

	result, err := db.Exec("INSERT INTO data VALUES (7, 'test')")
	checkError(err)
	lastInsertId, err := result.LastInsertId()
	log.Info("last insert id: ", lastInsertId)
	rowsAffected, err := result.RowsAffected()
	log.Info("rows affected: ", rowsAffected)

	checkError(err)

	row, err := db.Query("SELECT * FROM data")
	checkError(err)
	for row.Next() {
		var data Data_res
		err = row.Scan(&data.id, &data.name)
		checkError(err)
		log.Info(data)
	}
}
