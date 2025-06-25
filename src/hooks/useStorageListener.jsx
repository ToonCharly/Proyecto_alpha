import { useState, useEffect } from 'react';

// Hook personalizado para escuchar cambios en localStorage/sessionStorage
export const useStorageListener = (storageKey, defaultValue = null, storageType = 'localStorage') => {
  const storage = storageType === 'sessionStorage' ? sessionStorage : localStorage;
  
  const [value, setValue] = useState(() => {
    try {
      const item = storage.getItem(storageKey);
      return item ? JSON.parse(item) : defaultValue;
    } catch (error) {
      console.error(`Error al leer ${storageKey} de ${storageType}:`, error);
      return defaultValue;
    }
  });

  useEffect(() => {
    const handleStorageChange = (e) => {
      if (e.key === storageKey) {
        try {
          const newValue = e.newValue ? JSON.parse(e.newValue) : defaultValue;
          setValue(newValue);
        } catch (error) {
          console.error(`Error al parsear ${storageKey}:`, error);
          setValue(defaultValue);
        }
      }
    };

    // Escuchar cambios desde otras ventanas/pestañas
    window.addEventListener('storage', handleStorageChange);

    // Función personalizada para escuchar cambios en la misma ventana
    const handleCustomStorageChange = (event) => {
      if (event.detail && event.detail.key === storageKey) {
        setValue(event.detail.newValue);
      }
    };

    window.addEventListener('localStorageUpdate', handleCustomStorageChange);

    return () => {
      window.removeEventListener('storage', handleStorageChange);
      window.removeEventListener('localStorageUpdate', handleCustomStorageChange);
    };
  }, [storageKey, defaultValue, storageType]);

  return [value, setValue];
};

// Función auxiliar para disparar eventos de actualización manual
export const triggerStorageUpdate = (key, newValue) => {
  const event = new CustomEvent('localStorageUpdate', {
    detail: { key, newValue }
  });
  window.dispatchEvent(event);
};

// Hook específico para preferencias del sistema
export const usePreferenceListener = () => {
  const [appLogo] = useStorageListener('appLogo', '');
  const [sidebarTheme] = useStorageListener('sidebarTheme', null);
  const [userPanelTheme] = useStorageListener('userPanelTheme', null);
  
  // Estados para colores de botones - Admin
  const [adminActionButtonColor] = useStorageListener('actionButtonsColor', '#2e7d32');
  const [adminDeleteButtonColor] = useStorageListener('deleteButtonsColor', '#d32f2f');
  const [adminEditButtonColor] = useStorageListener('editButtonsColor', '#f57c00');
  const [adminFileSelectButtonColor] = useStorageListener('fileSelectButtonsColor', '#1976d2');
  
  // Estados para colores de botones - Facturación
  const [userActionButtonColor] = useStorageListener('factura_actionButtonsColor', '#2e7d32');
  const [userDeleteButtonColor] = useStorageListener('factura_deleteButtonsColor', '#d32f2f');
  const [userEditButtonColor] = useStorageListener('factura_editButtonsColor', '#f57c00');
  const [userFileSelectButtonColor] = useStorageListener('factura_fileSelectButtonsColor', '#1976d2');

  // Estados para fuentes - Admin
  const [adminFontId] = useStorageListener('appFontId', 'roboto');
  const [adminHeadingFontId] = useStorageListener('appHeadingFontId', 'roboto');
  
  // Estados para fuentes - Facturación
  const [userFontId] = useStorageListener('factura_appFontId', 'roboto');
  const [userHeadingFontId] = useStorageListener('factura_appHeadingFontId', 'roboto');

  return {
    // Logo
    appLogo,
    
    // Temas
    sidebarTheme,
    userPanelTheme,
    
    // Colores de botones - Admin
    adminActionButtonColor,
    adminDeleteButtonColor,
    adminEditButtonColor,
    adminFileSelectButtonColor,
    
    // Colores de botones - Facturación
    userActionButtonColor,
    userDeleteButtonColor,
    userEditButtonColor,
    userFileSelectButtonColor,
    
    // Fuentes - Admin
    adminFontId,
    adminHeadingFontId,
    
    // Fuentes - Facturación
    userFontId,
    userHeadingFontId
  };
};
