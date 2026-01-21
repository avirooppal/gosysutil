# ---------- build stage ----------
FROM golang:1.24.3-alpine AS build
WORKDIR /app

# install certs (needed for go mod download over https)
RUN apk add --no-cache ca-certificates

# deps
COPY go.mod go.sum ./
RUN go mod download

# source
COPY . .

# build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api ./cmd/api

# ---------- runtime stage ----------
FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=build /app/api .

EXPOSE 5001
CMD ["./api"]