package users

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/sovernut/bookstore_users-api/domain/users"
	"github.com/sovernut/bookstore_users-api/services"
	errors "github.com/sovernut/bookstore_users-api/utils/error"
)

func GetUser(c *gin.Context) {
	userId, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if userErr != nil {
		err := errors.NewBadRequestError("user id should be a number")
		c.JSON(err.Status, err)
		return
	}

	result, saveErr := services.GetUser(userId)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}
	c.JSON(http.StatusCreated, result)

}

func CreateUser(c *gin.Context) {
	var user users.User
	// bytes, err := ioutil.ReadAll(c.Request.Body)
	// if err != nil {
	// 	return
	// }
	// if err := json.Unmarshal(bytes, &user); err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		fmt.Println("Binding error", err)
		c.JSON(restErr.Status, restErr)
		return
	}
	result, saveErr := services.CreateUser(user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func SearchUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "not implement")

}
