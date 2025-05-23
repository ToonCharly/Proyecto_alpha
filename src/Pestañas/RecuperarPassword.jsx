import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import '../styles/Login.css';

const RecuperarPassword = () => {
  const [email, setEmail] = useState('');
  const [enviado, setEnviado] = useState(false);
  const [error, setError] = useState('');
  const [cargando, setCargando] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    // Validación básica
    if (!email || !email.includes('@')) {
      setError('Por favor ingresa un correo electrónico válido');
      return;
    }
    
    setCargando(true);
    setError('');
    
    try {
      const response = await fetch('http://localhost:8080/api/reset-password-request', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email }),
      });
      
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Error al procesar tu solicitud');
      }
      
      setEnviado(true);
    } catch (err) {
      setError(err.message || 'Ocurrió un error al procesar tu solicitud. Inténtalo de nuevo.');
      console.error('Error:', err);
    } finally {
      setCargando(false);
    }
  };

  return (
    <div className="login-container">
      <div className="login-form">
        <div className="logo-container">
          <div className="logo-circle"></div>
          <h2 className="logo-text">FACTS</h2>
        </div>
        
        <div className="login-form-content">
          {!enviado ? (
            <>
              <h2 className="sign-in-title">Recuperar contraseña</h2>
              <p className="login-description">
                Ingresa tu correo electrónico y te enviaremos instrucciones para restablecer tu contraseña.
              </p>
              
              {error && <div className="error-message">{error}</div>}
              
              <form onSubmit={handleSubmit}>
                <div className="form-group">
                  <label htmlFor="email">Correo electrónico</label>
                  <input
                    type="email"
                    id="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    placeholder="Ingresa tu correo electrónico"
                    required
                  />
                </div>
                
                <button 
                  type="submit" 
                  className="login-button"
                  disabled={cargando}
                >
                  {cargando ? 'Enviando...' : 'Enviar instrucciones'}
                </button>
              </form>
              
              <div className="form-footer">
                <Link to="/login" className="register-link">
                  Volver al inicio de sesión
                </Link>
              </div>
            </>
          ) : (
            <div className="success-message">
              <h2>¡Correo enviado!</h2>
              <p>
                Hemos enviado un correo a <strong>{email}</strong> con instrucciones para restablecer tu contraseña.
              </p>
              <p>
                Si no encuentras el correo, revisa en tu carpeta de spam.
              </p>
              <Link to="/login" className="login-button">
                Volver al inicio de sesión
              </Link>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default RecuperarPassword;