package models

import (
    "context"
    "database/sql"
    "time"
)

type Cuidado struct {
    ID           int64     `json:"id"`
    TipoCuidado  string    `json:"tipo_cuidado"`
    Descripcion  string    `json:"descripcion"`
    FechaCuidado time.Time `json:"fecha_cuidado"`
    MascotaID    int64     `json:"mascota_id"`
}

type CuidadoStore struct{ DB *sql.DB }

func (s CuidadoStore) Create(ctx context.Context, c *Cuidado) error {
    q := `INSERT INTO cuidados(tipo_cuidado, descripcion, fecha_cuidado, mascota_id)
          VALUES ($1,$2,$3,$4) RETURNING id`
    return s.DB.QueryRowContext(ctx, q, c.TipoCuidado, c.Descripcion, c.FechaCuidado, c.MascotaID).Scan(&c.ID)
}

func (s CuidadoStore) Get(ctx context.Context, id int64) (*Cuidado, error) {
    q := `SELECT id, tipo_cuidado, descripcion, fecha_cuidado, mascota_id FROM cuidados WHERE id=$1`
    var c Cuidado
    err := s.DB.QueryRowContext(ctx, q, id).Scan(&c.ID, &c.TipoCuidado, &c.Descripcion, &c.FechaCuidado, &c.MascotaID)
    if err != nil {
        return nil, err
    }
    return &c, nil
}

func (s CuidadoStore) ListByMascota(ctx context.Context, mascotaID int64) ([]Cuidado, error) {
    q := `SELECT id, tipo_cuidado, descripcion, fecha_cuidado, mascota_id FROM cuidados WHERE mascota_id=$1 ORDER BY fecha_cuidado DESC`
    rows, err := s.DB.QueryContext(ctx, q, mascotaID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    out := make([]Cuidado, 0)
    for rows.Next() {
        var c Cuidado
        if err := rows.Scan(&c.ID, &c.TipoCuidado, &c.Descripcion, &c.FechaCuidado, &c.MascotaID); err != nil {
            return nil, err
        }
        out = append(out, c)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }
    return out, nil
}

func (s CuidadoStore) Update(ctx context.Context, c *Cuidado) error {
    q := `UPDATE cuidados SET tipo_cuidado=$1, descripcion=$2, fecha_cuidado=$3, mascota_id=$4 WHERE id=$5`
    _, err := s.DB.ExecContext(ctx, q, c.TipoCuidado, c.Descripcion, c.FechaCuidado, c.MascotaID, c.ID)
    return err
}

func (s CuidadoStore) Delete(ctx context.Context, id int64) error {
    _, err := s.DB.ExecContext(ctx, `DELETE FROM cuidados WHERE id=$1`, id)
    return err
}
