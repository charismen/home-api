package service

import (
	"context"
	"testing"

	"github.com/charismen/home-api/internal/repository"
	"github.com/charismen/home-api/pkg/apiclient"
)

// TestNewItemService tests the NewItemService function
func TestNewItemService(t *testing.T) {
	// Create a real repository for testing
	repo := &repository.ItemRepository{}
	
	// Create a real API client for testing
	client := &apiclient.Client{}
	
	// Create the service with real dependencies
	service := NewItemService(repo, client)
	
	// Verify service was created correctly
	if service == nil {
		t.Fatal("Expected service to be created, got nil")
	}
}

// TestItemService_GetAllItems tests the GetAllItems method
func TestItemService_GetAllItems(t *testing.T) {
	// Skip this test in normal runs as it requires a real repository
	t.Skip("Skipping test that requires a real repository")
	
	// Create a real repository for testing
	repo := &repository.ItemRepository{}
	
	// Create a real API client for testing
	client := &apiclient.Client{}
	
	// Create the service with real dependencies
	service := NewItemService(repo, client)
	
	// Test getting items
	items, err := service.GetAllItems()
	
	// Just verify the method runs without error
	// In a real test, we would check the items returned
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	t.Logf("Retrieved %d items", len(items))
}

// TestItemService_SyncItems tests the SyncItems method
func TestItemService_SyncItems(t *testing.T) {
	// Skip this test in normal runs as it requires a real API client
	t.Skip("Skipping test that requires a real API client")
	
	// Create a real repository for testing
	repo := &repository.ItemRepository{}
	
	// Create a real API client for testing
	client := &apiclient.Client{}
	
	// Create the service with real dependencies
	service := NewItemService(repo, client)
	
	// Test syncing items
	err := service.SyncItems(context.Background())
	
	// Just verify the method runs without error
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}