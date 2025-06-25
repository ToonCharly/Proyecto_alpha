import React, { useState } from 'react';
        {subsecciones.map((subseccion) => (
          <button
            key={subseccion.id}
            className={`subseccion-card ${subseccionActiva === subseccion.id ? 'activa' : ''}`}
            onClick={() => setSubseccionActiva(subseccion.id)}
            style={{
              backgroundColor: subseccionActiva === subseccion.id ? subseccion.color : '#ffffff',
              borderColor: subseccionActiva === subseccion.id ? '#2196f3' : '#e0e0e0',
              transform: subseccionActiva === subseccion.id ? 'translateY(-2px)' : 'none',
              boxShadow: subseccionActiva === subseccion.id ? '0 4px 12px rgba(33, 150, 243, 0.2)' : '0 2px 4px rgba(0,0,0,0.1)'
            }}
          > '../STYLES/PreferenciasTabs.css';

// Componente para la navegación de subsecciones
function SubseccionNav({ subseccionActiva, setSubseccionActiva }) {
  const subsecciones = [
    { 
      id: 'botones', 
      nombre: '🎨 Colores y Botones', 
      descripcion: 'Personalizar colores de botones y tema',
      color: '#e3f2fd'
    },
    { 
      id: 'empresa', 
      nombre: '🏢 Logo y Empresa', 
      descripcion: 'Gestionar logo e información de empresa',
      color: '#f3e5f5'
    },
    { 
      id: 'plantillas', 
      nombre: '📄 Plantillas Word', 
      descripcion: 'Subir y gestionar plantillas de facturación',
      color: '#e8f5e8'
    },
    { 
      id: 'facturas', 
      nombre: '🧾 Config. Facturas', 
      descripcion: 'Configuraciones específicas de facturación',
      color: '#fff3e0'
    },
    { 
      id: 'avanzado', 
      nombre: '⚙️ Configuración Avanzada', 
      descripcion: 'Opciones avanzadas del sistema',
      color: '#fce4ec'
    }
  ];

  return (
    <div className="subseccion-nav-container">
      <h3 className="subseccion-nav-title">Configuraciones Disponibles:</h3>
      <div className="subseccion-nav-grid">
        {subsecciones.map((subseccion) => (
          <button
            key={subseccion.id}
            className={`subseccion-nav-card ${subseccionActiva === subseccion.id ? 'active' : ''}`}
            onClick={() => setSubseccionActiva(subseccion.id)}
          >
            <div className="subseccion-nav-icon">{subseccion.nombre}</div>
            <div className="subseccion-nav-desc">{subseccion.descripcion}</div>
          </button>
        ))}
      </div>
    </div>
  );
}

// Componente para la sección de Temas y Colores
function SeccionTemas({ panelActivo, selectedTheme, selectedThemeFactura, customColor, customColorFactura, handleThemeChange, themeOptions }) {
  return (
    <div className="seccion-content">
      <div className="info-card">
        <div className="card-header">
          <h2>🎨 Personalización de Colores</h2>
          <p>Cambia los colores del {panelActivo === 'admin' ? 'Panel Administrativo' : 'Panel de Facturación'}</p>
        </div>
        
        <div className="theme-options">
          <div className="info-group">
            <label htmlFor="theme-select">Color del Panel:</label>
            <select 
              id="theme-select" 
              value={panelActivo === 'admin' ? selectedTheme : selectedThemeFactura} 
              onChange={handleThemeChange}
              className="theme-select"
            >
              <option value="default">Azul por defecto</option>
              <option value="green">Verde</option>
              <option value="purple">Morado</option>
              <option value="orange">Naranja</option>
              <option value="red">Rojo</option>
              <option value="custom">Color personalizado</option>
            </select>
          </div>
          
          {((panelActivo === 'admin' && selectedTheme === 'custom') || 
            (panelActivo === 'facturacion' && selectedThemeFactura === 'custom')) && (
            <div className="info-group">
              <label htmlFor="custom-color">Color personalizado:</label>
              <input 
                type="color" 
                id="custom-color"
                value={panelActivo === 'admin' ? customColor : customColorFactura}
                onChange={handleThemeChange}
                className="color-picker"
              />
            </div>
          )}
          
          <div className="theme-preview-container">
            <h4>Vista Previa:</h4>
            <div className="theme-preview-flex">
              <div className="theme-preview-sidebar" 
                   style={{ 
                     backgroundColor: panelActivo === 'admin' 
                       ? (selectedTheme === 'custom' ? customColor : themeOptions.find(t => t.id === selectedTheme)?.color)
                       : (selectedThemeFactura === 'custom' ? customColorFactura : themeOptions.find(t => t.id === selectedThemeFactura)?.color)
                   }}>
                <div className="preview-menu-item">Inicio</div>
                <div className="preview-menu-item">Facturas</div>
                <div className="preview-menu-item active">Configuración</div>
              </div>
              <div className="theme-preview-content">
                <div className="preview-header">
                  {panelActivo === 'admin' ? 'Panel Administrativo' : 'Panel de Facturación'}
                </div>
                <div className="preview-content">Vista previa del tema</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

// Componente para la sección de Botones
function SeccionBotones({ panelActivo, actionButtonsColor, deleteButtonsColor, editButtonsColor, fileSelectButtonsColor, 
                         actionButtonsColorFactura, deleteButtonsColorFactura, editButtonsColorFactura, fileSelectButtonsColorFactura,
                         setActionButtonsColor, setDeleteButtonsColor, setEditButtonsColor, setFileSelectButtonsColor,
                         setActionButtonsColorFactura, setDeleteButtonsColorFactura, setEditButtonsColorFactura, setFileSelectButtonsColorFactura,
                         handleSaveButtonColors }) {
  
  const isAdmin = panelActivo === 'admin';
  
  return (
    <div className="seccion-content">
      <div className="info-card">
        <div className="card-header">
          <h2>🔘 Personalización de Botones</h2>
          <p>Personaliza los colores de los botones del {isAdmin ? 'Panel Administrativo' : 'Panel de Facturación'}</p>
        </div>
        
        <div className="button-color-grid">
          <div className="button-color-item">
            <label>Botones de Acción:</label>
            <input 
              type="color" 
              value={isAdmin ? actionButtonsColor : actionButtonsColorFactura}
              onChange={(e) => isAdmin ? setActionButtonsColor(e.target.value) : setActionButtonsColorFactura(e.target.value)}
            />
            <button 
              style={{ backgroundColor: isAdmin ? actionButtonsColor : actionButtonsColorFactura }}
              className="button-preview"
            >
              Guardar
            </button>
          </div>
          
          <div className="button-color-item">
            <label>Botones de Eliminar:</label>
            <input 
              type="color" 
              value={isAdmin ? deleteButtonsColor : deleteButtonsColorFactura}
              onChange={(e) => isAdmin ? setDeleteButtonsColor(e.target.value) : setDeleteButtonsColorFactura(e.target.value)}
            />
            <button 
              style={{ backgroundColor: isAdmin ? deleteButtonsColor : deleteButtonsColorFactura }}
              className="button-preview"
            >
              Eliminar
            </button>
          </div>
          
          <div className="button-color-item">
            <label>Botones de Editar:</label>
            <input 
              type="color" 
              value={isAdmin ? editButtonsColor : editButtonsColorFactura}
              onChange={(e) => isAdmin ? setEditButtonsColor(e.target.value) : setEditButtonsColorFactura(e.target.value)}
            />
            <button 
              style={{ backgroundColor: isAdmin ? editButtonsColor : editButtonsColorFactura }}
              className="button-preview"
            >
              Editar
            </button>
          </div>
          
          <div className="button-color-item">
            <label>Botones de Archivo:</label>
            <input 
              type="color" 
              value={isAdmin ? fileSelectButtonsColor : fileSelectButtonsColorFactura}
              onChange={(e) => isAdmin ? setFileSelectButtonsColor(e.target.value) : setFileSelectButtonsColorFactura(e.target.value)}
            />
            <button 
              style={{ backgroundColor: isAdmin ? fileSelectButtonsColor : fileSelectButtonsColorFactura }}
              className="button-preview"
            >
              Seleccionar
            </button>
          </div>
        </div>
        
        <button onClick={handleSaveButtonColors} className="save-button">
          Guardar Colores de Botones
        </button>
      </div>
    </div>
  );
}

// Componente para la sección de Empresa
function SeccionEmpresa({ logoImage, logoPreview, handleLogoUpload, handleLogoSave, localCompanyName, 
                        setLocalCompanyName, handleCompanySave, navbarBgColor, updateNavbarBgColor,
                        companyTextColor, updateCompanyTextColor }) {
  return (
    <div className="seccion-content">
      {/* Logo */}
      <div className="info-card">
        <div className="card-header">
          <h2>🖼️ Logo de la Empresa</h2>
          <p>Personaliza el logo que aparecerá en la aplicación</p>
        </div>
        
        <div className="logo-options">
          <div className="logo-upload">
            <label htmlFor="logo-upload">Seleccionar Logo:</label>
            <div className="file-upload-container">
              <input
                type="file"
                id="logo-upload"
                accept="image/*"
                onChange={handleLogoUpload}
                style={{ display: 'none' }}
              />
              <label htmlFor="logo-upload" className="file-upload-button">
                Elegir Archivo
              </label>
              {logoImage && <span className="file-name">{logoImage.name}</span>}
            </div>
          </div>
          
          {logoPreview && (
            <div className="logo-preview">
              <h4>Vista Previa:</h4>
              <img src={logoPreview} alt="Logo preview" className="logo-preview-img" />
            </div>
          )}
          
          <button onClick={handleLogoSave} className="save-button">
            Guardar Logo
          </button>
        </div>
      </div>
      
      {/* Información de Empresa */}
      <div className="info-card">
        <div className="card-header">
          <h2>🏢 Información de la Empresa</h2>
          <p>Configura el nombre y colores de la empresa</p>
        </div>
        
        <div className="company-options">
          <div className="info-group">
            <label htmlFor="company-name">Nombre de la Empresa:</label>
            <input
              type="text"
              id="company-name"
              value={localCompanyName}
              onChange={(e) => setLocalCompanyName(e.target.value)}
              placeholder="Ingresa el nombre de tu empresa"
            />
          </div>
          
          <div className="color-options-grid">
            <div className="info-group">
              <label htmlFor="navbar-bg-color">Color de Fondo del Navbar:</label>
              <input
                type="color"
                id="navbar-bg-color"
                value={navbarBgColor}
                onChange={(e) => updateNavbarBgColor(e.target.value)}
              />
            </div>
            
            <div className="info-group">
              <label htmlFor="company-text-color">Color del Texto de Empresa:</label>
              <input
                type="color"
                id="company-text-color"
                value={companyTextColor}
                onChange={(e) => updateCompanyTextColor(e.target.value)}
              />
            </div>
          </div>
          
          <button onClick={handleCompanySave} className="save-button">
            Guardar Información
          </button>
        </div>
      </div>
    </div>
  );
}

export { SubseccionNav, SeccionTemas, SeccionBotones, SeccionEmpresa };
