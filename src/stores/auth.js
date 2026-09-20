import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useAuthStore = defineStore('auth', () => {
  const user = ref(JSON.parse(localStorage.getItem('elsa_host_user') || 'null'));

  const login = (email, password) => {
    if (email === 'admin@demo.com' || email.trim() !== '') {
      const userData = {
        email: email || 'admin@demo.com',
        name: 'ELSA Host',
        initials: 'EL',
        role: 'admin'
      };
      user.value = userData;
      localStorage.setItem('elsa_host_user', JSON.stringify(userData));
      return true;
    }
    return false;
  };

  const logout = () => {
    user.value = null;
    localStorage.removeItem('elsa_host_user');
  };

  const isAuthenticated = () => !!user.value;

  return {
    user,
    login,
    logout,
    isAuthenticated
  };
});
