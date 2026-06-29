-- =============================================
-- TABLA: movimientos_inventario
-- Módulo de inventario - registro de entradas y salidas de productos
-- =============================================

CREATE TABLE IF NOT EXISTS `movimientos_inventario` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `id_producto` INT(11) NOT NULL,
  `id_usuario` INT(11) NOT NULL,
  `tipo` ENUM('entrada','salida') NOT NULL,
  `motivo` VARCHAR(50) NOT NULL,
  `cantidad` INT(11) NOT NULL,
  `observaciones` TEXT NULL,
  `id_referencia` INT(11) NULL DEFAULT NULL,
  `fecha` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_producto` (`id_producto`),
  INDEX `idx_fecha` (`fecha`),
  INDEX `idx_referencia` (`id_referencia`),
  INDEX `idx_tipo` (`tipo`),
  CONSTRAINT `fk_mov_producto` FOREIGN KEY (`id_producto`) REFERENCES `productos` (`id`),
  CONSTRAINT `fk_mov_usuario` FOREIGN KEY (`id_usuario`) REFERENCES `usuarios` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
