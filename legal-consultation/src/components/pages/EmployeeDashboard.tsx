import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { api } from '../../lib/api';
import type { Consultation, ListQuery } from '../../types';
import { Card } from '../ui/Card';
import { Button } from '../ui/Button';
import { Table } from '../ui/Table';
import { Modal } from '../ui/Modal';
import { Input } from '../ui/Input';
import { Select } from '../ui/Select';
import { Textarea } from '../ui/Textarea';

const statusLabels: Record<string, string> = {
  pending: '待处理',
  accepted: '已受理',
  in_progress: '处理中',
  waiting_supplement: '等待补充',
  completed: '已完成',
  closed: '已关闭',
};

const urgencyLabels: Record<string, string> = {
  low: '低',
  normal: '普通',
  high: '高',
  urgent: '紧急',
};

const categoryOptions = [
  { value: 'contract', label: '合同相关' },
  { value: 'labor', label: '劳动人事' },
  { value: 'ip', label: '知识产权' },
  { value: 'litigation', label: '诉讼支持' },
  { value: 'compliance', label: '合规审查' },
  { value: 'other', label: '其他' },
];

const urgencyOptions = [
  { value: 'low', label: '低' },
  { value: 'normal', label: '普通' },
  { value: 'high', label: '高' },
  { value: 'urgent', label: '紧急' },
];

export function EmployeeDashboard() {
  const { user, logout } = useAuth();
  const [consultations, setConsultations] = useState<Consultation[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const [newConsultation, setNewConsultation] = useState({
    title: '',
    description: '',
    category: 'contract',
    urgency: 'normal' as const,
  });

  useEffect(() => {
    loadConsultations();
  }, [page]);

  const loadConsultations = async () => {
    setIsLoading(true);
    try {
      const query: ListQuery = { page, page_size: pageSize };
      const response = await api.getConsultations(query);
      setConsultations(response.items);
      setTotal(response.total);
    } catch (error) {
      console.error('Failed to load consultations:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateConsultation = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await api.createConsultation(newConsultation);
      setIsCreateModalOpen(false);
      setNewConsultation({ title: '', description: '', category: 'contract', urgency: 'normal' });
      loadConsultations();
    } catch (error) {
      console.error('Failed to create consultation:', error);
    }
  };

  const columns = [
    { key: 'ticket_no', label: '工单号', width: '120px' },
    { key: 'title', label: '标题' },
    { key: 'category', label: '类型', width: '100px' },
    { key: 'urgency', label: '紧急程度', width: '100px' },
    { key: 'status', label: '状态', width: '100px' },
    { key: 'created_at', label: '创建时间', width: '180px' },
    { key: 'actions', label: '操作', width: '100px' },
  ];

  const totalPages = Math.ceil(total / pageSize);

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
          <h1 className="text-xl font-bold">法务部法律咨询系统 - 员工端</h1>
          <div className="flex items-center gap-4">
            <span className="text-gray-600">{user?.name} ({user?.department})</span>
            <Button variant="outline" size="sm" onClick={logout}>退出</Button>
          </div>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Stats */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
          <Card className="p-4">
            <p className="text-gray-500 text-sm">我的咨询</p>
            <p className="text-2xl font-bold">{total}</p>
          </Card>
          <Card className="p-4">
            <p className="text-gray-500 text-sm">处理中</p>
            <p className="text-2xl font-bold">
              {consultations.filter(c => ['accepted', 'in_progress'].includes(c.status)).length}
            </p>
          </Card>
          <Card className="p-4">
            <p className="text-gray-500 text-sm">已完成</p>
            <p className="text-2xl font-bold">
              {consultations.filter(c => c.status === 'completed').length}
            </p>
          </Card>
          <Card className="p-4">
            <p className="text-gray-500 text-sm">待评价</p>
            <p className="text-2xl font-bold">
              {consultations.filter(c => c.status === 'closed' && !c.rating).length}
            </p>
          </Card>
        </div>

        {/* Actions */}
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-lg font-semibold">我的咨询</h2>
          <Button onClick={() => setIsCreateModalOpen(true)}>+ 新建咨询</Button>
        </div>

        {/* Table */}
        <Card className="p-0 overflow-hidden">
          <Table columns={columns}>
            {consultations.map((c) => (
              <tr key={c.id} className="border-b last:border-b-0">
                <td className="px-4 py-3">{c.ticket_no}</td>
                <td className="px-4 py-3">{c.title}</td>
                <td className="px-4 py-3">
                  {categoryOptions.find(o => o.value === c.category)?.label || c.category}
                </td>
                <td className="px-4 py-3">
                  <span className={`px-2 py-1 rounded text-xs ${
                    c.urgency === 'urgent' ? 'bg-red-100 text-red-700' :
                    c.urgency === 'high' ? 'bg-orange-100 text-orange-700' :
                    'bg-gray-100 text-gray-700'
                  }`}>
                    {urgencyLabels[c.urgency]}
                  </span>
                </td>
                <td className="px-4 py-3">
                  <span className={`px-2 py-1 rounded text-xs ${
                    c.status === 'pending' ? 'bg-yellow-100 text-yellow-700' :
                    c.status === 'in_progress' ? 'bg-blue-100 text-blue-700' :
                    c.status === 'completed' ? 'bg-green-100 text-green-700' :
                    'bg-gray-100 text-gray-700'
                  }`}>
                    {statusLabels[c.status]}
                  </span>
                </td>
                <td className="px-4 py-3 text-sm text-gray-500">
                  {new Date(c.created_at).toLocaleString('zh-CN')}
                </td>
                <td className="px-4 py-3">
                  <Link to={`/consultation/${c.id}`}>
                    <Button variant="ghost" size="sm">查看</Button>
                  </Link>
                </td>
              </tr>
            ))}
          </Table>

          {/* Pagination */}
          <div className="px-4 py-3 border-t flex justify-between items-center">
            <p className="text-sm text-gray-500">
              共 {total} 条，第 {page}/{totalPages} 页
            </p>
            <div className="flex gap-2">
              <Button
                variant="outline"
                size="sm"
                disabled={page <= 1}
                onClick={() => setPage(p => p - 1)}
              >
                上一页
              </Button>
              <Button
                variant="outline"
                size="sm"
                disabled={page >= totalPages}
                onClick={() => setPage(p => p + 1)}
              >
                下一页
              </Button>
            </div>
          </div>
        </Card>
      </main>

      {/* Create Modal */}
      <Modal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        title="新建咨询"
      >
        <form onSubmit={handleCreateConsultation} className="space-y-4">
          <Input
            label="标题"
            value={newConsultation.title}
            onChange={(e) => setNewConsultation(p => ({ ...p, title: e.target.value }))}
            placeholder="请输入咨询标题"
            required
          />

          <Select
            label="类型"
            value={newConsultation.category}
            onChange={(e) => setNewConsultation(p => ({ ...p, category: e.target.value }))}
            options={categoryOptions}
          />

          <Select
            label="紧急程度"
            value={newConsultation.urgency}
            onChange={(e) => setNewConsultation(p => ({ ...p, urgency: e.target.value as any }))}
            options={urgencyOptions}
          />

          <Textarea
            label="详细描述"
            value={newConsultation.description}
            onChange={(e) => setNewConsultation(p => ({ ...p, description: e.target.value }))}
            placeholder="请详细描述您的法律咨询问题"
            rows={5}
            required
          />

          <div className="flex justify-end gap-2 pt-4">
            <Button variant="outline" type="button" onClick={() => setIsCreateModalOpen(false)}>
              取消
            </Button>
            <Button type="submit">提交</Button>
          </div>
        </form>
      </Modal>
    </div>
  );
}
