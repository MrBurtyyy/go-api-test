package server

import (
	"context"
	"flag"
	"fmt"
	"github.com/MrBurtyyy/go-api-test/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/yaml.v2"
	"net/http"
	"os"
)

var (
	cfgFile = flag.String("config", "etc/config.yaml", "path to configuration file")
	r       *chi.Mux
)

func Init() {
	if err := init2(); err != nil {
		fmt.Printf("failed to initialise server: %v\n", err)
	}
}

func init2() error {
	flag.Parse()
	cfgRaw, err := os.ReadFile(*cfgFile)
	if err != nil {
		return err
	}

	cfg := new(Config)
	if err = yaml.Unmarshal(cfgRaw, cfg); err != nil {
		return err
	}

	r = chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	dbCfg, err := pgxpool.ParseConfig(cfg.Postgres.ConnectionString)
	if err != nil {
		return err
	}
	db, err := pgxpool.NewWithConfig(context.Background(), dbCfg)
	if err != nil {
		return err
	}

	handler.Init(r, db)

	return nil
}

func ListenAndServe(addr string) {
	if err := http.ListenAndServe(addr, r); err != nil {
		fmt.Printf("failed to start HTTP server using '%v': %v", addr, err)
		os.Exit(1)
	}
}
