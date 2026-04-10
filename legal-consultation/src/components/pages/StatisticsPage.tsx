import { useState, useEffect } from 'react';
import { useAuth } from '../../context/AuthContext';
import { api } from '../../lib/api';
import type { OverviewStats, CategoryDistribution, ProcessingEfficiency, StaffWorkload } from '../../types';
import { Card } from '../ui/Card';
import { Button } from '../ui/Button';
import { Table } from '../ui/Table';

export function StatisticsPage() {
  const { user, logout } = useAuth();
  const [startDate, setStartDate] = useState(() => {
    const d = new Date();
    d.setMonth(d.getMonth() - 1);
    return d.toISOString().split('T')[0];
  });
  const [endDate, setEndDate] = useState(() => new Date().toISOString().split('T')[0]);

  const [overview, setOverview] = useState<OverviewStats | null>(null);
  const [categoryDist, setCategoryDist] = useState<CategoryDistribution[]>([]);
  const [efficiency, setEfficiency] = useState<ProcessingEfficiency[]>([]);
  const [workload, setWorkload] = useState<StaffWorkload[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    loadStatistics();
  }, [startDate, endDate]);

  const loadStatistics = async () => {
    setIsLoading(true);
    try {
      const [ov, cat, eff, wl] = await Promise.all([
        api.getOverviewStats(startDate, endDate),
        api.getCategoryDistribution(startDate, endDate),
        api.getProcessingEfficiency(startDate, endDate),
        api.getStaffWorkload(startDate, endDate),
      ]);
      setOverview(ov);
      setCategoryDist(cat);
      setEfficiency(eff);
      setWorkload(wl);
    } catch (error) {
      console.error('Failed to load statistics:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleExport = async () => {
    try {
      const result = await api.exportReport(startDate, endDate);
      // Open the export URL in a new window for download
      window.open(result.url, '_blank');
    } catch (error) {
      console.error('Failed to export:', error);
    }
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
          <h1 className="text-xl font-bold">统计分析</h1>
          <div className="flex items-center gap-4">
            <span className="text-gray-600">{user?.name}</span>
            <Button variant="outline" size="sm" onClick={logout}>退出</Button>
          </div>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Filters */}
        <Card className="mb-6">
          <div className="flex flex-wrap gap-4 items-end">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">开始日期</label>
              <input
                type="date"
                value={startDate}
                onChange={(e) => setStartDate(e.target.value)}
                className="border rounded px-3 py-2"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">结束日期</label>
              <input
                type="date"
                value={endDate}
                onChange={(e) => setEndDate(e.target.value)}
                className="border rounded px-3 py-2"
              />
            </div>
            <Button onClick={handleExport}>导出报表</Button>
          </div>
        </Card>

        {/* Overview Stats */}
        <div className="grid grid-cols-1 md:grid-cols-5 gap-4 mb-6">
          <Card className="p-4">
            <p className="text-gray-500 text-sm">总咨询量</p>
            <p className="text-2xl font-bold">{overview?.total_consultations || 0}</p>
          </Card>
          <Card className="p-4">
            <p className="text-gray-500 text-sm">待处理</p>
            <p className="text-2xl font-bold text-yellow-600">{overview?.pending_consultations || 0}</p>
          </Card>
          <Card className="p-4">
            <p className="text-gray-500 text-sm">已完成</p>
            <p className="text-2xl font-bold text-green-600">{overview?.completed_consultations || 0}</p>
          </Card>
          <Card className="p-4">
            <p className="text-gray-500 text-sm">平均响应时间</p>
            <p className="text-2xl font-bold">{overview?.avg_response_time || 0}h</p>
          </Card>
          <Card className="p-4">
            <p className="text-gray-500 text-sm">满意度</p>
            <p className="text-2xl font-bold text-blue-600">
              {((overview?.satisfaction_rate || 0) * 100).toFixed(1)}%
            </p>
          </Card>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Category Distribution */}
          <Card className="p-0 overflow-hidden">
            <div className="p-4 border-b">
              <h2 className="font-semibold">咨询类型分布</h2>
            </div>
            <Table columns={[
              { key: 'category', label: '类型' },
              { key: 'count', label: '数量', width: '100px' },
              { key: 'percentage', label: '占比', width: '100px' },
            ]}>
              {categoryDist.map((c) => (
                <tr key={c.category} className="border-b last:border-b-0">
                  <td className="px-4 py-3">{c.category}</td>
                  <td className="px-4 py-3 font-medium">{c.count}</td>
                  <td className="px-4 py-3 text-gray-500">{(c.percentage * 100).toFixed(1)}%</td>
                </tr>
              ))}
            </Table>
          </Card>

          {/* Staff Workload */}
          <Card className="p-0 overflow-hidden">
            <div className="p-4 border-b">
              <h2 className="font-semibold">法务工作量</h2>
            </div>
            <Table columns={[
              { key: 'staff_name', label: '姓名' },
              { key: 'total_tasks', label: '总任务', width: '80px' },
              { key: 'completed_tasks', label: '已完成', width: '80px' },
              { key: 'avg_time', label: '平均用时(h)', width: '100px' },
            ]}>
              {workload.map((w) => (
                <tr key={w.staff_id} className="border-b last:border-b-0">
                  <td className="px-4 py-3">{w.staff_name}</td>
                  <td className="px-4 py-3 font-medium">{w.total_tasks}</td>
                  <td className="px-4 py-3 text-green-600">{w.completed_tasks}</td>
                  <td className="px-4 py-3 text-gray-500">{w.avg_completion_time.toFixed(1)}</td>
                </tr>
              ))}
            </Table>
          </Card>
        </div>

        {/* Processing Efficiency */}
        <Card className="mt-6 p-0 overflow-hidden">
          <div className="p-4 border-b">
            <h2 className="font-semibold">每日处理效率</h2>
          </div>
          <Table columns={[
            { key: 'date', label: '日期', width: '150px' },
            { key: 'count', label: '处理数量', width: '120px' },
            { key: 'avg_time', label: '平均用时(h)', width: '120px' },
          ]}>
            {efficiency.map((e) => (
              <tr key={e.date} className="border-b last:border-b-0">
                <td className="px-4 py-3">{e.date}</td>
                <td className="px-4 py-3 font-medium">{e.count}</td>
                <td className="px-4 py-3 text-gray-500">{e.avg_time.toFixed(1)}</td>
              </tr>
            ))}
          </Table>
        </Card>
      </main>
    </div>
  );
}
