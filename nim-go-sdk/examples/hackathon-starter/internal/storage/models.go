package storage

import (
	"fmt"
	"strings"
)

// Employee represents an employee record in the database
type Employee struct {
	ID         int     `db:"id" json:"id"`
	FirstName  string  `db:"first_name" json:"firstName"`
	LastName   string  `db:"last_name" json:"lastName"`
	Recipient  string  `db:"recipient" json:"recipient"`
	Wage       float64 `db:"wage" json:"wage"`
	Department string  `db:"department" json:"department"`
}

// Validate checks if an Employee has valid field values
func (e *Employee) Validate() error {
	if strings.TrimSpace(e.FirstName) == "" {
		return fmt.Errorf("first name cannot be empty")
	}
	if strings.TrimSpace(e.LastName) == "" {
		return fmt.Errorf("last name cannot be empty")
	}
	if strings.TrimSpace(e.Recipient) == "" {
		return fmt.Errorf("recipient cannot be empty")
	}
	if e.Wage < 0 {
		return fmt.Errorf("wage cannot be negative")
	}
	if strings.TrimSpace(e.Department) == "" {
		return fmt.Errorf("department cannot be empty")
	}
	return nil
}
