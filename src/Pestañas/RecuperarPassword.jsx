import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import '../STYLES/RecoverPassword.css';

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
    <div className="recover-password-container">
      <div className="recover-password-card">
        <div className="recover-password-content">
          {!enviado ? (
            <>
              <h2 className="recover-password-title">Recuperar contraseña</h2>
              <p className="recover-password-description">
                Ingresa tu correo electrónico y te enviaremos instrucciones para restablecer tu contraseña.
              </p>
              
              {error && <div className="recover-password-error">{error}</div>}
              
              <form onSubmit={handleSubmit}>
                <div className="recover-password-form-group">
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
                  className="recover-password-button"
                  disabled={cargando}
                >
                  {cargando ? 'Enviando...' : 'Enviar'}
                </button>
              </form>
              
              <div className="recover-password-footer">
                <Link to="/login" className="recover-password-link">
                  Volver al inicio de sesión
                </Link>
              </div>
            </>
          ) : (
            <div className="recover-password-success">
              <h2>¡Correo enviado!</h2>
              <p>
                Hemos enviado un correo a <strong>{email}</strong> con instrucciones para restablecer tu contraseña.
              </p>
              <p>
                Si no encuentras el correo, revisa en tu carpeta de spam.
              </p>
              <Link to="/login" className="recover-password-button">
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