# Message Queue Service

A lightweight, topic-based message queue service built with Go. It exposes HTTP APIs to manage subscriptions, publish messages to topics, and retrieve messages for subscribers.

## Features

- Subscribe an HTTP endpoint to one or more topics.
- Unsubscribe from selected topics using a subscriber ID.
- Publish messages to a topic.
- Retrieve messages for a subscriber and topic, filtered by publication date.
- JSON-based REST API.

> **Delivery semantics:** The API includes a subscriber callback URL, but the interface alone does not establish whether messages are pushed to that URL. The `/v1/retrieve` endpoint supports pull-based consumption. Confirm the implementation before claiming push delivery, persistence, acknowledgments, or delivery guarantees.

## Requirements

- Go 1.26

## Getting started

Clone your repository (replace the placeholder with your GitHub URL):

```bash
git clone https://github.com/<your-username>/<your-repository>.git
cd <your-repository>
```

Run the service:

```bash
go run .
```

The service is expected to listen on **port 3000**.

## API reference

Base URL: `http://localhost:3000`

### 1. Subscribe

`POST /v1/subscribe`

Register a subscriber endpoint for one or more topics.

```bash
curl -X POST http://localhost:3000/v1/subscribe \
  -H 'Content-Type: application/json' \
  -d '{"endpoint":"http://localhost:4001/webhook","topics":["orders","payments"]}'
```

Example success response:

```json
{
  "data": {"subscriber_id": "unique-subscriber-id"},
  "message": "Subscribed successfully"
}
```

### 2. Unsubscribe

`POST /v1/unsubscribe`

Remove a subscriber from specified topics.

```bash
curl -X POST http://localhost:3000/v1/unsubscribe \
  -H 'Content-Type: application/json' \
  -d '{"subscriber_id":"unique-subscriber-id","topics":["orders"]}'
```

Example success response:

```json
{"message":"Unsubscribed successfully from the specified topics"}
```

### 3. Publish

`POST /v1/publish`

Publish a message to a topic.

```bash
curl -X POST http://localhost:3000/v1/publish \
  -H 'Content-Type: application/json' \
  -d '{"topic":"orders","message":"Order 123 created"}'
```

Example success response:

```json
{"message":"Message published successfully to all topic subscribers"}
```

> The response text is from the documented API contract; it should not be interpreted as a confirmed delivery guarantee without checking the implementation.

### 4. Retrieve messages

`GET /v1/retrieve`

Retrieve messages for a subscriber and topic published from a given date.

| Query parameter | Description |
| --- | --- |
| `topic` | Topic name |
| `subId` | Subscriber ID returned by the subscribe API |
| `startDate` | RFC3339 timestamp (for example, `2024-01-01T00:00:00Z`) |

```bash
curl -G http://localhost:3000/v1/retrieve \
  --data-urlencode 'topic=orders' \
  --data-urlencode 'subscriberId=unique-subscriber-id' \
  --data-urlencode 'startDate=2024-01-01T00:00:00Z'
```

Example success response:

```json
{
  "data": [
    {
      "ID": "message-id",
      "Body": "Order 123 created",
      "PUBLISHED_AT": "2024-01-01T00:00:00Z",
      "CONSUMED_AT": "2024-01-01T01:00:00Z"
    }
  ]
}
```

### Errors

API errors use the following JSON shape:

```json
{"error":"error-message"}
```

HTTP status codes depend on the handler implementation.

## How it works

```text
Publisher
   |
   | POST /v1/publish
   v
Message Queue Service
   |-- Topic: orders
   |     |-- Subscriber A
   |     `-- Subscriber B
   `-- Topic: payments
         `-- Subscriber A

Subscriber --> GET /v1/retrieve --> Message Queue Service
```

1. A subscriber registers an endpoint and its topics with `/v1/subscribe`.
2. A publisher submits a message to `/v1/publish` for a topic.
3. Subscribers can request messages through `/v1/retrieve`.
4. A subscriber can stop following topics using `/v1/unsubscribe`.

## Example workflow

1. Start the service using `go run .`.
2. Subscribe an endpoint to the `orders` topic and copy the returned `subscriber_id`.
3. Publish `Order 123 created` to `orders`.
4. Retrieve messages using the subscriber ID and a suitable `startDate`.
5. Unsubscribe when the subscriber no longer needs the topic.

## Future improvements

- Durable message storage and recovery after restart.
- Explicit acknowledgment and per-subscriber delivery state.
- Background push delivery workers with retries and backoff.
- Dead-letter queues for repeatedly failing deliveries.
- Message ordering, pagination, and retention policies.
- Authentication, authorization, and subscriber endpoint validation.
- Metrics, tracing, and load tests.

## License

Add a `LICENSE` file if you intend to distribute this project under an open-source license.
