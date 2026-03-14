-- Add tsvector column for full-text search on products
ALTER TABLE products
    ADD COLUMN search_vector tsvector;

-- Create a trigger function to update the search_vector column
CREATE OR REPLACE FUNCTION update_product_search_vector() RETURNS trigger AS
$$
BEGIN
    NEW.search_vector :=
            setweight(to_tsvector('english', coalesce(NEW.name, '')), 'A') ||
            setweight(to_tsvector('english', coalesce(NEW.description, '')), 'B') ||
            setweight(to_tsvector('english', coalesce(NEW.sku, '')), 'C');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create a trigger to call the function on insert and update
CREATE TRIGGER product_search_vector_update
    BEFORE INSERT OR UPDATE
    ON products
    FOR EACH ROW
EXECUTE FUNCTION update_product_search_vector();

-- Update existing rows to populate the search_vector column
UPDATE products
SET search_vector =
        setweight(to_tsvector('english', coalesce(name, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(description, '')), 'B') ||
        setweight(to_tsvector('english', coalesce(sku, '')), 'C');

-- Create a GIN index on the search_vector column for faster full-text search
CREATE INDEX idx_products_search_vector ON products USING GIN (search_vector);

-- Add comment to the search_vector column
COMMENT ON COLUMN products.search_vector IS 'Full-text search vector for products table with weighted fields: name (A), description (B), sku (C)';