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
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %v", err)
		}
	})
	return db
}

func TestEmployeeCRUD(t *testing.T) {
	db := newTestDB(t)

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

func TestEmployeeValidation(t *testing.T) {
	tests := []struct {
		name    string
		emp     *Employee
		wantErr bool
	}{
		{
			name: "valid employee",
			emp: &Employee{
				FirstName:  "John",
				LastName:   "Doe",
				Recipient:  "@john",
				Wage:       50000,
				Department: "Sales",
			},
			wantErr: false,
		},
		{
			name: "empty first name",
			emp: &Employee{
				FirstName:  "",
				LastName:   "Doe",
				Recipient:  "@john",
				Wage:       50000,
				Department: "Sales",
			},
			wantErr: true,
		},
		{
			name: "empty last name",
			emp: &Employee{
				FirstName:  "John",
				LastName:   "",
				Recipient:  "@john",
				Wage:       50000,
				Department: "Sales",
			},
			wantErr: true,
		},
		{
			name: "empty recipient",
			emp: &Employee{
				FirstName:  "John",
				LastName:   "Doe",
				Recipient:  "",
				Wage:       50000,
				Department: "Sales",
			},
			wantErr: true,
		},
		{
			name: "negative wage",
			emp: &Employee{
				FirstName:  "John",
				LastName:   "Doe",
				Recipient:  "@john",
				Wage:       -100,
				Department: "Sales",
			},
			wantErr: true,
		},
		{
			name: "empty department",
			emp: &Employee{
				FirstName:  "John",
				LastName:   "Doe",
				Recipient:  "@john",
				Wage:       50000,
				Department: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.emp.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetEmployeeInvalidID(t *testing.T) {
	db := newTestDB(t)

	_, err := db.GetEmployee(0)
	if err == nil {
		t.Fatalf("GetEmployee expected error for invalid id")
	}

	_, err = db.GetEmployee(-1)
	if err == nil {
		t.Fatalf("GetEmployee expected error for negative id")
	}
}

func TestDeleteEmployeeInvalidID(t *testing.T) {
	db := newTestDB(t)

	err := db.DeleteEmployee(0)
	if err == nil {
		t.Fatalf("DeleteEmployee expected error for invalid id")
	}
}

func TestGetEmployeesByDepartment(t *testing.T) {
	db := newTestDB(t)

	// Create multiple employees in different departments
	emps := []*Employee{
		{FirstName: "Alice", LastName: "Smith", Recipient: "@alice", Wage: 80000, Department: "Engineering"},
		{FirstName: "Bob", LastName: "Jones", Recipient: "@bob", Wage: 70000, Department: "Engineering"},
		{FirstName: "Charlie", LastName: "Brown", Recipient: "@charlie", Wage: 60000, Department: "Sales"},
	}

	for _, emp := range emps {
		_, err := db.CreateEmployee(emp)
		if err != nil {
			t.Fatalf("CreateEmployee failed: %v", err)
		}
	}

	// Get employees by department
	engEmployees, err := db.GetEmployeesByDepartment("Engineering")
	if err != nil {
		t.Fatalf("GetEmployeesByDepartment failed: %v", err)
	}
	if len(engEmployees) != 2 {
		t.Fatalf("GetEmployeesByDepartment expected 2, got %d", len(engEmployees))
	}

	salesEmployees, err := db.GetEmployeesByDepartment("Sales")
	if err != nil {
		t.Fatalf("GetEmployeesByDepartment failed: %v", err)
	}
	if len(salesEmployees) != 1 {
		t.Fatalf("GetEmployeesByDepartment expected 1, got %d", len(salesEmployees))
	}

	// Test with empty department
	_, err = db.GetEmployeesByDepartment("")
	if err == nil {
		t.Fatalf("GetEmployeesByDepartment expected error for empty department")
	}
}
