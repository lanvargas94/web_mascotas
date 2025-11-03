package main

import (
    "context"
    "log"
    "net/http"
    "time"

    "mascotas/internal/config"
    "mascotas/internal/database"
    httphandlers "mascotas/internal/http"
)

func main() {
    cfg := config.Load()
    port := cfg.Port
    dsn := cfg.DB_DSN

    ctx := context.Background()

    db, err := database.Open(dsn)
    if err != nil {
        log.Fatalf("db open error: %v", err)
    }
    defer db.Close()

    if err := database.Migrate(ctx, db); err != nil {
        log.Fatalf("db migrate error: %v", err)
    }

    // Build a simple handler with inlined CORS, logging and recovery.
    handler := httphandlers.NewRouter(db)

    srv := &http.Server{
        Addr:              ":" + port,
        Handler:           handler,
        ReadHeaderTimeout: 10 * time.Second,
        ReadTimeout:       15 * time.Second,
        WriteTimeout:      15 * time.Second,
        IdleTimeout:       60 * time.Second,
    }

    log.Printf("backend listening on :%s", port)
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("server error: %v", err)
    }
}
