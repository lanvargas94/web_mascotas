package http

import (
    "database/sql"
    "log"
    "net/http"
    "os"
    "strings"
)

// NewRouter builds the application's HTTP handler with a simple
// ServeMux and inlined cross-cutting concerns (CORS, logging, recovery)
// to keep things straightforward and student-friendly.
func NewRouter(db *sql.DB) http.Handler {
    h := NewHandlers(db)
    mux := http.NewServeMux()

    mux.HandleFunc("/health", h.Health)
    mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
        if err := db.Ping(); err != nil {
            http.Error(w, "not ready", http.StatusServiceUnavailable)
            return
        }
        w.WriteHeader(http.StatusNoContent)
    })

    // Mascotas collection
    mux.HandleFunc("/mascotas", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            h.ListMascotas(w, r)
        case http.MethodPost:
            h.CreateMascota(w, r)
        default:
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })

    // Mascotas item and nested cuidados
    mux.HandleFunc("/mascotas/", func(w http.ResponseWriter, r *http.Request) {
        path := r.URL.Path
        // /mascotas/{id} or /mascotas/{id}/cuidados
        if hasSuffix(path, "/cuidados") || hasSegment(path, "/cuidados/") {
            switch r.Method {
            case http.MethodGet:
                h.ListCuidadosByMascota(w, r)
            case http.MethodPost:
                h.CreateCuidadoForMascota(w, r)
            default:
                http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
            }
            return
        }
        switch r.Method {
        case http.MethodGet:
            h.GetMascota(w, r)
        case http.MethodPut:
            h.UpdateMascota(w, r)
        case http.MethodDelete:
            h.DeleteMascota(w, r)
        default:
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })

    // Cuidados by id
    mux.HandleFunc("/cuidados/", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            h.GetCuidado(w, r)
        case http.MethodPut:
            h.UpdateCuidado(w, r)
        case http.MethodDelete:
            h.DeleteCuidado(w, r)
        default:
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })

    // Wrap mux with simple handler that adds CORS, logging and recovery.
    return &appHandler{mux: mux}
}

// appHandler is a minimal wrapper that applies CORS headers,
// recovers from panics and logs requests.
type appHandler struct {
    mux *http.ServeMux
}

func (a *appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // CORS (very permissive for dev; restrict via ALLOWED_ORIGINS)
    origin := r.Header.Get("Origin")
    allowed := os.Getenv("ALLOWED_ORIGINS")
    if allowed == "" {
        if origin == "" {
            origin = "*"
        }
        w.Header().Set("Access-Control-Allow-Origin", origin)
    } else {
        ok := false
        for _, o := range strings.Split(allowed, ",") {
            if strings.TrimSpace(o) == origin {
                ok = true
                break
            }
        }
        if ok {
            w.Header().Set("Access-Control-Allow-Origin", origin)
        }
    }
    w.Header().Set("Vary", "Origin")
    w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    w.Header().Set("Access-Control-Allow-Credentials", "false")
    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusNoContent)
        return
    }

    // Recovery
    defer func() {
        if rec := recover(); rec != nil {
            log.Printf("panic: %v", rec)
            writeError(w, NewInternal("panic", "internal server error"))
        }
    }()

    // Logging (capture status)
    lrw := &loggingResponseWriter{ResponseWriter: w, status: http.StatusOK}
    a.mux.ServeHTTP(lrw, r)
    log.Printf("%s %s -> %d", r.Method, r.URL.Path, lrw.status)
}

type loggingResponseWriter struct {
    http.ResponseWriter
    status int
}

func (l *loggingResponseWriter) WriteHeader(code int) {
    l.status = code
    l.ResponseWriter.WriteHeader(code)
}

func hasSuffix(s, suf string) bool {
    if len(s) < len(suf) {
        return false
    }
    return s[len(s)-len(suf):] == suf
}

func hasSegment(s, seg string) bool {
    for i := 0; i+len(seg) <= len(s); i++ {
        if s[i:i+len(seg)] == seg {
            return true
        }
    }
    return false
}
