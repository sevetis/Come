// src/api/user.ts

import { UserRole } from "@/constants/roles";
import apiClient from "./client";

export interface User {
  id: number;
  username: string;
  email: string;
  avatar: string;
  role: UserRole
  banned: boolean;
}

// current user
export const getProfile = async () => {
  const response = await apiClient.get('/user/profile');
  return response.data.data;
};

export const updateProfile = async (updates: {username: string; email: string}) => {
  const response = await apiClient.put('/user/profile', updates)
  return response.data.data;
}

export const uploadAvatar = async (file: File): Promise<string> => {
  const formData = new FormData();
  formData.append('avatar', file);
  const response = await apiClient.post('/user/avatar', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  });
  return response.data.data;
}

export const getUser = async (id: number) => {
  const response = await apiClient.get(`/user/${id}`);
  return response.data.data;
}

export const getUsersBatch = async (ids: number[]) => {
  const response = await apiClient.get(`/user/batch?ids=${ids.join(',')}`);
  return response.data.data;
};

