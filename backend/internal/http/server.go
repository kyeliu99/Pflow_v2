package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/example/pflow/backend/internal/config"
	flowhttp "github.com/example/pflow/backend/internal/http/flow"
	workorderhttp "github.com/example/pflow/backend/internal/http/workorder"
)

type Server struct {
	engine *gin.Engine
	http   *http.Server
}

func NewServer(cfg config.Config, flowHandlers flowhttp.Handlers, workorderHandlers workorderhttp.Handlers) *Server {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	api := engine.Group("/api")
	{
		flow := api.Group("/flows")
		flow.GET("", flowHandlers.List)
		flow.POST("", flowHandlers.Create)
		flow.GET(":id", flowHandlers.Get)
		flow.PUT(":id", flowHandlers.Update)

		workorders := api.Group("/workorders")
		workorders.GET("", workorderHandlers.List)
		workorders.POST("", workorderHandlers.Create)
		workorders.GET(":id", workorderHandlers.Get)
		workorders.POST(":id/retry", workorderHandlers.Retry)
	}

	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port),
		Handler:      engine,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}

	return &Server{engine: engine, http: httpServer}
}

func (s *Server) Run() error {
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
