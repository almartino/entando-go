CREATE TABLE IF NOT EXISTS conferences(
    id SERIAL PRIMARY KEY,
    name VARCHAR(50),
    location VARCHAR(255)
);