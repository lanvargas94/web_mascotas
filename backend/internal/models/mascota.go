package models

import (
    "context"
    "database/sql"
    "time"
)

type Mascota struct {
    ID              int64     `json:"id"`
    Nombre          string    `json:"nombre"`
    Especie         string    `json:"especie"`
    Raza            string    `json:"raza"`
    FechaNacimiento time.Time `json:"fecha_nacimiento"`
    Sexo            string    `json:"sexo"`
}

type MascotaStore struct{ DB *sql.DB }

func (s MascotaStore) Create(ctx context.Context, m *Mascota) error {
    q := `INSERT INTO mascotas(nombre, especie, raza, fecha_nacimiento, sexo)
          VALUES ($1,$2,$3,$4,$5) RETURNING id`
    return s.DB.QueryRowContext(ctx, q, m.Nombre, m.Especie, m.Raza, m.FechaNacimiento, m.Sexo).Scan(&m.ID)
}

func (s MascotaStore) Get(ctx context.Context, id int64) (*Mascota, error) {
    q := `SELECT id, nombre, especie, raza, fecha_nacimiento, sexo FROM mascotas WHERE id=$1`
    var m Mascota
    err := s.DB.QueryRowContext(ctx, q, id).Scan(&m.ID, &m.Nombre, &m.Especie, &m.Raza, &m.FechaNacimiento, &m.Sexo)
    if err != nil {
        return nil, err
    }
    return &m, nil
}

func (s MascotaStore) List(ctx context.Context) ([]Mascota, error) {
    q := `SELECT id, nombre, especie, raza, fecha_nacimiento, sexo FROM mascotas ORDER BY id`
    rows, err := s.DB.QueryContext(ctx, q)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    out := make([]Mascota, 0)
    for rows.Next() {
        var m Mascota
        if err := rows.Scan(&m.ID, &m.Nombre, &m.Especie, &m.Raza, &m.FechaNacimiento, &m.Sexo); err != nil {
            return nil, err
        }
        out = append(out, m)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }
    return out, nil
}

func (s MascotaStore) ListPaged(ctx context.Context, limit, offset int64) ([]Mascota, error) {
    q := `SELECT id, nombre, especie, raza, fecha_nacimiento, sexo FROM mascotas ORDER BY id LIMIT $1 OFFSET $2`
    rows, err := s.DB.QueryContext(ctx, q, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    out := make([]Mascota, 0)
    for rows.Next() {
        var m Mascota
        if err := rows.Scan(&m.ID, &m.Nombre, &m.Especie, &m.Raza, &m.FechaNacimiento, &m.Sexo); err != nil {
            return nil, err
        }
        out = append(out, m)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }
    return out, nil
}

func (s MascotaStore) Update(ctx context.Context, m *Mascota) error {
    q := `UPDATE mascotas SET nombre=$1, especie=$2, raza=$3, fecha_nacimiento=$4, sexo=$5 WHERE id=$6`
    _, err := s.DB.ExecContext(ctx, q, m.Nombre, m.Especie, m.Raza, m.FechaNacimiento, m.Sexo, m.ID)
    return err
}

func (s MascotaStore) Delete(ctx context.Context, id int64) error {
    _, err := s.DB.ExecContext(ctx, `DELETE FROM mascotas WHERE id=$1`, id)
    return err
}
