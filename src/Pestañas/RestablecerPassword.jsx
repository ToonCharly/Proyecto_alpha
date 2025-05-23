import React, { useState, useEffect } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import '../styles/Login.css';

const RestablecerPassword = () => {
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [token, setToken] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);
  const [cargando, setCargando] = useState(false);
  const [tokenInvalido, setTokenInvalido] = useState(false);
  
  const location = useLocation();
  const navigate = useNavigate();
  
  // Extraer token de la URL al cargar la página
  useEffect(() => {
    const searchParams = new URLSearchParams(location.search);
    const tokenFromUrl = searchParams.get('token');
    
    if (!tokenFromUrl) {
      setTokenInvalido(true);
      setError('El enlace es inválido. Por favor solicita un nuevo enlace de recuperación.');
      return;
    }
    
    setToken(tokenFromUrl);
  }, [location]);
  
  const handleSubmit = async (e) => {
    e.preventDefault();
    
    // Validaciones
    if (password.length !== 5) { // Cambiar aquí
      setError('La contraseña debe tener exactamente 5 caracteres');
      return;
    }
    
    if (password !== confirmPassword) {
      setError('Las contraseñas no coinciden');
      return;
    }
    
    setCargando(true);
    setError('');
    
    try {
      const response = await fetch('http://localhost:8080/api/reset-password', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          token: token,
          newPassword: password
        }),
      });
      
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Error al restablecer contraseña');
      }
      
      setSuccess(true);
      
      // Redireccionar al login después de 3 segundos
      setTimeout(() => {
        navigate('/login');
      }, 3000);
      
    } catch (err) {
      setError(err.message || 'Ocurrió un error al restablecer tu contraseña');
      console.error('Error:', err);
    } finally {
      setCargando(false);
    }
  };
  
  // Si el token es inválido, mostrar mensaje de error
  if (tokenInvalido) {
    return (
      <div className="login-container">
        <div className="login-form">
          <div className="logo-container">
            <div className="logo-circle"></div>
            <h2 className="logo-text">FACTS</h2>
          </div>
          
          <div className="login-form-content">
            <div className="error-message">{error}</div>
            <Link to="/recuperar-password" className="login-button">
              Solicitar nuevo enlace
            </Link>
          </div>
        </div>
      </div>
    );
  }
  
  return (
    <div className="login-container">
      <div className="login-form">
        <div className="logo-container">
          <div className="logo-circle"></div>
          <h2 className="logo-text">FACTS</h2>
        </div>
        
        <div className="login-form-content">
          {!success ? (
            <>
              <h2 className="sign-in-title">Establecer nueva contraseña</h2>
              <p className="login-description">
                Ingresa tu nueva contraseña para tu cuenta.
              </p>
              
              {error && <div className="error-message">{error}</div>}
              
              <form onSubmit={handleSubmit}>
                <div className="form-group">
                  <label htmlFor="password">Nueva contraseña</label>
                  <input
                    type="password"
                    id="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder="Ingresa tu nueva contraseña"
                    required
                  />
                </div>
                
                <div className="form-group">
                  <label htmlFor="confirmPassword">Confirmar contraseña</label>
                  <input
                    type="password"
                    id="confirmPassword"
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    placeholder="Confirma tu nueva contraseña"
                    required
                  />
                </div>
                
                <button 
                  type="submit" 
                  className="login-button"
                  disabled={cargando}
                >
                  {cargando ? 'Actualizando...' : 'Guardar nueva contraseña'}
                </button>
              </form>
            </>
          ) : (
            <div className="success-message">
              <h2>¡Contraseña actualizada!</h2>
              <p>
                Tu contraseña ha sido actualizada correctamente.
              </p>
              <p>
                Serás redirigido al inicio de sesión en unos segundos...
              </p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default RestablecerPassword;