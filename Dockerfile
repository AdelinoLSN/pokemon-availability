# --- Estágio de Compilação ---
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copia os arquivos de módulo (gerados pelo seu setup.sh)
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante do código
COPY . .

# FORÇAMOS o output (-o) para a raiz do container de build como 'pokemon-app'
RUN CGO_ENABLED=0 GOOS=linux go build -o /pokemon-app ./app

# --- Estágio Final (Runner) ---
FROM alpine:3.20

# Instala certificados
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copiamos o binário do local EXATO onde ele foi gerado (/pokemon-app)
# para o diretório atual do runner (.)
COPY --from=builder /pokemon-app .
COPY --from=builder /app/data ./data

# O comando chmod agora encontrará o arquivo com o nome correto
RUN chmod +x ./pokemon-app

# Ajustamos o CMD para usar o nome correto
CMD ["./pokemon-app"]
