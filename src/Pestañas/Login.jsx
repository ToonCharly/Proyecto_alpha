import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import '../styles/Login.css';

const Form = () => {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    identifier: '',
    password: ''
  });
  const [modalErrores, setModalErrores] = useState([]);
  const [fieldErrors, setFieldErrors] = useState({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [currentSlide, setCurrentSlide] = useState(0);
  const [passwordVisible, setPasswordVisible] = useState(false);

  const slides = [
    "/gettyimages-2007594416-612x612.jpg",
    "/gettyimages-565975253-612x612.jpg",
    "/gettyimages-613241502-612x612.jpg",
  ];

  const nextSlide = () => {
    setCurrentSlide((prev) => (prev + 1) % slides.length);
  };

  useEffect(() => {
    // Solo configurar el timer si no hay errores mostrándose
    if (modalErrores.length === 0) {
      const timer = setInterval(() => {
        nextSlide();
      }, 6000);
      
      return () => clearInterval(timer);
    }
    // Añade modalErrores como dependencia para que el efecto se actualice cuando cambie
  }, [modalErrores]);

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

    if (name === 'password') {
      const limitedValue = value.slice(0, 8); // Cambiado de 5 a 8
      setFormData({ ...formData, [name]: limitedValue });
    } else {
      setFormData({ ...formData, [name]: value });
    }

    // Solo limpiar errores del campo específico que está siendo modificado
    if (fieldErrors[name]) {
      setFieldErrors((prev) => {
        const newErrors = { ...prev };
        delete newErrors[name];
        return newErrors;
      });
    }

    // Eliminar estas líneas que causan el reinicio del modal de errores
    // if (modalErrores.length > 0) {
    //   setModalErrores([]);
    // }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsSubmitting(true);
    setModalErrores([]);
    setFieldErrors({});

    const errores = [];
    const campoErrores = {};

    if (!formData.identifier || !formData.identifier.trim()) {
      errores.push('Debe ingresar un correo o nombre de usuario');
      campoErrores.identifier = true;
    }

    if (!formData.password) {
      errores.push('Debe ingresar una contraseña');
      campoErrores.password = true;
    } else if (formData.password.length !== 8) { // Cambiado de 5 a 8
      errores.push('La contraseña debe tener exactamente 8 caracteres'); // Mensaje actualizado
      campoErrores.password = true;
    }

    if (errores.length > 0) {
      setFieldErrors(campoErrores);
      setModalErrores(errores);
      setIsSubmitting(false);
      return;
    }

    try {
      const isEmail = formData.identifier.includes('@');
      console.log('Enviando datos al servidor:', {
        email: isEmail ? formData.identifier.trim() : '',
        username: !isEmail ? formData.identifier.trim() : '',
        password: formData.password,
      });

      const response = await fetch('http://localhost:8080/api/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: isEmail ? formData.identifier.trim() : '',
          username: !isEmail ? formData.identifier.trim() : '',
          password: formData.password,
        }),
      });

      console.log('Estado de la respuesta:', response.status);

      if (response.ok) {
        // Solo leer el cuerpo una vez
        const userData = await response.json();
        console.log('Datos del usuario:', userData);
        
        localStorage.setItem('userData', JSON.stringify(userData));
        
        // Redireccionar según el rol
        if (userData.role === 'admin') {
            navigate('/homeadmin');
        } else {
            navigate('/Home');
        }
      } else {
        // Para errores, conviene leer como texto
        const errorText = await response.text();
        console.error('Error en la respuesta del servidor:', errorText);
        setModalErrores(['Credenciales incorrectas']);
      }
    } catch (error) {
      console.error('Error al conectar con el servidor:', error);
      setModalErrores(['Error al conectar con el servidor']);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="login-page">
      <div className="login-container">
        <ModalError 
          errores={modalErrores} 
          onClose={() => setModalErrores([])} 
        />

        <div className="login-form-container">
          {/* Añadir el logo en la página de login */}
          <div className="login-logo"></div> {/* Este div mostrará el logo */}

          <div className="login-form-content">
            <h1 className="sign-in-title">Inicio de Sesion</h1>
            
            <div className="login-form">
              <div className="form-group">
                <label htmlFor="identifier" className="form-label">
                  Correo o Nombre de Usuario
                </label>
                <div className="input-container">
                  <input
                    id="identifier"
                    type="text"
                    name="identifier"
                    value={formData.identifier}
                    onChange={handleChange}
                    className={`form-input ${fieldErrors.identifier ? 'input-error' : ''}`}
                  />
                </div>
              </div>

              <div className="form-group">
                <label htmlFor="password" className="form-label">
                  Contraseña (8 caracteres) {/* Cambiado de 5 a 8 */}
                </label>
                <div className="input-container">
                  <input
                    id="password"
                    type={passwordVisible ? "text" : "password"}
                    name="password"
                    value={formData.password}
                    onChange={handleChange}
                    maxLength="8"
                    className={`form-input ${fieldErrors.password ? 'input-error' : ''}`}
                  />
                  <span 
                    className="password-toggle-icon"
                    onClick={() => setPasswordVisible(prev => !prev)}
                  >
                    <img 
                      src={passwordVisible ? "/visibilidad.png" : "/ojo.png"} 
                      alt={passwordVisible ? "Ocultar contraseña" : "Mostrar contraseña"} 
                      width="20" 
                      height="20" 
                    />
                  </span>
                </div>
              </div>

              <div className="forgot-password-container">
                <span 
                  className="forgot-password"
                  onClick={() => navigate('/recuperar-password')}  // CAMBIAR A ESTA RUTA
                >
                  Recuperar Contraseña
                </span>
              </div>

              <button
                onClick={handleSubmit}
                className="sign-in-button"
                disabled={isSubmitting}
              >
                {isSubmitting ? 'Cargando...' : 'Iniciar Sesión'}
              </button>
            </div>
          </div>

          <div className="signup-container">
            <p className="signup-text">
              No Tienes Una Cuenta? <span className="signup-link" onClick={() => navigate('/register')}>Registrate</span>
            </p>
          </div>
        </div>

        <div className="pattern-container">
          <div className="image-carousel">
            {slides.map((slide, index) => (
              <div 
                key={index} 
                className={`carousel-slide ${index === currentSlide ? 'active' : ''}`}
                style={{ backgroundImage: `url(${slide})` }}
              />
            ))}
          </div>

          <div className="pattern-lines">
            {[...Array(8)].map((_, i) => (
              <div
                key={i}
                className="pattern-line"
                style={{
                  width: `${40 + i * 10}%`,
                  height: `${40 + i * 10}%`,
                  transform: 'rotate(45deg)',
                }}
              ></div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default Form;