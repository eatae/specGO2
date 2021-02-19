package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// Article struct
type Product struct {
	Id      int `json:"id"`
	Item   	string `json:"item"`
	Amount  int `json:"amount"`
	Price 	string `json:"price"`
}

// ErrorMessage struct
type ErrorMessage struct {
	Message string `json:"Message"`
}

/* Local DataBase */
var productStore = []Product{
	{Id: 1, Item: "Wood chair", Amount: 8, Price: "12"},
	{Id: 2, Item: "Red-wood table", Amount: 4, Price: "26"},
}

/**
 * Show all
 * request: GET
 */
func ShowProducts(response http.ResponseWriter, request *http.Request) {
	fmt.Println("Hint: ShowAllProducts worked...")
	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(productStore) // пишем в Response все Products
	/** TEST **/
	//fmt.Println(Articles)
}

/**
 * Show one
 * request: GET
 * param: url/{id}
 */
func ShowProductById(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	find := false
	// parse id
	innerId, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		json.NewEncoder(response).Encode(err)
	}
	for _, product := range productStore {
		if product.Id == innerId {
			find = true
			json.NewEncoder(response).Encode(product)
		}
	}
	if !find {
		var err = ErrorMessage{Message: "Not found article with that ID"}
		response.WriteHeader(http.StatusNotFound)
		json.NewEncoder(response).Encode(err)
	}
	/** TEST **/
	//fmt.Println(productStore)
}

/**
 * Create
 * request: POST
 * param: raw json
 */
func CreateProduct(response http.ResponseWriter, request *http.Request) {
	/*{
			Id: 3,
			Item: "Steel chair",
			Amount: 10,
			Price: "8"
		}
	*/
	reqBody, _ := ioutil.ReadAll(request.Body)
	var product Product
	json.Unmarshal(reqBody, &product)
	/* add Article in DB */
	productStore = append(productStore, product)
	/* return new Article */
	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(product.Id)
	/** TEST **/
	// fmt.Println(productStore)
}

/**
 * Delete
 * request: DELETE
 * param: url/{id}
 */
func DeleteProduct(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	find := false
	innerId, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		json.NewEncoder(response).Encode(err)
	}

	for index, product := range productStore {
		if product.Id == innerId {
			find = true
			response.WriteHeader(http.StatusAccepted)
			productStore = append(productStore[:index], productStore[index+1:]...)
		}
	}
	if !find {
		var err = ErrorMessage{Message: "Not found article for delete with that ID"}
		response.WriteHeader(http.StatusNotFound)
		json.NewEncoder(response).Encode(err)
	}
}

/**
 * Update
 * request: PUT
 * param: url/{id}
 */
func UpdateProduct(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	find := false
	innerId, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		json.NewEncoder(response).Encode(err)
	}
	for index, product := range productStore {
		if product.Id == innerId {
			find = true
			reqBody, _ := ioutil.ReadAll(request.Body)
			response.WriteHeader(http.StatusAccepted) // Изменяем статус код на 202
			json.Unmarshal(reqBody, &productStore[index]) // перезаписываем всю информацию для статьи с Id
		}
	}

	if !find {
		response.WriteHeader(http.StatusNotFound)
		err := ErrorMessage{Message: "Not found article with that ID. Try use POST first"}
		json.NewEncoder(response).Encode(err)
	}
	/**** TEST ****/
	fmt.Println(productStore)
}


func main() {
	fmt.Println("REST API V2.0 worked....")
	/* init mux Router */
	router := mux.NewRouter().StrictSlash(true)

	/* route: show all */
	router.HandleFunc("/products", ShowProducts).Methods("GET")
	/* route: show one */
	router.HandleFunc("/product/{id}", ShowProductById).Methods("GET")
	/* create */
	router.HandleFunc("/product", CreateProduct).Methods("POST")
	/* delete */
	router.HandleFunc("/product/{id}", DeleteProduct).Methods("DELETE")
	/* update */
	router.HandleFunc("/product/{id}", UpdateProduct).Methods("PUT")

	/* serve */
	log.Fatal(http.ListenAndServe(":8051", router))
}
