package controller

import (
	"github.com/gin-gonic/gin"
)

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var users = map[string]User{}

func RegisterUser(c *gin.Context) {

}
