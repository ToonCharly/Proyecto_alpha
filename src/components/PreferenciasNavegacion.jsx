import React from 'react';
import '../STYLES/PreferenciasTabs.css';

// Componente principal para las subsecciones de Preferencias
const PreferenciasSubsecciones = ({ subseccionActiva, setSubseccionActiva }) => {  const subsecciones = [
    { 
      id: 'colores', 
      nombre: '🎨 Colores y Temas', 
      descripcion: 'Personalizar colores principales del sistema',
      color: '#e3f2fd'
    },
    { 
      id: 'botones', 
      nombre: '🔘 Botones', 
      descripcion: 'Estilos y colores de botones',
      color: '#f3e5f5'
    },
    { 
      id: 'tipografia', 
      nombre: '📝 Tipografía', 
      descripcion: 'Fuentes, tamaños y estilos de texto',
      color: '#e8f5e8'
    },
    { 
      id: 'empresa', 
      nombre: '🏢 Logo y Empresa', 
      descripcion: 'Gestionar logo e información empresarial',
      color: '#fff3e0'
    },
    { 
      id: 'plantillas', 
      nombre: '📄 Plantillas Word', 
      descripcion: 'Subir y gestionar plantillas de facturación',
      color: '#fce4ec'
    }
  ];

  return (
    <div className="subseccion-nav-container">
      <h3 className="subseccion-nav-title">📋 Configuraciones Disponibles:</h3>
      <div className="subseccion-nav-grid">
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
          >
            <div className="subseccion-content">
              <div className="subseccion-header">
                <h4 className="subseccion-nombre">{subseccion.nombre}</h4>
                {subseccionActiva === subseccion.id && (
                  <span className="check-icon">✓</span>
                )}
              </div>
              <p className="subseccion-descripcion">{subseccion.descripcion}</p>
            </div>
          </button>
        ))}
      </div>
    </div>
  );
};

export default PreferenciasSubsecciones;
