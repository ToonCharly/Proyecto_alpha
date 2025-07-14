import React, { useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import { PreferenciasProvider } from './context/PreferenciasContext';
import Form from './Pestañas/Login';
import RegisterForm from './Pestañas/Register';
import InicioFacturacion from './Pestañas/Pantalla_Principal';
import RecuperarPassword from './Pestañas/RecuperarPassword';
import Home from './Pestañas/Home'; 
import InformacionPersonal from './Pestañas/InformacionPersonal';
import Empresas from './Pestañas/Empresas'; 
import HistorialFacturas from './Pestañas/HistorialFacturas';
import HistorialEmisor from './Pestañas/HistorialEmisor';
import RestablecerPassword from './Pestañas/RestablecerPassword';
import HomeAdmin from './Pestañas/HomeAdmin';
import AdministrarUsuarios from './Pestañas/AdministrarUsuarios'; 
import DatosEmpresa from './Pestañas/DatosEmpresa';

function App() {
  // Agregar este useEffect para cargar el tema y logo al inicio
  useEffect(() => {
    // Cargar tema
    const savedTheme = localStorage.getItem('sidebarTheme');
    if (savedTheme) {
      try {
        const themeData = JSON.parse(savedTheme);
        document.documentElement.style.setProperty('--sidebar-color', themeData.color);
      } catch (error) {
        console.error('Error al aplicar el tema guardado:', error);
      }
    }
    
    // Cargar logo si existe
    const savedLogo = localStorage.getItem('appLogo');
    if (savedLogo) {
      document.documentElement.style.setProperty('--app-logo', `url(${savedLogo})`);
    }
    
    // Cargar color de texto si existe
    const savedTextColor = localStorage.getItem('companyTextColor');
    if (savedTextColor) {
      document.documentElement.style.setProperty('--company-text-color', savedTextColor);
    }

    // Cargar color de fondo del navbar
    const savedNavbarBgColor = localStorage.getItem('navbarBgColor');
    if (savedNavbarBgColor) {
      document.documentElement.style.setProperty('--navbar-bg-color', savedNavbarBgColor);
    }
  }, []);

  return (
    <div className="admin-panel">
      <AuthProvider>
        <PreferenciasProvider>
          <Router>
            <Routes>
              {/* Rutas públicas */}
              <Route path="/" element={<Form />} />
              <Route path="/login" element={<Navigate to="/" replace />} />
              <Route path="/register" element={<RegisterForm />} />
              <Route path="/recuperar-password" element={<RecuperarPassword />} />
              <Route path="/restablecer-password" element={<RestablecerPassword />} />
              
              {/* Rutas protegidas */}
              <Route path="/Home" element={<Home />} />
              <Route path="/homeadmin" element={<HomeAdmin />} />
              <Route path="/facturacion" element={<InicioFacturacion />} />
              <Route path="/informacion-personal" element={<InformacionPersonal />} />
              
              {/* Ruta para administrar empresas - usando el componente correcto */}
              <Route path="/empresas" element={<Empresas />} />
              <Route path="/empresas/:userId" element={<Empresas />} />
              
              {/* Nueva ruta para historial de facturas */}
              <Route path="/historial-facturas" element={<HistorialFacturas />} />
              <Route path="/historial-facturas/:userId" element={<HistorialFacturas />} />
              {/* Nueva ruta para historial de empresa emisora (solo admin empresa) */}
              <Route path="/historial-emisor" element={<HistorialEmisor />} />
              
              {/* Nueva ruta para administrar usuarios */}
              <Route path="/admin/usuarios" element={<AdministrarUsuarios />} />
              
              {/* Nueva ruta para el dashboard de la empresa */}
              <Route path="/dashboard" element={<DatosEmpresa />} />
            </Routes>
          </Router>
        </PreferenciasProvider>
      </AuthProvider>
    </div>
  );
}

export default App;