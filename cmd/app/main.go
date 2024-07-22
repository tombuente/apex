package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/tombuente/apex"
	"github.com/tombuente/apex/internal/accounting"
	"github.com/tombuente/apex/internal/logistics"
)

func main() {
	// "postgres://username:password@host:port/database_name"
	postgres, err := pgx.Connect(context.Background(), os.Getenv("DATABASE"))
	if err != nil {
		slog.Error("Unable to connect to Postgres:", "error", err)
		return
	}
	defer postgres.Close(context.Background())

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	logisticsDB := logistics.MakeDatabase(postgres)
	logisticsService := logistics.MakeService(logisticsDB)
	logisticsUIRouter, err := logistics.NewUIRouter(apex.TemplatesFS, logisticsService)
	if err != nil {
		slog.Error("Unable to create logistics UI router", "error", err)
		return
	}
	r.Mount("/logistics", logisticsUIRouter)

	accountingDB := accounting.MakeDatabase(postgres)
	accountingService := accounting.MakeService(accountingDB)
	accountingUIRouter, err := accounting.NewUIRouter(apex.TemplatesFS, accountingService)
	if err != nil {
		slog.Error("Unable to create accounting UI router", "error", err)
		return
	}
	r.Mount("/accounting", accountingUIRouter)

	staticHandler := http.FileServer(http.FS(apex.StaticFS))
	r.Handle("/static/*", staticHandler)

	slog.Info("Running...")
	http.ListenAndServe(":8080", r)
}
