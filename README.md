# Books API

Uma API REST em Go para gerenciar uma coleção de livros com suporte a filtros, paginação, ordenação e segurança.

## Stack Técnico

- **Linguagem:** Go 1.23
- **Framework:** Gin Gonic
- **Banco de Dados:** PostgreSQL 15
- **Containerização:** Docker & Docker Compose
- **Testes:** Go testing + testify
- **Monitoramento:** Prometheus (métricas de requisições)
- **TLS:** Certificados auto-assinados (desenvolvimento)

## Arquitetura

```
┌──────────────┐
│   Handler    │ HTTP Layer - Valida requests e retorna respostas
├──────────────┤
│   Service    │ Business Logic - Validação de dados e lógica de negócio
├──────────────┤
│ Repository   │ Data Access Layer - Queries SQL e acesso ao banco
├──────────────┤
│ PostgreSQL   │ Banco de Dados
└──────────────┘
```

## Setup Rápido

### Pré-requisitos

- Docker & Docker Compose
- OU
- Go 1.23+, PostgreSQL 15

### Opção 1: Docker Compose (Recomendado)

```bash
# Clonar repositório
git clone <repo>
cd books-api

# Iniciar tudo (banco + API)
docker compose up

# A API estará disponível em http://localhost:8080 (development)
# ou https://localhost:8443 se SERVER_TLS=true
```

### Opção 2: Desenvolvimento Local

```bash
# Instalar dependências
go mod download

# Gerar certificados TLS (desenvolvimento)
mkdir -p certs
openssl req -x509 -newkey rsa:2048 -keyout certs/key.pem -out certs/cert.pem \
  -days 365 -nodes -subj "/CN=localhost"

# IMPORTANTE: Iniciar PostgreSQL em outro terminal/janela
docker compose up postgres

# Em outro terminal: Executar migrations (automáticas ao iniciar a aplicação)
# As migrations são aplicadas automaticamente ao iniciar

# Configurar variáveis de ambiente (opcional)
cp .env.example .env.local

# Rodar a aplicação
go run cmd/api/main.go

# A API estará em http://localhost:8080 (padrão development)
```

**⚠️ Importante:** O PostgreSQL DEVE estar rodando antes de executar a aplicação. Use `docker compose up postgres` em um terminal separado.

## Variáveis de Ambiente

```bash
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=books_db

DATABASE_URL=postgres://postgres:postgres@localhost:5432/books_db?sslmode=disable
DATABASE_MAX_CONNECTIONS=25
DATABASE_SSL_MODE=disable

SERVER_ENV=development
SERVER_PORT=8080
SERVER_TLS=false

LOGGING_LEVEL=debug
LOGGING_FORMAT=json

CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001
```

## API Endpoints

### GET /api/v1/books
Lista todos os livros com filtros opcionais

**Query Parameters:**
- `author` (string): Filtro de autor (busca parcial, case-insensitive)
- `year` (int): Filtro de ano
- `masterpiece` (boolean): Filtro de obra-prima
- `sort` (string): Ordenação (name, author, masterpiece) - padrão: name
- `page` (int): Página (≥ 1) - padrão: 1
- `limit` (int): Itens por página (1-100) - padrão: 10

**Response (200):**
```json
{
  "data": [
    {
      "id": 1,
      "name": "One Hundred Years of Solitude",
      "author": "Gabriel García Márquez",
      "year": 1967,
      "masterpiece": true,
      "created_at": "2024-05-18T10:30:00Z",
      "updated_at": "2024-05-18T10:30:00Z"
    }
  ],
  "page": 1,
  "limit": 10,
  "total": 627
}
```

**Exemplos:**
```bash
# Listar todos os livros
curl http://localhost:8080/api/v1/books

# Filtrar por autor
curl "http://localhost:8080/api/v1/books?author=García"

# Filtrar por obra-prima
curl "http://localhost:8080/api/v1/books?masterpiece=true"

# Paginação
curl "http://localhost:8080/api/v1/books?page=2&limit=5"

# Múltiplos filtros
curl "http://localhost:8080/api/v1/books?author=Márquez&year=1967&sort=author"
```

### GET /api/v1/books/:id
Obtém um livro específico

```bash
curl http://localhost:8080/api/v1/books/1
```

**Response (200):**
```json
{
  "id": 1,
  "name": "One Hundred Years of Solitude",
  "author": "Gabriel García Márquez",
  "year": 1967,
  "masterpiece": true,
  "created_at": "2024-05-18T10:30:00Z",
  "updated_at": "2024-05-18T10:30:00Z"
}
```

### POST /api/v1/books
Cria um novo livro

**Body:**
```bash
curl -X POST http://localhost:8080/api/v1/books \
  -H "Content-Type: application/json" \
  -d '{
    "name": "One Hundred Years of Solitude",
    "author": "Gabriel García Márquez",
    "year": 1967,
    "masterpiece": true
  }'
```

**Validações:**
- `name` e `author` são obrigatórios
- `year` deve ser > 0 e ≤ ano atual
- `masterpiece` é opcional (padrão: false)

**Response (201):** Livro criado com ID gerado

### PUT /api/v1/books/:id
Atualiza um livro existente

```bash
curl -X PUT http://localhost:8080/api/v1/books/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Name",
    "author": "Updated Author",
    "year": 2024,
    "masterpiece": false
  }'
```

**Response (200):** Livro atualizado

### DELETE /api/v1/books/:id
Deleta um livro

```bash
curl -X DELETE http://localhost:8080/api/v1/books/1
```

**Response (204):** Sem conteúdo

## Health Check

```bash
curl http://localhost:8080/health
```

**Response (200):**
```json
{
  "status": "healthy"
}
```

## QA Testing - Postman Collection

Para testar todos os endpoints da API de forma estruturada:

### Importar a Collection

1. Abra o Postman
2. Clique em "Import" → "File"
3. Selecione `postman/Books-API.postman_collection.json`
4. Selecione `postman/Books-API.postman_environment.json`

### Configurar Ambiente

Após importar, certifique-se de:
1. Selecionar o environment "Production"
2. Verificar a variável `baseUrl`: `http://localhost:8080/api/v1`
3. Garantir que a API está rodando: `docker compose up`

### Testes Inclusos

- **00 - Health Check:** Verifica se a API está respondendo
- **01 - List All Books:** Testa listagem completa
- **02 - Create Book:** Cria um novo livro
- **03 - Get Book By ID:** Obtém um livro específico
- **04 - Update Book:** Atualiza um livro existente
- **05 - Delete Book:** Remove um livro
- **06 - Get Non-Existent Book:** Testa edge case (404)
- **07 - Create Book with Invalid Data:** Testa validação (400)
- **08 - List Books by Author Filter:** Testa filtros

Execute os testes na ordem apresentada para garantir fluxo correto (criar um livro antes de consultá-lo).

## Monitoramento

### Prometheus Metrics

A aplicação expõe métricas em `/metrics`:

```bash
curl http://localhost:8080/metrics
```

**Métricas disponíveis:**
- `http_requests_total` - Total de requisições HTTP (labels: method, endpoint, status)
- `http_request_duration_seconds` - Duração das requisições (labels: method, path)

## Seed do Banco de Dados

O banco é preenchido automaticamente na primeira execução com dados de `books.json`.

Para usar um arquivo customizado:
```bash
# No código, ajuste:
seed.SeedDatabase(svc, "seu_arquivo.json")
```

**Formato de JSON:**
```json
[
  {
    "name": "One Hundred Years of Solitude",
    "author": "Gabriel García Márquez",
    "year": 1967,
    "masterpiece": true
  }
]
```

## Estrutura do Projeto

```
books-api/
├── cmd/
│   └── api/
│       └── main.go           # Ponto de entrada
├── internal/
│   ├── config/               # Configuração da aplicação
│   ├── domain/               # Modelos de domínio (Book, erros)
│   ├── handler/              # HTTP handlers
│   ├── middleware/           # Middlewares (CORS, validação, etc)
│   ├── service/              # Lógica de negócio
│   ├── repository/           # Interfaces de acesso a dados
│   ├── infra/
│   │   ├── postgres/         # Implementação de repositório PostgreSQL
│   │   └── sqlc/             # Código gerado pelo sqlc
│   └── seed/                 # Seed do banco de dados
├── migrations/               # Migrations SQL
├── docker-compose.yaml       # Composição Docker
├── Dockerfile                # Container da aplicação
├── go.mod, go.sum           # Dependências Go
└── README.md
```

## Desenvolvimento

### Rodando Testes

```bash
# Todos os testes unitários
go test -short ./...

# Com cobertura
go test -short -cover ./...

# Teste específico
go test -short -run TestGetBookByID ./internal/handler

# Testes de integração (requer database postgres rodando)
go test ./...
```

**Nota:** Testes de integração são pulados por padrão com `-short`. Para rodá-los, é necessário ter PostgreSQL configurado conforme a seção [Setup Rápido](#setup-rápido).

### Gerando Código SQLC

Se modificar `queries.sql`, regenere o código:

```bash
# Instalar sqlc (se ainda não tiver)
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Gerar código
sqlc generate
```

## Migrations

As migrations são aplicadas automaticamente pelo Docker (via healthcheck do postgres).

Para adicionar uma nova migration:

1. Criar arquivo em `migrations/` com padrão: `NNN_descricao.sql`
2. Reiniciar os containers:
```bash
docker-compose down -v
docker-compose up
```

## Segurança

### TLS/HTTPS

- **Desenvolvimento (padrão):** HTTP sem TLS (SERVER_TLS=false)
- **Com TLS:** Certificados auto-assinados gerados automaticamente (SERVER_TLS=true)
- **Produção:** Use HTTPS com certificados de uma CA confiável (Let's Encrypt, etc) e configure SERVER_ENV=production

### Validação de Input

- Content-Type obrigatório para POST/PUT
- Tamanho máximo de payload: 10MB
- Validação de dados no service layer

### Rate Limiting

Proteção contra DoS com limite de requisições por IP.

### CORS

Configurável via variável `CORS_ALLOWED_ORIGINS`. Padrão: `http://localhost:3000,http://localhost:3001`

## Logging Estruturado

A aplicação usa `log/slog` (Go 1.21+) para logging estruturado em formato JSON:

```json
{
  "time": "2024-05-18T10:30:00.123Z",
  "level": "INFO",
  "msg": "request_completed",
  "method": "GET",
  "path": "/api/v1/books",
  "status": 200,
  "duration": "15.5ms"
}
```

## Graceful Shutdown

Ao receber SIGTERM ou SIGINT, a aplicação:
1. Para de aceitar novas conexões
2. Aguarda até 30 segundos para requisições em voo completarem
3. Fecha a conexão com o banco de dados
4. Encerra de forma limpa

## Troubleshooting

### "Connection refused" ao conectar ao banco

```bash
# Verificar se containers estão rodando
docker compose ps

# Ver logs do postgres
docker compose logs postgres
```

### Port 8080 já em uso

Altere em `compose.yaml` ou via variável de ambiente:
```bash
SERVER_PORT=8081 docker compose up
```

### Usar HTTPS em desenvolvimento

Definir variável de ambiente:
```bash
SERVER_TLS=true docker compose up
# A API estará em https://localhost:8080

# Com curl:
curl -k https://localhost:8080/api/v1/books
```

## Licença

MIT