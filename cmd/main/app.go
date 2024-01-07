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
	"github.com/stonik02/proxy_service/internal/persons"
	"github.com/stonik02/proxy_service/internal/roles"
	"github.com/stonik02/proxy_service/internal/token"
	"github.com/stonik02/proxy_service/pkg/db/postgresql"
	"github.com/stonik02/proxy_service/pkg/logging"
	"github.com/stonik02/proxy_service/pkg/middleware"
)

// TODO: JWT сделал, теперь сделать refresh и привязать к auth
// TODO: Сделать middleware на permission к разным ручкам
// TODO: Сделать возможность сразу присваивать несколько ролей
// TODO: Сделать валидацию данных

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

	rolesSQLClient := roles.NewPgClient(dbClient, &logger)
	personSQLClient := person.NewPgClient(dbClient, &logger)
	personRepository := person.NewRepository(&logger, personSQLClient)
	authRepository := auth.NewRepository(dbClient, &logger, personRepository)
	rolesRepository := roles.NewRepository(&logger, rolesSQLClient)
	tokenRepository := token.NewRepository(&logger, *cfg)

	checkPermissionMiddleware := middleware.AuthorizedRoleMiddleware{
		RolesRepository: rolesRepository,
		TokenRepository: tokenRepository,
	}

	logger.Info("register person handler")
	personHandler := person.NewHandler(logger, personRepository)
	personHandler.Register(router)

	logger.Info("register auth handler")
	authHandler := auth.NewHandler(logger, authRepository, checkPermissionMiddleware)
	authHandler.Register(router)

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

	fmt.Printf("cfg.Listen.Port = %s", cfg.Listen.Port)

	fmt.Println(listener.Addr())

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Infof("server is listening on port : %s", cfg.Listen.Port)
	logger.Fatal(server.Serve(listener))
}
