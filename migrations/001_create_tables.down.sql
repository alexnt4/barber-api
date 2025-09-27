--  Elminar índices 
DROP INDEX IF EXISTS idx_products_name;
DROP INDEX IF EXISTS idx_appointments_end_time;
DROP INDEX IF EXISTS idx_appointmets_start_time;
-- Eliminar tabla de relacion many-to-many
DROP TABLE IF EXISTS appoinments_products;
-- Eliminar tabla de citas
DROP TABLE IF EXISTS appointmets;
-- Eliminar tabla de productos
DROP TABLE IF EXISTS products;
