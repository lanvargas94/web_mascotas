# syntax=docker/dockerfile:1
FROM golang:1.22 AS build
WORKDIR /app
COPY go.mod .
# Copy go.sum if present (optional)
COPY go.sum* ./
RUN go mod download || true
# Now copy the rest of the source
COPY . .
# Ensure dependencies and go.sum are in sync
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/api

FROM gcr.io/distroless/base-debian12
WORKDIR /srv
COPY --from=build /app/server /srv/server
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/srv/server"]
