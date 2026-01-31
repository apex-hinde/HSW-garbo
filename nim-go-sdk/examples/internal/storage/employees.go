package storage

import (
	"database/sql"
	"fmt"
)

const employeeColumns = "id, first_name, last_name, recipient, wage, department"

// scanEmployee scans a row into an Employee struct
func scanEmployee(scanner interface{ Scan(...interface{}) error }, emp *Employee) error {
	return scanner.Scan(&emp.ID, &emp.FirstName, &emp.LastName, &emp.Recipient, &emp.Wage, &emp.Department)
}

// CreateEmployee inserts a new employee record
func (d *DB) CreateEmployee(emp *Employee) (int64, error) {
	query := `INSERT INTO employees (first_name, last_name, recipient, wage, department) VALUES (?, ?, ?, ?, ?)`

	result, err := d.conn.Exec(query, emp.FirstName, emp.LastName, emp.Recipient, emp.Wage, emp.Department)
	if err != nil {
		return 0, fmt.Errorf("error inserting employee: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert id: %w", err)
	}

	return id, nil
}

// GetEmployee retrieves an employee by ID
func (d *DB) GetEmployee(id int) (*Employee, error) {
	query := `SELECT ` + employeeColumns + ` FROM employees WHERE id = ?`

	emp := &Employee{}
	err := scanEmployee(d.conn.QueryRow(query, id), emp)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("employee not found")
		}
		return nil, fmt.Errorf("error querying employee: %w", err)
	}

	return emp, nil
}

// ListEmployees retrieves all employees
func (d *DB) ListEmployees() ([]*Employee, error) {
	query := `SELECT ` + employeeColumns + ` FROM employees`

	rows, err := d.conn.Query(query)
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
	query := `UPDATE employees SET first_name = ?, last_name = ?, recipient = ?, wage = ?, department = ? WHERE id = ?`

	result, err := d.conn.Exec(query, emp.FirstName, emp.LastName, emp.Recipient, emp.Wage, emp.Department, emp.ID)
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
	query := `DELETE FROM employees WHERE id = ?`

	result, err := d.conn.Exec(query, id)
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
	query := `SELECT ` + employeeColumns + ` FROM employees WHERE department = ?`

	rows, err := d.conn.Query(query, department)
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

	return employees, nil
}
