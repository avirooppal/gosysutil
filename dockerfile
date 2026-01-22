# ---------- Build Stage ----------
FROM golang:1.22-alpine AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o api ./cmd/api

# ---------- Runtime Stage ----------
FROM scratch
WORKDIR /app

COPY --from=build /app/api .

EXPOSE 5001
ENTRYPOINT ["./api"]
