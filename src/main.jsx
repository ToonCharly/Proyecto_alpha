import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import './App.css';
import { FacturaProvider } from './context/FacturaContext'; // Importa el proveedor del contexto

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <FacturaProvider> {/* Envuelve la aplicación con el proveedor */}
      <App />
    </FacturaProvider>
  </React.StrictMode>
);