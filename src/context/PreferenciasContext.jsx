import React, { createContext, useState, useEffect, useContext } from 'react';

// Crear el contexto
export const PreferenciasContext = createContext();

// Proveedor del contexto
export const PreferenciasProvider = ({ children }) => {
  const [companyName, setCompanyName] = useState('Empresa');
  const [companyTextColor, setCompanyTextColor] = useState('#000000');
  const [navbarBgColor, setNavbarBgColor] = useState('#ffffff');
  
  // Cargar preferencias guardadas
  useEffect(() => {
    const savedCompanyName = localStorage.getItem('companyName');
    if (savedCompanyName) {
      setCompanyName(savedCompanyName);
    }
    
    const savedTextColor = localStorage.getItem('companyTextColor');
    if (savedTextColor) {
      setCompanyTextColor(savedTextColor);
      document.documentElement.style.setProperty('--company-text-color', savedTextColor);
    }
    
    const savedNavbarBgColor = localStorage.getItem('navbarBgColor');
    if (savedNavbarBgColor) {
      setNavbarBgColor(savedNavbarBgColor);
      document.documentElement.style.setProperty('--navbar-bg-color', savedNavbarBgColor);
    }
  }, []);
  
  // Función para actualizar el nombre
  const updateCompanyName = (newName) => {
    setCompanyName(newName);
    localStorage.setItem('companyName', newName);
  };
  
  // Función para actualizar el color del texto
  const updateCompanyTextColor = (newColor) => {
    setCompanyTextColor(newColor);
    localStorage.setItem('companyTextColor', newColor);
    document.documentElement.style.setProperty('--company-text-color', newColor);
  };
  
  // Función para actualizar el color de fondo del navbar
  const updateNavbarBgColor = (newColor) => {
    setNavbarBgColor(newColor);
    localStorage.setItem('navbarBgColor', newColor);
    document.documentElement.style.setProperty('--navbar-bg-color', newColor);
  };
  
  return (
    <PreferenciasContext.Provider 
      value={{ 
        companyName, 
        updateCompanyName, 
        companyTextColor, 
        updateCompanyTextColor,
        navbarBgColor,
        updateNavbarBgColor
      }}
    >
      {children}
    </PreferenciasContext.Provider>
  );
};

// Hook personalizado para usar el contexto
export const usePreferencias = () => useContext(PreferenciasContext);