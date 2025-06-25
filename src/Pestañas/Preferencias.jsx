import React, { useState, useEffect } from 'react';
import { usePreferencias } from '../context/PreferenciasContext';
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
    // Usar el contexto para obtener los valores necesarios
  const { 
    companyName, 
    updateCompanyName, 
    companyTextColor, 
    updateCompanyTextColor,
    navbarBgColor,
    updateNavbarBgColor,
    baseFontSize,
    updateBaseFontSize,
    headingFontSize,
    updateHeadingFontSize
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
    updateBaseFontSize(newSize);
    showSuccessMessage(`Tamaño de fuente base actualizado a ${newSize}px`);
  };

  // Función para cambiar el tamaño de fuente de títulos
  const handleHeadingFontSizeChange = (e) => {
    const newSize = parseInt(e.target.value);
    updateHeadingFontSize(newSize);
    showSuccessMessage(`Tamaño de fuente de títulos actualizado a ${newSize}px`);
  };

  // Y lo mismo para handleHeadingFontFamilyChange

  // Estados para plantillas de factura
  const [plantillas, setPlantillas] = useState([]);
  const [plantillaWord, setPlantillaWord] = useState(null);
  const [plantillaActiva, setPlantillaActiva] = useState(null);
  const [loading, setLoading] = useState(false);
  const [descripcionPlantilla, setDescripcionPlantilla] = useState("");

  // Cargar plantillas del usuario
  useEffect(() => {
    const cargarPlantillas = async () => {
      try {
        setLoading(true);
        // Obtener ID de usuario del almacenamiento local o contexto
        const idUsuario = localStorage.getItem('userId') || "1"; // Valor por defecto si no hay usuario
        
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
      const idUsuario = localStorage.getItem('userId') || "1";
      
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
      const idUsuario = localStorage.getItem('userId') || "1";
      
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
      const idUsuario = localStorage.getItem('userId') || "1";
      
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
      )}
        {/* SECCIÓN DE TIPOGRAFÍA */}
      {subseccionActiva === 'tipografia' && (
        <div className="info-card">
          <div className="card-header">
            <h2>📝 Personalización de Tipografía - {panelActivo === 'admin' ? 'Panel Administrativo' : 'Panel de Facturación'}</h2>
          </div>
          
          <div className="typography-options">            <div className="info-group">
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
                padding: '20px', 
                border: '1px solid #ddd', 
                borderRadius: '8px',
                minHeight: '120px',
                width: '100%',
                maxWidth: '600px',
                backgroundColor: '#fafafa',
                fontFamily: fontOptions.find(f => f.id === (panelActivo === 'admin' ? selectedFont : selectedFontFactura))?.family,
                fontSize: `${baseFontSize}px`,
                display: 'flex',
                flexDirection: 'column',
                justifyContent: 'center'
              }}>
                <p style={{ margin: '0 0 8px 0', fontWeight: '500' }}>Vista previa de la fuente seleccionada</p>
                <p style={{ margin: '8px 0', textTransform: 'uppercase', letterSpacing: '1px' }}>HOLA A TODOS</p>
                <p style={{ margin: '8px 0' }}>hola a todos</p>
                <p style={{ margin: '8px 0 0 0', fontFamily: 'monospace' }}>0123456789</p>
              </div>
            </div>            <div className="info-group" style={{ marginTop: '25px' }}>
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
                padding: '20px', 
                border: '1px solid #ddd', 
                borderRadius: '8px',
                minHeight: '120px',
                width: '100%',
                maxWidth: '600px',
                backgroundColor: '#fafafa',
                display: 'flex',
                flexDirection: 'column',
                justifyContent: 'center'
              }}>
                <h3 style={{ 
                  margin: '0 0 12px 0', 
                  fontFamily: headingFontOptions.find(f => f.id === (panelActivo === 'admin' ? selectedHeadingFont : selectedHeadingFontFactura))?.family,
                  fontSize: `${headingFontSize}px`,
                  fontWeight: '600'
                }}>
                  Vista previa del título
                </h3>
                <p style={{ 
                  margin: '0',
                  fontFamily: fontOptions.find(f => f.id === (panelActivo === 'admin' ? selectedFont : selectedFontFactura))?.family,
                  fontSize: `${baseFontSize}px`,
                  color: '#666'
                }}>
                  Este es un texto normal con la fuente principal.
                </p>
              </div>
            </div>
              {/* CONTROLES DE TAMAÑO DE FUENTE */}
            <div className="info-group" style={{ marginTop: '25px' }}>
              <label htmlFor="base-font-size">Tamaño de Fuente Base:</label>
              <div style={{ display: 'flex', alignItems: 'center', gap: '15px', marginTop: '10px' }}>
                <input
                  type="range"
                  id="base-font-size"
                  min="12"
                  max="24"
                  value={baseFontSize}
                  onChange={handleBaseFontSizeChange}
                  style={{ flex: 1 }}
                />
                <span style={{ 
                  minWidth: '50px', 
                  textAlign: 'center',
                  padding: '5px 10px',
                  backgroundColor: '#f5f5f5',
                  borderRadius: '4px',
                  fontSize: '14px'
                }}>
                  {baseFontSize}px
                </span>
              </div>
              <div style={{ 
                marginTop: '15px', 
                padding: '20px', 
                border: '1px solid #ddd', 
                borderRadius: '8px',
                minHeight: '120px',
                width: '100%',
                maxWidth: '600px',
                backgroundColor: '#fafafa',
                fontSize: `${baseFontSize}px`,
                display: 'flex',
                flexDirection: 'column',
                justifyContent: 'center'
              }}>
                <p style={{ margin: '0 0 8px 0', fontWeight: '500' }}>Vista previa del tamaño de fuente base</p>
                <p style={{ margin: '0', color: '#666' }}>Este es un ejemplo de texto normal con el tamaño seleccionado.</p>
              </div>
            </div>
            
            <div className="info-group" style={{ marginTop: '25px' }}>
              <label htmlFor="heading-font-size">Tamaño de Fuente para Títulos:</label>
              <div style={{ display: 'flex', alignItems: 'center', gap: '15px', marginTop: '10px' }}>
                <input
                  type="range"
                  id="heading-font-size"
                  min="18"
                  max="36"
                  value={headingFontSize}
                  onChange={handleHeadingFontSizeChange}
                  style={{ flex: 1 }}
                />
                <span style={{ 
                  minWidth: '50px', 
                  textAlign: 'center',
                  padding: '5px 10px',
                  backgroundColor: '#f5f5f5',
                  borderRadius: '4px',
                  fontSize: '14px'
                }}>                  {headingFontSize}px
                </span>
              </div>
              <div style={{ 
                marginTop: '15px', 
                padding: '20px', 
                border: '1px solid #ddd', 
                borderRadius: '8px',
                minHeight: '120px',
                width: '100%',
                maxWidth: '600px',
                backgroundColor: '#fafafa',
                display: 'flex',
                flexDirection: 'column',
                justifyContent: 'center'
              }}>
                <h3 style={{ 
                  margin: '0 0 12px 0',
                  fontSize: `${headingFontSize}px`,
                  fontWeight: '600'
                }}>
                  Vista previa del título
                </h3>
                <p style={{ 
                  margin: '0',
                  fontSize: `${baseFontSize}px`,
                  color: '#666'
                }}>
                  Texto normal para comparación de tamaños.
                </p>
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
      <div className="info-card" style={{ marginTop: '30px' }}>
        <div className="card-header">
          <h2>📄 Plantillas de Factura</h2>
        </div>
        
        <div className="plantilla-options">
          {/* Plantilla activa */}
          {plantillaActiva && (
            <div className="plantilla-activa" style={{ 
              padding: '15px', 
              backgroundColor: '#f0f7ff', 
              borderRadius: '8px',
              marginBottom: '20px',
              border: '1px solid #90caf9'
            }}>
              <h3 style={{ margin: '0 0 10px 0', color: '#1976d2' }}>Plantilla Activa:</h3>
              <div style={{ display: 'flex', alignItems: 'center' }}>
                <span className="icono-archivo" style={{ 
                  marginRight: '15px',
                  fontSize: '28px',
                  color: '#1976d2'
                }}>
                  📄
                </span>
                <div>
                  <p style={{ margin: '0 0 5px 0', fontWeight: 'bold' }}>{plantillaActiva.nombre}</p>
                  {plantillaActiva.descripcion && (
                    <p style={{ margin: '0', fontSize: '0.9em', color: '#666' }}>{plantillaActiva.descripcion}</p>
                  )}
                  <p style={{ margin: '5px 0 0 0', fontSize: '0.8em', color: '#666' }}>
                    Fecha: {plantillaActiva.fecha_creacion}
                  </p>
                </div>
              </div>
            </div>
          )}
          
          {/* Información sobre plantilla activa */}
          {plantillaActiva ? (
            <div className="plantilla-activa-info" style={{
              backgroundColor: '#e6f7ff',
              border: '1px solid #91d5ff',
              borderRadius: '4px',
              padding: '10px',
              marginBottom: '20px'
            }}>
              <p style={{ margin: 0 }}>
                <strong>Plantilla Activa:</strong> {plantillaActiva.nombre}
              </p>
              <p style={{ margin: '5px 0 0 0', fontSize: '0.9em' }}>
                Esta plantilla se utilizará automáticamente para generar todas las facturas.
              </p>
            </div>
          ) : (
            <div className="plantilla-activa-info" style={{
              backgroundColor: '#fff2e8',
              border: '1px solid #ffbb96',
              borderRadius: '4px',
              padding: '10px',
              marginBottom: '20px'
            }}>
              <p style={{ margin: 0 }}>
                <strong>No hay plantilla activa.</strong> Se usará la plantilla por defecto para generar facturas.
              </p>
              <p style={{ margin: '5px 0 0 0', fontSize: '0.9em' }}>
                Activa una plantilla para utilizarla en todas las facturas.
              </p>
            </div>
          )}
          
          {/* Lista de plantillas */}
          {plantillas.length > 0 && (
            <div className="plantillas-lista" style={{ marginBottom: '25px' }}>
              <h3>Todas las Plantillas:</h3>
              <div style={{ 
                maxHeight: '250px',
                overflowY: 'auto',
                border: '1px solid #ddd',
                borderRadius: '4px',
                padding: '10px'
              }}>
                {plantillas.map(plantilla => (
                  <div key={plantilla.id} style={{ 
                    padding: '10px',
                    borderBottom: '1px solid #eee',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between',
                    backgroundColor: plantilla.activa ? '#f5f5f5' : 'transparent'
                  }}>
                    <div style={{ display: 'flex', alignItems: 'center', flex: 1 }}>
                      <span style={{ marginRight: '10px', fontSize: '20px' }}>📄</span>
                      <div>
                        <p style={{ margin: '0 0 3px 0', fontWeight: plantilla.activa ? 'bold' : 'normal' }}>
                          {plantilla.nombre}
                          {plantilla.activa && <span style={{ color: '#2e7d32', marginLeft: '8px' }}>• Activa</span>}
                        </p>
                        {plantilla.descripcion && (
                          <p style={{ margin: '0', fontSize: '0.85em', color: '#666' }}>
                            {plantilla.descripcion}
                          </p>
                        )}
                      </div>
                    </div>
                    <div>
                      {!plantilla.activa && (
                        <button 
                          onClick={() => activarPlantilla(plantilla.id)}
                          style={{ 
                            marginRight: '8px',
                            backgroundColor: actionButtonsColorFactura,
                            color: 'white',
                            border: 'none',
                            padding: '5px 10px',
                            borderRadius: '4px',
                            cursor: 'pointer'
                          }}
                        >
                          Activar
                        </button>
                      )}
                      <button 
                        onClick={() => eliminarPlantilla(plantilla.id)}
                        style={{ 
                          backgroundColor: deleteButtonsColorFactura,
                          color: 'white',
                          border: 'none',
                          padding: '5px 10px',
                          borderRadius: '4px',
                          cursor: 'pointer'
                        }}
                      >
                        Eliminar
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}
          
          {/* Subir nueva plantilla */}
          <div className="plantilla-upload">
            <h3>Subir Nueva Plantilla</h3>
            
            {/* Plantilla de ejemplo */}
            <div className="template-example" style={{ 
              marginBottom: '25px', 
              padding: '15px', 
              backgroundColor: '#f0f8ff', 
              borderRadius: '8px',
              border: '1px solid #c2e0ff'
            }}>
              <h4 style={{ margin: '0 0 10px 0', color: '#0066cc' }}>
                <span style={{ marginRight: '8px' }}>ℹ️</span>
                ¿No sabes cómo estructurar tu plantilla?
              </h4>
              
              <p style={{ margin: '0 0 12px 0', fontSize: '0.9em', color: '#444' }}>
                Descarga nuestra plantilla de ejemplo que muestra la estructura correcta y todas las variables disponibles para crear tus propias facturas personalizadas.
              </p>
              
              <div style={{ display: 'flex', alignItems: 'center' }}>
                {/* Actualizar icono y descripción */}
                <div style={{ 
                  width: '40px', 
                  height: '40px', 
                  backgroundColor: '#2b579a', // Color azul de Word
                  borderRadius: '4px',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  marginRight: '12px',
                  color: 'white',
                  fontWeight: 'bold',
                  fontSize: '16px'
                }}>
                  W
                </div>
                <div>
                  <p style={{ margin: '0 0 4px 0', fontWeight: 'bold' }}>
                    plantilla_ejemplo_factura.docx
                  </p>
                  <p style={{ margin: '0', fontSize: '0.8em', color: '#666' }}>
                    Documento Word con variables de ejemplo
                  </p>
                </div>                <button 
                  onClick={() => {
                    // Primer intento con fetch
                    descargarPlantillaEjemplo();
                    
                    // Si falla, intentar método directo después de 2 segundos
                    setTimeout(() => {
                      if (!document.querySelector('a[download="plantilla_ejemplo_factura.docx"]')) {
                        console.log('Intentando método de descarga alternativo...');
                        descargarPlantillaEjemploDirecto();
                      }
                    }, 2000);
                  }}
                  style={{ 
                    marginLeft: 'auto',
                    backgroundColor: actionButtonsColorFactura || '#0078d4',
                    color: 'white',
                    border: 'none',
                    padding: '8px 16px',
                    borderRadius: '4px',
                    cursor: 'pointer',
                    display: 'flex',
                    alignItems: 'center',
                    gap: '8px'
                  }}
                >
                  <span style={{ fontSize: '16px' }}>⬇️</span> 
                  Descargar Ejemplo
                </button>
              </div>
            </div>

            <div className="info-group" style={{ marginBottom: '15px' }}>
              <label htmlFor="descripcion-plantilla">Descripción (opcional):</label>
              <input 
                type="text"
                id="descripcion-plantilla"
                value={descripcionPlantilla}
                onChange={(e) => setDescripcionPlantilla(e.target.value)}
                className="theme-select"
                placeholder="Ej: Plantilla para facturas de servicios"
              />
            </div>
            
            <label htmlFor="plantilla-upload">Seleccionar Archivo Word:</label>
            <div className="file-upload-container">
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
                  backgroundColor: fileSelectButtonsColorFactura
                }}
              >
                {plantillaWord ? 'Cambiar plantilla' : 'Seleccionar plantilla'}
              </button>
              
              {plantillaWord && (
                <span className="file-name" style={{ marginLeft: '10px' }}>
                  {plantillaWord.name}
                </span>
              )}
            </div>
            
            <div className="file-requirements">
              <small>Formatos permitidos: DOC, DOCX. Tamaño máximo: 5MB</small>
            </div>
            
            <div className="plantilla-actions" style={{ 
              display: 'flex', 
              gap: '15px',
              marginTop: '20px' 
            }}>
              {loading && (
                <div className="loading-indicator" style={{ 
                  marginRight: '15px',
                  display: 'flex',
                  alignItems: 'center' 
                }}>
                  <div style={{
                    width: '20px',
                    height: '20px',
                    border: '3px solid #f3f3f3',
                    borderTop: '3px solid #3498db',
                    borderRadius: '50%',
                    animation: 'spin 1s linear infinite',
                    marginRight: '8px'
                  }}></div>
                  <span>Procesando...</span>
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
                  opacity: !plantillaWord || loading ? 0.7 : 1
                }}
              >
                Guardar y Activar Plantilla              </button>
            </div>
          </div>
        </div>
      </div>
      )}
    </div>
  );
}

export default Preferencias;
