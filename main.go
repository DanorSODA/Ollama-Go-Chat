package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
	"database/sql"
	_ "github.com/lib/pq"
)

const (
	MODEL_NAME = "tinyllama"
)

// Request structure for Ollama API
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// Response structure from Ollama API
type OllamaResponse struct {
	Response string `json:"response"`
}

// Add these new structs
type User struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Email       string    `json:"email"`
    Age         *int      `json:"age"`          // Pointer to allow NULL
    PhoneNumber string    `json:"phone_number"`
    Address     string    `json:"address"`
    Role        string    `json:"role"`
    IsActive    bool      `json:"is_active"`
    LastLogin   *time.Time `json:"last_login"`  // Pointer to allow NULL
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

var db *sql.DB

func initDB() error {
	connStr := "host=db port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	return db.Ping()
}

// CRUD operations
func createUser(name, email string, age *int, phone, address, role string) (*User, error) {
	var user User
	err := db.QueryRow(
		`INSERT INTO users (name, email, age, phone_number, address, role) 
		 VALUES ($1, $2, $3, $4, $5, $6) 
		 RETURNING id, name, email, age, phone_number, address, role, is_active, created_at`,
		name, email, age, phone, address, role,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Age, &user.PhoneNumber, 
		   &user.Address, &user.Role, &user.IsActive, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func getUser(id int) (*User, error) {
	var user User
	err := db.QueryRow(
		"SELECT id, name, email, created_at FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func updateUser(id int, name, email string) (*User, error) {
	var user User
	err := db.QueryRow(
		"UPDATE users SET name = $1, email = $2 WHERE id = $3 RETURNING id, name, email, created_at",
		name, email, id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func deleteUser(id int) error {
	result, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}
	return nil
}

// Function to send prompt to Ollama API and get a response
func generateText(prompt string) (string, error) {
	url := "http://localhost:11434/api/generate" // Ollama's local API endpoint

	requestBody := OllamaRequest{
		Model:  MODEL_NAME,
		Prompt: prompt,
		Stream: false,
	}

	// Convert request struct to JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	// Send HTTP request to Ollama API
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse the response
	var result OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Response, nil
}

func killExistingOllama() error {
	// Check if ollama is already running and kill it
	cmd := exec.Command("pkill", "ollama")
	if err := cmd.Run(); err != nil {
		// Ignore error as it might mean no process was found
		return nil
	}
	// Give it a moment to fully stop
	time.Sleep(500 * time.Millisecond)
	return nil
}

// Function to check if model exists and pull if needed
func ensureModelExists(modelName string) error {
	// Check if model exists using 'ollama list'
	cmd := exec.Command("ollama", "list")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error checking models: %v", err)
	}

	// If model name is not in the output, pull it
	if !bytes.Contains(output, []byte(modelName)) {
		fmt.Printf("Model %s not found. Pulling model...\n", modelName)
		cmd = exec.Command("ollama", "pull", modelName)
		cmd.Stdout = os.Stdout // Show pull progress
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error pulling model: %v", err)
		}
		fmt.Printf("Model %s successfully pulled!\n", modelName)
	}
	return nil
}

// Add this new function to search user by email
func getUserByEmail(email string) (*User, error) {
    var user User
    err := db.QueryRow(
        "SELECT id, name, email, age, phone_number, address, role, is_active, created_at FROM users WHERE email = $1",
        email,
    ).Scan(&user.ID, &user.Name, &user.Email, &user.Age, &user.PhoneNumber, &user.Address, &user.Role, &user.IsActive, &user.CreatedAt)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// Add this new function to get all users
func getAllUsers() ([]*User, error) {
    rows, err := db.Query("SELECT id, name, email, age, phone_number, address, role, is_active, created_at FROM users")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []*User
    for rows.Next() {
        var user User
        err := rows.Scan(
            &user.ID, 
            &user.Name, 
            &user.Email, 
            &user.Age, 
            &user.PhoneNumber, 
            &user.Address, 
            &user.Role, 
            &user.IsActive, 
            &user.CreatedAt,
        )
        if err != nil {
            return nil, err
        }
        users = append(users, &user)
    }
    return users, nil
}

// Update the parseAndExecuteDatabaseOperation function
func parseAndExecuteDatabaseOperation(response string) error {
    response = strings.ToLower(strings.TrimSpace(response))
    
    // Generic action detection
    isCreate := containsAny(response, []string{"create", "add", "new", "insert"})
    isRead := containsAny(response, []string{"get", "find", "show", "display", "search", "list", "fetch"})
    isUpdate := containsAny(response, []string{"update", "modify", "change", "edit"})
    isDelete := containsAny(response, []string{"delete", "remove", "drop"})
    
    // List all users (should be checked first)
    if isRead && containsAny(response, []string{"all", "everyone", "everybody", "users", "everything"}) {
        users, err := getAllUsers()
        if err != nil {
            return fmt.Errorf("failed to get users: %v", err)
        }
        
        if len(users) == 0 {
            fmt.Println("\nNo users found in the database.")
            return nil
        }
        
        fmt.Println("\nAll Users:")
        fmt.Println("----------------------------------------")
        for _, user := range users {
            fmt.Printf("ID: %d\n", user.ID)
            fmt.Printf("Name: %s\n", user.Name)
            fmt.Printf("Email: %s\n", user.Email)
            if user.Age != nil {
                fmt.Printf("Age: %d\n", *user.Age)
            }
            if user.PhoneNumber != "" {
                fmt.Printf("Phone: %s\n", user.PhoneNumber)
            }
            if user.Address != "" {
                fmt.Printf("Address: %s\n", user.Address)
            }
            if user.Role != "" {
                fmt.Printf("Role: %s\n", user.Role)
            }
            fmt.Printf("Active: %v\n", user.IsActive)
            fmt.Printf("Created: %v\n", user.CreatedAt.Format("2006-01-02 15:04:05"))
            fmt.Println("----------------------------------------")
        }
        return nil
    }
    
    // Create operation
    if isCreate {
        params := extractParameters(response)
        
        // Validate required fields for create operation only
        if params["name"] == "" || params["email"] == "" {
            return fmt.Errorf("name and email are required fields for creating a user")
        }
        
        // Convert age string to int pointer
        var agePtr *int
        if ageStr, ok := params["age"]; ok && ageStr != "" {
            age, err := strconv.Atoi(ageStr)
            if err == nil {
                agePtr = &age
            }
        }
        
        user, err := createUser(
            params["name"],
            params["email"],
            agePtr,
            params["phone"],
            params["address"],
            params["role"],
        )
        if err != nil {
            return fmt.Errorf("failed to create user: %v", err)
        }
        
        fmt.Printf("\nUser created successfully:\n")
        fmt.Printf("----------------------------------------\n")
        fmt.Printf("ID: %d\n", user.ID)
        fmt.Printf("Name: %s\n", user.Name)
        fmt.Printf("Email: %s\n", user.Email)
        if user.Age != nil {
            fmt.Printf("Age: %d\n", *user.Age)
        }
        if user.PhoneNumber != "" {
            fmt.Printf("Phone: %s\n", user.PhoneNumber)
        }
        if user.Role != "" {
            fmt.Printf("Role: %s\n", user.Role)
        }
        fmt.Printf("----------------------------------------\n")
        return nil
    }

    // Read operation
    if isRead {
        if strings.Contains(response, "email") {
            email := extractEmail(response)
            user, err := getUserByEmail(email)
            if err != nil {
                return err
            }
            fmt.Printf("Found user: %+v\n", user)
            return nil
        } else {
            id := extractUserId(response)
            user, err := getUser(id)
            if err != nil {
                return err
            }
            fmt.Printf("Found user: %+v\n", user)
            return nil
        }
    }

    // Update operation
    if isUpdate {
        params := extractParameters(response)
        id := extractUserId(response)
        if id == 0 {
            return fmt.Errorf("user ID is required for updates")
        }
        _, err := updateUser(id, params["name"], params["email"])
        if err != nil {
            return fmt.Errorf("failed to update user: %v", err)
        }
        fmt.Printf("User %d updated successfully\n", id)
        return nil
    }

    // Delete operation
    if isDelete {
        id := extractUserId(response)
        if id == 0 {
            return fmt.Errorf("user ID is required for deletion")
        }
        if err := deleteUser(id); err != nil {
            return fmt.Errorf("failed to delete user: %v", err)
        }
        fmt.Printf("User %d deleted successfully\n", id)
        return nil
    }

    return fmt.Errorf("could not understand the requested operation")
}

// Helper function to check if string contains any of the given phrases
func containsAny(s string, phrases []string) bool {
    for _, phrase := range phrases {
        if strings.Contains(s, phrase) {
            return true
        }
    }
    return false
}

func main() {
	// Initialize database connection
	if err := initDB(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Kill any existing Ollama processes
	if err := killExistingOllama(); err != nil {
		log.Printf("Warning: Could not kill existing Ollama process: %v", err)
	}

	// Start Ollama server with redirected output
	cmd := exec.Command("ollama", "serve")
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Start(); err != nil {
		log.Fatalf("Error starting Ollama server: %v", err)
	}

	// Give the server a moment to start
	time.Sleep(2 * time.Second)

	// Check and pull model if needed
	if err := ensureModelExists(MODEL_NAME); err != nil {
		log.Fatalf("Error ensuring model exists: %v", err)
	}

	fmt.Println("\nStarting AI-powered database management system")
	fmt.Println("You can interact with the database using natural language.")
	fmt.Println("Examples:")
	fmt.Println("- Add a new person named John with email john@example.com")
	fmt.Println("- Show me user with ID 1")
	fmt.Println("- Find the person with email john@example.com")
	fmt.Println("- Change user 1's email to new@example.com")
	fmt.Println("- Remove user with ID 2")
	fmt.Println("- Show all users")
	fmt.Println("Type 'quit' to exit")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\n> ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		if input == "quit" {
			break
		}

		// Construct a prompt that helps the AI understand the context
		prompt := `You are a database assistant. Based on the following request, generate a database operation.
		Available operations: create user, get user, update user, delete user, list all users.
		Request: %s
		Respond with the exact operation to perform, including all relevant parameters.
		For create operations, use format: create user with name=John Doe email=john@example.com
		For get operations, use format: get user with id=1 or get user with email=john@example.com
		For list operations, use format: list all users`

		response, err := generateText(prompt)
		if err != nil {
			fmt.Printf("Error generating response: %v\n", err)
			continue
		}

		// Parse and execute the AI's response
		if err := parseAndExecuteDatabaseOperation(response); err != nil {
			fmt.Printf("Error executing operation: %v\n", err)
			continue
		}
	}

	// Cleanup: kill the Ollama server process
	if err := cmd.Process.Kill(); err != nil {
		log.Printf("Error killing Ollama server: %v", err)
	}
}

// Update the extractParameters function with better regex patterns
func extractParameters(response string) map[string]string {
    params := make(map[string]string)
    
    // Extract name - look for "named" or "name" followed by text until next keyword
    nameMatch := regexp.MustCompile(`named?\s+([^,]+?)(?:\s+with|$|\s+age|\s+email|\s+phone|\s+role)`).FindStringSubmatch(response)
    if len(nameMatch) > 1 {
        params["name"] = strings.TrimSpace(nameMatch[1])
    }
    
    // Extract email - look for "email" followed by an email address
    emailMatch := regexp.MustCompile(`email\s+([^\s,]+)`).FindStringSubmatch(response)
    if len(emailMatch) > 1 {
        params["email"] = strings.TrimSpace(emailMatch[1])
    }
    
    // Extract age - look for "age" followed by numbers
    ageMatch := regexp.MustCompile(`age\s+(\d+)`).FindStringSubmatch(response)
    if len(ageMatch) > 1 {
        params["age"] = strings.TrimSpace(ageMatch[1])
    }
    
    // Extract phone - look for "phone" followed by numbers and possible symbols
    phoneMatch := regexp.MustCompile(`phone\s+([\d\+\-]+)`).FindStringSubmatch(response)
    if len(phoneMatch) > 1 {
        params["phone"] = strings.TrimSpace(phoneMatch[1])
    }
    
    // Extract role - look for "role" followed by quoted or unquoted text until end or next keyword
    roleMatch := regexp.MustCompile(`role\s+["']?([^"']+?)["']?(?:\s+|$)`).FindStringSubmatch(response)
    if len(roleMatch) > 1 {
        params["role"] = strings.TrimSpace(roleMatch[1])
    }
    
    // Extract address - look for "address" followed by quoted or unquoted text
    addressMatch := regexp.MustCompile(`address\s+["']?([^"']+?)["']?(?:\s+|$)`).FindStringSubmatch(response)
    if len(addressMatch) > 1 {
        params["address"] = strings.TrimSpace(addressMatch[1])
    }
    
    return params
}

// Helper function to extract user ID from AI response
func extractUserId(response string) int {
	idMatch := regexp.MustCompile(`id[=:]?\s*(\d+)`).FindStringSubmatch(response)
	if len(idMatch) > 1 {
		id, _ := strconv.Atoi(idMatch[1])
		return id
	}
	return 0
}

// Add this new helper function to extract email
func extractEmail(response string) string {
    emailMatch := regexp.MustCompile(`email[=:]?\s*["']?([^"',\s]+)["']?`).FindStringSubmatch(response)
    if len(emailMatch) > 1 {
        return strings.TrimSpace(emailMatch[1])
    }
    return ""
}
