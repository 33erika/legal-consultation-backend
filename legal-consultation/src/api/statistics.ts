import type {
  ApiResponse,
  OverviewStats,
  CategoryDistribution,
  ProcessingEfficiency,
  StaffWorkload,
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

export const statisticsApi = {
  // 获取概览统计
  getOverview: async (): Promise<OverviewStats> => {
    const response = await request<OverviewStats>('/statistics/overview');
    return response.data;
  },

  // 获取分类分布
  getCategoryDistribution: async (
    startDate?: string,
    endDate?: string
  ): Promise<CategoryDistribution[]> => {
    const params = new URLSearchParams();
    if (startDate) params.append('start_date', startDate);
    if (endDate) params.append('end_date', endDate);

    const response = await request<CategoryDistribution[]>(
      `/statistics/category-distribution?${params.toString()}`
    );
    return response.data;
  },

  // 获取处理效率
  getProcessingEfficiency: async (
    startDate?: string,
    endDate?: string
  ): Promise<ProcessingEfficiency[]> => {
    const params = new URLSearchParams();
    if (startDate) params.append('start_date', startDate);
    if (endDate) params.append('end_date', endDate);

    const response = await request<ProcessingEfficiency[]>(
      `/statistics/processing-efficiency?${params.toString()}`
    );
    return response.data;
  },

  // 导出统计数据
  export: async (startDate?: string, endDate?: string): Promise<Blob> => {
    const token = localStorage.getItem('token');
    const params = new URLSearchParams();
    if (startDate) params.append('start_date', startDate);
    if (endDate) params.append('end_date', endDate);

    const response = await fetch(
      `${API_BASE}/statistics/export?${params.toString()}`,
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );

    return response.blob();
  },

  // 获取员工工作量
  getStaffWorkload: async (
    startDate?: string,
    endDate?: string
  ): Promise<StaffWorkload[]> => {
    const params = new URLSearchParams();
    if (startDate) params.append('start_date', startDate);
    if (endDate) params.append('end_date', endDate);

    const response = await request<StaffWorkload[]>(
      `/statistics/staff-workload?${params.toString()}`
    );
    return response.data;
  },
};
