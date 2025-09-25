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

1. **`internal/pkg/config/config.go`** - Adicionadas configurações do Mercado Pago
2. **`internal/pkg/mercadopago/client.go`** - Cliente HTTP para comunicação com a API
3. **`internal/app/model/pagamento.go`** - Modelo de dados para pagamentos

### Funcionalidades do Cliente

- ✅ Criação de Orders
- ✅ Criação de Pagamentos (Cartão de Crédito/Débito)
- ✅ Criação de Pagamentos PIX
- ✅ Consulta de Pagamentos
- ✅ Cancelamento de Pagamentos
- ✅ Tratamento de Erros da API

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
- `bolbradesco` - Boleto
- `pec` - PEC

## Próximos Passos

1. Configure as variáveis de ambiente
2. Execute a migração do banco de dados (tabela `pagamentos`)
3. Implemente os contratos de request/response
4. Crie o repository de pagamentos
5. Implemente o service de pagamentos
6. Crie os handlers HTTP
7. Configure as rotas

## Migração do Banco de Dados

Execute o seguinte SQL para criar a tabela de pagamentos:

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
    momento_criacao TIMESTAMP NOT NULL DEFAULT NOW(),
    momento_atualizacao TIMESTAMP NOT NULL DEFAULT NOW(),
    momento_aprovacao TIMESTAMP,
    momento_cancelamento TIMESTAMP
);

CREATE INDEX idx_pagamentos_cliente_id ON pagamentos(cliente_id);
CREATE INDEX idx_pagamentos_empresa_id ON pagamentos(empresa_id);
CREATE INDEX idx_pagamentos_mercado_pago_order_id ON pagamentos(mercado_pago_order_id);
CREATE INDEX idx_pagamentos_mercado_pago_payment_id ON pagamentos(mercado_pago_payment_id);
CREATE INDEX idx_pagamentos_status ON pagamentos(status);
CREATE INDEX idx_pagamentos_metodo_pagamento ON pagamentos(metodo_pagamento);
CREATE INDEX idx_pagamentos_momento_criacao ON pagamentos(momento_criacao);
```
