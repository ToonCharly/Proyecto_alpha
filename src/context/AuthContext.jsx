import React, { createContext, useState, useContext, useEffect } from 'react';
import { migrationComplete } from '../utils/migrationUtils.jsx';

// Crear el contexto de autenticación
const AuthContext = createContext();

// Hook personalizado para usar el contexto de autenticación
export const useAuth = () => useContext(AuthContext);

export const AuthProvider = ({ children }) => {  // Migración automática al inicializar
  useEffect(() => {
    migrationComplete();
  }, []);

  const [user, setUser] = useState(() => {
    try {
      // CAMBIO: Usar sessionStorage en lugar de localStorage para evitar conflictos entre ventanas
      const savedUser = sessionStorage.getItem('user');
      return savedUser ? JSON.parse(savedUser) : null;
    } catch (error) {
      console.error('Error al recuperar usuario de sessionStorage:', error);
      sessionStorage.removeItem('user'); // Limpiar datos corruptos
      return null;
    }
  });
  
  const [authToken, setAuthToken] = useState(() => {
    // CAMBIO: Usar sessionStorage en lugar de localStorage
    return sessionStorage.getItem('authToken') || '';
  });

  // Función para iniciar sesión
  const login = (userData, token) => {
    if (!userData || !token) {
      console.error('Error: Datos de autenticación incompletos');
      return;
    }
    
    setUser(userData);
    setAuthToken(token);
  };

  // Función para cerrar sesión
  const logout = () => {
    setUser(null);
    setAuthToken('');
  };

  // Función para hacer peticiones autenticadas
  const fetchWithAuth = async (url, options = {}) => {
    if (!authToken) {
      throw new Error('No hay sesión activa');
    }
    
    const headers = {
      ...(options.headers || {}),
      'Authorization': `Bearer ${authToken}`,
    };
    
    try {
      const response = await fetch(url, { ...options, headers });
      
      // Manejo de errores de autenticación
      if (response.status === 401) {
        console.warn('Sesión expirada o token inválido');
        logout(); // Cerrar sesión automáticamente si el token expiró
        throw new Error('Su sesión ha expirado. Por favor, inicie sesión nuevamente.');
      }
      
      return response;
    } catch (error) {
      console.error('Error en fetchWithAuth:', error);
      throw error;
    }
  };
  // Sincronizar con sessionStorage (aislado por ventana/pestaña)
  useEffect(() => {
    if (user) {
      sessionStorage.setItem('user', JSON.stringify(user));
    } else {
      sessionStorage.removeItem('user');
    }
  }, [user]);
  
  useEffect(() => {
    if (authToken) {
      sessionStorage.setItem('authToken', authToken);
    } else {
      sessionStorage.removeItem('authToken');
    }
  }, [authToken]);

  return (
    <AuthContext.Provider value={{
      user,
      authToken,
      login,
      logout,
      fetchWithAuth,
      isAuthenticated: !!authToken
    }}>
      {children}
    </AuthContext.Provider>
  );
};