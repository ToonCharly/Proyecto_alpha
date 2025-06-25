import React, { createContext, useState, useEffect, useContext } from 'react';

// Crear el contexto
export const PreferenciasContext = createContext();

// Proveedor del contexto
export const PreferenciasProvider = ({ children }) => {
  const [companyName, setCompanyName] = useState('Empresa');
  const [companyTextColor, setCompanyTextColor] = useState('#000000');
  const [navbarBgColor, setNavbarBgColor] = useState('#ffffff');
  
  // Estados para tamaños de fuente globales
  const [baseFontSize, setBaseFontSize] = useState(16); // Tamaño base en px
  const [headingFontSize, setHeadingFontSize] = useState(24); // Tamaño de títulos en px
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
    
    // Cargar tamaños de fuente guardados
    const savedBaseFontSize = localStorage.getItem('baseFontSize');
    if (savedBaseFontSize) {
      const size = parseInt(savedBaseFontSize);
      setBaseFontSize(size);
      document.documentElement.style.setProperty('--base-font-size', `${size}px`);
    } else {
      document.documentElement.style.setProperty('--base-font-size', '16px');
    }
    
    const savedHeadingFontSize = localStorage.getItem('headingFontSize');
    if (savedHeadingFontSize) {
      const size = parseInt(savedHeadingFontSize);
      setHeadingFontSize(size);
      document.documentElement.style.setProperty('--heading-font-size', `${size}px`);
    } else {
      document.documentElement.style.setProperty('--heading-font-size', '24px');
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
  
  // Función para actualizar el tamaño base de fuente
  const updateBaseFontSize = (newSize) => {
    setBaseFontSize(newSize);
    localStorage.setItem('baseFontSize', newSize.toString());
    document.documentElement.style.setProperty('--base-font-size', `${newSize}px`);
  };
  
  // Función para actualizar el tamaño de fuente de títulos
  const updateHeadingFontSize = (newSize) => {
    setHeadingFontSize(newSize);
    localStorage.setItem('headingFontSize', newSize.toString());
    document.documentElement.style.setProperty('--heading-font-size', `${newSize}px`);
  };
    return (
    <PreferenciasContext.Provider 
      value={{ 
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
      }}
    >
      {children}
    </PreferenciasContext.Provider>
  );
};

// Hook personalizado para usar el contexto
export const usePreferencias = () => useContext(PreferenciasContext);