import React, { createContext, useState } from 'react';

export const FacturaContext = createContext();

export const FacturaProvider = ({ children }) => {
  const [historialFacturas, setHistorialFacturas] = useState([]);

  const agregarFactura = (factura) => {
    setHistorialFacturas((prev) => [...prev, factura]);
  };

  return (
    <FacturaContext.Provider value={{ historialFacturas, agregarFactura }}>
      {children}
    </FacturaContext.Provider>
  );
};