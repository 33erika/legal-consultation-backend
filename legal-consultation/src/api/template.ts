import type {
  ApiResponse,
  PaginatedResponse,
  Template,
  TemplateRequest,
  CreateTemplateRequestRequest,
  ListQuery,
} from '../types';

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

// 模板相关
export const templateApi = {
  // 获取模板列表
  getTemplates: async (query: { category?: string } = {}): Promise<Template[]> => {
    const params = new URLSearchParams();
    if (query.category) params.append('category', query.category);

    const response = await request<PaginatedResponse<Template>>(
      `/common/templates?${params.toString()}`
    );
    return response.data.items;
  },

  // 获取模板详情
  getTemplate: async (id: string): Promise<Template> => {
    const response = await request<Template>(`/common/templates/${id}`);
    return response.data;
  },

  // 更新模板使用次数
  incrementUsage: async (id: string): Promise<void> => {
    await request(`/common/templates/${id}/increment-usage`, {
      method: 'POST',
    });
  },
};

// 模板申请相关 (员工端)
export const templateRequestApi = {
  // 创建模板申请
  create: async (data: CreateTemplateRequestRequest): Promise<TemplateRequest> => {
    const response = await request<TemplateRequest>('/employee/template-requests', {
      method: 'POST',
      body: JSON.stringify(data),
    });
    return response.data;
  },

  // 获取我的模板申请列表
  getMyRequests: async (
    query: ListQuery = {}
  ): Promise<PaginatedResponse<TemplateRequest>> => {
    const params = new URLSearchParams();
    if (query.page) params.append('page', String(query.page));
    if (query.page_size) params.append('page_size', String(query.page_size));
    if (query.status) params.append('status', query.status);

    const response = await request<PaginatedResponse<TemplateRequest>>(
      `/employee/template-requests?${params.toString()}`
    );
    return response.data;
  },

  // 获取模板申请详情
  getRequest: async (id: string): Promise<TemplateRequest> => {
    const response = await request<TemplateRequest>(`/employee/template-requests/${id}`);
    return response.data;
  },
};

// 模板审批相关 (法务端)
export const templateApprovalApi = {
  // 获取待审批列表
  getPendingApprovals: async (
    query: ListQuery = {}
  ): Promise<PaginatedResponse<TemplateRequest>> => {
    const params = new URLSearchParams();
    if (query.page) params.append('page', String(query.page));
    if (query.page_size) params.append('page_size', String(query.page_size));

    const response = await request<PaginatedResponse<TemplateRequest>>(
      `/legal/template-requests/pending?${params.toString()}`
    );
    return response.data;
  },

  // 审批通过
  approve: async (id: string, draftUrl?: string): Promise<void> => {
    await request(`/legal/template-requests/${id}/approve`, {
      method: 'POST',
      body: JSON.stringify({ draft_url: draftUrl }),
    });
  },

  // 审批拒绝
  reject: async (id: string, reason: string): Promise<void> => {
    await request(`/legal/template-requests/${id}/reject`, {
      method: 'POST',
      body: JSON.stringify({ reason }),
    });
  },

  // 开始起草
  startDrafting: async (id: string): Promise<void> => {
    await request(`/legal/template-requests/${id}/start-drafting`, {
      method: 'POST',
    });
  },

  // 提交审核
  submitForReview: async (id: string, draftUrl: string): Promise<void> => {
    await request(`/legal/template-requests/${id}/submit-review`, {
      method: 'POST',
      body: JSON.stringify({ draft_url: draftUrl }),
    });
  },
};
