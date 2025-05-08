import { message } from 'antd';
import { login, getUserInfo } from './api';

interface UserInfo {
  username: string;
  token: string;
}

const TOKEN_KEY = 'oss_token';
const USER_INFO_KEY = 'oss_user_info';

export const getStoredUserInfo = (): UserInfo | null => {
  const stored = localStorage.getItem(USER_INFO_KEY);
  return stored ? JSON.parse(stored) : null;
};

export const setStoredUserInfo = (userInfo: UserInfo) => {
  localStorage.setItem(USER_INFO_KEY, JSON.stringify(userInfo));
  localStorage.setItem(TOKEN_KEY, userInfo.token);
};

export const clearStoredUserInfo = () => {
  localStorage.removeItem(USER_INFO_KEY);
  localStorage.removeItem(TOKEN_KEY);
};

export const getStoredToken = (): string | null => {
  return localStorage.getItem(TOKEN_KEY);
};

export const handleLogin = async (username: string, password: string): Promise<boolean> => {
  try {
    const { data } = await login({ username, password });
    const userInfo: UserInfo = {
      username,
      token: data.token,
    };
    setStoredUserInfo(userInfo);
    return true;
  } catch (error) {
    message.error('登录失败');
    return false;
  }
};

export const validateToken = async (): Promise<boolean> => {
  const token = getStoredToken();
  if (!token) return false;

  try {
    // 使用用户信息接口验证 token
    const { data } = await getUserInfo();
    // 更新存储的用户信息
    const userInfo = getStoredUserInfo();
    if (userInfo) {
      userInfo.username = data.username;
      setStoredUserInfo(userInfo);
    }
    return true;
  } catch (error) {
    clearStoredUserInfo();
    return false;
  }
}; 