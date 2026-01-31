package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/becomeliminal/nim-go-sdk/examples/hackathon-starter/internal/storage"
)

// Handler manages HTTP requests for employee operations
type Handler struct {
	db *storage.DB
}

// NewHandler creates a new API handler
func NewHandler(db *storage.DB) *Handler {
	return &Handler{db: db}
}

// RegisterRoutes sets up all API routes
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// Employee CRUD endpoints
	mux.HandleFunc("/api/employees", h.handleEmployees)
	mux.HandleFunc("/api/employees/", h.handleEmployeeByID)
	mux.HandleFunc("/api/employees/department/", h.handleEmployeesByDepartment)
}

// handleEmployees handles GET (list) and POST (create)
func (h *Handler) handleEmployees(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.listEmployees(w, r)
	case http.MethodPost:
		h.createEmployee(w, r)
	default:
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// handleEmployeeByID handles GET (read), PUT (update), and DELETE
func (h *Handler) handleEmployeeByID(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Extract ID from path
	path := strings.TrimPrefix(r.URL.Path, "/api/employees/")
	id, err := strconv.Atoi(path)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid employee id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getEmployee(w, r, id)
	case http.MethodPut:
		h.updateEmployee(w, r, id)
	case http.MethodDelete:
		h.deleteEmployee(w, r, id)
	default:
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// handleEmployeesByDepartment handles GET (list by department)
func (h *Handler) handleEmployeesByDepartment(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract department from path
	department := strings.TrimPrefix(r.URL.Path, "/api/employees/department/")
	if department == "" {
		respondError(w, http.StatusBadRequest, "department name required")
		return
	}

	employees, err := h.db.GetEmployeesByDepartment(department)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"employees":  employees,
		"count":      len(employees),
		"department": department,
	})
}

// listEmployees returns all employees
func (h *Handler) listEmployees(w http.ResponseWriter, r *http.Request) {
	employees, err := h.db.ListEmployees()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"employees": employees,
		"count":     len(employees),
	})
}

// createEmployee creates a new employee record
func (h *Handler) createEmployee(w http.ResponseWriter, r *http.Request) {
	var emp storage.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	id, err := h.db.CreateEmployee(&emp)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	emp.ID = int(id)
	respondJSON(w, http.StatusCreated, emp)
}

// getEmployee retrieves a single employee by ID
func (h *Handler) getEmployee(w http.ResponseWriter, r *http.Request, id int) {
	emp, err := h.db.GetEmployee(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, emp)
}

// updateEmployee updates an existing employee
func (h *Handler) updateEmployee(w http.ResponseWriter, r *http.Request, id int) {
	var emp storage.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	emp.ID = id
	if err := h.db.UpdateEmployee(&emp); err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, emp)
}

// deleteEmployee removes an employee record
func (h *Handler) deleteEmployee(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.db.DeleteEmployee(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"deleted": true,
		"id":      id,
	})
}

// Helper functions

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
