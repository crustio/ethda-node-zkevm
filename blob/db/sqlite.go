package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func NewBlobDB(path string) (BlobDB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	// Create the table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS blocks (hash TEXT PRIMARY KEY, data BLOB)`)
	if err != nil {
		return nil, err
	}

	return &SqliteBlobDB{
		db: db,
	}, nil
}

type SqliteBlobDB struct {
	db *sql.DB
}

func (s *SqliteBlobDB) Put(key []byte, value []byte) error {
	_, err := s.db.Exec(`INSERT INTO blocks (hash, data) VALUES (?, ?)`, string(key), value)
	if err != nil {
		return err
	}

	return nil
}

func (s *SqliteBlobDB) Get(key []byte) ([]byte, error) {
	var retrievedHash string
	var retrievedData []byte
	err := s.db.QueryRow(`SELECT hash, data FROM blocks WHERE hash = ?`, string(key)).Scan(&retrievedHash, &retrievedData)
	if err != nil {
		return nil, err
	}

	return retrievedData, nil
}

func (s *SqliteBlobDB) Has(key []byte) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM blocks WHERE hash = ?)`
	var exists bool
	err := s.db.QueryRow(query, string(key)).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// TODO Close db when pool is closed
func (s *SqliteBlobDB) Close() {
	s.db.Close()
}
