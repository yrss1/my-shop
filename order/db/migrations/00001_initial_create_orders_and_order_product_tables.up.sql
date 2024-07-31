DO $$
    BEGIN
        -- SEQUENCE --
        CREATE SEQUENCE IF NOT EXISTS order_id_seq
            START WITH 99999999
            INCREMENT BY 1;

        -- TABLES --
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

        -- DATA --
        INSERT INTO orders (user_id, total_price, status) VALUES
                                                              ((SELECT id FROM users WHERE name = 'John Doe'), 1199.98, 'processing'),
                                                              ((SELECT id FROM users WHERE name = 'Jane Smith'), 199.99, 'shipped'),
                                                              ((SELECT id FROM users WHERE name = 'John Doe'), 699.99, 'delivered');

        INSERT INTO order_product (order_id, product_id) VALUES
                                                             ((SELECT id FROM orders WHERE total_price = 1199.98), (SELECT id FROM products WHERE name = 'Laptop')),
                                                             ((SELECT id FROM orders WHERE total_price = 1199.98), (SELECT id FROM products WHERE name = 'Smartphone')),
                                                             ((SELECT id FROM orders WHERE total_price = 199.99), (SELECT id FROM products WHERE name = 'Headphones')),
                                                             ((SELECT id FROM orders WHERE total_price = 699.99), (SELECT id FROM products WHERE name = 'Smartphone'));

        COMMIT;
    END $$;
