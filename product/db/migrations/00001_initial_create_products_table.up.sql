DO $$
    BEGIN
        -- TABLES --
        CREATE TABLE IF NOT EXISTS products (
                                                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                                updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                                id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                                name VARCHAR(200) UNIQUE NOT NULL,
                                                description TEXT,
                                                price DECIMAL(10, 2) NOT NULL,
                                                category VARCHAR(100),
                                                quantity INTEGER NOT NULL
        );

        -- DATA --
        INSERT INTO products (name, description, price, category, quantity) VALUES
                                                                                ('Laptop', 'High performance laptop', 999.99, 'Electronics', 10),
                                                                                ('Smartphone', 'Latest model smartphone', 699.99, 'Electronics', 25),
                                                                                ('Headphones', 'Noise-cancelling headphones', 199.99, 'Accessories', 50),
                                                                                ('Coffee Maker', 'Automatic coffee maker', 89.99, 'Home Appliances', 30);

        COMMIT;
    END $$;
