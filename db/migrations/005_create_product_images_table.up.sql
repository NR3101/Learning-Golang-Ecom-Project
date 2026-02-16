CREATE TABLE product_images
(
    id         SERIAL PRIMARY KEY,
    product_id INTEGER      NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    url        VARCHAR(255) NOT NULL,
    alt_text   VARCHAR(255),
    is_primary BOOLEAN                  DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE index idx_product_images_product_id ON product_images (product_id);
CREATE index idx_product_images_is_primary ON product_images (is_primary);
CREATE index idx_product_images_deleted_at ON product_images (deleted_at);