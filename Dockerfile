# ============================================
# Stage 1: Build
# ============================================
FROM golang:1.25-alpine AS builder

# Instalar dependências necessárias
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copiar arquivos de dependência
COPY go.mod go.sum ./

# Baixar dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Build da aplicação
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o api cmd/api/main.go

# ============================================
# Stage 2: Runtime
# ============================================
FROM alpine:latest

# Instalar dependências de runtime
RUN apk --no-cache add ca-certificates tzdata openssl

WORKDIR /app

# Copiar certificados
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copiar binário compilado
COPY --from=builder /app/api .

# Copiar dados do seed
COPY --from=builder /app/books.json .

# Gerar certificados TLS auto-assinados se não existirem
RUN mkdir -p certs && \
    openssl req -x509 -newkey rsa:2048 -keyout certs/key.pem -out certs/cert.pem \
      -days 365 -nodes -subj "/CN=localhost" 2>/dev/null || true

# Expor porta
EXPOSE 8443

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider https://localhost:8443/books || exit 1

# Executar aplicação
CMD ["./api"]
