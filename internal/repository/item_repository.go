package repository

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/charismen/home-api/internal/models"
	"github.com/charismen/home-api/pkg/database"
)

// ItemRepository handles database operations for items
type ItemRepository struct {
	db *sql.DB
}

// NewItemRepository creates a new item repository
func NewItemRepository() *ItemRepository {
	return &ItemRepository{
		db: database.DB,
	}
}

// SaveItem saves an item to the database (idempotent operation)
func (r *ItemRepository) SaveItem(item map[string]interface{}) error {
	// Extract necessary fields from the item
	name, _ := item["name"].(string)
	
	// Convert the entire item to JSON for storage
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}
	
	// Use external_id for idempotency
	externalID := getStringValue(item, "id")
	
	// Use REPLACE INTO for idempotent writes
	_, err = r.db.Exec(`
		REPLACE INTO items 
		(name, type, external_id, data, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?)`,
		name, "pokemon", externalID, string(data), time.Now(), time.Now())
	
	return err
}

// GetAllItems retrieves all items from the database
func (r *ItemRepository) GetAllItems() ([]models.Item, error) {
	rows, err := r.db.Query(`
		SELECT id, name, type, external_id, data, created_at, updated_at 
		FROM items
		ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var items []models.Item
	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Type,
			&item.ExternalID,
			&item.Data,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	
	return items, nil
}

// Helper function to safely extract string values from map
func getStringValue(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		switch v := val.(type) {
		case string:
			return v
		case float64:
			return string(int(v))
		default:
			// Convert to JSON string
			bytes, _ := json.Marshal(v)
			return string(bytes)
		}
	}
	return ""
}