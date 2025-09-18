package repository

import (
	"database/sql"

	"github.com/charismen/home-api/pkg/database"
)

// OrderRepository handles database operations for orders
type OrderRepository struct {
	db *sql.DB
}

// NewOrderRepository creates a new order repository
func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		db: database.DB,
	}
}

// OrderStatusSummary represents a summary of orders by status
type OrderStatusSummary struct {
	Status    string  `json:"status"`
	Count     int     `json:"count"`
	TotalAmount float64 `json:"total_amount"`
}

// CustomerSpend represents a customer's total spend
type CustomerSpend struct {
	CustomerID string  `json:"customer_id"`
	TotalSpend float64 `json:"total_spend"`
}

// GetOrderSummaryByStatus returns the number of orders and total amount per status in the last 30 days
func (r *OrderRepository) GetOrderSummaryByStatus() ([]OrderStatusSummary, error) {
	rows, err := r.db.Query(`
		SELECT 
			status, 
			COUNT(*) as count, 
			SUM(amount) as total_amount
		FROM orders
		WHERE created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
		GROUP BY status
		ORDER BY status
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var summaries []OrderStatusSummary
	for rows.Next() {
		var summary OrderStatusSummary
		err := rows.Scan(&summary.Status, &summary.Count, &summary.TotalAmount)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}
	
	return summaries, nil
}

// GetTopCustomersBySpend returns the top 5 customers by total spend
func (r *OrderRepository) GetTopCustomersBySpend() ([]CustomerSpend, error) {
	rows, err := r.db.Query(`
		SELECT 
			customer_id, 
			SUM(amount) as total_spend
		FROM orders
		WHERE status = 'PAID'
		GROUP BY customer_id
		ORDER BY total_spend DESC
		LIMIT 5
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var customers []CustomerSpend
	for rows.Next() {
		var customer CustomerSpend
		err := rows.Scan(&customer.CustomerID, &customer.TotalSpend)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	
	return customers, nil
}