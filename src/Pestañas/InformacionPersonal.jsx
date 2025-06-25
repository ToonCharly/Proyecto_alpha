import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import '../STYLES/Informacion_Personal.css'; 

const InformacionPersonal = () => {
  const navigate = useNavigate();
  const [userData, setUserData] = useState(null);
  const [isEditing, setIsEditing] = useState(false);
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    phone: '',
  });
  const [modalErrores, setModalErrores] = useState([]);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [fieldErrors, setFieldErrors] = useState({});

  useEffect(() => {
    // CAMBIO: Usar sessionStorage en lugar de localStorage para evitar conflictos entre ventanas
    const storedUserData = sessionStorage.getItem('userData');
    if (!storedUserData) {
      navigate('/login');
      return;
    }

    try {
      const parsedUserData = JSON.parse(storedUserData);
      setUserData(parsedUserData);
      setFormData({
        username: parsedUserData.username || '',
        email: parsedUserData.email || '',
        phone: parsedUserData.phone || '(Sin teléfono registrado)',
      });
      fetchUserDetails(parsedUserData.email || parsedUserData.username);
    } catch (error) {
      console.error("Error al procesar los datos del usuario:", error);
      navigate('/login');
    }
  }, [navigate]);

  const fetchUserDetails = async (identifier) => {
    try {
      const response = await fetch(`http://localhost:8080/api/usuario?identifier=${encodeURIComponent(identifier)}`);
      if (!response.ok) {
        throw new Error("No se pudo obtener la información del usuario");
      }
      const userDetails = await response.json();
      setFormData(prevFormData => ({
        username: userDetails.username || prevFormData.username,
        email: userDetails.email || prevFormData.email,
        phone: userDetails.phone || prevFormData.phone,
      }));
    } catch (error) {
      console.error("Error al obtener detalles del usuario:", error);
    }
  };

  const ModalError = ({ errores, onClose }) => {
    if (!errores || errores.length === 0) return null;
    return (
      <div className="modal-overlay">
        <div className="modal-contenido">
          <h3>Por favor corrige los siguientes errores:</h3>
          <ul className="lista-errores">
            {errores.map((error, index) => (
              <li key={index}>
                <span className="punto-error">•</span> {error}
              </li>
            ))}
          </ul>
          <button className="modal-boton" onClick={onClose}>
            Cerrar
          </button>
        </div>
      </div>
    );
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));

    if (fieldErrors[name]) {
      setFieldErrors(prev => {
        const newErrors = { ...prev };
        delete newErrors[name];
        return newErrors;
      });
    }

    if (modalErrores.length > 0) {
      setModalErrores([]);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsSubmitting(true);
    setModalErrores([]);
    setFieldErrors({});

    const errores = [];
    const campoErrores = {};

    if (!formData.username || !formData.username.trim()) {
      errores.push('El nombre de usuario es obligatorio');
      campoErrores.username = true;
    }

    if (!formData.email || !formData.email.trim()) {
      errores.push('El correo electrónico es obligatorio');
      campoErrores.email = true;
    } else if (!/\S+@\S+\.\S+/.test(formData.email.trim())) {
      errores.push('El formato del correo electrónico no es válido');
      campoErrores.email = true;
    }

    if (errores.length > 0) {
      setFieldErrors(campoErrores);
      setModalErrores(errores);
      setIsSubmitting(false);
      return;
    }

    try {
      const updateData = {
        email: formData.email.trim(),
        username: formData.username.trim(),
        phone: formData.phone,
      };

      const response = await fetch('http://localhost:8080/api/actualizar_usuario', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(updateData),
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || "Error al actualizar los datos");
      }

      // Actualizar localStorage con los nuevos datos
      const updatedUserData = {
        ...userData,
        username: formData.username,
        email: formData.email,
        phone: formData.phone,
      };

      localStorage.setItem('userData', JSON.stringify(updatedUserData));
      setUserData(updatedUserData);
      setIsEditing(false);
      alert("Datos actualizados correctamente");

      // Disparar evento para notificar al componente Home
      const profileEvent = new CustomEvent('profileUpdated');
      window.dispatchEvent(profileEvent);
    } catch (error) {
      console.error("Error al actualizar la información:", error);
      if (error.message.includes("Failed to fetch") || error.message.includes("Network Error")) {
        const updatedUserData = {
          ...userData,
          username: formData.username,
          email: formData.email,
          phone: formData.phone,
        };
        localStorage.setItem('userData', JSON.stringify(updatedUserData));
        setUserData(updatedUserData);
        setIsEditing(false);
        alert("Datos guardados localmente (sin conexión al servidor)");
      } else {
        setModalErrores([`Error al actualizar: ${error.message}`]);
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  if (!userData) {
    return <div className="loading">Cargando información...</div>;
  }

  return (
    <div className="info-personal-container" style={{ marginTop: '60px', marginLeft: '290px' }}>
      <ModalError 
        errores={modalErrores} 
        onClose={() => setModalErrores([])} 
      />

      <h1 className="titulo">Empresas Registradas</h1>

      <div className="info-card">
        <div className="card-header">
          <h2>Mis Datos</h2>
          {!isEditing ? (
            <button 
              className="btn-editar" 
              onClick={() => setIsEditing(true)}
            >
              Editar Información
            </button>
          ) : (
            <button 
              className="btn-cancelar" 
              onClick={() => {
                setIsEditing(false);
                const storedData = JSON.parse(sessionStorage.getItem('userData')); // Migrado a sessionStorage
                setFormData({
                  username: storedData.username || '',
                  email: storedData.email || '',
                  phone: storedData.phone || '(Sin teléfono registrado)',
                });
                setFieldErrors({});
              }}
            >
              Cancelar
            </button>
          )}
        </div>

        <form onSubmit={handleSubmit}>
          <div className="info-grid">
            <div className="info-group">
              <label htmlFor="username">Nombre de Usuario: <span className="required">*</span></label>
              <input
                type="text"
                id="username"
                name="username"
                value={formData.username}
                onChange={handleChange}
                className={fieldErrors.username ? 'input-error' : ''}
                disabled={!isEditing}
                required
              />
            </div>

            <div className="info-group">
              <label htmlFor="email">Correo Electrónico: <span className="required">*</span></label>
              <input
                type="email"
                id="email"
                name="email"
                value={formData.email}
                onChange={handleChange}
                className={fieldErrors.email ? 'input-error' : ''}
                disabled={!isEditing}
                required
              />
            </div>

            <div className="info-group">
              <label htmlFor="phone">Teléfono:</label>
              {isEditing ? (
                <input
                  type="tel"
                  id="phone"
                  name="phone"
                  value={formData.phone}
                  onChange={handleChange}
                  className={fieldErrors.phone ? 'input-error' : ''}
                  placeholder="Ingresa tu número de teléfono"
                />
              ) : (
                <div className="read-only-field">
                  {formData.phone || '(Sin teléfono registrado)'}
                </div>
              )}
            </div>
          </div>

          {isEditing && (
            <div className="form-actions">
              <button 
                type="submit" 
                className="btn-guardar"
                disabled={isSubmitting}
              >
                {isSubmitting ? 'Guardando...' : 'Guardar Cambios'}
              </button>
            </div>
          )}
        </form>
      </div>
    </div>
  );
};

export default InformacionPersonal;