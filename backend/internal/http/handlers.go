package http

import (
    "encoding/json"
    "errors"
    "net/http"
    "net/url"
    "strconv"
    "strings"
    "time"

    "mascotas/internal/models"
    "github.com/go-playground/validator/v10"
    "database/sql"
)

type Handlers struct {
    DB           *sql.DB
    validate     *validator.Validate
}

func NewHandlers(db *sql.DB) *Handlers {
    return &Handlers{
        DB:       db,
        validate: validator.New(),
    }
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
    respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Mascotas

func (h *Handlers) ListMascotas(w http.ResponseWriter, r *http.Request) {
    limit, offset, appErr := parsePagination(r.URL.Query())
    if appErr != nil {
        writeError(w, appErr)
        return
    }
    list, err := (models.MascotaStore{DB: h.DB}).ListPaged(r.Context(), limit, offset)
    if err != nil {
        writeError(w, err)
        return
    }
    respondJSON(w, http.StatusOK, list)
}

func (h *Handlers) CreateMascota(w http.ResponseWriter, r *http.Request) {
    var in struct {
        Nombre          string `json:"nombre" validate:"required,min=2,max=100"`
        Especie         string `json:"especie" validate:"required,oneof=Perro Gato Conejo"`
        Raza            string `json:"raza" validate:"required,min=2,max=100"`
        FechaNacimiento string `json:"fecha_nacimiento" validate:"required,datetime=2006-01-02"`
        Sexo            string `json:"sexo" validate:"required,oneof=Macho Hembra"`
    }
    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        writeError(w, NewBadRequest("invalid_json", "JSON inválido"))
        return
    }
    if err := h.validate.Struct(in); err != nil {
        if verrs, ok := err.(validator.ValidationErrors); ok {
            writeError(w, AppError{Code: "validation_error", Status: http.StatusBadRequest, Msg: "Datos inválidos", Fields: mapFieldErrors(verrs)})
        } else {
            writeError(w, NewBadRequest("validation_error", "Datos inválidos"))
        }
        return
    }
    dob, err := parseDate(in.FechaNacimiento)
    if err != nil {
        writeError(w, NewBadRequest("invalid_date", "fecha_nacimiento debe ser YYYY-MM-DD"))
        return
    }
    m := &models.Mascota{Nombre: in.Nombre, Especie: in.Especie, Raza: in.Raza, FechaNacimiento: dob, Sexo: in.Sexo}
    if err := (models.MascotaStore{DB: h.DB}).Create(r.Context(), m); err != nil {
        writeError(w, err)
        return
    }
    respondJSON(w, http.StatusCreated, m)
}

func (h *Handlers) GetMascota(w http.ResponseWriter, r *http.Request) {
    id, err := idFromPath(r.URL.Path)
    if err != nil {
        http.Error(w, "invalid id", http.StatusBadRequest)
        return
    }
    m, err := (models.MascotaStore{DB: h.DB}).Get(r.Context(), id)
    if err != nil {
        writeError(w, err)
        return
    }
    respondJSON(w, http.StatusOK, m)
}

func (h *Handlers) UpdateMascota(w http.ResponseWriter, r *http.Request) {
    id, err := idFromPath(r.URL.Path)
    if err != nil {
        http.Error(w, "invalid id", http.StatusBadRequest)
        return
    }
    var in struct {
        Nombre          string `json:"nombre" validate:"required,min=2,max=100"`
        Especie         string `json:"especie" validate:"required,oneof=Perro Gato Conejo"`
        Raza            string `json:"raza" validate:"required,min=2,max=100"`
        FechaNacimiento string `json:"fecha_nacimiento" validate:"required,datetime=2006-01-02"`
        Sexo            string `json:"sexo" validate:"required,oneof=Macho Hembra"`
    }
    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        writeError(w, NewBadRequest("invalid_json", "JSON inválido"))
        return
    }
    if err := h.validate.Struct(in); err != nil {
        if verrs, ok := err.(validator.ValidationErrors); ok {
            writeError(w, AppError{Code: "validation_error", Status: http.StatusBadRequest, Msg: "Datos inválidos", Fields: mapFieldErrors(verrs)})
        } else {
            writeError(w, NewBadRequest("validation_error", "Datos inválidos"))
        }
        return
    }
    dob, err := parseDate(in.FechaNacimiento)
    if err != nil {
        writeError(w, NewBadRequest("invalid_date", "fecha_nacimiento debe ser YYYY-MM-DD"))
        return
    }
    m := &models.Mascota{ID: id, Nombre: in.Nombre, Especie: in.Especie, Raza: in.Raza, FechaNacimiento: dob, Sexo: in.Sexo}
    if err := (models.MascotaStore{DB: h.DB}).Update(r.Context(), m); err != nil {
        writeError(w, err)
        return
    }
    respondJSON(w, http.StatusOK, m)
}

func (h *Handlers) DeleteMascota(w http.ResponseWriter, r *http.Request) {
    id, err := idFromPath(r.URL.Path)
    if err != nil {
        http.Error(w, "invalid id", http.StatusBadRequest)
        return
    }
    if err := (models.MascotaStore{DB: h.DB}).Delete(r.Context(), id); err != nil {
        writeError(w, err)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

// Cuidados

func (h *Handlers) ListCuidadosByMascota(w http.ResponseWriter, r *http.Request) {
    mascotaID, err := idFromNested(r.URL.Path)
    if err != nil {
        writeError(w, NewBadRequest("invalid_id", "ID de mascota inválido"))
        return
    }
    list, err := (models.CuidadoStore{DB: h.DB}).ListByMascota(r.Context(), mascotaID)
    if err != nil {
        writeError(w, err)
        return
    }
    respondJSON(w, http.StatusOK, list)
}

func (h *Handlers) CreateCuidadoForMascota(w http.ResponseWriter, r *http.Request) {
    mascotaID, err := idFromNested(r.URL.Path)
    if err != nil {
        writeError(w, NewBadRequest("invalid_id", "ID de mascota inválido"))
        return
    }
    var in struct {
        TipoCuidado  string `json:"tipo_cuidado" validate:"required,oneof=Vacunación Desparasitación 'Consulta Veterinaria' Baño"`
        Descripcion  string `json:"descripcion" validate:"required,min=2,max=500"`
        FechaCuidado string `json:"fecha_cuidado" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
    }
    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        writeError(w, NewBadRequest("invalid_json", "JSON inválido"))
        return
    }
    if err := h.validate.Struct(in); err != nil {
        if verrs, ok := err.(validator.ValidationErrors); ok {
            writeError(w, AppError{Code: "validation_error", Status: http.StatusBadRequest, Msg: "Datos inválidos", Fields: mapFieldErrors(verrs)})
        } else {
            writeError(w, NewBadRequest("validation_error", "Datos inválidos"))
        }
        return
    }
    t, err := time.Parse(time.RFC3339, in.FechaCuidado)
    if err != nil {
        writeError(w, NewBadRequest("invalid_datetime", "fecha_cuidado debe ser RFC3339"))
        return
    }
    c := &models.Cuidado{TipoCuidado: in.TipoCuidado, Descripcion: in.Descripcion, FechaCuidado: t, MascotaID: mascotaID}
    if err := (models.CuidadoStore{DB: h.DB}).Create(r.Context(), c); err != nil {
        writeError(w, err)
        return
    }
    respondJSON(w, http.StatusCreated, c)
}

func (h *Handlers) GetCuidado(w http.ResponseWriter, r *http.Request) {
    id, err := idFromPath(r.URL.Path)
    if err != nil {
        http.Error(w, "invalid id", http.StatusBadRequest)
        return
    }
    c, err := (models.CuidadoStore{DB: h.DB}).Get(r.Context(), id)
    if err != nil {
        writeError(w, err)
        return
    }
    respondJSON(w, http.StatusOK, c)
}

func (h *Handlers) UpdateCuidado(w http.ResponseWriter, r *http.Request) {
    id, err := idFromPath(r.URL.Path)
    if err != nil {
        http.Error(w, "invalid id", http.StatusBadRequest)
        return
    }
    var in struct {
        TipoCuidado  string `json:"tipo_cuidado" validate:"required,oneof=Vacunación Desparasitación 'Consulta Veterinaria' Baño"`
        Descripcion  string `json:"descripcion" validate:"required,min=2,max=500"`
        FechaCuidado string `json:"fecha_cuidado" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
        MascotaID    int64  `json:"mascota_id" validate:"required,gt=0"`
    }
    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        writeError(w, NewBadRequest("invalid_json", "JSON inválido"))
        return
    }
    if err := h.validate.Struct(in); err != nil {
        if verrs, ok := err.(validator.ValidationErrors); ok {
            writeError(w, AppError{Code: "validation_error", Status: http.StatusBadRequest, Msg: "Datos inválidos", Fields: mapFieldErrors(verrs)})
        } else {
            writeError(w, NewBadRequest("validation_error", "Datos inválidos"))
        }
        return
    }
    t, err := time.Parse(time.RFC3339, in.FechaCuidado)
    if err != nil {
        writeError(w, NewBadRequest("invalid_datetime", "fecha_cuidado debe ser RFC3339"))
        return
    }
    c := &models.Cuidado{ID: id, TipoCuidado: in.TipoCuidado, Descripcion: in.Descripcion, FechaCuidado: t, MascotaID: in.MascotaID}
    if err := (models.CuidadoStore{DB: h.DB}).Update(r.Context(), c); err != nil {
        writeError(w, err)
        return
    }
    respondJSON(w, http.StatusOK, c)
}

func (h *Handlers) DeleteCuidado(w http.ResponseWriter, r *http.Request) {
    id, err := idFromPath(r.URL.Path)
    if err != nil {
        http.Error(w, "invalid id", http.StatusBadRequest)
        return
    }
    if err := (models.CuidadoStore{DB: h.DB}).Delete(r.Context(), id); err != nil {
        writeError(w, err)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

// helpers

func respondJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(v)
}

func idFromPath(path string) (int64, error) {
    // expect .../{id}
    idx := strings.LastIndexByte(path, '/')
    if idx < 0 || idx+1 >= len(path) {
        return 0, errors.New("missing id")
    }
    return strconv.ParseInt(path[idx+1:], 10, 64)
}

func idFromNested(path string) (int64, error) {
    // expect /mascotas/{id}/cuidados
    parts := strings.Split(strings.Trim(path, "/"), "/")
    if len(parts) < 3 {
        return 0, errors.New("missing nested id")
    }
    return strconv.ParseInt(parts[1], 10, 64)
}

func parseDate(s string) (time.Time, error) {
    // expect YYYY-MM-DD
    return time.Parse("2006-01-02", s)
}

func parsePagination(q url.Values) (int64, int64, error) {
    limit := int64(50)
    offset := int64(0)
    if v := q.Get("limit"); v != "" {
        n, err := strconv.ParseInt(v, 10, 64)
        if err != nil || n <= 0 || n > 200 {
            return 0, 0, NewBadRequest("invalid_limit", "limit debe ser 1..200")
        }
        limit = n
    }
    if v := q.Get("offset"); v != "" {
        n, err := strconv.ParseInt(v, 10, 64)
        if err != nil || n < 0 {
            return 0, 0, NewBadRequest("invalid_offset", "offset debe ser >= 0")
        }
        offset = n
    }
    return limit, offset, nil
}

func mapFieldErrors(verrs validator.ValidationErrors) []FieldError {
    out := make([]FieldError, 0, len(verrs))
    for _, ve := range verrs {
        out = append(out, FieldError{Field: ve.Field(), Message: ve.Error()})
    }
    return out
}
