import React, { useState, useEffect } from 'react';
import '../styles/Empresas.css';

const Empresas = () => {
  const [showForm, setShowForm] = useState(false);
  const [empresas, setEmpresas] = useState([]);
  const [regimenesFiscales, setRegimenesFiscales] = useState([]);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [confirmDelete, setConfirmDelete] = useState(null);
  const [confirmEdit, setConfirmEdit] = useState(null);
  const [formData, setFormData] = useState({
    rfc: '',
    razon_social: '',
    regimen_fiscal: '',
    direccion: '',
    codigo_postal: '',
    pais: '',
    estado: '',
    localidad: '',
    municipio: '',
    colonia: '',
  });

  // Obtener usuario del objeto userData almacenado en localStorage
  const getUserId = () => {
    const userDataString = sessionStorage.getItem('userData'); // Migrado a sessionStorage para evitar conflictos
    if (!userDataString) return null;
    
    try {
      const userData = JSON.parse(userDataString);
      return userData.id; // Asegurarse que la propiedad sea 'id' y no otra
    } catch (e) {
      console.error("Error al parsear userData:", e);
      return null;
    }
  };
  
  const userId = getUserId();

  // Obtener empresas
  const fetchEmpresas = async () => {
    if (!userId) {
      setError('No hay un usuario activo. Por favor inicie sesión.');
      return;
    }

    try {
      const res = await fetch(`http://localhost:8080/api/empresas?id_usuario=${userId}`);
      if (!res.ok) {
        throw new Error('Error al obtener las empresas');
      }
      const data = await res.json();
      setEmpresas(data || []);
    } catch (error) {
      console.error('Error al cargar empresas:', error);
      setError('Hubo un problema al cargar las empresas. Por favor, intenta nuevamente.');
    }
  };

  // Obtener regímenes fiscales
  const fetchRegimenesFiscales = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/regimenes-fiscales');
      if (!response.ok) {
        throw new Error('Error al obtener los regímenes fiscales');
      }
      const data = await response.json();
      const filteredData = data.filter((regimen) => regimen.descripcion !== 'Sin regimen');
      setRegimenesFiscales(filteredData || []);
    } catch (error) {
      console.error('Error al cargar los regímenes fiscales:', error);
      setError('Hubo un problema al cargar los regímenes fiscales. Por favor, intenta nuevamente.');
    }
  };

  // Cargar datos al montar el componente
  useEffect(() => {
    fetchRegimenesFiscales();
    if (userId) {
      fetchEmpresas();
    }
  }, [userId]);

  // Limpiar mensajes después de 5 segundos
  useEffect(() => {
    if (success || error) {
      const timer = setTimeout(() => {
        setSuccess('');
        setError('');
      }, 5000);
      return () => clearTimeout(timer);
    }
  }, [success, error]);

  // Manejar cambios en el formulario
  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
  };

  // Modificado: Manejar envío del formulario para mostrar confirmación
  const handleSubmit = (e) => {
    e.preventDefault();
    
    if (!userId) {
      setError('No hay un usuario activo. Por favor inicie sesión.');
      return;
    }

    // Si es una edición, mostrar confirmación
    if (formData.id) {
      setConfirmEdit({
        data: formData,
        name: formData.razon_social
      });
    } else {
      // Si es nuevo registro, proceder directamente
      saveEmpresa(formData);
    }
  };

  // Nuevo: Función para guardar empresa (separada de handleSubmit)
  const saveEmpresa = async (empresaData) => {
    setError('');
    setSuccess('');
    setConfirmEdit(null); // Limpiar confirmación
    
    const dataToSend = {
      id_usuario: userId,
      ...empresaData,
    };

    try {
      // Verificar si el servidor está respondiendo
      try {
        const pingResponse = await fetch('http://localhost:8080/api/ping', {
          method: 'GET',
          headers: {
            'Accept': 'application/json',
          },
          signal: AbortSignal.timeout(3000)
        });
        
        if (!pingResponse.ok) {
          throw new Error('Servidor no responde correctamente');
        }
      } catch (pingError) {
        console.error('Error al verificar estado del servidor:', pingError);
        setError('El servidor no está respondiendo. Verifica que esté en ejecución.');
        return;
      }

      let response;
      
      // Si formData tiene un ID, estamos editando una empresa existente
      if (empresaData.id) {
        // Enviar solicitud PUT para actualizar
        response = await fetch(`http://localhost:8080/api/empresas/${empresaData.id}`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(dataToSend),
        });
      } else {
        // Enviar solicitud POST para crear nueva
        response = await fetch('http://localhost:8080/api/empresas', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(dataToSend),
        });
      }

      if (response.ok) {
        const result = await response.json();
        
        if (empresaData.id) {
          // Actualizar la empresa en el arreglo local
          setEmpresas(prevEmpresas => 
            prevEmpresas.map(empresa => 
              empresa.id === empresaData.id ? result : empresa
            )
          );
          setSuccess(`Empresa "${result.razon_social}" actualizada correctamente`);
        } else {
          // Agregar la nueva empresa al arreglo
          setEmpresas(prevEmpresas => [...prevEmpresas, result]);
          setSuccess(`Empresa "${result.razon_social}" registrada correctamente`);
        }
        
        // Limpiar el formulario y cerrarlo
        setShowForm(false);
        setFormData({
          rfc: '',
          razon_social: '',
          regimen_fiscal: '',
          direccion: '',
          codigo_postal: '',
          pais: '',
          estado: '',
          localidad: '',
          municipio: '',
          colonia: '',
        });
      } else {
        const errorText = await response.text();
        console.error('Error al guardar la empresa:', errorText);
        setError(`Error al ${empresaData.id ? 'actualizar' : 'guardar'} la empresa: ${errorText}`);
      }
    } catch (error) {
      console.error('Error al conectar con el backend:', error);
      setError('Error al conectar con el servidor. Por favor, intenta nuevamente.');
    }
  };

  // Confirmación antes de desvincular
  const handleConfirmDelete = (id) => {
    const empresa = empresas.find(emp => emp.id === id);
    if (empresa) {
      setConfirmDelete({
        id: id,
        name: empresa.razon_social
      });
    }
  };

  // Cancelar confirmación de desvinculación
  const handleCancelDelete = () => {
    setConfirmDelete(null);
  };

  // Cancelar confirmación de edición
  const handleCancelEdit = () => {
    setConfirmEdit(null);
  };

  // Manejar desvinculación de una empresa
  const handleDesvincular = async (id) => {
    setError('');
    setSuccess('');
    setConfirmDelete(null);
    
    try {
      // Verificar si el servidor está respondiendo
      try {
        const pingResponse = await fetch('http://localhost:8080/api/ping', {
          method: 'GET',
          headers: {
            'Accept': 'application/json',
          },
          // Timeout corto para verificar que el servidor esté vivo
          signal: AbortSignal.timeout(3000)
        });
        
        if (!pingResponse.ok) {
          throw new Error('Servidor no responde correctamente');
        }
      } catch (pingError) {
        console.error('Error al verificar estado del servidor:', pingError);
        setError('El servidor no está respondiendo. Verifica que esté en ejecución.');
        return;
      }
      
      // Eliminar la empresa
      const response = await fetch(`http://localhost:8080/api/empresas/${id}`, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
        }
      });
      
      if (response.ok) {
        // Obtener el nombre antes de eliminar para el mensaje
        const empresaName = empresas.find(emp => emp.id === id)?.razon_social || '';
        
        // Actualizar la lista de empresas sin recargar
        setEmpresas(prevEmpresas => prevEmpresas.filter(empresa => empresa.id !== id));
        
        // Mensaje de éxito
        setSuccess(`Empresa "${empresaName}" desvinculada correctamente`);
      } else {
        const errorData = await response.text();
        console.error('Error de respuesta del servidor:', errorData);
        setError(`Error al desvincular empresa: ${errorData || 'El servidor rechazó la solicitud'}`);
      }
    } catch (error) {
      console.error('Error en la operación de desvinculación:', error);
      setError(`Error de conexión: ${error.message}`);
    }
  };

  // Manejar edición de una empresa
  const handleEditar = (id) => {
    const empresa = empresas.find((empresa) => empresa.id === id);
    if (empresa) {
      // Importante: asegúrate de incluir todos los campos, especialmente el ID
      setFormData({
        ...empresa
      });
      setShowForm(true);
    }
  };

  // Modal de confirmación para desvinculación - Cambio de color del botón Cancelar
  const ConfirmationModal = () => {
    if (!confirmDelete) return null;
    
    return (
      <div className="modal-overlay">
        <div className="modal-content" style={{
          backgroundColor: 'white',
          padding: '20px',
          borderRadius: '8px',
          boxShadow: '0 4px 8px rgba(0, 0, 0, 0.2)',
          maxWidth: '500px',
          width: '90%'
        }}>
          <h3>Confirmar Eliminación</h3>
          <p>
            ¿Estás seguro de que deseas eliminar la empresa <strong>{confirmDelete.name}</strong>?
            Esta acción no podrá deshacerse.
          </p>
          <div style={{
            display: 'flex',
            justifyContent: 'center',
            gap: '40px',
            marginTop: '25px'
          }}>
            <button 
              type="button" 
              className="close-modal"
              onClick={handleCancelDelete}
              style={{
                padding: '10px 20px',
                minWidth: '120px',
                borderRadius: '5px',
                backgroundColor: '#2c3e50', /* Restored original blue */
                color: 'white',
                border: 'none'
              }}
            >
              Cancelar
            </button>
            <button 
              type="button" 
              onClick={() => handleDesvincular(confirmDelete.id)}
              style={{
                padding: '10px 20px',
                minWidth: '120px',
                borderRadius: '5px',
                backgroundColor: '#a93226', /* Dark red */
                color: 'white',
                border: 'none'
              }}
            >
              Confirmar
            </button>
          </div>
        </div>
      </div>
    );
  };

  // Modal de confirmación para edición - Cambio de color del botón Cancelar
  const EditConfirmationModal = () => {
    if (!confirmEdit) return null;
    
    return (
      <div className="modal-overlay" style={{ zIndex: 1100 }}>
        <div className="modal-content" style={{
          backgroundColor: 'white',
          padding: '20px',
          borderRadius: '8px',
          boxShadow: '0 4px 8px rgba(0, 0, 0, 0.2)',
          maxWidth: '500px',
          width: '90%',
          position: 'relative',
          zIndex: 1101
        }}>
          <h3>Confirmar Cambios</h3>
          <p>
            ¿Estás seguro de guardar los cambios para la empresa <strong>{confirmEdit.name}</strong>?
          </p>
          <div style={{
            display: 'flex',
            justifyContent: 'center',
            gap: '40px',
            marginTop: '25px'
          }}>
            <button 
              type="button" 
              className="close-modal"
              onClick={handleCancelEdit}
              style={{
                padding: '10px 20px',
                minWidth: '120px',
                borderRadius: '5px',
                backgroundColor: '#2c3e50', /* Restored original blue */
                color: 'white',
                border: 'none'
              }}
            >
              Cancelar
            </button>
            <button 
              type="button" 
              onClick={() => saveEmpresa(confirmEdit.data)}
              style={{
                padding: '10px 20px',
                minWidth: '120px',
                borderRadius: '5px',
                backgroundColor: '#196f3d', /* Dark green */
                color: 'white',
                border: 'none'
              }}
            >
              Confirmar
            </button>
          </div>
        </div>
      </div>
    );
  };

  return (
    <div className="empresas-container" style={{ marginTop: '60px', marginLeft: '230px' }}>
      
      {error && <div className="error-message">{error}</div>}
      {success && <div className="success-message">{success}</div>}
      <ConfirmationModal />
      <EditConfirmationModal />

      {showForm && (
        <div className="modal-overlay">
          <div className="modal-content">
            <form onSubmit={handleSubmit} className="empresa-form">
              <h2>{formData.id ? 'Editar Empresa' : 'Registrar Nueva Empresa'}</h2>
              <label>
                RFC:
                <input
                  type="text"
                  name="rfc"
                  value={formData.rfc}
                  onChange={handleChange}
                  required
                />
              </label>
              <label>
                Razón Social:
                <input
                  type="text"
                  name="razon_social"
                  value={formData.razon_social}
                  onChange={handleChange}
                  required
                />
              </label>
              <label>
                Régimen Fiscal:
                <select
                  name="regimen_fiscal"
                  value={formData.regimen_fiscal}
                  onChange={handleChange}
                  required
                >
                  <option value="">Seleccione un régimen fiscal</option>
                  {regimenesFiscales.map((regimen) => (
                    <option key={regimen.id} value={regimen.id}>
                      {regimen.codigo} - {regimen.descripcion}
                    </option>
                  ))}
                </select>
              </label>
              <label>
                Dirección:
                <input
                  type="text"
                  name="direccion"
                  value={formData.direccion}
                  onChange={handleChange}
                  required
                />
              </label>
              <label>
                Código Postal:
                <input
                  type="text"
                  name="codigo_postal"
                  value={formData.codigo_postal}
                  onChange={handleChange}
                  required
                />
              </label>
              <label>
                País:
                <input
                  type="text"
                  name="pais"
                  value={formData.pais}
                  onChange={handleChange}
                  required
                />
              </label>
              <label>
                Estado:
                <input
                  type="text"
                  name="estado"
                  value={formData.estado}
                  onChange={handleChange}
                  required
                />
              </label>
              <label>
                Localidad:
                <input
                  type="text"
                  name="localidad"
                  value={formData.localidad}
                  onChange={handleChange}
                  required
                />
              </label>
              <label>
                Municipio:
                <input
                  type="text"
                  name="municipio"
                  value={formData.municipio}
                  onChange={handleChange}
                  required
                />
              </label>
              <label>
                Colonia:
                <input
                  type="text"
                  name="colonia"
                  value={formData.colonia}
                  onChange={handleChange}
                  required
                />
              </label>
              <button type="submit">{formData.id ? 'Guardar Cambios' : 'Registrar Empresa'}</button>
              <button
                type="button"
                className="close-modal"
                onClick={() => setShowForm(false)}
              >
                Cancelar
              </button>
            </form>
          </div>
        </div>
      )}

      <h1 className="titulo">Empresas Registradas</h1>

      <div className="table-card">
        <div style={{ textAlign: 'right', marginBottom: '15px' }}>
          <button onClick={() => {
            // Reset form data before showing the form for a new company
            setFormData({
              rfc: '',
              razon_social: '',
              regimen_fiscal: '',
              direccion: '',
              codigo_postal: '',
              pais: '',
              estado: '',
              localidad: '',
              municipio: '',
              colonia: '',
            });
            setShowForm(!showForm);
          }}>
            {showForm ? 'Cancelar' : 'Agregar'}
          </button>
        </div>
        <table>
          <thead>
            <tr>
              <th>Razón Social</th>
              <th>Acciones</th>
            </tr>
          </thead>
          <tbody>
            {empresas && empresas.length > 0 ? (
              empresas.map((empresa) => (
                <tr key={empresa.id}>
                  <td>
                    <div className="empresa-info">
                      <div className="empresa-nombre">{empresa.rfc}</div>
                      <div className="empresa-detalles">{empresa.razon_social}</div>
                      <div className="empresa-detalles">{empresa.direccion}</div>
                    </div>
                  </td>
                  <td>
                    <button onClick={() => handleEditar(empresa.id)}>Editar</button>
                    <button onClick={() => handleConfirmDelete(empresa.id)}>Eliminar</button>
                  </td>
                </tr>
              ))
            ) : (
              <tr>
                <td colSpan="2">No hay empresas registradas.</td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default Empresas;