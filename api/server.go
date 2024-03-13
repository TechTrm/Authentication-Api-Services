package api

import (
	db "github.com/TechTrm/Authentication-Api-Services/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our AuthApi service.
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(store *db.Store) *Server {

	server := &Server{
		store: store,
	}
	router := gin.Default()

	router.POST("/v1/api/users/register", server.createUser)
	// router.POST("/users/login", server.loginUser)
	// router.POST("/users/logout", server.logOutUser)

	// router.GET("/users/:id", server.getAccount)
	// router.GET("/users", server.listAccount)
	// router.PUT("/users", server.updateAccount)
	// router.DELETE("/users/:id", server.deleteAccount)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
