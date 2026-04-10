import { useState, useEffect } from 'react';
import { Link, useParams } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { api } from '../../lib/api';
import type { Consultation } from '../../types';
import { Card } from '../ui/Card';
import { Button } from '../ui/Button';
import { Textarea } from '../ui/Textarea';
import { Modal } from '../ui/Modal';

const statusLabels: Record<string, string> = {
  pending: '待处理',
  accepted: '已受理',
  in_progress: '处理中',
  waiting_supplement: '等待补充',
  completed: '已完成',
  closed: '已关闭',
};

export function ConsultationDetail() {
  const { id } = useParams<{ id: string }>();
  const { user } = useAuth();
  const [consultation, setConsultation] = useState<Consultation | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [replyContent, setReplyContent] = useState('');
  const [isInternal, setIsInternal] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isRateModalOpen, setIsRateModalOpen] = useState(false);
  const [rating, setRating] = useState(5);
  const [showSupplementModal, setShowSupplementModal] = useState(false);
  const [supplementContent, setSupplementContent] = useState('');

  useEffect(() => {
    loadConsultation();
  }, [id]);

  const loadConsultation = async () => {
    if (!id) return;
    setIsLoading(true);
    try {
      const data = await api.getConsultation(id);
      setConsultation(data);
    } catch (error) {
      console.error('Failed to load consultation:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleReply = async () => {
    if (!id || !replyContent.trim()) return;
    setIsSubmitting(true);
    try {
      await api.replyConsultation(id, replyContent, isInternal);
      setReplyContent('');
      setIsInternal(false);
      loadConsultation();
    } catch (error) {
      console.error('Failed to reply:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleRequestSupplement = async () => {
    if (!id || !supplementContent.trim()) return;
    setIsSubmitting(true);
    try {
      await api.requestSupplement(id, supplementContent);
      setSupplementContent('');
      setShowSupplementModal(false);
      loadConsultation();
    } catch (error) {
      console.error('Failed to request supplement:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleClose = async () => {
    if (!id) return;
    try {
      await api.closeConsultation(id);
      loadConsultation();
    } catch (error) {
      console.error('Failed to close:', error);
    }
  };

  const handleRate = async () => {
    if (!id) return;
    try {
      await api.rateConsultation(id, rating);
      setIsRateModalOpen(false);
      loadConsultation();
    } catch (error) {
      console.error('Failed to rate:', error);
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <p className="text-gray-500">加载中...</p>
      </div>
    );
  }

  if (!consultation) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <p className="text-gray-500">咨询不存在</p>
      </div>
    );
  }

  const isEmployee = user?.role === 'employee';
  const isLegalStaff = ['legal_staff', 'legal_head'].includes(user?.role || '');

  return (
    <div className="min-h-screen bg-gray-100">
      <header className="bg-white shadow">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex justify-between items-center">
            <div>
              <Link to={isEmployee ? '/' : '/legal'} className="text-blue-600 hover:underline mb-2 inline-block">
                返回列表
              </Link>
              <h1 className="text-xl font-bold">{consultation.title}</h1>
              <p className="text-gray-500 text-sm">
                工单号: {consultation.ticket_no} |
                创建时间: {new Date(consultation.created_at).toLocaleString('zh-CN')}
              </p>
            </div>
            <div className="flex gap-2">
              {consultation.status === 'completed' && isEmployee && !consultation.rating && (
                <Button onClick={() => setIsRateModalOpen(true)}>评价</Button>
              )}
              {consultation.status === 'completed' && isEmployee && (
                <Button variant="outline" onClick={handleClose}>关闭</Button>
              )}
            </div>
          </div>
        </div>
      </header>

      <main className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Info Card */}
        <Card className="mb-6">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <p className="text-gray-500 text-sm">咨询类型</p>
              <p className="font-medium">{consultation.category}</p>
            </div>
            <div>
              <p className="text-gray-500 text-sm">紧急程度</p>
              <p className="font-medium">{consultation.urgency}</p>
            </div>
            <div>
              <p className="text-gray-500 text-sm">当前状态</p>
              <p className="font-medium">{statusLabels[consultation.status]}</p>
            </div>
            <div>
              <p className="text-gray-500 text-sm">发起人</p>
              <p className="font-medium">{consultation.user_name} ({consultation.user_department})</p>
            </div>
            {consultation.legal_staff_name && (
              <div className="col-span-2">
                <p className="text-gray-500 text-sm">处理法务</p>
                <p className="font-medium">{consultation.legal_staff_name}</p>
              </div>
            )}
          </div>
        </Card>

        {/* Description */}
        <Card className="mb-6">
          <h3 className="font-semibold mb-2">问题描述</h3>
          <p className="text-gray-700 whitespace-pre-wrap">{consultation.description}</p>
        </Card>

        {/* Replies */}
        <Card className="mb-6">
          <h3 className="font-semibold mb-4">沟通记录</h3>
          <div className="space-y-4">
            {consultation.replies?.map((reply) => (
              <div
                key={reply.id}
                className={`p-4 rounded-lg ${
                  reply.is_internal ? 'bg-yellow-50 border border-yellow-200' : 'bg-gray-50'
                }`}
              >
                <div className="flex justify-between items-start mb-2">
                  <div>
                    <span className="font-medium">{reply.user_name}</span>
                    {reply.is_internal && (
                      <span className="ml-2 px-2 py-0.5 bg-yellow-200 text-yellow-800 text-xs rounded">
                        内部
                      </span>
                    )}
                  </div>
                  <span className="text-gray-400 text-sm">
                    {new Date(reply.created_at).toLocaleString('zh-CN')}
                  </span>
                </div>
                <p className="text-gray-700 whitespace-pre-wrap">{reply.content}</p>
              </div>
            ))}
            {(!consultation.replies || consultation.replies.length === 0) && (
              <p className="text-gray-500 text-center py-4">暂无沟通记录</p>
            )}
          </div>
        </Card>

        {/* Reply Form (Legal Staff) */}
        {isLegalStaff && consultation.status !== 'closed' && consultation.status !== 'completed' && (
          <Card className="mb-6">
            <h3 className="font-semibold mb-4">回复</h3>
            <Textarea
              value={replyContent}
              onChange={(e) => setReplyContent(e.target.value)}
              placeholder="请输入回复内容..."
              rows={4}
            />
            <div className="flex justify-between items-center mt-4">
              <label className="flex items-center gap-2">
                <input
                  type="checkbox"
                  checked={isInternal}
                  onChange={(e) => setIsInternal(e.target.checked)}
                />
                <span className="text-sm">内部回复（员工不可见）</span>
              </label>
              <div className="flex gap-2">
                <Button
                  variant="outline"
                  onClick={() => setShowSupplementModal(true)}
                >
                  要求补充资料
                </Button>
                <Button
                  onClick={handleReply}
                  disabled={!replyContent.trim() || isSubmitting}
                >
                  发送回复
                </Button>
              </div>
            </div>
          </Card>
        )}

        {/* Rating */}
        {consultation.rating && (
          <Card className="mb-6">
            <h3 className="font-semibold mb-2">满意度评价</h3>
            <div className="flex items-center gap-1">
              {[1, 2, 3, 4, 5].map((star) => (
                <span
                  key={star}
                  className={`text-2xl ${star <= consultation.rating! ? 'text-yellow-400' : 'text-gray-300'}`}
                >
                  *
                </span>
              ))}
              <span className="ml-2 text-gray-600">
                {consultation.rating === 5 ? '非常满意' :
                 consultation.rating === 4 ? '满意' :
                 consultation.rating === 3 ? '一般' :
                 consultation.rating === 2 ? '不满意' : '非常不满意'}
              </span>
            </div>
          </Card>
        )}
      </main>

      {/* Rate Modal */}
      <Modal
        isOpen={isRateModalOpen}
        onClose={() => setIsRateModalOpen(false)}
        title="评价服务"
      >
        <div className="space-y-4">
          <p className="text-gray-600">请对本次法律咨询服务进行评价：</p>
          <div className="flex justify-center gap-2">
            {[1, 2, 3, 4, 5].map((star) => (
              <button
                key={star}
                onClick={() => setRating(star)}
                className={`text-4xl transition-colors ${
                  star <= rating ? 'text-yellow-400' : 'text-gray-300'
                }`}
              >
                *
              </button>
            ))}
          </div>
          <div className="flex justify-end gap-2 pt-4">
            <Button variant="outline" onClick={() => setIsRateModalOpen(false)}>
              取消
            </Button>
            <Button onClick={handleRate}>提交评价</Button>
          </div>
        </div>
      </Modal>

      {/* Supplement Modal */}
      <Modal
        isOpen={showSupplementModal}
        onClose={() => setShowSupplementModal(false)}
        title="要求补充资料"
      >
        <Textarea
          label="请说明需要补充的内容"
          value={supplementContent}
          onChange={(e) => setSupplementContent(e.target.value)}
          placeholder="请详细说明需要补充哪些资料..."
          rows={4}
        />
        <div className="flex justify-end gap-2 pt-4">
          <Button variant="outline" onClick={() => setShowSupplementModal(false)}>
            取消
          </Button>
          <Button
            onClick={handleRequestSupplement}
            disabled={!supplementContent.trim() || isSubmitting}
          >
            发送
          </Button>
        </div>
      </Modal>
    </div>
  );
}
