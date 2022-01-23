package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/garugaru/knowledge/server/api"
	"github.com/garugaru/knowledge/server/conf"
	"github.com/garugaru/knowledge/server/data"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var (
		gracefulTimeout time.Duration
		confPath        string
		profile         bool
	)

	flag.DurationVar(&gracefulTimeout, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.StringVar(&confPath, "config", "config.yml", "configuration path default ./config.yml")
	flag.BoolVar(&profile, "profile", false, "enable pprof server")
	flag.Parse()

	ctx := context.Background()

	config, err := conf.FromYaml(confPath)
	if err != nil {
		log.Fatal(err)
	}

	catalog, err := createCatalog(config.Catalog)
	if err != nil {
		log.Fatal(err)
	}

	if err := catalog.Init(); err != nil {
		log.Fatal(err)
	}

	apiService := api.New(api.Config{
		EnableMetrics: true,
	}, catalog)

	apiServer := apiService.Server(api.ServeOpts{})

	go func() {
		if err := apiServer.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	if profile {
		serveProfiler()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	<-c

	log.Println("shutting down")

	ctx, cancel := context.WithTimeout(ctx, gracefulTimeout)
	defer cancel()

	if err := apiServer.Shutdown(ctx); err != nil {
		logrus.Warn(err)
	}
}

func createCatalog(catalog conf.Catalog) (data.Catalog, error) {
	db, err := createDB(catalog.Database)
	if err != nil {
		return nil, err
	}
	return data.NewDBCatalog(db), nil
}

func createDB(database conf.Database) (*gorm.DB, error) {
	switch database.Type {
	case conf.DatabaseTypeMySql:
		dsn, present := database.Params["dsn"]
		if !present {
			return nil, errors.New("dsn parameter must be defined using mysql database")
		}
		return gorm.Open(mysql.Open(dsn.(string)), &gorm.Config{})
	case conf.DatabaseTypePostgres:
		dsn, present := database.Params["dsn"]
		if !present {
			return nil, errors.New("dsn parameter must be defined using postgres database")
		}
		return gorm.Open(postgres.Open(dsn.(string)), &gorm.Config{})
	case conf.DatabaseTypeSQLite:
		dbPath, present := database.Params["db_path"]
		if !present {
			return nil, errors.New("db_path parameter must be defined using sqlite database")
		}
		return gorm.Open(sqlite.Open(dbPath.(string)), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unknown database type %s", database.Type)
	}
}

func serveProfiler() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}
