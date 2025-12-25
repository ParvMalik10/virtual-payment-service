# Virtual Payment Service (VPS)

A fault-tolerant financial transaction engine built in **Go (Golang)**. This service simulates a core banking backend, ensuring data integrity through **ACID compliant transactions** and preventing double-spending via **Idempotency keys**.

## 🚀 Key Features

* **Idempotency:** Implements a robust specific mechanism to handle duplicate requests (e.g., due to network retries). If a client sends the same `Idempotency-Key` twice, the system returns the cached successful response without re-processing the transaction.
* **ACID Compliance:** Uses strict database locking strategies. Transfers are atomic; if the credit fails after the debit, the entire transaction rolls back to prevent financial inconsistency.
* **Concurrency Safe:** Designed to handle multiple concurrent requests without race conditions on user balances.
* **Persistent Ledger:** detailed transaction logging using SQLite (scalable to PostgreSQL).

## 🛠 Tech Stack

* **Language:** Go (Golang)
* **Database:** SQLite3 (embedded) / SQL
* **Architecture:** RESTful API with Middleware logic

## ⚙️ How to Run Locally

### Prerequisites
* Go 1.18+ installed

### Installation
1.  Clone the repository:
    ```bash
    git clone [https://github.com/ParvMalik10/virtual-payment-service.git](https://github.com/ParvMalik10/virtual-payment-service.git)
    cd virtual-payment-service
    ```
2.  Install dependencies:
    ```bash
    go mod tidy
    ```
3.  Start the server:
    ```bash
    go run main.go
    ```
    *The server will start on port 8080.*

## 🧪 API Usage & Testing

You can test the API using `curl`.

### 1. Make a Payment (Success)
Send a request with a unique Idempotency Key:
```bash
curl -v -H "Idempotency-Key: key_unique_123" http://localhost:8080/payment
