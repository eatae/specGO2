package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"specGo2/examples/crudV7/models"
)

// GetUsers
//
// method:GET
func GetUsers(context *gin.Context) {
	var users []models.User
	err := models.GetAllUsers(&users)
	if err != nil {
		context.AbortWithStatus(http.StatusNotFound)
	} else {
		context.JSON(http.StatusOK, users)
	}
}

// GetUserById
//
// method:GET
// param: id
func GetUserById(context *gin.Context) {
	id := context.Params.ByName("id")
	var user models.User
	err := models.GetUserById(&user)
	if err != nil {
		context.AbortWithStatus(http.StatusNotFound)
	} else {
		context.JSON(http.StatusOK, user)
	}
}

// CreateUser
//
// method:POST
func CreateUser(context *gin.Context) {
	var user models.User
	context.BindJSON(&user)
	err := models.CreateUser(&user)
	if err != nil {
		fmt.Println(err.Error())
		context.AbortWithStatus(http.StatusNotFound)
	} else {
		context.JSON(http.StatusCreated, user)
	}
}

// UpdateUser
//
// method:PUT
func UpdateUser(context *gin.Context) {

}

// DeleteUser
//
// method:DELETE
// param: id
func DeleteUser(context *gin.Context) {

}
