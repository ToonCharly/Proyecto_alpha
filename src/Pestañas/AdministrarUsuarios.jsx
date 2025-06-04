import React, { useState, useEffect } from 'react';
import '../STYLES/AdministrarUsuarios.css';

function AdministrarUsuarios() {
  const [usuarios, setUsuarios] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const [mensaje, setMensaje] = useState(null);
  const [currentUser, setCurrentUser] = useState(null);
  const [confirmAction, setConfirmAction] = useState(null);

  // Depurar al inicio del componente
  useEffect(() => {
    console.log("Componente AdministrarUsuarios iniciado");
    const storedUserData = localStorage.getItem('userData');
    console.log("Datos almacenados:", storedUserData);
    
    if (storedUserData) {
      try {
        const parsedUserData = JSON.parse(storedUserData);
        console.log("Datos parseados:", parsedUserData);
        
        setCurrentUser(parsedUserData);
        
        // Cargar usuarios sin pasar token
        fetchUsuarios();
      } catch (err) {
        console.error("Error al procesar datos del usuario:", err);
        setError("Error al cargar datos de usuario");
        setIsLoading(false);
      }
    } else {
      console.error("No hay datos de usuario en localStorage");
      setError("No has iniciado sesión o tus credenciales han expirado");
      setIsLoading(false);
    }
  }, []);

  // Función modificada para obtener la lista de usuarios
  const fetchUsuarios = async () => {
    try {
      setIsLoading(true);
      setError(null);
      
      // Hacer petición sin token por ahora
      const response = await fetch('http://localhost:8080/api/usuarios');
      
      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(`Error (${response.status}): ${errorText}`);
      }
      
      const data = await response.json();
      console.log("Datos recibidos:", data);
      
      // Verificar si los datos son un array
      if (!Array.isArray(data)) {
        console.error("Datos recibidos no son un array:", data);
        throw new Error("Formato de datos incorrecto");
      }
      
      setUsuarios(data);
    } catch (err) {
      console.error("Error completo:", err);
      setError('Error al cargar la lista de usuarios: ' + err.message);
    } finally {
      setIsLoading(false);
    }
  };

  // Función para cambiar el rol de un usuario
  const cambiarRolUsuario = async (userId, hacerAdmin) => {
    try {
      setIsLoading(true);
      
      const response = await fetch(`http://localhost:8080/api/usuarios/${userId}/rol`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${currentUser?.token}`
        },
        body: JSON.stringify({ isAdmin: hacerAdmin })
      });
      
      if (!response.ok) {
        throw new Error(`No se pudo ${hacerAdmin ? 'asignar' : 'remover'} el rol de administrador`);
      }
      
      // Actualizar la lista de usuarios localmente
      setUsuarios(usuarios.map(usuario => 
        usuario.id === userId ? { ...usuario, isAdmin: hacerAdmin } : usuario
      ));
      
      setMensaje({
        tipo: 'exito',
        texto: `Se ha ${hacerAdmin ? 'asignado' : 'removido'} correctamente el rol de administrador`
      });
      
      // Ocultar el mensaje después de 3 segundos
      setTimeout(() => {
        setMensaje(null);
      }, 3000);
      
    } catch (err) {
      console.error("Error al cambiar rol:", err);
      setMensaje({
        tipo: 'error',
        texto: err.message
      });
    } finally {
      setConfirmAction(null);
      setIsLoading(false);
    }
  };

  // Función para confirmar acción
  const mostrarConfirmacion = (usuario, hacerAdmin) => {
    setConfirmAction({
      usuario,
      hacerAdmin,
      mensaje: `¿Estás seguro de ${hacerAdmin ? 'hacer administrador' : 'quitar permisos de administrador'} a ${usuario.nombre || usuario.email}?`
    });
  };

  return (
    <div className="seccion-container">
      {/* Mensaje de éxito o error */}
      {mensaje && (
        <div className={`mensaje-banner ${mensaje.tipo}`}>
          {mensaje.texto}
          <button 
            onClick={() => setMensaje(null)} 
            className="cerrar-mensaje"
          >
            ×
          </button>
        </div>
      )}
      
      {/* Modal de confirmación */}
      {confirmAction && (
        <div className="modal-confirmacion">
          <div className="modal-contenido">
            <h3>Confirmar acción</h3>
            <p>{confirmAction.mensaje}</p>
            <div className="modal-acciones">
              <button 
                className="btn-cancelar" 
                onClick={() => setConfirmAction(null)}
              >
                Cancelar
              </button>
              <button 
                className={confirmAction.hacerAdmin ? "btn-confirmar" : "btn-peligro"}
                onClick={() => cambiarRolUsuario(confirmAction.usuario.id, confirmAction.hacerAdmin)}
              >
                Confirmar
              </button>
            </div>
          </div>
        </div>
      )}

      <div className="seccion-header">
        <h2>Administración de Usuarios</h2>
        <button 
          className="btn-refresh" 
          onClick={() => fetchUsuarios()} // Sin pasar currentUser
          disabled={isLoading}
        >
          <i className="fas fa-sync-alt"></i> Actualizar
        </button>
      </div>
      
      {error && <div className="error-message">{error}</div>}
      
      <div className="seccion-content">
        {isLoading && usuarios.length === 0 ? (
          <div className="loading-spinner">
            <div className="spinner"></div>
            <p>Cargando usuarios...</p>
          </div>
        ) : (
          <div className="tabla-usuarios-container">
            <table className="tabla-usuarios">
              <thead>
                <tr>
                  <th>Nombre</th>
                  <th>Correo Electrónico</th>
                  <th>Rol</th>
                  <th>Acciones</th>
                </tr>
              </thead>
              <tbody>
                {usuarios.length === 0 ? (
                  <tr>
                    <td colSpan="4" className="no-usuarios">No hay usuarios registrados</td>
                  </tr>
                ) : (
                  usuarios.map(usuario => (
                    <tr key={usuario.id} className={usuario.isAdmin ? 'fila-admin' : ''}>
                      <td>{usuario.nombre || '(Sin nombre)'}</td>
                      <td>{usuario.email}</td>
                      <td>
                        <span className={`badge ${usuario.isAdmin ? 'badge-admin' : 'badge-usuario'}`}>
                          {usuario.isAdmin ? 'Administrador' : 'Usuario'}
                        </span>
                      </td>
                      <td>
                        {currentUser?.id !== usuario.id && (
                          <button 
                            className={usuario.isAdmin ? "btn-quitar-admin" : "btn-hacer-admin"}
                            onClick={() => mostrarConfirmacion(usuario, !usuario.isAdmin)}
                            disabled={isLoading}
                          >
                            {usuario.isAdmin ? 'Quitar Admin' : 'Hacer Admin'}
                          </button>
                        )}
                        {currentUser?.id === usuario.id && (
                          <span className="usuario-actual">Usuario actual</span>
                        )}
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}

export default AdministrarUsuarios;