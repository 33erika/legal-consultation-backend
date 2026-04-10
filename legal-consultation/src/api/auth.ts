import type {
  ApiResponse,
  LoginRequest,
  LoginResponse,
  User,
} from '../types';

const API_BASE = '/api/v1';

function getToken(): string | null {
  return localStorage.getItem('token');
}

function setToken(token: string): void {
  localStorage.setItem('token', token);
}

function clearToken(): void {
  localStorage.removeItem('token');
}

async function request<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<ApiResponse<T>> {
  const token = getToken();
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...(options.headers || {}),
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers,
  });

  const data: ApiResponse<T> = await response.json();

  if (data.code !== 0) {
    throw new Error(data.message || 'Request failed');
  }

  return data;
}

// 认证相关
export const authApi = {
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await request<LoginResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify(data),
    });
    setToken(response.data.token);
    return response.data;
  },

  logout: async (): Promise<void> => {
    await request('/auth/logout', { method: 'POST' });
    clearToken();
  },

  getProfile: async (): Promise<User> => {
    const response = await request<User>('/auth/profile');
    return response.data;
  },

  updateProfile: async (data: Partial<User>): Promise<User> => {
    const response = await request<User>('/auth/profile', {
      method: 'PUT',
      body: JSON.stringify(data),
    });
    return response.data;
  },

  changePassword: async (oldPassword: string, newPassword: string): Promise<void> => {
    await request('/auth/change-password', {
      method: 'POST',
      body: JSON.stringify({ old_password: oldPassword, new_password: newPassword }),
    });
  },
};

export { getToken, clearToken };
