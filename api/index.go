package handler

import (
	"net/http"
	"os"
	"sync"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/i18n"
	"github.com/QuantumNous/new-api/logger"
	"github.com/QuantumNous/new-api/middleware"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/oauth"
	perfmetrics "github.com/QuantumNous/new-api/pkg/perf_metrics"
	"github.com/QuantumNous/new-api/router"
	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/service/authz"
	"github.com/QuantumNous/new-api/setting/ratio_setting"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var (
	appOnce sync.Once
	app     *gin.Engine
)

func initApp() {
	appOnce.Do(func() {
		common.InitEnv()
		logger.SetupLogger()
		ratio_setting.InitRatioSettings()
		service.InitHttpClient()
		service.InitTokenEncoders()

		if err := model.InitDB(); err != nil {
			common.FatalLog("failed to initialize database: " + err.Error())
			return
		}
		if err := authz.Init(model.DB); err != nil {
			common.FatalLog("failed to initialize authorization: " + err.Error())
			return
		}
		model.CheckSetup()
		model.InitOptionMap()
		common.CleanupOldCacheFiles()
		model.GetPricing()
		if err := model.InitLogDB(); err != nil {
			common.FatalLog("failed to initialize log database: " + err.Error())
			return
		}
		if err := common.InitRedisClient(); err != nil {
			common.SysError("redis init failed: " + err.Error())
		}
		perfmetrics.Init()
		common.StartSystemMonitor()

		if err := i18n.Init(); err != nil {
			common.SysError("i18n init failed: " + err.Error())
		}
		i18n.SetUserLangLoader(model.GetUserLanguage)

		if err := oauth.LoadCustomProviders(); err != nil {
			common.SysError("oauth load failed: " + err.Error())
		}

		gin.SetMode(gin.ReleaseMode)
		app = gin.New()
		app.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}))
		app.Use(middleware.RequestId())
		app.Use(middleware.PoweredBy())
		app.Use(middleware.I18n())
		middleware.SetUpLogger(app)

		store := cookie.NewStore([]byte(common.SessionSecret))
		store.Options(sessions.Options{
			Path:     "/",
			MaxAge:   2592000,
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
		})
		app.Use(sessions.Sessions("session", store))

		// Skip embedded frontend - Vercel serves static files separately
		os.Setenv("FRONTEND_BASE_URL", "/")
		router.SetRouter(app, router.ThemeAssets{})
	})
}

func Handler(w http.ResponseWriter, r *http.Request) {
	initApp()
	app.ServeHTTP(w, r)
}
