import React, { useState, useEffect } from 'react';
import { usePreferencias } from '../context/PreferenciasContext';
import '../STYLES/Preferencias.css';

function Preferencias() {
  // Estados para manejo de temas
  const [selectedTheme, setSelectedTheme] = useState('default');
  const [customColor, setCustomColor] = useState('#000000');
  const [mensaje, setMensaje] = useState(null);
  
  // Estado para gestión de logo
  const [logoImage, setLogoImage] = useState(null);
  const [logoPreview, setLogoPreview] = useState('');
  
  // Usar el contexto para obtener los valores necesarios
  const { 
    companyName, 
    updateCompanyName, 
    companyTextColor, 
    updateCompanyTextColor,
    navbarBgColor,
    updateNavbarBgColor
  } = usePreferencias();
  
  // Estados locales para editar
  const [localCompanyName, setLocalCompanyName] = useState(companyName);
  
  // Actualizar estado local cuando cambien los valores del contexto
  useEffect(() => {
    setLocalCompanyName(companyName);
  }, [companyName]);

  // Cargar configuraciones guardadas
  useEffect(() => {
    // Cargar tema
    const savedTheme = localStorage.getItem('sidebarTheme');
    if (savedTheme) {
      try {
        const themeData = JSON.parse(savedTheme);
        setSelectedTheme(themeData.id);
        if (themeData.id === 'custom') {
          setCustomColor(themeData.color);
        }
        
        // Aplicar el tema guardado
        document.documentElement.style.setProperty('--sidebar-color', themeData.color);
      } catch (error) {
        console.error('Error al cargar el tema guardado:', error);
      }
    }
    
    // Cargar logo si existe
    const savedLogo = localStorage.getItem('appLogo');
    if (savedLogo) {
      setLogoPreview(savedLogo);
      document.documentElement.style.setProperty('--app-logo', `url(${savedLogo})`);
    }
    
    // Cargar colores de botones
    const savedActionButtonsColor = localStorage.getItem('actionButtonsColor');
    if (savedActionButtonsColor) {
      setActionButtonsColor(savedActionButtonsColor);
      document.documentElement.style.setProperty('--action-button-color', savedActionButtonsColor);
    }
    
    const savedDeleteButtonsColor = localStorage.getItem('deleteButtonsColor');
    if (savedDeleteButtonsColor) {
      setDeleteButtonsColor(savedDeleteButtonsColor);
      document.documentElement.style.setProperty('--delete-button-color', savedDeleteButtonsColor);
    }
    
    const savedEditButtonsColor = localStorage.getItem('editButtonsColor');
    if (savedEditButtonsColor) {
      setEditButtonsColor(savedEditButtonsColor);
      document.documentElement.style.setProperty('--edit-button-color', savedEditButtonsColor);
    }
    
    const savedFileSelectButtonsColor = localStorage.getItem('fileSelectButtonsColor');
    if (savedFileSelectButtonsColor) {
      setFileSelectButtonsColor(savedFileSelectButtonsColor);
      document.documentElement.style.setProperty('--file-select-button-color', savedFileSelectButtonsColor);
    }
    
    // Cargar tipografía guardada
    const savedFontId = localStorage.getItem('appFontId');
    if (savedFontId) {
      setSelectedFont(savedFontId);
      const fontFamily = localStorage.getItem('appFontFamily');
      if (fontFamily) {
        document.documentElement.style.setProperty('--app-font-family', fontFamily);
      }
    }
    
    const savedHeadingFontId = localStorage.getItem('appHeadingFontId');
    if (savedHeadingFontId) {
      setSelectedHeadingFont(savedHeadingFontId);
      const headingFontFamily = localStorage.getItem('appHeadingFontFamily');
      if (headingFontFamily) {
        document.documentElement.style.setProperty('--app-heading-font-family', headingFontFamily);
      }
    }
  }, []);
  
  // Función para manejar cambios en el nombre (actualiza en tiempo real)
  const handleCompanyNameChange = (e) => {
    const newName = e.target.value;
    setLocalCompanyName(newName);
    updateCompanyName(newName); // Actualiza inmediatamente el contexto
  };
  
  // Función para manejar cambios en el color del texto
  const handleTextColorChange = (e) => {
    updateCompanyTextColor(e.target.value); // Actualiza inmediatamente el contexto
  };
  
  // Añadir el manejador para el color del navbar
  const handleNavbarBgColorChange = (e) => {
    updateNavbarBgColor(e.target.value);
  };

  // Aplicar tema a la aplicación
  const applyTheme = (themeId, color) => {
    const themeColor = themeId === 'custom' ? color : themeOptions.find(t => t.id === themeId)?.color;
    
    if (themeColor) {
      document.documentElement.style.setProperty('--sidebar-color', themeColor);
      
      // Guardar preferencia de tema en localStorage
      localStorage.setItem('sidebarTheme', JSON.stringify({
        id: themeId,
        color: themeId === 'custom' ? color : themeColor
      }));
      
      showSuccessMessage('Color del panel actualizado correctamente');
    }
  };

  // Función para seleccionar un archivo de imagen
  const handleLogoChange = (e) => {
    const file = e.target.files[0];
    
    if (file) {
      if (file.size > 1024 * 1024) { // 1MB límite
        showErrorMessage('La imagen es demasiado grande. Tamaño máximo: 1MB');
        return;
      }
      
      const allowedTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/svg+xml'];
      if (!allowedTypes.includes(file.type)) {
        showErrorMessage('Formato no soportado. Use: JPG, PNG, GIF o SVG');
        return;
      }
      
      setLogoImage(file);
      
      // Crear vista previa
      const reader = new FileReader();
      reader.onloadend = () => {
        setLogoPreview(reader.result);
      };
      reader.readAsDataURL(file);
    }
  };

  // Guardar el logo
  const saveLogo = () => {
    if (logoPreview) {
      localStorage.setItem('appLogo', logoPreview);
      document.documentElement.style.setProperty('--app-logo', `url(${logoPreview})`);
      showSuccessMessage('Logo guardado correctamente');
    }
  };

  // Eliminar el logo
  const removeLogo = () => {
    setLogoImage(null);
    setLogoPreview('');
    localStorage.removeItem('appLogo');
    document.documentElement.style.removeProperty('--app-logo');
    showSuccessMessage('Logo eliminado correctamente');
  };

  // Funciones para mensajes
  const showSuccessMessage = (text) => {
    setMensaje({ tipo: 'exito', texto: text });
    setTimeout(() => {
      setMensaje(null);
    }, 3000);
  };

  const showErrorMessage = (text) => {
    setMensaje({ tipo: 'error', texto: text });
    setTimeout(() => {
      setMensaje(null);
    }, 5000);
  };

  const handleThemeChange = (e) => {
    const newTheme = e.target.value;
    setSelectedTheme(newTheme);
    applyTheme(newTheme, customColor);
  };

  const handleCustomColorChange = (e) => {
    const newColor = e.target.value;
    setCustomColor(newColor);
    if (selectedTheme === 'custom') {
      applyTheme('custom', newColor);
    }
  };

  // Añade estos estados después de los estados existentes
  const [actionButtonsColor, setActionButtonsColor] = useState('#2e7d32'); // Color para guardar/descargar
  const [deleteButtonsColor, setDeleteButtonsColor] = useState('#d32f2f'); // Color para eliminar
  const [editButtonsColor, setEditButtonsColor] = useState('#1976d2'); // Color para botones de editar
  const [fileSelectButtonsColor, setFileSelectButtonsColor] = useState('#455a64'); // Color para seleccionar archivos

  // Función para manejar el cambio de color de los botones de acción
  const handleActionButtonsColorChange = (e) => {
    const newColor = e.target.value;
    setActionButtonsColor(newColor);
    
    // Asegurar que el color se aplica inmediatamente
    document.documentElement.style.setProperty('--action-button-color', newColor);
    
    // Calcular un color más oscuro para hover
    const darkerColor = getDarkerColor(newColor);
    document.documentElement.style.setProperty('--action-button-color-dark', darkerColor);
    
    // Guardar en localStorage
    localStorage.setItem('actionButtonsColor', newColor);
    localStorage.setItem('actionButtonsColorDark', darkerColor);
    
    showSuccessMessage('Color de botones de acción actualizado');
  };

  // Función para manejar el cambio de color de los botones de eliminar
  const handleDeleteButtonsColorChange = (e) => {
    const newColor = e.target.value;
    setDeleteButtonsColor(newColor);
    localStorage.setItem('deleteButtonsColor', newColor);
    document.documentElement.style.setProperty('--delete-button-color', newColor);
    showSuccessMessage('Color de botones de eliminar actualizado');
  };

  // Función para manejar el cambio de color de los botones de editar
  const handleEditButtonsColorChange = (e) => {
    const newColor = e.target.value;
    setEditButtonsColor(newColor);
    localStorage.setItem('editButtonsColor', newColor);
    document.documentElement.style.setProperty('--edit-button-color', newColor);
    showSuccessMessage('Color de botones de editar actualizado');
  };

  // Función para manejar el cambio de color de los botones de seleccionar archivo
  const handleFileSelectButtonsColorChange = (e) => {
    const newColor = e.target.value;
    setFileSelectButtonsColor(newColor);
    
    // Asegurar que el color se aplica inmediatamente
    document.documentElement.style.setProperty('--file-select-button-color', newColor);
    
    // Calcular un color más oscuro para hover
    const darkerColor = getDarkerColor(newColor);
    document.documentElement.style.setProperty('--file-select-button-color-dark', darkerColor);
    
    // Guardar en localStorage
    localStorage.setItem('fileSelectButtonsColor', newColor);
    localStorage.setItem('fileSelectButtonsColorDark', darkerColor);
    
    showSuccessMessage('Color de botones de selección de archivo actualizado');
  };

  // Función auxiliar para calcular un color más oscuro para hover
  function getDarkerColor(hexColor) {
    // Convertir el color hex a RGB
    const r = parseInt(hexColor.slice(1, 3), 16);
    const g = parseInt(hexColor.slice(3, 5), 16);
    const b = parseInt(hexColor.slice(5, 7), 16);
    
    // Oscurecer multiplicando por 0.8
    const darkerR = Math.floor(r * 0.8);
    const darkerG = Math.floor(g * 0.8);
    const darkerB = Math.floor(b * 0.8);
    
    // Convertir de nuevo a hex
    return `#${darkerR.toString(16).padStart(2, '0')}${darkerG.toString(16).padStart(2, '0')}${darkerB.toString(16).padStart(2, '0')}`;
  }

  // Opciones de temas
  const themeOptions = [
    { id: 'default', name: 'Tema Predeterminado', color: '#455a64' },
    { id: 'dark', name: 'Oscuro', color: '#37474f' },
    { id: 'red', name: 'Rojo', color: '#d32f2f' },
    { id: 'green', name: 'Verde', color: '#2e7d32' },
    { id: 'blue', name: 'Azul', color: '#1976d2' },
    { id: 'purple', name: 'Púrpura', color: '#6a1b9a' },
    { id: 'custom', name: 'Personalizado', color: '#000000' }
  ];

  // Opciones de fuentes
  const fontOptions = [
    { id: 'roboto', name: 'Roboto (Predeterminado)', family: "'Roboto', sans-serif" },
    { id: 'openSans', name: 'Open Sans', family: "'Open Sans', sans-serif" },
    { id: 'lato', name: 'Lato', family: "'Lato', sans-serif" },
    { id: 'montserrat', name: 'Montserrat', family: "'Montserrat', sans-serif" },
    { id: 'poppins', name: 'Poppins', family: "'Poppins', sans-serif" },
    { id: 'sourceSansPro', name: 'Source Sans Pro', family: "'Source Sans Pro', sans-serif" },
    { id: 'raleway', name: 'Raleway', family: "'Raleway', sans-serif" }
  ];

  // Opciones de fuentes para títulos (puedes usar las mismas o diferentes)
  const headingFontOptions = [
    { id: 'roboto', name: 'Roboto (Predeterminado)', family: "'Roboto', sans-serif" },
    { id: 'openSans', name: 'Open Sans', family: "'Open Sans', sans-serif" },
    { id: 'lato', name: 'Lato', family: "'Lato', sans-serif" },
    { id: 'montserrat', name: 'Montserrat', family: "'Montserrat', sans-serif" },
    { id: 'poppins', name: 'Poppins', family: "'Poppins', sans-serif" },
    { id: 'playfairDisplay', name: 'Playfair Display', family: "'Playfair Display', serif" },
    { id: 'merriweather', name: 'Merriweather', family: "'Merriweather', serif" }
  ];

  // Estados para manejar las selecciones de fuentes
  const [selectedFont, setSelectedFont] = useState('roboto');
  const [selectedHeadingFont, setSelectedHeadingFont] = useState('roboto');

  // Añade esto al principio del componente o en un archivo de inicialización global
  useEffect(() => {
    // Establecer valores predeterminados para las variables CSS
    if (!document.documentElement.style.getPropertyValue('--action-button-color')) {
      document.documentElement.style.setProperty('--action-button-color', '#2e7d32');
      document.documentElement.style.setProperty('--action-button-color-dark', '#1b5e20');
    }
    
    if (!document.documentElement.style.getPropertyValue('--delete-button-color')) {
      document.documentElement.style.setProperty('--delete-button-color', '#d32f2f');
      document.documentElement.style.setProperty('--delete-button-color-dark', '#b71c1c');
    }
    
    if (!document.documentElement.style.getPropertyValue('--edit-button-color')) {
      document.documentElement.style.setProperty('--edit-button-color', '#1976d2');
      document.documentElement.style.setProperty('--edit-button-color-dark', '#0d47a1');
    }
    
    if (!document.documentElement.style.getPropertyValue('--file-select-button-color')) {
      document.documentElement.style.setProperty('--file-select-button-color', '#455a64');
      document.documentElement.style.setProperty('--file-select-button-color-dark', '#37474f');
    }
    
    // Cargar valores guardados del localStorage
    const savedActionButtonsColor = localStorage.getItem('actionButtonsColor');
    if (savedActionButtonsColor) {
      document.documentElement.style.setProperty('--action-button-color', savedActionButtonsColor);
      setActionButtonsColor(savedActionButtonsColor);
    }
    
    // Repetir para los demás colores
  }, []);
  
  // Añade estas funciones

  // Función para cambiar la fuente principal
  const handleFontChange = (e) => {
    const fontId = e.target.value;
    setSelectedFont(fontId);
    
    const fontFamily = fontOptions.find(f => f.id === fontId)?.family;
    if (fontFamily) {
      document.documentElement.style.setProperty('--app-font-family', fontFamily);
      localStorage.setItem('appFontFamily', fontFamily);
      localStorage.setItem('appFontId', fontId);
      showSuccessMessage('Tipografía principal actualizada');
    }
  };

  // Función para cambiar la fuente de los títulos
  const handleHeadingFontChange = (e) => {
    const fontId = e.target.value;
    setSelectedHeadingFont(fontId);
    
    const fontFamily = headingFontOptions.find(f => f.id === fontId)?.family;
    if (fontFamily) {
      document.documentElement.style.setProperty('--app-heading-font-family', fontFamily);
      localStorage.setItem('appHeadingFontFamily', fontFamily);
      localStorage.setItem('appHeadingFontId', fontId);
      showSuccessMessage('Tipografía de títulos actualizada');
    }
  };

  return (
    <div className="info-personal-container" style={{ marginTop: '60px', marginLeft: '290px' }}>
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
      
      <h1 className="titulo">Preferencias del Sistema</h1>
      
      {/* SECCIÓN DE TEMA */}
      <div className="info-card">
        <div className="card-header">
          <h2>Personalización de Colores</h2>
        </div>
        
        <div className="theme-options">
          <div className="info-group">
            <label htmlFor="theme-select">Color del Panel Administrativo:</label>
            <select 
              id="theme-select" 
              value={selectedTheme} 
              onChange={handleThemeChange}
              className="theme-select"
            >
              {themeOptions.map(theme => (
                <option key={theme.id} value={theme.id}>{theme.name}</option>
              ))}
            </select>
          </div>
          
          {selectedTheme === 'custom' && (
            <div className="info-group">
              <label htmlFor="custom-color">Color Personalizado:</label>
              <input 
                type="color" 
                id="custom-color" 
                value={customColor} 
                onChange={handleCustomColorChange}
                className="color-picker"
              />
            </div>
          )}
          
          <div className="theme-preview-container">
            <h3>Vista Previa</h3>
            <div className="theme-preview-flex">
              <div className="theme-preview-sidebar" 
                   style={{ backgroundColor: selectedTheme === 'custom' 
                           ? customColor 
                           : themeOptions.find(t => t.id === selectedTheme)?.color }}>
                <div className="preview-menu-item">Inicio</div>
                <div className="preview-menu-item">Facturas</div>
                <div className="preview-menu-item active">Configuración</div>
              </div>
              <div className="theme-preview-content">
                <div className="preview-header">Panel Administrativo</div>
                <div className="preview-content">Vista previa del tema</div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      {/* SECCIÓN DE LOGO */}
      <div className="info-card" style={{ marginTop: '30px' }}>
        <div className="card-header">
          <h2>Logo de la Empresa</h2>
        </div>
        
        <div className="logo-options">
          <div className="logo-upload">
            <label htmlFor="logo-upload">Seleccionar Logo:</label>
            <div className="file-upload-container">
              <input
                type="file"
                id="logo-upload"
                accept="image/png, image/jpeg, image/gif, image/svg+xml"
                onChange={handleLogoChange}
                style={{ display: 'none' }}
              />
              <button 
                type="button" 
                onClick={() => document.getElementById('logo-upload').click()}
                className="file-input-button"
              >
                {logoImage ? 'Cambiar imagen' : 'Seleccionar imagen'}
              </button>
              
              {logoImage && (
                <span className="file-name">
                  {logoImage.name}
                </span>
              )}
            </div>
            
            <div className="file-requirements">
              <small>Formatos permitidos: JPG, PNG, GIF, SVG. Tamaño máximo: 1MB</small>
            </div>
            
            <div className="logo-preview-container">
              <h3>Vista Previa del Logo</h3>
              <div className="logo-preview">
                {logoPreview ? (
                  <img src={logoPreview} alt="Logo preview" className="preview-image" />
                ) : (
                  <div className="no-logo">No hay logo seleccionado</div>
                )}
              </div>
            </div>
            
            <div className="logo-actions">
              <button 
                type="button" 
                onClick={saveLogo}
                className="btn-guardar"
                disabled={!logoPreview}
              >
                Guardar Logo
              </button>
              <button 
                type="button" 
                onClick={removeLogo}
                className="btn-eliminar"
                disabled={!logoPreview}
              >
                Eliminar Logo
              </button>
            </div>
          </div>
        </div>
      </div>
      
      {/* SECCIÓN DE INFORMACIÓN DE EMPRESA */}
      <div className="info-card" style={{ marginTop: '30px' }}>
        <div className="card-header">
          <h2>Información de la Empresa</h2>
        </div>
        
        <div className="company-options">
          <div className="info-group">
            <label htmlFor="company-name">Nombre de la Empresa:</label>
            <input
              type="text"
              id="company-name"
              value={localCompanyName}
              onChange={handleCompanyNameChange}
              className="theme-select"
              placeholder="Ingrese el nombre de su empresa"
            />
            <p className="field-help">Este nombre aparecerá en el encabezado del portal administrativo.</p>
          </div>
          
          <div className="info-group" style={{ marginTop: '20px' }}>
            <label htmlFor="text-color">Color del Texto:</label>
            <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
              <input
                type="color"
                id="text-color"
                value={companyTextColor}
                onChange={handleTextColorChange}
                className="color-picker"
              />
              <span style={{ color: companyTextColor, fontWeight: 'bold' }}>
                Vista previa del color
              </span>
            </div>
            <p className="field-help">Este color se aplicará al texto del encabezado.</p>
          </div>
          
          {/* NUEVA OPCIÓN: Color de fondo del navbar */}
          <div className="info-group" style={{ marginTop: '20px' }}>
            <label htmlFor="navbar-bg-color">Color de fondo del panel:</label>
            <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
              <input
                type="color"
                id="navbar-bg-color"
                value={navbarBgColor}
                onChange={handleNavbarBgColorChange}
                className="color-picker"
              />
              <div style={{ 
                backgroundColor: navbarBgColor, 
                padding: '8px 15px', 
                borderRadius: '4px',
                border: '1px solid #ddd'
              }}>
                Vista previa del panel
              </div>
            </div>
            <p className="field-help">Este color se aplicará al fondo del panel de navegación.</p>
          </div>
          
          {/* Vista previa completa */}
          <div className="preview-container" style={{ marginTop: '25px', padding: '15px', backgroundColor: '#f5f5f5', borderRadius: '4px' }}>
            <h3>Vista Previa:</h3>
            <div style={{ 
              padding: '15px',
              backgroundColor: navbarBgColor,
              borderRadius: '4px',
              border: '1px solid #ddd',
              marginTop: '10px'
            }}>
              <span style={{ 
                fontSize: '1.5rem', 
                fontWeight: 'bold', 
                color: companyTextColor 
              }}>
                Portal Administrativo de {localCompanyName}
              </span>
            </div>
          </div>
        </div>
      </div>
      
      {/* NUEVA SECCIÓN: PERSONALIZACIÓN DE BOTONES */}
      <div className="info-card" style={{ marginTop: '30px' }}>
        <div className="card-header">
          <h2>Personalización de Botones</h2>
        </div>
        
        <div className="buttons-options">
          <div className="info-group">
            <label htmlFor="action-buttons-color">Color de botones de acción (Guardar, Descargar):</label>
            <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
              <input
                type="color"
                id="action-buttons-color"
                value={actionButtonsColor}
                onChange={handleActionButtonsColorChange}
                className="color-picker"
              />
              <button style={{ 
                backgroundColor: actionButtonsColor, 
                color: 'white',
                padding: '8px 16px',
                borderRadius: '4px',
                border: 'none',
                cursor: 'pointer'
              }}>
                Vista previa
              </button>
            </div>
            <p className="field-help">Este color se aplicará a botones de guardar y descargar.</p>
          </div>
          
          <div className="info-group" style={{ marginTop: '20px' }}>
            <label htmlFor="delete-buttons-color">Color de botones de eliminar:</label>
            <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
              <input
                type="color"
                id="delete-buttons-color"
                value={deleteButtonsColor}
                onChange={handleDeleteButtonsColorChange}
                className="color-picker"
              />
              <button style={{ 
                backgroundColor: deleteButtonsColor, 
                color: 'white',
                padding: '8px 16px',
                borderRadius: '4px',
                border: 'none',
                cursor: 'pointer'
              }}>
                Vista previa
              </button>
            </div>
            <p className="field-help">Este color se aplicará a botones de eliminar.</p>
          </div>
          
          <div className="info-group" style={{ marginTop: '20px' }}>
            <label htmlFor="edit-buttons-color">Color de botones de editar:</label>
            <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
              <input
                type="color"
                id="edit-buttons-color"
                value={editButtonsColor}
                onChange={handleEditButtonsColorChange}
                className="color-picker"
              />
              <button style={{ 
                backgroundColor: editButtonsColor, 
                color: 'white',
                padding: '8px 16px',
                borderRadius: '4px',
                border: 'none',
                cursor: 'pointer'
              }}>
                Vista previa
              </button>
            </div>
            <p className="field-help">Este color se aplicará a botones de editar información y datos.</p>
          </div>
          
          <div className="info-group" style={{ marginTop: '20px' }}>
            <label htmlFor="file-select-buttons-color">Color de botones de selección de archivo:</label>
            <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
              <input
                type="color"
                id="file-select-buttons-color"
                value={fileSelectButtonsColor}
                onChange={handleFileSelectButtonsColorChange}
                className="color-picker"
              />
              <button style={{ 
                backgroundColor: fileSelectButtonsColor, 
                color: 'white',
                padding: '8px 16px',
                borderRadius: '4px',
                border: 'none',
                cursor: 'pointer'
              }}>
                Vista previa
              </button>
            </div>
            <p className="field-help">Este color se aplicará a botones para seleccionar archivos .key y .cer.</p>
          </div>
        </div>
      </div>
      
      {/* NUEVA SECCIÓN: PERSONALIZACIÓN DE TIPOGRAFÍA */}
      <div className="info-card" style={{ marginTop: '30px' }}>
        <div className="card-header">
          <h2>Personalización de Tipografía</h2>
        </div>
        
        <div className="typography-options">
          <div className="info-group">
            <label htmlFor="font-select">Fuente Principal:</label>
            <select 
              id="font-select" 
              value={selectedFont} 
              onChange={handleFontChange}
              className="theme-select"
            >
              {fontOptions.map(font => (
                <option key={font.id} value={font.id}>{font.name}</option>
              ))}
            </select>            
            <div style={{ 
              marginTop: '15px', 
              padding: '15px', 
              border: '1px solid #ddd', 
              borderRadius: '4px',
              fontFamily: fontOptions.find(f => f.id === selectedFont)?.family
            }}>
              <p style={{ margin: 0 }}>Vista previa de la fuente seleccionada.</p>
              <p style={{ margin: '10px 0 0 0' }}>HOLA A TODOS</p>
              <p style={{ margin: '5px 0 0 0' }}>hola a todos</p>
              <p style={{ margin: '5px 0 0 0' }}>0123456789</p>
            </div>
          </div>
          
          <div className="info-group" style={{ marginTop: '25px' }}>
            <label htmlFor="heading-font-select">Fuente para Títulos:</label>
            <select 
              id="heading-font-select" 
              value={selectedHeadingFont} 
              onChange={handleHeadingFontChange}
              className="theme-select"
            >
              {headingFontOptions.map(font => (
                <option key={font.id} value={font.id}>{font.name}</option>
              ))}
            </select>
            
            <div style={{ 
              marginTop: '15px', 
              padding: '15px', 
              border: '1px solid #ddd', 
              borderRadius: '4px'
            }}>
              <h3 style={{ 
                margin: '0 0 10px 0', 
                fontFamily: headingFontOptions.find(f => f.id === selectedHeadingFont)?.family 
              }}>
                Vista previa del título con la fuente seleccionada
              </h3>
              <p style={{ 
                margin: 0,
                fontFamily: fontOptions.find(f => f.id === selectedFont)?.family 
              }}>
                Este es un texto normal con la fuente principal.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Preferencias;