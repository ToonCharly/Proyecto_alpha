-- VERIFICAR ESTRUCTURA DE TABLAS PARA IDENTIFICAR COLUMNAS DISPONIBLES

-- 1. Verificar estructura de la tabla crm_pedidos_det
DESCRIBE optimus.crm_pedidos_det;

-- 2. Verificar estructura de la tabla crm_pedidos
DESCRIBE optimus.crm_pedidos;

-- 3. Verificar estructura de la tabla crm_productos
DESCRIBE optimus.crm_productos;

-- 4. Verificar estructura de la tabla crm_impuestos
DESCRIBE optimus.crm_impuestos;

-- 5. Consulta alternativa para ver las columnas (en caso de que DESCRIBE no funcione)
SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE, COLUMN_DEFAULT 
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_SCHEMA = 'optimus' AND TABLE_NAME = 'crm_pedidos_det'
ORDER BY ORDINAL_POSITION;

-- 6. Ver datos de ejemplo de crm_pedidos_det para entender la estructura
SELECT * FROM optimus.crm_pedidos_det LIMIT 5;

-- 7. Ver datos de ejemplo de crm_impuestos para entender la estructura
SELECT * FROM optimus.crm_impuestos LIMIT 5;
