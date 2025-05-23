import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import '../styles/Pantalla_Principal.css';

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
  const [modalError, setModalError] = useState('');
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
      setModalError(['Datos actualizados correctamente']);
    } catch (error) {
      setModalError(['Error al actualizar los datos: ' + error.message]);
    }
  };

  // Componente de pantalla de carga con mensaje personalizable
  const PantallaDeCarga = ({ mensaje = "Generando factura, por favor espera..." }) => (
    <div className="pantalla-carga">
      <div className="spinner"></div>
      <p>{mensaje}</p>
    </div>
  );

  const ModalError = ({ texto, onClose }) => {
    if (!texto || texto.length === 0) return null;
  
    return (
      <div className="modal-overlay">
        <div className="modal-contenido">
          <ul className="lista-errores">
            {Array.isArray(texto) ? (
              texto.map((error, index) => (
                <li key={index}>
                  <span className="punto-error">•</span> {error}
                </li>
              ))
            ) : (
              <li>
                <span className="punto-error">•</span> {texto}
              </li>
            )}
          </ul>
          <button className="modal-boton" onClick={onClose}>
            Cerrar
          </button>
        </div>
      </div>
    );
  };

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
      setModalError(['La clave del ticket debe tener al menos 32 caracteres alfanuméricos.']);
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
      setVentas(data.ventas || []);
      
      // Si no hay ventas, mostrar mensaje
      if (!data.ventas || data.ventas.length === 0) {
        setModalError(['No se encontraron ventas para la clave de ticket proporcionada']);
      }

      console.log("Datos completos:", data);
      if (data.ventas && data.ventas.length > 0) {
        console.log("Ejemplo de venta:", data.ventas[0]);
        console.log("Propiedades disponibles:", Object.keys(data.ventas[0]));
        console.log("Valor de clave_producto:", data.ventas[0].clave_producto);
        console.log("Valor de claveProducto:", data.ventas[0].claveProducto);
      }

      // En la función buscarVentas, después de recibir la respuesta:
      console.log("Datos de ventas completos:", data.ventas);
      if (data.ventas && data.ventas.length > 0) {
        console.log("Primera venta - propiedades:", Object.keys(data.ventas[0]));
        console.log("Primera venta - sat_clave:", data.ventas[0].sat_clave);
        console.log("Primera venta - sat_medida:", data.ventas[0].sat_medida);
      }
    } catch (err) {
      console.error('Error al buscar ventas:', err);
      setModalError(['Error al buscar ventas: ' + err.message]);
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
      setModalError(['Por favor, ingrese una clave de ticket válida.']);
      return;
    }
    
    buscarVentas(claveTicket);
  };

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
      setModalError(errores);
      return;
    }
  
    setGenerando(true);
    setModalError([]);

    try {
      const formDataToSend = new FormData();
      const facturaData = {
        rfc: formData.rfc,
        razon_social: formData.razonSocial,
        direccion: formData.direccion,
        codigo_postal: formData.codigoPostal,
        pais: formData.pais,
        estado: formData.estado,
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

      const response = await fetch('http://localhost:8080/api/generar_factura', {
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

      // Guardar en el historial después de generar la factura exitosamente
      await guardarEnHistorial(facturaData);
      
      // Mensaje de éxito
      setModalError(['Factura generada correctamente y guardada en el historial']);
      
    } catch (error) {
      setModalError(['Error al generar la factura: ' + error.message]);
    } finally {
      setGenerando(false);
    }
  };

  // Función para guardar la factura en el historial
  const guardarEnHistorial = async (facturaData) => {
    try {
      let userId = getUserId();

      // Verificar que el usuario existe en la base de datos
      try {
        const userResponse = await fetch(`http://localhost:8080/api/usuarios/${userId}`);
        if (!userResponse.ok) {
          console.error(`El usuario con ID ${userId} no existe en la base de datos`);
          // Usar un ID de usuario por defecto que sepas que existe
          userId = 1; // O cualquier ID que sepas que existe
        }
      } catch (error) {
        console.error("Error al verificar usuario:", error);
      }
      
      // Versión mejorada con verificación
      let descripcionProductos = "";
if (ventas && ventas.length > 0) {
  // Solo incluir los nombres de los productos, nada más
  descripcionProductos = ventas.map(venta => 
    venta.producto || 'Producto sin nombre'
  ).join('; ');
  console.log("DESCRIPCIÓN GENERADA: " + descripcionProductos);
} else {
  descripcionProductos = "No hay detalle de productos disponible";
  console.log("NO HAY VENTAS PARA DESCRIPCIÓN");
}
      
      console.log("Descripción generada:", descripcionProductos);
      
      const historialData = {
        id_usuario: userId,
        rfc_receptor: facturaData.rfc,
        razon_social_receptor: facturaData.razon_social,
        clave_ticket: facturaData.clave_ticket,
        total: facturaData.total,
        uso_cfdi: facturaData.uso_cfdi,
        observaciones: facturaData.observaciones || '',
        estado: 'G'
      };
      
      console.log("Enviando datos al historial:", historialData);
      console.log("JSON a enviar:", JSON.stringify(historialData, null, 2));
      
      const response = await fetch('http://localhost:8080/api/historial-facturas', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(historialData),
      });
      
      const respuestaTexto = await response.text();
      console.log("Respuesta del servidor:", respuestaTexto);
      
      if (!response.ok) {
        console.error('Error al guardar en historial:', respuestaTexto);
      }
    } catch (error) {
      console.error('Error al guardar en historial:', error);
    }
  };

  return (
    <div className="empresa-container" style={{ marginTop: '60px', marginLeft: '290px' }}>
      {generando && <PantallaDeCarga mensaje="Generando factura, por favor espera..." />}
      {buscandoVentas && <PantallaDeCarga mensaje="Buscando ventas, por favor espera..." />}
      <ModalError texto={modalError} onClose={() => setModalError('')} />
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
      <th>Descuento</th>
      <th>Total</th>
    </tr>
  </thead>
  <tbody>
    {ventas.map((venta, index) => (
      <tr key={index}>
        <td>{venta.codigo_producto || 'N/A'}</td>  
        <td>{venta.producto}</td>                  
        <td>{venta.sat_clave || 'No disponible'}</td>     
        <td>{venta.sat_medida || 'No disponible'}</td>
        <td>{venta.cantidad}</td>
        <td>${venta.precio.toFixed(2)}</td>
        <td>${venta.descuento.toFixed(2)}</td>
        <td>${venta.total.toFixed(2)}</td>
      </tr>
    ))}
  </tbody>
</table>
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