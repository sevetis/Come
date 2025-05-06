// src/api/post.ts

import apiClient from "./client";

export interface Post {
  id: number;
  title: string;
  content: string;
  authorId: number;
  createdAt: string;
  updatedAt: string;
}

export interface Comment {
  id: number;
  postId: number;
  authorId: number;
  content: string;
  createdAt: string;
  updatedAt: string;
}

export const getPostsPaginated = async (page: number, pageSize: number) => {
  const response = await apiClient.get(`/post/batch?page=${page}&pageSize=${pageSize}`);
  return response.data;
};

export const getPost = async (id: number): Promise<Post> => {
  const response = await apiClient.get(`/post/${id}`, {});
  return response.data;
}

export const createPost = async (post: {title: string; content: string}): Promise<Post> => {
  const response = await apiClient.post('/post/create', post)
  return response.data;
}

export const updatePost = async (id: number, updates: {title: string; content: string}): Promise<Post> => {
  const response = await apiClient.put(`/post/${id}`, updates);
  return response.data;
}

export const deletePost = async (id: number) => {
  const response = await apiClient.delete(`/post/${id}`);
  return response.data;
}

export const getPostComments = async (postId: number): Promise<Comment[]> => {
  const response = await apiClient.get(`/post/${postId}/comments`);
  return response.data || [];
};

export const createComment = async (postId: number, content: string): Promise<Comment> => {
  const response = await apiClient.post(`/post/${postId}/comment`, { content });
  return response.data;
};
