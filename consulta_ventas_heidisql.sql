-- ============================================================================
-- CONSULTA PARA VER DATOS DE VENTAS EN HEIDISQL
-- Reemplaza '19a5a142aca218a7ff5deec2831d6834' por tu serie real
-- ============================================================================

-- CONSULTA PRINCIPAL - DATOS COMPLETOS CON DEPURACIÓN
SELECT 
    p.id_pedido, 
    p.clave_pedido, 
    d.descripcion AS producto, 
    d.cantidad, 
    d.precio, 
    COALESCE(d.precio_o, d.precio) AS precio_o,           
    d.iva AS iva_pedido,                
    d.descuento,
    d.idproducto AS codigo_producto, 
    pr.clave AS categoria_producto,
    pr.sat_clave,           
    pr.sat_medida,
    pr.idempresa,
    -- Impuestos desde crm_impuestos (usando MAX para evitar duplicados)
    COALESCE(MAX(imp.iva), 0) AS iva_config,
    COALESCE(MAX(imp.ieps1), 0) AS ieps1_config,
    COALESCE(MAX(imp.ieps2), 0) AS ieps2_config,
    COALESCE(MAX(imp.ieps3), 0) AS ieps3_config,
    -- Conteo de configuraciones para el diagnóstico
    COUNT(imp.idiva) AS configs_count,
    -- Información adicional para depuración
    pr.idempresa AS empresa_producto,
    GROUP_CONCAT(DISTINCT imp.idempresa) AS empresas_impuestos,
    -- Diagnóstico mejorado
    CASE 
        WHEN pr.idproducto IS NULL THEN 'Sin producto en crm_productos'
        WHEN pr.sat_clave IS NULL OR pr.sat_clave = '' OR pr.sat_clave = '0' THEN 'Sin clave SAT'
        WHEN pr.sat_medida IS NULL OR pr.sat_medida = '' THEN 'Sin unidad SAT'
        ELSE 'Datos completos'
    END AS diagnostico_config,
    -- Determinar origen de impuestos
    CASE 
        WHEN COUNT(imp.idiva) > 0 THEN 'crm_impuestos'
        ELSE 'ticket_original'
    END AS origen_impuestos,
    -- IVA final que se usará
    CASE 
        WHEN COUNT(imp.idiva) > 0 THEN COALESCE(MAX(imp.iva), 0)
        ELSE d.iva
    END AS iva_final,
    -- Cálculos finales
    (d.cantidad * d.precio) - d.descuento AS subtotal,
    ((d.cantidad * d.precio) - d.descuento) * 
    (CASE 
        WHEN COUNT(imp.idiva) > 0 THEN COALESCE(MAX(imp.iva), 0)
        ELSE d.iva
    END / 100.0) AS iva_importe,
    ((d.cantidad * d.precio) - d.descuento) * 
    (1 + (CASE 
        WHEN COUNT(imp.idiva) > 0 THEN COALESCE(MAX(imp.iva), 0)
        ELSE d.iva
    END / 100.0)) AS total_final
FROM optimus.crm_pedidos p 
JOIN optimus.crm_pedidos_det d ON p.id_pedido = d.id_pedido 
LEFT JOIN optimus.crm_productos pr ON d.idproducto = pr.idproducto 
LEFT JOIN optimus.crm_impuestos imp ON pr.idempresa = imp.idempresa
WHERE p.clave_pedido = '19a5a142aca218a7ff5deec2831d6834'
GROUP BY p.id_pedido, p.clave_pedido, d.descripcion, d.cantidad, d.precio, 
         d.precio_o, d.iva, d.descuento, d.idproducto, pr.clave, 
         pr.sat_clave, pr.sat_medida, pr.idempresa
ORDER BY d.idproducto;

-- ============================================================================
-- CONSULTA SIMPLIFICADA - SOLO DATOS BÁSICOS
-- ============================================================================
SELECT 
    d.idproducto AS 'Código',
    d.descripcion AS 'Producto',
    d.cantidad AS 'Cant',
    d.precio AS 'Precio',
    d.iva AS 'IVA Ticket',
    COALESCE(MAX(imp.iva), 0) AS 'IVA Config',
    COUNT(imp.idiva) AS 'Configs',
    pr.sat_clave AS 'Clave SAT',
    pr.sat_medida AS 'Unidad SAT',
    CASE 
        WHEN COUNT(imp.idiva) > 0 THEN 'Configurado'
        ELSE 'Sin configurar'
    END AS 'Estado'
FROM optimus.crm_pedidos p 
JOIN optimus.crm_pedidos_det d ON p.id_pedido = d.id_pedido 
LEFT JOIN optimus.crm_productos pr ON d.idproducto = pr.idproducto 
LEFT JOIN optimus.crm_impuestos imp ON pr.idempresa = imp.idempresa
WHERE p.clave_pedido = '19a5a142aca218a7ff5deec2831d6834'
GROUP BY d.idproducto, d.descripcion, d.cantidad, d.precio, d.iva, pr.sat_clave, pr.sat_medida
ORDER BY d.idproducto;

-- ============================================================================
-- VERIFICAR SI EL PEDIDO EXISTE
-- ============================================================================
SELECT 
    'VERIFICACIÓN DE PEDIDO' as tipo,
    COUNT(*) as productos_encontrados,
    p.clave_pedido
FROM optimus.crm_pedidos p 
JOIN optimus.crm_pedidos_det d ON p.id_pedido = d.id_pedido 
WHERE p.clave_pedido = '19a5a142aca218a7ff5deec2831d6834'
GROUP BY p.clave_pedido;

-- ============================================================================
-- VER CONFIGURACIONES DE IMPUESTOS DISPONIBLES
-- ============================================================================
SELECT 
    'CONFIGURACIONES DISPONIBLES' as tipo,
    imp.idempresa,
    imp.iva,
    imp.ieps1,
    imp.ieps2,
    imp.ieps3,
    COUNT(pr.idproducto) as productos_con_esta_empresa
FROM optimus.crm_impuestos imp
LEFT JOIN optimus.crm_productos pr ON imp.idempresa = pr.idempresa
GROUP BY imp.idempresa, imp.iva, imp.ieps1, imp.ieps2, imp.ieps3
ORDER BY imp.idempresa;

-- ============================================================================
-- CONSULTA COMPLETA CON IEPS - TODOS LOS IMPUESTOS
-- ============================================================================
SELECT 
    p.id_pedido, 
    p.clave_pedido, 
    d.descripcion AS producto, 
    d.cantidad, 
    d.precio, 
    COALESCE(d.precio_o, d.precio) AS precio_o,           
    d.iva AS iva_pedido,                
    d.descuento,
    d.idproducto AS codigo_producto, 
    pr.clave AS categoria_producto,
    pr.sat_clave,           
    pr.sat_medida,
    pr.idempresa,
    -- Impuestos desde crm_impuestos
    COALESCE(MAX(imp.iva), 0) AS iva_config,
    COALESCE(MAX(imp.ieps1), 0) AS ieps1_config,
    COALESCE(MAX(imp.ieps2), 0) AS ieps2_config,
    COALESCE(MAX(imp.ieps3), 0) AS ieps3_config,
    -- Conteo de configuraciones
    COUNT(imp.idiva) AS configs_count,
    -- Información de depuración
    pr.idempresa AS empresa_producto,
    GROUP_CONCAT(DISTINCT imp.idempresa) AS empresas_impuestos,
    -- Diagnóstico
    CASE 
        WHEN pr.idproducto IS NULL THEN 'Sin producto en crm_productos'
        WHEN pr.sat_clave IS NULL OR pr.sat_clave = '' OR pr.sat_clave = '0' THEN 'Sin clave SAT'
        WHEN pr.sat_medida IS NULL OR pr.sat_medida = '' THEN 'Sin unidad SAT'
        ELSE 'Datos completos'
    END AS diagnostico_config,
    -- Origen de impuestos
    CASE 
        WHEN COUNT(imp.idiva) > 0 THEN 'crm_impuestos'
        ELSE 'ticket_original'
    END AS origen_impuestos,
    -- Impuestos finales que se usarán
    CASE 
        WHEN COUNT(imp.idiva) > 0 THEN COALESCE(MAX(imp.iva), 0)
        ELSE d.iva
    END AS iva_final,
    CASE 
        WHEN COUNT(imp.idiva) > 0 THEN COALESCE(MAX(imp.ieps1), 0)
        ELSE 0
    END AS ieps1_final,
    CASE 
        WHEN COUNT(imp.idiva) > 0 THEN COALESCE(MAX(imp.ieps2), 0)
        ELSE 0
    END AS ieps2_final,
    CASE 
        WHEN COUNT(imp.idiva) > 0 THEN COALESCE(MAX(imp.ieps3), 0)
        ELSE 0
    END AS ieps3_final,
    -- Cálculos con todos los impuestos
    (d.cantidad * d.precio) - d.descuento AS subtotal,
    -- IVA
    ((d.cantidad * d.precio) - d.descuento) * 
    (CASE 
        WHEN COUNT(imp.idiva) > 0 THEN COALESCE(MAX(imp.iva), 0)
        ELSE d.iva
    END / 100.0) AS iva_importe,
    -- IEPS1
    ((d.cantidad * d.precio) - d.descuento) * 
    (CASE 
        WHEN COUNT(imp.idiva) > 0 THEN COALESCE(MAX(imp.ieps1), 0)
        ELSE 0
    END / 100.0) AS ieps1_importe,
    -- IEPS2
    ((d.cantidad * d.precio) - d.descuento) * 
    (CASE 
        WHEN COUNT(imp.idiva) > 0 THEN COALESCE(MAX(imp.ieps2), 0)
        ELSE 0
    END / 100.0) AS ieps2_importe,
    -- IEPS3
    ((d.cantidad * d.precio) - d.descuento) * 
    (CASE 
        WHEN COUNT(imp.idiva) > 0 THEN COALESCE(MAX(imp.ieps3), 0)
        ELSE 0
    END / 100.0) AS ieps3_importe,
    -- Total con todos los impuestos
    ((d.cantidad * d.precio) - d.descuento) * 
    (1 + 
        (CASE WHEN COUNT(imp.idiva) > 0 THEN COALESCE(MAX(imp.iva), 0) ELSE d.iva END / 100.0) +
        (CASE WHEN COUNT(imp.idiva) > 0 THEN COALESCE(MAX(imp.ieps1), 0) ELSE 0 END / 100.0) +
        (CASE WHEN COUNT(imp.idiva) > 0 THEN COALESCE(MAX(imp.ieps2), 0) ELSE 0 END / 100.0) +
        (CASE WHEN COUNT(imp.idiva) > 0 THEN COALESCE(MAX(imp.ieps3), 0) ELSE 0 END / 100.0)
    ) AS total_con_todos_impuestos
FROM optimus.crm_pedidos p 
JOIN optimus.crm_pedidos_det d ON p.id_pedido = d.id_pedido 
LEFT JOIN optimus.crm_productos pr ON d.idproducto = pr.idproducto 
LEFT JOIN optimus.crm_impuestos imp ON pr.idempresa = imp.idempresa
WHERE p.clave_pedido = '19a5a142aca218a7ff5deec2831d6834'
GROUP BY p.id_pedido, p.clave_pedido, d.descripcion, d.cantidad, d.precio, 
         d.precio_o, d.iva, d.descuento, d.idproducto, pr.clave, 
         pr.sat_clave, pr.sat_medida, pr.idempresa
ORDER BY d.idproducto;

-- ============================================================================
-- CONSULTA RESUMEN CON IEPS - MÁS FÁCIL DE LEER
-- ============================================================================
SELECT 
    d.idproducto AS 'Código',
    d.descripcion AS 'Producto',
    d.cantidad AS 'Cant',
    d.precio AS 'Precio',
    ROUND((d.cantidad * d.precio) - d.descuento, 2) AS 'Subtotal',
    -- Impuestos originales del ticket
    d.iva AS 'IVA Ticket',
    -- Impuestos desde configuración
    COALESCE(MAX(imp.iva), 0) AS 'IVA Config',
    COALESCE(MAX(imp.ieps1), 0) AS 'IEPS1 Config',
    COALESCE(MAX(imp.ieps2), 0) AS 'IEPS2 Config',
    COALESCE(MAX(imp.ieps3), 0) AS 'IEPS3 Config',
    -- Impuestos finales (los que se usarán)
    CASE 
        WHEN COUNT(imp.idiva) > 0 THEN COALESCE(MAX(imp.iva), 0)
        ELSE d.iva
    END AS 'IVA Final',
    CASE 
        WHEN COUNT(imp.idiva) > 0 THEN COALESCE(MAX(imp.ieps1), 0)
        ELSE 0
    END AS 'IEPS1 Final',
    CASE 
        WHEN COUNT(imp.idiva) > 0 THEN COALESCE(MAX(imp.ieps2), 0)
        ELSE 0
    END AS 'IEPS2 Final',
    CASE 
        WHEN COUNT(imp.idiva) > 0 THEN COALESCE(MAX(imp.ieps3), 0)
        ELSE 0
    END AS 'IEPS3 Final',
    -- Total con todos los impuestos
    ROUND(((d.cantidad * d.precio) - d.descuento) * 
    (1 + 
        (CASE WHEN COUNT(imp.idiva) > 0 THEN COALESCE(MAX(imp.iva), 0) ELSE d.iva END / 100.0) +
        (CASE WHEN COUNT(imp.idiva) > 0 THEN COALESCE(MAX(imp.ieps1), 0) ELSE 0 END / 100.0) +
        (CASE WHEN COUNT(imp.idiva) > 0 THEN COALESCE(MAX(imp.ieps2), 0) ELSE 0 END / 100.0) +
        (CASE WHEN COUNT(imp.idiva) > 0 THEN COALESCE(MAX(imp.ieps3), 0) ELSE 0 END / 100.0)
    ), 2) AS 'Total Final',
    -- Estado
    COUNT(imp.idiva) AS 'Configs',
    CASE 
        WHEN COUNT(imp.idiva) > 0 THEN 'Configurado'
        ELSE 'Sin configurar'
    END AS 'Estado',
    pr.sat_clave AS 'Clave SAT',
    pr.sat_medida AS 'Unidad SAT'
FROM optimus.crm_pedidos p 
JOIN optimus.crm_pedidos_det d ON p.id_pedido = d.id_pedido 
LEFT JOIN optimus.crm_productos pr ON d.idproducto = pr.idproducto 
LEFT JOIN optimus.crm_impuestos imp ON pr.idempresa = imp.idempresa
WHERE p.clave_pedido = '19a5a142aca218a7ff5deec2831d6834'
GROUP BY d.idproducto, d.descripcion, d.cantidad, d.precio, d.iva, pr.sat_clave, pr.sat_medida
ORDER BY d.idproducto;

-- ============================================================================
-- INSTRUCCIONES:
-- 1. Reemplaza '19a5a142aca218a7ff5deec2831d6834' por tu serie real
-- 2. Ejecuta cada consulta por separado en HeidiSQL
-- 3. La primera consulta te mostrará todos los datos completos
-- 4. La segunda consulta te mostrará un resumen más fácil de leer
-- 5. La tercera verificará si el pedido existe
-- 6. La cuarta te mostrará las configuraciones de impuestos disponibles
-- 7. La quinta consulta muestra todos los datos incluyendo IEPS
-- 8. La sexta consulta es un resumen incluyendo IEPS
-- ============================================================================
