# Start with Ubuntu as base image
FROM --platform=linux/amd64 ubuntu:22.04

# Install required dependencies including PostgreSQL client
RUN apt-get update && apt-get install -y \
    curl \
    wget \
    git \
    golang \
    ca-certificates \
    postgresql-client \
    && rm -rf /var/lib/apt/lists/*

# Install Ollama with proper verification
RUN curl -fsSL https://github.com/ollama/ollama/releases/download/v0.1.29/ollama-linux-amd64 -o /usr/bin/ollama \
    && chmod +x /usr/bin/ollama \
    && /usr/bin/ollama --version

# Create app directory
WORKDIR /app

# Copy Go source code
COPY main.go .
COPY init.sql .

# Initialize Go module and get dependencies
RUN go mod init chat-app \
    && go mod tidy \
    && go get github.com/lib/pq

# Build the Go application
RUN go build -o chat-app main.go

# Create startup script with better process handling
RUN echo '#!/bin/bash\n\
# Start Ollama in the background\n\
/usr/bin/ollama serve &\n\
OLLAMA_PID=$!\n\
\n\
# Wait for Ollama to start\n\
sleep 5\n\
\n\
# Wait for PostgreSQL to be ready\n\
while ! pg_isready -h db -p 5432 -U postgres; do\n\
    echo "Waiting for PostgreSQL..."\n\
    sleep 1\n\
done\n\
\n\
# Initialize database\n\
PGPASSWORD=postgres psql -h db -U postgres -d postgres -f /app/init.sql\n\
\n\
# Pull the model and wait for completion\n\
echo "Pulling tinyllama model..."\n\
/usr/bin/ollama pull tinyllama\n\
\n\
# Start the chat application\n\
echo "Starting chat application..."\n\
./chat-app\n\
\n\
# If chat-app exits, kill Ollama\n\
kill $OLLAMA_PID\n\
wait $OLLAMA_PID' > /app/start.sh && chmod +x /app/start.sh

# Expose Ollama's default port
EXPOSE 11434

# Set the startup script as the entry point
ENTRYPOINT ["/bin/bash", "/app/start.sh"] 