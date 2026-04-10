import type { ApiResponse, PaginatedResponse, User, ListQuery } from '../types';

const API_BASE = '/api/v1';

async function request<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<ApiResponse<T>> {
  const token = localStorage.getItem('token');
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

export const adminApi = {
  // 获取用户列表
  getUsers: async (query: ListQuery = {}): Promise<PaginatedResponse<User>> => {
    const params = new URLSearchParams();
    if (query.page) params.append('page', String(query.page));
    if (query.page_size) params.append('page_size', String(query.page_size));
    if (query.keyword) params.append('keyword', query.keyword);
    if (query.status) params.append('status', query.status);

    const response = await request<PaginatedResponse<User>>(
      `/admin/users?${params.toString()}`
    );
    return response.data;
  },

  // 创建用户
  createUser: async (data: {
    username: string;
    password: string;
    name: string;
    email: string;
    role: string;
    department?: string;
    position?: string;
  }): Promise<User> => {
    const response = await request<User>('/admin/users', {
      method: 'POST',
      body: JSON.stringify(data),
    });
    return response.data;
  },

  // 更新用户
  updateUser: async (id: string, data: Partial<User>): Promise<User> => {
    const response = await request<User>(`/admin/users/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
    return response.data;
  },

  // 重置密码
  resetPassword: async (id: string): Promise<{ new_password: string }> => {
    const response = await request<{ new_password: string }>(
      `/admin/users/${id}/reset-password`,
      {
        method: 'POST',
      }
    );
    return response.data;
  },

  // 启用/禁用用户
  toggleUserStatus: async (id: string): Promise<void> => {
    await request(`/admin/users/${id}/toggle-status`, {
      method: 'POST',
    });
  },

  // 删除用户
  deleteUser: async (id: string): Promise<void> => {
    await request(`/admin/users/${id}`, {
      method: 'DELETE',
    });
  },
};
