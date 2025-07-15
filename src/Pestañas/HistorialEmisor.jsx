import React, { useEffect, useState, useRef, useCallback } from 'react';
import '../STYLES/HistorialFacturas.css';
import '../STYLES/Notificaciones.css';

function HistorialEmisor() {
  const [facturas, setFacturas] = useState([]);
  const [cargando, setCargando] = useState(true);
  const [error, setError] = useState(null);
  const [notificaciones, setNotificaciones] = useState([]);
  const notificacionIdRef = useRef(1);

  // Notificación rápida
  const mostrarNotificacion = (mensaje, tipo = 'info') => {
    const id = notificacionIdRef.current++;
    setNotificaciones(prev => [...prev, { id, mensaje, tipo }]);
    setTimeout(() => {
      setNotificaciones(prev => prev.filter(n => n.id !== id));
    }, 5000);
  };

  // Obtener id_usuario del usuario logueado
  const getUsuarioId = () => {
    const userDataString = sessionStorage.getItem('userData');
    if (!userDataString) return null;
    try {
      const userData = JSON.parse(userDataString);
      return userData.id || userData.id_usuario;
    } catch {
      return null;
    }
  };

  // Mueve cargarHistorial fuera de useEffect para poder reutilizarla en el botón de refrescar
  const cargarHistorial = useCallback(async () => {
    const idUsuario = getUsuarioId();
    if (!idUsuario) {
      setError("No se pudo identificar el usuario actual.");
      setCargando(false);
      return;
    }
    setCargando(true);
    setError(null);
    try {
      const response = await fetch(`http://localhost:8080/api/facturas-empresa-activa?id_usuario=${idUsuario}`);
      if (!response.ok) throw new Error('Error al cargar historial');
      const data = await response.json();
      if (Array.isArray(data)) {
        setFacturas(data);
        if (data.length === 0) {
          mostrarNotificacion('No hay facturas generadas a nombre de la empresa vinculada.', 'info');
        }
      } else {
        setFacturas([]);
        mostrarNotificacion('No hay facturas generadas a nombre de la empresa vinculada.', 'info');
      }
    } catch (err) {
      setError(err.message);
    } finally {
      setCargando(false);
    }
  }, []);

  useEffect(() => {
    cargarHistorial();
  }, [cargarHistorial]);

  // Formatear fecha
  const formatearFecha = (fechaStr) => {
    if (!fechaStr) return '';
    try {
      const fecha = new Date(fechaStr);
      if (isNaN(fecha.getTime())) return fechaStr;
      return fecha.toLocaleString('es-MX', {
        year: 'numeric',
        month: 'short',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        hour12: false
      });
    } catch {
      return fechaStr;
    }
  };

  return (
    <div className="historial-facturas-container" style={{ marginTop: '60px', marginLeft: '290px' }}>
      {/* Notificaciones */}
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

      <h1 className="titulo">Historial de Facturas</h1>

      <div className="info-card">
        <div className="card-header">
          <h2>Facturas Generadas a Nombre de la Empresa</h2>
          <button
            className="btn-refresh"
            onClick={cargarHistorial}
            title="Actualizar lista de facturas"
          >
            🔄 Actualizar
          </button>
        </div>

        <div className="card-content">
          {cargando && (
            <div className="loading-spinner">
              <div className="spinner"></div>
              <p>Cargando facturas...</p>
            </div>
          )}

          {!cargando && facturas.length === 0 && !error && (
            <div className="no-facturas">
              <p>No hay facturas generadas a nombre de esta empresa.</p>
            </div>
          )}

          {!cargando && facturas.length > 0 && (
            <div className="tabla-facturas-container">
              <table className="tabla-facturas">
                <thead>
                  <tr>
                    <th>Fecha</th>
                    <th>Folio</th>
                    <th>RFC Cliente</th>
                    <th>Razón Social</th>
                    <th>Total</th>
                    <th>Usuario que generó</th>
                    <th>Email usuario</th>
                  </tr>
                </thead>
                <tbody>
                  {facturas.map((factura) => (
                    <tr key={factura.id}>
                      <td>{formatearFecha(factura.fecha_emision || factura.fecha)}</td>
                      <td>{factura.folio}</td>
                      <td>{factura.rfc_receptor || factura.rfc_cliente}</td>
                      <td>{factura.razon_social_receptor || factura.cliente}</td>
                      <td>$ {parseFloat(factura.total).toLocaleString('es-MX', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</td>
                      <td>{factura.usuario?.nombre ?? '-'}</td>
                      <td>{factura.usuario?.email ?? '-'}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
              <div className="tabla-footer">
                <span>Total de facturas: {facturas.length}</span>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default HistorialEmisor;