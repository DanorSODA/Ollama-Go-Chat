CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    age INT,
    phone_number VARCHAR(20),
    address TEXT,
    role VARCHAR(50) DEFAULT 'user',
    is_active BOOLEAN DEFAULT true,
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create an index on email for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Insert some sample data with ON CONFLICT DO NOTHING
INSERT INTO users (name, email, age, phone_number, address, role) 
VALUES 
    ('John Doe', 'john@example.com', 30, '+1234567890', '123 Main St, City', 'admin'),
    ('Jane Smith', 'jane@example.com', 25, '+0987654321', '456 Oak Ave, Town', 'user')
ON CONFLICT (email) DO NOTHING; 