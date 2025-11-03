package http

import (
    "database/sql"
    "encoding/json"
    "errors"
    "net/http"
)

type AppError struct {
    Code   string `json:"code"`
    Status int    `json:"-"`
    Msg    string `json:"message"`
    Fields []FieldError `json:"fields,omitempty"`
}

func (e AppError) Error() string { return e.Msg }

func NewBadRequest(code, msg string) AppError  { return AppError{Code: code, Status: http.StatusBadRequest, Msg: msg} }
func NewNotFound(code, msg string) AppError    { return AppError{Code: code, Status: http.StatusNotFound, Msg: msg} }
func NewConflict(code, msg string) AppError    { return AppError{Code: code, Status: http.StatusConflict, Msg: msg} }
func NewInternal(code, msg string) AppError    { return AppError{Code: code, Status: http.StatusInternalServerError, Msg: msg} }

func writeError(w http.ResponseWriter, err error) {
    var app AppError
    switch {
    case errors.As(err, &app):
        respondErrorJSON(w, app.Status, app)
    case errors.Is(err, sql.ErrNoRows):
        respondErrorJSON(w, http.StatusNotFound, AppError{Code: "not_found", Msg: "recurso no encontrado"})
    default:
        respondErrorJSON(w, http.StatusInternalServerError, AppError{Code: "internal_error", Msg: err.Error()})
    }
}

func respondErrorJSON(w http.ResponseWriter, status int, e AppError) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(map[string]any{"error": e})
}

type FieldError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}
