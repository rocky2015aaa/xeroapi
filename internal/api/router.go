package api

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/shmoulana/xeroapi/internal/api/handlers"
)

func NewRouter(handler *handlers.Handler) http.Handler {
	router := gin.Default()
	// Get the absolute path to the templates directory
	absPath, err := filepath.Abs("assets/templates/*")
	if err != nil {
		panic(err)
	}

	// Load the templates using the absolute path
	router.LoadHTMLGlob(absPath)

	router.GET("/", handler.IndexHandler)
	router.GET("/ping", handler.Ping)
	router.GET("/login", handler.LoginHandler)
	router.GET("/refresh", handler.HandleTokenRefreshRequest)
	router.GET(os.Getenv(EnvXeroCallbackAPI), handler.HandleOAuthCallback)

	v1 := router.Group("/api/v1")
	v1.POST("/items", handler.CreateItems)
	v1.POST("/contacts", handler.CreateContacts)
	v1.POST("/tracking_catergory", handler.CreateTrackingCategory)
	v1.POST("/:tracking_category_ID/option", handler.CreateTrackingOption)
	v1.POST("/account", handler.CreateAccount)

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"message": "Method Not Allowed",
		})
	})

	return router
}
