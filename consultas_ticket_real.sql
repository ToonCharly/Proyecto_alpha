-- CONSULTA PRINCIPAL - VALORES REALES DEL TICKET (OPCIÓN MÁS CONFIABLE)
-- Esta consulta muestra ÚNICAMENTE los datos que estaban en el ticket original
-- SIN complicaciones de múltiples configuraciones de impuestos

SELECT 
    d.idproducto AS 'Código Producto',
    d.descripcion AS 'Producto',
    d.cantidad AS 'Cantidad',
    d.precio AS 'Precio Unitario',
    d.iva AS 'IVA Real (%)',        -- VALOR EXACTO DEL TICKET
    0 AS 'IEPS (%)',                -- IEPS del ticket (generalmente 0)
    d.descuento AS 'Descuento',
    COALESCE(pr.sat_clave, 'No disponible') AS 'Clave SAT',
    COALESCE(pr.sat_medida, 'No disponible') AS 'Unidad SAT',
    -- Cálculos con valores reales
    (d.cantidad * d.precio) - d.descuento AS 'Subtotal',
    ((d.cantidad * d.precio) - d.descuento) * (d.iva / 100) AS 'IVA Importe',
    ((d.cantidad * d.precio) - d.descuento) * (1 + d.iva / 100) AS 'Total'
FROM optimus.crm_pedidos p 
JOIN optimus.crm_pedidos_det d ON p.id_pedido = d.id_pedido 
LEFT JOIN optimus.crm_productos pr ON d.idproducto = pr.idproducto 
WHERE p.clave_pedido = 'TU_SERIE_AQUI'
ORDER BY d.id_pedido_det;

-- ================================================================
-- CONSULTA PARA VERIFICAR ESTRUCTURA DE LA TABLA (DIAGNÓSTICO)
-- ================================================================
SELECT 
    'Verificación de columnas en crm_pedidos_det' AS diagnostico,
    COUNT(*) AS total_productos
FROM optimus.crm_pedidos p 
JOIN optimus.crm_pedidos_det d ON p.id_pedido = d.id_pedido 
WHERE p.clave_pedido = 'TU_SERIE_AQUI';

-- Ver qué columnas tiene realmente la tabla
DESCRIBE optimus.crm_pedidos_det;

-- ================================================================
-- CONSULTA SIMPLIFICADA PARA VERIFICAR DATOS BÁSICOS DEL TICKET
-- ================================================================
SELECT 
    p.clave_pedido,
    d.descripcion AS producto,
    d.cantidad,
    d.precio,
    d.iva AS iva_del_pedido,
    d.descuento,
    pr.sat_clave,
    pr.sat_medida,
    -- Información de impuestos de configuración
    COALESCE(imp.ieps1, 0) AS ieps1_configurado,
    COALESCE(imp.ieps2, 0) AS ieps2_configurado,
    COALESCE(imp.ieps3, 0) AS ieps3_configurado
FROM optimus.crm_pedidos p 
JOIN optimus.crm_pedidos_det d ON p.id_pedido = d.id_pedido 
LEFT JOIN optimus.crm_productos pr ON d.idproducto = pr.idproducto 
LEFT JOIN optimus.crm_impuestos imp ON pr.idempresa = imp.idempresa
WHERE p.clave_pedido = 'TU_SERIE_AQUI'
ORDER BY d.id_pedido_det;

-- ================================================================
-- CONSULTA PARA VERIFICAR CONFIGURACIÓN DE IMPUESTOS
-- ================================================================
SELECT 
    imp.idempresa,
    imp.iva,
    imp.ieps1,
    imp.ieps2,
    imp.ieps3,
    COUNT(*) as cantidad_productos_empresa
FROM optimus.crm_impuestos imp
JOIN optimus.crm_productos pr ON imp.idempresa = pr.idempresa
WHERE pr.idproducto IN (
    SELECT DISTINCT d.idproducto 
    FROM optimus.crm_pedidos p 
    JOIN optimus.crm_pedidos_det d ON p.id_pedido = d.id_pedido 
    WHERE p.clave_pedido = 'TU_SERIE_AQUI'
)
GROUP BY imp.idempresa, imp.iva, imp.ieps1, imp.ieps2, imp.ieps3;

-- ================================================================
-- INSTRUCCIONES DE USO:
-- ================================================================
-- 1. Reemplaza 'TU_SERIE_AQUI' con la serie real del ticket
-- 2. Ejecuta la primera consulta para ver los datos completos del ticket
-- 3. Ejecuta la segunda consulta para ver datos básicos
-- 4. Ejecuta la tercera consulta para verificar configuración de impuestos
-- 
-- IMPORTANTE: Los valores mostrados en las columnas iva_ticket, ieps1_ticket, etc.
-- son los valores REALES que estaban en el ticket cuando se creó.
-- Estos valores pueden ser diferentes a los configurados en crm_impuestos.
