CREATE TABLE IF NOT EXISTS orders (
    id VARCHAR(36) PRIMARY KEY,
    price NUMERIC,
    tax NUMERIC,
    final_price NUMERIC
)