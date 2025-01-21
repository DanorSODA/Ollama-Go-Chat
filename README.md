# Ollama Go Chat

A simple chat application built with Go that interfaces with Ollama AI models. This application provides a lightweight command-line interface for interacting with various Ollama models, with built-in response time tracking and Docker support.

## Features

- Interactive command-line chat interface
- Automatic model downloading and management
- Response time tracking for performance monitoring
- Support for multiple Ollama models

## Requirements

- Go 1.16 or higher
- Docker (for containerized usage)
- Ollama (for local development)

## Quick Start with Docker

1. Clone the repository:

```bash
git clone git@github.com:DanorSODA/Ollama-Go-Chat.git
cd Ollama-Go-Chat
```

2. Build the Docker image:

```bash
docker build -t ollama-go-chat .
```

3. Run the container:

```bash
docker run -it ollama-go-chat
```

## Local Development

1. Install Ollama from https://ollama.ai

2. Clone the repository:

```bash
git clone git@github.com:DanorSODA/Ollama-Go-Chat.git
cd Ollama-Go-Chat
```

3. Run the Go application:

```bash
go run main.go
```

## Usage

Once running, you can:

- Enter prompts at the `>` prompt
- See response times for each interaction
- Type `quit` to exit the application

Example interaction:

```
Starting chat with tinyllama model
Enter your prompt (or 'quit' to exit):
> What is Go programming?

Generating response...

----------------------------------------
AI Response:
[Model response will appear here]

Response time: 1.23 seconds
----------------------------------------
```

## Configuration

The application uses TinyLlama by default for optimal performance on most systems. You can modify the model by changing the `MODEL_NAME` constant in `main.go`:

```go
const (
    MODEL_NAME = "tinyllama"  // Change to your preferred model
)
```

Available models include:

- `tinyllama` (Default, fastest)
- `orca-mini` (Light and fast)
- `phi` (Balanced)
- `neural-chat` (More capable)
- `mistral` (Most capable)

## Performance Considerations

- Model response times vary based on your hardware capabilities
