package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/stonik02/proxy_service/internal/auth"
	"github.com/stonik02/proxy_service/internal/config"
	"github.com/stonik02/proxy_service/internal/person"
	"github.com/stonik02/proxy_service/internal/roles"
	"github.com/stonik02/proxy_service/pkg/logging"
	"github.com/stonik02/proxy_service/pkg/logging/db/postgresql"
	_ "github.com/stonik02/proxy_service/pkg/logging/db/postgresql"
)

// TODO: Сделать возможность сразу присваивать несколько ролей
// TODO: Сделать jwt
// TODO: Сделать валидацию данных
// TODO: Сделать middleware на permission к разным ручкам

func main() {
	logger := logging.GetLogger()

	logger.Info("create router")
	router := httprouter.New()

	cfg := config.GetConfig()

	// init db postgreSQL
	dbClient, err := postgresql.NewClient(context.TODO(), 5, cfg.Storage)
	if err != nil {
		logger.Fatalf("fatal error: %s", err)
	}

	personRepository := person.NewRepository(dbClient, &logger)
	logger.Info("register person handler")
	personHandler := person.NewHandler(logger, personRepository)
	personHandler.Register(router)

	authRepository := auth.NewRepository(dbClient, &logger, personRepository)
	logger.Info("register auth handler")
	authHandler := auth.NewHandler(logger, authRepository)
	authHandler.Register(router)

	rolesRepository := roles.NewRepository(dbClient, &logger)
	logger.Info("register roles handler")
	rolesHandler := roles.NewHandler(logger, rolesRepository)
	rolesHandler.Register(router)

	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start application")

	listener, err := net.Listen("tcp", cfg.Listen.Port)
	if err != nil {
		panic(err)

	}

	fmt.Println(listener.Addr())

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Infof("server is listening on port : %s", cfg.Listen.Port)
	logger.Fatal(server.Serve(listener))
}
