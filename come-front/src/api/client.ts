// src/api/client.ts

import axios from 'axios';

const apiClient = axios.create({
  baseURL: '/api/',
  headers: { 'Content-Type': 'application/json' },
});

apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = token;
  }
  return config;
}, (error) => Promise.reject(error));

export default apiClient;

