import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import '../STYLES/NavBar.css';

const NavBar = () => {
  const navigate = useNavigate();
  const [username, setUsername] = useState('Usuario');
  
  // Obtener el nombre de usuario al cargar el componente
  useEffect(() => {
    const getUserData = () => {
      try {
        const userData = JSON.parse(sessionStorage.getItem('userData')); // Migrado a sessionStorage para evitar conflictos
        if (userData && userData.username) {
          setUsername(userData.username);
        }
      } catch (error) {
        console.error('Error al leer datos de usuario:', error);
      }
    };

    getUserData();
    
    // Escuchar cambios en localStorage
    window.addEventListener('storage', getUserData);
    
    return () => {
      window.removeEventListener('storage', getUserData);
    };
  }, []);
  
  // Función para cerrar sesión
  const handleLogout = () => {
    localStorage.removeItem('userData');
    setUsername('Usuario');
    navigate('/login');
  };

  // Mostrar menú especial solo para administradores de empresa emisora
  let isAdminEmpresa = false;
  try {
    const userData = JSON.parse(sessionStorage.getItem('userData'));
    isAdminEmpresa = userData && userData.rol === 'ADMIN_EMPRESA';
  } catch {
    // Si ocurre un error al parsear userData, simplemente no mostramos el menú admin
  }

  return (
    <nav className="navbar">
      <div className="logo">
        <Link to="/home">
          <img 
            src="/upscalemedia-transformed.png" 
            alt="Logo" 
            style={{
              maxHeight: '40px',
              maxWidth: '150px',
              objectFit: 'contain'
            }}
          />
        </Link>
      </div>

      <div className="nav-links">
        <Link to="/home">Inicio</Link>
        <Link to="/facturas">Facturas</Link>
        <Link to="/historial-facturas">Historial de Facturas</Link>
        {isAdminEmpresa && (
          <Link to="/historial-emisor">Historial de Empresa</Link>
        )}
      </div>

      <div className="user-section">
        <span className="username">Hola, {username}</span>
        <div className="user-menu">
          <button className="profile-button">
            👤
          </button>
          <div className="dropdown-menu">
            <Link to="/perfil">Mi Perfil</Link>
            <Link to="/configuracion">Configuración</Link>
            <button onClick={handleLogout}>Cerrar Sesión</button>
          </div>
        </div>
      </div>
    </nav>
  );
};

export default NavBar;