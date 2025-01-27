# Ollama Go Chat with Database Management

A Go application that combines Ollama AI models with PostgreSQL database management, allowing natural language interactions for database operations. This application provides an AI-powered interface for performing CRUD operations on a user database.

## Features

- Natural language database operations (Create, Read, Update, Delete)
- PostgreSQL database integration
- AI-powered command interpretation
- Docker Compose setup with persistent storage
- Interactive command-line interface

## Requirements

- Docker and Docker Compose
- Git

## Quick Start with Docker Compose

1. Clone the repository:

```bash
git clone git@github.com:DanorSODA/Ollama-Go-Chat.git
cd Ollama-Go-Chat
```

2. Start the application using Docker Compose:

```bash
docker compose up --build
```

3. In a new terminal, attach to the running container to interact with the application:

```bash
docker exec -it ollama-go-demo-app-1 /bin/bash -c "./chat-app"
```

## Usage Examples

Once connected, you can interact with the database using natural language. Here are some example commands:

### Create Users

```
> Create a new user named John Doe with email john@example.com age 30 phone +1234567890 role developer
```

### Read Users

```
> Show all users
> Find user with ID 1
> Find the person with email john@example.com
```

### Update Users

```
> Update user 1's email to new@example.com
```

### Delete Users

```
> Delete user with ID 1
```

## Database Schema

The application uses a PostgreSQL database with the following user schema:

- `id`: Serial Primary Key
- `name`: String (required)
- `email`: String (required, unique)
- `age`: Integer (optional)
- `phone_number`: String (optional)
- `address`: Text (optional)
- `role`: String (optional)
- `is_active`: Boolean
- `created_at`: Timestamp
- `updated_at`: Timestamp

## Architecture

The application consists of three main components:

1. **Go Application**: Handles user input and database operations
2. **Ollama AI**: Interprets natural language commands
3. **PostgreSQL Database**: Stores user data

## Docker Components

- `app`: Main application container (Go + Ollama)
- `db`: PostgreSQL database container
- Persistent volume for database data

## Development

### Local Setup

1. Install Docker and Docker Compose
2. Clone the repository
3. Modify `docker-compose.yml` for any custom configurations
4. Run with `docker compose up --build`

### Customization

- Modify `init.sql` to change the database schema
- Update `main.go` to add new database operations
- Adjust the AI model in `const MODEL_NAME` (default: "tinyllama")

## Stopping the Application

To stop the application and clean up:

```bash
docker compose down
```

To remove all data including the database volume:

```bash
docker compose down -v
```
