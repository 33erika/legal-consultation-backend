import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { api } from '../../lib/api';
import type { Consultation, DashboardStats } from '../../types';
import { Card } from '../ui/Card';
import { Button } from '../ui/Button';
import { Table } from '../ui/Table';
import { Select } from '../ui/Select';

export function LegalDashboard() {
  const { user, logout } = useAuth();
  const [dashboardStats, setDashboardStats] = useState<DashboardStats | null>(null);
  const [poolConsultations, setPoolConsultations] = useState<Consultation[]>([]);
  const [myTasks, setMyTasks] = useState<Consultation[]>([]);
  const [urgency, setUrgency] = useState('');
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    loadDashboard();
  }, [urgency]);

  const loadDashboard = async () => {
    setIsLoading(true);
    try {
      const [stats, pool, tasks] = await Promise.all([
        api.getDashboard(),
        api.getConsultationPool({ urgency: urgency || undefined }),
        api.getMyTasks(),
      ]);
      setDashboardStats(stats);
      setPoolConsultations(pool.items);
      setMyTasks(tasks.items);
    } catch (error) {
      console.error('Failed to load dashboard:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleAccept = async (id: string) => {
    try {
      await api.acceptConsultation(id);
      loadDashboard();
    } catch (error) {
      console.error('Failed to accept:', error);
    }
  };

  const urgencyOptions = [
    { value: '', label: '全部' },
    { value: 'urgent', label: '紧急' },
    { value: 'high', label: '高' },
    { value: 'normal', label: '普通' },
    { value: 'low', label: '低' },
  ];

  const statusLabels: Record<string, string> = {
    pending: '待处理',
    accepted: '已受理',
    in_progress: '处理中',
    waiting_supplement: '等待补充',
    completed: '已完成',
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <p className="text-gray-500">加载中...</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
          <h1 className="text-xl font-bold">法务部工作台</h1>
          <div className="flex items-center gap-4">
            <span className="text-gray-600">{user?.name}</span>
            <Button variant="outline" size="sm" onClick={logout}>退出</Button>
          </div>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Stats */}
        <div className="grid grid-cols-1 md:grid-cols-5 gap-4 mb-6">
          <Card className="p-4 bg-blue-50 border-blue-200">
            <p className="text-blue-600 text-sm">待受理</p>
            <p className="text-2xl font-bold text-blue-700">{dashboardStats?.consultation_stats.pending || 0}</p>
          </Card>
          <Card className="p-4 bg-yellow-50 border-yellow-200">
            <p className="text-yellow-600 text-sm">处理中</p>
            <p className="text-2xl font-bold text-yellow-700">{dashboardStats?.consultation_stats.in_progress || 0}</p>
          </Card>
          <Card className="p-4 bg-green-50 border-green-200">
            <p className="text-green-600 text-sm">今日完成</p>
            <p className="text-2xl font-bold text-green-700">{dashboardStats?.consultation_stats.completed_today || 0}</p>
          </Card>
          <Card className="p-4 bg-purple-50 border-purple-200">
            <p className="text-purple-600 text-sm">待审批模板</p>
            <p className="text-2xl font-bold text-purple-700">{dashboardStats?.template_stats.pending_approval || 0}</p>
          </Card>
          <Card className="p-4 bg-orange-50 border-orange-200">
            <p className="text-orange-600 text-sm">起草中模板</p>
            <p className="text-2xl font-bold text-orange-700">{dashboardStats?.template_stats.drafting || 0}</p>
          </Card>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Pool */}
          <Card className="p-0 overflow-hidden">
            <div className="p-4 border-b flex justify-between items-center">
              <h2 className="font-semibold">待受理咨询</h2>
              <Select
                value={urgency}
                onChange={(e) => setUrgency(e.target.value)}
                options={urgencyOptions}
                className="w-32"
              />
            </div>
            <Table columns={[
              { key: 'title', label: '标题' },
              { key: 'urgency', label: '紧急程度', width: '80px' },
              { key: 'action', label: '操作', width: '80px' },
            ]}>
              {poolConsultations.slice(0, 5).map((c) => (
                <tr key={c.id} className="border-b last:border-b-0">
                  <td className="px-4 py-3">
                    <p className="font-medium truncate">{c.title}</p>
                    <p className="text-xs text-gray-500">{c.user_name} - {c.user_department}</p>
                  </td>
                  <td className="px-4 py-3">
                    <span className={`px-2 py-1 rounded text-xs ${
                      c.urgency === 'urgent' ? 'bg-red-100 text-red-700' :
                      c.urgency === 'high' ? 'bg-orange-100 text-orange-700' :
                      'bg-gray-100 text-gray-700'
                    }`}>
                      {c.urgency === 'urgent' ? '紧急' : c.urgency === 'high' ? '高' : '普通'}
                    </span>
                  </td>
                  <td className="px-4 py-3">
                    <Button size="sm" onClick={() => handleAccept(c.id)}>受理</Button>
                  </td>
                </tr>
              ))}
            </Table>
            {poolConsultations.length > 5 && (
              <div className="p-3 border-t text-center">
                <Link to="/legal/pool" className="text-blue-600 hover:underline text-sm">
                  查看全部 {poolConsultations.length} 个
                </Link>
              </div>
            )}
          </Card>

          {/* My Tasks */}
          <Card className="p-0 overflow-hidden">
            <div className="p-4 border-b">
              <h2 className="font-semibold">我的任务</h2>
            </div>
            <Table columns={[
              { key: 'title', label: '标题' },
              { key: 'status', label: '状态', width: '100px' },
              { key: 'action', label: '操作', width: '80px' },
            ]}>
              {myTasks.slice(0, 5).map((c) => (
                <tr key={c.id} className="border-b last:border-b-0">
                  <td className="px-4 py-3">
                    <p className="font-medium truncate">{c.title}</p>
                    <p className="text-xs text-gray-500">{c.ticket_no}</p>
                  </td>
                  <td className="px-4 py-3">
                    <span className={`px-2 py-1 rounded text-xs ${
                      c.status === 'accepted' ? 'bg-blue-100 text-blue-700' :
                      c.status === 'in_progress' ? 'bg-yellow-100 text-yellow-700' :
                      c.status === 'waiting_supplement' ? 'bg-orange-100 text-orange-700' :
                      'bg-gray-100 text-gray-700'
                    }`}>
                      {statusLabels[c.status]}
                    </span>
                  </td>
                  <td className="px-4 py-3">
                    <Link to={`/consultation/${c.id}`}>
                      <Button variant="ghost" size="sm">处理</Button>
                    </Link>
                  </td>
                </tr>
              ))}
            </Table>
            {myTasks.length > 5 && (
              <div className="p-3 border-t text-center">
                <Link to="/legal/tasks" className="text-blue-600 hover:underline text-sm">
                  查看全部 {myTasks.length} 个
                </Link>
              </div>
            )}
          </Card>
        </div>

        {/* Quick Links */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mt-6">
          <Link to="/legal/pool">
            <Card className="p-4 text-center hover:shadow-md transition-shadow cursor-pointer">
              <p className="text-lg font-semibold">咨询池</p>
              <p className="text-gray-500 text-sm">受理新咨询</p>
            </Card>
          </Link>
          <Link to="/legal/tasks">
            <Card className="p-4 text-center hover:shadow-md transition-shadow cursor-pointer">
              <p className="text-lg font-semibold">我的任务</p>
              <p className="text-gray-500 text-sm">处理中的咨询</p>
            </Card>
          </Link>
          <Link to="/legal/templates">
            <Card className="p-4 text-center hover:shadow-md transition-shadow cursor-pointer">
              <p className="text-lg font-semibold">模板审批</p>
              <p className="text-gray-500 text-sm">待审批 {dashboardStats?.template_stats.pending_approval || 0}</p>
            </Card>
          </Link>
          <Link to="/statistics">
            <Card className="p-4 text-center hover:shadow-md transition-shadow cursor-pointer">
              <p className="text-lg font-semibold">统计分析</p>
              <p className="text-gray-500 text-sm">查看数据报表</p>
            </Card>
          </Link>
        </div>
      </main>
    </div>
  );
}
