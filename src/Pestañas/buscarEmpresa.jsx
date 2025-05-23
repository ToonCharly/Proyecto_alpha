import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom'; // Importar useNavigate

function InicioFacturacion() {
  const navigate = useNavigate(); // Hook para redirigir a otra página
  const [criterioBusqueda, setCriterioBusqueda] = useState('');
  const [empresa, setEmpresa] = useState(null); // Ahora se utiliza para mostrar datos de la empresa
  const [error, setError] = useState(null);
  const [cargando, setCargando] = useState(false);
  const [formData, setFormData] = useState({
    rfc: '',
    razonSocial: '',
  });
  const [datosPrincipales, setDatosPrincipales] = useState(false); // Renombrado para evitar confusión
  const [datosAdicionales, setDatosAdicionales] = useState(false); // Renombrado para evitar confusión

 const buscarEmpresa = async (e) => {
  e.preventDefault();
  setError(null);
  setCargando(true);

  try {
    const response = await fetch(`http://localhost:8080/api/factura?criterio=${encodeURIComponent(criterioBusqueda)}`);
    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.error || 'Empresa no encontrada');
    }

    setEmpresa(data); // Se utiliza para mostrar los datos de la empresa
    setFormData((prevFormData) => ({
      ...prevFormData,
      rfc: data.rfc || criterioBusqueda.toUpperCase(),
      razonSocial: data.razon_social || '',
    }));
    setDatosPrincipales(true); // Mostrar datos principales
    setDatosAdicionales(true); // Mostrar datos adicionales
  } catch (error) {
    setError(`Empresa no encontrada: ${error.message}. Redirigiendo a Administrar Empresas...`);
    setEmpresa(null);

    // Redirigir a la página de Administrar Empresas después de 3 segundos
    setTimeout(() => {
      navigate('/administrar-empresas'); // Cambia la ruta según tu configuración
    }, 3000);
  } finally {
    setCargando(false);
  }
};

  return (
    <div>
      <form onSubmit={buscarEmpresa}>
        <input
          type="text"
          value={criterioBusqueda}
          onChange={(e) => setCriterioBusqueda(e.target.value)}
          placeholder="Ingresa RFC o Razón Social"
        />
        <button type="submit" disabled={cargando}>
          {cargando ? 'Buscando...' : 'Buscar'}
        </button>
      </form>
      {error && <p>{error}</p>}

      {empresa && datosPrincipales && (
        <div>
          <h2>Datos Principales</h2>
          <p>RFC: {formData.rfc}</p>
          <p>Razón Social: {formData.razonSocial}</p>
        </div>
      )}

      {empresa && datosAdicionales && (
        <div>
          <h2>Datos Adicionales</h2>
          <p>Más información sobre la empresa...</p>
        </div>
      )}
    </div>
  );
}

export default InicioFacturacion;