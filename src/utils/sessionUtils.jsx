// Utilidad para migrar datos de localStorage a sessionStorage
// Esto evita conflictos entre múltiples ventanas/pestañas con diferentes usuarios

export const migrateToSessionStorage = () => {
  try {
    // Lista de keys que necesitan ser migradas de localStorage a sessionStorage
    const keysToMigrate = [
      'userData',
      'user', 
      'authToken'
    ];

    let migrated = false;

    keysToMigrate.forEach(key => {
      // Si existe en localStorage pero no en sessionStorage, migrar
      const localValue = localStorage.getItem(key);
      const sessionValue = sessionStorage.getItem(key);
      
      if (localValue && !sessionValue) {
        sessionStorage.setItem(key, localValue);
        console.log(`✅ Migrado ${key} de localStorage a sessionStorage`);
        migrated = true;
      }
    });

    // Opcional: Limpiar localStorage después de migrar para evitar conflictos futuros
    if (migrated) {
      console.log('🔄 Migración completada. Datos de usuario ahora aislados por ventana.');
      
      // Comentado por seguridad - descomentar si quieres limpiar localStorage
      // keysToMigrate.forEach(key => {
      //   localStorage.removeItem(key);
      // });
    }

    return migrated;
  } catch (error) {
    console.error('❌ Error durante la migración:', error);
    return false;
  }
};

// Función para verificar si hay conflictos de sesión
export const checkSessionConflict = () => {
  try {
    const localUserData = localStorage.getItem('userData');
    const sessionUserData = sessionStorage.getItem('userData');
    
    if (localUserData && sessionUserData) {
      const localUser = JSON.parse(localUserData);
      const sessionUser = JSON.parse(sessionUserData);
      
      // Si son usuarios diferentes, hay conflicto
      if (localUser.id !== sessionUser.id) {
        console.warn('⚠️ Detectado conflicto de sesión entre ventanas');
        return {
          hasConflict: true,
          localUser: localUser,
          sessionUser: sessionUser
        };
      }
    }
    
    return { hasConflict: false };
  } catch (error) {
    console.error('Error al verificar conflictos de sesión:', error);
    return { hasConflict: false };
  }
};

// Función para obtener información del usuario de manera segura
export const getSafeUserData = () => {
  try {
    // Priorizar sessionStorage sobre localStorage
    let userData = sessionStorage.getItem('userData');
    
    // Si no hay en sessionStorage, intentar migrar desde localStorage
    if (!userData) {
      const localUserData = localStorage.getItem('userData');
      if (localUserData) {
        sessionStorage.setItem('userData', localUserData);
        userData = localUserData;
        console.log('🔄 Migrado userData a sessionStorage automáticamente');
      }
    }
    
    return userData ? JSON.parse(userData) : null;
  } catch (error) {
    console.error('Error al obtener datos de usuario:', error);
    return null;
  }
};

// Función para limpiar datos de sesión al cerrar
export const clearSessionData = () => {
  try {
    sessionStorage.removeItem('userData');
    sessionStorage.removeItem('user');
    sessionStorage.removeItem('authToken');
    console.log('🧹 Datos de sesión limpiados');
  } catch (error) {
    console.error('Error al limpiar datos de sesión:', error);
  }
};
