package api

import (
	"fmt"

	db "github.com/TechTrm/Authentication-Api-Services/db/sqlc"
	"github.com/TechTrm/Authentication-Api-Services/token"
	"github.com/TechTrm/Authentication-Api-Services/util"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our AuthApi service.
type Server struct {
	config  util.Config
	store  *db.Store
	tokenMaker token.Maker
	router *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, store *db.Store) (*Server, error)  {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)

	if err != nil {
		return nil, fmt.Errorf("cannot create token maker : %w", err)
    }

	server := &Server{
		config: config,
		store: store,
		tokenMaker: tokenMaker,
	}
	
	server.setUpRouter()
	return server, nil
}

func (server *Server) setUpRouter() {
	router := gin.Default()

	router.POST("/v1/api/users/register", server.createUser)
	router.POST("/v1/api/users/login", server.loginUser)
	router.POST("/v1/api/tokens/refresh_token", server.renewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.GET("/v1/api/users", server.listUsers)
	authRoutes.POST("/v1/api/users/logout", server.Logout)

	// router.GET("/users/:id", server.getAccount)
	// router.GET("/users", server.listAccount)
	// router.PUT("/users", server.updateAccount)
	// router.DELETE("/users/:id", server.deleteAccount)

	server.router = router
}
	

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
