package users

import (
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
		return
	}

	result, saveErr := services.UserServices.CreateUser(user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}
	c.JSON(http.StatusCreated, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func Get(c *gin.Context) {
	// before get resources we need to validate authentication
	if restErr := oauth.AuthenticateRequest(c.Request); restErr != nil {
		c.JSON(restErr.Status, restErr)
		return
	}
	//caller is not authorized to access this resource
	if callerId := oauth.GetCallerId(c.Request); callerId == 0 {
		restErr := rest_errors.NewUnauthorizedError("can not get callerId")
		c.JSON(restErr.Status, restErr)
		return
	}

	userId, restErr := getUserId(c.Param("user_id"))
	if restErr != nil {
		c.JSON(restErr.Status, restErr)
		return
	}
	user, restErr := services.UserServices.GetUser(userId)
	if restErr != nil {
		c.JSON(restErr.Status, restErr)
		return
	}

	//check if current access token owner is equal to passed user_id param
	//means is logined in user retrieve his own info, we return full user info
	if oauth.GetCallerId(c.Request) == user.Id {
		c.JSON(http.StatusOK, user.Marshall(false))
		return
	}
	//otherwise, we check is public request or private request
	//return full info only when is private internal call
	c.JSON(http.StatusOK, user.Marshall(oauth.IsPublic(c.Request)))
}

func Update(c *gin.Context) {
	userId, restErr := getUserId(c.Param("user_id"))
	if restErr != nil {
		c.JSON(restErr.Status, restErr)
		return
	}
	//validate json body pass in
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	user.Id = userId
	isPartial := c.Request.Method == http.MethodPatch

	result, restErr := services.UserServices.UpdateUser(isPartial, user)
	if restErr != nil {
		c.JSON(restErr.Status, restErr)
		return
	}
	c.JSON(http.StatusOK, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func Delete(c *gin.Context) {
	userId, restErr := getUserId(c.Param("user_id"))
	if restErr != nil {
		c.JSON(restErr.Status, restErr)
		return
	}
	if restErr := services.UserServices.DeleteUser(userId); restErr != nil {
		c.JSON(restErr.Status, restErr)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func Search(c *gin.Context) {
	status := c.Query("status")
	users, restErr := services.UserServices.Search(status)
	if restErr != nil {
		c.JSON(restErr.Status, restErr)
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
		return
	}
	user, restErr := services.UserServices.LoginUser(request)
	if restErr != nil {
		c.JSON(restErr.Status, restErr)
		return
	}
	c.JSON(http.StatusOK, user.Marshall(false))
}

func getUserId(userIdParam string) (int64, *rest_errors.RestErr) {
	//convert string to decimal
	userId, err := strconv.ParseInt(userIdParam, 10, 64)
	if err != nil {
		logger.Error(err.Error())
		return 0, rest_errors.NewBadRequestError("user id should be a number")
	}
	return userId, nil
}
