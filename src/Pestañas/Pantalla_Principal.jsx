import React, { useState, useEffect, useRef } from 'react';
import { Link } from 'react-router-dom';
import '../styles/Pantalla_Principal.css';
import '../styles/Notificaciones.css'; // Si ya existe este archivo de estilos

const formatCurrency = (amount) => {
  // Si no hay valor o es 0, mostrar $0.00
  if (!amount) return '$ 0.00';
  
  // Formatear con separadores de miles y 2 decimales
  return `$ ${new Intl.NumberFormat('es-MX', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  }).format(amount)}`;
};

function InicioFacturacion() {
  const [empresas, setEmpresas] = useState([]);
  const [cargandoEmpresas, setCargandoEmpresas] = useState(true);
  const [empresaSeleccionadaId, setEmpresaSeleccionadaId] = useState('');
  const [empresa, setEmpresa] = useState(null);
  const [error, setError] = useState(null);
  const [generando, setGenerando] = useState(false);
  const [buscandoVentas, setBuscandoVentas] = useState(false);
  const [usoCfdi, setUsoCfdi] = useState('');
  const [ventas, setVentas] = useState([]);
  const [usuarioNoAutenticado, setUsuarioNoAutenticado] = useState(false);
  
  // Estado para el modo de edición
  const [isEditing, setIsEditing] = useState(false);
  const [originalFormData, setOriginalFormData] = useState(null);

  const [ticketData, setTicketData] = useState({
    claveTicket: '',
    totalTicket: '',
  });

  const [formData, setFormData] = useState({
    razonSocial: '',
    direccion: '',
    codigoPostal: '',
    pais: '',
    estado: '',
    localidad: '',
    municipio: '',
    colonia: '',
    observaciones: '',
    rfc: '',
  });

  const [showDatosPrincipales, setShowDatosPrincipales] = useState(false);
  const [showDatosAdicionales, setShowDatosAdicionales] = useState(false);

  const camposObligatoriosForm = [
    'razonSocial', 'direccion', 'codigoPostal', 'pais', 'estado',
    'localidad', 'municipio', 'colonia'
  ];

  const camposObligatoriosTicket = ['claveTicket', 'totalTicket'];

  // Obtiene ID de usuario del sessionStorage (evita conflictos entre ventanas)
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

  // Cargar empresas del usuario
  useEffect(() => {
    const fetchEmpresas = async () => {
      const userId = getUserId();
      if (!userId) {
        setCargandoEmpresas(false);
        setUsuarioNoAutenticado(true);
        return;
      }

      try {
        const response = await fetch(`http://localhost:8080/api/empresas?id_usuario=${userId}`);
        if (!response.ok) {
          throw new Error('Error al obtener las empresas');
        }
        const data = await response.json();
        setEmpresas(data || []);
      } catch (err) {
        console.error('Error al cargar empresas:', err);
        setError('Error al cargar las empresas. Por favor, intenta nuevamente.');
      } finally {
        setCargandoEmpresas(false);
      }
    };

    fetchEmpresas();
  }, []);

  const seleccionarEmpresa = (e) => {
    const empresaId = e.target.value;
    setEmpresaSeleccionadaId(empresaId);
    setIsEditing(false); // Reiniciar modo de edición al cambiar de empresa
    
    if (!empresaId) {
      setEmpresa(null);
      setShowDatosPrincipales(false);
      setShowDatosAdicionales(false);
      return;
    }
    
    const empresaSeleccionada = empresas.find(emp => emp.id.toString() === empresaId);
    if (empresaSeleccionada) {
      const newFormData = {
        rfc: empresaSeleccionada.rfc || '',
        razonSocial: empresaSeleccionada.razon_social || '',
        direccion: empresaSeleccionada.direccion || '',
        codigoPostal: empresaSeleccionada.codigo_postal || '',
        pais: empresaSeleccionada.pais || '',
        estado: empresaSeleccionada.estado || '',
        localidad: empresaSeleccionada.localidad || '',
        municipio: empresaSeleccionada.municipio || '',
        colonia: empresaSeleccionada.colonia || '',
        observaciones: '',
      };
      
      setFormData(newFormData);
      setOriginalFormData({...newFormData}); // Guardar datos originales para cancelar
      setEmpresa(empresaSeleccionada);
      setShowDatosPrincipales(true);
      setShowDatosAdicionales(true);
    }
  };

  // Agregar función para manejar el modo de edición
  const handleEditMode = () => {
    setIsEditing(true);
  };

  // Agregar función para cancelar el modo de edición
  const handleCancelEdit = () => {
    setIsEditing(false);
    setFormData({...originalFormData}); // Restaurar datos originales
  };

  // Agregar función para guardar cambios
  const handleSaveChanges = async () => {
    // Aquí agregarías validación y llamada a API para guardar cambios
    // Por ahora, solo actualizamos la UI
    try {
      // Simular llamada a API
      // const response = await fetch('...', { method: 'POST', body: ... });
      
      // Actualizar datos originales
      setOriginalFormData({...formData});
      setIsEditing(false);
      mostrarNotificacion('Datos actualizados correctamente', 'success');
    } catch (error) {
      mostrarNotificacion('Error al actualizar los datos: ' + error.message, 'error');
    }
  };

  // Componente de pantalla de carga con mensaje personalizable
  const PantallaDeCarga = ({ mensaje = "Generando factura, por favor espera..." }) => (
    <div className="pantalla-carga">
      <div className="spinner"></div>
      <p>{mensaje}</p>
    </div>
  );

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData((prevFormData) => ({
      ...prevFormData,
      [name]: value,
    }));
  };

  // Buscar ventas en la base de datos optimus
  const buscarVentas = async (claveTicket) => {
    if (!claveTicket || claveTicket.length < 30 || !/^[a-zA-Z0-9]+$/.test(claveTicket)) {
      mostrarNotificacion('La clave del ticket debe tener al menos 32 caracteres alfanuméricos.', 'error');
      return;
    }
    
    setBuscandoVentas(true); // Mostrar pantalla de carga
  
    try {
      const response = await fetch(`http://localhost:8080/api/ventas?serie=${encodeURIComponent(claveTicket)}`);
      
      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || 'Error al buscar las ventas.');
      }
    
      const data = await response.json();
      console.log("Datos de ventas recibidos:", data);
      
      // Mapear los datos para incluir información de impuestos más detallada
      const ventasConImpuestos = data.ventas || [];
      setVentas(ventasConImpuestos);
      
      // Si no hay ventas, mostrar mensaje
      if (!data.ventas || data.ventas.length === 0) {
        mostrarNotificacion('No se encontraron ventas para la clave de ticket proporcionada', 'error');
      } else {
        // Mostrar información detallada de impuestos si está disponible
        console.log("Datos de impuestos en las ventas:", {
          primeraVenta: data.ventas[0],
          propiedadesDisponibles: Object.keys(data.ventas[0]),
          datosImpuestos: {
            iva: data.ventas[0].iva,
            ieps1: data.ventas[0].ieps1,
            ieps2: data.ventas[0].ieps2,
            ieps3: data.ventas[0].ieps3
          }
        });
      }
    } catch (err) {
      console.error('Error al buscar ventas:', err);
      mostrarNotificacion('Error al buscar ventas: ' + err.message, 'error');
      setVentas([]);
    } finally {
      setBuscandoVentas(false); // Ocultar pantalla de carga
    }
  };

  // Modificar la función handleTicketChange
const handleTicketChange = (e) => {
  const { name, value } = e.target;
  
  if (name === 'totalTicket') {
    // Cuando el campo tiene el foco, solo permitir entrada numérica básica
    if (!value.includes('$')) {
      // Permitir solo dígitos y un punto decimal sin formateo adicional
      let numericInput = value.replace(/[^\d.]/g, '');
      
      // Si hay múltiples puntos decimales, quedarse solo con el primero
      if (numericInput.split('.').length > 2) {
        const parts = numericInput.split('.');
        numericInput = parts[0] + '.' + parts.slice(1).join('');
      }
      
      // Guardar el valor tal como lo escribió el usuario
      setTicketData(prev => ({
        ...prev,
        [name]: numericInput
      }));
    } else {
      // Si de alguna manera el valor incluye $, limpiarlo
      const numericValue = value.replace(/[$,\s]/g, '');
      setTicketData(prev => ({
        ...prev,
        [name]: numericValue
      }));
    }
  } else {
    // Para otros campos, mantener comportamiento original
    setTicketData(prev => ({
      ...prev,
      [name]: value
    }));
  }
};

  // Iniciar búsqueda al hacer clic en el botón
  const handleBuscarVenta = () => {
    const claveTicket = ticketData.claveTicket;
    if (!claveTicket) {
      mostrarNotificacion('Por favor, ingrese una clave de ticket válida.', 'error');
      return;
    }
    
    buscarVentas(claveTicket);
  };

  // Añadir estos nuevos estados para las notificaciones toast
  const [notificaciones, setNotificaciones] = useState([]);
  const notificacionIdRef = useRef(1);
  
  // Función para mostrar notificaciones estilo toast
  const mostrarNotificacion = (mensaje, tipo = 'error') => {
    const id = notificacionIdRef.current++;
    const nuevaNotificacion = {
      id,
      mensaje,
      tipo, // 'success' o 'error'
    };
    
    setNotificaciones(prev => [...prev, nuevaNotificacion]);
    
    // Auto-eliminar después de 5 segundos
    setTimeout(() => {
      setNotificaciones(prev => prev.filter(n => n.id !== id));
    }, 5000);
  };
  
  // Función para resetear el formulario después de facturar
  const resetearFormulario = () => {
    // Mantener la lista de empresas y el usuario seleccionado
    // pero resetear todo lo demás
    setEmpresaSeleccionadaId('');
    setEmpresa(null);
    setShowDatosPrincipales(false);
    setShowDatosAdicionales(false);
    setIsEditing(false);
    setVentas([]);
    
    setFormData({
      razonSocial: '',
      direccion: '',
      codigoPostal: '',
      pais: '',
      estado: '',
      localidad: '',
      municipio: '',
      colonia: '',
      observaciones: '',
      rfc: '',
    });
    
    setTicketData({
      claveTicket: '',
      totalTicket: '',
    });
    
    setUsoCfdi('');
  };
  
  // Modificar la función handleFacturar
  const handleFacturar = async () => {
    const errores = [];
  
    [...camposObligatoriosForm, ...camposObligatoriosTicket].forEach((campo) => {
      const valor = campo in ticketData ? ticketData[campo] : formData[campo];
      if (!valor || valor.toString().trim() === '') {
        errores.push(`Por favor completa el campo obligatorio: ${campo}`);
      }
    });
  
    if (!usoCfdi) {
      errores.push('Por favor selecciona el Uso CFDI');
    }
  
    if (errores.length > 0) {
      // Mostrar errores con el nuevo sistema de notificaciones
      errores.forEach(error => mostrarNotificacion(error, 'error'));
      return;
    }
  
    setGenerando(true);

    try {
      const formDataToSend = new FormData();
      const userId = getUserId(); // Obtener el ID del usuario actual
      const facturaData = {
        idempresa: empresa?.id || 0, // ID de la empresa (puede ser 0 si no hay empresa)
        id_usuario: userId, // ID del usuario que está generando la factura
        receptor_rfc: formData.rfc,
        receptor_razon_social: formData.razonSocial,
        direccion: formData.direccion,
        codigo_postal: formData.codigoPostal,
        pais: formData.pais,
        estado: formData.estado ? parseInt(formData.estado, 10) : 0,
        localidad: formData.localidad,
        municipio: formData.municipio,
        colonia: formData.colonia,
        observaciones: formData.observaciones,
        uso_cfdi: usoCfdi,
        regimen_fiscal: empresa?.regimen_fiscal || '',
        clave_ticket: ticketData.claveTicket,
        total: parseFloat(ticketData.totalTicket || 0),
      };

      formDataToSend.append('datos', JSON.stringify(facturaData));

      const response = await fetch('http://localhost:8080/api/generar-factura', {
        method: 'POST',
        body: formDataToSend,
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(`Error al generar la factura: ${response.status} ${errorText}`);
      }

      const blob = await response.blob();

      if (blob.size === 0) {
        throw new Error('El servidor respondió con un archivo vacío');
      }

      // Para descargar la factura
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = 'factura.zip';
      link.style.display = 'none';
      document.body.appendChild(link);
      link.click();

      setTimeout(() => {
        window.URL.revokeObjectURL(url);
        document.body.removeChild(link);
      }, 100);

      // El backend ya guarda automáticamente en el historial al generar la factura
      
      // NUEVO: Guardar las ventas en la base de datos automáticamente
      if (ventas && ventas.length > 0) {
        try {
          await guardarVentasEnBD();
          mostrarNotificacion('Factura generada y ventas guardadas en BD correctamente', 'success');
        } catch (ventasError) {
          console.error('Error al guardar ventas en BD:', ventasError);
          mostrarNotificacion('Factura generada correctamente, pero error al guardar ventas en BD', 'warning');
        }
      } else {
        mostrarNotificacion('Factura generada correctamente', 'success');
      }
      
      // NUEVO: Resetear el formulario después de facturar exitosamente
      resetearFormulario();
      
    } catch (error) {
      // Mostrar error con el nuevo sistema de notificaciones
      mostrarNotificacion('Error al generar la factura: ' + error.message, 'error');
    } finally {
      setGenerando(false);
    }
  };

  // Reemplaza la función guardarEnHistorial

// Función para guardar ventas en la base de datos
const guardarVentasEnBD = async (mostrarNotif = false) => {
  if (!ventas || ventas.length === 0) {
    if (mostrarNotif) {
      mostrarNotificacion('No hay ventas para guardar', 'error');
    }
    throw new Error('No hay ventas para guardar');
  }

  // Solo mostrar loading si se llama manualmente
  if (mostrarNotif) {
    setGenerando(true);
  }
  
  try {
    const datosParaGuardar = {
      serie: ticketData.claveTicket,
      ventas: ventas.map(venta => ({
        clave_producto: (venta.codigo_producto || '').toString(),
        descripcion: venta.producto || '',
        clave_sat: venta.sat_clave || '',
        unidad_sat: venta.sat_medida || '',
        cantidad: Math.floor(venta.cantidad) || 0, // Convertir a entero
        precio_unitario: parseFloat(venta.precio) || 0,
        descuento: parseFloat(venta.descuento) || 0,
        total: parseFloat(venta.total) || 0
      }))
    };

    console.log('Datos a enviar a BD:', datosParaGuardar);

    const response = await fetch('http://localhost:8080/api/ventas/guardar', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(datosParaGuardar)
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Error ${response.status}: ${errorText}`);
    }

    const resultado = await response.json();
    console.log('Respuesta del servidor BD:', resultado);

    if (mostrarNotif) {
      mostrarNotificacion(
        `✅ ${resultado.insertados} ventas guardadas correctamente en la base de datos`, 
        'success'
      );
    }

    return resultado;

  } catch (error) {
    console.error('Error al guardar ventas:', error);
    if (mostrarNotif) {
      mostrarNotificacion('Error al guardar ventas: ' + error.message, 'error');
    }
    throw error; // Re-lanzar el error para que lo maneje handleFacturar
  } finally {
    if (mostrarNotif) {
      setGenerando(false);
    }
  }
};

  return (
    <div className="empresa-container" style={{ marginTop: '60px', marginLeft: '290px' }}>
      {generando && <PantallaDeCarga mensaje="Generando factura, por favor espera..." />}
      {buscandoVentas && <PantallaDeCarga mensaje="Buscando ventas, por favor espera..." />}
      
      {/* Componente de notificaciones estilo toast */}
      <div className="notificaciones-container">
        {notificaciones.map(notif => (
          <div 
            key={notif.id} 
            className={`notificacion notificacion-${notif.tipo}`}
            onClick={() => setNotificaciones(prev => prev.filter(n => n.id !== notif.id))}
          >
            <div className="notificacion-contenido">
              <span className="notificacion-icono">
                {notif.tipo === 'success' ? '✓' : '✕'}
              </span>
              <span className="notificacion-mensaje">{notif.mensaje}</span>
            </div>
          </div>
        ))}
      </div>
      
      <h1 className="titulo">Panel de Facturación</h1>
      
      {/* Selector de empresas */}
      <div className="tarjeta-empresa">
        <h2 className="titulo-empresa">Seleccionar Empresa</h2>
        
        {cargandoEmpresas ? (
          <p className="text-center">Cargando empresas...</p>
        ) : usuarioNoAutenticado ? (
          <div className="no-empresas">
            <p>No hay un usuario activo. Por favor inicie sesión para continuar.</p>
            <Link to="/login" className="boton-registrar">Iniciar Sesión</Link>
          </div>
        ) : error ? (
          <p className="mensaje-error">{error}</p>
        ) : empresas.length === 0 ? (
          <div className="no-empresas">
            <p>
              No tienes ninguna empresa registrada para timbrar, por favor agregue una en la sección de{' '}
              <span 
                className="enlace-administrar"
                onClick={() => {
                  // Dispatch custom event to communicate with parent Home component
                  const event = new CustomEvent('navigateToSection', { 
                    detail: { section: 'administrarEmpresas' } 
                  });
                  window.dispatchEvent(event);
                }}
              >
                Administrar Empresas
              </span>.
            </p>
          </div>
        ) : (
          <div className="selector-empresas">
            <label className="etiqueta" htmlFor="empresa-select">
              Empresa a facturar <span className="required-mark">*</span>
            </label>
            <select
              id="empresa-select"
              value={empresaSeleccionadaId}
              onChange={seleccionarEmpresa}
              className="input-rfc"
              required
            >
              <option value="">Selecciona una empresa</option>
              {empresas.map((emp) => (
                <option key={emp.id} value={emp.id}>
                  {emp.razon_social} - {emp.rfc}
                </option>
              ))}
            </select>
          </div>
        )}
      </div>

      {empresa && showDatosPrincipales && (
        <div className="tarjeta-empresa animate-card">
          <h2 className="titulo-empresa">Datos Principales</h2>
          <div className="header-actions">
            {!isEditing ? (
              <button 
                className="btn-editar" 
                onClick={handleEditMode}
              >
                Editar Información
              </button>
            ) : (
              <button 
                className="btn-cancelar" 
                onClick={handleCancelEdit}
              >
                Cancelar
              </button>
            )}
          </div>
          <div className="info-principal fila-campos">
            <div className="campo-info">
              <label className="etiqueta">RFC:</label>
              {isEditing ? (
                <input
                  type="text"
                  name="rfc"
                  value={formData.rfc}
                  onChange={handleInputChange}
                  className="input-rfc"
                />
              ) : (
                <div className="read-only-field">{formData.rfc}</div>
              )}
            </div>
            <div className="campo-info">
              <label className="etiqueta">Razón Social:</label>
              {isEditing ? (
                <input
                  type="text"
                  name="razonSocial"
                  value={formData.razonSocial}
                  onChange={handleInputChange}
                  className="input-rfc"
                />
              ) : (
                <div className="read-only-field">{formData.razonSocial}</div>
              )}
            </div>
          </div>
        </div>
      )}

      {empresa && showDatosAdicionales && (
        <div className="tarjeta-empresa">
          <h2 className="titulo-empresa">Datos Adicionales</h2>
          <div className="info-principal grid-campos">
            {/* Mostrar todos los campos excepto observaciones */}
            {[
              'direccion', 'codigoPostal', 'pais', 'estado',
              'localidad', 'municipio', 'colonia'
            ].map((campo) => (
              <div className="campo-adicional" key={campo}>
                <label className="etiqueta-adicional" htmlFor={campo}>
                  {campo === 'codigoPostal' ? 'Código Postal' : 
                   campo.charAt(0).toUpperCase() + campo.slice(1).replace(/([A-Z])/g, ' $1')}
                  {camposObligatoriosForm.includes(campo) ? (
                    <span style={{ color: 'red' }}> *</span>
                  ) : null}
                </label>
                {isEditing ? (
                  <input
                    type="text"
                    id={campo}
                    name={campo}
                    value={formData[campo]}
                    onChange={handleInputChange}
                    className="input-rfc"
                  />
                ) : (
                  <div className="read-only-field">{formData[campo]}</div>
                )}
              </div>
            ))}
            
            {/* CFDI como campo de ancho completo */}
            <div className="campo-adicional campo-completo">
              <label className="etiqueta-adicional" htmlFor="usoCfdi">
                Uso CFDI <span style={{ color: 'red' }}>*</span>
              </label>
              <select
                id="usoCfdi"
                value={usoCfdi}
                onChange={(e) => setUsoCfdi(e.target.value)}
                className="input-rfc"
                required
              >
                <option value="">Seleccione una opción</option>
                <option value="G03">G03 - Gastos en general</option>
                <option value="G01">G01 - Adquisición de mercancías</option>
                <option value="G02">G02 - Devoluciones, descuentos o bonificaciones</option>
                <option value="P01">P01 - Por definir</option>
                <option value="I01">I01 - Construcciones</option>
                <option value="I02">I02 - Mobiliario y equipo</option>
                <option value="I03">I03 - Equipo de transporte</option>
                <option value="I04">I04 - Equipo de cómputo</option>
                <option value="D01">D01 - Honorarios médicos, dentales y gastos hospitalarios</option>
                <option value="D02">D02 - Gastos médicos por incapacidad o discapacidad</option>
                <option value="D03">D03 - Gastos funerarios</option>
                <option value="D04">D04 - Donativos</option>
                <option value="D05">D05 - Intereses reales por créditos hipotecarios</option>
                <option value="D06">D06 - Aportaciones voluntarias al SAR</option>
                <option value="D07">D07 - Primas por seguros de gastos médicos</option>
                <option value="D08">D08 - Gastos de transporte escolar obligatorio</option>
                <option value="D09">D09 - Depósitos en cuentas para el ahorro, plans de retiro</option>
                <option value="D10">D10 - Pagos por servicios educativos (colegiaturas)</option>
              </select>
            </div>
            
            {/* Observaciones como campo de ancho completo */}
            <div className="campo-adicional campo-completo">
              <label className="etiqueta-adicional" htmlFor="observaciones">
                Observaciones
              </label>
              <input
                type="text"
                id="observaciones"
                name="observaciones"
                value={formData.observaciones}
                onChange={handleInputChange}
                className="input-rfc"
                placeholder="Observaciones (opcional)"
              />
            </div>
          </div>
          
          {isEditing && (
            <div className="form-actions">
              <button 
                type="button" 
                className="btn-guardar"
                onClick={handleSaveChanges}
              >
                Guardar Cambios
              </button>
            </div>
          )}
        </div>
      )}

      {empresa && (
        <div className="tarjeta-empresa">
          <h2 className="titulo-empresa">Datos del Ticket</h2>
          <div className="info-principal campos-ticket">
            <div className="campo-info">
              <label className="etiqueta" htmlFor="claveTicket">
                Clave del Ticket <span style={{ color: 'red' }}>*</span>
              </label>
              <div className="input-con-boton">
                <input
                  type="text"
                  id="claveTicket"
                  name="claveTicket"
                  value={ticketData.claveTicket}
                  onChange={handleTicketChange}
                  className="input-rfc"
                />
                <button 
                  type="button"
                  onClick={handleBuscarVenta}
                  className="boton-buscar-ticket"
                  disabled={buscandoVentas}
                >
                  {buscandoVentas ? "Buscando..." : "Buscar"}
                </button>
              </div>
            </div>
            <div className="campo-info">
              <label className="etiqueta" htmlFor="totalTicket">
                Total del Ticket <span style={{ color: 'red' }}>*</span>
              </label>
              <input
                type="text"
                id="totalTicket"
                name="totalTicket"
                value={document.activeElement === document.getElementById('totalTicket') 
                  ? ticketData.totalTicket 
                  : ticketData.totalTicket ? formatCurrency(ticketData.totalTicket) : ''}
                onChange={handleTicketChange}
                onFocus={(e) => {
                  // Al hacer foco, mostrar solo el valor numérico sin ningún formato
                  e.target.value = ticketData.totalTicket;
                }}
                onBlur={(e) => {
                  // Solo formatear cuando pierde el foco, NO durante la edición
                  if (ticketData.totalTicket) {
                    e.target.value = formatCurrency(ticketData.totalTicket);
                  } else {
                    e.target.value = '';
                  }
                }}
                className="input-rfc"
                placeholder="$ 0.00"
              />
            </div>
          </div>

          {ventas.length > 0 && (
            <div className="tabla-productos">
              <h3 className="titulo-productos">Detallado de productos</h3>
              {/* Tabla corregida sin columna de acciones extra */}
<table className="tabla-ventas">
  <thead>
    <tr>
      <th>Clave Prod/Ser</th>
      <th>Producto</th>
      <th>Clave SAT</th>  
      <th>Unidad SAT</th> 
      <th>Cantidad</th>
      <th>Precio</th>
      <th>IVA (%)</th>
      <th>IEPS (%)</th>
      <th>Descuento</th>
    </tr>
  </thead>
  <tbody>
    {ventas.map((venta, index) => {
      // Obtener los valores de impuestos del producto
      const ivaProducto = venta.iva || 0;
      const ieps1Producto = venta.ieps1 || 0;
      const ieps2Producto = venta.ieps2 || 0;
      const ieps3Producto = venta.ieps3 || 0;
      const totalIEPS = ieps1Producto + ieps2Producto + ieps3Producto;
      
      return (
        <tr key={index}>
          <td>{venta.codigo_producto || 'N/A'}</td>  
          <td>{venta.producto}</td>                  
          <td>{venta.sat_clave || 'No disponible'}</td>     
          <td>{venta.sat_medida || 'No disponible'}</td>
          <td>{venta.cantidad}</td>
          <td>${venta.precio.toFixed(2)}</td>
          <td style={{
            backgroundColor: ivaProducto > 0 ? '#e8f5e8' : '#fff3cd',
            fontWeight: 'bold',
            color: ivaProducto > 0 ? '#2d5016' : '#856404'
          }}>
            {ivaProducto > 0 ? `${ivaProducto.toFixed(1)}%` : '0.0%'}
          </td>
          <td style={{
            backgroundColor: totalIEPS > 0 ? '#e1f5fe' : '#f5f5f5',
            fontWeight: totalIEPS > 0 ? 'bold' : 'normal',
            color: totalIEPS > 0 ? '#01579b' : '#666'
          }}>
            {totalIEPS > 0 ? `${totalIEPS.toFixed(1)}%` : '0.0%'}
          </td>
          <td>${venta.descuento.toFixed(2)}</td>
        </tr>
      );
    })}
    {/* Fila del total */}
    <tr style={{ backgroundColor: '#f8f9fa', fontWeight: 'bold', borderTop: '2px solid #dee2e6' }}>
      <td colSpan="8" style={{ textAlign: 'right', padding: '12px' }}>
        TOTAL DE TODOS LOS PRODUCTOS:
      </td>
      <td style={{ fontSize: '16px', color: '#28a745' }}>
        ${ventas.reduce((total, venta) => total + venta.total, 0).toFixed(2)}
      </td>
    </tr>
  </tbody>
</table>

{/* Tabla de diagnóstico temporal - solo visible cuando hay datos */}
{ventas.length > 0 && (
  <div style={{ marginTop: '20px', padding: '15px', backgroundColor: '#f8f9fa', borderRadius: '5px' }}>
    <h4 style={{ color: '#6c757d', marginBottom: '15px' }}>🔍 Diagnóstico de Configuración</h4>
    <div style={{ maxHeight: '300px', overflowY: 'auto' }}>
      <table style={{ width: '100%', fontSize: '12px', border: '1px solid #dee2e6' }}>
        <thead>
          <tr style={{ backgroundColor: '#e9ecef' }}>
            <th style={{ padding: '8px', border: '1px solid #dee2e6' }}>Producto</th>
            <th style={{ padding: '8px', border: '1px solid #dee2e6' }}>Estado Config</th>
            <th style={{ padding: '8px', border: '1px solid #dee2e6' }}>Configs Imp</th>
            <th style={{ padding: '8px', border: '1px solid #dee2e6' }}>Empresa Prod</th>
            <th style={{ padding: '8px', border: '1px solid #dee2e6' }}>Empresa Imp</th>
            <th style={{ padding: '8px', border: '1px solid #dee2e6' }}>IVA</th>
            <th style={{ padding: '8px', border: '1px solid #dee2e6' }}>IEPS Total</th>
            <th style={{ padding: '8px', border: '1px solid #dee2e6' }}>SAT Clave</th>
            <th style={{ padding: '8px', border: '1px solid #dee2e6' }}>SAT Medida</th>
          </tr>
        </thead>
        <tbody>
          {ventas.map((venta, index) => {
            const totalIEPS = (venta.ieps1 || 0) + (venta.ieps2 || 0) + (venta.ieps3 || 0);
            const tieneProblema = venta.diagnostico_config !== 'Configuración completa';
            const tieneMultiplesConfigs = venta.cantidad_configs > 1;
            
            return (
              <tr key={index} style={{ 
                backgroundColor: tieneProblema ? '#fff3cd' : (tieneMultiplesConfigs ? '#e1f5fe' : '#d4edda')
              }}>
                <td style={{ padding: '6px', border: '1px solid #dee2e6', fontSize: '11px' }}>
                  {venta.producto.substring(0, 30)}...
                </td>
                <td style={{ 
                  padding: '6px', 
                  border: '1px solid #dee2e6',
                  fontWeight: 'bold',
                  color: tieneProblema ? '#856404' : (tieneMultiplesConfigs ? '#0c5460' : '#155724')
                }}>
                  {venta.diagnostico_config || 'N/A'}
                </td>
                <td style={{ 
                  padding: '6px', 
                  border: '1px solid #dee2e6',
                  fontWeight: tieneMultiplesConfigs ? 'bold' : 'normal',
                  color: tieneMultiplesConfigs ? '#d32f2f' : '#333'
                }}>
                  {venta.cantidad_configs || 0}
                  {tieneMultiplesConfigs && ' ⚠️'}
                </td>
                <td style={{ padding: '6px', border: '1px solid #dee2e6' }}>
                  {venta.empresa_producto || 'N/A'}
                </td>
                <td style={{ padding: '6px', border: '1px solid #dee2e6' }}>
                  {venta.empresa_impuesto || 'N/A'}
                </td>
                <td style={{ padding: '6px', border: '1px solid #dee2e6' }}>
                  {(venta.iva || 0).toFixed(1)}%
                </td>
                <td style={{ padding: '6px', border: '1px solid #dee2e6' }}>
                  {totalIEPS.toFixed(1)}%
                </td>
                <td style={{ padding: '6px', border: '1px solid #dee2e6' }}>
                  {venta.sat_clave || 'No config'}
                </td>
                <td style={{ padding: '6px', border: '1px solid #dee2e6' }}>
                  {venta.sat_medida || 'No config'}
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
    <div style={{ marginTop: '10px', fontSize: '11px', color: '#6c757d' }}>
      <strong>Leyenda:</strong> 
      <span style={{ backgroundColor: '#d4edda', padding: '2px 6px', marginLeft: '5px' }}>Verde = Configuración completa</span>
      <span style={{ backgroundColor: '#e1f5fe', padding: '2px 6px', marginLeft: '5px' }}>Azul = Múltiples configuraciones (se usa MAX)</span>
      <span style={{ backgroundColor: '#fff3cd', padding: '2px 6px', marginLeft: '5px' }}>Amarillo = Problemas de configuración</span>
    </div>
  </div>
)}
              
              {/* Información sobre productos para facturar */}
              <div className="info-guardado-automatico" style={{
                marginTop: '15px',
                textAlign: 'center',
                padding: '10px',
                backgroundColor: '#e7f3ff',
                border: '1px solid #b3d7ff',
                borderRadius: '8px',
                fontSize: '0.9em',
                color: '#0066cc'
              }}>
                <span style={{ fontSize: '0.8em', color: '#666' }}>
                  {ventas.length} producto{ventas.length !== 1 ? 's' : ''} listo{ventas.length !== 1 ? 's' : ''} para facturar
                </span>
              </div>
            </div>
          )}

          {/* Tabla de diagnóstico temporal - solo visible cuando hay datos */}
          {ventas.length > 0 && (
            <div style={{ marginTop: '20px', padding: '15px', backgroundColor: '#f8f9fa', borderRadius: '5px' }}>
              <h4 style={{ color: '#6c757d', marginBottom: '15px' }}>🔍 Diagnóstico de Configuración</h4>
              <div style={{ maxHeight: '300px', overflowY: 'auto' }}>
                <table style={{ width: '100%', fontSize: '12px', border: '1px solid #dee2e6' }}>
                  <thead>
                    <tr style={{ backgroundColor: '#e9ecef' }}>
                      <th style={{ padding: '8px', border: '1px solid #dee2e6' }}>Producto</th>
                      <th style={{ padding: '8px', border: '1px solid #dee2e6' }}>Estado Config</th>
                      <th style={{ padding: '8px', border: '1px solid #dee2e6' }}>Empresa Prod</th>
                      <th style={{ padding: '8px', border: '1px solid #dee2e6' }}>Empresa Imp</th>
                      <th style={{ padding: '8px', border: '1px solid #dee2e6' }}>IVA</th>
                      <th style={{ padding: '8px', border: '1px solid #dee2e6' }}>IEPS Total</th>
                      <th style={{ padding: '8px', border: '1px solid #dee2e6' }}>SAT Clave</th>
                      <th style={{ padding: '8px', border: '1px solid #dee2e6' }}>SAT Medida</th>
                    </tr>
                  </thead>
                  <tbody>
                    {ventas.map((venta, index) => {
                      const totalIEPS = (venta.ieps1 || 0) + (venta.ieps2 || 0) + (venta.ieps3 || 0);
                      const tieneProblema = venta.diagnostico_config !== 'Configuración completa';
                      
                      return (
                        <tr key={index} style={{ 
                          backgroundColor: tieneProblema ? '#fff3cd' : '#d4edda' 
                        }}>
                          <td style={{ padding: '6px', border: '1px solid #dee2e6', fontSize: '11px' }}>
                            {venta.producto.substring(0, 30)}...
                          </td>
                          <td style={{ 
                            padding: '6px', 
                            border: '1px solid #dee2e6',
                            fontWeight: 'bold',
                            color: tieneProblema ? '#856404' : '#155724'
                          }}>
                            {venta.diagnostico_config || 'N/A'}
                          </td>
                          <td style={{ padding: '6px', border: '1px solid #dee2e6' }}>
                            {venta.empresa_producto || 'N/A'}
                          </td>
                          <td style={{ padding: '6px', border: '1px solid #dee2e6' }}>
                            {venta.empresa_impuesto || 'N/A'}
                          </td>
                          <td style={{ padding: '6px', border: '1px solid #dee2e6' }}>
                            {(venta.iva || 0).toFixed(1)}%
                          </td>
                          <td style={{ padding: '6px', border: '1px solid #dee2e6' }}>
                            {totalIEPS.toFixed(1)}%
                          </td>
                          <td style={{ padding: '6px', border: '1px solid #dee2e6' }}>
                            {venta.sat_clave || 'No config'}
                          </td>
                          <td style={{ padding: '6px', border: '1px solid #dee2e6' }}>
                            {venta.sat_medida || 'No config'}
                          </td>
                        </tr>
                      );
                    })}
                  </tbody>
                </table>
              </div>
              <div style={{ marginTop: '10px', fontSize: '11px', color: '#6c757d' }}>
                <strong>Leyenda:</strong> 
                <span style={{ backgroundColor: '#d4edda', padding: '2px 6px', marginLeft: '5px' }}>Verde = Configuración completa</span>
                <span style={{ backgroundColor: '#fff3cd', padding: '2px 6px', marginLeft: '5px' }}>Amarillo = Problemas de configuración</span>
              </div>
            </div>
          )}
        </div>
      )}

      {empresa && (
        <div className="boton-confirmar">
          <button 
            className="boton-facturar" 
            onClick={handleFacturar}
            disabled={generando}
          >
            {generando ? "Generando factura..." : "FACTURAR"}
          </button>
        </div>
      )}

    </div>
  );
}

export default InicioFacturacion;