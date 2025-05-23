import React, { useState, useEffect } from 'react';
import '../styles/HistorialFacturas.css';

function HistorialFacturas() {
  const [facturas, setFacturas] = useState([]);
  const [cargando, setCargando] = useState(true);
  const [error, setError] = useState(null);
  const [usuarioNoAutenticado, setUsuarioNoAutenticado] = useState(false);
  const [success, setSuccess] = useState('');

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

  // Obtiene ID de usuario del localStorage
  const getUserId = () => {
    const userDataString = localStorage.getItem('userData');
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

  return (
    <div className="empresas-container" style={{ marginTop: '60px', marginLeft: '230px' }}>
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
          <h1 className="titulo">Facturas Generadas</h1>
          
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