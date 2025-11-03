CREATE TABLE IF NOT EXISTS mascotas (
  id BIGSERIAL PRIMARY KEY,
  nombre TEXT NOT NULL,
  especie TEXT NOT NULL,
  raza TEXT NOT NULL,
  fecha_nacimiento DATE NOT NULL,
  sexo TEXT NOT NULL
);

-- Ensure constraints for valid enums
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_constraint WHERE conname = 'mascotas_especie_check'
  ) THEN
    ALTER TABLE mascotas
      ADD CONSTRAINT mascotas_especie_check CHECK (especie IN ('Perro','Gato','Conejo'));
  END IF;
  IF NOT EXISTS (
    SELECT 1 FROM pg_constraint WHERE conname = 'mascotas_sexo_check'
  ) THEN
    ALTER TABLE mascotas
      ADD CONSTRAINT mascotas_sexo_check CHECK (sexo IN ('Macho','Hembra'));
  END IF;
END$$;

CREATE TABLE IF NOT EXISTS cuidados (
  id BIGSERIAL PRIMARY KEY,
  tipo_cuidado TEXT NOT NULL,
  descripcion TEXT NOT NULL,
  fecha_cuidado TIMESTAMPTZ NOT NULL,
  mascota_id BIGINT NOT NULL REFERENCES mascotas(id) ON DELETE CASCADE
);

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_constraint WHERE conname = 'cuidados_tipo_check'
  ) THEN
    ALTER TABLE cuidados
      ADD CONSTRAINT cuidados_tipo_check CHECK (tipo_cuidado IN ('Vacunación','Desparasitación','Consulta Veterinaria','Baño'));
  END IF;
END$$;

CREATE INDEX IF NOT EXISTS idx_cuidados_mascota_id ON cuidados(mascota_id);
CREATE INDEX IF NOT EXISTS idx_cuidados_fecha ON cuidados(fecha_cuidado DESC);
