package db

import (
	"fmt"
	"goipban/config"
	"goipban/models"
	"os"
	"time"

	"database/sql"

	_ "modernc.org/sqlite"
)

// BannedDB manages the SQLite database
type BannedDB struct {
	db       *sql.DB
	filepath string
}

// NewBanDB creates or opens the SQLite database and initializes the table
func ConnectToBannedDB() (*BannedDB, error) {
	dbFilePath := config.BannedDB.FilePath
	db, err := sql.Open("sqlite", dbFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create the banned_ips table if it doesn't exist
	createTableQuery := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		%s TEXT PRIMARY KEY,
		%s DATETIME
	)`,
		config.BannedDB.Table,
		config.BannedDB.IPCol,
		config.BannedDB.BannedFromCol,
	)

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &BannedDB{db: db, filepath: dbFilePath}, nil
}

// AddBanEntry adds a new banned IP entry to the database
func (b *BannedDB) InsertBannedIP(banEntry models.BannedIPEntry) error {
	query := fmt.Sprintf(`
		INSERT OR REPLACE INTO %s (%s, %s)
		VALUES (?, ?)
	`,
		config.BannedDB.Table,
		config.BannedDB.IPCol,
		config.BannedDB.BannedFromCol,
	)

	_, err := b.db.Exec(query, banEntry.IP, banEntry.BannedFrom)
	if err != nil {
		return fmt.Errorf("failed to add banned IP: %w", err)
	}

	return nil
}

func (b *BannedDB) RemoveBannedIP(ip string) error {
	// Prepare the DELETE query within the transaction
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE %s = ?
	`,
		config.BannedDB.Table,
		config.BannedDB.IPCol,
	)

	_, err := b.db.Exec(query, ip)
	if err != nil {
		return fmt.Errorf("failed to remove banned IP: %w", err)
	}
	return nil
}

func (b *BannedDB) IsBannedIP(ip string) (bool, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(1) FROM %s WHERE %s = ?
	`,
		config.BannedDB.Table,
		config.BannedDB.IPCol,
	)

	var count int
	err := b.db.QueryRow(query, ip).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check banned IP: %w", err)
	}
	return count > 0, nil
}

// GetBannedIPsBefore retrieves items from the table where BanTime is less than the given time
func (b *BannedDB) GetIPsBannedBefore(beforeTime time.Time) ([]models.BannedIPEntry, error) {
	// Query to select items with BanTime less than the given value
	query := fmt.Sprintf(`
		SELECT %s, %s
		FROM %s
		WHERE %s < ?
	`,
		config.BannedDB.IPCol,
		config.BannedDB.BannedFromCol,
		config.BannedDB.Table,
		config.BannedDB.BannedFromCol)

	// Execute the query
	rows, err := b.db.Query(query, beforeTime)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Parse the results into a slice of BannedIPEntry
	var entries []models.BannedIPEntry
	for rows.Next() {
		var entry models.BannedIPEntry
		if err := rows.Scan(&entry.IP, &entry.BannedFrom); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		entries = append(entries, entry)
	}

	// Check for errors that occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return entries, nil
}

// Close closes the database connection
func (b *BannedDB) Close() error {
	return b.db.Close()
}

// delete delete the database connection
func (b *BannedDB) deleteTestDB() error {
	if b == nil {
		return nil
	}
	err := b.db.Close()
	if err != nil {
		return fmt.Errorf("failed to close db: %w", err)
	}
	err = os.Remove(b.filepath)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
