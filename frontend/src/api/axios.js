// axios实例配置
import axios from 'axios';
import API_BASE_URL from './config';

// 创建axios实例
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
});

// 请求拦截器
apiClient.interceptors.request.use(
  config => {
    // 从localStorage获取token
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  error => {
    return Promise.reject(error);
  }
);

// 响应拦截器
apiClient.interceptors.response.use(
  response => {
    return response;
  },
  error => {
    // 统一处理错误
    if (error.response) {
      // 服务器返回错误状态码
      console.error('API Error:', error.response.data);
    } else if (error.request) {
      // 请求已发出但没有收到响应
      console.error('Network Error:', error.request);
    } else {
      // 请求配置出错
      console.error('Request Error:', error.message);
    }
    return Promise.reject(error);
  }
);

export default apiClient;