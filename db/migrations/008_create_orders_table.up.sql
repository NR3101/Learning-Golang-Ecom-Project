CREATE TYPE order_status AS ENUM ('pending', 'confirmed', 'shipped', 'delivered', 'cancelled');

CREATE TABLE orders
(
    id           SERIAL PRIMARY KEY,
    user_id      INTEGER        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    total_amount DECIMAL(10, 2) NOT NULL,
    status       order_status             DEFAULT 'pending',
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMP WITH TIME ZONE
);

CREATE index idx_orders_user_id ON orders (user_id);
CREATE index idx_orders_status ON orders (status);
CREATE index idx_orders_deleted_at ON orders (deleted_at);