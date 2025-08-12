# 📨 Serviço de Notificações

Serviço simples e eficiente para envio de notificações através de múltiplos canais. Atualmente com suporte completo para e-mails e pronto para expansão.

## 📤 Como Enviar Notificações

### 1. E-mail
```json
{
  "userId": 123,
  "channel": "EMAIL",
  "subject": "Assunto do E-mail",
  "recipient": "usuario@exemplo.com",
  "payload": {
    "html": "Olá, {nome}!Seja bem-vindo à nossa plataforma.",
    "plainText": "Olá! Seja bem-vindo à nossa plataforma."
  }
}
```

### 2. SMS (em breve)
```json
{
  "userId": 123,
  "channel": "SMS",
  "recipient": "+5511999999999",
  "payload": {
    "message": "Sua conta foi criada com sucesso!"
  }
}
```

### 3. WhatsApp (em breve)
```json
{
  "userId": 123,
  "channel": "WHATSAPP",
  "recipient": "+5511999999999",
  "payload": {
    "message": "Seu pedido #123 foi enviado!"
  }
}
```

### 4. Notificação Push (em breve)
```json
{
  "userId": 123,
  "channel": "PUSH_NOTIFICATION",
  "recipient": "device-token-abc123",
  "payload": {
    "deviceToken": "device-token-abc123",
    "title": "Novo Aviso",
    "body": "Você tem uma nova mensagem não lida",
    "data": {
      "tipo": "mensagem",
      "id": "123"
    }
  }
}
```

## 📝 Campos Obrigatórios

- `userId`: ID do usuário destinatário (número inteiro positivo)
- `channel`: Canal de envio (EMAIL, SMS, WHATSAPP, PUSH_NOTIFICATION)
- `recipient`: Destinatário (e-mail, número de telefone ou token de dispositivo)
- `payload`: Conteúdo da mensagem (estrutura varia conforme o canal)

## ✅ Resposta de Sucesso

```json
{
  "success": true,
  "messageId": "abc123xyz",
  "status": "enviado"
}
```

## 🚀 Começando

1. Instale as dependências:
```bash
   npm install
   ```

2. Configure as variáveis de ambiente no arquivo `.env`

3. Inicie o serviço:
```bash
   npm run start:dev
   ```

A documentação completa da API estará disponível em `http://localhost:3001/api` quando o serviço estiver em execução.
