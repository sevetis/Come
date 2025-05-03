// src/api/client.ts

import apiClient from "./client";

export interface AuthPayload {
  email: string;
  password: string;
}

export interface RegisterPayload extends AuthPayload {
  username: string;
}

export interface LoginPayload {
  user_id: number;
  role: number;
  exp: number;
}

export const login = async (data: AuthPayload) => {
  const response = await apiClient.post('/user/login', data);
  return response.data;
};

export const register = async (data: RegisterPayload) => {
  const response = await apiClient.post('/user/register', data);
  return response.data;
};
