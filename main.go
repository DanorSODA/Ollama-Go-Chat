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
	"time"
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

func main() {
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

	fmt.Printf("\nStarting chat with %s model\n", MODEL_NAME)
	fmt.Println("Enter your prompt (or 'quit' to exit):")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\n> ")
		if !scanner.Scan() {
			break
		}

		prompt := scanner.Text()
		if prompt == "quit" {
			break
		}

		fmt.Println("\nGenerating response...\n")
		
		// Start timer
		start := time.Now()
		
		response, err := generateText(prompt)
		if err != nil {
			log.Printf("Error generating text: %v", err)
			continue
		}

		// Calculate elapsed time
		elapsed := time.Since(start)

		fmt.Println("----------------------------------------")
		fmt.Println("AI Response:")
		fmt.Println(response)
		fmt.Printf("\nResponse time: %.2f seconds\n", elapsed.Seconds())
		fmt.Println("----------------------------------------")
	}

	// Cleanup: kill the Ollama server process
	if err := cmd.Process.Kill(); err != nil {
		log.Printf("Error killing Ollama server: %v", err)
	}
}
