package main

import (
	"encoding/json"
	"fmt"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Article struct
type Article struct {
	Id      string `json:"Id"`
	Title   string `json:"Title"`
	Author  string `json:"Author"`
	Content string `json:"Content"`
}

// ErrorMessage struct
type ErrorMessage struct {
	Message string `json:"Message"`
}

/* Local DataBase */
var Articles = []Article{
	{Id: "1", Title: "First title", Author: "First author", Content: "First content"},
	{Id: "2", Title: "Second title", Author: "Second author", Content: "Second content"},
}

type User struct {
	Id       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

var JWTSecretKey = []byte("secret")

/* Local DataBase */
var Users = []User{
	{1, "Bob", "2222"},
	{2, "Alice", "1234"},
}

/**
 * Get Token
 * request: POST
 */
func GetToken(writer http.ResponseWriter, requestPtr *http.Request) {
	// получаем данные пользователя из запроса
	reqBody, _ := ioutil.ReadAll(requestPtr.Body)
	var innerUser User
	// инициализируем пользователя
	json.Unmarshal(reqBody, &innerUser)

	// ищем пользователя
	for _, user := range Users {
		/*
		 * если находим - передаём токен
		 */
		if innerUser.Login == user.Login && innerUser.Password == user.Password {
			// создаём токен
			token := jwt.New(jwt.SigningMethodHS256)
			// claims - нужен для конфигурации токена
			claims := token.Claims.(jwt.MapClaims)
			claims["admin"] = true
			claims["name"] = "New User"
			// время существования токена
			claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
			// получаем токен
			tokenString, err := token.SignedString(JWTSecretKey)
			if err != nil {
				log.Fatal(err)
			}
			// отправляем пользователю
			writer.Write([]byte(tokenString))
			return
		}
	}

	/*
	 * Если не найден - отвечаем Unauthorize
	 */
	writer.WriteHeader(http.StatusUnauthorized)
	writer.Write([]byte("User is not Found."))
}

/* token Validator (middleware)
 * Это распознаватель, который нам подтверждает, что используя наш SecretKey был создан токен.
 * То есть токен наш.
 */
var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return JWTSecretKey, nil
	},
})

/**
 * Show all
 * request: GET
 */
func ShowArticles(writer http.ResponseWriter, requestPtr *http.Request) {
	fmt.Println("Hint: ShowAllArticles worked...")
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(Articles) // пишем в Response все Articles в виде JSON
	/** TEST **/
	//fmt.Println(Articles)
}

/**
 * Show one
 * request: GET
 * param: url/{id}
 */
func ShowArticleById(writer http.ResponseWriter, requestPtr *http.Request) {
	vars := mux.Vars(requestPtr)
	find := false

	for _, article := range Articles {
		if article.Id == vars["id"] {
			find = true
			json.NewEncoder(writer).Encode(article)
		}
	}
	if !find {
		var err = ErrorMessage{Message: "Not found article with that ID"}
		writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(writer).Encode(err)
	}
	/** TEST **/
	//fmt.Println(Articles)
}

/**
 * Create
 * request: POST
 * param: raw json
 */
func CreateArticle(writer http.ResponseWriter, requestPtr *http.Request) {
	/*
		{
			"Id" : "3",
			"Title" : "Title from json POST method",
			"Author" : "Any Name",
			"Content" : "Content from json POST method"
		}
	*/
	reqBody, _ := ioutil.ReadAll(requestPtr.Body)
	var article Article
	json.Unmarshal(reqBody, &article)
	/* add Article in DB */
	Articles = append(Articles, article)
	/* return new Article */
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(article.Id)
	/** TEST **/
	//fmt.Println(Articles)
}

/**
 * Delete
 * request: DELETE
 * param: url/{id}
 */
func DeleteArticle(writer http.ResponseWriter, requestPtr *http.Request) {
	vars := mux.Vars(requestPtr)
	find := false
	for index, article := range Articles {
		if article.Id == vars["id"] {
			find = true
			writer.WriteHeader(http.StatusAccepted)
			Articles = append(Articles[:index], Articles[index+1:]...)
		}
	}
	if !find {
		var err = ErrorMessage{Message: "Not found article for delete with that ID"}
		writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(writer).Encode(err)
	}
}

/**
 * Update
 * request: PUT
 * param: url/{id}
 */
func UpdateArticle(writer http.ResponseWriter, requestPtr *http.Request) {
	vars := mux.Vars(requestPtr)
	idKey := vars["id"] // строка
	find := false

	for index, article := range Articles {
		if article.Id == idKey {
			find = true
			reqBody, _ := ioutil.ReadAll(requestPtr.Body)
			writer.WriteHeader(http.StatusAccepted)   // Изменяем статус код на 202
			json.Unmarshal(reqBody, &Articles[index]) // перезаписываем всю информацию для статьи с Id
		}
	}

	if !find {
		writer.WriteHeader(http.StatusNotFound)
		err := ErrorMessage{Message: "Not found article with that ID. Try use POST first"}
		json.NewEncoder(writer).Encode(err)
	}
	/**** TEST ****/
	fmt.Println(Articles)
}

func main() {
	fmt.Println("REST API V2.0 worked....")
	/* init mux Router */
	router := mux.NewRouter().StrictSlash(true)

	/* get token */
	router.HandleFunc("/auth", GetToken).Methods("POST")

	/* show all */
	router.HandleFunc("/articles", ShowArticles).Methods("GET")

	/* show one */
	router.HandleFunc("/article/{id}", ShowArticleById).Methods("GET")

	/* create (with token) */
	//router.HandleFunc("/article", CreateArticle).Methods("POST")
	router.Handle("/article", jwtMiddleware.Handler(http.HandlerFunc(CreateArticle))).Methods("POST")

	/* delete */
	router.HandleFunc("/article/{id}", DeleteArticle).Methods("DELETE")

	/* update */
	router.HandleFunc("/article/{id}", UpdateArticle).Methods("PUT")

	/* serve */
	log.Fatal(http.ListenAndServe(":8050", router))
}
