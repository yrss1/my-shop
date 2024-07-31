DO $$
    BEGIN
        -- TABLES --
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
        INSERT INTO payments (user_id, order_id, amount, status) VALUES
                                                                     ((SELECT id FROM users WHERE name = 'John Doe'), (SELECT id FROM orders WHERE total_price = 1199.98), '1199.98', 'completed'),
                                                                     ((SELECT id FROM users WHERE name = 'Jane Smith'), (SELECT id FROM orders WHERE total_price = 199.99), '199.99', 'completed'),
                                                                     ((SELECT id FROM users WHERE name = 'John Doe'), (SELECT id FROM orders WHERE total_price = 699.99), '699.99', 'pending');

        COMMIT;
    END $$;
