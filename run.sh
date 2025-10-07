#!/bin/sh

# Configurações básicas da aplicação
export DEBUG=false
export HTTP_SERVER_READ_TIMEOUT=20
export HTTP_SERVER_WRITE_TIMEOUT=60
export HTTP_SERVER_IDLE_TIMEOUT=120
export HTTP_SERVER_PORT=:1450

# Configurações do Banco de Dados
export DATABASE_POSTGRES_HOST=postgres
export DATABASE_POSTGRES_PORT=6432
export DATABASE_POSTGRES_POOL_MAX_LIFETIME_CONNECTION=300
export DATABASE_POSTGRES_NAME=jampa_trip_db
export DATABASE_POSTGRES_USER=jampa_trip_user
export DATABASE_POSTGRES_PASSWORD=jampa_trip_password
export DATABASE_POSTGRES_LOG=/workspace/logs/database.log

# Configurações JWT
export JWT_SECRET=jampa_trip_jwt_secret_key_2024_very_secure
export JWT_ACCESS_TOKEN_EXPIRATION=15m
export JWT_REFRESH_TOKEN_EXPIRATION=168h

# Configurações Redis
export REDIS_HOST=localhost
export REDIS_PORT=6379
export REDIS_PASSWORD=""
export REDIS_DB=0

# Configurações do Mercado Pago
export MERCADO_PAGO_ACCESS_TOKEN=your_access_token_here
export MERCADO_PAGO_PUBLIC_KEY=your_public_key_here
export MERCADO_PAGO_WEBHOOK_SECRET=your_webhook_secret_here
export MERCADO_PAGO_ENVIRONMENT=sandbox
export MERCADO_PAGO_BASE_URL=https://api.mercadopago.com

go run cmd/main.go