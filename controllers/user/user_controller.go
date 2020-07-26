package user

import (
	"net/http"
	"strconv"

	"github.com/olmuz/bookstore_users-api/domain/users"
	"github.com/olmuz/bookstore_users-api/services"
	"github.com/olmuz/bookstore_users-api/utils/errors"

	"github.com/gin-gonic/gin"
)

func getUserID(userIDparam string) (int64, *errors.RestErr) {
	userID, err := strconv.ParseInt(userIDparam, 10, 64)
	if err != nil {
		err := errors.NewBadRequestError("user_id has to be a type of number")
		return 0, err
	}
	return userID, nil
}

func Create(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, &restErr)
		return
	}

	result, restErr := services.CreateUser(user)
	if restErr != nil {
		c.JSON(restErr.Status, restErr)
		return
	}
	c.JSON(http.StatusCreated, result.Marshal(c.GetHeader("X-Public") == "true"))
}

func Get(c *gin.Context) {
	userID, idErr := getUserID(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}

	user, getErr := services.GetUser(userID)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}
	c.JSON(http.StatusOK, user.Marshal(c.GetHeader("X-Public") == "true"))
}

func Update(c *gin.Context) {
	userID, idErr := getUserID(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}

	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, &restErr)
		return
	}

	user.ID = userID

	partial := c.Request.Method == http.MethodPatch
	updatedUser, updateErr := services.UpdateUser(partial, user)
	if updateErr != nil {
		c.JSON(updateErr.Status, updateErr)
	}
	c.JSON(http.StatusOK, updatedUser.Marshal(c.GetHeader("X-Public") == "true"))
	return
}

func Delete(c *gin.Context) {
	userID, idErr := getUserID(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}
	if deleteErr := services.DeleteUser(userID); deleteErr != nil {
		c.JSON(deleteErr.Status, deleteErr)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
	return
}

func Search(c *gin.Context) {
	status := c.Query("status")
	users, err := services.SearchUser(status)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, users.Marshal(c.GetHeader("X-Public") == "true"))
	return
}
