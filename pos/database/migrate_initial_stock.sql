-- ==============================================
-- MIGRATION: Backfill initial stock as inventory entries
--
-- Context: ctrCreateProduct() previously did not register an inventory
-- movement when a product was created with stock > 0. This migration
-- creates one "Ajuste positivo" entry for every product that currently
-- has stock > 0 and has NO existing movement in movimientos_inventario.
--
-- Products that already have movements (e.g. CER-002, CER-003) are
-- intentionally excluded — their audit trail is already correct.
--
-- Run once. Re-running is safe: the WHERE NOT IN guard prevents duplicates.
-- ==============================================

INSERT INTO movimientos_inventario
    (id_producto, id_usuario, tipo, motivo, cantidad, observaciones, id_referencia)
SELECT
    p.id,
    1,                                          -- Administrador (id=1)
    'entrada',
    'Ajuste positivo',
    p.stock,
    'Migración: stock inicial del producto',
    NULL
FROM productos p
WHERE p.stock > 0
  AND p.id NOT IN (
      SELECT DISTINCT id_producto FROM movimientos_inventario
  );
