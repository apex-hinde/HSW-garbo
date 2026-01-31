package storage

import (
	"path/filepath"
	"testing"
)

func newTestDB(t *testing.T) *DB {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "employees.db")
	db, err := NewDB(path)
	if err != nil {
		t.Fatalf("NewDB failed: %v", err)
	}
	return db
}

func TestEmployeeCRUD(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	emp := &Employee{
		FirstName:  "Ada",
		LastName:   "Lovelace",
		Recipient:  "@ada",
		Wage:       123.45,
		Department: "Engineering",
	}

	id, err := db.CreateEmployee(emp)
	if err != nil {
		t.Fatalf("CreateEmployee failed: %v", err)
	}
	if id == 0 {
		t.Fatalf("CreateEmployee returned invalid id")
	}

	got, err := db.GetEmployee(int(id))
	if err != nil {
		t.Fatalf("GetEmployee failed: %v", err)
	}
	if got.FirstName != emp.FirstName || got.LastName != emp.LastName || got.Recipient != emp.Recipient || got.Wage != emp.Wage || got.Department != emp.Department {
		t.Fatalf("GetEmployee mismatch: got=%+v want=%+v", got, emp)
	}

	employees, err := db.ListEmployees()
	if err != nil {
		t.Fatalf("ListEmployees failed: %v", err)
	}
	if len(employees) != 1 {
		t.Fatalf("ListEmployees expected 1, got %d", len(employees))
	}

	got.Recipient = "@ada-updated"
	got.Wage = 200.00
	got.Department = "Research"
	if err := db.UpdateEmployee(got); err != nil {
		t.Fatalf("UpdateEmployee failed: %v", err)
	}

	updated, err := db.GetEmployee(got.ID)
	if err != nil {
		t.Fatalf("GetEmployee after update failed: %v", err)
	}
	if updated.Recipient != "@ada-updated" || updated.Wage != 200.00 || updated.Department != "Research" {
		t.Fatalf("UpdateEmployee mismatch: got=%+v", updated)
	}

	byDept, err := db.GetEmployeesByDepartment("Research")
	if err != nil {
		t.Fatalf("GetEmployeesByDepartment failed: %v", err)
	}
	if len(byDept) != 1 {
		t.Fatalf("GetEmployeesByDepartment expected 1, got %d", len(byDept))
	}

	if err := db.DeleteEmployee(updated.ID); err != nil {
		t.Fatalf("DeleteEmployee failed: %v", err)
	}
	if _, err := db.GetEmployee(updated.ID); err == nil {
		t.Fatalf("GetEmployee expected error after delete")
	}
}

