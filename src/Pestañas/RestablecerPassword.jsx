import React, { useState, useEffect } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import '../STYLES/RecoverPassword.css';

const RestablecerPassword = () => {
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [token, setToken] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);
  const [cargando, setCargando] = useState(false);
  const [tokenInvalido, setTokenInvalido] = useState(false);
  const [passwordsMatch, setPasswordsMatch] = useState(true);
  
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
  
  // Verificar coincidencia de contraseñas en tiempo real
  useEffect(() => {
    if (confirmPassword) {
      setPasswordsMatch(password === confirmPassword);
    } else {
      setPasswordsMatch(true); // No mostrar error cuando el campo está vacío
    }
  }, [password, confirmPassword]);
  
  const handleSubmit = async (e) => {
    e.preventDefault();
    
    // Validaciones
    if (password.length > 8) {
      setError('La contraseña no puede tener más de 8 caracteres');
      return;
    }
    
    if (password.length === 0) {
      setError('La contraseña no puede estar vacía');
      return;
    }
    
    if (password !== confirmPassword) {
      setError('Las contraseñas no coinciden');
      return;
    }
    
    // Validación adicional: verificar que no contenga espacios
    if (password.includes(' ')) {
      setError('La contraseña no puede contener espacios en blanco');
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
      <div className="recover-password-container">
        <div className="recover-password-card">
          <div className="recover-password-content">
            <div className="recover-password-error">{error}</div>
            <Link to="/recuperar-password" className="recover-password-button">
              Solicitar nuevo enlace
            </Link>
          </div>
        </div>
      </div>
    );
  }
  
  return (
    <div className="recover-password-container">
      <div className="recover-password-card">
        <div className="recover-password-content">
          {!success ? (
            <>
              <h2 className="recover-password-title">Establecer nueva contraseña</h2>
              <p className="recover-password-description">
                Ingresa tu nueva contraseña para tu cuenta.
              </p>
              
              {error && <div className="recover-password-error">{error}</div>}
              
              <form onSubmit={handleSubmit}>
                <div className="recover-password-form-group">
                  <label htmlFor="password">Nueva contraseña <span className="password-limit">(máximo 8 caracteres)</span></label>
                  <input
                    type="password"
                    id="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder="Ingresa tu nueva contraseña"
                    maxLength={8}
                    required
                  />
                  <div className="character-count">{password.length}/8 caracteres</div>
                </div>
                
                <div className="recover-password-form-group">
                  <label htmlFor="confirmPassword">Confirmar contraseña</label>
                  <input
                    type="password"
                    id="confirmPassword"
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    placeholder="Confirma tu nueva contraseña"
                    maxLength={8}
                    required
                    className={!passwordsMatch ? "password-mismatch" : ""}
                  />
                  {!passwordsMatch && (
                    <div className="password-feedback">Las contraseñas no coinciden</div>
                  )}
                </div>
                
                <button 
                  type="submit" 
                  className="recover-password-button"
                  disabled={cargando || password.length === 0 || !passwordsMatch}
                >
                  {cargando ? 'Actualizando...' : 'Guardar nueva contraseña'}
                </button>
              </form>
            </>
          ) : (
            <div className="recover-password-success">
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