package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/becomeliminal/nim-go-sdk/examples/hackathon-starter/internal/storage"
)

// TestListEmployees tests the GET /api/employees endpoint
func TestListEmployees(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create test data
	emp1 := &storage.Employee{FirstName: "Alice", LastName: "Smith", Recipient: "@alice", Wage: 75000, Department: "Engineering"}
	emp2 := &storage.Employee{FirstName: "Bob", LastName: "Jones", Recipient: "@bob", Wage: 80000, Department: "Sales"}

	id1, err := db.CreateEmployee(emp1)
	if err != nil {
		t.Fatalf("failed to create employee fixture: %v", err)
	}
	id2, err := db.CreateEmployee(emp2)
	if err != nil {
		t.Fatalf("failed to create employee fixture: %v", err)
	}

	handler := NewHandler(db)

	// Create request
	req := httptest.NewRequest("GET", "/api/employees", nil)
	w := httptest.NewRecorder()

	handler.handleEmployees(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Parse response
	var result map[string]interface{}
	json.NewDecoder(w.Body).Decode(&result)

	if count, ok := result["count"].(float64); !ok || count != 2 {
		t.Errorf("expected 2 employees, got %v", result["count"])
	}

	t.Logf("✓ Listed %d employees (IDs: %d, %d)", int(result["count"].(float64)), id1, id2)
}

// TestCreateEmployee tests the POST /api/employees endpoint
func TestCreateEmployee(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	handler := NewHandler(db)

	emp := storage.Employee{
		FirstName:  "Charlie",
		LastName:   "Brown",
		Recipient:  "@charlie",
		Wage:       85000,
		Department: "Engineering",
	}

	body, err := json.Marshal(emp)
	if err != nil {
		t.Fatalf("failed to marshal employee: %v", err)
	}
	req := httptest.NewRequest("POST", "/api/employees", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.handleEmployees(w, req)

	// Check status code
	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	// Parse response
	var created storage.Employee
	json.NewDecoder(w.Body).Decode(&created)

	if created.FirstName != emp.FirstName {
		t.Errorf("expected name %s, got %s", emp.FirstName, created.FirstName)
	}

	if created.ID == 0 {
		t.Error("expected employee ID to be set")
	}

	t.Logf("✓ Created employee (ID: %d, Name: %s)", created.ID, created.FirstName)
}

// TestGetEmployee tests the GET /api/employees/{id} endpoint
func TestGetEmployee(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create test data
	emp := &storage.Employee{FirstName: "Diana", LastName: "Prince", Recipient: "@diana", Wage: 90000, Department: "Management"}
	id, err := db.CreateEmployee(emp)
	if err != nil {
		t.Fatalf("failed to create employee fixture: %v", err)
	}

	handler := NewHandler(db)

	// Create request
	req := httptest.NewRequest("GET", "/api/employees/1", nil)
	w := httptest.NewRecorder()

	handler.handleEmployeeByID(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Parse response
	var retrieved storage.Employee
	json.NewDecoder(w.Body).Decode(&retrieved)

	if retrieved.FirstName != emp.FirstName {
		t.Errorf("expected name %s, got %s", emp.FirstName, retrieved.FirstName)
	}

	t.Logf("✓ Retrieved employee (ID: %d, Name: %s)", id, retrieved.FirstName)
}

// TestUpdateEmployee tests the PUT /api/employees/{id} endpoint
func TestUpdateEmployee(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create test data
	emp := &storage.Employee{FirstName: "Eve", LastName: "Wilson", Recipient: "@eve", Wage: 70000, Department: "HR"}
	id, err := db.CreateEmployee(emp)
	if err != nil {
		t.Fatalf("failed to create employee fixture: %v", err)
	}

	handler := NewHandler(db)

	// Update the employee
	updated := storage.Employee{
		FirstName:  "Eve",
		LastName:   "Wilson",
		Recipient:  "@eve",
		Wage:       75000, // Changed wage
		Department: "HR",
	}

	body, err := json.Marshal(updated)
	if err != nil {
		t.Fatalf("failed to marshal updated employee: %v", err)
	}
	req := httptest.NewRequest("PUT", "/api/employees/1", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.handleEmployeeByID(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Parse response
	var result storage.Employee
	json.NewDecoder(w.Body).Decode(&result)

	if result.Wage != 75000 {
		t.Errorf("expected wage 75000, got %v", result.Wage)
	}

	t.Logf("✓ Updated employee (ID: %d, New wage: %.0f)", id, result.Wage)
}

// TestDeleteEmployee tests the DELETE /api/employees/{id} endpoint
func TestDeleteEmployee(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create test data
	emp := &storage.Employee{FirstName: "Frank", LastName: "Miller", Recipient: "@frank", Wage: 80000, Department: "Engineering"}
	id, err := db.CreateEmployee(emp)
	if err != nil {
		t.Fatalf("failed to create employee fixture: %v", err)
	}

	handler := NewHandler(db)

	// Delete the employee
	req := httptest.NewRequest("DELETE", "/api/employees/1", nil)
	w := httptest.NewRecorder()

	handler.handleEmployeeByID(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Parse response
	var result map[string]interface{}
	json.NewDecoder(w.Body).Decode(&result)

	if deleted, ok := result["deleted"].(bool); !ok || !deleted {
		t.Errorf("expected deleted to be true, got %v", result["deleted"])
	}

	t.Logf("✓ Deleted employee (ID: %d)", id)
}

// TestListEmployeesByDepartment tests the GET /api/employees/department/{department} endpoint
func TestListEmployeesByDepartment(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create test data
	emp1 := &storage.Employee{FirstName: "Grace", LastName: "Lee", Recipient: "@grace", Wage: 85000, Department: "Engineering"}
	emp2 := &storage.Employee{FirstName: "Henry", LastName: "Chen", Recipient: "@henry", Wage: 80000, Department: "Engineering"}
	emp3 := &storage.Employee{FirstName: "Iris", LastName: "Kumar", Recipient: "@iris", Wage: 70000, Department: "Sales"}

	db.CreateEmployee(emp1)
	db.CreateEmployee(emp2)
	db.CreateEmployee(emp3)

	handler := NewHandler(db)

	// Create request
	req := httptest.NewRequest("GET", "/api/employees/department/Engineering", nil)
	w := httptest.NewRecorder()

	handler.handleEmployeesByDepartment(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Parse response
	var result map[string]interface{}
	json.NewDecoder(w.Body).Decode(&result)

	if count, ok := result["count"].(float64); !ok || count != 2 {
		t.Errorf("expected 2 employees in Engineering, got %v", result["count"])
	}

	t.Logf("✓ Listed %d employees in Engineering department", int(result["count"].(float64)))
}

// TestGetEmployeeNotFound tests error handling for non-existent employee
func TestGetEmployeeNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	handler := NewHandler(db)

	// Try to get non-existent employee
	req := httptest.NewRequest("GET", "/api/employees/999", nil)
	w := httptest.NewRecorder()

	handler.handleEmployeeByID(w, req)

	// Check status code
	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}

	t.Logf("✓ Correctly returned 404 for non-existent employee")
}

// TestInvalidEmployeeID tests error handling for invalid ID format
func TestInvalidEmployeeID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	handler := NewHandler(db)

	// Try with invalid ID
	req := httptest.NewRequest("GET", "/api/employees/invalid", nil)
	w := httptest.NewRecorder()

	handler.handleEmployeeByID(w, req)

	// Check status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	t.Logf("✓ Correctly returned 400 for invalid employee ID")
}

// TestCreateEmployeeInvalidData tests error handling for invalid input
func TestCreateEmployeeInvalidData(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	handler := NewHandler(db)

	// Create employee with invalid data (missing required field)
	invalidEmp := storage.Employee{
		FirstName:  "Jack",
		LastName:   "", // Missing last name
		Recipient:  "@jack",
		Wage:       80000,
		Department: "Sales",
	}

	body, err := json.Marshal(invalidEmp)
	if err != nil {
		t.Fatalf("failed to marshal invalid employee: %v", err)
	}
	req := httptest.NewRequest("POST", "/api/employees", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.handleEmployees(w, req)

	// Check status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	t.Logf("✓ Correctly rejected invalid employee data")
}

// Helper function to set up test database
func setupTestDB(t *testing.T) *storage.DB {
	db, err := storage.NewDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	return db
}
