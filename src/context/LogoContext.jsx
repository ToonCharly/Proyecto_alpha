import React, { createContext, useContext, useState, useEffect } from 'react';

// Crear el contexto para el logo
const LogoContext = createContext();

// Hook personalizado para usar el contexto
export const useLogo = () => {
  const context = useContext(LogoContext);
  if (!context) {
    throw new Error('useLogo debe ser usado dentro de un LogoProvider');
  }
  return context;
};

// Función auxiliar para limpiar el userId
const limpiarUserId = (userId) => {
  if (!userId || userId === 'default') return '1';
  
  // Limpiar el formato "default:1" a solo "1"
  if (userId.includes(':')) {
    return userId.split(':')[1];
  }
  
  return userId;
};

// Proveedor del contexto
export const LogoProvider = ({ children }) => {
  const [logoActivo, setLogoActivo] = useState(null);
  const [logoData, setLogoData] = useState(null);
  const [loading, setLoading] = useState(false);

  // Función para cargar el logo activo
  const cargarLogoActivo = async (forceRefresh = false) => {
    if (loading && !forceRefresh) return;
    
    setLoading(true);
    try {
      const rawUserId = localStorage.getItem('userId') || '1';
      const idUsuario = limpiarUserId(rawUserId);
      const idUsuarioNum = parseInt(idUsuario) || 1;
      
      const response = await fetch(`http://localhost:8080/api/logos/obtener-activo-json?id_usuario=${idUsuarioNum}`);
      
      if (response.ok) {
        const data = await response.json();
        if (data.exists && data.imagen_base64) {
          const logoDataUrl = `data:${data.tipo};base64,${data.imagen_base64}`;
          setLogoActivo(logoDataUrl);
          setLogoData(data);
        } else {
          setLogoActivo(null);
          setLogoData(null);
        }
      } else {
        setLogoActivo(null);
        setLogoData(null);
      }
    } catch (error) {
      console.error('Error al cargar logo activo:', error);
      setLogoActivo(null);
      setLogoData(null);
    } finally {
      setLoading(false);
    }
  };

  // Función para refrescar el logo (llamada desde Preferencias)
  const refrescarLogo = () => {
    cargarLogoActivo(true);
  };

  // Función para limpiar el logo (cuando se elimina)
  const limpiarLogo = () => {
    setLogoActivo(null);
    setLogoData(null);
  };

  // Cargar el logo al inicializar
  useEffect(() => {
    cargarLogoActivo();
  }, []);

  // Escuchar cambios en el localStorage del userId
  useEffect(() => {
    const handleStorageChange = (event) => {
      if (event.key === 'userId') {
        cargarLogoActivo(true);
      }
    };

    window.addEventListener('storage', handleStorageChange);
    return () => window.removeEventListener('storage', handleStorageChange);
  }, []);

  const value = {
    logoActivo,
    logoData,
    loading,
    refrescarLogo,
    limpiarLogo,
    cargarLogoActivo
  };

  return (
    <LogoContext.Provider value={value}>
      {children}
    </LogoContext.Provider>
  );
};

export default LogoContext;
