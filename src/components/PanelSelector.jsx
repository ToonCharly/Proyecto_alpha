import React, { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import '../STYLES/PanelSelector.css';

const PanelSelector = ({ currentPanel }) => {
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef(null);
  const navigate = useNavigate();

  // Cerrar el menú cuando se hace clic fuera de él
  useEffect(() => {
    const handleClickOutside = (event) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setIsOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  // Cambiar de panel
  const switchPanel = (panelType) => {
    if (panelType === 'admin') {
      navigate('/admin');
    } else {
      navigate('/home');
    }
    setIsOpen(false);
  };

  return (
    <div className="panel-selector-container" ref={dropdownRef}>
      <button 
        className="panel-selector-button"
        onClick={() => setIsOpen(!isOpen)}
      >
        <span>{currentPanel === 'admin' ? 'Panel Administrativo' : 'Panel de Facturación'}</span>
        <i className={`panel-selector-icon ${isOpen ? 'open' : ''}`}>▼</i>
      </button>
      
      {isOpen && (
        <div className="panel-selector-dropdown">
          <div 
            className={`panel-option ${currentPanel === 'admin' ? 'active' : ''}`}
            onClick={() => switchPanel('admin')}
          >
            <i className="panel-icon admin-icon">🔧</i>
            <span>Panel Administrativo</span>
          </div>
          <div 
            className={`panel-option ${currentPanel === 'facturacion' ? 'active' : ''}`}
            onClick={() => switchPanel('facturacion')}
          >
            <i className="panel-icon facturacion-icon">📋</i>
            <span>Panel de Facturación</span>
          </div>
        </div>
      )}
    </div>
  );
};

export default PanelSelector;