package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/charismen/home-api/internal/models"
	"github.com/charismen/home-api/internal/repository"
	"github.com/charismen/home-api/pkg/apiclient"
	"github.com/charismen/home-api/pkg/redis"
)

const (
	// CacheKeyItems is the Redis key for cached items
	CacheKeyItems = "api:items"
	// CacheTTL is the TTL for cached items
	CacheTTL = 5 * time.Minute
)

// ItemService handles business logic for items
type ItemService struct {
	repo      *repository.ItemRepository
	apiClient *apiclient.Client
}

// NewItemService creates a new item service
func NewItemService(repo *repository.ItemRepository, apiClient *apiclient.Client) *ItemService {
	return &ItemService{
		repo:      repo,
		apiClient: apiClient,
	}
}

// SyncItems fetches items from the external API and stores them in the database
func (s *ItemService) SyncItems(ctx context.Context) error {
	// Fetch items from the external API
	items, err := s.apiClient.FetchItems(ctx, 20)
	if err != nil {
		return err
	}

	// Store each item in the database
	for _, item := range items {
		err := s.repo.SaveItem(item)
		if err != nil {
			log.Printf("Error saving item: %v", err)
			// Continue with other items
		}
	}

	// Invalidate cache
	redis.Delete(CacheKeyItems)

	return nil
}

// GetAllItems retrieves all items, using cache if available
func (s *ItemService) GetAllItems() ([]models.Item, error) {
	// Try to get from cache first
	cachedData, err := redis.Get(CacheKeyItems)
	if err == nil {
		// Cache hit
		var items []models.Item
		err = json.Unmarshal([]byte(cachedData), &items)
		if err == nil {
			return items, nil
		}
		// If unmarshal fails, continue to get from DB
	}

	// Cache miss or error, get from database
	items, err := s.repo.GetAllItems()
	if err != nil {
		return nil, err
	}

	// Store in cache for next time
	itemsJSON, err := json.Marshal(items)
	if err == nil {
		redis.Set(CacheKeyItems, string(itemsJSON), CacheTTL)
	}

	return items, nil
}