import React, { useState, useEffect, useRef } from 'react';
import '../styles/HistorialFacturas.css';
import '../styles/Notificaciones.css'; 

function HistorialFacturas() {
  const [facturas, setFacturas] = useState([]);
  const [cargando, setCargando] = useState(true);
  const [error, setError] = useState(null);
  const [usuarioNoAutenticado, setUsuarioNoAutenticado] = useState(false);
  const [success, setSuccess] = useState('');
  
  // Añadir estados para las notificaciones
  const [notificaciones, setNotificaciones] = useState([]);
  const notificacionIdRef = useRef(1);
  
  // Definir la función mostrarNotificacion
  const mostrarNotificacion = (mensaje, tipo = 'info') => {
    const id = notificacionIdRef.current++;
    const nuevaNotificacion = {
      id,
      mensaje,
      tipo, // 'success', 'error', 'info'
    };
    
    setNotificaciones(prev => [...prev, nuevaNotificacion]);
    
    // Auto-eliminar después de 5 segundos
    setTimeout(() => {
      setNotificaciones(prev => prev.filter(n => n.id !== id));
    }, 5000);
  };
  
  // Extraer cargarFacturas fuera del useEffect para que sea accesible
  const cargarFacturas = async () => {
    setCargando(true);
    setError(null);
    
    try {      // Obtener ID del usuario desde sessionStorage (migrado para evitar conflictos entre ventanas)
      const userData = JSON.parse(sessionStorage.getItem('userData'));
      if (!userData || !userData.id) {
        setUsuarioNoAutenticado(true);
        setCargando(false);
        return;
      }
      
      console.log("Cargando facturas para usuario ID:", userData.id);
      
      const response = await fetch(`http://localhost:8080/api/historial_facturas?id_usuario=${userData.id}`);
      
      if (!response.ok) {
        throw new Error(`Error al cargar facturas: ${response.status}`);
      }
      
      const data = await response.json();
      console.log("Datos de facturas recibidos:", data);
      
      // Verificar la estructura de los datos
      if (Array.isArray(data)) {
        setFacturas(data);
      } else if (data && Array.isArray(data.facturas)) {
        setFacturas(data.facturas);
      } else {
        console.warn("Formato de respuesta inesperado:", data);
        setFacturas([]);
      }
    } catch (error) {
      console.error("Error al cargar facturas:", error);
      setError(`Error al cargar facturas: ${error.message}`);
    } finally {
      setCargando(false);
    }
  };

  // Función para obtener el estado formateado para mostrar
  const getEstadoDisplay = (estadoCodigo) => {
    if (!estadoCodigo) return 'Generada'; // Handle NULL/undefined
    
    switch(estadoCodigo) {
      case 'G':
        return 'Generada';
      case 'P':
        return 'Pendiente';
      case 'C':
        return 'Cancelada';
      case 'E':
        return 'Error';
      default:
        return estadoCodigo || 'Generada';
    }
  };

  // Función para obtener la clase CSS basada en el estado
  const getEstadoClassName = (estado) => {
    if (!estado) return 'estado-g'; // Default for NULL values
    return `estado-${estado.toLowerCase()}`;
  };
  // Obtiene ID de usuario del sessionStorage (migrado para evitar conflictos)
  const getUserId = () => {
    const userDataString = sessionStorage.getItem('userData');
    if (!userDataString) return null;
    
    try {
      const userData = JSON.parse(userDataString);
      return userData.id;
    } catch (e) {
      console.error("Error al parsear userData:", e);
      return null;
    }
  };

  // Cargar facturas del usuario
  useEffect(() => {
    const fetchFacturas = async () => {
      const userId = getUserId();
      if (!userId) {
        setCargando(false);
        setUsuarioNoAutenticado(true);
        return;
      }

      try {
        // Sin parámetros de búsqueda
        const url = `http://localhost:8080/api/historial-facturas?id_usuario=${userId}`;
        
        setCargando(true);
        const response = await fetch(url);
        if (!response.ok) {
          throw new Error('Error al obtener el historial de facturas');
        }
        
        const data = await response.json();
        setFacturas(data || []);
        console.log("Facturas cargadas:", data);
      } catch (err) {
        console.error('Error al cargar historial de facturas:', err);
        setError('Error al cargar el historial de facturas. Por favor, intenta nuevamente.');
      } finally {
        setCargando(false);
      }
    };

    fetchFacturas();
  }, []); // Cargar al montar el componente

  // Limpiar mensajes después de 5 segundos
  useEffect(() => {
    if (success || error) {
      const timer = setTimeout(() => {
        setSuccess('');
        setError('');
      }, 5000);
      return () => clearTimeout(timer);
    }
  }, [success, error]);

  // Modificar la función formatearFecha para usar la zona horaria del usuario
  const formatearFecha = (fechaStr) => {
    if (!fechaStr) return 'Fecha no disponible';
    
    try {
      let fecha;
      
      // Convertir la fecha a objeto Date
      if (typeof fechaStr === 'string' && fechaStr.match(/^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$/)) {
        const [datePart, timePart] = fechaStr.split(' ');
        const [year, month, day] = datePart.split('-');
        const [hour, minute, second] = timePart.split(':');
        
        // Crear fecha en UTC
        fecha = new Date(Date.UTC(year, month - 1, day, hour, minute, second));
      } else {
        fecha = new Date(fechaStr);
      }
      
      if (!fecha || isNaN(fecha.getTime())) {
        return String(fechaStr);
      }
      
      // Usar toLocaleDateString para formatear según la zona horaria del navegador
      const options = {
        day: '2-digit',
        month: 'short',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
        hour12: false
      };
      
      // Esto ajustará automáticamente a la zona horaria local del usuario
      const fechaFormateada = fecha.toLocaleDateString('es-MX', options);
      
      // Personalizamos un poco más para obtener exactamente el formato deseado
      // Convertir "may." a "May" (primeras 3 letras en mayúscula inicial)
      return fechaFormateada.replace(/(\d{2})\s*de\s*([a-z]{3})\.?\s*de\s*(\d{4})/, (_, dia, mes, anio) => {
        return `${dia}/${mes.charAt(0).toUpperCase() + mes.slice(1, 3)}/${anio}`;
      });
    } catch (e) {
      console.error("Error al formatear fecha:", e, "Valor recibido:", fechaStr);
      return String(fechaStr);
    }
  };

  // Función para descargar una factura del historial
  async function descargarFactura(facturaId) {
    try {
      setSuccess('Descargando factura...');
      const response = await fetch(`http://localhost:8080/api/descargar-factura/${facturaId}`);
      if (!response.ok) {
        throw new Error('Error al descargar la factura');
      }
      
      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `factura-${facturaId}.zip`;
      link.click();
      
      setTimeout(() => window.URL.revokeObjectURL(url), 100);
      setSuccess('Factura descargada correctamente');
    } catch (error) {
      setError('Error al descargar la factura: ' + error.message);
    }
  }

  // Modificar el useEffect para usar la función cargarFacturas
  useEffect(() => {
    cargarFacturas();
    
    // Actualizar facturas cada 30 segundos
    const intervalo = setInterval(cargarFacturas, 30000);
    return () => clearInterval(intervalo);
  }, []);

  return (
    <div className="empresas-container" style={{ marginTop: '60px', marginLeft: '230px' }}>
      {/* Añadir componente de notificaciones */}
      <div className="notificaciones-container">
        {notificaciones.map(notif => (
          <div 
            key={notif.id} 
            className={`notificacion notificacion-${notif.tipo}`}
            onClick={() => setNotificaciones(prev => prev.filter(n => n.id !== notif.id))}
          >
            <div className="notificacion-contenido">
              <span className="notificacion-icono">
                {notif.tipo === 'success' ? '✓' : notif.tipo === 'error' ? '✕' : 'ℹ️'}
              </span>
              <span className="notificacion-mensaje">{notif.mensaje}</span>
            </div>
          </div>
        ))}
      </div>
      
      {error && <div className="error-message">{error}</div>}
      {success && <div className="success-message">{success}</div>}

      {cargando ? (
        <div className="loading-message">Cargando facturas...</div>
      ) : usuarioNoAutenticado ? (
        <div className="error-message">
          No hay un usuario activo. Por favor inicie sesión para ver sus facturas.
        </div>
      ) : (
        <>
          {/* Header con título centrado y botón de actualizar con color */}
          <div style={{ 
            marginBottom: '20px',
            position: 'relative',
            textAlign: 'center'
          }}>
            <h1 className="titulo" style={{ 
              textAlign: 'center', 
              marginBottom: '10px',
              position: 'relative',
              display: 'inline-block'
            }}>
              Facturas Generadas
            </h1>
            
            <div style={{
              height: '2px',
              backgroundColor: '#1890ff',
              width: '100%',
              marginBottom: '20px'
            }}></div>
            
            <button
              onClick={() => {
                mostrarNotificacion('Actualizando lista de facturas...', 'info');
                cargarFacturas();
              }}
              style={{
                position: 'absolute',
                right: '0',
                top: '0',
                padding: '8px 16px',
                backgroundColor: '#1890ff',
                color: 'white',
                border: 'none',
                borderRadius: '4px',
                cursor: 'pointer',
                display: 'flex',
                alignItems: 'center',
                gap: '5px',
                boxShadow: '0 2px 0 rgba(0,0,0,0.045)'
              }}
            >
              <span style={{ transform: 'rotate(90deg)' }}>↻</span>
              Actualizar
            </button>
          </div>
          
          {facturas.length === 0 ? (
            <div className="table-card">
              <p style={{ 
                textAlign: 'center', 
                padding: '30px 0',
                fontSize: '16px',
                color: '#666'
              }}>
                No tienes facturas generadas en tu historial.
              </p>
            </div>
          ) : (
            <div className="table-card">
              <table>
                <thead>
                  <tr>
                    <th>Fecha</th>
                    <th>RFC</th>
                    <th>Razón Social</th>
                    <th>Total</th>
                    <th>Estado</th>
                    <th>Acciones</th>
                  </tr>
                </thead>
                <tbody>
                  {facturas.map((factura) => (
                    <tr key={factura.id}>
                      <td>{formatearFecha(factura.fecha_generacion)}</td>
                      <td>{factura.rfc_receptor}</td>
                      <td>{factura.razon_social_receptor}</td>
                      <td>$ {factura.total.toLocaleString('es-MX', {minimumFractionDigits: 2, maximumFractionDigits: 2})}</td>
                      <td>
                        <span className={getEstadoClassName(factura.estado)}>
                          {getEstadoDisplay(factura.estado)}
                        </span>
                      </td>
                      <td className="acciones">
                        <button 
                          className="file-download-button action-button"
                          onClick={() => descargarFactura(factura.id)}
                          title="Descargar factura (ZIP con PDF y XML)"
                        >
                          📥 Descargar
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </>
      )}
    </div>
  );
}

export default HistorialFacturas;