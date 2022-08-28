-- migrate:up
CREATE TABLE IF NOT EXISTS tickets (
  ID SERIAL PRIMARY KEY,
  body TEXT,
  status VARCHAR(255) NOT NULL,
  error_details TEXT,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

CREATE INDEX tickets_status_idx ON tickets (status);

CREATE TABLE IF NOT EXISTS products (
  product_id TEXT NOT NULL PRIMARY KEY,
  name TEXT NOT NULL,
  -- maybe we should use an integer for prices, here for simplicity I chosed double
  price DOUBLE PRECISION,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

CREATE INDEX products_id_idx ON products (product_id);

-- migrate:down
DROP TABLE tickets;
DROP TABLE products;
