# Web_mascotas

Aplicación  para gestionar mascotas y sus cuidados. Incluye:
- Backend en Go  con lógica en handlers y acceso a datos directo mediante stores en `models`.
- Frontend en Next.js 14 (TypeScript + TailwindCSS + SWR).
- PostgreSQL 13 con Docker.


## Requisitos
- Docker Desktop (Compose v2)


## Estructura

```
web_arquitecturamonolitica/
  backend/
    cmd/api/main.go
    internal/
      config/
        config.go
      database/
        db.go
        migrate.go
        migrations/0001_init.sql
      http/
        handlers.go
        router.go
        errors.go
      models/
        cuidado.go
        mascota.go
    go.mod
    Dockerfile
  frontend/
    app/
      layout.tsx
      page.tsx
      error.tsx
    components/
      PetForm.tsx
      PetList.tsx
      CareList.tsx
      Skeleton.tsx
      Toast.tsx
    lib/
      api/client.ts
      hooks.ts
    public/
    styles/
      globals.css
    next.config.mjs
    package.json
    tsconfig.json
    postcss.config.js
    tailwind.config.ts
    Dockerfile
  docker-compose.yml
```

## Puertos
- Frontend: `3000` (host) → Next.js
- Backend: `8080` (host) → API Go
- Base de datos: `5436` (host) → Postgres (redirige a `5432` del contenedor)

## Levantar con Docker Compose
1. En la raíz del proyecto: `docker compose up -d --build`
2. Verificación rápida:
   - Backend: `docker compose logs --follow backend` ("backend listening on :8080")
   - DB: `docker compose logs --follow db` (healthcheck OK)
   - Frontend: `docker compose logs --follow frontend` (Ready en 3000)
3. Probar:
   - Web: `http://localhost:3000`
   - Salud API: `http://localhost:8080/health`
   - Ready API: `http://localhost:8080/ready` (204)
4. Apagar: `docker compose down`

## Variables de entorno
- Backend
  - `PORT` (por defecto `8080`)
  - `DB_DSN` (por defecto apunta al servicio `db` en Compose)
  - `ALLOWED_ORIGINS` (por defecto `http://localhost:3000`)
- Frontend
  - `NEXT_PUBLIC_API_URL` (por defecto `http://localhost:8080`)

## Endpoints
- `GET /health` → `{ "status": "ok" }`
- `GET /ready` → 204
- Mascotas: `GET /mascotas?limit&offset`, `POST /mascotas`, `GET/PUT/DELETE /mascotas/{id}`
- Cuidados: `GET /mascotas/{id}/cuidados`, `POST /mascotas/{id}/cuidados`, `GET/PUT/DELETE /cuidados/{id}`

