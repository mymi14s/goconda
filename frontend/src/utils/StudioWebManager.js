// StudioWebManager.js
import axios from 'axios';
import alertify from 'alertifyjs';
import Utils from './Utils';

const validator = new Utils.Validator();

const isDev = window.location.port === '3000';


const devBaseURL = 'http://localhost:8080'; // HTTPS in dev for cross-site cookie
const prodBaseURL = ''; // same-origin in prod
class StudioWebManager {
  constructor(baseURL = isDev ? devBaseURL : prodBaseURL) {
    this.api = axios.create({
      baseURL,
      headers: { 'Content-Type': 'application/json' },
      xsrfCookieName: 'XSRF-TOKEN',   // cookie name from server
      xsrfHeaderName: 'X-XSRF-TOKEN', // axios auto-adds this header
      withCredentials: true,
    });

    this.api.interceptors.response.use(
      (response) => {
        if (![200, 201, 202].includes(response.status)) {
          const msg = response.data?.error || 'Unexpected response';
          alertify.error(msg);
          return Promise.reject({ ...response, message: msg });
        }
        return response;
      },
      (error) => {
        const msg = error?.response?.data?.error || error.message || 'API request failed';
        validator.error(msg);
        // if (error?.response?.status === 401) {
        //   window.location.href = '/';
        // }
        return Promise.reject(error);
      }
    );
  }

  async initCsrf() {
    // hit once on app start so the server can set XSRF-TOKEN
    try { await this.api.get('/csrf'); } catch {}
  }

  async login(email, password) {
    if (!email || !password) return validator.error('Email and password are required!');
    if (!validator.isValidEmail(email)) return validator.error('Email is invalid!');

    const res = await this.api.post('/api/v1/auth/login', { email, password });
    alertify.success('Logged in successfully');
    return res.data; // contains safe user info only
  }

  async set_sio_sid(sid) {
    if (!sid) return validator.error('Socket IO SID not found!');
    await this.api.post('/user/sio-sid', { sid });
    sessionStorage.setItem('sio_sid', sid); // harmless
  }

  async logout() {
    await this.api.post('/api/v1/auth/logout');
    alertify.success('Logged out successfully');
    window.location.href = '/';
  }

  async fetchSettings() {
    const res = await this.api.get('/api/settings/');
    return res.data;
  }
}

export default StudioWebManager;
