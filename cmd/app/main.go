package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tombuente/apex/internal/accounting"
	"github.com/tombuente/apex/internal/logistics"
	"github.com/tombuente/apex/internal/static"
)

func main() {
	postgres, err := pgx.Connect(context.Background(), os.Getenv("DATABASE"))
	if err != nil {
		slog.Error("Unable to connect to Postgres:", "error", err)
		return
	}
	defer postgres.Close(context.Background())

	db := sqlx.MustConnect("sqlite3", "data.sqlite")

	_, err = db.ExecContext(context.Background(), logistics.Schema)
	if err != nil {
		fmt.Println("Unable to load logistics schema:", err)
		return
	}
	_, err = db.ExecContext(context.Background(), logistics.Fixture)
	if err != nil {
		fmt.Println("Unable to load logistics fixture:", err)
		return
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	logisticsDB := logistics.NewDatabase(db)
	logisticsService := logistics.NewService(logisticsDB)

	logisticsAPIRouter := logistics.NewAPIRouter(logisticsService)

	r.Route("/api", func(r chi.Router) {
		r.Mount("/logistics", logisticsAPIRouter)
	})

	logisticsUIRouter, err := logistics.NewUIRouter(logisticsService)
	if err != nil {
		fmt.Println("Unable to create UI router:", err)
		return
	}

	r.Mount("/logistics", logisticsUIRouter)

	accountingQueries := accounting.NewDatabase(postgres)
	accountingService := accounting.NewService(accountingQueries)
	accountingUIRouter, err := accounting.NewUIRouter(accountingService)
	if err != nil {
		fmt.Println("Unable to create UI router:", err)
		return
	}

	r.Mount("/accounting", accountingUIRouter)

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "static"))
	filesRouter, err := static.NewStaticRouter(filesDir)
	if err != nil {
		fmt.Println("Unable to setup static files: ", err)
	}
	r.Mount("/static", filesRouter)

	slog.Info("Running...")
	http.ListenAndServe(":8080", r)
}
