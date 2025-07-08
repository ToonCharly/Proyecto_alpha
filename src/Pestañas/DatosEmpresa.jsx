import React, { useState, useEffect, useCallback } from 'react';
import '../STYLES/DatosEmpresa.css';

function DatosEmpresa() {
  console.log("🔄 DatosEmpresa componente iniciando...");
  const [isEditing, setIsEditing] = useState(true); 
  const [datosFiscales, setDatosFiscales] = useState({
    rfc: '',
    nombreComercial: '',
    razonSocial: '',
    direccion1: '',
    colonia: '',
    cp: '',
    ciudad: '',
    estado: '',
    metodoPago: '',
    regimenFiscal: '',
    descripcionRegimen: '',
    tipoPago: '',
    claveArchivoCSD: '',
    condicionPago: '',
    serie: ''
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

  // Estado para controlar la visibilidad del buscador de RFC
  const [mostrarBuscadorRFC, setMostrarBuscadorRFC] = useState(true);

  // Estado para controlar el modal de confirmación de restablecimiento
  const [mostrarModalRestablecer, setMostrarModalRestablecer] = useState(false);

  // Obtener datos del usuario actual
  const [userData, setUserData] = useState(null);
  const [isAdmin, setIsAdmin] = useState(false); 
  const [isLoading, setIsLoading] = useState(true); 

  // Función para guardar datos temporalmente en localStorage
  const guardarDatosLocales = useCallback((datos) => {
    try {
      const userId = userData?.id;
      if (userId) {
        localStorage.setItem(`datosFiscales_${userId}`, JSON.stringify(datos));
      }
    } catch (error) {
      console.error('Error al guardar datos locales:', error);
    }
  }, [userData]);

  // Función para obtener datos fiscales existentes - silenciosa para carga inicial
  const fetchDatosFiscales = useCallback(async (identifier) => {
    try {
      setIsSubmitting(true);
      console.log("🔍 Cargando datos fiscales para usuario:", identifier);
      
      const response = await fetch(`http://localhost:8080/api/datos-fiscales?id_usuario=${identifier || ''}`);
      
      if (response.ok) {
        const data = await response.json();
        console.log("📊 Datos fiscales recibidos:", data);
        
        // Si hay datos, cargarlos en el formulario
        if (data && Object.keys(data).length > 0) {
          const nuevosData = {
            rfc: data.rfcEmisor || data.rfc || '',
            razonSocial: data.razonSocial || '',
            direccion1: data.direccionFiscal || data.direccion1 || '',
            cp: data.cp || data.codigo_postal || '',
            claveArchivoCSD: data.claveArchivoCSD || data.clave_csd || '',
            regimenFiscal: data.regimenFiscal || '',
            serie: data.serie_df || data.serie || '',
            // Mantener estos campos si no vienen del servidor (para preservar datos del buscador)
            nombreComercial: data.nombreComercial || '',
            colonia: data.colonia || '',
            ciudad: data.ciudad || '',
            estado: data.estado || '',
            metodoPago: data.metodoPago || '',
            descripcionRegimen: data.descripcionRegimen || '',
            tipoPago: data.tipoPago || '',
            condicionPago: data.condicionPago || ''
          };
          
          setDatosFiscales(nuevosData);
          
          // Guardar en localStorage
          if (userData?.id) {
            try {
              localStorage.setItem(`datosFiscales_${userData.id}`, JSON.stringify(nuevosData));
            } catch (error) {
              console.error('Error al guardar datos locales:', error);
            }
          }
          
          setIsEditing(false); // Solo salimos del modo edición si hay datos
          console.log("✅ Datos fiscales cargados correctamente");
        } else {
          console.log("ℹ️ No hay datos fiscales previos, mantener modo edición");
          setIsEditing(true);
        }
      } else {
        console.log("ℹ️ No se encontraron datos fiscales, mantener modo edición");
        setIsEditing(true);
      }
    } catch (error) {
      console.error("❌ Error al cargar datos fiscales:", error);
      setIsEditing(true);
    } finally {
      setIsSubmitting(false);
    }
  }, [userData]);

  useEffect(() => {
    const storedUserData = sessionStorage.getItem('userData');
    if (storedUserData) {
      try {
        const parsedUserData = JSON.parse(storedUserData);
        console.log("👤 Usuario encontrado en sessionStorage:", parsedUserData);
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
          setIsEditing(false);
        } else {
          // Si es admin, cargar datos fiscales
          console.log("🔄 Verificando datos fiscales locales...");
          
          try {
            const datosGuardados = localStorage.getItem(`datosFiscales_${parsedUserData.id}`);
            if (datosGuardados) {
              const datosLocales = JSON.parse(datosGuardados);
              console.log("📂 Datos fiscales encontrados en localStorage, cargándolos...");
              setDatosFiscales(datosLocales);
              setIsEditing(false);
            } else {
              // Si no hay datos locales, cargar desde el servidor
              console.log("🔄 Cargando datos fiscales desde el servidor...");
              // No llamar fetchDatosFiscales aquí para evitar dependencias circulares
              // En su lugar, marcar que necesitamos cargar desde el servidor
              setIsEditing(true);
            }
          } catch (error) {
            console.error('Error al cargar datos locales:', error);
            setIsEditing(true);
          }
        }
      } catch (err) {
        console.error("❌ Error al procesar datos del usuario:", err);
        setMensaje({
          tipo: 'error',
          texto: 'Error al cargar datos del usuario. Por favor, inicia sesión nuevamente.'
        });
      }
    } else {
      console.log("⚠️ No se encontraron datos de usuario en sessionStorage");
      setMensaje({
        tipo: 'error',
        texto: 'No hay sesión activa. Por favor, inicia sesión.'
      });
    }
    setIsLoading(false); // Terminar la carga inicial
  }, []);

  // Cargar datos desde el servidor si no hay datos locales y es admin
  useEffect(() => {
    if (userData && isAdmin && isEditing && !datosFiscales.rfc) {
      console.log("🔄 Cargando datos fiscales desde el servidor...");
      if (userData.id) {
        fetchDatosFiscales(userData.id);
      } else if (userData.email || userData.username) {
        fetchDatosFiscales(userData.email || userData.username);
      }
    }
  }, [userData, isAdmin, isEditing, datosFiscales.rfc, fetchDatosFiscales]);

  // Controlar la visibilidad del buscador de RFC
  useEffect(() => {
    // Si no hay RFC en los datos fiscales, mostrar el buscador
    // Si hay RFC, ocultar el buscador (datos ya cargados)
    if (!datosFiscales.rfc || datosFiscales.rfc === '') {
      setMostrarBuscadorRFC(true);
    } else {
      setMostrarBuscadorRFC(false);
    }
  }, [datosFiscales.rfc]);

  // Manejar cambios en los campos de texto
  const handleChange = (e) => {
    const { name, value } = e.target;
    
    // Validación específica para el campo serie
    if (name === 'serie') {
      // Máximo 25 caracteres, solo letras sin acentos, números, sin caracteres especiales ni ñ
      const serieRegex = /^[A-Za-z0-9]*$/;
      
      // Si el valor supera 25 caracteres o contiene caracteres no válidos, no actualizar
      if (value.length > 25 || !serieRegex.test(value)) {
        return; // No actualizar el estado si no cumple las validaciones
      }
    }
    
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
    if (!datosFiscales.rfc) {
      errores.rfc = true;
      mensajesError.push('El RFC es obligatorio');
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
    
    // Validación de archivos CSD - solo obligatorios si no hay datos previos guardados
    const hayDatosPrevios = datosFiscales.rfc && datosFiscales.razonSocial;
    if (!hayDatosPrevios) {
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
      formDataToSend.append('rfc', datosFiscales.rfc);
      formDataToSend.append('nombre_comercial', datosFiscales.nombreComercial);
      formDataToSend.append('razon_social', datosFiscales.razonSocial);
      formDataToSend.append('direccion_fiscal', datosFiscales.direccion1); // Cambiar key para compatibilidad
      formDataToSend.append('direccion', datosFiscales.direccion1); // Agregar también como direccion para datos administrativos
      formDataToSend.append('colonia', datosFiscales.colonia);
      formDataToSend.append('codigo_postal', datosFiscales.cp);
      formDataToSend.append('ciudad', datosFiscales.ciudad);
      formDataToSend.append('estado', datosFiscales.estado);
      formDataToSend.append('metodo_pago', datosFiscales.metodoPago);
      formDataToSend.append('clave_csd', datosFiscales.claveArchivoCSD);
      formDataToSend.append('regimen_fiscal', datosFiscales.regimenFiscal);
      formDataToSend.append('descripcion_regimen', datosFiscales.descripcionRegimen);
      formDataToSend.append('tipo_pago', datosFiscales.tipoPago);
      formDataToSend.append('condicion_pago', datosFiscales.condicionPago);
      formDataToSend.append('serie_df', datosFiscales.serie); // Agregar serie de datos fiscales
      
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
        const responseData = await response.json();
        console.log("✅ Datos guardados exitosamente:", responseData);
        
        // Guardar los datos en localStorage para persistencia entre recargas
        guardarDatosLocales(datosFiscales);
        
        // Mantener los datos actuales que el usuario acaba de guardar
        // NO recargar desde el servidor para mantener todos los campos
        setIsEditing(false); // Cambiar a modo visualización
        setMensaje({
          tipo: 'exito',
          texto: 'Datos fiscales guardados correctamente'
        });
        
        // Ocultar el buscador de RFC después de guardar datos fiscales
        setMostrarBuscadorRFC(false);
        
        // NO recargar los datos automáticamente para evitar pérdida de información
        // Los datos ya están correctos en el estado local
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

  // Función para restablecer todos los datos
  const restablecerDatos = async () => {
    if (!isAdmin) {
      setMensaje({
        tipo: 'error',
        texto: 'Solo los administradores pueden restablecer los datos fiscales'
      });
      return;
    }

    // Mostrar el modal de confirmación en lugar de window.confirm
    setMostrarModalRestablecer(true);
  };

  // Función para confirmar el restablecimiento desde el modal
  const confirmarRestablecimiento = async () => {
    setMostrarModalRestablecer(false);

    try {
      setIsSubmitting(true);
      
      // Limpiar datos en el frontend
      setDatosFiscales({
        rfc: '',
        nombreComercial: '',
        razonSocial: '',
        direccion1: '',
        colonia: '',
        cp: '',
        ciudad: '',
        estado: '',
        metodoPago: '',
        regimenFiscal: '',
        descripcionRegimen: '',
        tipoPago: '',
        claveArchivoCSD: '',
        condicionPago: ''
      });
      
      // Limpiar archivos CSD
      setCsdKey(null);
      setCsdCer(null);
      
      // Limpiar datos guardados en localStorage
      if (userData && userData.id) {
        try {
          localStorage.removeItem(`datosFiscales_${userData.id}`);
          console.log("🗑️ Datos fiscales eliminados del localStorage");
        } catch (error) {
          console.error('Error al limpiar localStorage:', error);
        }
      }
      
      // Limpiar campos de búsqueda
      setRfcBusqueda('');
      setEmpresaEncontrada(null);
      setMostrarModalEmpresa(false);
      
      // Activar modo edición
      setIsEditing(true);
      
      // Mostrar el buscador de RFC después de restablecer
      setMostrarBuscadorRFC(true);
      
      // Llamar al backend para eliminar los datos de la base de datos
      if (userData && userData.id) {
        try {
          const response = await fetch('http://localhost:8080/api/restablecer-datos-fiscales', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({
              id_usuario: userData.id
            })
          });
          
          if (response.ok) {
            setMensaje({
              tipo: 'exito',
              texto: 'Datos fiscales restablecidos correctamente. Usa el buscador de RFC para cargar nueva información.'
            });
          } else {
            // Aunque falle el backend, los datos del frontend ya se limpiaron
            setMensaje({
              tipo: 'info',
              texto: 'Datos restablecidos localmente. Usa el buscador de RFC para cargar nueva información.'
            });
          }
        } catch (error) {
          console.error('Error al restablecer en el servidor:', error);
          // Los datos ya se limpiaron en el frontend
          setMensaje({
            tipo: 'info',
            texto: 'Datos restablecidos localmente. Usa el buscador de RFC para cargar nueva información.'
          });
        }
      }
      
    } catch (error) {
      console.error('Error al restablecer datos:', error);
      setMensaje({
        tipo: 'error',
        texto: 'Error al restablecer los datos. Por favor, intenta nuevamente.'
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  // Función para cargar empresa encontrada en el formulario
  const cargarEmpresaEnFormulario = () => {
    if (empresaEncontrada) {
      setDatosFiscales({
        rfc: empresaEncontrada.rfc || '',
        nombreComercial: empresaEncontrada.nombre_comercial || '',
        razonSocial: empresaEncontrada.razon_social || '',
        direccion1: empresaEncontrada.direccion1 || '',
        colonia: empresaEncontrada.colonia || '',
        cp: empresaEncontrada.cp || '',
        ciudad: empresaEncontrada.ciudad || '',
        estado: empresaEncontrada.estado || '',
        metodoPago: empresaEncontrada.metodo_pago || '',
        regimenFiscal: empresaEncontrada.c_regimenfiscal || '',
        descripcionRegimen: empresaEncontrada.descripcion_regimen || '',
        tipoPago: empresaEncontrada.tipo_pago || '',
        condicionPago: empresaEncontrada.condicion_pago || '',
        claveArchivoCSD: datosFiscales.claveArchivoCSD // Mantener la clave CSD existente
      });
      setMostrarModalEmpresa(false);
      setIsEditing(true); // Activar modo edición para que puedan guardarse los datos
      
      // Ocultar el buscador de RFC después de cargar datos
      setMostrarBuscadorRFC(false);
      
      // Mostrar mensaje de éxito
      setMensaje({
        tipo: 'exito',
        texto: 'Datos de la empresa cargados correctamente. Puedes editarlos y guardarlos.'
      });
    }
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
        throw new Error('Datos no encontrados en el sistema');
      }

      const dataRFC = await responseRFC.json();
      console.log("🔍 Datos del RFC obtenidos:", dataRFC); // Debug temporal
      
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
        descripcion_regimen: dataRFC.descripcion_regimen, // Incluir la descripción del régimen fiscal
        idtipopago: dataRFC.idtipopago, // Incluir el ID del tipo de pago
        tipo_pago: dataRFC.tipo_pago, // Incluir el tipo de pago
        idcondicion: dataRFC.idcondicion, // Incluir el ID de la condición de pago
        condicion_pago: dataRFC.condicion_pago // Incluir la condición de pago
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
    <div className="info-personal-container" style={{ marginTop: '0px', marginLeft: '290px' }}>
      {/* Indicador de carga inicial */}
      {isLoading && (
        <div className="loading-container" style={{ textAlign: 'center', padding: '50px' }}>
          <p>Cargando información fiscal...</p>
        </div>
      )}

      {/* Contenido principal solo cuando no está cargando */}
      {!isLoading && (
        <>
          {modalErrores.length > 0 && (
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

      {/* Buscador de Empresas por RFC - solo se muestra si mostrarBuscadorRFC es true */}
      {mostrarBuscadorRFC && (
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
      )}

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
                <div className="info-item">
                  <label>Tipo de Pago:</label>
                  <span>{empresaEncontrada.tipo_pago || 'No disponible'}</span>
                </div>
                <div className="info-item">
                  <label>Condición de Pago:</label>
                  <span>{empresaEncontrada.condicion_pago || 'No disponible'}</span>
                </div>
              </div>
            </div>
            <div className="modal-footer">
              <button 
                className="btn-cargar-datos"
                onClick={cargarEmpresaEnFormulario}
              >
                Cargar Estos Datos
              </button>
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

      {/* Modal de Confirmación de Restablecimiento */}
      {mostrarModalRestablecer && (
        <div className="modal-overlay">
          <div className="modal-restablecer">
            <h2>¿Restablecer Datos Fiscales?</h2>
            <p>¿Estás seguro de que quieres restablecer todos los datos fiscales? Esta acción no se puede deshacer y eliminará toda la información fiscal actual.</p>
            <div className="modal-buttons">
              <button 
                className="btn btn-secondary" 
                onClick={() => setMostrarModalRestablecer(false)}
              >
                Cancelar
              </button>
              <button 
                className="btn btn-danger" 
                onClick={confirmarRestablecimiento}
                disabled={isSubmitting}
              >
                {isSubmitting ? 'Restableciendo...' : 'Restablecer'}
              </button>
            </div>
          </div>
        </div>
      )}

      <div className="info-card">
        <div className="card-header">
          <h2>Datos Fiscales de la Empresa</h2>
          <div className="header-buttons">
            {!isEditing && isAdmin ? (
              <button 
                className="btn-restablecer" 
                onClick={restablecerDatos}
                disabled={isSubmitting}
              >
                Restablecer
              </button>
            ) : isEditing && isAdmin ? (
              <>
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
                          rfc: '',
                          nombreComercial: '',
                          razonSocial: '',
                          direccion1: '',
                          colonia: '',
                          cp: '',
                          ciudad: '',
                          estado: '',
                          metodoPago: '',
                          regimenFiscal: '',
                          descripcionRegimen: '',
                          tipoPago: '',
                          claveArchivoCSD: '',
                          condicionPago: ''
                        });
                      }
                    }
                    setCsdKey(null);
                    setCsdCer(null);
                    setFieldErrors({});
                    setIsEditing(false);
                  }}
                >
                  Cancelar
                </button>
                <button 
                  className="btn-restablecer" 
                  onClick={restablecerDatos}
                  disabled={isSubmitting}
                >
                  Restablecer
                </button>
              </>
            ) : null}
          </div>
        </div>

        <form onSubmit={handleSubmit}>
          <div className="info-grid">
            <div className="info-group">
              <label htmlFor="rfc">RFC: <span className="required">*</span></label>
              {isEditing ? (
                <input
                  type="text"
                  id="rfc"
                  name="rfc"
                  value={datosFiscales.rfc}
                  onChange={handleChange}
                  className={fieldErrors.rfc ? 'input-error' : ''}
                  disabled={!isEditing}
                  placeholder="Ej: ABC123456DE7"
                  required
                />
              ) : (
                <div className="read-only-field">
                  {datosFiscales.rfc || '—'}
                </div>
              )}
            </div>
            
            <div className="info-group">
              <label htmlFor="nombreComercial">Nombre Comercial:</label>
              {isEditing ? (
                <input
                  type="text"
                  id="nombreComercial"
                  name="nombreComercial"
                  value={datosFiscales.nombreComercial}
                  onChange={handleChange}
                  disabled={!isEditing}
                  placeholder="Nombre comercial de la empresa"
                />
              ) : (
                <div className="read-only-field">
                  {datosFiscales.nombreComercial || '—'}
                </div>
              )}
            </div>
            
            <div className="info-group">
              <label htmlFor="razonSocial">Razón Social: <span className="required">*</span></label>
              {isEditing ? (
                <input
                  type="text"
                  id="razonSocial"
                  name="razonSocial"
                  value={datosFiscales.razonSocial}
                  onChange={handleChange}
                  className={fieldErrors.razonSocial ? 'input-error' : ''}
                  disabled={!isEditing}
                  placeholder="Razón social completa"
                  required
                />
              ) : (
                <div className="read-only-field">
                  {datosFiscales.razonSocial || '—'}
                </div>
              )}
            </div>
            
            <div className="info-group">
              <label htmlFor="direccion1">Dirección:</label>
              {isEditing ? (
                <input
                  type="text"
                  id="direccion1"
                  name="direccion1"
                  value={datosFiscales.direccion1}
                  onChange={handleChange}
                  disabled={!isEditing}
                  placeholder="Ingresa la dirección"
                />
              ) : (
                <div className="read-only-field">
                  {datosFiscales.direccion1 || '—'}
                </div>
              )}
            </div>
            
            <div className="info-group">
              <label htmlFor="colonia">Colonia:</label>
              {isEditing ? (
                <input
                  type="text"
                  id="colonia"
                  name="colonia"
                  value={datosFiscales.colonia}
                  onChange={handleChange}
                  disabled={!isEditing}
                  placeholder="Ingresa la colonia"
                />
              ) : (
                <div className="read-only-field">
                  {datosFiscales.colonia || '—'}
                </div>
              )}
            </div>
            
            <div className="info-group">
              <label htmlFor="cp">Código Postal: <span className="required">*</span></label>
              {isEditing ? (
                <input
                  type="text"
                  id="cp"
                  name="cp"
                  value={datosFiscales.cp}
                  onChange={handleChange}
                  className={fieldErrors.cp ? 'input-error' : ''}
                  disabled={!isEditing}
                  placeholder="Ej: 01000"
                  maxLength="5"
                  required
                />
              ) : (
                <div className="read-only-field">
                  {datosFiscales.cp || '—'}
                </div>
              )}
            </div>
            
            <div className="info-group">
              <label htmlFor="ciudad">Ciudad:</label>
              {isEditing ? (
                <input
                  type="text"
                  id="ciudad"
                  name="ciudad"
                  value={datosFiscales.ciudad}
                  onChange={handleChange}
                  disabled={!isEditing}
                  placeholder="Ingresa la ciudad"
                />
              ) : (
                <div className="read-only-field">
                  {datosFiscales.ciudad || '—'}
                </div>
              )}
            </div>
            
            <div className="info-group">
              <label htmlFor="estado">Estado:</label>
              {isEditing ? (
                <input
                  type="text"
                  id="estado"
                  name="estado"
                  value={datosFiscales.estado}
                  onChange={handleChange}
                  disabled={!isEditing}
                  placeholder="Ingresa el estado"
                />
              ) : (
                <div className="read-only-field">
                  {datosFiscales.estado || '—'}
                </div>
              )}
            </div>
            
            <div className="info-group">
              <label htmlFor="metodoPago">Método de Pago:</label>
              {isEditing ? (
                <input
                  type="text"
                  id="metodoPago"
                  name="metodoPago"
                  value={datosFiscales.metodoPago}
                  onChange={handleChange}
                  disabled={!isEditing}
                  placeholder="Ej: Transferencia electrónica"
                />
              ) : (
                <div className="read-only-field">
                  {datosFiscales.metodoPago || '—'}
                </div>
              )}
            </div>

            <div className="info-group">
              <label htmlFor="tipoPago">Tipo de Pago:</label>
              {isEditing ? (
                <input
                  type="text"
                  id="tipoPago"
                  name="tipoPago"
                  value={datosFiscales.tipoPago}
                  onChange={handleChange}
                  disabled={!isEditing}
                  placeholder="Ej: Contado"
                />
              ) : (
                <div className="read-only-field">
                  {datosFiscales.tipoPago || '—'}
                </div>
              )}
            </div>

            <div className="info-group">
              <label htmlFor="condicionPago">Condición de Pago:</label>
              {isEditing ? (
                <input
                  type="text"
                  id="condicionPago"
                  name="condicionPago"
                  value={datosFiscales.condicionPago}
                  onChange={handleChange}
                  disabled={!isEditing}
                  placeholder="Ej: Inmediato"
                />
              ) : (
                <div className="read-only-field">
                  {datosFiscales.condicionPago || '—'}
                </div>
              )}
            </div>
            
            <div className="info-group">
              <label htmlFor="serie">Serie:</label>
              {isEditing ? (
                <div>
                  <input
                    type="text"
                    id="serie"
                    name="serie"
                    value={datosFiscales.serie}
                    onChange={handleChange}
                    disabled={!isEditing}
                    placeholder="Solo letras y números, máx 25 caracteres (ej: A, B, SERIE1)"
                    maxLength="25"
                  />
                  <div style={{ fontSize: '0.8em', color: '#666', marginTop: '4px' }}>
                    Solo letras (sin acentos), números. Máximo 25 caracteres. No se permiten: ñ, acentos, símbolos.
                  </div>
                </div>
              ) : (
                <div className="read-only-field">
                  {datosFiscales.serie || '—'}
                </div>
              )}
            </div>
            
            {/* Mensaje informativo sobre certificados CSD */}
            <div className="info-message" style={{ gridColumn: "1 / span 2", marginTop: "10px", marginBottom: "20px" }}>
              <span className="info-icon">ℹ️</span>
              <span><strong style={{ fontWeight: "800" }}>CERTIFICADOS CSD</strong> son dos archivos tramitados en la página del SAT, no es la FIEL.</span>
            </div>

            {/* Sección de archivos CSD - Primera fila */}
            <div className="csd-files-section" style={{ gridColumn: "1 / span 2", marginBottom: "20px" }}>
              <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "30px" }}>
                {/* CSD KEY */}
                <div className="info-group">
                  <label htmlFor="csdKey">
                    CSD KEY: 
                    {(!datosFiscales.rfc || !datosFiscales.razonSocial) && <span className="required">*</span>}
                  </label>
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
                      {csdKey ? `📄 ${csdKey.name}` : '—'}
                    </div>
                  )}
                </div>

                {/* CSD CER */}
                <div className="info-group">
                  <label htmlFor="csdCer">
                    CSD CER: 
                    {(!datosFiscales.rfc || !datosFiscales.razonSocial) && <span className="required">*</span>}
                  </label>
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
                      {csdCer ? `📄 ${csdCer.name}` : '—'}
                    </div>
                  )}
                </div>
              </div>
            </div>

            {/* Sección de configuración CSD - Segunda fila */}
            <div className="csd-config-section" style={{ gridColumn: "1 / span 2" }}>
              <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "30px" }}>
                {/* Clave del Archivo CSD */}
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
                      placeholder="Clave de los certificados CSD"
                    />
                  ) : (
                    <div className="read-only-field">
                      {datosFiscales.claveArchivoCSD ? '••••••••' : '—'}
                    </div>
                  )}
                </div>
                
                {/* Régimen Fiscal */}
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
                    <div className="read-only-field">
                      {datosFiscales.regimenFiscal && datosFiscales.descripcionRegimen
                        ? `${datosFiscales.regimenFiscal} - ${datosFiscales.descripcionRegimen}`
                        : datosFiscales.regimenFiscal 
                          ? `${datosFiscales.regimenFiscal} - ${getRegimenFiscalLabel(datosFiscales.regimenFiscal)}`
                          : '—'
                      }
                    </div>
                  )}
                </div>
              </div>
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
        </>
      )}
    </div>
  );
}

export default DatosEmpresa;