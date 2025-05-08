import axios from 'axios';
import { getStoredToken } from './auth';

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
});

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    const token = getStoredToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // 处理未授权错误
      window.location.href = '/ui/login';
    }
    return Promise.reject(error);
  }
);

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
}

export interface ObjectMetadata {
  size: number;
  lastModified: string;
  contentType: string;
  etag: string;
}

export const login = (data: LoginRequest) =>
  api.post('/login', data);

export const getUserInfo = () =>
  api.get('/oss/user/info');

export const listBuckets = () =>
  api.get('/oss/buckets');

export const createBucket = (bucket: string) =>
  api.put(`/oss/${bucket}`);

export const deleteBucket = (bucket: string) =>
  api.delete(`/oss/${bucket}`);

export const listObjects = (bucket: string) =>
  api.get(`/oss/${bucket}`);

export const uploadObject = (bucket: string, key: string, file: File) => {
  const formData = new FormData();
  formData.append('file', file);
  return api.put(`/oss/${bucket}/${key}`, formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};

export const downloadObject = (bucket: string, key: string) =>
  api.get(`/oss/${bucket}/${key}`, {
    responseType: 'blob',
  });

export const deleteObject = (bucket: string, key: string) =>
  api.delete(`/oss/${bucket}/${key}`);

export const getObjectMetadata = (bucket: string, key: string) =>
  api.head(`/oss/${bucket}/${key}`);

export const copyObject = (sourceBucket: string, sourceKey: string, targetBucket: string, targetKey: string) =>
  api.put(`/oss/${targetBucket}/${targetKey}/copy`, {
    sourceBucket,
    sourceKey,
  });

export const getObjectURL = (bucket: string, key: string) =>
  api.get(`/oss/${bucket}/${key}/url`); 