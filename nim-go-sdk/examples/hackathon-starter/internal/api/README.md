# Employee Directory REST API

This directory contains the REST API handlers for managing the employee directory. The API follows standard RESTful conventions and returns JSON responses.

## Base Path
All endpoints are relative to `/api/employees`.

## Data Models

### Employee Object
When an employee is returned, it follows this structure:

```json
{
  "ID": 1,
  "FirstName": "Alice",
  "LastName": "Smith",
  "Recipient": "@alice",
  "Wage": 75000,
  "Department": "Engineering"
}
```

*Note: Field names are capitalized as they map directly to the Go struct fields.*

---

## Endpoints

### 1. List All Employees
Returns a list of all employees in the directory.

- **Method:** `GET`
- **URL:** `/api/employees`
- **Success Response:**
  - **Code:** 200 OK
  - **Content:**
    ```json
    {
      "employees": [
        { "ID": 1, "FirstName": "Alice", ... },
        { "ID": 2, "FirstName": "Bob", ... }
      ],
      "count": 2
    }
    ```

### 2. Create Employee
Adds a new employee record.

- **Method:** `POST`
- **URL:** `/api/employees`
- **Input Body:** (JSON)
  ```json
  {
    "FirstName": "Charlie",
    "LastName": "Brown",
    "Recipient": "@charlie",
    "Wage": 85000,
    "Department": "Engineering"
  }
  ```
- **Success Response:**
  - **Code:** 201 Created
  - **Content:** The created Employee object (including assigned ID).

### 3. Get Employee by ID
Retrieves details for a specific employee.

- **Method:** `GET`
- **URL:** `/api/employees/{id}`
- **Success Response:**
  - **Code:** 200 OK
  - **Content:** Employee object.
- **Error Response:**
  - **Code:** 404 Not Found
  - **Content:** `{"error": "employee not found"}`

### 4. Update Employee
Updates an existing employee record.

- **Method:** `PUT`
- **URL:** `/api/employees/{id}`
- **Input Body:** (JSON) Full or partial employee object.
- **Success Response:**
  - **Code:** 200 OK
  - **Content:** Updated Employee object.

### 5. Delete Employee
Removes an employee from the directory.

- **Method:** `DELETE`
- **URL:** `/api/employees/{id}`
- **Success Response:**
  - **Code:** 200 OK
  - **Content:**
    ```json
    {
      "deleted": true,
      "id": 1
    }
    ```

### 6. List by Department
Filters employees by their department name.

- **Method:** `GET`
- **URL:** `/api/employees/department/{department_name}`
- **Success Response:**
  - **Code:** 200 OK
  - **Content:**
    ```json
    {
      "employees": [...],
      "count": 5,
      "department": "Engineering"
    }
    ```

---

## Example Usage with curl

Assuming the server is running on `localhost:8080`:

### List all employees
```bash
curl http://localhost:8080/api/employees
```

### Create a new employee
```bash
curl -X POST http://localhost:8080/api/employees \
  -H "Content-Type: application/json" \
  -d '{
    "FirstName": "Alice",
    "LastName": "Smith",
    "Recipient": "@alice",
    "Wage": 75000,
    "Department": "Engineering"
  }'
```

### Get employee by ID
```bash
curl http://localhost:8080/api/employees/1
```

### Update an employee
```bash
curl -X PUT http://localhost:8080/api/employees/1 \
  -H "Content-Type: application/json" \
  -d '{
    "FirstName": "Alice",
    "LastName": "Smith",
    "Recipient": "@alice",
    "Wage": 80000,
    "Department": "Engineering"
  }'
```

### Delete an employee
```bash
curl -X DELETE http://localhost:8080/api/employees/1
```

### List employees by department
```bash
curl http://localhost:8080/api/employees/department/Engineering
```

---

## Error Handling
In case of an error, the API returns a relevant HTTP status code and a JSON body:

```json
{
  "error": "description of the error"
}
```

Common status codes:
- `400 Bad Request`: Invalid input or ID format.
- `404 Not Found`: Resource does not exist.
- `405 Method Not Allowed`: Incorrect HTTP method used.
- `500 Internal Server Error`: Database or server-side failure.
