DO $$
    BEGIN
        -- EXTENSIONS --
        CREATE EXTENSION IF NOT EXISTS pgcrypto;

        -- SEQUENCE --
        CREATE SEQUENCE IF NOT EXISTS order_id_seq
            START WITH 99999999
            INCREMENT BY 1;

        -- TABLES --
        CREATE TABLE IF NOT EXISTS users (
                                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                             updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                             id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                             name VARCHAR(100) NOT NULL,
                                             email VARCHAR(100) UNIQUE NOT NULL,
                                             password VARCHAR(255) NOT NULL,
                                             address VARCHAR(255),
                                             role VARCHAR(50) DEFAULT 'customer'
        );

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

        CREATE TABLE IF NOT EXISTS orders (
                                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                              updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                              id BIGINT PRIMARY KEY DEFAULT nextval('order_id_seq'),
                                              user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                                              total_price DECIMAL(10, 2) NOT NULL,
                                              status VARCHAR(50) NOT NULL
        );

        CREATE TABLE IF NOT EXISTS order_product (
                                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                                     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                                     order_id BIGINT NOT NULL,
                                                     product_id UUID NOT NULL,
                                                     PRIMARY KEY (order_id, product_id),
                                                     FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
                                                     FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
        );

        CREATE TABLE IF NOT EXISTS payments (
                                                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                                updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                                id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                                user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                                                order_id BIGINT REFERENCES orders(id) ON DELETE CASCADE,
                                                amount VARCHAR(200) NOT NULL,
                                                status VARCHAR(50) NOT NULL
        );

        -- DATA --
        INSERT INTO users (name, email, password, address, role) VALUES
                                                           ('John Doe', 'john.doe@example.com', '1233','123 Main St, Springfield', 'customer'),
                                                           ('Jane Smith', 'jane.smith@example.com', '1233','456 Oak St, Metropolis', 'customer'),
                                                           ('Alice Johnson', 'alice.johnson@example.com', '1233','789 Pine St, Gotham', 'admin'),
                                                           ('Bob Brown', 'bob.brown@example.com', '1233','101 Maple St, Star City', 'customer');

        INSERT INTO products (name, description, price, category, quantity) VALUES
                                                                                ('Laptop', 'High performance laptop', 999.99, 'Electronics', 10),
                                                                                ('Smartphone', 'Latest model smartphone', 699.99, 'Electronics', 25),
                                                                                ('Headphones', 'Noise-cancelling headphones', 199.99, 'Accessories', 50),
                                                                                ('Coffee Maker', 'Automatic coffee maker', 89.99, 'Home Appliances', 30);

        INSERT INTO orders (user_id, total_price, status) VALUES
                                                              ((SELECT id FROM users WHERE name = 'John Doe'), 1199.98, 'processing'),
                                                              ((SELECT id FROM users WHERE name = 'Jane Smith'), 199.99, 'shipped'),
                                                              ((SELECT id FROM users WHERE name = 'John Doe'), 699.99, 'delivered');

        INSERT INTO order_product (order_id, product_id) VALUES
                                                             ((SELECT id FROM orders WHERE total_price = 1199.98), (SELECT id FROM products WHERE name = 'Laptop')),
                                                             ((SELECT id FROM orders WHERE total_price = 1199.98), (SELECT id FROM products WHERE name = 'Smartphone')),
                                                             ((SELECT id FROM orders WHERE total_price = 199.99), (SELECT id FROM products WHERE name = 'Headphones')),
                                                             ((SELECT id FROM orders WHERE total_price = 699.99), (SELECT id FROM products WHERE name = 'Smartphone'));

        INSERT INTO payments (user_id, order_id, amount, status) VALUES
                                                                     ((SELECT id FROM users WHERE name = 'John Doe'), (SELECT id FROM orders WHERE total_price = 1199.98), '1199.98', 'completed'),
                                                                     ((SELECT id FROM users WHERE name = 'Jane Smith'), (SELECT id FROM orders WHERE total_price = 199.99), '199.99', 'completed'),
                                                                     ((SELECT id FROM users WHERE name = 'John Doe'), (SELECT id FROM orders WHERE total_price = 699.99), '699.99', 'pending');

        COMMIT;
    END $$;




-- DO $$
--     BEGIN
--         -- EXTENSIONS --
--         CREATE EXTENSION IF NOT EXISTS pgcrypto;
--
--         -- TABLES --
--         CREATE TABLE IF NOT EXISTS users (
--                                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--                                              updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--                                              id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--                                              name VARCHAR(100) NOT NULL,
--                                              email VARCHAR(100) UNIQUE NOT NULL,
--                                              address VARCHAR(255),
--                                              role VARCHAR(50) NOT NULL
--         );
--
--         -- DATA --
--         INSERT INTO users (name, email, address, role) VALUES
--                                                            ('John Doe', 'john.doe@example.com', '123 Main St, Springfield', 'customer'),
--                                                            ('Jane Smith', 'jane.smith@example.com', '456 Oak St, Metropolis', 'customer'),
--                                                            ('Alice Johnson', 'alice.johnson@example.com', '789 Pine St, Gotham', 'admin'),
--                                                            ('Bob Brown', 'bob.brown@example.com', '101 Maple St, Star City', 'customer');
--
--         COMMIT;
--     END $$;
