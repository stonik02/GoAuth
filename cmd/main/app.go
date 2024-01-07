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
	person "github.com/stonik02/proxy_service/internal/persons"
	"github.com/stonik02/proxy_service/internal/roles"
	"github.com/stonik02/proxy_service/internal/token"
	utils "github.com/stonik02/proxy_service/internal/util/middleware"
	"github.com/stonik02/proxy_service/pkg/db/postgresql"
	"github.com/stonik02/proxy_service/pkg/logging"
	"github.com/stonik02/proxy_service/pkg/middleware"
)

// TODO: JWT сделал, теперь сделать refresh
// TODO: В middleware для update user сделать проверку, что это сам юзер или админ??
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
	utilsSQLClient := utils.NewPgClient(dbClient, &logger)
	utilsRepository := utils.NewRepository(&logger, utilsSQLClient)
	personRepository := person.NewRepository(&logger, personSQLClient)
	tokenRepository := token.NewRepository(&logger, *cfg)
	authRepository := auth.NewRepository(dbClient, &logger, personRepository, tokenRepository)
	rolesRepository := roles.NewRepository(&logger, rolesSQLClient)

	checkPermissionMiddleware := middleware.AuthorizedRoleMiddleware{
		UtilsRepository: *utilsRepository,
		TokenRepository: tokenRepository,
		Cfg:             *cfg,
		Logger:          &logger,
	}

	logger.Info("register person handler")
	personHandler := person.NewHandler(logger, personRepository, checkPermissionMiddleware)
	personHandler.Register(router)

	logger.Info("register auth handler")
	authHandler := auth.NewHandler(logger, authRepository, checkPermissionMiddleware)
	authHandler.Register(router)

	logger.Info("register roles handler")
	rolesHandler := roles.NewHandler(logger, rolesRepository, checkPermissionMiddleware)
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
