import React, { useState, useContext, useEffect } from 'react';
import '../STYLES/Home.css';
import InicioFacturacion from './Pantalla_Principal';
import { FacturaContext } from '../context/FacturaContext';
import InformacionPersonal from './InformacionPersonal';
import Empresas from './Empresas'; 
import HistorialFacturas from './HistorialFacturas';
import { usePreferencias } from '../context/PreferenciasContext'; // Añadir esta importación

// Definición correcta de los iconos
const FileIcon = () => <img src="/bill_13140958.png" alt="Crear Factura" className="icon" />;
const ClockIcon = () => <img src="/transaction-history_18281961.png" alt="Historial de Facturas" className="icon" />;
const RecoverIcon = () => <img src="/office-building_4300058.png" alt="Administrar Empresas" className="icon" />;
const MenuIcon = () => <img src="/icono_menu.png" alt="Menú" className="icon" />;
const CloseIcon = () => <img src="/icono_menun.png" alt="Cerrar" className="icon" />;
const UserIcon = () => <img src="/user-profile_4803060.png" alt="Usuario" className="icon" />;
const LogoutIcon = () => <img src="/sign-out_6461685.png" alt="Cerrar Sesión" className="icon" />;

const Home = () => {
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [showLogoutModal, setShowLogoutModal] = useState(false);
  const [activeSection, setActiveSection] = useState('crearFactura');
  const [username, setUsername] = useState('');
  const [userData, setUserData] = useState(null);
  const [companyLogo, setCompanyLogo] = useState(''); // Estado para almacenar el logo

  // Usar el contexto de preferencias para obtener los valores personalizados
  const { companyName, companyTextColor, navbarBgColor } = usePreferencias();

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

  // Cargar datos de usuario y logo al montar el componente
  useEffect(() => {
    loadUserData();
    
    // Cargar el logo de la empresa
    const savedLogo = localStorage.getItem('appLogo');
    if (savedLogo) {
      setCompanyLogo(savedLogo);
    }
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

  // Cargar el tema personalizado para el panel de facturación
  useEffect(() => {
    // Cargar tema del panel de usuario configurado por el admin
    const savedUserTheme = localStorage.getItem('userPanelTheme');
    if (savedUserTheme) {
      try {
        const themeData = JSON.parse(savedUserTheme);
        
        // Aplicar color del sidebar
        if (themeData.color || themeData.sidebarColor) {
          document.documentElement.style.setProperty('--user-sidebar-color', themeData.sidebarColor || themeData.color);
        }
        
        // Aplicar colores de acento y otros elementos si existen
        if (themeData.accentColor) {
          document.documentElement.style.setProperty('--user-accent-color', themeData.accentColor);
        }
      } catch (error) {
        console.error('Error al aplicar el tema de usuario:', error);
      }
    }
    
    // Cargar colores de botones específicos del panel de facturación
    const prefix = 'factura_';
    
    // Botones de acción
    const actionButtonsColor = localStorage.getItem(`${prefix}actionButtonsColor`);
    if (actionButtonsColor) {
      document.documentElement.style.setProperty('--user-action-button-color', actionButtonsColor);
    }
    
    // Botones de eliminar
    const deleteButtonsColor = localStorage.getItem(`${prefix}deleteButtonsColor`);
    if (deleteButtonsColor) {
      document.documentElement.style.setProperty('--user-delete-button-color', deleteButtonsColor);
    }
    
    // Botones de editar
    const editButtonsColor = localStorage.getItem(`${prefix}editButtonsColor`);
    if (editButtonsColor) {
      document.documentElement.style.setProperty('--user-edit-button-color', editButtonsColor);
    }
    
    // Botones de selección de archivo
    const fileSelectButtonsColor = localStorage.getItem(`${prefix}fileSelectButtonsColor`);
    if (fileSelectButtonsColor) {
      document.documentElement.style.setProperty('--user-file-select-button-color', fileSelectButtonsColor);
    }
    
    // Fuentes personalizadas
    const fontFamily = localStorage.getItem(`${prefix}appFontFamily`);
    if (fontFamily) {
      document.documentElement.style.setProperty('--user-app-font-family', fontFamily);
    }
    
    const headingFontFamily = localStorage.getItem(`${prefix}appHeadingFontFamily`);
    if (headingFontFamily) {
      document.documentElement.style.setProperty('--user-app-heading-font-family', headingFontFamily);
    }
  }, []);

  useEffect(() => {
    // Función para obtener los datos del usuario
    const fetchUserData = async () => {
      try {
        // Si guardas el token en localStorage
        const token = localStorage.getItem('token');
        
        if (token) {
          // Hacer la petición al servidor para obtener datos del usuario
          const response = await fetch('/api/usuario/perfil', {
            headers: {
              'Authorization': `Bearer ${token}`
            }
          });
          
          if (response.ok) {
            const data = await response.json();
            setUserData(data);
            
            // Opcional: guardar en localStorage para acceso rápido
            localStorage.setItem('userData', JSON.stringify(data));
          }
        } else {
          // Alternativamente, intenta obtener del localStorage si ya existe
          const cachedUserData = localStorage.getItem('userData');
          if (cachedUserData) {
            setUserData(JSON.parse(cachedUserData));
          }
        }
      } catch (error) {
        console.error('Error al obtener datos del usuario:', error);
      }
    };
    
    fetchUserData();
  }, []);
  
  const toggleSidebar = () => {
    setSidebarOpen(!sidebarOpen);
  };

  const handleLogout = () => {
    localStorage.removeItem('userData');
    window.location.href = '/login';
  };

  // Contar el número de facturas en el historial
  const facturasCount = historialFacturas ? historialFacturas.length : 0;

  return (
    <div className="app-container user-panel">
      <header className="navbar" style={{ backgroundColor: navbarBgColor || '#fff' }}>
        <div className="navbar-left">
          <button onClick={toggleSidebar} className="menu-button">
            {sidebarOpen ? <CloseIcon /> : <MenuIcon />}
          </button>
          
          <span 
            className="welcome-text" 
            style={{ 
              fontSize: '2rem', 
              fontWeight: 'bold',
              color: companyTextColor || 'inherit'
            }}
          >
            Portal de Facturación de {companyName || '(Empresa)'}
          </span>
        </div>
        
        {/* Agregar navbar-right para contenido alineado a la derecha */}
        <div className="navbar-right">
          {/* Logo de la empresa en el lado derecho */}
          {companyLogo && (
            <div className="company-logo">
              <img src={companyLogo} alt="Logo de la empresa" style={{ height: '40px', maxWidth: '150px' }} />
            </div>
          )}
        </div>
      </header>

      <div className="content-container">
        <aside className={`sidebar ${sidebarOpen ? 'open' : 'closed'}`}>
          <nav className="sidebar-nav">
            <div className="nav-section">
              <div 
                className={`nav-item ${activeSection === 'crearFactura' ? 'active' : ''}`}
                onClick={() => setActiveSection('crearFactura')}
              >
                <div className="nav-item-content">
                  <FileIcon />
                  <span>Crear Factura</span>
                </div>
              </div>
            </div>

            <div className="nav-section">
              <div 
                className={`nav-item ${activeSection === 'historialFacturas' ? 'active' : ''}`}
                onClick={() => setActiveSection('historialFacturas')}
              >
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
              <div 
                className={`nav-item ${activeSection === 'administrarEmpresas' ? 'active' : ''}`}
                onClick={() => setActiveSection('administrarEmpresas')}
              >
                <div className="nav-item-content">
                  <RecoverIcon />
                  <span>Administrar Empresas</span>
                </div>
              </div>
            </div>
          </nav>
          
          <div className="config-section">
            <div 
              className={`nav-item ${activeSection === 'informacionPersonal' ? 'active' : ''}`}
              onClick={() => setActiveSection('informacionPersonal')}
            >
              <div className="nav-item-content">
                <UserIcon />
                <span>Información Personal</span>
              </div>
            </div>
          </div>

          <div className="logout-section">
            <div className="nav-item" onClick={() => setShowLogoutModal(true)}>
              <div className="nav-item-content">
                <LogoutIcon />
                <span>Cerrar Sesión</span>
              </div>
            </div>
          </div>
          
          {/* Información del usuario */}
          {userData && (
            <div className="user-info-section">
              <div className="user-info-content">
                <div className="user-avatar">
                  {userData.nombre 
                    ? userData.nombre.charAt(0).toUpperCase() 
                    : (username ? username.charAt(0).toUpperCase() : 'U')}
                </div>
                <div className="user-details">
                  <div className="user-name">{userData.nombre || username || 'Usuario'}</div>
                  <div className="user-email">{userData.email || ''}</div>
                </div>
              </div>
            </div>
          )}
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
              <button 
                className="btn cancel" 
                onClick={() => setShowLogoutModal(false)}
                style={{ backgroundColor: 'var(--user-delete-button-color, #d32f2f)' }}
              >
                Cancelar
              </button>
              <button 
                className="btn confirm" 
                onClick={handleLogout}
                style={{ backgroundColor: 'var(--user-action-button-color, #2e7d32)' }}
              >
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