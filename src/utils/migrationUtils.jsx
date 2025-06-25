export const migrationComplete = () => {
  try {
    // Verificar si ya se ejecutó la migración
    const migrationCompleted = sessionStorage.getItem('migration_completed');
    
    if (!migrationCompleted) {
      console.log('🔄 Iniciando migración de localStorage a sessionStorage...');
      
      // Lista de keys relacionadas con datos de usuario
      const userDataKeys = [
        'userData',
        'user', 
        'authToken'
      ];

      let migrated = false;

      // Migrar datos de usuario
      userDataKeys.forEach(key => {
        const localValue = localStorage.getItem(key);
        
        if (localValue) {
          // Solo migrar si no existe en sessionStorage
          const sessionValue = sessionStorage.getItem(key);
          if (!sessionValue) {
            sessionStorage.setItem(key, localValue);
            console.log(`✅ Migrado: ${key}`);
            migrated = true;
          }
          
          // Limpiar localStorage para evitar conflictos futuros
          localStorage.removeItem(key);
          console.log(`🧹 Limpiado localStorage: ${key}`);
        }
      });

      // Marcar migración como completada
      sessionStorage.setItem('migration_completed', 'true');
      
      if (migrated) {
        console.log('✅ Migración completada exitosamente');
        console.log('🔒 Los datos de usuario ahora están aislados por ventana');
      } else {
        console.log('ℹ️ No se encontraron datos para migrar');
      }
    }
    
    return true;
  } catch (error) {
    console.error('❌ Error durante la migración:', error);
    return false;
  }
};

// Función para verificar el estado actual de las sesiones
export const checkCurrentSession = () => {
  const sessionUser = sessionStorage.getItem('userData');
  const localUser = localStorage.getItem('userData');
  
  console.log('📊 Estado actual de la sesión:');
  console.log('- sessionStorage userData:', sessionUser ? '✅ Presente' : '❌ Ausente');
  console.log('- localStorage userData:', localUser ? '⚠️ Presente (debería estar limpio)' : '✅ Limpio');
  
  if (localUser && sessionUser) {
    console.warn('⚠️ ADVERTENCIA: Hay datos duplicados en localStorage y sessionStorage');
    console.warn('💡 Recomendación: Limpiar localStorage manualmente');
  }
  
  return {
    hasSessionData: !!sessionUser,
    hasLocalData: !!localUser,
    isDuplicated: !!(sessionUser && localUser)
  };
};

// Función para limpiar manualmente localStorage si es necesario
export const cleanLocalStorage = () => {
  try {
    const userDataKeys = ['userData', 'user', 'authToken'];
    
    userDataKeys.forEach(key => {
      if (localStorage.getItem(key)) {
        localStorage.removeItem(key);
        console.log(`🧹 Limpiado manualmente: ${key}`);
      }
    });
    
    console.log('✅ Limpieza manual completada');
    return true;
  } catch (error) {
    console.error('❌ Error durante la limpieza manual:', error);
    return false;
  }
};
