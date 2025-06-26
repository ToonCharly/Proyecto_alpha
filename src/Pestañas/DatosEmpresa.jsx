import React, { useState, useEffect } from 'react';
// Eliminar la importación de useNavigate si no se va a usar
// import { useNavigate } from 'react-router-dom';
import '../STYLES/DatosEmpresa.css';

function DatosEmpresa() {
  const [isEditing, setIsEditing] = useState(true); 
  const [datosFiscales, setDatosFiscales] = useState({
    claveSat: '',
    rfcEmisor: '',
    razonSocial: '',
    cp: '',
    direccionFiscal: '',
    claveArchivoCSD: '',
    regimenFiscal: ''
  });
  
  // Estados para el buscador de empresas
  const [rfcBusqueda, setRfcBusqueda] = useState('');
  const [empresaEncontrada, setEmpresaEncontrada] = useState(null);
  const [mostrarModalEmpresa, setMostrarModalEmpresa] = useState(false);
  const [buscandoEmpresa, setBuscandoEmpresa] = useState(false);
  const [errorBusqueda, setErrorBusqueda] = useState('');
  
  // Estados para manejar los archivos
  const [csdKey, setCsdKey] = useState(null);
  const [csdCer, setCsdCer] = useState(null);
  
  // Estado para mensajes y carga
  const [mensaje, setMensaje] = useState(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [fieldErrors, setFieldErrors] = useState({});
  const [modalErrores, setModalErrores] = useState([]);

  // Obtener datos del usuario actual
  const [userData, setUserData] = useState(null);
  const [isAdmin, setIsAdmin] = useState(false); // Añadir este estado después de tus declaraciones de estado existentes

  useEffect(() => {
    const storedUserData = sessionStorage.getItem('userData'); // Migrado a sessionStorage para evitar conflictos
    if (storedUserData) {
      try {
        const parsedUserData = JSON.parse(storedUserData);
        setUserData(parsedUserData);
        
        // Verificar si el usuario es administrador
        const isUserAdmin = parsedUserData.role === 'admin';
        setIsAdmin(isUserAdmin);
        
        // Si no es admin, mostrar mensaje
        if (!isUserAdmin) {
          setMensaje({
            tipo: 'error',
            texto: 'Solo los administradores pueden editar los datos fiscales'
          });
          setIsEditing(false); // Desactivar edición para no-admin
        } else if (parsedUserData.email || parsedUserData.username) {
          try {
            fetchDatosFiscales(parsedUserData.id || parsedUserData.email || parsedUserData.username);
          } catch {
            console.log("No hay datos previos o no se pudo conectar - modo edición activado");
          }
        }
      } catch (err) {
        console.error("Error al procesar datos del usuario:", err);
      }
    }
  }, []);

  // Función para obtener datos fiscales existentes - silenciosa para carga inicial
  const fetchDatosFiscales = async (identifier) => {
    try {
      setIsSubmitting(true);
      // Usar el parámetro identifier en lugar de ignorarlo
      const response = await fetch(`http://localhost:8080/api/datos-fiscales?id_usuario=${identifier || ''}`);
      if (response.ok) {
        const data = await response.json();
        setDatosFiscales({
          claveSat: data.claveSat || '',
          rfcEmisor: data.rfcEmisor || '',
          razonSocial: data.razonSocial || '',
          cp: data.cp || '',
          direccionFiscal: data.direccionFiscal || '',
          claveArchivoCSD: data.claveArchivoCSD || '',
          regimenFiscal: data.regimenFiscal || ''
        });
        setIsEditing(false); // Solo salimos del modo edición si hay datos
      } else {
        // Si no hay datos, mantenemos el modo edición sin mostrar error
        setIsEditing(true);
      }
    } catch {
      setIsEditing(true);
    } finally {
      setIsSubmitting(false);
    }
  };

  // Manejar cambios en los campos de texto
  const handleChange = (e) => {
    const { name, value } = e.target;
    setDatosFiscales(prev => ({
      ...prev,
      [name]: value
    }));

    // Limpiar errores cuando el usuario modifica un campo
    if (fieldErrors[name]) {
      setFieldErrors(prev => {
        const newErrors = { ...prev };
        delete newErrors[name];
        return newErrors;
      });
    }
  };

  // Manejar cambios en los archivos
  const handleFileChange = (e, fileType) => {
    const file = e.target.files[0];
    if (fileType === 'key') {
      setCsdKey(file);
    } else if (fileType === 'cer') {
      setCsdCer(file);
    }
  };

  // Manejar envío del formulario
  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsSubmitting(true);
    setModalErrores([]);
    setMensaje(null); // Limpiar mensaje anterior
    setFieldErrors({});
    
    // Verificar si el usuario es administrador
    if (!isAdmin) {
      setMensaje({
        tipo: 'error',
        texto: 'Solo los administradores pueden editar los datos fiscales'
      });
      return;
    }
    
    // Validación de campos requeridos
    const errores = {};
    const mensajesError = [];
    
    // Validar campos obligatorios con mensajes específicos
    if (!datosFiscales.rfcEmisor) {
      errores.rfcEmisor = true;
      mensajesError.push('El RFC del Emisor es obligatorio');
    }
    
    if (!datosFiscales.razonSocial) {
      errores.razonSocial = true;
      mensajesError.push('La Razón Social es obligatoria');
    }
    
    if (!datosFiscales.cp) {
      errores.cp = true;
      mensajesError.push('El Código Postal es obligatorio');
    }
    
    if (!datosFiscales.regimenFiscal) {
      errores.regimenFiscal = true;
      mensajesError.push('Debe seleccionar un Régimen Fiscal');
    }
    
    // Validación de archivos solo si es un nuevo registro y no hay RFC
    if (!datosFiscales.rfcEmisor) {
      if (!csdKey) {
        errores.csdKey = true;
        mensajesError.push('Debe cargar el archivo CSD KEY');
      }
      if (!csdCer) {
        errores.csdCer = true;
        mensajesError.push('Debe cargar el archivo CSD CER');
      }
    }
    
    // Si hay errores, mostrarlos y detener el envío
    if (Object.keys(errores).length > 0) {
      setFieldErrors(errores);
      setIsSubmitting(false);
      
      // Crear mensajes de error para la modal
      const mensajesModal = mensajesError.map(mensaje => ({ mensaje }));
      setModalErrores(mensajesModal);
      
      // Mensaje general para el banner
      setMensaje({ 
        tipo: 'error', 
        texto: 'Por favor completa todos los campos marcados como obligatorios' 
      });
      return;
    }
    
    try {
      // Crear FormData para enviar los archivos
      const formDataToSend = new FormData();
      
      // AÑADIR ESTA LÍNEA - Incluir el ID del usuario actual
      // Verificación explícita de userData
      if (!userData) {
        console.error("Error crítico: userData es null o undefined");
        setMensaje({
          tipo: 'error',
          texto: 'Sesión inválida. Por favor, vuelve a iniciar sesión.'
        });
        setIsSubmitting(false);
        return;
      }

      // Verificación explícita del ID
      if (!userData.id) {
        console.error("Error crítico: userData no contiene ID", userData);
        setMensaje({
          tipo: 'error',
          texto: 'No se pudo identificar el usuario. Por favor, vuelve a iniciar sesión.'
        });
        setIsSubmitting(false);
        return;
      }

      // Convertir a entero y verificar
      const userId = parseInt(userData.id, 10);
      if (isNaN(userId) || userId <= 0) {
        console.error("Error crítico: ID de usuario inválido", userData.id);
        setMensaje({
          tipo: 'error',
          texto: 'ID de usuario inválido. Por favor, vuelve a iniciar sesión.'
        });
        setIsSubmitting(false);
        return;
      }
      
      formDataToSend.append('id_usuario', userId);
      console.log("ID convertido a número:", userId);
      
      // Resto de los campos que ya estás enviando
      formDataToSend.append('rfc', datosFiscales.rfcEmisor);
      formDataToSend.append('razon_social', datosFiscales.razonSocial);
      formDataToSend.append('direccion_fiscal', datosFiscales.direccionFiscal);
      formDataToSend.append('codigo_postal', datosFiscales.cp);
      formDataToSend.append('clave_csd', datosFiscales.claveArchivoCSD);
      formDataToSend.append('regimen_fiscal', datosFiscales.regimenFiscal);
      
      // Agregar archivos solo si se seleccionaron
      if (csdKey) formDataToSend.append('csdKey', csdKey);
      if (csdCer) formDataToSend.append('csdCer', csdCer);
      
      // Inspeccionar el FormData
      console.log("Contenido del FormData a enviar:");
      for (let [key, value] of formDataToSend.entries()) {
        console.log(`${key}: ${value}`);
      }
      
      // Enviar los datos al servidor
      const response = await fetch('http://localhost:8080/api/actualizar-datos-fiscales', {
        method: 'POST',
        body: formDataToSend
      });
      
      // Mejorar el manejo de errores
      if (!response.ok) {
        // Obtener el texto completo de la respuesta para mejor diagnóstico
        const responseText = await response.text();
        console.error("Respuesta de error del servidor:", responseText);
        
        try {
          // Intentar parsear como JSON si es posible
          const errorData = JSON.parse(responseText);
          throw new Error(errorData.error || 'Error en los datos enviados');
        } catch {
          // Quitar el parámetro '_' completamente si no lo usas
          throw new Error(`Error ${response.status}: ${responseText || 'Error al guardar los datos'}`);
        }
      } else {
        // Procesar la respuesta exitosa
        await response.json(); // Procesar la respuesta sin asignarla
        
        // Actualizar estado local con los datos que acabas de guardar 
        // en lugar de esperar a que se carguen desde el servidor
        setDatosFiscales(prev => ({
          ...prev,
          // No reinicies los valores - mantén los actuales que el usuario acaba de guardar
        }));
        
        setIsEditing(false); // Cambiar a modo visualización
        setMensaje({
          tipo: 'exito',
          texto: 'Datos fiscales guardados correctamente'
        });
        
        // Eliminar el timeout que está causando problemas
        // setTimeout(() => {
        //   fetchDatosFiscales(userData.id);
        // }, 1000);
      }
    } catch (error) {
      console.error('Error al guardar datos:', error);
      
      // Solo mostrar el mensaje de error específico, sin mencionar problemas de conexión
      setModalErrores([{ mensaje: 'Revisa los datos ingresados e intenta nuevamente' }]);
      setMensaje({ 
        tipo: 'error', 
        texto: 'No se pudieron guardar los datos. Verifica que todos los campos sean correctos.' 
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  // Función auxiliar para obtener la etiqueta del régimen fiscal
  const getRegimenFiscalLabel = (clave) => {
    const regimenes = {
      '601': 'General de Ley Personas Morales',
      '603': 'Personas Morales con Fines no Lucrativos',
      '605': 'Sueldos y Salarios e Ingresos Asimilados a Salarios',
      '606': 'Arrendamiento',
      '608': 'Demás ingresos',
      '612': 'Personas Físicas con Actividades Empresariales y Profesionales',
      '620': 'Sociedades Cooperativas de Producción que optan por diferir sus ingresos',
      '621': 'Incorporación Fiscal',
      '622': 'Actividades Agrícolas, Ganaderas, Silvícolas y Pesqueras',
      '625': 'Régimen de las Actividades Empresariales con ingresos a través de Plataformas Tecnológicas',
      '626': 'Régimen Simplificado de Confianza'
    };
    
    return regimenes[clave] || '';
  };

  // Función para buscar empresa por RFC
  const buscarEmpresaPorRFC = async () => {
    if (!rfcBusqueda.trim()) {
      setErrorBusqueda('Por favor ingresa un RFC');
      return;
    }

    setBuscandoEmpresa(true);
    setErrorBusqueda('');

    try {
      // Paso 1: Buscar en adm_empresas_rfc para obtener el idempresa
      const responseRFC = await fetch(`http://localhost:8080/api/buscar-empresa-rfc?rfc=${rfcBusqueda.trim()}`);
      
      if (!responseRFC.ok) {
        throw new Error('RFC no encontrado en el sistema');
      }

      const dataRFC = await responseRFC.json();
      
      if (!dataRFC.idempresa) {
        throw new Error('No se encontró una empresa asociada a este RFC');
      }

      // Paso 2: Buscar en adm_empresa usando el idempresa
      const responseEmpresa = await fetch(`http://localhost:8080/api/empresa-detalle?idempresa=${dataRFC.idempresa}`);
      
      if (!responseEmpresa.ok) {
        throw new Error('No se pudieron obtener los detalles de la empresa');
      }

      const empresaData = await responseEmpresa.json();
      
      // Guardar la información de la empresa encontrada
      setEmpresaEncontrada({
        ...empresaData,
        rfc: rfcBusqueda.trim(), // Incluir el RFC buscado
        metodo_pago: dataRFC.metodo_pago, // Incluir el método de pago del RFC
        c_regimenfiscal: dataRFC.c_regimenfiscal, // Incluir la clave del régimen fiscal
        descripcion_regimen: dataRFC.descripcion_regimen // Incluir la descripción del régimen fiscal
      });
      
      // Mostrar el modal con la información
      setMostrarModalEmpresa(true);

    } catch (error) {
      console.error('Error al buscar empresa:', error);
      setErrorBusqueda(error.message || 'Error al buscar la empresa');
    } finally {
      setBuscandoEmpresa(false);
    }
  };

  return (
<div className="info-personal-container" style={{ marginTop: '0px', marginLeft: '290px' }}>      {modalErrores.length > 0 && (
        <div className="modal-errores">
          <div className="modal-contenido">
            <h3>Campos pendientes:</h3>
            <ul>
              {modalErrores.map((error, index) => (
                <li key={index}>{error.mensaje}</li>
              ))}
            </ul>
            <button onClick={() => setModalErrores([])}>Cerrar</button>
          </div>
        </div>
      )}

      {/* Mostrar mensaje de éxito o error */}
      {mensaje && (
        <div className={`mensaje-banner ${mensaje.tipo}`}>
          {mensaje.texto}
          <button 
            onClick={() => setMensaje(null)} 
            className="cerrar-mensaje"
          >
            ×
          </button>
        </div>
      )}

      <h1 className="titulo">Información Fiscal</h1>

      {/* Buscador de Empresas por RFC */}
      <div className="buscador-card">
        <div className="card-header">
          <h2>Buscar Empresa por RFC</h2>
        </div>
        <div className="buscador-content">
          <div className="buscador-grupo">
            <label htmlFor="rfcBusqueda">RFC de la Empresa:</label>
            <div className="buscador-input-group">
              <input
                type="text"
                id="rfcBusqueda"
                value={rfcBusqueda}
                onChange={(e) => setRfcBusqueda(e.target.value.toUpperCase())}
                placeholder="Ej: ABC123456DE7"
                maxLength="13"
                disabled={buscandoEmpresa}
              />
              <button
                type="button"
                className="btn-buscar"
                onClick={buscarEmpresaPorRFC}
                disabled={buscandoEmpresa || !rfcBusqueda.trim()}
              >
                {buscandoEmpresa ? 'Buscando...' : 'Buscar'}
              </button>
            </div>
            {errorBusqueda && (
              <div className="error-busqueda">{errorBusqueda}</div>
            )}
          </div>
        </div>
      </div>

      {/* Modal de Información de la Empresa */}
      {mostrarModalEmpresa && empresaEncontrada && (
        <div className="modal-overlay">
          <div className="modal-empresa">
            <div className="modal-header">
              <h3>Información de la Empresa</h3>
              <button 
                className="btn-cerrar-modal"
                onClick={() => setMostrarModalEmpresa(false)}
              >
                ×
              </button>
            </div>
            <div className="modal-body">
              <div className="empresa-info-grid">
                <div className="info-item">
                  <label>RFC:</label>
                  <span>{empresaEncontrada.rfc}</span>
                </div>
                <div className="info-item">
                  <label>Nombre Comercial:</label>
                  <span>{empresaEncontrada.nombre_comercial || 'No disponible'}</span>
                </div>
                <div className="info-item">
                  <label>Razón Social:</label>
                  <span>{empresaEncontrada.razon_social || 'No disponible'}</span>
                </div>
                <div className="info-item">
                  <label>Dirección:</label>
                  <span>{empresaEncontrada.direccion1 || 'No disponible'}</span>
                </div>
                <div className="info-item">
                  <label>Colonia:</label>
                  <span>{empresaEncontrada.colonia || 'No disponible'}</span>
                </div>
                <div className="info-item">
                  <label>Código Postal:</label>
                  <span>{empresaEncontrada.cp || 'No disponible'}</span>
                </div>
                <div className="info-item">
                  <label>Ciudad:</label>
                  <span>{empresaEncontrada.ciudad || 'No disponible'}</span>
                </div>
                <div className="info-item">
                  <label>Estado:</label>
                  <span>{empresaEncontrada.estado || 'No disponible'}</span>
                </div>
                <div className="info-item">
                  <label>Método de Pago:</label>
                  <span>{empresaEncontrada.metodo_pago || 'No disponible'}</span>
                </div>
                <div className="info-item">
                  <label>Régimen Fiscal:</label>
                  <span>
                    {empresaEncontrada.c_regimenfiscal && empresaEncontrada.descripcion_regimen 
                      ? `${empresaEncontrada.c_regimenfiscal} - ${empresaEncontrada.descripcion_regimen}`
                      : 'No disponible'
                    }
                  </span>
                </div>
              </div>
            </div>
            <div className="modal-footer">
              <button 
                className="btn-cerrar"
                onClick={() => setMostrarModalEmpresa(false)}
              >
                Cerrar
              </button>
            </div>
          </div>
        </div>
      )}

      <div className="info-card">
        <div className="card-header">
          <h2>Datos Fiscales de la Empresa</h2>
          {!isEditing && isAdmin ? (
            <button 
              className="btn-editar" 
              onClick={() => setIsEditing(true)}
            >
              Editar Datos 
            </button>
          ) : isEditing && isAdmin ? (
            <button 
              className="btn-cancelar" 
              onClick={() => {
                // Al cancelar, intentar recuperar datos o simplemente limpiar el formulario
                if (userData && !isSubmitting) {
                  try {
                    fetchDatosFiscales(userData.id || userData.email || userData.username);
                  } catch {
                    // Resetear formulario
                    setDatosFiscales({
                      claveSat: '',
                      rfcEmisor: '',
                      razonSocial: '',
                      cp: '',
                      direccionFiscal: '',
                      claveArchivoCSD: '',
                      regimenFiscal: ''
                    });
                  }
                }
                setCsdKey(null);
                setCsdCer(null);
                setFieldErrors({});
              }}
            >
              Cancelar
            </button>
          ) : null}
        </div>

        <form onSubmit={handleSubmit}>
          <div className="info-grid">
            <div className="info-group">
              <label htmlFor="rfcEmisor">RFC del Emisor: <span className="required">*</span></label>
              <input
                type="text"
                id="rfcEmisor"
                name="rfcEmisor"
                value={datosFiscales.rfcEmisor}
                onChange={handleChange}
                className={fieldErrors.rfcEmisor ? 'input-error' : ''}
                disabled={!isEditing}
                required
              />
            </div>
            
            <div className="info-group">
              <label htmlFor="razonSocial">Razón Social: <span className="required">*</span></label>
              <input
                type="text"
                id="razonSocial"
                name="razonSocial"
                value={datosFiscales.razonSocial}
                onChange={handleChange}
                className={fieldErrors.razonSocial ? 'input-error' : ''}
                disabled={!isEditing}
                required
              />
            </div>
            
            <div className="info-group">
              <label htmlFor="direccionFiscal">Dirección Fiscal:</label>
              {isEditing ? (
                <input
                  type="text"
                  id="direccionFiscal"
                  name="direccionFiscal"
                  value={datosFiscales.direccionFiscal}
                  onChange={handleChange}
                  className={fieldErrors.direccionFiscal ? 'input-error' : ''}
                  disabled={!isEditing}
                />
              ) : (
                <div className="read-only-field">
                  {datosFiscales.direccionFiscal || '(No registrado)'}
                </div>
              )}
            </div>
            
            <div className="info-group">
              <label htmlFor="cp">CP: <span className="required">*</span></label>
              <input
                type="text"
                id="cp"
                name="cp"
                value={datosFiscales.cp}
                onChange={handleChange}
                className={fieldErrors.cp ? 'input-error' : ''}
                disabled={!isEditing}
                required
              />
            </div>
            
            {/* Mensaje informativo sobre certificados CSD */}
            <div className="info-message" style={{ gridColumn: "1 / span 2", marginTop: "10px", marginBottom: "10px" }}>
              <span className="info-icon">ℹ️</span>
              <span><strong style={{ fontWeight: "800" }}>CERTIFICADOS CSD</strong> son dos archivos tramitados en la página del SAT, no es la FIEL.</span>
            </div>

            {/* Mover los selectores de archivos CSD aquí */}
            <div className="info-group">
              <label htmlFor="csdKey">CSD KEY: <span className="required">*</span></label>
              {isEditing ? (
                <div className="file-input-container">
                  <input
                    type="file"
                    id="csdKey"
                    name="csdKey"
                    accept=".key"
                    onChange={(e) => handleFileChange(e, 'key')}
                    style={{ display: 'none' }}
                  />
                  <button 
                    type="button" 
                    onClick={() => document.getElementById('csdKey').click()}
                    className="file-select-btn file-select-button"
                  >
                    Seleccionar archivo
                  </button>
                  <span className="file-status">
                    {csdKey ? csdKey.name : 'Sin archivos seleccionados'}
                  </span>
                </div>
              ) : (
                <div className="read-only-field">
                  {csdKey ? 'Archivo cargado' : 'Sin archivo'}
                </div>
              )}
            </div>

            <div className="info-group">
              <label htmlFor="csdCer">CSD CER: <span className="required">*</span></label>
              {isEditing ? (
                <div className="file-input-container">
                  <input
                    type="file"
                    id="csdCer"
                    name="csdCer"
                    accept=".cer"
                    onChange={(e) => handleFileChange(e, 'cer')}
                    style={{ display: 'none' }}
                  />
                  <button 
                    type="button" 
                    onClick={() => document.getElementById('csdCer').click()}
                    className="file-select-btn file-select-button"
                  >
                    Seleccionar archivo
                  </button>
                  <span className="file-status">
                    {csdCer ? csdCer.name : 'Sin archivos seleccionados'}
                  </span>
                </div>
              ) : (
                <div className="read-only-field">
                  {csdCer ? 'Archivo cargado' : 'Sin archivo'}
                </div>
              )}
            </div>

            <div className="info-group">
              <label htmlFor="claveArchivoCSD">Clave del Archivo CSD:</label>
              {isEditing ? (
                <input
                  type="password"
                  id="claveArchivoCSD"
                  name="claveArchivoCSD"
                  value={datosFiscales.claveArchivoCSD}
                  onChange={handleChange}
                  className={fieldErrors.claveArchivoCSD ? 'input-error' : ''}
                />
              ) : (
                <div className="read-only-field" style={{ padding: '10px 0' }}>
                  {datosFiscales.claveArchivoCSD ? '••••••••' : '(No registrado)'}
                </div>
              )}
            </div>
            
            <div className="info-group">
              <label htmlFor="regimenFiscal">Régimen Fiscal: <span className="required">*</span></label>
              {isEditing ? (
                <select
                  id="regimenFiscal"
                  name="regimenFiscal"
                  value={datosFiscales.regimenFiscal}
                  onChange={handleChange}
                  className={fieldErrors.regimenFiscal ? 'input-error' : ''}
                  required
                >
                  <option value="">Seleccione una opción</option>
                  <option value="601">601 - General de Ley Personas Morales</option>
                  <option value="603">603 - Personas Morales con Fines no Lucrativos</option>
                  <option value="605">605 - Sueldos y Salarios e Ingresos Asimilados a Salarios</option>
                  <option value="606">606 - Arrendamiento</option>
                  <option value="608">608 - Demás ingresos</option>
                  <option value="612">612 - Personas Físicas con Actividades Empresariales y Profesionales</option>
                  <option value="620">620 - Sociedades Cooperativas de Producción que optan por diferir sus ingresos</option>
                  <option value="621">621 - Incorporación Fiscal</option>
                  <option value="622">622 - Actividades Agrícolas, Ganaderas, Silvícolas y Pesqueras</option>
                  <option value="625">625 - Régimen de las Actividades Empresariales con ingresos a través de Plataformas Tecnológicas</option>
                  <option value="626">626 - Régimen Simplificado de Confianza</option>
                </select>
              ) : (
                <div className="read-only-field" style={{ padding: '10px 0' }}>
                  {datosFiscales.regimenFiscal 
                    ? `${datosFiscales.regimenFiscal} - ${getRegimenFiscalLabel(datosFiscales.regimenFiscal)}` 
                    : '(No registrado)'}
                </div>
              )}
            </div>
          </div>
          
          {isEditing && (
            <div className="form-actions">
              <button 
                type="submit" 
                className="btn-guardar"
                disabled={isSubmitting}
              >
                {isSubmitting ? 'Guardando...' : 'Guardar Datos Fiscales'}
              </button>
            </div>
          )}
        </form>
      </div>
    </div>
  );
}

export default DatosEmpresa;