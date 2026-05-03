-- ============================================
-- POS SYSTEM - Script de creacion de base de datos
-- Motor: MySQL 9.6.0
-- Fecha: 2 de mayo de 2026
-- ============================================

-- Crear base de datos
CREATE DATABASE IF NOT EXISTS pos
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_general_ci;

USE pos;

-- ============================================
-- TABLA: usuarios
-- Descripcion: Usuarios del sistema (admin, vendedores)
-- ============================================
CREATE TABLE usuarios (
  id INT AUTO_INCREMENT PRIMARY KEY,
  usuario VARCHAR(100) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  nombre VARCHAR(150) NOT NULL,
  perfil ENUM('Administrador', 'Especial', 'Vendedor') NOT NULL DEFAULT 'Vendedor',
  foto VARCHAR(255) DEFAULT NULL,
  ultimo_login DATETIME DEFAULT NULL,
  estado TINYINT(1) NOT NULL DEFAULT 1,
  fecha_creacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- ============================================
-- TABLA: categorias
-- Descripcion: Categorias de productos
-- ============================================
CREATE TABLE categorias (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nombre VARCHAR(100) NOT NULL,
  descripcion TEXT DEFAULT NULL,
  estado TINYINT(1) NOT NULL DEFAULT 1,
  fecha_creacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- ============================================
-- TABLA: productos
-- Descripcion: Catalogo de productos
-- Relacion: categoria (N:1)
-- ============================================
CREATE TABLE productos (
  id INT AUTO_INCREMENT PRIMARY KEY,
  codigo VARCHAR(50) NOT NULL UNIQUE,
  codigo_barras VARCHAR(50) DEFAULT NULL,
  descripcion VARCHAR(255) NOT NULL,
  id_categoria INT NOT NULL,
  precio_compra DECIMAL(10,2) NOT NULL DEFAULT 0.00,
  precio_venta DECIMAL(10,2) NOT NULL DEFAULT 0.00,
  stock INT NOT NULL DEFAULT 0,
  stock_minimo INT NOT NULL DEFAULT 5,
  imagen VARCHAR(255) DEFAULT NULL,
  estado TINYINT(1) NOT NULL DEFAULT 1,
  fecha_creacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (id_categoria) REFERENCES categorias(id)
    ON UPDATE CASCADE ON DELETE RESTRICT
) ENGINE=InnoDB;

-- ============================================
-- TABLA: clientes
-- Descripcion: Registro de clientes
-- ============================================
CREATE TABLE clientes (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nombre VARCHAR(150) NOT NULL,
  documento VARCHAR(20) DEFAULT NULL,
  email VARCHAR(100) DEFAULT NULL,
  telefono VARCHAR(20) DEFAULT NULL,
  direccion VARCHAR(255) DEFAULT NULL,
  fecha_nacimiento DATE DEFAULT NULL,
  total_compras DECIMAL(10,2) NOT NULL DEFAULT 0.00,
  fecha_creacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- ============================================
-- TABLA: ventas
-- Descripcion: Registro de ventas/transacciones
-- Relaciones: usuario (N:1), cliente (N:1)
-- ============================================
CREATE TABLE ventas (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_usuario INT NOT NULL,
  id_cliente INT DEFAULT NULL,
  codigo_venta VARCHAR(50) NOT NULL UNIQUE,
  subtotal DECIMAL(10,2) NOT NULL DEFAULT 0.00,
  impuesto DECIMAL(10,2) NOT NULL DEFAULT 0.00,
  descuento DECIMAL(10,2) NOT NULL DEFAULT 0.00,
  total DECIMAL(10,2) NOT NULL DEFAULT 0.00,
  metodo_pago ENUM('Efectivo', 'Tarjeta', 'Transferencia', 'Mixto') NOT NULL DEFAULT 'Efectivo',
  estado ENUM('completada', 'cancelada', 'pendiente') NOT NULL DEFAULT 'completada',
  fecha TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (id_usuario) REFERENCES usuarios(id)
    ON UPDATE CASCADE ON DELETE RESTRICT,
  FOREIGN KEY (id_cliente) REFERENCES clientes(id)
    ON UPDATE CASCADE ON DELETE SET NULL
) ENGINE=InnoDB;

-- ============================================
-- TABLA: detalle_ventas
-- Descripcion: Items individuales de cada venta
-- Relaciones: venta (N:1), producto (N:1)
-- ============================================
CREATE TABLE detalle_ventas (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_venta INT NOT NULL,
  id_producto INT NOT NULL,
  cantidad INT NOT NULL,
  precio_unitario DECIMAL(10,2) NOT NULL,
  descuento DECIMAL(10,2) NOT NULL DEFAULT 0.00,
  subtotal DECIMAL(10,2) NOT NULL,
  FOREIGN KEY (id_venta) REFERENCES ventas(id)
    ON UPDATE CASCADE ON DELETE CASCADE,
  FOREIGN KEY (id_producto) REFERENCES productos(id)
    ON UPDATE CASCADE ON DELETE RESTRICT
) ENGINE=InnoDB;

-- ============================================
-- INDICES adicionales para optimizar consultas
-- ============================================
CREATE INDEX idx_productos_categoria ON productos(id_categoria);
CREATE INDEX idx_productos_codigo_barras ON productos(codigo_barras);
CREATE INDEX idx_ventas_fecha ON ventas(fecha);
CREATE INDEX idx_ventas_usuario ON ventas(id_usuario);
CREATE INDEX idx_ventas_cliente ON ventas(id_cliente);
CREATE INDEX idx_ventas_estado ON ventas(estado);
CREATE INDEX idx_detalle_venta ON detalle_ventas(id_venta);
CREATE INDEX idx_detalle_producto ON detalle_ventas(id_producto);
CREATE INDEX idx_clientes_documento ON clientes(documento);

-- ============================================
-- DATOS INICIALES
-- ============================================

-- Usuario administrador (password: admin123)
INSERT INTO usuarios (usuario, password, nombre, perfil, estado)
VALUES ('admin', '$2y$12$0FTpJzLnCtLjz0SQ1r1h6OwqYAEnbEG0C.LHijJ4pEYnewrWPqmNi', 'Administrador del Sistema', 'Administrador', 1);

-- Categorias de ejemplo
INSERT INTO categorias (nombre, descripcion) VALUES
  ('Bebidas', 'Refrescos, jugos, aguas y bebidas en general'),
  ('Alimentos', 'Snacks, comida preparada y productos alimenticios'),
  ('Limpieza', 'Productos de limpieza para el hogar'),
  ('Cuidado Personal', 'Higiene y cuidado personal'),
  ('Varios', 'Productos varios y miscelaneos');

-- Cliente generico (para ventas sin cliente registrado)
INSERT INTO clientes (nombre, documento)
VALUES ('Publico General', '0000000000');
