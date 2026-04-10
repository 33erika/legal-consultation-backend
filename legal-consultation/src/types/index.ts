// API 响应类型
export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data: T;
}

export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  page_size: number;
}

// 用户相关
export type UserRole = 'employee' | 'supervisor' | 'legal_staff' | 'legal_head' | 'admin';

export interface User {
  id: string;
  username: string;
  name: string;
  email: string;
  phone?: string;
  department?: string;
  role: UserRole;
  position?: string;
  status: 'active' | 'inactive';
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

// 咨询相关
export type ConsultationStatus =
  | 'pending'
  | 'accepted'
  | 'in_progress'
  | 'waiting_supplement'
  | 'completed'
  | 'closed';

export type ConsultationUrgency = 'low' | 'normal' | 'high' | 'urgent';

export interface Consultation {
  id: string;
  ticket_no: string;
  title: string;
  description: string;
  category: string;
  status: ConsultationStatus;
  urgency: ConsultationUrgency;
  attachments: Attachment[];
  user_id: string;
  user_name: string;
  user_department: string;
  legal_staff_id?: string;
  legal_staff_name?: string;
  replies?: Reply[];
  rating?: number;
  created_at: string;
  updated_at: string;
  closed_at?: string;
}

export interface Reply {
  id: string;
  consultation_id: string;
  content: string;
  is_internal: boolean;
  attachments?: Attachment[];
  user_id: string;
  user_name: string;
  user_role: UserRole;
  created_at: string;
}

export interface CreateConsultationRequest {
  title: string;
  description: string;
  category: string;
  urgency: ConsultationUrgency;
  attachment_ids?: string[];
}

// 附件相关
export interface Attachment {
  id: string;
  filename: string;
  url: string;
  size: number;
  mime_type: string;
  created_at: string;
}

// 模板申请相关
export type TemplateRequestStatus =
  | 'pending_approval'
  | 'drafting'
  | 'pending_review'
  | 'approved'
  | 'rejected';

export interface TemplateRequest {
  id: string;
  ticket_no: string;
  title: string;
  contract_type: string;
  party_a: string;
  party_b: string;
  description: string;
  status: TemplateRequestStatus;
  attachments: Attachment[];
  user_id: string;
  user_name: string;
  user_department: string;
  legal_staff_id?: string;
  legal_staff_name?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateTemplateRequestRequest {
  title: string;
  contract_type: string;
  party_a: string;
  party_b: string;
  description: string;
  attachment_ids?: string[];
}

// 合同模板相关
export interface Template {
  id: string;
  name: string;
  contract_type: string;
  description: string;
  file_url: string;
  version: number;
  is_active: boolean;
  usage_count: number;
  created_at: string;
  updated_at: string;
}

// 统计相关
export interface OverviewStats {
  total_consultations: number;
  pending_consultations: number;
  completed_consultations: number;
  avg_response_time: number;
  satisfaction_rate: number;
}

export interface CategoryDistribution {
  category: string;
  count: number;
  percentage: number;
}

export interface ProcessingEfficiency {
  date: string;
  count: number;
  avg_time: number;
}

export interface StaffWorkload {
  staff_id: string;
  staff_name: string;
  total_tasks: number;
  completed_tasks: number;
  avg_completion_time: number;
}

// 法务工作台
export interface DashboardStats {
  consultation_stats: {
    pending: number;
    in_progress: number;
    completed_today: number;
    total: number;
  };
  template_stats: {
    pending_approval: number;
    drafting: number;
    pending_review: number;
  };
}

// 列表查询参数
export interface ListQuery {
  page?: number;
  page_size?: number;
  keyword?: string;
  status?: string;
  urgency?: string;
  start_date?: string;
  end_date?: string;
}
