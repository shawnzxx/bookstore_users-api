package users

import (
	"encoding/json"
	"github.com/shawnzxx/bookstore_utils-go/app_logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shawnzxx/bookstore_oauth-go/oauth"
	"github.com/shawnzxx/bookstore_users-api/domain/users"
	"github.com/shawnzxx/bookstore_users-api/services"
	"github.com/shawnzxx/bookstore_utils-go/rest_errors"
)

var (
	logger = app_logger.GetLogger()
)

func Create(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		logError(restErr)
		return
	}

	result, saveErr := services.UserServices.CreateUser(user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		logError(saveErr)
		return
	}
	c.JSON(http.StatusCreated, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func Get(c *gin.Context) {
	// before get resources we need to validate authentication
	if err := oauth.AuthenticateRequest(c.Request); err != nil {
		c.JSON(err.Status, err)
		logError(err)
		return
	}
	//caller is not authorized to access this resource
	if callerId := oauth.GetCallerId(c.Request); callerId == 0 {
		err := rest_errors.NewUnauthorizedError("can not get callerId")
		c.JSON(err.Status, err)
		logError(err)
		return
	}

	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		logError(idErr)
		return
	}
	user, getErr := services.UserServices.GetUser(userId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		logError(getErr)
		return
	}

	//check if current access token owner is equal to passed user_id param
	//means is logined in user retrieve his own info, we return full user info
	if oauth.GetCallerId(c.Request) == user.Id {
		c.JSON(http.StatusOK, user.Marshall(false))
		return
	}
	//otherwise we check is public request or private request
	//return full info only when is private internal call
	c.JSON(http.StatusOK, user.Marshall(oauth.IsPublic(c.Request)))
}

func Update(c *gin.Context) {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		logError(idErr)
		return
	}
	//validate json body pass in
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		logError(restErr)
		return
	}

	user.Id = userId
	isPartial := c.Request.Method == http.MethodPatch

	result, updateErr := services.UserServices.UpdateUser(isPartial, user)
	if updateErr != nil {
		c.JSON(updateErr.Status, updateErr)
		logError(updateErr)
		return
	}
	c.JSON(http.StatusOK, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func Delete(c *gin.Context) {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		logError(idErr)
		return
	}
	if deleteErr := services.UserServices.DeleteUser(userId); deleteErr != nil {
		c.JSON(deleteErr.Status, deleteErr)
		logError(deleteErr)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func Search(c *gin.Context) {
	status := c.Query("status")
	users, err := services.UserServices.Search(status)
	if err != nil {
		c.JSON(err.Status, err)
		logError(err)
		return
	}
	result := users.Marshall(c.GetHeader("X-Public") == "true")
	c.JSON(http.StatusOK, result)
}

func Login(c *gin.Context) {
	var request users.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		logError(restErr)
		return
	}
	user, err := services.UserServices.LoginUser(request)
	if err != nil {
		c.JSON(err.Status, err)
		logError(err)
		return
	}
	c.JSON(http.StatusOK, user.Marshall(false))
}

func getUserId(userIdParam string) (int64, *rest_errors.RestErr) {
	//convert string to decimal
	userId, err := strconv.ParseInt(userIdParam, 10, 64)
	if err != nil {
		logError(err)
		return 0, rest_errors.NewBadRequestError("user id should be a number")
	}
	return userId, nil
}

func logError(err interface{}) {
	errByte, _ := json.Marshal(err)
	logger.Error("bookstore_users-api error %v", string(errByte))
}
