import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import '../STYLES/DatosEmpresa.css';

function DatosEmpresa() {
  const navigate = useNavigate();
  const [isEditing, setIsEditing] = useState(false); // Por defecto en modo visualización
  const [datosFiscales, setDatosFiscales] = useState({
    claveSat: '',
    rfcEmisor: '',
    razonSocial: '',
    cp: '',
    direccionFiscal: '',
    claveArchivoCSD: '',
    regimenFiscal: ''
  });
  
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

  useEffect(() => {
    // Obtener datos del usuario del localStorage
    const storedUserData = localStorage.getItem('userData');
    if (storedUserData) {
      try {
        const parsedUserData = JSON.parse(storedUserData);
        setUserData(parsedUserData);
        
        // Si hay un usuario, intentar obtener sus datos fiscales
        if (parsedUserData.email || parsedUserData.username) {
          fetchDatosFiscales(parsedUserData.email || parsedUserData.username);
        }
      } catch (error) {
        console.error("Error al procesar datos del usuario:", error);
      }
    }
  }, []);

  // Función para obtener datos fiscales existentes
  const fetchDatosFiscales = async (identifier) => {
    try {
      setIsSubmitting(true);
      setMensaje(null); // Limpiar mensaje anterior
      const response = await fetch(`http://localhost:8080/api/empresa/datos-fiscales?identifier=${encodeURIComponent(identifier)}`);
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
        setIsEditing(false); // Si hay datos, mostrar en modo visualización
      }
    } catch (error) {
      console.error("Error al obtener datos fiscales:", error);
      setModalErrores([...modalErrores, { 
        mensaje: `Error al cargar datos fiscales: ${error.message}` 
      }]);
      setMensaje({
        tipo: 'error',
        texto: `Error al cargar datos fiscales: ${error.message}`
      });
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
    
    // Validación de campos requeridos
    const errores = {};
    if (!datosFiscales.rfcEmisor) errores.rfcEmisor = true;
    if (!datosFiscales.razonSocial) errores.razonSocial = true;
    if (!datosFiscales.cp) errores.cp = true;
    if (!datosFiscales.regimenFiscal) errores.regimenFiscal = true;
    
    // Validación de archivos si es un nuevo registro
    if (!datosFiscales.rfcEmisor) {
      if (!csdKey) errores.csdKey = true;
      if (!csdCer) errores.csdCer = true;
    }
    
    if (Object.keys(errores).length > 0) {
      setFieldErrors(errores);
      setIsSubmitting(false);
      const errorMsg = 'Por favor completa todos los campos requeridos';
      setModalErrores([...modalErrores, { mensaje: errorMsg }]);
      setMensaje({ tipo: 'error', texto: errorMsg });
      return;
    }
    
    try {
      // Crear FormData para enviar los archivos
      const formDataToSend = new FormData();
      
      // Agregar los datos del formulario
      Object.keys(datosFiscales).forEach(key => {
        formDataToSend.append(key, datosFiscales[key]);
      });
      
      // Agregar identificador del usuario
      if (userData && (userData.email || userData.username)) {
        formDataToSend.append('identifier', userData.email || userData.username);
      }
      
      // Agregar los archivos solo si se han seleccionado
      if (csdKey) formDataToSend.append('csdKey', csdKey);
      if (csdCer) formDataToSend.append('csdCer', csdCer);
      
      // Enviar los datos al servidor
      const response = await fetch('http://localhost:8080/api/empresa/datos-fiscales', {
        method: 'POST',
        body: formDataToSend,
      });
      
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Error al guardar los datos');
      }
      
      setIsEditing(false);
      setMensaje({
        tipo: 'exito',
        texto: 'Datos fiscales guardados correctamente'
      });
      
      // Usar navigate para redirigir después de un guardado exitoso
      setTimeout(() => {
        navigate('/dashboard');
      }, 2000);
      
    } catch (error) {
      console.error('Error al guardar datos:', error);
      const errorMsg = `Error: ${error.message}`;
      setModalErrores([...modalErrores, { mensaje: errorMsg }]);
      setMensaje({ tipo: 'error', texto: errorMsg });
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

  return (
    <div className="info-personal-container" style={{ marginTop: '60px', marginLeft: '290px' }}>
      {modalErrores.length > 0 && (
        <div className="modal-errores">
          <div className="modal-contenido">
            <h3>Errores:</h3>
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

      <div className="info-card">
        <div className="card-header">
          <h2>Datos Fiscales de la Empresa</h2>
          {!isEditing ? (
            <button 
              className="btn-editar" 
              onClick={() => setIsEditing(true)}
            >
              Editar Datos 
            </button>
          ) : (
            <button 
              className="btn-cancelar" 
              onClick={() => {
                setIsEditing(false);
                if (userData) {
                  fetchDatosFiscales(userData.email || userData.username);
                }
                setCsdKey(null);
                setCsdCer(null);
                setFieldErrors({});
              }}
            >
              Cancelar
            </button>
          )}
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