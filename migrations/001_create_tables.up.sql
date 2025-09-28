-- Crear tabla de productos 
CREATE TABLE products (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  price DECIMAL(10, 2) NOT NULL,
  description VARCHAR(500),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Crear tabla de citas 
CREATE TABLE appointments (
  id SERIAL PRIMARY KEY,
  cliente_name VARCHAR(100) NOT NULL,
  start_time TIMESTAMP NOT NULL,
  end_time TIMESTAMP NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

--  Creear tabla de relacion many-to-many entre citas y productos
CREATE TABLE appointment_product (
  appointment_id INTEGER REFERENCES appointments(id) ON DELETE CASCADE,
  product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
  PRIMARY KEY (appointment_id, product_id)
);


-- Crear indices para mejorar el rendmiento
CREATE INDEX idx_appointments_start_time ON appointments(start_time);
CREATE INDEX idx_appointments_end_time ON appointments(end_time);
CREATE INDEX idx_products_name ON products(name);

-- Insertar algunos productos de ejemplo
INSERT INTO products (name, price, description) VALUES
('Corte de Cabello', 15000.00, 'Corte tradicional'),
('Corte de Barba', 8000.00, 'Arreglo y perfilado de barba')













