

```markdown
# Real-Time Comment System with Kafka, Go, and Fiber

This project is a real-time streaming comment system built using:
- **Go**: for building the producer and consumer services
- **Fiber**: for exposing a REST API
- **Kafka**: for message queuing and streaming
- **Docker Compose**: for orchestrating Kafka and Zookeeper

It simulates real-time features like commenting on a live stream (e.g., Instagram Live) using Kafka as the messaging backbone.

---

## Project Structure

```

Go-kafka/
│
├── producer.go        # Fiber API server that produces messages to Kafka
├── consumer.go        # Kafka consumer that listens to the 'comments' topic
├── docker-compose.yml # Sets up Kafka and Zookeeper containers
└── README.md          # This file

````

---

## Getting Started

### Prerequisites

- Docker & Docker Compose
- Go 1.18+
- Postman or curl (for testing)

---

### Step-by-Step Setup

#### 1. Clone the repository

```bash
git clone [click on this link](https://github.com/PriyankaMandal3/Real_time_Comment.git) "
cd go-kafka-comments
````

#### 2. Start Kafka and Zookeeper using Docker Compose

```bash
docker-compose up -d
```

This runs:

* Kafka broker on `localhost:29092`
* Zookeeper on `localhost:2181`

---

#### 3. Start the Producer (API)

```bash
go run producer.go
```

Server will start at: `http://localhost:3000`

---

#### 4. Start the Consumer

In a new terminal:

```bash
go run consumer.go
```

You’ll see logs like:

```
consumer started
Received message Count: 1: | Topic (comments) | Message({"text":"Hello!"})
```

---

## Testing with Postman

### POST a Comment

* **Method:** `POST`
* **URL:** `http://localhost:3000/api/v1/comments`
* **Body:** `raw` -> `JSON`

```json
{
  "text": "Hello from Postman!"
}
```

### Expected Response:

```json
{
  "success": true,
  "message": "Comment pushed successfully",
  "comment": {
    "text": "Hello from Postman!"
  }
}
```

This sends the comment to the Kafka topic `comments`.

The consumer terminal will show the message received.

---

## Important Concepts & FAQs

### Where is the topic created?

> The topic `comments` is **created automatically** when the producer sends a message, thanks to Kafka’s `auto.create.topics.enable=true` (default).

You don’t explicitly define the topic in code.

---

### Where is the message stored?

> Kafka stores it in the `comments` topic and your **consumer reads it from there**.

If not stored in a database or file, **you won’t be able to retrieve it via a `GET` endpoint.**

---

### Why can’t I see messages in Postman?

> Postman only sees the result of the `POST` request. Kafka is a **streaming system**, not a database — it doesn’t persist data for external query unless you add **a database, file, or in-memory store**.

---

### What if I want to change ports?

* To change the **Fiber server port**, update:

  ```go
  app.Listen(":3000")
  ```

* To change Kafka’s port, modify the `docker-compose.yml`:

  ```yaml
  ports:
    - "29092:29092"
  KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092,PLAINTEXT_HOST://127.0.0.1:29092
  ```

Update all client configs (`producer.go`, `consumer.go`) to match the new port.


## Future Enhancements

You can extend the project to:

* Add message persistence to a **JSON file or database**
* Add a `GET /comments` endpoint to retrieve recent comments
* Integrate **Kafdrop** UI to visually inspect Kafka messages
* Deploy with Dockerized Go services and Kafka cluster

---

