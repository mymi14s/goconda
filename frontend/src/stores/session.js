import { defineStore } from 'pinia'

export const useSessionStore = defineStore('session', {
  state: () => ({
    user: null,
    is_authenticated: null,
    token: null,
    setting: null
  }),
  actions: {
    setUser(user) {
        console.log(user)
      this.user = user.user;
      this.is_authenticated = user.is_authenticated;
    },
    setToken(token) {
      this.token = token
    },
    setSetting(setting) {
      this.setting = setting
    },
    clearSession() {
      this.user = null
      this.token = null
    },
     persist: true
  }
})
