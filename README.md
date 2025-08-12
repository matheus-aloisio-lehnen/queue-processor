# 🚀 Event-Driven Queue Processor (Publisher + Dispatcher)

A lightweight microservice built with **Event-Driven Design** and **Pub/Sub**.  
It exposes an **HTTP Publisher API** that receives a validated payload and publishes it to a **topic**.  
A **Subscriber/Dispatcher** consumes messages from Pub/Sub and **routes each event to all matching handlers**, acting as a centralized event dispatcher.

```
[Client] → HTTP POST /publish → [Publisher] → Pub/Sub (topic)
                                              ↓
                                      [Subscriber] → [Dispatcher] → [Handler A]
                                                                         [Handler B]
                                                                         [Handler N]
```

---

## ✅ Key Features

- **HTTP → Pub/Sub publisher** with input validation.
- **Fan-out dispatcher**: every handler that supports a topic is invoked.
- **Loose coupling**: add/remove handlers without changing core flow.
- **Contract-first** payload with clear validation errors.
- **Production-minded**: guidance for idempotency, retries, DLQ, and observability.

---

## 🧩 Payload Contract (HTTP Publisher)

### Go DTOs (reference)
```go
type MetaDto struct {
  AccessToken *string `json:"access_token,omitempty" validate:"omitempty"`
  Topic       string  `json:"topic" validate:"required"`
}

type InputDto struct {
  Meta MetaDto     `json:"meta" validate:"required"`
  Data interface{} `json:"data" validate:"required"`
}
```

**Validation messages (translated):**
- `Meta.Topic.required` → **"Topic is required."**
- `Data.required` → **"Data is required."**

### Example request (JSON)
```json
{
  "meta": {
    "access_token": "optional-token-or-null",
    "topic": "user.created"
  },
  "data": {
    "id": 123,
    "name": "Ada Lovelace",
    "email": "ada@example.com"
  }
}
```

> Notes  
> • `meta.topic` (**required**) is the Pub/Sub topic to publish to.  
> • `meta.access_token` (**optional**) can be used for guardrails or downstream auth.  
> • `data` (**required**) is free-form and forwarded **as-is** to subscribers.

---

## 🌐 HTTP Publisher API

**Endpoint**
```
POST /publish
Content-Type: application/json
```

**Success (200 OK or 202 Accepted, depending on your server setup):**
```json
{
  "success": true,
  "topic": "user.created",
  "messageId": "abc123xyz"
}
```

**Validation error (422 Unprocessable Entity):**
```json
{
  "success": false,
  "errors": {
    "meta.topic": "Topic is required.",
    "data": "Data is required."
  }
}
```

**Example `curl`:**
```bash
curl -X POST http://localhost:3001/publish   -H "Content-Type: application/json"   -d '{
    "meta": { "access_token": null, "topic": "user.created" },
    "data": { "id": 123, "name": "Ada Lovelace" }
  }'
```

---

## 📡 Subscriber & Event Dispatcher

The subscriber listens to Pub/Sub subscriptions and, for **each** incoming message:

1. Parses `{ meta.topic, data }`.
2. **Finds all handlers** that support the `topic`.
3. **Dispatches** the event to every matching handler (fan-out).
4. Acks/Nacks according to handler outcomes and retry strategy.

### Conceptual interfaces (example)
```go
type Event struct {
  Topic string
  Data  any
}

type Handler interface {
  // Return true if this handler wants to process the given topic.
  Supports(topic string) bool

  // Handle the event (make it idempotent).
  Handle(ctx context.Context, e Event) error
}

type Dispatcher interface {
  Register(h Handler)
  Dispatch(ctx context.Context, e Event) error // calls all matching handlers
}
```

### Example handler (sketch)
```go
type SendWelcomeEmailHandler struct{}

func (h SendWelcomeEmailHandler) Supports(topic string) bool {
  return topic == "user.created"
}

func (h SendWelcomeEmailHandler) Handle(ctx context.Context, e Event) error {
  // e.Data → { id, name, email, ... }
  // Call email service, template engine, etc.
  return nil
}
```

> **Fan-out behavior:** The dispatcher invokes **every** handler whose `Supports(topic)` returns `true`.  
> Use this to implement notifications, audit logs, projections, caches, and integrations in parallel.

---

## ⚙️ Configuration (Environment)

Adjust to your broker/provider. Common variables:

```
HTTP_PORT=3001
PUBSUB_BROKER_URL=<broker-endpoint-or-connection-string>
TOPIC_PREFIX=            # optional, e.g., "dev."
SUBSCRIPTION_NAME=all-events-dispatch
PUBLISH_TIMEOUT_MS=5000  # example
MAX_RETRIES=5            # example
DEAD_LETTER_TOPIC=dlq.all-events  # optional
```

---

## 🏃 Running Locally

> Commands may vary with your layout; adapt as needed.

**1) Start the Publisher (HTTP API)**
```bash
go run ./cmd/publisher
# or: make run-publisher
```

**2) Start the Subscriber/Dispatcher**
```bash
go run ./cmd/subscriber
# or: make run-subscriber
```

**3) Publish a test event**
```bash
curl -X POST http://localhost:3001/publish   -H "Content-Type: application/json"   -d '{"meta":{"topic":"demo.event"},"data":{"hello":"world"}}'
```

---

## 🗂 Suggested Project Structure

```
/cmd
  /publisher        # HTTP server (POST /publish)
  /subscriber       # Pub/Sub consumer + dispatcher
/core
  /domain           # Event, interfaces, handler contracts
  /application      # Orchestrations, use cases
  /infra            # Pub/Sub client, logger, http, config, validation
    /validation     # Validation helpers (go-playground/validator)
/internal or /pkg   # Optionally shared components
```

---

## 🛡️ Design & Reliability Notes

- **Idempotency:** Make handlers idempotent (e.g., based on `messageId` or business keys) to tolerate retries.
- **Delivery semantics:** Assume **at-least-once** delivery from most Pub/Sub providers.
- **Poison messages:** Consider Dead-Letter Topics/Queues (DLQ) for repeated failures.
- **Backoff & retries:** Configure exponential backoff at the subscriber or broker level.
- **Security:** Validate input at the edge; `access_token` is optional and can be used for authentication/authorization.
- **Observability:** Add request IDs and message IDs to logs/metrics; export handler timings and error counts.

---

## 🔧 Extending the System (Add a New Handler)

1. **Create a handler** implementing `Supports(topic)` and `Handle(ctx, e)`.
2. **Register** the handler in the subscriber/dispatcher composition root.
3. **Publish** events to the relevant `topic`.
4. **Test** with real payloads and verify idempotency + retries.

---

## 📜 License

This project is distributed under the terms of the MIT License (or your chosen license).