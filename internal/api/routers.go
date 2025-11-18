package api

import (
	"example-otelgo-miniservice/internal/api/dto"
	"example-otelgo-miniservice/internal/api/middleware"
	"example-otelgo-miniservice/internal/application/service"
	"example-otelgo-miniservice/internal/config"
	"example-otelgo-miniservice/internal/port"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router struct {
	jwe    *service.JWE
	engine *gin.Engine
	routes struct {
		authRoute *AuthRoute
	}
}

func NewRouter(
	jwe *service.JWE,
	// Routes
	authRouter *AuthRoute,
) (port.Router, error) {
	gin.SetMode(gin.ReleaseMode)

	r := &Router{
		jwe:    jwe,
		engine: gin.Default(),
	}

	// Routes
	r.routes.authRoute = authRouter

	if err := r.RegisterRoutes(); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Router) Start() error {

	proxies := []string{
		"127.0.0.1",
	}

	if err := r.engine.SetTrustedProxies(proxies); err != nil {
		return err
	}

	return r.engine.Run(fmt.Sprintf(":%s", config.GetEnv("APP_PORT", "8080")))
}

func (r *Router) RegisterRoutes() error {
	r.engine.Use(gin.Recovery())
	r.engine.Use(middleware.OTELMiddleware())

	r.engine.GET("/ping", func(c *gin.Context) {
		RespJsonSuccess(c, http.StatusOK, nil, "pong", "0")
	})

	api := r.engine.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			_ = v1
		}

		auth := api.Group("/auth")
		{
			auth.POST("/login", r.routes.authRoute.Login)
			auth.POST("/logout", r.routes.authRoute.Logout)
			auth.POST("/refresh", r.routes.authRoute.RefreshToken)
			auth.POST("/profile", r.routes.authRoute.Profile)
		}
	}

	r.engine.NoRoute(func(c *gin.Context) {
		RespJsonSuccess(c, http.StatusNotFound, nil, dto.MsgPageNotFound, "0")
	})

	return nil
}

func RespJsonSuccess(c *gin.Context, statusCode int, result any, msg, code string) {
	resp := gin.H{
		"code":    code,
		"message": msg,
	}

	if result != nil {
		resp["result"] = result
	}

	c.JSON(statusCode, resp)
}

func RespJsonError(c *gin.Context, statusCode int, err, code string) {
	c.JSON(statusCode, gin.H{
		"code":  code,
		"error": err,
	})
}
