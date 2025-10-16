# ===== builder =====
FROM golang:1.22-alpine AS builder
WORKDIR /src

# pasang dep yang sering berubah dulu biar cache efisien
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY . .

# build
ARG APP_NAME=transaction-api
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o /out/app ./cmd/server

# ===== runner =====
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /out/app /app
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app"]
