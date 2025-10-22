package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/pflow/backend/internal/camunda"
	"github.com/example/pflow/backend/internal/config"
	"github.com/example/pflow/backend/internal/flow"
	httpserver "github.com/example/pflow/backend/internal/http"
	flowhttp "github.com/example/pflow/backend/internal/http/flow"
	workorderhttp "github.com/example/pflow/backend/internal/http/workorder"
	"github.com/example/pflow/backend/internal/mq"
	"github.com/example/pflow/backend/internal/persistence"
	"github.com/example/pflow/backend/internal/workorder"
)

func main() {
	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := persistence.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("close db: %v", err)
		}
	}()

	publisher, err := mq.NewPublisher(cfg.Queue)
	if err != nil {
		log.Printf("queue disabled: %v", err)
	}
	if publisher != nil {
		defer func() {
			if err := publisher.Close(); err != nil {
				log.Printf("close publisher: %v", err)
			}
		}()
	}

	camundaClient := camunda.NewClient(cfg.Camunda)
	runtime := camunda.NewRuntime(camundaClient.HTTP())

	flowRepo := flow.NewRepository(db.DB)
	flowService := flow.NewService(flowRepo, camundaClient, publisher)

	workorderRepo := workorder.NewRepository(db.DB)
	flowReader := workorder.FlowServiceAdapter{Service: flowService}
	workorderService := workorder.NewService(workorderRepo, flowReader, runtime, publisher)

	server := httpserver.NewServer(cfg,
		flowhttp.Handlers{Service: flowService},
		workorderhttp.Handlers{Service: workorderService},
	)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Run()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("server shutdown error: %v", err)
		}
	case err := <-errCh:
		if err != nil {
			log.Printf("server error: %v", err)
		}
	}

	log.Println("server stopped")
}
