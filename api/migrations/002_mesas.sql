-- Migration: Create mesas table and add id_mesa to ventas
-- Run against the POS database after 001_initial.sql

-- New table: mesas
CREATE TABLE mesas (
  id INT AUTO_INCREMENT PRIMARY KEY,
  numero INT NOT NULL UNIQUE,
  nombre VARCHAR(100) DEFAULT NULL,
  capacidad INT NOT NULL DEFAULT 4,
  estado ENUM('libre','ocupada','reservada','pagando') NOT NULL DEFAULT 'libre',
  id_venta_activa INT DEFAULT NULL,
  id_mesero INT DEFAULT NULL,
  fecha_creacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  fecha_actualizacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_mesas_estado (estado),
  FOREIGN KEY (id_venta_activa) REFERENCES ventas(id) ON UPDATE CASCADE ON DELETE SET NULL,
  FOREIGN KEY (id_mesero) REFERENCES usuarios(id) ON UPDATE CASCADE ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Add id_mesa column to ventas
ALTER TABLE ventas ADD COLUMN id_mesa INT DEFAULT NULL;
ALTER TABLE ventas ADD CONSTRAINT fk_ventas_mesa FOREIGN KEY (id_mesa) REFERENCES mesas(id)
  ON UPDATE CASCADE ON DELETE SET NULL;
CREATE INDEX idx_ventas_mesa ON ventas(id_mesa);

-- Add 'abierta' status to ventas.estado
ALTER TABLE ventas MODIFY COLUMN estado
  ENUM('completada','cancelada','pendiente','abierta') NOT NULL DEFAULT 'completada';
