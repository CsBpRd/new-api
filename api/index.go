package main

import (
	"net/http"
	"os"
	"sync"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/controller"
	"github.com/QuantumNous/new-api/middleware"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/router"
	"github.com/gin-gonic/gin"
)

var (
	appOnce sync.Once
	app     *gin.Engine
)

func initApp() {
	appOnce.Do(func() {
		// Initialize environment
		common.InitEnv()

		// Initialize database
		model.InitDB()

		// Initialize disk data
		common.InitDiskData()

		// Initialize OAuth
		controller.InitOAuth()

		// Initialize Casbin
		middleware.InitCasbinEnforcer()

		// Create Gin engine
		gin.SetMode(gin.ReleaseMode)
		app = gin.New()

		// Setup router
		router.SetRouter(app, nil)
	})
}

// Handler is the entry point for Vercel serverless function
func Handler(w http.ResponseWriter, r *http.Request) {
	initApp()
	app.ServeHTTP(w, r)
}
