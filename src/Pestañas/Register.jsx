import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import '../styles/Register.css';

const RegisterForm = () => {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    phone: '',
    password: ''
  });

  const [modalErrores, setModalErrores] = useState([]);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [fieldErrors, setFieldErrors] = useState({});
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
    // Añade modalErrores como dependencia
  }, [modalErrores]);

  const handleChange = (e) => {
    const { name, value } = e.target;

    if (name === 'password') {
      const limitedValue = value.slice(0, 8); // Cambiado de 5 a 8
      setFormData({ ...formData, [name]: limitedValue });
    } else {
      setFormData({ ...formData, [name]: value });
    }

    // Solo limpiar errores del campo específico, no todos los errores
    if (fieldErrors[name]) {
      setFieldErrors((prev) => {
        const newErrors = { ...prev };
        delete newErrors[name];
        return newErrors;
      });
    }
  };

  // Componente ModalError con clases específicas
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

  const validatePhone = (phone) => {
    const cleanPhone = phone.replace(/\D/g, '');
    return cleanPhone.length >= 10;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    setIsSubmitting(true);
    setModalErrores([]);
    setFieldErrors({});

    const errores = [];
    const campoErrores = {};

    try {
      const userData = {
        username: formData.username.trim(),
        email: formData.email.trim(),
        phone: formData.phone.replace(/\D/g, ''),
        password: formData.password
      };  // Se agregó esta llave de cierre que faltaba

      if (!userData.username) {
        errores.push("Debe ingresar un nombre de usuario");
        campoErrores.username = true;
      }

      if (!userData.email) {
        errores.push("Debe ingresar un correo electrónico");
        campoErrores.email = true;
      }

      if (!userData.phone) {
        errores.push("Debe ingresar un número de teléfono");
        campoErrores.phone = true;
      }

      if (!userData.password) {
        errores.push("Debe ingresar una contraseña");
        campoErrores.password = true;
      }

      if (userData.password && userData.password.length !== 8) {
        errores.push("La contraseña debe tener exactamente 8 caracteres"); // Mensaje actualizado
        campoErrores.password = true;
      }

      if (userData.phone && !validatePhone(userData.phone)) {
        errores.push("El número de teléfono debe tener al menos 10 dígitos");
        campoErrores.phone = true;
      }

      if (errores.length > 0) {
        setFieldErrors(campoErrores);
        setModalErrores(errores);
        setIsSubmitting(false);
        return;
      }

      // First, make the registration request
      console.log("Sending registration data:", userData);
      const response = await fetch('http://localhost:8080/api/registrar_usuario', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(userData),
      });

      // Get the response data
      const responseText = await response.text();
      console.log("Raw server response:", responseText);
      
      // Try to parse the response
      let responseData = {};
      try {
        if (responseText) {
          responseData = JSON.parse(responseText);
          console.log("Parsed server response:", responseData);
        }
      } catch (parseError) {
        console.error("Failed to parse server response:", parseError);
      }

      if (!response.ok) {
        if (response.status === 409) {
          setModalErrores([responseData?.message || "Usuario o correo ya registrado."]);
        } else {
          setModalErrores([responseData?.message || responseText || "Error al registrar al usuario."]);
        }
        setIsSubmitting(false);
        return;
      }

      // After successful registration, we need to get the user's ID
      // Either from the registration response or by making a separate login request
      let userId = null;
      
      // Try to get ID from registration response first
      if (responseData.id) {
        userId = responseData.id;
      } else if (responseData.usuario && responseData.usuario.id) {
        userId = responseData.usuario.id;
      } else if (responseData.user && responseData.user.id) {
        userId = responseData.user.id;
      }
      
      // If we didn't get ID from registration, try to login to get it
      if (!userId) {
        console.log("User ID not found in registration response, attempting login");
        try {
          const loginResponse = await fetch('http://localhost:8080/api/login', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({
              email: userData.email,
              password: userData.password
            }),
          });
          
          if (loginResponse.ok) {
            const loginData = await loginResponse.json();
            console.log("Login response:", loginData);
            
            // Try to extract user ID from login response
            if (loginData.id) {
              userId = loginData.id;
            } else if (loginData.usuario && loginData.usuario.id) {
              userId = loginData.usuario.id;
            } else if (loginData.user && loginData.user.id) {
              userId = loginData.user.id;
            }
          } else {
            console.error("Failed to login after registration");
          }
        } catch (loginError) {
          console.error("Error during login after registration:", loginError);
        }
      }
      
      // Store the user data in localStorage with ID if available
      const userDataToStore = {
        username: userData.username,
        email: userData.email,
        phone: userData.phone
      };
      
      if (userId) {
        userDataToStore.id = userId;
      }
      
      console.log("Storing user data:", userDataToStore);
      localStorage.setItem('userData', JSON.stringify(userDataToStore));
      
      // Show success message and navigate
      alert('Usuario registrado exitosamente');
      
      // Add delay to ensure localStorage is updated
      setTimeout(() => {
        navigate('/Home');
      }, 100);
    } catch (error) {
      console.error('Error de conexión:', error);
      setModalErrores([`Error al conectar con el servidor: ${error.message}`]);
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
          <div className="login-form-content">
            <div className="logo-container" style={{ 
              display: 'flex', 
              justifyContent: 'center', 
              alignItems: 'center',
              marginBottom: '20px'
            }}>
              <div className="registro-logo"></div>
            </div>

            <h1 className="sign-in-title">Registro de Usuario</h1>
            
            <form className="login-form" onSubmit={handleSubmit} noValidate>
              <div className="form-group">
                <label htmlFor="username" className="form-label">
                  Nombre de Usuario
                </label>
                <div className="input-container">
                  <input
                    id="username"
                    type="text"
                    name="username"
                    value={formData.username}
                    onChange={handleChange}
                    className={`form-input ${fieldErrors.username ? 'input-error' : ''}`}
                  />
                </div>
              </div>

              <div className="form-group">
                <label htmlFor="email" className="form-label">
                  Correo Electrónico
                </label>
                <div className="input-container">
                  <input
                    id="email"
                    type="email"
                    name="email"
                    value={formData.email}
                    onChange={handleChange}
                    className={`form-input ${fieldErrors.email ? 'input-error' : ''}`}
                  />
                </div>
              </div>

              <div className="form-group">
                <label htmlFor="phone" className="form-label">
                  Número de Teléfono
                </label>
                <div className="input-container">
                  <input
                    id="phone"
                    type="tel"
                    name="phone"
                    value={formData.phone}
                    onChange={handleChange}
                    className={`form-input ${fieldErrors.phone ? 'input-error' : ''}`}
                  />
                </div>
              </div>

              <div className="form-group">
                <label htmlFor="password" className="form-label">
                  Contraseña (8 caracteres) 
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

              <button
                type="submit"
                className="sign-in-button"
                disabled={isSubmitting}
              >
                {isSubmitting ? 'Registrando...' : 'Registrarse'}
              </button>
            </form>
          </div>

          <div className="signup-container">
            <p className="signup-text">
              ¿Ya tienes una cuenta? <span className="signup-link" onClick={() => navigate('/login')}>Iniciar Sesión</span>
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
        </div>
      </div>
    </div>
  );
};

export default RegisterForm;