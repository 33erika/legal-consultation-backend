import type {
  ApiResponse,
  PaginatedResponse,
  Consultation,
  CreateConsultationRequest,
  Reply,
  ListQuery,
} from '../types';

// 员工端 API
export const employeeApi = {
  // 创建咨询
  createConsultation: async (
    data: CreateConsultationRequest
  ): Promise<Consultation> => {
    const response = await request<Consultation>('/employee/consultations', {
      method: 'POST',
      body: JSON.stringify(data),
    });
    return response.data;
  },

  // 获取我的咨询列表
  getMyConsultations: async (
    query: ListQuery = {}
  ): Promise<PaginatedResponse<Consultation>> => {
    const params = new URLSearchParams();
    if (query.page) params.append('page', String(query.page));
    if (query.page_size) params.append('page_size', String(query.page_size));
    if (query.status) params.append('status', query.status);
    if (query.urgency) params.append('urgency', query.urgency);
    if (query.start_date) params.append('start_date', query.start_date);
    if (query.end_date) params.append('end_date', query.end_date);

    const response = await request<PaginatedResponse<Consultation>>(
      `/employee/consultations?${params.toString()}`
    );
    return response.data;
  },

  // 获取咨询详情
  getConsultation: async (id: string): Promise<Consultation> => {
    const response = await request<Consultation>(`/employee/consultations/${id}`);
    return response.data;
  },

  // 评价咨询
  rateConsultation: async (id: string, rating: number): Promise<void> => {
    await request(`/employee/consultations/${id}/rate`, {
      method: 'POST',
      body: JSON.stringify({ rating }),
    });
  },

  // 关闭咨询
  closeConsultation: async (id: string): Promise<void> => {
    await request(`/employee/consultations/${id}/close`, {
      method: 'POST',
    });
  },
};

// 法务专员端 API
export const legalApi = {
  // 获取待处理咨询列表
  getPendingConsultations: async (
    query: ListQuery = {}
  ): Promise<PaginatedResponse<Consultation>> => {
    const params = new URLSearchParams();
    if (query.page) params.append('page', String(query.page));
    if (query.page_size) params.append('page_size', String(query.page_size));
    if (query.status) params.append('status', query.status);
    if (query.urgency) params.append('urgency', query.urgency);

    const response = await request<PaginatedResponse<Consultation>>(
      `/legal/consultations?${params.toString()}`
    );
    return response.data;
  },

  // 获取已分配给我的咨询列表
  getMyAssignedConsultations: async (
    query: ListQuery = {}
  ): Promise<PaginatedResponse<Consultation>> => {
    const params = new URLSearchParams();
    if (query.page) params.append('page', String(query.page));
    if (query.page_size) params.append('page_size', String(query.page_size));

    const response = await request<PaginatedResponse<Consultation>>(
      `/legal/consultations/my?${params.toString()}`
    );
    return response.data;
  },

  // 接受咨询
  acceptConsultation: async (id: string): Promise<void> => {
    await request(`/legal/consultations/${id}/accept`, {
      method: 'POST',
    });
  },

  // 回复咨询
  replyConsultation: async (
    id: string,
    content: string,
    isInternal: boolean = false
  ): Promise<Reply> => {
    const response = await request<Reply>(`/legal/consultations/${id}/reply`, {
      method: 'POST',
      body: JSON.stringify({ content, is_internal: isInternal }),
    });
    return response.data;
  },

  // 要求补充资料
  requestSupplement: async (id: string, content: string): Promise<void> => {
    await request(`/legal/consultations/${id}/supplement`, {
      method: 'POST',
      body: JSON.stringify({ content }),
    });
  },

  // 完成咨询
  completeConsultation: async (id: string, content: string): Promise<void> => {
    await request(`/legal/consultations/${id}/complete`, {
      method: 'POST',
      body: JSON.stringify({ content }),
    });
  },

  // 获取工作台统计
  getDashboardStats: async () => {
    const response = await request('/legal/dashboard');
    return response.data;
  },

  // 获取咨询详情
  getConsultation: async (id: string): Promise<Consultation> => {
    const response = await request<Consultation>(`/legal/consultations/${id}`);
    return response.data;
  },
};

// 上传附件
export const uploadApi = {
  upload: async (file: File): Promise<{ id: string; url: string; filename: string }> => {
    const formData = new FormData();
    formData.append('file', file);

    const token = localStorage.getItem('token');
    const response = await fetch('/api/v1/common/upload', {
      method: 'POST',
      headers: {
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
      },
      body: formData,
    });

    const data = await response.json();
    if (data.code !== 0) {
      throw new Error(data.message || 'Upload failed');
    }
    return data.data;
  },
};

// 辅助函数
async function request<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<ApiResponse<T>> {
  const token = localStorage.getItem('token');
  const headers: HeadersInit = {
    ...(options.headers || {}),
  };

  // 不设置 Content-Type，让浏览器自动设置 multipart/form-data
  if (!(options.body instanceof FormData)) {
    headers['Content-Type'] = 'application/json';
  }

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(`${endpoint.startsWith('/') ? '' : '/api/v1'}${endpoint}`, {
    ...options,
    headers,
  });

  const data: ApiResponse<T> = await response.json();

  if (data.code !== 0) {
    throw new Error(data.message || 'Request failed');
  }

  return data;
}
