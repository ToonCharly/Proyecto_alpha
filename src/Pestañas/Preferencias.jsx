import React, { useState, useEffect } from 'react';
import PreferenciasNavegacion from '../components/PreferenciasNavegacion';
import '../STYLES/Preferencias.css';
import '../STYLES/PreferenciasSubsecciones.css';

function Preferencias() {
  // Estado para controlar qué panel se está personalizando
  const [panelActivo, setPanelActivo] = useState('admin'); // 'admin' o 'facturacion'    // NUEVO: Estado para controlar las subsecciones
  const [subseccionActiva, setSubseccionActiva] = useState('colores'); // 'colores', 'botones', 'tipografia', 'empresa', 'plantillas'
  
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
  // Estados para información de la empresa
  const [localCompanyName, setLocalCompanyName] = useState('Mi Empresa');
  const [companyTextColor, setCompanyTextColor] = useState('#000000');
  const [navbarBgColor, setNavbarBgColor] = useState('#ffffff');
  const [baseFontSize, setBaseFontSize] = useState(16);
  const [headingFontSize, setHeadingFontSize] = useState(24);
  
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
  
  // Cargar configuraciones guardadas
  useEffect(() => {
    // Cargar configuraciones de la empresa
    const savedCompanyName = localStorage.getItem('companyName');
    if (savedCompanyName) {
      setLocalCompanyName(savedCompanyName);
    }
    
    const savedCompanyTextColor = localStorage.getItem('companyTextColor');
    if (savedCompanyTextColor) {
      setCompanyTextColor(savedCompanyTextColor);
    }
    
    const savedNavbarBgColor = localStorage.getItem('navbarBgColor');
    if (savedNavbarBgColor) {
      setNavbarBgColor(savedNavbarBgColor);
    }
    
    const savedBaseFontSize = localStorage.getItem('baseFontSize');
    if (savedBaseFontSize) {
      setBaseFontSize(parseInt(savedBaseFontSize));
    }
    
    const savedHeadingFontSize = localStorage.getItem('headingFontSize');
    if (savedHeadingFontSize) {
      setHeadingFontSize(parseInt(savedHeadingFontSize));
    }
    
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
    const cssPrefix = panel === 'admin' ? 'admin-' : 'user-';
    
    const savedFontId = localStorage.getItem(`${storagePrefix}appFontId`);
    const savedFontFamily = localStorage.getItem(`${storagePrefix}appFontFamily`);
    
    if (savedFontId) {
      panel === 'admin'
        ? setSelectedFont(savedFontId)
        : setSelectedFontFactura(savedFontId);
    }
    
    if (savedFontFamily) {
      // Aplicar variables CSS específicas del panel
      document.documentElement.style.setProperty(`--${cssPrefix}app-font-family`, savedFontFamily);
      // También aplicar la variable global
      document.documentElement.style.setProperty(`--app-font-family`, savedFontFamily);
    }
    
    const savedHeadingFontId = localStorage.getItem(`${storagePrefix}appHeadingFontId`);
    const savedHeadingFontFamily = localStorage.getItem(`${storagePrefix}appHeadingFontFamily`);
    
    if (savedHeadingFontId) {
      panel === 'admin'
        ? setSelectedHeadingFont(savedHeadingFontId)
        : setSelectedHeadingFontFactura(savedHeadingFontId);
    }
    
    if (savedHeadingFontFamily) {
      // Aplicar variables CSS específicas del panel
      document.documentElement.style.setProperty(`--${cssPrefix}app-heading-font-family`, savedHeadingFontFamily);
      // También aplicar la variable global
      document.documentElement.style.setProperty(`--app-heading-font-family`, savedHeadingFontFamily);
    }
  };
  
  // Función para manejar cambios en el nombre (actualiza en tiempo real)
  const handleCompanyNameChange = (e) => {
    const newName = e.target.value;
    setLocalCompanyName(newName);
    localStorage.setItem('companyName', newName);
  };
  
  // Función para manejar cambios en el color del texto
  const handleTextColorChange = (e) => {
    const newColor = e.target.value;
    setCompanyTextColor(newColor);
    localStorage.setItem('companyTextColor', newColor);
    document.documentElement.style.setProperty('--company-text-color', newColor);
  };
  
  // Añadir el manejador para el color del navbar
  const handleNavbarBgColorChange = (e) => {
    const newColor = e.target.value;
    setNavbarBgColor(newColor);
    localStorage.setItem('navbarBgColor', newColor);
    document.documentElement.style.setProperty('--navbar-bg-color', newColor);
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

      // Solo permitir PNG y SVG (puedes agregar más si el backend lo soporta)
      const allowedTypes = ['image/png', 'image/svg+xml'];
      if (!allowedTypes.includes(file.type)) {
        // Mostrar ventana emergente personalizada
        window.alert('Formato de imagen NO permitido. Solo se aceptan: PNG y SVG.');
        // También mostrar el mensaje de error en la UI
        showErrorMessage('Formato de imagen NO permitido. Solo se aceptan: PNG y SVG.');
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
  const saveLogo = async () => {
    if (!logoPreview) return;
    
    try {
      setLoading(true);
      
      // Primero actualizar la UI localmente sin esperar al servidor
      localStorage.setItem('appLogo', logoPreview);
      
      // Actualizar la variable CSS global
      document.documentElement.style.setProperty('--app-logo', `url(${logoPreview})`);
      
      // Emitir un evento personalizado para que otros componentes se actualicen
      const logoEvent = new CustomEvent('logoUpdated', { 
        detail: { logoUrl: logoPreview } 
      });
      window.dispatchEvent(logoEvent);
      
      // Almacenar logo en localStorage es suficiente para actualizarlo localmente
      // El evento logoUpdated ya notifica a otros componentes
      
      // Opcionalmente, intentar guardar en el servidor
      if (logoImage) {
        const formData = new FormData();
        formData.append('logo', logoImage);
        
        const response = await fetch('http://localhost:8080/api/guardar-logo', {
          method: 'POST',
          body: formData
        });
        
        if (!response.ok) {
          const errorText = await response.text();
          console.warn('Logo guardado localmente, pero no en el servidor:', errorText);
          // No lanzamos error para no interrumpir la experiencia del usuario
        }
      }
      
      showSuccessMessage('Logo guardado y aplicado correctamente');
    } catch (error) {
      console.error('Error al guardar logo:', error);
      // Aún así mostramos éxito porque el logo se guardó localmente
      showSuccessMessage('Logo aplicado correctamente (no se pudo guardar en el servidor)');
    } finally {
      setLoading(false);
    }
  };

  // Eliminar el logo (compartido para ambos paneles)
  const removeLogo = async () => {
    try {
      setLoading(true);
      
      // Primero actualizar la UI localmente
      setLogoImage(null);
      setLogoPreview('');
      localStorage.removeItem('appLogo');
      document.documentElement.style.removeProperty('--app-logo');
      
      // Limpiar el input de archivo
      if (document.getElementById('logo-upload')) {
        document.getElementById('logo-upload').value = '';
      }
      
      // Emitir evento para que otros componentes se actualicen
      const logoEvent = new CustomEvent('logoUpdated', { 
        detail: { logoUrl: null } 
      });
      window.dispatchEvent(logoEvent);
      
      // Opcionalmente, comunicar al servidor
      try {
        const response = await fetch('http://localhost:8080/api/eliminar-logo', {
          method: 'DELETE'
        });
        
        if (!response.ok) {
          console.warn('Logo eliminado localmente, pero no en el servidor');
        }
      } catch (serverError) {
        console.error('Error al comunicar con el servidor:', serverError);
      }
      
      showSuccessMessage('Logo eliminado correctamente');
    } catch (error) {
      showErrorMessage(`Error al eliminar logo: ${error.message}`);
    } finally {
      setLoading(false);
    }
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
      // Actualizar variables CSS específicas del panel
      document.documentElement.style.setProperty(`--${cssPrefix}app-font-family`, fontFamily);
      // También actualizar la variable global para compatibilidad
      document.documentElement.style.setProperty(`--app-font-family`, fontFamily);
      
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
      // Actualizar variables CSS específicas del panel
      document.documentElement.style.setProperty(`--${cssPrefix}app-heading-font-family`, fontFamily);
      // También actualizar la variable global para compatibilidad
      document.documentElement.style.setProperty(`--app-heading-font-family`, fontFamily);
      
      localStorage.setItem(`${storagePrefix}appHeadingFontFamily`, fontFamily);
      localStorage.setItem(`${storagePrefix}appHeadingFontId`, fontId);
      showSuccessMessage(`Tipografía de títulos del panel de ${panelActivo === 'admin' ? 'administración' : 'facturación'} actualizada`);
    }
  };

  // Función para cambiar el tamaño de fuente base
  const handleBaseFontSizeChange = (e) => {
    const newSize = parseInt(e.target.value);
    setBaseFontSize(newSize);
    localStorage.setItem('baseFontSize', newSize.toString());
    document.documentElement.style.setProperty('--base-font-size', `${newSize}px`);
    showSuccessMessage(`Tamaño de fuente base actualizado a ${newSize}px`);
  };

  // Función para cambiar el tamaño de fuente de títulos
  const handleHeadingFontSizeChange = (e) => {
    const newSize = parseInt(e.target.value);
    setHeadingFontSize(newSize);
    localStorage.setItem('headingFontSize', newSize.toString());
    document.documentElement.style.setProperty('--heading-font-size', `${newSize}px`);
    showSuccessMessage(`Tamaño de fuente de títulos actualizado a ${newSize}px`);
  };

  // Y lo mismo para handleHeadingFontFamilyChange

  // Estados para plantillas de factura
  const [plantillas, setPlantillas] = useState([]);
  const [plantillaWord, setPlantillaWord] = useState(null);
  const [plantillaActiva, setPlantillaActiva] = useState(null);
  const [loading, setLoading] = useState(false);
  const [descripcionPlantilla, setDescripcionPlantilla] = useState("");

  // Estados para logo de plantillas
  const [logoPlantilla, setLogoPlantilla] = useState(null);
  const [logoPlantillaPreview, setLogoPlantillaPreview] = useState('');
  const [logoPlantillaId, setLogoPlantillaId] = useState(null);

  // Cargar plantillas del usuario
  useEffect(() => {
    const cargarPlantillas = async () => {
      try {
        setLoading(true);
        // Obtener ID de usuario del almacenamiento local o contexto
        const rawUserId = localStorage.getItem('userId') || "1"; // Valor por defecto si no hay usuario
        const idUsuario = limpiarUserId(rawUserId);
        
        const response = await fetch(`http://localhost:8080/api/plantillas/listar?id_usuario=${idUsuario}`);
        
        if (!response.ok) {
          throw new Error("Error al cargar plantillas");
        }
        
        const data = await response.json();
        setPlantillas(data.plantillas || []);
        
        // Identificar plantilla activa
        const activa = data.plantillas?.find(p => p.activa);
        if (activa) {
          setPlantillaActiva(activa);
        }
      } catch (error) {
        console.error("Error:", error);
        showErrorMessage(`Error al cargar plantillas: ${error.message}`);
      } finally {
        setLoading(false);
      }
    };
    
    cargarPlantillas();
  }, []);

  // Manejar selección de archivo
  const handlePlantillaChange = (e) => {
    const file = e.target.files[0];
    
    if (file) {
      if (file.size > 5 * 1024 * 1024) { // 5MB límite
        showErrorMessage('La plantilla es demasiado grande. Tamaño máximo: 5MB');
        return;
      }
      
      const allowedTypes = [
        'application/msword', 
        'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
      ];
      
      if (!allowedTypes.includes(file.type) && 
          !file.name.endsWith('.doc') && 
          !file.name.endsWith('.docx')) {
        showErrorMessage('Formato no soportado. Use archivos Word (.doc o .docx)');
        return;
      }
      
      setPlantillaWord(file);
    }
  };

  // Guardar plantilla
  const savePlantilla = async () => {
    if (!plantillaWord) return;
    
    try {
      setLoading(true);
      
      // Obtener ID de usuario
      const rawUserId = localStorage.getItem('userId') || "1";
      const idUsuario = limpiarUserId(rawUserId);
      
      // Crear FormData
      const formData = new FormData();
      formData.append('plantilla', plantillaWord);
      formData.append('descripcion', descripcionPlantilla);
      formData.append('activar', 'true'); // Activar automáticamente la nueva plantilla
      
      // Enviar al servidor
      const response = await fetch(`http://localhost:8080/api/plantillas/subir?id_usuario=${idUsuario}`, {
        method: 'POST',
        body: formData
      });
      
      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText);
      }
      
// Opción 1: Ignorar explícitamente el valor
await response.json();
      
      // Recargar lista de plantillas
      const responseList = await fetch(`http://localhost:8080/api/plantillas/listar?id_usuario=${idUsuario}`);
      if (responseList.ok) {
        const listData = await responseList.json();
        setPlantillas(listData.plantillas || []);
        
        // Actualizar plantilla activa
        const activa = listData.plantillas?.find(p => p.activa);
        if (activa) {
          setPlantillaActiva(activa);
        }
      }
      
      // Limpiar formulario
      setPlantillaWord(null);
      setDescripcionPlantilla("");
      if (document.getElementById('plantilla-upload')) {
        document.getElementById('plantilla-upload').value = '';
      }
      
      showSuccessMessage('Plantilla guardada y activada correctamente');
    } catch (error) {
      showErrorMessage(`Error al guardar plantilla: ${error.message}`);
    } finally {
      setLoading(false);
    }
  };

  // Activar plantilla
  const activarPlantilla = async (idPlantilla) => {
    try {
      setLoading(true);
      
      // Obtener ID de usuario
      const rawUserId = localStorage.getItem('userId') || "1";
      const idUsuario = limpiarUserId(rawUserId);
      
      const response = await fetch('http://localhost:8080/api/plantillas/activar', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          id_usuario: parseInt(idUsuario),
          id_plantilla: idPlantilla
        })
      });
      
      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText);
      }
      
      // Actualizar estado local
      setPlantillas(plantillas.map(p => ({
        ...p,
        activa: p.id === idPlantilla
      })));
      
      // Actualizar plantilla activa
      const nuevaActiva = plantillas.find(p => p.id === idPlantilla);
      if (nuevaActiva) {
        setPlantillaActiva({...nuevaActiva, activa: true});
      }
      
      showSuccessMessage('Plantilla activada correctamente');
    } catch (error) {
      showErrorMessage(`Error al activar plantilla: ${error.message}`);
    } finally {
      setLoading(false);
    }
  };

  // Eliminar plantilla
  const eliminarPlantilla = async (idPlantilla) => {
    if (!confirm("¿Estás seguro de eliminar esta plantilla?")) {
      return;
    }
    
    try {
      setLoading(true);
      
      // Obtener ID de usuario
      const rawUserId = localStorage.getItem('userId') || "1";
      const idUsuario = limpiarUserId(rawUserId);
      
      const response = await fetch('http://localhost:8080/api/plantillas/eliminar', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          id_usuario: parseInt(idUsuario),
          id_plantilla: idPlantilla
        })
      });
      
      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText);
      }
      
      // Actualizar estado local
      const nuevasPlantillas = plantillas.filter(p => p.id !== idPlantilla);
      setPlantillas(nuevasPlantillas);
      
      // Si la plantilla eliminada era la activa, actualizar
      if (plantillaActiva && plantillaActiva.id === idPlantilla) {
        setPlantillaActiva(null);
      }
      
      showSuccessMessage('Plantilla eliminada correctamente');
    } catch (error) {
      showErrorMessage(`Error al eliminar plantilla: ${error.message}`);
    } finally {
      setLoading(false);
    }
  };

  // Función para descargar plantilla de ejemplo
const descargarPlantillaEjemplo = () => {
  try {
    setLoading(true);
    
    fetch('http://localhost:8080/api/plantillas/ejemplo')      .then(response => {
        if (!response.ok) {
          throw new Error(`Error: ${response.status} ${response.statusText}`);
        }
        
        // Logging para debug
        console.log('Response headers:', {
          contentType: response.headers.get('content-type'),
          contentLength: response.headers.get('content-length'),
          contentDisposition: response.headers.get('content-disposition')
        });
        
        return response.blob();
      })
      .then(blob => {
        // Verificar tamaño del blob
        console.log('Tamaño del archivo descargado:', blob.size, 'bytes');
        console.log('Tipo del blob:', blob.type);
        
        if (blob.size < 5000) {
          console.error('Archivo muy pequeño, posiblemente corrupto');
          throw new Error('El archivo descargado parece estar corrupto (muy pequeño)');
        }
        
        // Crear enlace de descarga
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.style.display = 'none';
        a.href = url;
        a.download = 'plantilla_ejemplo_factura.docx';
        
        document.body.appendChild(a);
        a.click();
        
        // Limpiar después de un breve delay
        setTimeout(() => {
          window.URL.revokeObjectURL(url);
          document.body.removeChild(a);
        }, 100);
        
        showSuccessMessage('Plantilla Word descargada correctamente. Ábrela con Microsoft Word o una aplicación compatible.');
      })
      .catch(error => {
        console.error('Error:', error);
        showErrorMessage(`Error al descargar: ${error.message}`);
      })
      .finally(() => {
        setLoading(false);
      });
  } catch (error) {
    console.error('Error general:', error);
    showErrorMessage('Error al procesar la solicitud');
    setLoading(false);
  }
};

// Función alternativa de descarga usando window.open
const descargarPlantillaEjemploDirecto = () => {
  try {
    setLoading(true);
    
    // Método alternativo: abrir en nueva ventana para descarga directa
    const url = 'http://localhost:8080/api/plantillas/ejemplo';
    
    // Crear un enlace temporal y hacer clic
    const link = document.createElement('a');
    link.href = url;
    link.download = 'plantilla_ejemplo_factura.docx';
    link.target = '_blank';
    link.rel = 'noopener noreferrer';
    
    // Agregar al DOM temporalmente
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    
    showSuccessMessage('Descarga iniciada. Verifica tu carpeta de descargas.');
    
  } catch (error) {
    console.error('Error en descarga directa:', error);
    showErrorMessage('Error al iniciar la descarga');
  } finally {
    setLoading(false);
  }
};

  // Función para cargar el logo de plantilla desde la base de datos y forzar recarga (anti-cache)
  const cargarLogoPlantilla = async () => {
    try {
      let rawUserId = localStorage.getItem('userId') || '1';
      const idUsuario = limpiarUserId(rawUserId);
      const idUsuarioNum = parseInt(idUsuario) || 1;
      // Agregar timestamp para evitar cache
      const url = `http://localhost:8080/api/logos/obtener-activo-json?id_usuario=${idUsuarioNum}&t=${Date.now()}`;
      const response = await fetch(url);
      if (response.ok) {
        const data = await response.json();
        if (data.exists && data.imagen_base64) {
          // NO agregar timestamp a data URL, solo usar el base64 puro
          const logoDataUrl = `data:${data.tipo};base64,${data.imagen_base64}`;
          setLogoPlantillaPreview(logoDataUrl);
          if (data.id) {
            setLogoPlantillaId(data.id);
          }
          localStorage.removeItem('logoPlantillaPreview');
        } else {
          setLogoPlantillaPreview('');
          setLogoPlantillaId(null);
          localStorage.removeItem('logoPlantillaPreview');
        }
      } else {
        setLogoPlantillaPreview('');
        setLogoPlantillaId(null);
        localStorage.removeItem('logoPlantillaPreview');
      }
    } catch (error) {
      setLogoPlantillaPreview('');
      setLogoPlantillaId(null);
      localStorage.removeItem('logoPlantillaPreview');
      console.error('Error al cargar logo desde la base de datos:', error);
    }
  };

  // Funciones para manejar el logo de plantillas
  const handleLogoPlantillaChange = async (e) => {
    const file = e.target.files[0];
    if (file) {
      if (file.size > 2 * 1024 * 1024) { // 2MB límite
        showErrorMessage('La imagen del logo es demasiado grande. Tamaño máximo: 2MB');
        return;
      }
      // Solo permitir PNG y SVG (puedes agregar más si el backend lo soporta)
      const allowedTypes = ['image/png', 'image/svg+xml'];
      if (!allowedTypes.includes(file.type)) {
        window.alert('Formato de imagen NO permitido. Solo se aceptan: PNG y SVG.');
        showErrorMessage('Formato de imagen NO permitido. Solo se aceptan: PNG y SVG.');
        return;
      }
      setLogoPlantilla(file);
      // Crear vista previa
      const reader = new FileReader();
      reader.onloadend = async () => {
        setLogoPlantillaPreview(reader.result);
        // Guardar únicamente en la base de datos
        try {
          const rawUserId = localStorage.getItem('userId') || '1';
          const idUsuario = limpiarUserId(rawUserId);
          // Convertir archivo a base64
          const base64 = reader.result.split(',')[1];
          const logoData = {
            id_usuario: parseInt(idUsuario),
            nombre_logo: file.name,
            tipo_mime: file.type,
            imagen_base64: base64
          };
          const response = await fetch(`http://localhost:8080/api/logos/subir`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify(logoData)
          });
          if (response.ok) {
            const result = await response.json();
            if (result.data && result.data.id) {
              setLogoPlantillaId(result.data.id);
            }
            showSuccessMessage('Logo de plantilla guardado correctamente en la base de datos');
            // Recargar logo desde la base de datos para forzar actualización y evitar cache
            await cargarLogoPlantilla();
          } else {
            const errorData = await response.json();
            showErrorMessage(`Error al guardar logo: ${errorData.error || 'Error desconocido'}`);
          }
        } catch (error) {
          console.error('Error al guardar logo en servidor:', error);
          showErrorMessage('Error al guardar logo en la base de datos');
        }
      };
      reader.readAsDataURL(file);
    }
  };

  // Eliminar logo de plantillas
  const eliminarLogoPlantilla = async () => {
    try {
      const rawUserId = localStorage.getItem('userId') || '1';
      const idUsuario = limpiarUserId(rawUserId);
      if (!logoPlantillaId) {
        showErrorMessage('No hay logo para eliminar');
        return;
      }
      const response = await fetch(`http://localhost:8080/api/logos/eliminar`, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ 
          id_logo: logoPlantillaId,
          id_usuario: parseInt(idUsuario) 
        })
      });
      if (response.ok) {
        setLogoPlantilla(null);
        setLogoPlantillaPreview('');
        setLogoPlantillaId(null);
        // Limpiar el input de archivo
        if (document.getElementById('logo-plantilla-upload')) {
          document.getElementById('logo-plantilla-upload').value = '';
        }
        // Recargar logo desde la base de datos para forzar actualización y evitar cache
        await cargarLogoPlantilla();
        showSuccessMessage('Logo de plantilla eliminado correctamente de la base de datos');
      } else {
        const errorData = await response.json();
        showErrorMessage(`Error al eliminar logo: ${errorData.error || 'Error desconocido'}`);
      }
    } catch (error) {
      console.error('Error al eliminar logo:', error);
      showErrorMessage(`Error al eliminar logo: ${error.message}`);
    }
  };

  // Cargar logo de plantilla guardado al inicializar
  useEffect(() => {
    cargarLogoPlantilla();
  }, []);

  // Eliminar logo de plantillas
  // ...existing code... (eliminarLogoPlantilla duplicada, ya está definida más abajo)

  // Cargar logo de plantilla guardado al inicializar
  useEffect(() => {
    const cargarLogoPlantilla = async () => {
      // Cargar únicamente desde la base de datos
      try {
        let rawUserId = localStorage.getItem('userId') || '1';
        const idUsuario = limpiarUserId(rawUserId);
        
        // Asegurar que sea un número válido
        const idUsuarioNum = parseInt(idUsuario) || 1;
        
        const response = await fetch(`http://localhost:8080/api/logos/obtener-activo-json?id_usuario=${idUsuarioNum}`);
        
        if (response.ok) {
          const data = await response.json();
          if (data.exists && data.imagen_base64) {
            const logoDataUrl = `data:${data.tipo};base64,${data.imagen_base64}`;
            setLogoPlantillaPreview(logoDataUrl);
            // Guardar el ID del logo para poder eliminarlo después
            if (data.id) {
              setLogoPlantillaId(data.id);
            }
          }
        }
      } catch (error) {
        console.error('Error al cargar logo desde la base de datos:', error);
      }
    };
    
    cargarLogoPlantilla();
  }, []);

  // Función auxiliar para limpiar el userId
  const limpiarUserId = (userId) => {
    if (!userId || userId === 'default') return '1';
    
    // Limpiar el formato "default:1" a solo "1"
    if (userId.includes(':')) {
      return userId.split(':')[1];
    }
    
    return userId;
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
          </div>        </div>
          <div className="panel-indicator">
          Estás personalizando: <strong>{panelActivo === 'admin' ? 'PANEL ADMINISTRATIVO' : 'PANEL DE FACTURACIÓN'}</strong>
        </div>
      </div>
        {/* NAVEGACIÓN DE SUBSECCIONES MEJORADA */}
      <PreferenciasNavegacion 
        subseccionActiva={subseccionActiva}
        setSubseccionActiva={setSubseccionActiva}
      />      {/* SECCIÓN DE COLORES Y TEMAS */}
      {subseccionActiva === 'colores' && (
        <div className="info-card">
          <div className="card-header">
            <h2>🎨 Personalización de Colores y Temas</h2>
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
            </div>          </div>
        </div>
      </div>      )}
        {/* SECCIÓN DE BOTONES */}
      {subseccionActiva === 'botones' && (
        <div className="info-card">
          <div className="card-header">
            <h2>🔘 Personalización de Botones - {panelActivo === 'admin' ? 'Panel Administrativo' : 'Panel de Facturación'}</h2>
          </div>
          
          <div className="buttons-options">
            <div className="button-customization-grid">
              {/* Botones de Acción */}
              <div className="button-color-group">
                <div className="button-color-header">
                  <div className="button-color-icon action-icon">💾</div>
                  <h3 className="button-color-title">Botones de Acción</h3>
                </div>
                <div className="color-preview-container">
                  <div className="color-input-wrapper">
                    <div className="color-picker-button">
                      <input
                        type="color"
                        id="action-buttons-color"
                        value={panelActivo === 'admin' ? actionButtonsColor : actionButtonsColorFactura}
                        onChange={handleActionButtonsColorChange}
                      />
                      <div 
                        className="color-display" 
                        style={{ backgroundColor: panelActivo === 'admin' ? actionButtonsColor : actionButtonsColorFactura }}
                      ></div>
                    </div>
                  </div>
                  <button 
                    className="preview-button"
                    style={{ 
                      backgroundColor: panelActivo === 'admin' ? actionButtonsColor : actionButtonsColorFactura
                    }}
                  >
                    Guardar
                  </button>
                </div>
                <p className="button-description">
                  Color para botones de guardar, descargar y otras acciones principales
                </p>
              </div>

              {/* Botones de Eliminar */}
              <div className="button-color-group">
                <div className="button-color-header">
                  <div className="button-color-icon delete-icon">🗑️</div>
                  <h3 className="button-color-title">Botones de Eliminar</h3>
                </div>
                <div className="color-preview-container">
                  <div className="color-input-wrapper">
                    <div className="color-picker-button">
                      <input
                        type="color"
                        id="delete-buttons-color"
                        value={panelActivo === 'admin' ? deleteButtonsColor : deleteButtonsColorFactura}
                        onChange={handleDeleteButtonsColorChange}
                      />
                      <div 
                        className="color-display" 
                        style={{ backgroundColor: panelActivo === 'admin' ? deleteButtonsColor : deleteButtonsColorFactura }}
                      ></div>
                    </div>
                  </div>
                  <button 
                    className="preview-button"
                    style={{ 
                      backgroundColor: panelActivo === 'admin' ? deleteButtonsColor : deleteButtonsColorFactura
                    }}
                  >
                    Eliminar
                  </button>
                </div>
                <p className="button-description">
                  Color para botones de eliminar registros y datos
                </p>
              </div>

              {/* Botones de Editar */}
              <div className="button-color-group">
                <div className="button-color-header">
                  <div className="button-color-icon edit-icon">✏️</div>
                  <h3 className="button-color-title">Botones de Editar</h3>
                </div>
                <div className="color-preview-container">
                  <div className="color-input-wrapper">
                    <div className="color-picker-button">
                      <input
                        type="color"
                        id="edit-buttons-color"
                        value={panelActivo === 'admin' ? editButtonsColor : editButtonsColorFactura}
                        onChange={handleEditButtonsColorChange}
                      />
                      <div 
                        className="color-display" 
                        style={{ backgroundColor: panelActivo === 'admin' ? editButtonsColor : editButtonsColorFactura }}
                      ></div>
                    </div>
                  </div>
                  <button 
                    className="preview-button"
                    style={{ 
                      backgroundColor: panelActivo === 'admin' ? editButtonsColor : editButtonsColorFactura
                    }}
                  >
                    Editar
                  </button>
                </div>
                <p className="button-description">
                  Color para botones de editar información y modificar datos
                </p>
              </div>

              {/* Botones de Selección de Archivo */}
              <div className="button-color-group">
                <div className="button-color-header">
                  <div className="button-color-icon file-icon">📁</div>
                  <h3 className="button-color-title">Botones de Archivo</h3>
                </div>
                <div className="color-preview-container">
                  <div className="color-input-wrapper">
                    <div className="color-picker-button">
                      <input
                        type="color"
                        id="file-select-buttons-color"
                        value={panelActivo === 'admin' ? fileSelectButtonsColor : fileSelectButtonsColorFactura}
                        onChange={handleFileSelectButtonsColorChange}
                      />
                      <div 
                        className="color-display" 
                        style={{ backgroundColor: panelActivo === 'admin' ? fileSelectButtonsColor : fileSelectButtonsColorFactura }}
                      ></div>
                    </div>
                  </div>
                  <button 
                    className="preview-button"
                    style={{ 
                      backgroundColor: panelActivo === 'admin' ? fileSelectButtonsColor : fileSelectButtonsColorFactura
                    }}
                  >
                    Seleccionar
                  </button>
                </div>
                <p className="button-description">
                  Color para botones de selección de archivos .key y .cer
                </p>
              </div>
            </div>
          </div>
        </div>
      )}
        {/* SECCIÓN DE TIPOGRAFÍA */}
      {subseccionActiva === 'tipografia' && (
        <div className="info-card">
          <div className="card-header">
            <h2>📝 Personalización de Tipografía - {panelActivo === 'admin' ? 'Panel Administrativo' : 'Panel de Facturación'}</h2>
          </div>
          
          <div className="typography-options">
            {/* Selectores de Fuente */}
            <div className="typography-control-group">
              <div className="typography-selector">
                <label htmlFor="font-select">Fuente Principal</label>
                <select 
                  id="font-select" 
                  value={panelActivo === 'admin' ? selectedFont : selectedFontFactura} 
                  onChange={handleFontChange}
                  className="typography-select"
                >
                  {fontOptions.map(font => (
                    <option key={font.id} value={font.id}>{font.name}</option>
                  ))}
                </select>
                <div className="typography-preview">
                  <p className="typography-preview-title">Vista previa de fuente principal</p>
                  <div 
                    className="typography-preview-samples"
                    style={{ 
                      fontFamily: fontOptions.find(f => f.id === (panelActivo === 'admin' ? selectedFont : selectedFontFactura))?.family,
                      fontSize: `${baseFontSize}px`
                    }}
                  >
                    <p style={{ fontWeight: '500' }}>Texto con énfasis</p>
                    <p style={{ textTransform: 'uppercase', letterSpacing: '1px' }}>TEXTO EN MAYÚSCULAS</p>
                    <p>Texto normal de párrafo</p>
                    <p style={{ fontFamily: 'monospace' }}>0123456789</p>
                  </div>
                </div>
              </div>

              <div className="typography-selector">
                <label htmlFor="heading-font-select">Fuente para Títulos</label>
                <select 
                  id="heading-font-select" 
                  value={panelActivo === 'admin' ? selectedHeadingFont : selectedHeadingFontFactura} 
                  onChange={handleHeadingFontChange}
                  className="typography-select"
                >
                  {headingFontOptions.map(font => (
                    <option key={font.id} value={font.id}>{font.name}</option>
                  ))}
                </select>
                <div className="typography-preview">
                  <p className="typography-preview-title">Vista previa de títulos</p>
                  <div className="typography-preview-samples">
                    <h3 style={{ 
                      margin: '0 0 12px 0', 
                      fontFamily: headingFontOptions.find(f => f.id === (panelActivo === 'admin' ? selectedHeadingFont : selectedHeadingFontFactura))?.family,
                      fontSize: `${headingFontSize}px`,
                      fontWeight: '600'
                    }}>
                      Título Principal
                    </h3>
                    <p style={{ 
                      margin: '0',
                      fontFamily: fontOptions.find(f => f.id === (panelActivo === 'admin' ? selectedFont : selectedFontFactura))?.family,
                      fontSize: `${baseFontSize}px`,
                      color: '#666'
                    }}>
                      Texto de contenido normal.
                    </p>
                  </div>
                </div>
              </div>
            </div>

            <div className="typography-section-divider"></div>
            {/* Controles de Tamaño de Fuente */}
            <div className="typography-control-group">
              <div className="typography-size-control">
                <label htmlFor="base-font-size">Tamaño de Fuente Base</label>
                <div className="typography-size-slider">
                  <input
                    type="range"
                    id="base-font-size"
                    min="12"
                    max="24"
                    value={baseFontSize}
                    onChange={handleBaseFontSizeChange}
                    className="typography-slider"
                  />
                  <span className="typography-size-value">
                    {baseFontSize}px
                  </span>
                </div>
                <div className="typography-preview">
                  <p className="typography-preview-title">Vista previa del tamaño base</p>
                  <div 
                    className="typography-preview-samples"
                    style={{ fontSize: `${baseFontSize}px` }}
                  >
                    <p style={{ margin: '0 0 8px 0', fontWeight: '500' }}>Texto con énfasis</p>
                    <p style={{ margin: '0', color: '#666' }}>Este es un ejemplo de texto normal con el tamaño seleccionado.</p>
                  </div>
                </div>
              </div>
              
              <div className="typography-size-control">
                <label htmlFor="heading-font-size">Tamaño de Fuente para Títulos</label>
                <div className="typography-size-slider">
                  <input
                    type="range"
                    id="heading-font-size"
                    min="18"
                    max="36"
                    value={headingFontSize}
                    onChange={handleHeadingFontSizeChange}
                    className="typography-slider"
                  />
                  <span className="typography-size-value">
                    {headingFontSize}px
                  </span>
                </div>
                <div className="typography-preview">
                  <p className="typography-preview-title">Vista previa del título</p>
                  <div className="typography-preview-samples">
                    <h3 style={{ 
                      margin: '0 0 12px 0',
                      fontSize: `${headingFontSize}px`,
                      fontWeight: '600'
                    }}>
                      Título de Ejemplo
                    </h3>
                    <p style={{ 
                      margin: '0',
                      fontSize: `${baseFontSize}px`,
                      color: '#666'
                    }}>
                      Texto normal para comparación.
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
      
      {/* SECCIÓN DE LOGO Y EMPRESA (SIMPLIFICADA) */}
      {subseccionActiva === 'empresa' && (
        <>
          {/* LOGO */}
          <div className="info-card" style={{ marginTop: '30px' }}>
            <div className="card-header">
              <h2>🏢 Logo de la Empresa</h2>
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
          <h2>🏢 Información de la Empresa</h2>
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
          </div>        </div>
      </div>
        </>
      )}
      
      {/* SECCIÓN DE PLANTILLAS */}
      {subseccionActiva === 'plantillas' && (
      <div className="info-card plantilla-modern" style={{ marginTop: '30px' }}>
        <div className="card-header">
          <h2 style={{display:'flex',alignItems:'center',gap:'10px'}}>
            <span style={{fontSize:'2.2rem',color:'#2b579a'}}>📄</span> Plantillas de Factura
          </h2>
        </div>
        <div className="plantilla-options">
          {/* Info de ayuda y ejemplo */}
          <div className="plantilla-ayuda-modern" style={{
            background:'#f0f8ff',
            border:'1.5px solid #c2e0ff',
            borderRadius:'10px',
            padding:'22px 24px',
            marginBottom:'28px',
            display:'flex',
            alignItems:'flex-start',
            gap:'18px',
            boxShadow:'0 2px 8px rgba(33,150,243,0.06)'}}>
            <div style={{fontSize:'2.1rem',marginRight:'8px',color:'#1976d2',flexShrink:0}}></div>
            <div style={{flex:1}}>
              <div style={{fontWeight:'bold',fontSize:'1.13rem',marginBottom:'6px'}}>¿No sabes cómo estructurar tu plantilla?</div>
              <div style={{fontSize:'1rem',color:'#444',marginBottom:'10px'}}>Descarga nuestra plantilla de ejemplo que muestra la estructura correcta y todas las variables disponibles para crear tus propias facturas personalizadas.</div>
              <div style={{display:'flex',alignItems:'center',gap:'14px'}}>
                <div style={{
                  width:'44px',height:'44px',background:'#2b579a',borderRadius:'6px',display:'flex',alignItems:'center',justifyContent:'center',color:'white',fontWeight:'bold',fontSize:'1.7rem',boxShadow:'0 2px 8px #90caf9',flexShrink:0
                }}>W</div>
                <div style={{flex:1}}>
                  <div style={{fontWeight:'bold',fontSize:'1.05rem',marginBottom:'2px'}}>plantilla_ejemplo_factura.docx</div>
                  <div style={{fontSize:'0.93em',color:'#666'}}>Documento Word con variables de ejemplo</div>
                </div>
                <button 
                  onClick={() => {
                    descargarPlantillaEjemplo();
                    setTimeout(() => {
                      if (!document.querySelector('a[download="plantilla_ejemplo_factura.docx"]')) {
                        descargarPlantillaEjemploDirecto();
                      }
                    }, 2000);
                  }}
                  style={{
                    backgroundColor: actionButtonsColorFactura || '#0078d4',
                    color: 'white',
                    border: 'none',
                    padding: '10px 18px',
                    borderRadius: '6px',
                    cursor: 'pointer',
                    fontWeight:'bold',
                    fontSize:'1.08rem',
                    display:'flex',alignItems:'center',gap:'7px',
                    boxShadow:'0 2px 8px #90caf9'
                  }}
                >
                  <span style={{fontSize:'1.3rem'}}>⬇️</span>
                  Descargar Ejemplo
                </button>
              </div>
            </div>
          </div>

          {/* 1. Subir/cambiar plantilla Word */}
          <div className="plantilla-upload-modern" style={{
            background:'#f9f9fb',
            border:'1.5px solid #e0e0e0',
            borderRadius:'14px',
            padding:'38px 48px',
            marginBottom:'32px',
            boxShadow:'0 2px 14px rgba(0,0,0,0.09)',
            maxWidth:'950px',
            minWidth:'520px',
            minHeight:'220px',
            width:'100%',
            marginRight:'auto',
            marginLeft:'auto',
            marginTop:'0',
            display:'flex',
            flexDirection:'column',
            justifyContent:'center',
          }}>
            <h3 style={{margin:'0 0 10px 0',fontSize:'16px',color:'#333',display:'flex',alignItems:'center',gap:'7px'}}>
              <span style={{fontSize:'1.2rem'}}>📄</span> Subir Plantilla Word
            </h3>
            <p style={{margin:'0 0 10px 0',fontSize:'0.93em',color:'#666'}}>Selecciona tu archivo Word personalizado para facturas.</p>
            <div style={{display:'flex',alignItems:'center',gap:'10px',marginBottom:'10px'}}>
              <input
                type="file"
                id="plantilla-upload"
                accept=".doc,.docx,application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document"
                onChange={handlePlantillaChange}
                style={{ display: 'none' }}
              />
              <button 
                type="button" 
                onClick={() => document.getElementById('plantilla-upload').click()}
                className="file-input-button"
                style={{ 
                  backgroundColor: fileSelectButtonsColorFactura,
                  fontWeight:'bold',fontSize:'0.98em',padding:'7px 16px',borderRadius:'6px',border:'none',cursor:'pointer'
                }}
              >
                {plantillaWord ? 'Cambiar plantilla' : 'Seleccionar plantilla'}
              </button>
              {plantillaWord && (
                <span className="file-name" style={{ marginLeft: '10px',fontSize:'0.97em',color:'#1976d2',fontWeight:'bold' }}>
                  {plantillaWord.name}
                </span>
              )}
            </div>
            <div className="file-requirements" style={{marginBottom:'8px'}}>
              <small style={{color:'#888'}}>Formatos permitidos: DOC, DOCX. Máx: 5MB</small>
            </div>
            <div className="info-group" style={{marginBottom:'8px'}}>
              <div style={{display:'flex',flexDirection:'column',alignItems:'flex-start',textAlign:'left',width:'100%'}}>
                <span style={{fontWeight:'bold',fontSize:'0.97em',marginBottom:'0'}}>
                  Descripción <span style={{fontWeight:400,fontSize:'0.97em',color:'#888',marginLeft:'7px'}}>(opcional)</span>
                </span>
              </div>
              <div style={{display:'flex',justifyContent:'flex-start',width:'100%'}}>
                <input 
                  type="text"
                  id="descripcion-plantilla"
                  value={descripcionPlantilla}
                  onChange={(e) => setDescripcionPlantilla(e.target.value)}
                  className="theme-select"
                  placeholder="Ej: Plantilla para facturas de servicios"
                  style={{
                    marginTop: '4px',
                    fontSize: '0.97em',
                    width: '100%',
                    minWidth: '320px',
                    maxWidth: '700px',
                    padding: '10px 14px',
                    borderRadius: '7px',
                    border: '1.5px solid #b0b0b0',
                    background: '#fff',
                    boxSizing: 'border-box',
                    fontWeight: 500,
                    textAlign: 'left'
                  }}
                />
              </div>
            </div>
            <div className="plantilla-actions" style={{ 
              display: 'flex', 
              gap: '10px',
              marginTop: '10px' 
            }}>
              {loading && (
                <div className="loading-indicator" style={{ 
                  marginRight: '10px',
                  display: 'flex',
                  alignItems: 'center' 
                }}>
                  <div style={{
                    width: '18px',
                    height: '18px',
                    border: '3px solid #f3f3f3',
                    borderTop: '3px solid #3498db',
                    borderRadius: '50%',
                    animation: 'spin 1s linear infinite',
                    marginRight: '6px'
                  }}></div>
                  <span style={{fontSize:'0.97em'}}>Procesando...</span>
                </div>
              )}
              <button 
                type="button" 
                onClick={savePlantilla}
                disabled={!plantillaWord || loading}
                style={{ 
                  backgroundColor: actionButtonsColorFactura,
                  color: 'white',
                  padding: '8px 16px',
                  borderRadius: '4px',
                  border: 'none',
                  cursor: !plantillaWord || loading ? 'not-allowed' : 'pointer',
                  opacity: !plantillaWord || loading ? 0.7 : 1,
                  fontWeight:'bold',fontSize:'1em'
                }}
              >
                Guardar y Activar Plantilla
              </button>
            </div>
          </div>

          {/* 2. Lista de plantillas */}
          {plantillas.length > 0 && (
            <div className="plantillas-lista-modern" style={{
              display: 'flex',
              flexWrap: 'wrap',
              gap: '18px',
              marginBottom: '32px',
              marginTop: '18px',
              justifyContent: 'center', // Centra las cards horizontalmente
              alignItems: 'flex-start', // Opcional: alinea arriba si hay varias filas
              width: '100%'
            }}>
              {plantillas.map(plantilla => (
                <div key={plantilla.id} className="plantilla-card" style={{
                  flex: '1 1 320px',
                  minWidth: '320px',
                  maxWidth: '420px',
                  background: plantilla.activa ? '#e3f2fd' : '#fff',
                  border: plantilla.activa ? '2px solid #1976d2' : '1.5px solid #e0e0e0',
                  borderRadius: '12px',
                  padding: '18px 20px',
                  boxShadow: plantilla.activa ? '0 2px 8px #90caf9' : '0 1px 4px rgba(0,0,0,0.04)',
                  display: 'flex',
                  flexDirection: 'column',
                  gap: '8px',
                  position: 'relative',
                  margin: '0 auto' // Centra la card si hay espacio extra
                }}>
                  <div style={{display:'flex',alignItems:'center',gap:'12px'}}>
                    <span style={{fontSize:'1.7rem',color:'#2b579a'}}></span>
                    <div style={{fontWeight:'bold',fontSize:'1.1rem',color:plantilla.activa?'#1976d2':'#333'}}>{plantilla.nombre}</div>
                    {plantilla.activa && <span style={{color:'#2e7d32',fontWeight:'bold',marginLeft:'8px'}}>• Activa</span>}
                  </div>
                  {plantilla.descripcion && <div style={{fontSize:'0.98em',color:'#555'}}>{plantilla.descripcion}</div>}
                  <div style={{fontSize:'0.92em',color:'#888'}}>Fecha: {plantilla.fecha_creacion}</div>
                  <div style={{display:'flex',gap:'10px',marginTop:'8px'}}>
                    {!plantilla.activa && (
                      <button className="action-button" style={{padding:'7px 16px',borderRadius:'6px',fontWeight:'bold',fontSize:'0.98em',border:'none',cursor:'pointer'}} onClick={()=>activarPlantilla(plantilla.id)}>Activar</button>
                    )}
                    <button className="delete-button" style={{padding:'7px 16px',borderRadius:'6px',fontWeight:'bold',fontSize:'0.98em',border:'none',cursor:'pointer'}} onClick={()=>eliminarPlantilla(plantilla.id)}>Eliminar</button>
                  </div>
                </div>
              ))}
            </div>
          )}

          {/* Card: Logo para Plantillas */}
          <div className="plantilla-upload-modern" style={{
            background:'#f9f9fb',
            border:'1.5px solid #e0e0e0',
            borderRadius:'14px',
            padding:'38px 48px',
            marginBottom:'32px',
            boxShadow:'0 2px 14px rgba(0,0,0,0.09)',
            maxWidth:'950px',
            minWidth:'520px',
            minHeight:'220px',
            width:'100%',
            marginRight:'auto',
            marginLeft:'auto',
            marginTop:'0',
            display:'flex',
            flexDirection:'column',
            justifyContent:'center',
          }}>
            <h3 style={{margin:'0 0 10px 0',fontSize:'16px',color:'#333',display:'flex',alignItems:'center',gap:'7px'}}>
              <span style={{fontSize:'1.2rem'}}>📷</span> Logo para Plantillas
            </h3>
            <p style={{margin:'0 0 10px 0',fontSize:'0.93em',color:'#666'}}>Este logo aparecerá en todas las facturas generadas con plantillas de Word.</p>
            <div style={{display:'flex',alignItems:'center',gap:'10px',marginBottom:'10px'}}>
              <input
                type="file"
                id="logo-plantilla-upload"
                accept="image/png, image/jpeg, image/gif, image/svg+xml"
                onChange={handleLogoPlantillaChange}
                style={{ display: 'none' }}
              />
              <button 
                type="button" 
                onClick={() => document.getElementById('logo-plantilla-upload').click()}
                className="file-input-button"
                style={{ 
                  backgroundColor: fileSelectButtonsColorFactura,
                  fontWeight:'bold',fontSize:'0.98em',padding:'7px 16px',borderRadius:'6px',border:'none',cursor:'pointer'
                }}
              >
                {logoPlantilla ? 'Cambiar logo' : 'Seleccionar logo'}
              </button>
              {logoPlantilla && (
                <span className="file-name" style={{ marginLeft: '10px',fontSize:'0.97em',color:'#1976d2',fontWeight:'bold' }}>
                  {logoPlantilla.name}
                </span>
              )}
            </div>
            <div className="file-requirements" style={{marginBottom:'8px'}}>
              <small style={{color:'#888'}}>Formatos permitidos: JPG, PNG, GIF, SVG. Máx: 2MB. Recomendado: 300x100px</small>
            </div>
            {logoPlantillaPreview && (
              <div className="logo-preview-container" style={{ marginTop: '10px', textAlign: 'center', display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
                {logoPlantillaPreview.startsWith('data:image') ? (
                  <div className="logo-preview" style={{
                    border: '1px solid #ddd',
                    borderRadius: '10px',
                    padding: '18px',
                    backgroundColor: '#f9f9f9',
                    display: 'inline-block',
                    maxWidth: '480px',
                    minWidth: '320px',
                    minHeight: '160px',
                    height: 'auto',
                    textAlign: 'center',
                  }}>
                    <img
                      src={logoPlantillaPreview}
                      alt="Logo preview"
                      style={{
                        maxWidth: '100%',
                        maxHeight: '120px',
                        minHeight: '60px',
                        objectFit: 'contain',
                        display: 'block',
                        margin: '0 auto',
                        background: '#fff',
                        borderRadius: '8px',
                        boxShadow: '0 1px 8px rgba(0,0,0,0.10)'
                      }}
                    />
                  </div>
                ) : (
                  <div style={{ color: '#888', fontSize: '0.98em', padding: '12px' }}>No se pudo mostrar la vista previa del logo.</div>
                )}
                <div style={{ width: '100%', display: 'flex', justifyContent: 'center' }}>
                  <button
                    type="button"
                    onClick={eliminarLogoPlantilla}
                    style={{
                      marginTop: '18px',
                      backgroundColor: deleteButtonsColorFactura,
                      color: 'white',
                      border: 'none',
                      padding: '8px 22px',
                      borderRadius: '6px',
                      cursor: 'pointer',
                      fontSize: '15px',
                      fontWeight: 'bold',
                      boxShadow: '0 1px 6px rgba(0,0,0,0.07)'
                    }}
                  >
                    Eliminar Logo
                  </button>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
      )}
    </div>
  );
}

export default Preferencias;
