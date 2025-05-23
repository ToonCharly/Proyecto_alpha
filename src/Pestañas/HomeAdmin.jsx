import React, { useState, useContext, useEffect } from 'react';
import { usePreferencias } from '../context/PreferenciasContext';
import '../STYLES/HomeAdmin.css';
import { FacturaContext } from '../context/FacturaContext';
import InformacionPersonal from './InformacionPersonal';
import HistorialFacturas from './HistorialFacturas';
import Preferencias from './Preferencias';
import DatosEmpresa from './DatosEmpresa';

// iconos
const ClockIcon = () => <img src="/icono_historial.png" alt="Historial de Facturas" className="icon" />;
const MenuIcon = () => <img src="/icono_menu.png" alt="Menú" className="icon" />;
const CloseIcon = () => <img src="/icono_menun.png" alt="Cerrar" className="icon" />;
const ConfigIcon = () => <img src="settings_1550664.png" alt="Configuración" className="icon" />;
const BackIcon = () => <img src="/left-arrow_11880191.png" alt="Regresar" className="icon" />; 
const InfoIcon = () => <img src="/user_667429.png" alt="Información" className="icon" />; 
const LogoutIcon = () => <img src="/sign-out_6461685.png" alt="Cerrar Sesión" className="icon" />;
const EmpresaIcon = () => <img src="/office-building_4300059.png" alt="Empresa" className="icon" />;
const PreferenciasIcon = () => <img src="/control-panel_12765560.png" alt="Preferencias" className="icon" />;

const HomeAdmin = () => {
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [showLogoutModal, setShowLogoutModal] = useState(false);
  const [activeSection, setActiveSection] = useState('historialFacturas');
  const [configMode, setConfigMode] = useState(false);
  // Añadir estado para userData
  const [userData, setUserData] = useState({ username: '', email: '' });
  
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
      const storedUserData = JSON.parse(localStorage.getItem('userData'));
      if (storedUserData) {
        setUserData({
          username: storedUserData.username || '',
          email: storedUserData.email || ''
        });
      }
    } catch (error) {
      console.error('Error loading user data:', error);
    }

    return () => {
      window.removeEventListener('navigateToSection', handleNavigateEvent);
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
  }, []);

  const toggleSidebar = () => {
    setSidebarOpen(!sidebarOpen);
  };

  const handleLogout = () => {
    localStorage.removeItem('userData');
    window.location.href = '/login';
  };

  const handleNavigation = (section) => {
    if (section === 'configuracion') {
      setConfigMode(true);
      setActiveSection('informacionPersonal');
    } else {
      setActiveSection(section);
    }
  };
  
  const exitConfigMode = () => {
    setConfigMode(false);
    setActiveSection('historialFacturas');
  };

  const facturasCount = historialFacturas ? historialFacturas.length : 0;

  // Asegúrate de que la función renderActiveSection incluya el caso para datosEmpresa
  const renderActiveSection = () => {
    switch (activeSection) {
      case 'historialFacturas':
        return <HistorialFacturas />;
      case 'informacionPersonal':
        return <InformacionPersonal />;
      case 'datosEmpresa': // Importante: usa exactamente este nombre
        return <DatosEmpresa />;
      case 'preferencias':
        return <Preferencias />;
      default:
        return <HistorialFacturas />;
    }
  };

  return (
    <div className="app-container">
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
          <div className="company-logo"></div>
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
                    <div className="nav-item" onClick={() => setActiveSection('historialFacturas')}>
                      <div className="nav-item-content">
                        <ClockIcon />
                        <span>Historial de Facturas</span>
                        {facturasCount > 0 && (
                          <span className="facturas-badge" style={{
                            backgroundColor: '#4caf50',
                            color: 'white',
                            borderRadius: '50%',
                            padding: '2px 8px',
                            fontSize: '0.8rem',
                            marginLeft: '8px',
                            display: 'inline-flex',
                            alignItems: 'center',
                            justifyContent: 'center'
                          }}>
                            {facturasCount}
                          </span>
                        )}
                      </div>
                    </div>
                  </div>
                </div>
                
                {/* Footer options - Configuration and Logout */}
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
                  
                  {/* User info section - responsive version */}
                  <div className="user-info-section">
                    <div className={`user-info-content ${!sidebarOpen ? 'collapsed' : ''}`}>
                      {sidebarOpen ? (
                        // Full user info when sidebar is open
                        <>
                          <div className="user-name">{userData.username}</div>
                          <div className="user-email">{userData.email}</div>
                        </>
                      ) : (
                        // Compact user info (initials) when sidebar is collapsed
                        <div className="user-initials">
                          {userData.username ? userData.username.charAt(0).toUpperCase() : '?'}
                        </div>
                      )}
                    </div>
                  </div>
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
                    <div className="nav-item" onClick={() => setActiveSection('informacionPersonal')}>
                      <div className="nav-item-content">
                        <InfoIcon />
                        <span>Información del Usuario</span>
                      </div>
                    </div>
                  </div>
                  
                  <div className="nav-section">
                    <div className="nav-item" onClick={() => setActiveSection('datosEmpresa')}>
                      <div className="nav-item-content">
                        <EmpresaIcon />
                        <span>Información de la Empresa</span>
                      </div>
                    </div>
                  </div>
                  
                  <div className="nav-section">
                    <div className="nav-item" onClick={() => setActiveSection('preferencias')}>
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