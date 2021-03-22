package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Items -------------------------------------------------

// Article struct
type Article struct {
	Id      string `json:"Id"`
	Title   string `json:"Title"`
	Author  string `json:"Author"`
	Content string `json:"Content"`
}

// User struct
type User struct {
	Id       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

// ResponseMessage struct
const successStatus string = "success"
const errorStatus string = "error"

type ResponseMessage struct {
	Status  string      `json:"Status"`
	Message interface{} `json:"Message"`
}

func (resMsg *ResponseMessage) setError(msg interface{}) {
	resMsg.Status = errorStatus
	resMsg.Message = msg
}

func (resMsg *ResponseMessage) setSuccess(msg interface{}) {
	resMsg.Status = successStatus
	resMsg.Message = msg
}

// Db settings struct
type DbConf struct {
	Port     string
	User     string
	Password string
	DBname   string
	SSlmode  string
}

// DB instance
var DB *sql.DB

// Postgres conf
var pgConf = DbConf{
	Port:     "54320",
	User:     "postgres",
	Password: "2222",
	DBname:   "specgo",
	SSlmode:  "disable",
}

var connectStr = fmt.Sprintf("user=%v password=%v port= %v dbname=%v sslmode=%v",
	pgConf.User, pgConf.Password, pgConf.Port, pgConf.DBname, pgConf.SSlmode)

// SecretKey
var JWTSecretKey = []byte("secret")

// --------------------------------------------------------
// Functions
// --------------------------------------------------------

// Insert User
func InsertUser(innerUser User) ResponseMessage {
	var resMsg ResponseMessage
	// Check Login exists
	if userLoginExists(innerUser.Login) {
		resMsg.setError("Login already exists.")
		return resMsg
	}
	// Query
	_, err := DB.Exec("INSERT INTO users (login, password) VALUES ($1, $2)", innerUser.Login, innerUser.Password)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	return resMsg
}

// Check User login in DB
func userLoginExists(innerLogin string) bool {
	var count int
	err := DB.QueryRow("SELECT count(id) FROM users WHERE login = $1", innerLogin).Scan(&count)
	if err != nil {
		panic(err)
	}
	if count > 0 {
		return true
	}
	return false
}

// Check User Pass and Login in DB
func userExists(innerUser User) bool {
	var count int
	err := DB.QueryRow("SELECT count(id) FROM users WHERE login = $1 AND password = $2",
		innerUser.Login, innerUser.Password).Scan(&count)
	if err != nil {
		panic(err)
	}
	if count > 0 {
		return true
	}
	return false
}

// --------------------------------------------------------
// Actions
// --------------------------------------------------------

// Register user
//
// request: POST
func Register(writer http.ResponseWriter, request *http.Request) {
	reqBody, _ := ioutil.ReadAll(request.Body)
	var innerUser User
	json.Unmarshal(reqBody, &innerUser)
	// set Header JSON
	writer.Header().Set("Content-Type", "application/json")
	// Insert
	resMsg := InsertUser(innerUser)

	if resMsg.Status == errorStatus {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(resMsg)
		return
	}
	resMsg.setSuccess("User registered.")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(resMsg)
}

// Get Token User
//
// request: POST
func GetToken(writer http.ResponseWriter, requestPtr *http.Request) {
	var innerUser User
	var resMsg ResponseMessage
	reqBody, _ := ioutil.ReadAll(requestPtr.Body)
	json.Unmarshal(reqBody, &innerUser)
	// set Header JSON
	writer.Header().Set("Content-Type", "application/json")
	if !userExists(innerUser) {
		resMsg.setError("Password or login does not exist")
		// set Header Status Unauthorized
		writer.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(writer).Encode(resMsg)
		return
	}
	// create Token
	token := jwt.New(jwt.SigningMethodHS256)
	// claims - нужен для конфигурации токена
	claims := token.Claims.(jwt.MapClaims)
	claims["admin"] = true
	claims["name"] = "New User"
	// Token lifetime
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	// получаем токен
	tokenString, err := token.SignedString(JWTSecretKey)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	// send Token
	resMsg.setSuccess(tokenString)
	json.NewEncoder(writer).Encode(resMsg)
}

// token Validator (middleware)
// Это распознаватель, который нам подтверждает, что используя наш SecretKey был создан токен.
// То есть токен наш.
var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return JWTSecretKey, nil
	},
})

// Show all users
//
// request: GET
func ShowUsers(writer http.ResponseWriter, requestPtr *http.Request) {
	rows, err := DB.Query("SELECT * FROM users")
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	defer rows.Close()

	user := User{}
	users := make([]User, 0)
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Login, &user.Password)
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
			return
		}
		users = append(users, user)
	}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(users)
}

// Articles -----------------------------------------------

// Show all articles
//
// request: GET
func ShowArticles(writer http.ResponseWriter, requestPtr *http.Request) {
	rows, err := DB.Query("SELECT * FROM articles")
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	defer rows.Close()

	article := Article{}
	articles := make([]Article, 0)
	for rows.Next() {
		err := rows.Scan(&article.Id, &article.Title, &article.Author, &article.Content)
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
			return
		}
		articles = append(articles, article)
	}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(articles)
}

// Show one article
//
// request: GET
// param: url/{id}
func ShowArticleById(writer http.ResponseWriter, requestPtr *http.Request) {
	vars := mux.Vars(requestPtr)
	row := DB.QueryRow("SELECT * FROM articles WHERE id = $1", vars["id"])
	article := Article{}
	err := row.Scan(&article.Id, &article.Title, &article.Author, &article.Content)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(article)
}

// Create article
//
// request: POST
// param: raw json
func CreateArticle(writer http.ResponseWriter, requestPtr *http.Request) {
	reqBody, _ := ioutil.ReadAll(requestPtr.Body)
	article := Article{}
	json.Unmarshal(reqBody, &article)
	/* add Article in DB */
	err := DB.QueryRow("INSERT INTO articles (title, author, content) VALUES ($1, $2, $3) RETURNING id",
		article.Title, article.Author, article.Content).Scan(&article.Id)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
		return
	}
	var resMsg ResponseMessage
	resMsg.setSuccess(article)
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(resMsg)
}

// Delete article
//
//request: DELETE
//param: url/{id}
func DeleteArticle(writer http.ResponseWriter, requestPtr *http.Request) {
	vars := mux.Vars(requestPtr)
	var resMsg ResponseMessage
	result, err := DB.Exec("DELETE FROM articles WHERE id = $1", vars["id"])
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
		return
	}
	if rows, _ := result.RowsAffected(); rows < 1 {
		resMsg.setError("Not found article for delete with that ID")
		writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(writer).Encode(resMsg)
		return
	}

	resMsg.setSuccess("Article successfully deleted.")
	json.NewEncoder(writer).Encode(resMsg)
}

// Update update
//
// request: PUT
// param: url/{id}
func UpdateArticle(writer http.ResponseWriter, requestPtr *http.Request) {
	vars := mux.Vars(requestPtr)
	reqBody, _ := ioutil.ReadAll(requestPtr.Body)
	resMsg := ResponseMessage{}
	innerArticle := Article{}

	json.Unmarshal(reqBody, &innerArticle)
	row := DB.QueryRow("SELECT * FROM articles WHERE id = $1", vars["id"])
	dbArticle := Article{}
	err := row.Scan(&dbArticle.Id, &dbArticle.Title, &dbArticle.Author, &dbArticle.Content)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
		return
	}
	// update
	if innerArticle.Content != "" {
		dbArticle.Content = innerArticle.Content
	}
	if innerArticle.Author != "" {
		dbArticle.Author = innerArticle.Author
	}
	if innerArticle.Title != "" {
		dbArticle.Title = innerArticle.Title
	}
	_, err = DB.Exec(`UPDATE articles 
								SET title = $1,
								author = $2,
								content = $3
								WHERE id = $4`,
		dbArticle.Title, dbArticle.Author, dbArticle.Content, dbArticle.Id)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
		return
	}
	resMsg.setSuccess("Article successfully updated.")
	json.NewEncoder(writer).Encode(resMsg)
}

// --------------------------------------------------------
// Tests
// --------------------------------------------------------

// Custom test1
func Test1(writer http.ResponseWriter, request *http.Request) {

	rows, err := DB.Query("SELECT * FROM users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.Id, &user.Login, &user.Password)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(user)
	}
}

// Custom test2
func Test2(writer http.ResponseWriter, request *http.Request) {
	result := userLoginExists("Bob")
	fmt.Println(result)
}

// Custom test3
func Test3(writer http.ResponseWriter, request *http.Request) {

}

func main() {

	fmt.Println("REST API V2.0 worked....")

	// open DB connect
	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		log.Fatal(err)
	}
	DB = db
	defer db.Close()

	/* init mux Router */
	router := mux.NewRouter().StrictSlash(true)

	/* User registration */
	router.HandleFunc("/register", Register).Methods("POST")

	/* User auth, get token */
	router.HandleFunc("/auth", GetToken).Methods("POST")

	/* Users show all */
	router.HandleFunc("/users", ShowUsers).Methods("GET")

	/* Articles show all */
	router.HandleFunc("/articles", ShowArticles).Methods("GET")

	/* show one */
	router.HandleFunc("/article/{id}", ShowArticleById).Methods("GET")

	/* create (with token) */
	router.Handle("/article", jwtMiddleware.Handler(http.HandlerFunc(CreateArticle))).Methods("POST")

	/* delete */
	router.HandleFunc("/article/{id}", DeleteArticle).Methods("DELETE")

	/* update */
	router.HandleFunc("/article/{id}", UpdateArticle).Methods("PUT")

	/* tests */
	router.HandleFunc("/test1", Test1).Methods("GET")
	router.HandleFunc("/test2", Test2).Methods("GET")
	router.HandleFunc("/test3", Test3).Methods("GET")

	/* serve */
	log.Fatal(http.ListenAndServe(":8050", router))
}
