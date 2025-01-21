# Start with Ubuntu as base image
FROM --platform=linux/amd64 ubuntu:22.04

# Install required dependencies
RUN apt-get update && apt-get install -y \
    curl \
    wget \
    git \
    golang \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Install Ollama with proper verification
RUN curl -fsSL https://github.com/ollama/ollama/releases/download/v0.1.29/ollama-linux-amd64 -o /usr/bin/ollama \
    && chmod +x /usr/bin/ollama \
    && /usr/bin/ollama --version

# Create app directory
WORKDIR /app

# Copy Go source code
COPY main.go .

# Build the Go application
RUN go build -o chat-app main.go

# Create startup script that pulls the model first
RUN echo '#!/bin/bash\n\
/usr/bin/ollama serve &\n\
sleep 5\n\
# Pull the model before starting the app\n\
/usr/bin/ollama pull tinyllama\n\
./chat-app' > /app/start.sh && chmod +x /app/start.sh

# Expose Ollama's default port
EXPOSE 11434

# Set the startup script as the entry point
ENTRYPOINT ["/bin/bash", "/app/start.sh"] 