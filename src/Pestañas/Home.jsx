import React, { useState, useContext, useEffect } from 'react';
import '../STYLES/Home.css';
import InicioFacturacion from './Pantalla_Principal';
import { FacturaContext } from '../context/FacturaContext';
import InformacionPersonal from './InformacionPersonal';
import Empresas from './Empresas'; 
import HistorialFacturas from './HistorialFacturas';

// Definición correcta de los iconos
const FileIcon = () => <img src="/icono_factura.png" alt="Crear Factura" className="icon" />;
const ClockIcon = () => <img src="/icono_historial.png" alt="Historial de Facturas" className="icon" />;
const RecoverIcon = () => <img src="/icono_recuperar.png" alt="Administrar Empresas" className="icon" />;
const MenuIcon = () => <img src="/icono_menu.png" alt="Menú" className="icon" />;
const CloseIcon = () => <img src="/icono_menun.png" alt="Cerrar" className="icon" />;
const UserIcon = () => <img src="/icono_user.png" alt="Usuario" className="icon" />;

const Home = () => {
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [userMenuOpen, setUserMenuOpen] = useState(false);
  const [showLogoutModal, setShowLogoutModal] = useState(false);
  const [activeSection, setActiveSection] = useState('crearFactura');
  const [username, setUsername] = useState('');

  const { historialFacturas } = useContext(FacturaContext);
  
  // Función para cargar el nombre de usuario desde localStorage
  const loadUserData = () => {
    try {
      const userData = JSON.parse(localStorage.getItem('userData'));
      if (userData && userData.username) {
        setUsername(userData.username);
      }
    } catch (error) {
      console.error('Error al leer datos de usuario:', error);
    }
  };

  // Cargar datos de usuario al montar el componente
  useEffect(() => {
    loadUserData();
  }, []);

  // Escuchar eventos de navegación y actualización de perfil
  useEffect(() => {
    // Para navegación entre secciones
    const handleNavigateEvent = (event) => {
      if (event.detail && event.detail.section) {
        setActiveSection(event.detail.section);
      }
    };
    
    // Para actualización de información de perfil
    const handleProfileUpdate = () => {
      console.log('Evento de actualización de perfil detectado');
      loadUserData(); // Recargar datos de usuario desde localStorage
    };
    
    window.addEventListener('navigateToSection', handleNavigateEvent);
    window.addEventListener('profileUpdated', handleProfileUpdate);
    
    // Limpiar event listeners
    return () => {
      window.removeEventListener('navigateToSection', handleNavigateEvent);
      window.removeEventListener('profileUpdated', handleProfileUpdate);
    };
  }, []);

  const toggleSidebar = () => {
    setSidebarOpen(!sidebarOpen);
  };

  const toggleUserMenu = () => {
    setUserMenuOpen(!userMenuOpen);
  };

  const handleLogout = () => {
    localStorage.removeItem('userData');
    window.location.href = '/login';
  };

  const handleNavigation = (section) => {
    setActiveSection(section);
    setUserMenuOpen(false);
  };

  // Contar el número de facturas en el historial
  const facturasCount = historialFacturas ? historialFacturas.length : 0;

  return (
    <div className="app-container">
      <header className="navbar">
        <div className="navbar-left">
          <button onClick={toggleSidebar} className="menu-button">
            {sidebarOpen ? <CloseIcon /> : <MenuIcon />}
          </button>
          <span className="welcome-text" style={{ fontSize: '2rem', fontWeight: 'bold' }}>
            Portal de  Facturación de (Empresa)
          </span>
        </div>

        <div className="navbar-right">
          <div className="user-menu" onClick={toggleUserMenu} style={{ cursor: 'pointer' }}>
            <UserIcon />
            <span className="user-name">{username || 'Usuario'}</span>
          </div>

          {userMenuOpen && (
            <div className="user-dropdown">
              <div
                className="dropdown-item"
                onClick={() => handleNavigation('informacionPersonal')}
                style={{ cursor: 'pointer' }}
              >
                Información Personal
              </div>
              <div
                className="dropdown-item"
                onClick={() => setShowLogoutModal(true)}
                style={{ cursor: 'pointer' }}
              >
                Cerrar Sesión
              </div>
            </div>
          )}
        </div>
      </header>

      <div className="content-container">
        <aside className={`sidebar ${sidebarOpen ? 'open' : 'closed'}`}>
          <nav className="sidebar-nav">
            <div className="nav-section">
              <div className="nav-item" onClick={() => setActiveSection('crearFactura')}>
                <div className="nav-item-content">
                  <FileIcon />
                  <span>Crear Factura</span>
                </div>
              </div>
            </div>

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

            <div className="nav-section">
              <div className="nav-item" onClick={() => setActiveSection('administrarEmpresas')}>
                <div className="nav-item-content">
                  <RecoverIcon />
                  <span>Administrar Empresas</span>
                </div>
              </div>
            </div>
          </nav>
        </aside>

        <main className="main-content">
          {activeSection === 'crearFactura' && <InicioFacturacion />}
          {activeSection === 'historialFacturas' && <HistorialFacturas />}
          {activeSection === 'administrarEmpresas' && <Empresas />}
          {activeSection === 'informacionPersonal' && <InformacionPersonal />}
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

export default Home;