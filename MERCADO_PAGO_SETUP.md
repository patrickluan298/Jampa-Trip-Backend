# Configuração do Mercado Pago - Jampa Trip Backend

## Variáveis de Ambiente Necessárias

Adicione as seguintes variáveis ao seu arquivo `.env`:

```env
# Configurações do Mercado Pago
# Obtenha essas credenciais em: https://www.mercadopago.com.br/developers/panel/credentials
MERCADO_PAGO_ACCESS_TOKEN=your_access_token_here
MERCADO_PAGO_PUBLIC_KEY=your_public_key_here
MERCADO_PAGO_WEBHOOK_SECRET=your_webhook_secret_here
MERCADO_PAGO_ENVIRONMENT=sandbox
MERCADO_PAGO_BASE_URL=https://api.mercadopago.com
```

## Como Obter as Credenciais

1. Acesse o [Painel de Desenvolvedores do Mercado Pago](https://www.mercadopago.com.br/developers/panel/credentials)
2. Faça login com sua conta Mercado Pago
3. Crie uma nova aplicação ou selecione uma existente
4. Copie as credenciais:
   - **Access Token**: Token de acesso para autenticação na API
   - **Public Key**: Chave pública para uso no frontend
   - **Webhook Secret**: Chave secreta para validar webhooks

## Ambientes

- **sandbox**: Ambiente de testes (recomendado para desenvolvimento)
- **production**: Ambiente de produção

## URLs da API

- **Sandbox**: `https://api.mercadopago.com`
- **Production**: `https://api.mercadopago.com`

## Estrutura Implementada

### Arquivos Criados/Modificados

1. **`pkg/config/config.go`** - Adicionadas configurações do Mercado Pago
2. **`pkg/mercadopago/client.go`** - Cliente HTTP para comunicação com a API
3. **`internal/model/pagamento.go`** - Modelo de dados para pagamentos
4. **`internal/service/pagamento.go`** - Lógica de negócio para pagamentos
5. **`internal/handler/pagamento.go`** - Handlers HTTP para pagamentos
6. **`internal/repository/pagamento.go`** - Camada de acesso a dados
7. **`internal/contract/pagamento_request.go`** - Contratos de request
8. **`internal/contract/pagamento_response.go`** - Contratos de response

### Funcionalidades Implementadas

- ✅ Criação de Orders
- ✅ Criação de Pagamentos (Cartão de Crédito/Débito)
- ✅ Criação de Pagamentos PIX
- ✅ Autorização de Cartões (Pré-autorização)
- ✅ Gestão de Cartões dos Clientes
- ✅ Consulta de Pagamentos
- ✅ Atualização de Pagamentos
- ✅ Tratamento de Erros da API

### Endpoints Implementados

- `POST /jampa-trip/api/v1/pagamentos/cartao-credito` - Pagamento com cartão de crédito
- `POST /jampa-trip/api/v1/pagamentos/cartao-debito` - Pagamento com cartão de débito
- `POST /jampa-trip/api/v1/pagamentos/pix` - Pagamento via PIX
- `GET /jampa-trip/api/v1/pagamentos` - Buscar pagamentos
- `GET /jampa-trip/api/v1/pagamentos/:id` - Obter pagamento por ID
- `PUT /jampa-trip/api/v1/pagamentos/:id` - Atualizar pagamento

### Status de Pagamento Suportados

- `pending` - Pendente
- `approved` - Aprovado
- `authorized` - Autorizado
- `in_process` - Em Processamento
- `in_mediation` - Em Mediação
- `rejected` - Rejeitado
- `cancelled` - Cancelado
- `refunded` - Reembolsado
- `charged_back` - Estornado

### Métodos de Pagamento Suportados

- `credit_card` - Cartão de Crédito
- `debit_card` - Cartão de Débito
- `pix` - PIX

## Configuração Atual

O sistema de pagamentos está **totalmente implementado** e funcional. Não são necessários passos adicionais de configuração além das variáveis de ambiente.

### Funcionalidades Disponíveis

1. **Autorização de Cartões**: Sistema de pré-autorização para cartões de crédito
2. **Pagamentos Diretos**: Cartão de crédito, débito e PIX
3. **Gestão de Cartões**: Cadastro e gerenciamento de cartões dos clientes
4. **Consulta de Pagamentos**: Busca e visualização de pagamentos
5. **Atualização de Status**: Atualização automática de status dos pagamentos

## Migração do Banco de Dados

A tabela de pagamentos é criada automaticamente pelo GORM com base no modelo `Pagamento`. O schema inclui:

```sql
CREATE TABLE pagamentos (
    id SERIAL PRIMARY KEY,
    cliente_id INTEGER NOT NULL REFERENCES clientes(id),
    empresa_id INTEGER NOT NULL REFERENCES empresas(id),
    mercado_pago_order_id VARCHAR(255) UNIQUE,
    mercado_pago_payment_id VARCHAR(255),
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    status_detail VARCHAR(100),
    valor DECIMAL(10,2) NOT NULL,
    moeda VARCHAR(3) NOT NULL DEFAULT 'BRL',
    metodo_pagamento VARCHAR(50) NOT NULL,
    descricao TEXT,
    numero_parcelas INTEGER DEFAULT 1,
    token_cartao VARCHAR(255),
    chave_pix VARCHAR(255),
    qr_code TEXT,
    -- Campos específicos para cartão
    last_four_digits VARCHAR(4),
    first_six_digits VARCHAR(6),
    payment_method_id VARCHAR(50),
    issuer_id VARCHAR(50),
    cardholder_name VARCHAR(255),
    captured BOOLEAN DEFAULT FALSE,
    transaction_amount_refunded DECIMAL(10,2) DEFAULT 0,
    -- Timestamps
    momento_criacao TIMESTAMP NOT NULL DEFAULT NOW(),
    momento_atualizacao TIMESTAMP NOT NULL DEFAULT NOW(),
    momento_aprovacao TIMESTAMP,
    momento_cancelamento TIMESTAMP,
    momento_autorizacao TIMESTAMP,
    momento_captura TIMESTAMP
);
```

## Exemplo de Uso

Consulte o arquivo `EXEMPLO_PAGAMENTO_CARTAO_CREDITO.md` para exemplos práticos de como usar a API de pagamentos.

## Monitoramento

O sistema inclui logs estruturados para monitoramento de pagamentos e integração com o Mercado Pago. Todos os eventos são registrados com detalhes completos para facilitar o debug e auditoria.