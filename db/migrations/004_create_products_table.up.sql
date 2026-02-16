CREATE TABLE products
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255)        NOT NULL,
    description TEXT,
    price       DECIMAL(10, 2)      NOT NULL,
    category_id INTEGER             NOT NULL REFERENCES categories (id) ON DELETE CASCADE,
    stock       INTEGER                  DEFAULT 0,
    sku         VARCHAR(100) UNIQUE NOT NULL,
    is_active   BOOLEAN                  DEFAULT TRUE,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP WITH TIME ZONE
);

CREATE index idx_products_category_id ON products (category_id);
CREATE index idx_products_sku ON products (sku);
CREATE index idx_products_is_active ON products (is_active);
CREATE index idx_products_deleted_at ON products (deleted_at);
