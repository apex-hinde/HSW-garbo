// Hackathon Starter: Complete AI Financial Agent
// Build intelligent financial tools with nim-go-sdk + Liminal banking APIs
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/becomeliminal/nim-go-sdk/core"
	"github.com/becomeliminal/nim-go-sdk/examples/hackathon-starter/internal/api"
	"github.com/becomeliminal/nim-go-sdk/examples/hackathon-starter/internal/storage"
	"github.com/becomeliminal/nim-go-sdk/executor"
	"github.com/becomeliminal/nim-go-sdk/server"
	"github.com/becomeliminal/nim-go-sdk/tools"
	"github.com/joho/godotenv"
)

func main() {
	// ============================================================================
	// CONFIGURATION
	// ============================================================================
	// Load .env file if it exists (optional - will use system env vars if not found)
	_ = godotenv.Load()

	// Load configuration from environment variables
	// Create a .env file or export these in your shell

	anthropicKey := os.Getenv("ANTHROPIC_API_KEY")
	if anthropicKey == "" {
		log.Fatal("❌ ANTHROPIC_API_KEY environment variable is required")
	}

	liminalBaseURL := os.Getenv("LIMINAL_BASE_URL")
	if liminalBaseURL == "" {
		liminalBaseURL = "https://api.liminal.cash"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	employeesDBPath := os.Getenv("EMPLOYEES_DB_PATH")
	if employeesDBPath == "" {
		employeesDBPath = "employees.db"
	}

	// ============================================================================
	// LIMINAL EXECUTOR SETUP
	// ============================================================================
	// The HTTPExecutor handles all API calls to Liminal banking services.
	// Authentication is handled automatically via JWT tokens passed from the
	// frontend login flow (email/OTP). No API key needed!

	liminalExecutor := executor.NewHTTPExecutor(executor.HTTPExecutorConfig{
		BaseURL: liminalBaseURL,
	})
	log.Println("✅ Liminal API configured")

	// ============================================================================
	// EMPLOYEE DATABASE SETUP
	// ============================================================================
	// Local SQLite database for employee directory tools

	db, err := storage.NewDB(employeesDBPath)
	if err != nil {
		log.Fatalf("❌ Failed to initialize employee database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("⚠️ Failed to close employee database: %v", err)
		}
	}()

	log.Println("✅ Employee database configured")

	// ============================================================================
	// SERVER SETUP
	// ============================================================================
	// Create the nim-go-sdk server with Claude AI
	// The server handles WebSocket connections and manages conversations
	// Authentication is automatic: JWT tokens from the login flow are extracted
	// from WebSocket connections and forwarded to Liminal API calls

	srv, err := server.New(server.Config{
		AnthropicKey:    anthropicKey,
		SystemPrompt:    hackathonSystemPrompt,
		Model:           "claude-sonnet-4-20250514",
		MaxTokens:       4096,
		LiminalExecutor: liminalExecutor, // SDK automatically handles JWT extraction and forwarding
	})
	if err != nil {
		log.Fatal(err)
	}

	// ============================================================================
	// ADD LIMINAL BANKING TOOLS
	// ============================================================================
	// These are the 9 core Liminal tools that give your AI access to real banking:
	//
	// READ OPERATIONS (no confirmation needed):
	//   1. get_balance - Check wallet balance
	//   2. get_savings_balance - Check savings positions and APY
	//   3. get_vault_rates - Get current savings rates
	//   4. get_transactions - View transaction history
	//   5. get_profile - Get user profile info
	//   6. search_users - Find users by display tag
	//
	// WRITE OPERATIONS (require user confirmation):
	//   7. send_money - Send money to another user
	//   8. deposit_savings - Deposit funds into savings
	//   9. withdraw_savings - Withdraw funds from savings

	srv.AddTools(tools.LiminalTools(liminalExecutor)...)
	log.Println("✅ Added 9 Liminal banking tools")

	// ============================================================================
	// ADD CUSTOM TOOLS
	// ============================================================================
	// This is where you'll add your hackathon project's custom tools!
	// Below is an example spending analyzer tool to get you started.

	srv.AddTool(countEmployeeCount(liminalExecutor, *db))
	srv.AddTool(isPayrollDone(liminalExecutor, *db))
	srv.AddTool(doPayroll(liminalExecutor, *db))

	log.Println("✅ Added custom spending analyzer tool")

	// Employee management tools (CRUD + department lookup)
	srv.AddTools(createEmployeeTools(db)...)
	log.Println("✅ Added employee management tools")

	srv.AddTool(cashFlowAnalysisTool())
	srv.AddTool(analyseCashFlow(transactions, days))
	log.Println("✅ Added custom cash flow insight and projection tool")

	// TODO: Add more custom tools here!
	// Examples:
	//   - Savings goal tracker
	//   - Budget alerts
	//   - Spending category analyzer
	//   - Bill payment predictor
	//   - Cash flow forecaster

	// ============================================================================
	// REST API SETUP
	// ============================================================================
	// Create HTTP mux for REST API endpoints
	mux := http.NewServeMux()

	// Register employee API routes
	apiHandler := api.NewHandler(db)
	apiHandler.RegisterRoutes(mux)
	log.Println("✅ Added REST API endpoints for employees")

	// Register WebSocket handler from nim-go-sdk server
	mux.HandleFunc("/ws", srv.HandleWebSocket)

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// ============================================================================
	// START SERVER
	// ============================================================================

	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("🚀 Hackathon Starter Server Running")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("📡 WebSocket endpoint: ws://localhost:%s/ws", port)
	log.Printf("💚 Health check: http://localhost:%s/health", port)
	log.Printf("🔌 REST API: http://localhost:%s/api/employees", port)
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("Ready for connections! Start your frontend with: cd frontend && npm run dev")
	log.Println()

	// Start HTTP server with custom mux
	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// ============================================================================
// SYSTEM PROMPT
// ============================================================================
// This prompt defines your AI agent's personality and behavior
// Customize this to match your hackathon project's focus!

const hackathonSystemPrompt = `You are Nim, a friendly AI financial assistant built for the Liminal Vibe Banking Hackathon.

WHAT YOU DO:
You help users manage their money using Liminal's stablecoin banking platform. You can check balances, review transactions, send money, and manage savings - all through natural conversation.

CONVERSATIONAL STYLE:
- Be warm, friendly, and conversational - not robotic
- Use casual language when appropriate, but stay professional about money
- Ask clarifying questions when something is unclear
- Remember context from earlier in the conversation
- Explain things simply without being condescending

WHEN TO USE TOOLS:
- Use tools immediately for simple queries ("what's my balance?")
- For actions, gather all required info first ("send $50 to @alice")
- Always confirm before executing money movements
- Don't use tools for general questions about how things work

MONEY MOVEMENT RULES (IMPORTANT):
- ALL money movements require explicit user confirmation
- Show a clear summary before confirming:
  * send_money: "Send $50 USD to @alice"
  * deposit_savings: "Deposit $100 USD into savings"
  * withdraw_savings: "Withdraw $50 USD from savings"
- Never assume amounts or recipients
- Always use the exact currency the user specified

AVAILABLE BANKING TOOLS:
- Check wallet balance (get_balance)
- Check savings balance and APY (get_savings_balance)
- View savings rates (get_vault_rates)
- View transaction history (get_transactions)
- Get profile info (get_profile)
- Search for users (search_users)
- Send money (send_money) - requires confirmation
- Deposit to savings (deposit_savings) - requires confirmation
- Withdraw from savings (withdraw_savings) - requires confirmation

CUSTOM ANALYTICAL TOOLS:
- Analyze spending patterns (analyze_spending)

EMPLOYEE DIRECTORY TOOLS:
- Add employee (create_employee)
- Get employee by id (get_employee)
- List employees (list_employees)
- Update employee (update_employee)
- Delete employee (delete_employee)
- List employees by department (list_employees_by_department)

TIPS FOR GREAT INTERACTIONS:
- Proactively suggest relevant actions ("Want me to move some to savings?")
- Explain the "why" behind suggestions
- Celebrate financial wins ("Nice! Your savings earned $5 this month!")
- Be encouraging about savings goals
- Make finance feel less intimidating

Remember: You're here to make banking delightful and help users build better financial habits!`

// ============================================================================
// CUSTOM TOOL: SPENDING ANALYZER
// ============================================================================
// This is an example custom tool that demonstrates how to:
// 1. Define tool parameters with JSON schema
// 2. Call other Liminal tools from within your tool
// 3. Process and analyze the data
// 4. Return useful insights
//
// Use this as a template for your own hackathon tools!
func countEmployeeCount(liminalExecutor core.ToolExecutor, db storage.DB) core.Tool {
	return tools.New("count_employees").
		Description("counts the total number of employees").
		Handler(func(ctx context.Context, toolParams *core.ToolParams) (*core.ToolResult, error) {
			L, e := (db).ListEmployees()
			if e != nil {
				log.Fatal()
			}

			result := map[string]interface{}{
				"employee count": len(L),
			}

			return &core.ToolResult{
				Success: true,
				Data:    result,
			}, nil
		}).
		Build()
}
func doPayroll(liminalExecutor core.ToolExecutor, db storage.DB) core.Tool {
	return tools.New("fulfill_remaining_payroll").
		Description("does the payroll for all employees").
		Handler(func(ctx context.Context, toolParams *core.ToolParams) (*core.ToolResult, error) {
			txRequest := map[string]interface{}{
				"limit": 100, // Get up to 100 transactions
			}

			txRequestJSON, _ := json.Marshal(txRequest)
			txResponse, err := liminalExecutor.Execute(ctx, &core.ExecuteRequest{
				UserID:    toolParams.UserID,
				Tool:      "get_transactions",
				Input:     txRequestJSON,
				RequestID: toolParams.RequestID,
			})
			if err != nil {
				return &core.ToolResult{
					Success: false,
					Error:   fmt.Sprintf("failed to fetch transactions: %v", err),
				}, nil
			}
			if !txResponse.Success {
				return &core.ToolResult{
					Success: false,
					Error:   fmt.Sprintf("transaction fetch failed: %s", txResponse.Error),
				}, nil
			}
			var transactions []map[string]interface{}
			var txData map[string]interface{}
			if err := json.Unmarshal(txResponse.Data, &txData); err == nil {
				if txArray, ok := txData["transactions"].([]interface{}); ok {
					for _, tx := range txArray {
						if txMap, ok := tx.(map[string]interface{}); ok {
							transactions = append(transactions, txMap)
						}
					}
				}
			}
			newTransactions := removeOld(transactions)
			peopleWhoNeedToBePaid := checkPayments(newTransactions, db)

			L, e := (db).ListEmployees()
			if e != nil {
				log.Fatal()
			}
			paymentRequests := []map[string]interface{}{}
			for _, v := range L {
				if slices.Contains(peopleWhoNeedToBePaid, v.Recipient) {

					payRequest := map[string]interface{}{
						"recipient": v.Recipient,
						"amount":    v.Wage,
						"currency":  "USD",
						"note":      v.Recipient + " payroll",
					}
					paymentRequests = append(paymentRequests, payRequest)

				}
			}
			result := map[string]interface{}{
				"payment requests": paymentRequests,
			}
			log.Println(result)

			return &core.ToolResult{
				Success: true,
				Data:    result,
			}, nil
		}).Build()
}
func isPayrollDone(liminalExecutor core.ToolExecutor, db storage.DB) core.Tool {
	return tools.New("payroll_check").
		Description("checks if all payroll is done").
		Handler(func(ctx context.Context, toolParams *core.ToolParams) (*core.ToolResult, error) {

			txRequest := map[string]interface{}{
				"limit": 100, // Get up to 100 transactions
			}

			txRequestJSON, _ := json.Marshal(txRequest)
			txResponse, err := liminalExecutor.Execute(ctx, &core.ExecuteRequest{
				UserID:    toolParams.UserID,
				Tool:      "get_transactions",
				Input:     txRequestJSON,
				RequestID: toolParams.RequestID,
			})
			if err != nil {
				return &core.ToolResult{
					Success: false,
					Error:   fmt.Sprintf("failed to fetch transactions: %v", err),
				}, nil
			}
			if !txResponse.Success {
				return &core.ToolResult{
					Success: false,
					Error:   fmt.Sprintf("transaction fetch failed: %s", txResponse.Error),
				}, nil
			}
			var transactions []map[string]interface{}
			var txData map[string]interface{}
			if err := json.Unmarshal(txResponse.Data, &txData); err == nil {
				if txArray, ok := txData["transactions"].([]interface{}); ok {
					for _, tx := range txArray {
						if txMap, ok := tx.(map[string]interface{}); ok {
							transactions = append(transactions, txMap)
						}
					}
				}
			}
			newTransactions := removeOld(transactions)
			answer := checkPayments(newTransactions, db)
			result := map[string]interface{}{
				"people who have not been paid": answer,
			}
			return &core.ToolResult{
				Success: true,
				Data:    result,
			}, nil

		}).Build()

}
func removeOld(transactions []map[string]interface{}) []map[string]interface{} {

	var newTransactions []map[string]interface{}
	for _, v := range transactions {
		log.Println(v["createdAt"])
		var timeString string = v["createdAt"].(string)
		var time, e = time.Parse(time.RFC3339, timeString)
		if e != nil {
			log.Println("time not formated correctly")
		}
		log.Println("time:", time)
		log.Println((time).After(time.Local().AddDate(0, 0, -14)))
		if (time).After(time.Local().AddDate(0, 0, -14)) {
			newTransactions = append(newTransactions, v)
		}

	}
	return newTransactions
}

func checkPayments(transactions []map[string]interface{}, db storage.DB) []string {
	var acc []string
	var notes []string
	log.Println("transactions length:", len(transactions))

	for _, v := range transactions {
		var note string = v["note"].(string)
		log.Println("note:", note)
		words := strings.Fields(note)
		if len(words) == 2 {
			var word string = words[0]
			notes = append(notes, word)
		}
	}

	L, e := db.ListEmployees()
	if e != nil {
		log.Fatal()
	}
	for _, valueFromDb := range L {

		var recipient string = valueFromDb.Recipient
		if !slices.Contains(notes, recipient) {
			acc = append(acc, recipient)
		}

	}
	return acc
}

// ============================================================================
// EMPLOYEE DIRECTORY TOOLS
// ============================================================================

func createEmployeeTools(db *storage.DB) []core.Tool {
	return []core.Tool{
		createEmployeeTool(db),
		getEmployeeTool(db),
		listEmployeesTool(db),
		updateEmployeeTool(db),
		deleteEmployeeTool(db),
		listEmployeesByDepartmentTool(db),
	}
}

func createEmployeeTool(db *storage.DB) core.Tool {
	handler := func(ctx context.Context, toolParams *core.ToolParams) (*core.ToolResult, error) {
		var params struct {
			FirstName  string  `json:"first_name"`
			LastName   string  `json:"last_name"`
			Recipient  string  `json:"recipient"`
			Wage       float64 `json:"wage"`
			Department string  `json:"department"`
		}
		if err := json.Unmarshal(toolParams.Input, &params); err != nil {
			return &core.ToolResult{Success: false, Error: fmt.Sprintf("invalid input: %v", err)}, nil
		}

		emp := &storage.Employee{
			FirstName:  params.FirstName,
			LastName:   params.LastName,
			Recipient:  params.Recipient,
			Wage:       params.Wage,
			Department: params.Department,
		}

		id, err := db.CreateEmployee(emp)
		if err != nil {
			return &core.ToolResult{Success: false, Error: err.Error()}, nil
		}

		emp.ID = int(id)
		return &core.ToolResult{Success: true, Data: emp}, nil
	}

	return tools.New("create_employee").
		Description("Create a new employee record in the local employee directory.").
		Schema(tools.ObjectSchema(map[string]interface{}{
			"first_name": tools.StringProperty("Employee first name"),
			"last_name":  tools.StringProperty("Employee last name"),
			"recipient":  tools.StringProperty("Employee recipient handle, e.g. @ada"),
			"wage":       tools.NumberProperty("Employee wage (non-negative)"),
			"department": tools.StringProperty("Employee department"),
		})).
		Handler(handler).
		Build()
}

func getEmployeeTool(db *storage.DB) core.Tool {
	handler := func(ctx context.Context, toolParams *core.ToolParams) (*core.ToolResult, error) {
		var params struct {
			ID int `json:"id"`
		}
		if err := json.Unmarshal(toolParams.Input, &params); err != nil {
			return &core.ToolResult{Success: false, Error: fmt.Sprintf("invalid input: %v", err)}, nil
		}

		emp, err := db.GetEmployee(params.ID)
		if err != nil {
			return &core.ToolResult{Success: false, Error: err.Error()}, nil
		}

		return &core.ToolResult{Success: true, Data: emp}, nil
	}

	return tools.New("get_employee").
		Description("Get an employee record by id.").
		Schema(tools.ObjectSchema(map[string]interface{}{
			"id": tools.IntegerProperty("Employee id"),
		})).
		Handler(handler).
		Build()
}

func listEmployeesTool(db *storage.DB) core.Tool {
	handler := func(ctx context.Context, toolParams *core.ToolParams) (*core.ToolResult, error) {
		emps, err := db.ListEmployees()
		if err != nil {
			return &core.ToolResult{Success: false, Error: err.Error()}, nil
		}
		return &core.ToolResult{Success: true, Data: map[string]interface{}{"employees": emps, "count": len(emps)}}, nil
	}

	return tools.New("list_employees").
		Description("List all employees in the directory.").
		Schema(tools.ObjectSchema(map[string]interface{}{})).
		Handler(handler).
		Build()
}

func updateEmployeeTool(db *storage.DB) core.Tool {
	handler := func(ctx context.Context, toolParams *core.ToolParams) (*core.ToolResult, error) {
		var params struct {
			ID         int     `json:"id"`
			FirstName  string  `json:"first_name"`
			LastName   string  `json:"last_name"`
			Recipient  string  `json:"recipient"`
			Wage       float64 `json:"wage"`
			Department string  `json:"department"`
		}
		if err := json.Unmarshal(toolParams.Input, &params); err != nil {
			return &core.ToolResult{Success: false, Error: fmt.Sprintf("invalid input: %v", err)}, nil
		}

		emp := &storage.Employee{
			ID:         params.ID,
			FirstName:  params.FirstName,
			LastName:   params.LastName,
			Recipient:  params.Recipient,
			Wage:       params.Wage,
			Department: params.Department,
		}

		if err := db.UpdateEmployee(emp); err != nil {
			return &core.ToolResult{Success: false, Error: err.Error()}, nil
		}

		return &core.ToolResult{Success: true, Data: emp}, nil
	}

	return tools.New("update_employee").
		Description("Update an existing employee record by id.").
		Schema(tools.ObjectSchema(map[string]interface{}{
			"id":         tools.IntegerProperty("Employee id"),
			"first_name": tools.StringProperty("Employee first name"),
			"last_name":  tools.StringProperty("Employee last name"),
			"recipient":  tools.StringProperty("Employee recipient handle, e.g. @ada"),
			"wage":       tools.NumberProperty("Employee wage (non-negative)"),
			"department": tools.StringProperty("Employee department"),
		})).
		Handler(handler).
		Build()
}

func deleteEmployeeTool(db *storage.DB) core.Tool {
	handler := func(ctx context.Context, toolParams *core.ToolParams) (*core.ToolResult, error) {
		var params struct {
			ID int `json:"id"`
		}
		if err := json.Unmarshal(toolParams.Input, &params); err != nil {
			return &core.ToolResult{Success: false, Error: fmt.Sprintf("invalid input: %v", err)}, nil
		}

		if err := db.DeleteEmployee(params.ID); err != nil {
			return &core.ToolResult{Success: false, Error: err.Error()}, nil
		}

		return &core.ToolResult{Success: true, Data: map[string]interface{}{"deleted": true, "id": params.ID}}, nil
	}

	return tools.New("delete_employee").
		Description("Delete an employee record by id.").
		Schema(tools.ObjectSchema(map[string]interface{}{
			"id": tools.IntegerProperty("Employee id"),
		})).
		Handler(handler).
		Build()
}

func listEmployeesByDepartmentTool(db *storage.DB) core.Tool {
	handler := func(ctx context.Context, toolParams *core.ToolParams) (*core.ToolResult, error) {
		var params struct {
			Department string `json:"department"`
		}
		if err := json.Unmarshal(toolParams.Input, &params); err != nil {
			return &core.ToolResult{Success: false, Error: fmt.Sprintf("invalid input: %v", err)}, nil
		}

		emps, err := db.GetEmployeesByDepartment(params.Department)
		if err != nil {
			return &core.ToolResult{Success: false, Error: err.Error()}, nil
		}

		return &core.ToolResult{Success: true, Data: map[string]interface{}{"employees": emps, "count": len(emps), "department": params.Department}}, nil
	}

	return tools.New("list_employees_by_department").
		Description("List employees in a specific department.").
		Schema(tools.ObjectSchema(map[string]interface{}{
			"department": tools.StringProperty("Department name"),
		})).
		Handler(handler).
		Build()
}

func cashFlowAnalysisTool() core.Tool {
	//Use analyseCashFlow somewhere
}

// OLSModel represents a simple linear regression model
type OLSModel struct {
	Weight float64
	Bias   float64
}

type DayData struct {
	DayNumber int
	Date      time.Time
	Amount    float64
	Count     int
}

// Fit trains the model using analytical OLS solution
// β = (X^T X)^(-1) X^T y
func (m *OLSModel) Fit(X, y []float64) error {
	if len(X) != len(y) {
		return fmt.Errorf("X and y must have the same length")
	}

	n := float64(len(X))
	if n == 0 {
		return fmt.Errorf("need at least one data point")
	}

	// Calculate means
	var sumX, sumY, sumXY, sumX2 float64
	for i := 0; i < len(X); i++ {
		sumX += X[i]
		sumY += y[i]
		sumXY += X[i] * y[i]
		sumX2 += X[i] * X[i]
	}

	// Analytical OLS solution
	// weight = (n*Σ(xy) - Σx*Σy) / (n*Σ(x²) - (Σx)²)
	// bias = (Σy - weight*Σx) / n
	denominator := n*sumX2 - sumX*sumX
	if math.Abs(denominator) < 1e-10 {
		return fmt.Errorf("cannot fit model: data is too uniform")
	}

	m.Weight = (n*sumXY - sumX*sumY) / denominator
	m.Bias = (sumY - m.Weight*sumX) / n

	return nil
}

// Predict makes a single prediction
func (m *OLSModel) Predict(x float64) float64 {
	return m.Weight*x + m.Bias
}

// PredictBatch makes batch predictions
func (m *OLSModel) PredictBatch(x []float64) []float64 {
	predictions := make([]float64, len(x))
	for i, val := range x {
		predictions[i] = m.Predict(val)
	}
	return predictions
}

// CalculateR2 calculates R-squared score
func (m *OLSModel) CalculateR2(X, y []float64) float64 {
	predictions := m.PredictBatch(X)

	// Calculate mean of y
	var sumY float64
	for _, val := range y {
		sumY += val
	}
	meanY := sumY / float64(len(y))

	// Calculate SS_res and SS_tot
	var ssRes, ssTot float64
	for i := 0; i < len(y); i++ {
		ssRes += math.Pow(y[i]-predictions[i], 2)
		ssTot += math.Pow(y[i]-meanY, 2)
	}

	if ssTot == 0 {
		return 0
	}

	return 1 - (ssRes / ssTot)
}

// CalculateMSE calculates Mean Squared Error
func (m *OLSModel) CalculateMSE(X, y []float64) float64 {
	predictions := m.PredictBatch(X)

	var sumSquaredError float64
	for i := 0; i < len(y); i++ {
		error := y[i] - predictions[i]
		sumSquaredError += error * error
	}

	return sumSquaredError / float64(len(y))
}

// Split regressions every 30 days, to display variation in trends over the course of months
func MonthSlicer(days []DayData, income []float64) []OLSModel {
	if len(days) == 0 {
		return []OLSModel{}
	}

	window := 30
	regs := make([]OLSModel, 0, len(days)/window+1)

	// start represents the earliest day selected, or of the 30-day window
	// as such, it increments by the window size (+30)
	for start := 0; start < len(days); start += window {
		end := start + window
		if end > len(days) {
			end = len(days)
		}

		// Build X and y for this window
		size := end - start
		X := make([]float64, size)
		y := make([]float64, size)

		for i := 0; i < size; i++ {
			X[i] = float64(days[start+i].DayNumber)
			y[i] = income[start+i]
		}

		// Create an OLS regression for this time frame
		model := OLSModel{}
		err := model.Fit(X, y)
		if err == nil {
			regs = append(regs, model)
		}
	}

	return regs
}

// Takes an array of OLS regressions and calculates a weighted average of the weights, and normal average of the bias to project fo the next 30 days
func ProjectCashFlow(regs []OLSModel, weighted bool) OLSModel {
	if len(regs) == 1 {
		return regs[0]
	}

	var sumW float64
	var sumB float64
	var sumWeights float64

	for i := 0; i < len(regs); i++ {
		w := regs[i].Weight
		b := regs[i].Bias

		if weighted {
			// weight factor increases with the recency of the time slice
			factor := float64(1+i) / float64(len(regs))
			sumW += w * factor
			sumWeights += factor
		} else {
			sumW += w
			sumWeights += 1
		}

		sumB += b
	}

	avgW := sumW / sumWeights
	avgB := sumB / float64(len(regs))

	projection := OLSModel{Weight: avgW, Bias: avgB}

	return projection
}

// analyseCashFlow processes transaction data and returns insights
func analyseCashFlow(transactions []map[string]interface{}, days int) map[string]interface{} {
	// Handle edge cases
	if len(transactions) == 0 {
		return map[string]interface{}{
			"error": "no transactions provided",
		}
	}

	if days <= 0 {
		days = 7 // Default to 7 days prediction
	}

	// ========================================================================
	// STEP 1: Parse and aggregate transaction data by day
	// ========================================================================

	dayMap := make(map[string]*DayData)
	var minDate time.Time

	for i, txn := range transactions {
		// Extract date (support multiple formats)
		var txnDate time.Time
		var err error

		if dateStr, ok := txn["date"].(string); ok {
			// Try multiple date formats
			formats := []string{
				"2006-01-02",
				"2006-01-02T15:04:05Z",
				"2006-01-02 15:04:05",
				time.RFC3339,
			}

			for _, format := range formats {
				txnDate, err = time.Parse(format, dateStr)
				if err == nil {
					break
				}
			}

			if err != nil {
				log.Printf("Warning: could not parse date for transaction %d: %v", i, dateStr)
				continue
			}
		} else if dateTime, ok := txn["date"].(time.Time); ok {
			txnDate = dateTime
		} else {
			log.Printf("Warning: transaction %d has invalid date format", i)
			continue
		}

		// Extract amount
		var amount float64
		switch v := txn["amount"].(type) {
		case float64:
			amount = v
		case float32:
			amount = float64(v)
		case int:
			amount = float64(v)
		case int64:
			amount = float64(v)
		default:
			log.Printf("Warning: transaction %d has invalid amount format", i)
			continue
		}

		// Track minimum date
		if minDate.IsZero() || txnDate.Before(minDate) {
			minDate = txnDate
		}

		// Aggregate by day
		dateKey := txnDate.Format("2006-01-02")
		if _, exists := dayMap[dateKey]; !exists {
			dayMap[dateKey] = &DayData{
				Date:   txnDate,
				Amount: 0,
				Count:  0,
			}
		}
		dayMap[dateKey].Amount += amount
		dayMap[dateKey].Count++
	}

	if len(dayMap) == 0 {
		return map[string]interface{}{
			"error": "no valid transactions after parsing",
		}
	}

	// Convert map to slice and sort by date
	var dayDataSlice []DayData
	for _, data := range dayMap {
		dayDataSlice = append(dayDataSlice, *data)
	}

	sort.Slice(dayDataSlice, func(i, j int) bool {
		return dayDataSlice[i].Date.Before(dayDataSlice[j].Date)
	})

	// Assign day numbers (1-indexed from first transaction)
	for i := range dayDataSlice {
		daysSinceStart := int(dayDataSlice[i].Date.Sub(minDate).Hours() / 24)
		dayDataSlice[i].DayNumber = daysSinceStart + 1
	}

	// ========================================================================
	// STEP 2: Prepare data for OLS regression
	// ========================================================================
	X := make([]float64, len(dayDataSlice))
	y := make([]float64, len(dayDataSlice))

	for i, data := range dayDataSlice {
		X[i] = float64(data.DayNumber)
		y[i] = data.Amount
	}

	// ========================================================================
	// STEP 3: Fit OLS model
	// ========================================================================
	model := &OLSModel{}
	err := model.Fit(X, y)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("failed to fit model: %v", err),
		}
	}

	// ========================================================================
	// STEP 4: Calculate model statistics
	// ========================================================================
	r2 := model.CalculateR2(X, y)
	mse := model.CalculateMSE(X, y)
	rmse := math.Sqrt(mse)

	// Calculate trend
	var trend string
	if model.Weight > 0.01 {
		trend = "increasing"
	} else if model.Weight < -0.01 {
		trend = "decreasing"
	} else {
		trend = "stable"
	}

	// ========================================================================
	// STEP 5: Make future predictions
	// ========================================================================
	lastDayNumber := dayDataSlice[len(dayDataSlice)-1].DayNumber
	lastDate := dayDataSlice[len(dayDataSlice)-1].Date

	futurePredictions := make([]map[string]interface{}, days)
	for i := 0; i < days; i++ {
		futureDayNumber := lastDayNumber + i + 1
		futureDate := lastDate.AddDate(0, 0, i+1)
		predictedAmount := model.Predict(float64(futureDayNumber))

		futurePredictions[i] = map[string]interface{}{
			"day":              futureDayNumber,
			"date":             futureDate.Format("2006-01-02"),
			"predicted_amount": math.Round(predictedAmount*100) / 100, // Round to 2 decimals
		}
	}

	// ========================================================================
	// STEP 6: Calculate historical statistics
	// ========================================================================
	var totalAmount, totalCount float64
	var minAmount, maxAmount float64
	minAmount = y[0]
	maxAmount = y[0]

	for i, data := range dayDataSlice {
		totalAmount += data.Amount
		totalCount += float64(data.Count)

		if y[i] < minAmount {
			minAmount = y[i]
		}
		if y[i] > maxAmount {
			maxAmount = y[i]
		}
	}

	avgAmount := totalAmount / float64(len(dayDataSlice))
	avgTransactionsPerDay := totalCount / float64(len(dayDataSlice))

	// ========================================================================
	// STEP 7: Format historical data
	// ========================================================================
	historicalData := make([]map[string]interface{}, len(dayDataSlice))
	for i, data := range dayDataSlice {
		predicted := model.Predict(float64(data.DayNumber))
		residual := data.Amount - predicted

		historicalData[i] = map[string]interface{}{
			"day":               data.DayNumber,
			"date":              data.Date.Format("2006-01-02"),
			"actual_amount":     math.Round(data.Amount*100) / 100,
			"predicted_amount":  math.Round(predicted*100) / 100,
			"residual":          math.Round(residual*100) / 100,
			"transaction_count": data.Count,
		}
	}

	// ========================================================================
	// STEP 8: Build response
	// ========================================================================
	result := map[string]interface{}{
		"model": map[string]interface{}{
			"equation":  fmt.Sprintf("y = %.2f + %.2fx", model.Bias, model.Weight),
			"weight":    math.Round(model.Weight*100) / 100,
			"bias":      math.Round(model.Bias*100) / 100,
			"r_squared": math.Round(r2*10000) / 10000,
			"mse":       math.Round(mse*100) / 100,
			"rmse":      math.Round(rmse*100) / 100,
		},
		"insights": map[string]interface{}{
			"trend":                    trend,
			"total_days":               len(dayDataSlice),
			"total_amount":             math.Round(totalAmount*100) / 100,
			"total_transactions":       int(totalCount),
			"avg_amount_per_day":       math.Round(avgAmount*100) / 100,
			"avg_transactions_per_day": math.Round(avgTransactionsPerDay*100) / 100,
			"min_daily_amount":         math.Round(minAmount*100) / 100,
			"max_daily_amount":         math.Round(maxAmount*100) / 100,
			"date_range": map[string]string{
				"start": dayDataSlice[0].Date.Format("2006-01-02"),
				"end":   dayDataSlice[len(dayDataSlice)-1].Date.Format("2006-01-02"),
			},
		},
		"predictions":     futurePredictions,
		"historical_data": historicalData,
	}

	return result
}

// ============================================================================
// HACKATHON IDEAS
// ============================================================================
// Here are some ideas for custom tools you could build:
//
// 1. SAVINGS GOAL TRACKER
//    - Track progress toward savings goals
//    - Calculate how long until goal is reached
//    - Suggest optimal deposit amounts
//
// 2. BUDGET ANALYZER
//    - Set spending limits by category
//    - Alert when approaching limits
//    - Compare actual vs. planned spending
//
// 3. RECURRING PAYMENT DETECTOR
//    - Identify subscription payments
//    - Warn about upcoming bills
//    - Suggest savings opportunities
//
// 4. CASH FLOW FORECASTER
//    - Predict future balance based on patterns
//    - Identify potential low balance periods
//    - Suggest when to save vs. spend
//
// 5. SMART SAVINGS ADVISOR
//    - Analyze spare cash available
//    - Recommend savings deposits
//    - Calculate interest projections
//
// 6. SPENDING INSIGHTS
//    - Categorize spending automatically
//    - Compare to typical user patterns
//    - Highlight unusual activity
//
// 7. FINANCIAL HEALTH SCORE
//    - Calculate overall financial wellness
//    - Track improvements over time
//    - Provide actionable recommendations
//
// 8. PEER COMPARISON (anonymous)
//    - Compare savings rate to anonymized peers
//    - Show percentile rankings
//    - Motivate better habits
//
// 9. TAX ESTIMATION
//    - Track potential tax obligations
//    - Suggest amounts to set aside
//    - Generate tax reports
//
// 10. EMERGENCY FUND BUILDER
//     - Calculate needed emergency fund size
//     - Track progress toward goal
//     - Suggest automated savings plan
//
// ============================================================================
