package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

const employeeColumns = "id, first_name, last_name, recipient, wage, department"

const (
	// Query timeouts to prevent hanging queries
	defaultTimeout = 10 * time.Second
)

// scanEmployee scans a row into an Employee struct
func scanEmployee(scanner interface{ Scan(...interface{}) error }, emp *Employee) error {
	return scanner.Scan(&emp.ID, &emp.FirstName, &emp.LastName, &emp.Recipient, &emp.Wage, &emp.Department)
}

// CreateEmployee inserts a new employee record
func (d *DB) CreateEmployee(emp *Employee) (int64, error) {
	if err := emp.Validate(); err != nil {
		return 0, fmt.Errorf("validation error: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `INSERT INTO employees (first_name, last_name, recipient, wage, department) VALUES (?, ?, ?, ?, ?)`

	result, err := d.conn.ExecContext(ctx, query, emp.FirstName, emp.LastName, emp.Recipient, emp.Wage, emp.Department)
	if err != nil {
		return 0, fmt.Errorf("error inserting employee: %w", err)
	}

	return result.LastInsertId()
}

// GetEmployee retrieves an employee by ID
func (d *DB) GetEmployee(id int) (*Employee, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid employee id: %d", id)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `SELECT ` + employeeColumns + ` FROM employees WHERE id = ?`

	emp := &Employee{}
	err := scanEmployee(d.conn.QueryRowContext(ctx, query, id), emp)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("employee not found")
		}
		return nil, fmt.Errorf("error querying employee: %w", err)
	}

	return emp, nil
}

// ListEmployees retrieves all employees with a default limit for safety
func (d *DB) ListEmployees() ([]*Employee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `SELECT ` + employeeColumns + ` FROM employees`

	rows, err := d.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying employees: %w", err)
	}
	defer rows.Close()

	var employees []*Employee
	for rows.Next() {
		emp := &Employee{}
		err := scanEmployee(rows, emp)
		if err != nil {
			return nil, fmt.Errorf("error scanning employee: %w", err)
		}
		employees = append(employees, emp)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating employees: %w", err)
	}

	return employees, nil
}

// UpdateEmployee updates an existing employee record
func (d *DB) UpdateEmployee(emp *Employee) error {
	if err := emp.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	if emp.ID <= 0 {
		return fmt.Errorf("invalid employee id: %d", emp.ID)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `UPDATE employees SET first_name = ?, last_name = ?, recipient = ?, wage = ?, department = ? WHERE id = ?`

	result, err := d.conn.ExecContext(ctx, query, emp.FirstName, emp.LastName, emp.Recipient, emp.Wage, emp.Department, emp.ID)
	if err != nil {
		return fmt.Errorf("error updating employee: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("employee not found")
	}

	return nil
}

// DeleteEmployee removes an employee by ID
func (d *DB) DeleteEmployee(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid employee id: %d", id)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `DELETE FROM employees WHERE id = ?`

	result, err := d.conn.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting employee: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("employee not found")
	}

	return nil
}

// GetEmployeesByDepartment retrieves all employees in a department
func (d *DB) GetEmployeesByDepartment(department string) ([]*Employee, error) {
	if department == "" {
		return nil, fmt.Errorf("department cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `SELECT ` + employeeColumns + ` FROM employees WHERE department = ? ORDER BY last_name, first_name`

	rows, err := d.conn.QueryContext(ctx, query, department)
	if err != nil {
		return nil, fmt.Errorf("error querying employees by department: %w", err)
	}
	defer rows.Close()

	var employees []*Employee
	for rows.Next() {
		emp := &Employee{}
		err := scanEmployee(rows, emp)
		if err != nil {
			return nil, fmt.Errorf("error scanning employee: %w", err)
		}
		employees = append(employees, emp)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating employees: %w", err)
	}

	return employees, nil
}
