import React, { useState, useEffect, useRef, useCallback } from 'react';
import '../STYLES/HistorialFacturas.css';
import '../STYLES/Notificaciones.css'; 

function HistorialFacturas() {
  const [facturas, setFacturas] = useState([]);
  const [cargando, setCargando] = useState(true);
  const [error, setError] = useState(null);
  const [usuarioNoAutenticado, setUsuarioNoAutenticado] = useState(false);
  const [success, setSuccess] = useState('');
  
  // Estados para paginación
  const [paginaActual, setPaginaActual] = useState(1);
  const [totalPaginas, setTotalPaginas] = useState(1);
  const [totalFacturas, setTotalFacturas] = useState(0);
  const [tieneSiguiente, setTieneSiguiente] = useState(false);
  const [tieneAnterior, setTieneAnterior] = useState(false);
  const [limite] = useState(10); // Valor fijo sin selector - 10 facturas por página
  
  // Estados para búsqueda
  const [criterioBusqueda, setCriterioBusqueda] = useState('folio'); // 'folio', 'rfc_receptor', 'razon_social_receptor'
  const [valorBusqueda, setValorBusqueda] = useState('');
  const [buscando, setBuscando] = useState(false);
  
  // Estados para autocompletado de razón social
  const [sugerenciasRazonSocial, setSugerenciasRazonSocial] = useState([]);
  const [mostrarSugerencias, setMostrarSugerencias] = useState(false);
  const [indiceSugerenciaSeleccionada, setIndiceSugerenciaSeleccionada] = useState(-1);
  const [todasLasRazonesSociales, setTodasLasRazonesSociales] = useState([]);
  
  // Añadir estados para las notificaciones
  const [notificaciones, setNotificaciones] = useState([]);
  const notificacionIdRef = useRef(1);
  
  // Definir la función mostrarNotificacion
  const mostrarNotificacion = (mensaje, tipo = 'info') => {
    const id = notificacionIdRef.current++;
    const nuevaNotificacion = {
      id,
      mensaje,
      tipo, 
    };
    
    setNotificaciones(prev => [...prev, nuevaNotificacion]);
    
    // Auto-eliminar después de 5 segundos
    setTimeout(() => {
      setNotificaciones(prev => prev.filter(n => n.id !== id));
    }, 5000);
  };
  
  // Extraer cargarFacturas fuera del useEffect para que sea accesible
  const cargarFacturas = useCallback(async (pagina = 1) => {
    setCargando(true);
    setError(null);
    
    try {      
      // Obtener ID del usuario desde sessionStorage (migrado para evitar conflictos entre ventanas)
      const userData = JSON.parse(sessionStorage.getItem('userData'));
      if (!userData || !userData.id) {
        console.log("❌ NO HAY USUARIO AUTENTICADO");
        setUsuarioNoAutenticado(true);
        setCargando(false);
        return;
      }
      
      console.log("👤 USUARIO AUTENTICADO - ID:", userData.id);
      console.log("👤 DATOS COMPLETOS DEL USUARIO:", userData);
      
      const url = `http://localhost:8080/api/historial-facturas?id_usuario=${userData.id}&pagina=${pagina}&limite=${limite}`;
      console.log("🔗 URL DE CONSULTA:", url);
      
      console.log("Cargando facturas para usuario ID:", userData.id, "página:", pagina);
      
      const response = await fetch(url);
      
      if (!response.ok) {
        throw new Error(`Error al cargar facturas: ${response.status}`);
      }
      
      const data = await response.json();
      console.log("📊 DATOS RECIBIDOS DEL BACKEND:", data);
      console.log("📊 TIPO DE DATOS:", typeof data);
      
      // Verificar la estructura de los datos paginados
      if (data && data.facturas && Array.isArray(data.facturas)) {
        setFacturas(data.facturas);
        
        // Actualizar información de paginación
        if (data.paginacion) {
          setPaginaActual(data.paginacion.pagina_actual);
          setTotalPaginas(data.paginacion.total_paginas);
          setTotalFacturas(data.paginacion.total_facturas);
          setTieneSiguiente(data.paginacion.tiene_siguiente);
          setTieneAnterior(data.paginacion.tiene_anterior);
        }
      } else {
        console.warn("Formato de respuesta inesperado:", data);
        setFacturas([]);
        setPaginaActual(1);
        setTotalPaginas(1);
        setTotalFacturas(0);
        setTieneSiguiente(false);
        setTieneAnterior(false);
      }
    } catch (error) {
      console.error("Error al cargar facturas:", error);
      setError(`Error al cargar facturas: ${error.message}`);
    } finally {
      setCargando(false);
    }
  }, [limite]); // Solo depende de limite, no de la función

  // Función para realizar búsqueda
  const realizarBusqueda = async (pagina = 1) => {
    const userId = getUserId();
    if (!userId) {
      setUsuarioNoAutenticado(true);
      return;
    }

    setBuscando(true);
    setError(null);

    try {
      // Construir parámetros de búsqueda
      const params = new URLSearchParams({ 
        id_usuario: userId,
        pagina: pagina.toString(),
        limite: limite.toString()
      });
      
      // Solo agregar el criterio de búsqueda si hay un valor
      if (valorBusqueda.trim()) {
        params.append(criterioBusqueda, valorBusqueda.trim());
      }

      const url = `http://localhost:8080/api/buscar-facturas?${params.toString()}`;
      
      const response = await fetch(url);
      if (!response.ok) {
        throw new Error('Error al buscar facturas');
      }
      
      const data = await response.json();
      
      // Verificar la estructura de los datos paginados
      if (data && data.facturas && Array.isArray(data.facturas)) {
        setFacturas(data.facturas);
        
        // Actualizar información de paginación
        if (data.paginacion) {
          setPaginaActual(data.paginacion.pagina_actual);
          setTotalPaginas(data.paginacion.total_paginas);
          setTotalFacturas(data.paginacion.total_facturas);
          setTieneSiguiente(data.paginacion.tiene_siguiente);
          setTieneAnterior(data.paginacion.tiene_anterior);
        }
        
        if (data.facturas.length === 0 && pagina === 1) {
          mostrarNotificacion('No se encontraron facturas con los criterios especificados', 'info');
        } else if (pagina === 1) {
          mostrarNotificacion(`Se encontraron ${data.paginacion?.total_facturas || data.facturas.length} factura(s)`, 'success');
        }
      } else {
        setFacturas([]);
        setPaginaActual(1);
        setTotalPaginas(1);
        setTotalFacturas(0);
        setTieneSiguiente(false);
        setTieneAnterior(false);
        mostrarNotificacion('No se encontraron facturas con los criterios especificados', 'info');
      }
    } catch (err) {
      console.error('Error al buscar facturas:', err);
      setError('Error al buscar facturas. Por favor, intenta nuevamente.');
    } finally {
      setBuscando(false);
    }
  };

  // Función auxiliar para realizar búsqueda con un valor específico
  const realizarBusquedaConValor = async (valor, pagina = 1) => {
    const userId = getUserId();
    if (!userId) {
      setUsuarioNoAutenticado(true);
      return;
    }

    setBuscando(true);
    setError(null);

    try {
      // Construir parámetros de búsqueda con el valor específico
      const params = new URLSearchParams({ 
        id_usuario: userId,
        pagina: pagina.toString(),
        limite: limite.toString()
      });
      
      if (valor && valor.trim()) {
        params.append(criterioBusqueda, valor.trim());
      }

      const url = `http://localhost:8080/api/buscar-facturas?${params.toString()}`;
      
      const response = await fetch(url);
      if (!response.ok) {
        throw new Error('Error al buscar facturas');
      }
      
      const data = await response.json();
      
      // Manejar respuesta paginada
      if (data && data.facturas && Array.isArray(data.facturas)) {
        setFacturas(data.facturas);
        
        // Actualizar información de paginación
        if (data.paginacion) {
          setPaginaActual(data.paginacion.pagina_actual);
          setTotalPaginas(data.paginacion.total_paginas);
          setTotalFacturas(data.paginacion.total_facturas);
          setTieneSiguiente(data.paginacion.tiene_siguiente);
          setTieneAnterior(data.paginacion.tiene_anterior);
        }
        
        if (data.facturas.length === 0) {
          mostrarNotificacion('No se encontraron facturas con los criterios especificados', 'info');
        } else {
          mostrarNotificacion(`Se encontraron ${data.paginacion ? data.paginacion.total_facturas : data.facturas.length} factura(s)`, 'success');
        }
      } else {
        setFacturas([]);
        setPaginaActual(1);
        setTotalPaginas(1);
        setTotalFacturas(0);
        setTieneSiguiente(false);
        setTieneAnterior(false);
        mostrarNotificacion('No se encontraron facturas con los criterios especificados', 'info');
      }
    } catch (err) {
      console.error('Error al buscar facturas:', err);
      setError('Error al buscar facturas. Por favor, intenta nuevamente.');
    } finally {
      setBuscando(false);
    }
  };

  // Función para cargar todas las razones sociales únicas
  const cargarTodasLasRazonesSociales = useCallback(async () => {
    try {
      const userId = getUserId();
      if (!userId) return;

      // Usar el endpoint correcto y obtener todas las facturas para extraer las razones sociales
      const response = await fetch(`http://localhost:8080/api/historial-facturas?id_usuario=${userId}&limite=1000`);
      if (!response.ok) return;

      const data = await response.json();
      let facturas = [];
      
      // Extraer facturas según la estructura de respuesta
      if (data && data.facturas && Array.isArray(data.facturas)) {
        facturas = data.facturas;
      } else if (Array.isArray(data)) {
        facturas = data;
      }
      
      // Extraer todas las razones sociales únicas
      const razonesSociales = [...new Set(
        facturas
          .map(factura => factura.razon_social_receptor)
          .filter(razon => razon && razon.trim() !== '')
      )].sort(); // Ordenar alfabéticamente

      setTodasLasRazonesSociales(razonesSociales);
    } catch (error) {
      console.error('Error al cargar razones sociales:', error);
    }
  }, []); // No tiene dependencias externas

  // Función para filtrar sugerencias localmente (como Spotify)
  const filtrarSugerenciasLocal = (termino) => {
    if (!termino || termino.length < 3) {
      setSugerenciasRazonSocial([]);
      setMostrarSugerencias(false);
      return;
    }

    const terminoLower = termino.toLowerCase();
    const sugerenciasFiltradas = todasLasRazonesSociales
      .filter(razon => razon.toLowerCase().includes(terminoLower))
      .slice(0, 8); // Limitar a 8 sugerencias

    setSugerenciasRazonSocial(sugerenciasFiltradas);
    setMostrarSugerencias(sugerenciasFiltradas.length > 0);
  };

  // Función para manejar el cambio en el input de búsqueda
  const manejarCambioBusqueda = (valor) => {
    setValorBusqueda(valor);
    setIndiceSugerenciaSeleccionada(-1);
    
    // Solo buscar sugerencias si el criterio es razón social
    if (criterioBusqueda === 'razon_social_receptor') {
      filtrarSugerenciasLocal(valor); // Usar filtrado local instantáneo
    }
  };

  // Función para seleccionar una sugerencia
  const seleccionarSugerencia = (sugerencia) => {
    setValorBusqueda(sugerencia);
    setMostrarSugerencias(false);
    setSugerenciasRazonSocial([]);
    setIndiceSugerenciaSeleccionada(-1);
    
    // Automáticamente ejecutar la búsqueda con la razón social seleccionada
    realizarBusquedaConValor(sugerencia);
  };

  // Función para manejar teclas en el input
  const manejarTeclas = (e) => {
    if (!mostrarSugerencias || sugerenciasRazonSocial.length === 0) {
      if (e.key === 'Enter') {
        realizarBusqueda();
      }
      return;
    }

    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault();
        setIndiceSugerenciaSeleccionada(prev => 
          prev < sugerenciasRazonSocial.length - 1 ? prev + 1 : prev
        );
        break;
      case 'ArrowUp':
        e.preventDefault();
        setIndiceSugerenciaSeleccionada(prev => prev > 0 ? prev - 1 : -1);
        break;
      case 'Enter':
        e.preventDefault();
        if (indiceSugerenciaSeleccionada >= 0) {
          seleccionarSugerencia(sugerenciasRazonSocial[indiceSugerenciaSeleccionada]);
        } else {
          realizarBusqueda();
        }
        break;
      case 'Escape':
        setMostrarSugerencias(false);
        setIndiceSugerenciaSeleccionada(-1);
        break;
      default:
        break;
    }
  };

  // Función para limpiar búsqueda y cargar todas las facturas
  const limpiarBusqueda = () => {
    setValorBusqueda('');
    setCriterioBusqueda('folio');
    setMostrarSugerencias(false);
    setSugerenciasRazonSocial([]);
    setIndiceSugerenciaSeleccionada(-1);
    cargarFacturas(1); // Reiniciar desde la primera página
  };

  // Efecto para limpiar sugerencias cuando cambia el criterio
  useEffect(() => {
    setMostrarSugerencias(false);
    setSugerenciasRazonSocial([]);
    setIndiceSugerenciaSeleccionada(-1);
    // Limpiar el campo de búsqueda al cambiar cualquier criterio
    setValorBusqueda('');
  }, [criterioBusqueda]); // Solo cuando cambia el criterio de búsqueda

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

  // Cargar facturas del usuario al montar el componente
  useEffect(() => {
    cargarFacturas(1);
  }, []); // Solo al montar

  // Cargar razones sociales una sola vez al montar el componente
  useEffect(() => {
    cargarTodasLasRazonesSociales();
  }, []); // Solo al montar

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

  // Funciones de paginación
  const cambiarPagina = (nuevaPagina) => {
    if (nuevaPagina !== paginaActual && nuevaPagina >= 1 && nuevaPagina <= totalPaginas) {
      if (valorBusqueda.trim()) {
        // Si hay búsqueda activa, usar búsqueda paginada
        realizarBusqueda(nuevaPagina);
      } else {
        // Si no hay búsqueda, cargar facturas normales
        cargarFacturas(nuevaPagina);
      }
    }
  };

  const generarBotonesPagina = () => {
    const botones = [];
    const maxBotones = 5;
    
    if (totalPaginas <= maxBotones) {
      // Si hay pocas páginas, mostrar todas
      for (let i = 1; i <= totalPaginas; i++) {
        botones.push(i);
      }
    } else {
      // Lógica más compleja para muchas páginas
      if (paginaActual <= 3) {
        // Estamos al principio
        for (let i = 1; i <= 4; i++) {
          botones.push(i);
        }
        botones.push('...');
        botones.push(totalPaginas);
      } else if (paginaActual >= totalPaginas - 2) {
        // Estamos al final
        botones.push(1);
        botones.push('...');
        for (let i = totalPaginas - 3; i <= totalPaginas; i++) {
          botones.push(i);
        }
      } else {
        // Estamos en el medio
        botones.push(1);
        botones.push('...');
        for (let i = paginaActual - 1; i <= paginaActual + 1; i++) {
          botones.push(i);
        }
        botones.push('...');
        botones.push(totalPaginas);
      }
    }
    
    return botones;
  };

  return (
    <div className="historial-facturas-container" style={{ marginTop: '60px', marginLeft: '290px' }}>
      {/* Componente de notificaciones */}
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

      <h1 className="titulo">Historial de Facturas</h1>

      <div className="info-card">
        <div className="card-header">
          <h2>Facturas Generadas</h2>
          <button
            className="btn-refresh"
            onClick={() => {
              mostrarNotificacion('Actualizando lista de facturas...', 'info');
              // Si hay un valor de búsqueda, actualizar con búsqueda; si no, cargar todas
              if (valorBusqueda.trim()) {
                realizarBusqueda();
              } else {
                cargarFacturas(1);
              }
            }}
            title="Actualizar lista de facturas"
          >
            🔄 Actualizar
          </button>
        </div>

        {/* Barra de búsqueda */}
        <div className="busqueda-container">
          <div className="busqueda-avanzada">
            <div className="selector-criterio">
              <label htmlFor="criterio">Buscar por:</label>
              <select
                id="criterio"
                value={criterioBusqueda}
                onChange={(e) => setCriterioBusqueda(e.target.value)}
                className="select-criterio"
              >
                <option value="folio">Folio</option>
                <option value="rfc_receptor">RFC</option>
                <option value="razon_social_receptor">Razón Social</option>
              </select>
            </div>
            <div className="campo-busqueda-unico">
              <div className="input-container-con-sugerencias">
                <input
                  type="text"
                  placeholder={
                    criterioBusqueda === 'folio' ? 'Buscar por folio (ej: F000001)' :
                    criterioBusqueda === 'rfc_receptor' ? 'Buscar por RFC' :
                    'Buscar por razón social'
                  }
                  value={valorBusqueda}
                  onChange={(e) => manejarCambioBusqueda(e.target.value)}
                  onKeyDown={manejarTeclas}
                  onFocus={() => {
                    if (criterioBusqueda === 'razon_social_receptor' && valorBusqueda.length >= 3) {
                      filtrarSugerenciasLocal(valorBusqueda);
                    }
                  }}
                  onBlur={() => {
                    // Delay para permitir clicks en sugerencias
                    setTimeout(() => setMostrarSugerencias(false), 150);
                  }}
                  className="input-busqueda"
                />
                {mostrarSugerencias && sugerenciasRazonSocial.length > 0 && (
                  <div className="sugerencias-container">
                    {sugerenciasRazonSocial.map((sugerencia, index) => (
                      <div
                        key={index}
                        className={`sugerencia-item ${index === indiceSugerenciaSeleccionada ? 'seleccionada' : ''}`}
                        onClick={() => seleccionarSugerencia(sugerencia)}
                        onMouseEnter={() => setIndiceSugerenciaSeleccionada(index)}
                      >
                        {sugerencia}
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </div>
            <div className="busqueda-botones">
              <button
                className="btn-buscar"
                onClick={realizarBusqueda}
                disabled={buscando}
              >
                {buscando ? '🔍 Buscando...' : '🔍 Buscar'}
              </button>
              <button
                className="btn-limpiar"
                onClick={limpiarBusqueda}
                disabled={buscando}
              >
                🗑️ Limpiar
              </button>
            </div>
          </div>
        </div>

        <div className="card-content">
          {cargando && (
            <div className="loading-spinner">
              <div className="spinner"></div>
              <p>Cargando facturas...</p>
            </div>
          )}
          
          {!cargando && usuarioNoAutenticado && (
            <div className="error-message">
              No hay un usuario activo. Por favor inicie sesión para ver sus facturas.
            </div>
          )}
          
          {!cargando && !usuarioNoAutenticado && facturas.length === 0 && (
            <div className="no-facturas">
              <p>No tienes facturas generadas en tu historial.</p>
            </div>
          )}
          
          {!cargando && !usuarioNoAutenticado && facturas.length > 0 && (
            <div className="tabla-facturas-container">
              <table className="tabla-facturas">
                <thead>
                  <tr>
                    <th>Fecha</th>
                    <th>Folio</th>
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
                      <td className="folio-cell">
                        {factura.folio && factura.folio.trim() !== '' 
                          ? factura.folio 
                          : <span className="folio-faltante">Sin folio</span>
                        }
                      </td>
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
                          className="btn-descargar"
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
              
              {/* Controles de paginación */}
              {totalPaginas > 1 && (
                <div className="pagination-controls">
                  <div className="pagination-info">
                    <span>
                      Página {paginaActual} de {totalPaginas} 
                      ({totalFacturas} facturas en total)
                    </span>
                  </div>
                  <div className="pagination-buttons">
                    <button 
                      className="btn-pagination"
                      disabled={!tieneAnterior}
                      onClick={() => cambiarPagina(paginaActual - 1)}
                    >
                      ← Anterior
                    </button>
                    
                    {/* Números de página */}
                    {generarBotonesPagina().map((pagina, index) => (
                      <button
                        key={index}
                        className={`btn-pagination ${pagina === paginaActual ? 'active' : ''} ${pagina === '...' ? 'dots' : ''}`}
                        disabled={pagina === '...'}
                        onClick={() => pagina !== '...' && cambiarPagina(pagina)}
                      >
                        {pagina}
                      </button>
                    ))}
                    
                    <button 
                      className="btn-pagination"
                      disabled={!tieneSiguiente}
                      onClick={() => cambiarPagina(paginaActual + 1)}
                    >
                      Siguiente →
                    </button>
                  </div>
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default HistorialFacturas;