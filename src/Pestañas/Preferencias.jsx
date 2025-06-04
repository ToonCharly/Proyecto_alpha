import React, { useState, useEffect } from 'react';
import { usePreferencias } from '../context/PreferenciasContext';
import '../STYLES/Preferencias.css';

function Preferencias() {
  // Estado para controlar qué panel se está personalizando
  const [panelActivo, setPanelActivo] = useState('admin'); // 'admin' o 'facturacion'
  
  // Estados para manejo de temas - Panel Admin
  const [selectedTheme, setSelectedTheme] = useState('default');
  const [customColor, setCustomColor] = useState('#000000');
  
  // Estados para manejo de temas - Panel Facturación
  const [selectedThemeFactura, setSelectedThemeFactura] = useState('default');
  const [customColorFactura, setCustomColorFactura] = useState('#000000');
  
  const [mensaje, setMensaje] = useState(null);
  
  // Estado para gestión de logo (compartido)
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
  
  // Estados para colores de botones - Panel Admin
  const [actionButtonsColor, setActionButtonsColor] = useState('#2e7d32');
  const [deleteButtonsColor, setDeleteButtonsColor] = useState('#d32f2f');
  const [editButtonsColor, setEditButtonsColor] = useState('#1976d2');
  const [fileSelectButtonsColor, setFileSelectButtonsColor] = useState('#455a64');
  
  // Estados para colores de botones - Panel Facturación
  const [actionButtonsColorFactura, setActionButtonsColorFactura] = useState('#2e7d32');
  const [deleteButtonsColorFactura, setDeleteButtonsColorFactura] = useState('#d32f2f');
  const [editButtonsColorFactura, setEditButtonsColorFactura] = useState('#1976d2');
  const [fileSelectButtonsColorFactura, setFileSelectButtonsColorFactura] = useState('#455a64');
  
  // Estados para fuentes - Panel Admin
  const [selectedFont, setSelectedFont] = useState('roboto');
  const [selectedHeadingFont, setSelectedHeadingFont] = useState('roboto');
  
  // Estados para fuentes - Panel Facturación
  const [selectedFontFactura, setSelectedFontFactura] = useState('roboto');
  const [selectedHeadingFontFactura, setSelectedHeadingFontFactura] = useState('roboto');
  
  // Actualizar estado local cuando cambien los valores del contexto
  useEffect(() => {
    setLocalCompanyName(companyName);
  }, [companyName]);

  // Cargar configuraciones guardadas
  useEffect(() => {
    // Cargar tema del panel admin
    const savedTheme = localStorage.getItem('sidebarTheme');
    if (savedTheme) {
      try {
        const themeData = JSON.parse(savedTheme);
        setSelectedTheme(themeData.id);
        if (themeData.id === 'custom') {
          setCustomColor(themeData.color);
        }
      } catch (error) {
        console.error('Error al cargar el tema guardado:', error);
      }
    }
    
    // Cargar tema del panel facturación
    const savedThemeFactura = localStorage.getItem('userPanelTheme');
    if (savedThemeFactura) {
      try {
        const themeData = JSON.parse(savedThemeFactura);
        setSelectedThemeFactura(themeData.id);
        if (themeData.id === 'custom') {
          setCustomColorFactura(themeData.color);
        }
      } catch (error) {
        console.error('Error al cargar el tema de facturación:', error);
      }
    }
    
    // Cargar logo compartido
    const savedLogo = localStorage.getItem('appLogo');
    if (savedLogo) {
      setLogoPreview(savedLogo);
    }
    
    // Cargar colores de botones - Admin
    loadButtonColors('admin');
    
    // Cargar colores de botones - Facturación
    loadButtonColors('factura');
    
    // Cargar tipografías - Admin
    loadFonts('admin');
    
    // Cargar tipografías - Facturación
    loadFonts('factura');
    
  }, []);
  
  // Función para cargar colores de botones según el panel
  const loadButtonColors = (panel) => {
    const storagePrefix = panel === 'admin' ? '' : 'factura_';
    const cssPrefix = panel === 'admin' ? 'admin-' : 'user-';
    
    const savedActionButtonsColor = localStorage.getItem(`${storagePrefix}actionButtonsColor`);
    if (savedActionButtonsColor) {
      panel === 'admin' 
        ? setActionButtonsColor(savedActionButtonsColor)
        : setActionButtonsColorFactura(savedActionButtonsColor);
        
      // Aplicar la variable CSS correcta según el panel
      document.documentElement.style.setProperty(`--${cssPrefix}action-button-color`, savedActionButtonsColor);
    }
    
    const savedDeleteButtonsColor = localStorage.getItem(`${storagePrefix}deleteButtonsColor`);
    if (savedDeleteButtonsColor) {
      panel === 'admin'
        ? setDeleteButtonsColor(savedDeleteButtonsColor)
        : setDeleteButtonsColorFactura(savedDeleteButtonsColor);
        
      document.documentElement.style.setProperty(`--${cssPrefix}delete-button-color`, savedDeleteButtonsColor);
    }
    
    const savedEditButtonsColor = localStorage.getItem(`${storagePrefix}editButtonsColor`);
    if (savedEditButtonsColor) {
      panel === 'admin'
        ? setEditButtonsColor(savedEditButtonsColor)
        : setEditButtonsColorFactura(savedEditButtonsColor);
    }
    
    const savedFileSelectButtonsColor = localStorage.getItem(`${storagePrefix}fileSelectButtonsColor`);
    if (savedFileSelectButtonsColor) {
      panel === 'admin'
        ? setFileSelectButtonsColor(savedFileSelectButtonsColor)
        : setFileSelectButtonsColorFactura(savedFileSelectButtonsColor);
    }
  };
  
  // Función para cargar fuentes según el panel
  const loadFonts = (panel) => {
    const storagePrefix = panel === 'admin' ? '' : 'factura_';
    
    const savedFontId = localStorage.getItem(`${storagePrefix}appFontId`);
    if (savedFontId) {
      panel === 'admin'
        ? setSelectedFont(savedFontId)
        : setSelectedFontFactura(savedFontId);
    }
    
    const savedHeadingFontId = localStorage.getItem(`${storagePrefix}appHeadingFontId`);
    if (savedHeadingFontId) {
      panel === 'admin'
        ? setSelectedHeadingFont(savedHeadingFontId)
        : setSelectedHeadingFontFactura(savedHeadingFontId);
    }
  };
  
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

  // Aplicar tema a la aplicación según el panel seleccionado
  const applyTheme = (themeId, color) => {
    const themeColor = themeId === 'custom' ? color : themeOptions.find(t => t.id === themeId)?.color;
    
    if (themeColor) {
      if (panelActivo === 'admin') {
        // Aplicar al panel admin
        document.documentElement.style.setProperty('--sidebar-color', themeColor);
        document.documentElement.style.setProperty('--admin-accent-color', themeColor);
        
        // Guardar preferencia de tema en localStorage
        localStorage.setItem('sidebarTheme', JSON.stringify({
          id: themeId,
          color: themeColor,
          textColor: companyTextColor // Usar variable existente en lugar de adminTheme.textColor
        }));
      } else {
        // Aplicar al panel de facturación
        document.documentElement.style.setProperty('--user-sidebar-color', themeColor);
        document.documentElement.style.setProperty('--user-accent-color', themeColor);
        
        // Guardar preferencia de tema en localStorage
        localStorage.setItem('userPanelTheme', JSON.stringify({
          id: themeId,
          color: themeColor,
          sidebarColor: themeColor,
          textColor: companyTextColor, // Usar variable existente en lugar de facturaTheme.textColor
          navbarColor: navbarBgColor, // Usar variable existente en lugar de facturaTheme.navbarColor
          accentColor: themeColor,
          buttonColor: actionButtonsColorFactura // Usar variable existente en lugar de facturaTheme.buttonColor
        }));
      }
      
      showSuccessMessage(`Color del panel de ${panelActivo === 'admin' ? 'administración' : 'facturación'} actualizado correctamente`);
    }
  };
  
  // Funciones para el selector de panel
  const handlePanelChange = (panel) => {
    setPanelActivo(panel);
  };

  // Función para seleccionar un archivo de imagen (logo compartido)
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

  // Guardar el logo (compartido para ambos paneles)
  const saveLogo = () => {
    if (logoPreview) {
      localStorage.setItem('appLogo', logoPreview);
      document.documentElement.style.setProperty('--app-logo', `url(${logoPreview})`);
      showSuccessMessage('Logo guardado correctamente');
    }
  };

  // Eliminar el logo (compartido para ambos paneles)
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

  // Manejadores de cambio de tema según el panel activo
  const handleThemeChange = (e) => {
    const newTheme = e.target.value;
    
    if (panelActivo === 'admin') {
      setSelectedTheme(newTheme);
      applyTheme(newTheme, customColor);
    } else {
      setSelectedThemeFactura(newTheme);
      applyTheme(newTheme, customColorFactura);
    }
  };

  const handleCustomColorChange = (e) => {
    const newColor = e.target.value;
    
    if (panelActivo === 'admin') {
      setCustomColor(newColor);
      if (selectedTheme === 'custom') {
        applyTheme('custom', newColor);
      }
    } else {
      setCustomColorFactura(newColor);
      if (selectedThemeFactura === 'custom') {
        applyTheme('custom', newColor);
      }
    }
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

  // Opciones de temas (compartidas para ambos paneles)
  const themeOptions = [
    { id: 'default', name: 'Tema Predeterminado', color: '#455a64' },
    { id: 'dark', name: 'Oscuro', color: '#37474f' },
    { id: 'red', name: 'Rojo', color: '#d32f2f' },
    { id: 'green', name: 'Verde', color: '#2e7d32' },
    { id: 'blue', name: 'Azul', color: '#1976d2' },
    { id: 'purple', name: 'Púrpura', color: '#6a1b9a' },
    { id: 'custom', name: 'Personalizado', color: '#000000' }
  ];

  // Opciones de fuentes (compartidas para ambos paneles)
  const fontOptions = [
    { id: 'roboto', name: 'Roboto (Predeterminado)', family: "'Roboto', sans-serif" },
    { id: 'lato', name: 'Lato', family: "'Lato', sans-serif" },
    { id: 'montserrat', name: 'Montserrat', family: "'Montserrat', sans-serif" },
    { id: 'openSans', name: 'Open Sans', family: "'Open Sans', sans-serif" },
    { id: 'poppins', name: 'Poppins', family: "'Poppins', sans-serif" },
    { id: 'raleway', name: 'Raleway', family: "'Raleway', sans-serif" },
    { id: 'sourceSansPro', name: 'Source Sans Pro', family: "'Source Sans Pro', sans-serif" }
  ];

  // Opciones de fuentes para títulos (compartidas para ambos paneles)
  const headingFontOptions = [
    { id: 'roboto', name: 'Roboto (Predeterminado)', family: "'Roboto', sans-serif" },
    { id: 'lato', name: 'Lato', family: "'Lato', sans-serif" },
    { id: 'merriweather', name: 'Merriweather', family: "'Merriweather', serif" },
    { id: 'montserrat', name: 'Montserrat', family: "'Montserrat', sans-serif" },
    { id: 'openSans', name: 'Open Sans', family: "'Open Sans', sans-serif" },
    { id: 'playfairDisplay', name: 'Playfair Display', family: "'Playfair Display', serif" },
    { id: 'poppins', name: 'Poppins', family: "'Poppins', sans-serif" }
  ];

  // Funciones para manejar cambios en los colores de botones según el panel activo
  const handleActionButtonsColorChange = (e) => {
    const newColor = e.target.value;
    const storagePrefix = panelActivo === 'admin' ? '' : 'factura_';
    const cssPrefix = panelActivo === 'admin' ? '' : 'user-';
    
    if (panelActivo === 'admin') {
      setActionButtonsColor(newColor);
    } else {
      setActionButtonsColorFactura(newColor);
    }
    
    // Aplicar inmediatamente
    document.documentElement.style.setProperty(`--${cssPrefix}action-button-color`, newColor);
    
    // Calcular color oscuro para hover
    const darkerColor = getDarkerColor(newColor);
    document.documentElement.style.setProperty(`--${cssPrefix}action-button-color-dark`, darkerColor);
    
    // Guardar en localStorage
    localStorage.setItem(`${storagePrefix}actionButtonsColor`, newColor);
    localStorage.setItem(`${storagePrefix}actionButtonsColorDark`, darkerColor);
    
    showSuccessMessage(`Color de botones de acción del panel de ${panelActivo === 'admin' ? 'administración' : 'facturación'} actualizado`);
  };

  // Función para manejar cambios en el color de botones de eliminar
  const handleDeleteButtonsColorChange = (e) => {
    const newColor = e.target.value;
    const storagePrefix = panelActivo === 'admin' ? '' : 'factura_';
    const cssPrefix = panelActivo === 'admin' ? '' : 'user-';
    
    if (panelActivo === 'admin') {
      setDeleteButtonsColor(newColor);
    } else {
      setDeleteButtonsColorFactura(newColor);
    }
    
    // Aplicar inmediatamente
    document.documentElement.style.setProperty(`--${cssPrefix}delete-button-color`, newColor);
    
    // Guardar en localStorage
    localStorage.setItem(`${storagePrefix}deleteButtonsColor`, newColor);
    
    showSuccessMessage(`Color de botones de eliminar del panel de ${panelActivo === 'admin' ? 'administración' : 'facturación'} actualizado`);
  };

  // Función para manejar cambios en el color de botones de editar
  const handleEditButtonsColorChange = (e) => {
    const newColor = e.target.value;
    const storagePrefix = panelActivo === 'admin' ? '' : 'factura_';
    const cssPrefix = panelActivo === 'admin' ? '' : 'user-';
    
    if (panelActivo === 'admin') {
      setEditButtonsColor(newColor);
    } else {
      setEditButtonsColorFactura(newColor);
    }
    
    // Aplicar inmediatamente
    document.documentElement.style.setProperty(`--${cssPrefix}edit-button-color`, newColor);
    
    // Guardar en localStorage
    localStorage.setItem(`${storagePrefix}editButtonsColor`, newColor);
    
    showSuccessMessage(`Color de botones de editar del panel de ${panelActivo === 'admin' ? 'administración' : 'facturación'} actualizado`);
  };

  // Función para manejar cambios en el color de botones de seleccionar archivo
  const handleFileSelectButtonsColorChange = (e) => {
    const newColor = e.target.value;
    const storagePrefix = panelActivo === 'admin' ? '' : 'factura_';
    const cssPrefix = panelActivo === 'admin' ? '' : 'user-';
    
    if (panelActivo === 'admin') {
      setFileSelectButtonsColor(newColor);
    } else {
      setFileSelectButtonsColorFactura(newColor);
    }
    
    // Aplicar inmediatamente
    document.documentElement.style.setProperty(`--${cssPrefix}file-select-button-color`, newColor);
    
    // Calcular color oscuro para hover
    const darkerColor = getDarkerColor(newColor);
    document.documentElement.style.setProperty(`--${cssPrefix}file-select-button-color-dark`, darkerColor);
    
    // Guardar en localStorage
    localStorage.setItem(`${storagePrefix}fileSelectButtonsColor`, newColor);
    localStorage.setItem(`${storagePrefix}fileSelectButtonsColorDark`, darkerColor);
    
    showSuccessMessage(`Color de botones de selección de archivo del panel de ${panelActivo === 'admin' ? 'administración' : 'facturación'} actualizado`);
  };

  // Funciones para cambiar las fuentes según el panel activo
  const handleFontChange = (e) => {
    const fontId = e.target.value;
    const storagePrefix = panelActivo === 'admin' ? '' : 'factura_';
    const cssPrefix = panelActivo === 'admin' ? 'admin-' : 'user-';
    
    if (panelActivo === 'admin') {
      setSelectedFont(fontId);
    } else {
      setSelectedFontFactura(fontId);
    }
    
    const fontFamily = fontOptions.find(f => f.id === fontId)?.family;
    if (fontFamily) {
      // Ahora actualizará correctamente --admin-app-font-family
      document.documentElement.style.setProperty(`--${cssPrefix}app-font-family`, fontFamily);
      localStorage.setItem(`${storagePrefix}appFontFamily`, fontFamily);
      localStorage.setItem(`${storagePrefix}appFontId`, fontId);
      showSuccessMessage(`Tipografía principal del panel de ${panelActivo === 'admin' ? 'administración' : 'facturación'} actualizada`);
    }
  };

  // Función para cambiar la fuente de los títulos
  const handleHeadingFontChange = (e) => {
    const fontId = e.target.value;
    const storagePrefix = panelActivo === 'admin' ? '' : 'factura_';
    const cssPrefix = panelActivo === 'admin' ? 'admin-' : 'user-';
    
    if (panelActivo === 'admin') {
      setSelectedHeadingFont(fontId);
    } else {
      setSelectedHeadingFontFactura(fontId);
    }
    
    const fontFamily = headingFontOptions.find(f => f.id === fontId)?.family;
    if (fontFamily) {
      // Ahora actualizará correctamente --admin-app-heading-font-family
      document.documentElement.style.setProperty(`--${cssPrefix}app-heading-font-family`, fontFamily);
      localStorage.setItem(`${storagePrefix}appHeadingFontFamily`, fontFamily);
      localStorage.setItem(`${storagePrefix}appHeadingFontId`, fontId);
      showSuccessMessage(`Tipografía de títulos del panel de ${panelActivo === 'admin' ? 'administración' : 'facturación'} actualizada`);
    }
  };

  // Y lo mismo para handleHeadingFontFamilyChange

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
      
      {/* SELECTOR DE PANEL */}
      <div className="panel-selector-container">
        <h2>Selecciona el panel a personalizar:</h2>
        <div className="panel-selector-tabs">
          <div 
            className={`panel-tab ${panelActivo === 'admin' ? 'active' : ''}`}
            onClick={() => handlePanelChange('admin')}
          >
            Panel Administrativo
          </div>
          <div 
            className={`panel-tab ${panelActivo === 'facturacion' ? 'active' : ''}`}
            onClick={() => handlePanelChange('facturacion')}
          >
            Panel de Facturación
          </div>
        </div>
        
        <div className="panel-indicator">
          Estás personalizando: <strong>{panelActivo === 'admin' ? 'PANEL ADMINISTRATIVO' : 'PANEL DE FACTURACIÓN'}</strong>
        </div>
      </div>
      
      {/* SECCIÓN DE TEMA */}
      <div className="info-card">
        <div className="card-header">
          <h2>Personalización de Colores</h2>
        </div>
        
        <div className="theme-options">
          <div className="info-group">
            <label htmlFor="theme-select">Color del Panel {panelActivo === 'admin' ? 'Administrativo' : 'de Facturación'}:</label>
            <select 
              id="theme-select" 
              value={panelActivo === 'admin' ? selectedTheme : selectedThemeFactura} 
              onChange={handleThemeChange}
              className="theme-select"
            >
              {themeOptions.map(theme => (
                <option key={theme.id} value={theme.id}>{theme.name}</option>
              ))}
            </select>
          </div>
          
          {(panelActivo === 'admin' && selectedTheme === 'custom') || 
           (panelActivo === 'facturacion' && selectedThemeFactura === 'custom') ? (
            <div className="info-group">
              <label htmlFor="custom-color">Color Personalizado:</label>
              <input 
                type="color" 
                id="custom-color" 
                value={panelActivo === 'admin' ? customColor : customColorFactura} 
                onChange={handleCustomColorChange}
                className="color-picker"
              />
            </div>
          ) : null}
          
          <div className="theme-preview-container">
            <h3>Vista Previa</h3>
            <div className="theme-preview-flex">
              <div className="theme-preview-sidebar" 
                   style={{ backgroundColor: panelActivo === 'admin' 
                     ? (selectedTheme === 'custom' ? customColor : themeOptions.find(t => t.id === selectedTheme)?.color)
                     : (selectedThemeFactura === 'custom' ? customColorFactura : themeOptions.find(t => t.id === selectedThemeFactura)?.color) 
                   }}>
                <div className="preview-menu-item">Inicio</div>
                <div className="preview-menu-item">Facturas</div>
                <div className="preview-menu-item active">Configuración</div>
              </div>
              <div className="theme-preview-content">
                <div className="preview-header">
                  {panelActivo === 'admin' ? 'Panel Administrativo' : 'Panel de Facturación'}
                </div>
                <div className="preview-content">Vista previa del tema</div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      {/* SECCIÓN DE LOGO (Compartida) */}
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
      
      {/* SECCIÓN DE INFORMACIÓN DE EMPRESA (Compartida) */}
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
            <p className="field-help">Este nombre aparecerá en el encabezado del portal.</p>
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
                Portal {panelActivo === 'admin' ? 'Administrativo' : 'de Facturación'} de {localCompanyName}
              </span>
            </div>
          </div>
        </div>
      </div>
      
      {/* SECCIÓN DE PERSONALIZACIÓN DE BOTONES */}
      <div className="info-card" style={{ marginTop: '30px' }}>
        <div className="card-header">
          <h2>Personalización de Botones - {panelActivo === 'admin' ? 'Panel Administrativo' : 'Panel de Facturación'}</h2>
        </div>
        
        <div className="buttons-options">
          <div className="info-group">
            <label htmlFor="action-buttons-color">Color de botones de acción (Guardar, Descargar):</label>
            <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
              <input
                type="color"
                id="action-buttons-color"
                value={panelActivo === 'admin' ? actionButtonsColor : actionButtonsColorFactura}
                onChange={handleActionButtonsColorChange}
                className="color-picker"
              />
              <button style={{ 
                backgroundColor: panelActivo === 'admin' ? actionButtonsColor : actionButtonsColorFactura, 
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
                value={panelActivo === 'admin' ? deleteButtonsColor : deleteButtonsColorFactura}
                onChange={handleDeleteButtonsColorChange}
                className="color-picker"
              />
              <button style={{ 
                backgroundColor: panelActivo === 'admin' ? deleteButtonsColor : deleteButtonsColorFactura, 
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
                value={panelActivo === 'admin' ? editButtonsColor : editButtonsColorFactura}
                onChange={handleEditButtonsColorChange}
                className="color-picker"
              />
              <button style={{ 
                backgroundColor: panelActivo === 'admin' ? editButtonsColor : editButtonsColorFactura, 
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
                value={panelActivo === 'admin' ? fileSelectButtonsColor : fileSelectButtonsColorFactura}
                onChange={handleFileSelectButtonsColorChange}
                className="color-picker"
              />
              <button style={{ 
                backgroundColor: panelActivo === 'admin' ? fileSelectButtonsColor : fileSelectButtonsColorFactura, 
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
      
      {/* SECCIÓN DE PERSONALIZACIÓN DE TIPOGRAFÍA */}
      <div className="info-card" style={{ marginTop: '30px' }}>
        <div className="card-header">
          <h2>Personalización de Tipografía - {panelActivo === 'admin' ? 'Panel Administrativo' : 'Panel de Facturación'}</h2>
        </div>
        
        <div className="typography-options">
          <div className="info-group">
            <label htmlFor="font-select">Fuente Principal:</label>
            <select 
              id="font-select" 
              value={panelActivo === 'admin' ? selectedFont : selectedFontFactura} 
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
              fontFamily: fontOptions.find(f => f.id === (panelActivo === 'admin' ? selectedFont : selectedFontFactura))?.family
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
              value={panelActivo === 'admin' ? selectedHeadingFont : selectedHeadingFontFactura} 
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
                fontFamily: headingFontOptions.find(f => f.id === (panelActivo === 'admin' ? selectedHeadingFont : selectedHeadingFontFactura))?.family 
              }}>
                Vista previa del título con la fuente seleccionada
              </h3>
              <p style={{ 
                margin: 0,
                fontFamily: fontOptions.find(f => f.id === (panelActivo === 'admin' ? selectedFont : selectedFontFactura))?.family 
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