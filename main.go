package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
)

var DB *sql.DB

func initDB() {
	var err error
	// Open connection to a file-based database called 'ledger.db'
	DB, err = sql.Open("sqlite3", "./ledger.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create the tables for Accounts and Transactions
	// The 'transactions' table is crucial for Idempotency checks
	schema := `
    CREATE TABLE IF NOT EXISTS accounts (
        id TEXT PRIMARY KEY,
        balance INTEGER NOT NULL
    );
    CREATE TABLE IF NOT EXISTS transactions (
        idempotency_key TEXT PRIMARY KEY,
        from_user TEXT,
        to_user TEXT,
        amount INTEGER
    );`

	_, err = DB.Exec(schema)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	// Seed dummy data: User A has $100, User B has $0
	// We use INSERT OR IGNORE so we don't duplicate data on restart
	DB.Exec(`INSERT OR IGNORE INTO accounts (id, balance) VALUES ('user_a', 100)`)
	DB.Exec(`INSERT OR IGNORE INTO accounts (id, balance) VALUES ('user_b', 0)`)
}

func paymentHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Get the Idempotency Key (The unique ID for this click/request)
	idempotencyKey := r.Header.Get("Idempotency-Key")
	if idempotencyKey == "" {
		http.Error(w, "Missing Idempotency-Key header", 400)
		return
	}

	// 2. IDEMPOTENCY CHECK: Have we processed this key before?
	var existingKey string
	err := DB.QueryRow("SELECT idempotency_key FROM transactions WHERE idempotency_key = ?", idempotencyKey).Scan(&existingKey)
	
	if err == nil {
		// If we found it, STOP. Do not charge again.
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Transaction already processed (Cached Response)"))
		return
	}

	// 3. START ATOMIC TRANSACTION (ACID Compliance)
	tx, err := DB.Begin()
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}

	// 4. EXECUTE TRANSFER (Hardcoded: User A -> User B, $10)
	// In a real app, you would parse the amount from the JSON body
	_, err = tx.Exec("UPDATE accounts SET balance = balance - 10 WHERE id = 'user_a'")
	if err != nil {
		tx.Rollback() // Safety switch: Undo if this fails
		return
	}

	_, err = tx.Exec("UPDATE accounts SET balance = balance + 10 WHERE id = 'user_b'")
	if err != nil {
		tx.Rollback() // Safety switch: Undo if this fails
		return
	}

	// 5. SAVE THE KEY so we don't process it again
	_, err = tx.Exec("INSERT INTO transactions (idempotency_key, from_user, to_user, amount) VALUES (?, 'user_a', 'user_b', 10)", idempotencyKey)
	if err != nil {
		tx.Rollback()
		return
	}

	// 6. COMMIT (Make changes permanent)
	err = tx.Commit()
	if err != nil {
		http.Error(w, "Commit failed", 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Payment Successful"))
}

func main() {
	initDB()
	http.HandleFunc("/payment", paymentHandler)
	
	// Updated output message to be generic and professional
	fmt.Println("Virtual Payment Service (VPS) is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
