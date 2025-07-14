import React, { useState, useContext, useEffect } from 'react';
import { usePreferencias } from '../context/PreferenciasContext';
import '../STYLES/HomeAdmin.css';
import { FacturaContext } from '../context/FacturaContext';
import InformacionPersonal from './InformacionPersonal';
import HistorialFacturas from './HistorialFacturas';
import Preferencias from './Preferencias';
import DatosEmpresa from './DatosEmpresa';
import AdministrarUsuarios from './AdministrarUsuarios';
import HistorialEmisor from './HistorialEmisor';

// iconos
const ClockIcon = () => <img src="/transaction-history_18281961.png" alt="Historial de Facturas" className="icon" />;
const MenuIcon = () => <img src="/icono_menu.png" alt="Menú" className="icon" />;
const CloseIcon = () => <img src="/icono_menun.png" alt="Cerrar" className="icon" />;
const ConfigIcon = () => <img src="/cogwheel_16577871.png" alt="Configuración" className="icon" />;
const BackIcon = () => <img src="/left_3734793.png" alt="Regresar" className="icon" />; 
const InfoIcon = () => <img src="/user-profile_4803060.png" alt="Información" className="icon" />; 
const LogoutIcon = () => <img src="/sign-out_6461685.png" alt="Cerrar Sesión" className="icon" />;
const EmpresaIcon = () => <img src="/information_17930919.png" alt="Empresa" className="icon" />;
const PreferenciasIcon = () => <img src="/setting_8220652.png" alt="Preferencias" className="icon" />;
const UsersIcon = () => <img src="/perfil.png" alt="Administrar Usuarios" className="icon" />; // Nuevo icono

const HomeAdmin = () => {
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [showLogoutModal, setShowLogoutModal] = useState(false);
  const [activeSection, setActiveSection] = useState('historialFacturas');
  const [configMode, setConfigMode] = useState(false);
  const [userData, setUserData] = useState({ username: '', email: '' });
  const [companyLogo, setCompanyLogo] = useState(''); // Añadir estado para el logo
  
  // Usar el hook usePreferencias para obtener valores y funciones del contexto
  const { companyName, companyTextColor } = usePreferencias();
  
  const { historialFacturas } = useContext(FacturaContext);
  
  useEffect(() => {
    const handleNavigateEvent = (event) => {
      if (event.detail && event.detail.section) {
        setActiveSection(event.detail.section);
      }
    };
    
    window.addEventListener('navigateToSection', handleNavigateEvent);
    
    // Load user data
    try {
      const storedUserData = JSON.parse(sessionStorage.getItem('userData')); // Migrado a sessionStorage
      if (storedUserData) {
        setUserData({
          username: storedUserData.username || '',
          email: storedUserData.email || ''
        });
      }
    } catch (error) {
      console.error('Error loading user data:', error);
    }

    // Cargar el logo de la empresa
    const savedLogo = localStorage.getItem('appLogo');
    if (savedLogo) {
      setCompanyLogo(savedLogo);
    }
    return () => {
      window.removeEventListener('navigateToSection', handleNavigateEvent);
    };
  }, []);

  // Listener para actualización automática del logo
  useEffect(() => {
    // Función para manejar cambios en localStorage
    const handleStorageChange = (e) => {
      if (e.key === 'appLogo') {
        setCompanyLogo(e.newValue || '');
      }
    };

    // Escuchar cambios desde otras ventanas/pestañas
    window.addEventListener('storage', handleStorageChange);

    // Para cambios en la misma ventana, interceptar setItem
    const originalSetItem = localStorage.setItem;
    localStorage.setItem = function(key, value) {
      originalSetItem.apply(this, arguments);
      if (key === 'appLogo') {
        setCompanyLogo(value || '');
      }
    };

    // Para cambios en la misma ventana, interceptar removeItem
    const originalRemoveItem = localStorage.removeItem;
    localStorage.removeItem = function(key) {
      originalRemoveItem.apply(this, arguments);
      if (key === 'appLogo') {
        setCompanyLogo('');
      }
    };

    // Cleanup al desmontar el componente
    return () => {
      window.removeEventListener('storage', handleStorageChange);
      localStorage.setItem = originalSetItem;
      localStorage.removeItem = originalRemoveItem;
    };
  }, []);

  // Agregar este useEffect a tu componente Home
  useEffect(() => {
    // Cargar tema desde localStorage
    const savedTheme = localStorage.getItem('sidebarTheme');
    if (savedTheme) {
      try {
        const themeData = JSON.parse(savedTheme);
        document.documentElement.style.setProperty('--sidebar-color', themeData.color);
      } catch (error) {
        console.error('Error al aplicar el tema guardado:', error);
      }
    }

    // Cargar fuentes del panel admin (sin prefijo)
    const fontFamily = localStorage.getItem('appFontFamily');
    if (fontFamily) {
      document.documentElement.style.setProperty('--admin-app-font-family', fontFamily);
    }
    
    const headingFontFamily = localStorage.getItem('appHeadingFontFamily');
    if (headingFontFamily) {
      document.documentElement.style.setProperty('--admin-app-heading-font-family', headingFontFamily);
    }
  }, []);

  const toggleSidebar = () => {
    setSidebarOpen(!sidebarOpen);
  };

  const handleLogout = () => {
    localStorage.removeItem('userData');
    window.location.href = '/login';
  };

  // Centraliza la navegación
  // Nueva lógica: Mantener configMode=true mientras se navega entre subsecciones de configuración
  const handleNavigation = (section) => {
    if (section === 'configuracion') {
      setConfigMode(true);
      setActiveSection('informacionPersonal');
    } else if (
      section === 'informacionPersonal' ||
      section === 'datosEmpresa' ||
      section === 'preferencias'
    ) {
      setConfigMode(true);
      setActiveSection(section);
    } else {
      setConfigMode(false);
      setActiveSection(section);
    }
  };
  
  const exitConfigMode = () => {
    setConfigMode(false);
    setActiveSection('historialFacturas');
  };

const facturasCount = historialFacturas ? historialFacturas.length : 0;

  const renderActiveSection = () => {
    switch (activeSection) {
      case 'historialFacturas':
        return <HistorialFacturas />;
      case 'historialEmisor':
        return <HistorialEmisor />;
      case 'administrarUsuarios':
        return <AdministrarUsuarios />;
      case 'informacionPersonal':
        return <InformacionPersonal />;
      case 'datosEmpresa':
        return <DatosEmpresa />;
      case 'preferencias':
        return <Preferencias />;
      default:
        return <HistorialFacturas />;
    }
  };

  return (
    <div className="app-container admin-panel">
      <header className="navbar">
        <div className="navbar-left">
          <button onClick={toggleSidebar} className="menu-button">
            {sidebarOpen ? <CloseIcon /> : <MenuIcon />}
          </button>
          <span className="welcome-text" style={{ fontSize: '2rem', fontWeight: 'bold', color: companyTextColor }}>
            Portal Administrativo de {companyName}
          </span>
        </div>
        
        <div className="navbar-right">
          {companyLogo && (
            <div className="company-logo">
              <img src={companyLogo} alt="Logo de la empresa" style={{ height: '40px', maxWidth: '150px' }} />
            </div>
          )}
        </div>
      </header>

      <div className="content-container">
        <aside className={`sidebar ${sidebarOpen ? 'open' : 'closed'}`}>
          <nav className="sidebar-nav" style={{ 
            display: 'flex', 
            flexDirection: 'column', 
            height: '100%',
            justifyContent: 'space-between' 
          }}>
            {!configMode ? (
              <>
                <div>
                  <div className="nav-section">
                    <div className="nav-item" onClick={() => handleNavigation('historialFacturas')}>
                      <div className="nav-item-content">
                        <ClockIcon />
                        <span>Historial de Facturas</span>
                        {facturasCount > 0 && (
                          <span className="facturas-badge">
                            {facturasCount}
                          </span>
                        )}
                      </div>
                    </div>
                  </div>
                  <div className="nav-section">
                    <div className="nav-item" onClick={() => handleNavigation('historialEmisor')}>
                      <div className="nav-item-content">
                        <ClockIcon />
                        <span>Historial Emisor</span>
                      </div>
                    </div>
                  </div>
                  
                  <div className="nav-section">
                    <div className="nav-item" onClick={() => handleNavigation('administrarUsuarios')}>
                      <div className="nav-item-content">
                        <UsersIcon />
                        <span>Administrar Usuarios</span>
                      </div>
                    </div>
                  </div>
                </div>
                
                <div>
                  <div className="nav-section">
                    <div className="nav-item" onClick={() => handleNavigation('configuracion')}>
                      <div className="nav-item-content">
                        <ConfigIcon />
                        <span>Configuración</span>
                      </div>
                    </div>
                  </div>
                  
                  <div className="nav-section">
                    <div className="nav-item" onClick={() => setShowLogoutModal(true)}>
                      <div className="nav-item-content">
                        <LogoutIcon />
                        <span>Cerrar Sesión</span>
                      </div>
                    </div>
                  </div>
                  
                  {userData && (
                    <div className="user-info-section">
                      <div className="user-info-content">
                        <div className="user-avatar">
                          {userData.username 
                            ? userData.username.charAt(0).toUpperCase() 
                            : (userData.email ? userData.email.charAt(0).toUpperCase() : 'A')}
                        </div>
                        <div className="user-details">
                          <div className="user-name">Administrador</div>
                          <div className="user-email">{userData.email || ''}</div>
                        </div>
                      </div>
                    </div>
                  )}
                </div>
              </>
            ) : (
              <>
                <div>
                  <div className="nav-section">
                    <div className="nav-item" onClick={exitConfigMode}>
                      <div className="nav-item-content">
                        <BackIcon />
                        <span>Regresar</span>
                      </div>
                    </div>
                  </div>
                  
                  <div className="nav-section" style={{ marginTop: '20px' }}>
                    <div className="nav-item" onClick={() => handleNavigation('informacionPersonal')}>
                      <div className="nav-item-content">
                        <InfoIcon />
                        <span>Información del Usuario</span>
                      </div>
                    </div>
                  </div>
                  
                  <div className="nav-section">
                    <div className="nav-item" onClick={() => handleNavigation('datosEmpresa')}>
                      <div className="nav-item-content">
                        <EmpresaIcon />
                        <span>Información de la Empresa</span>
                      </div>
                    </div>
                  </div>
                  
                  <div className="nav-section">
                    <div className="nav-item" onClick={() => handleNavigation('preferencias')}>
                      <div className="nav-item-content">
                        <PreferenciasIcon />
                        <span>Preferencias</span>
                      </div>
                    </div>
                  </div>
                </div>
              </>
            )}
          </nav>
        </aside>

        <main className="main-content">
          {renderActiveSection()}
        </main>
      </div>

      {showLogoutModal && (
        <div className="modal-overlay">
          <div className="modal">
            <h2>¿Cerrar Sesión?</h2>
            <p>¿Estás seguro de que deseas cerrar sesión?</p>
            <div className="modal-buttons">
              <button className="btn cancel" onClick={() => setShowLogoutModal(false)}>
                Cancelar
              </button>
              <button className="btn confirm" onClick={handleLogout}>
                Cerrar Sesión
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default HomeAdmin;