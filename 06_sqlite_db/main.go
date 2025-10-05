package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "./app_data.db"
const tableName = "user_messages"

// createTable executes the SQL command to create the messages table.
func createTable(db *sql.DB) error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT NOT NULL,
			data TEXT NOT NULL
		);
	`, tableName)

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error executing CREATE TABLE: %w", err)
	}

	log.Printf("Table '%s' created successfully.", tableName)
	return nil
}

// insertMessage safely inserts a new row into the database.
func insertMessage(db *sql.DB, userID string, data string) error {
	// 1. --- DEFINE SQL WITH PLACEHOLDERS ---
	query := fmt.Sprintf("INSERT INTO %s (user_id, data) VALUES (?, ?)", tableName)

	// 2. --- EXECUTE INSERTION WITH VALUES ---
	_, err := db.Exec(query, userID, data)
	if err != nil {
		return fmt.Errorf("error inserting message: %w", err)
	}

	log.Printf("Inserted message for user '%s'.", userID)
	return nil
}

// queryMessages retrieves all messages for a given user.
func queryMessages(db *sql.DB, userID string) error {
	log.Printf("\nQuerying messages for user '%s'...", userID)

	// SQL Correction: Select ALL columns (id, user_id, data) and use a placeholder.
	query := fmt.Sprintf("SELECT id, user_id, data FROM %s WHERE user_id = ?", tableName)

	// 1. --- EXECUTE QUERY (using secure placeholder '?') ---
	// We pass the userID as a separate argument, preventing SQL Injection.
	rows, err := db.Query(query, userID)
	if err != nil {
		return fmt.Errorf("error querying messages: %w", err)
	}

	// 2. --- DEFER CLOSE ROWS (CRUCIAL!) ---
	// Always close the rows to free up database resources.
	defer rows.Close()

	// 3. --- ITERATE AND SCAN ROWS ---
	for rows.Next() {
		var id int
		var data string     // The JSON string we stored
		var dbUserID string // The user_id stored in the row

		// 4. --- SCAN COLUMNS INTO VARIABLES ---
		// The number and type MUST match the SELECT statement exactly.
		if err := rows.Scan(&id, &dbUserID, &data); err != nil {
			return fmt.Errorf("error scanning row: %w", err)
		}

		log.Printf(" | ID: %d | User: %s | Data: %s", id, dbUserID, data)
	}

	// Always check for errors that occurred during row iteration (e.g., network error)
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error during row iteration: %w", err)
	}

	return nil
}

func main() {
	fmt.Println("--- Project 6: SQLite Database (`database/sql`) ---")

	// Clean up any old database file for a fresh start
	os.Remove(dbPath)
	log.Printf("Starting fresh. Database file: %s", dbPath)

	// 1. --- OPEN THE DATABASE CONNECTION & PING ---
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	log.Println("Database opened and connection verified.")

	// 2. --- TASK 2: CREATE THE TABLE ---
	if err := createTable(db); err != nil {
		log.Fatalf("Table creation failed: %v", err)
	}

	// 3. --- TASK 3: INSERT DUMMY DATA ---
	if err := insertMessage(db, "user-A", `{"token_type":"access","value":"abc-123"}`); err != nil {
		log.Fatalf("Insert failed: %v", err)
	}
	if err := insertMessage(db, "user-B", `{"token_type":"refresh","value":"xyz-789"}`); err != nil {
		log.Fatalf("Insert failed: %v", err)
	}
	if err := insertMessage(db, "user-A", `{"status":"pending","last_attempt":1672531200}`); err != nil {
		log.Fatalf("Insert failed: %v", err)
	}

	// 4. --- TASK 4: EXECUTE THE QUERY ---
	if err := queryMessages(db, "user-A"); err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	log.Println("\nDatabase operations complete. Run: go run ./06_sqlite_db")
}
