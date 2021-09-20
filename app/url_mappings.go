package app

import (
	"github.com/shawnzxx/bookstore_users-api/controllers/ping"
	"github.com/shawnzxx/bookstore_users-api/controllers/users"
)

func mapUrls() {
	_router.GET("/ping", ping.Ping)

	_router.POST("/users", users.Create)
	_router.GET("/users/:user_id", users.Get)
	_router.PUT("/users/:user_id", users.Update)
	_router.PATCH("/users/:user_id", users.Update)
	_router.DELETE("/users/:user_id", users.Delete)
	_router.POST("/users/login", users.Login)

	_router.GET("/internal/users/search", users.Search)
}
