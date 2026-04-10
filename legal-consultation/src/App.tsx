import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './context/AuthContext';
import {
  LoginPage,
  EmployeeDashboard,
  ConsultationDetail,
  LegalDashboard,
  StatisticsPage,
} from './components/pages';

// Protected Route wrapper
function ProtectedRoute({ children, allowedRoles }: { children: React.ReactNode; allowedRoles?: string[] }) {
  const { isAuthenticated, isLoading, user } = useAuth();

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <p className="text-gray-500">加载中...</p>
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  if (allowedRoles && user && !allowedRoles.includes(user.role)) {
    // Redirect based on role
    if (user.role === 'employee') {
      return <Navigate to="/" replace />;
    }
    if (['legal_staff', 'legal_head'].includes(user.role)) {
      return <Navigate to="/legal" replace />;
    }
    if (user.role === 'admin') {
      return <Navigate to="/admin" replace />;
    }
    return <Navigate to="/" replace />;
  }

  return <>{children}</>;
}

// Home redirect based on role
function HomeRedirect() {
  const { user, isLoading, isAuthenticated } = useAuth();

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <p className="text-gray-500">加载中...</p>
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  if (!user) return null;

  switch (user.role) {
    case 'employee':
    case 'supervisor':
      return <Navigate to="/dashboard" replace />;
    case 'legal_staff':
    case 'legal_head':
      return <Navigate to="/legal" replace />;
    case 'admin':
      return <Navigate to="/admin" replace />;
    default:
      return <Navigate to="/dashboard" replace />;
  }
}

function AppRoutes() {
  const { user } = useAuth();

  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />

      {/* Employee routes */}
      <Route
        path="/"
        element={
          <ProtectedRoute allowedRoles={['employee', 'supervisor', 'legal_staff', 'legal_head', 'admin']}>
            <HomeRedirect />
          </ProtectedRoute>
        }
      />
      <Route
        path="/dashboard"
        element={
          <ProtectedRoute allowedRoles={['employee', 'supervisor']}>
            <EmployeeDashboard />
          </ProtectedRoute>
        }
      />

      {/* Legal staff routes */}
      <Route
        path="/legal"
        element={
          <ProtectedRoute allowedRoles={['legal_staff', 'legal_head']}>
            <LegalDashboard />
          </ProtectedRoute>
        }
      />

      {/* Shared routes */}
      <Route
        path="/consultation/:id"
        element={
          <ProtectedRoute>
            <ConsultationDetail />
          </ProtectedRoute>
        }
      />

      {/* Statistics */}
      <Route
        path="/statistics"
        element={
          <ProtectedRoute allowedRoles={['legal_staff', 'legal_head', 'admin']}>
            <StatisticsPage />
          </ProtectedRoute>
        }
      />

      {/* Admin routes - placeholder for now */}
      <Route
        path="/admin"
        element={
          <ProtectedRoute allowedRoles={['admin']}>
            <div className="min-h-screen bg-gray-100 p-8">
              <h1 className="text-2xl font-bold">管理后台</h1>
              <p className="text-gray-500 mt-2">管理员功能开发中...</p>
            </div>
          </ProtectedRoute>
        }
      />

      {/* Catch all */}
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <AppRoutes />
      </AuthProvider>
    </BrowserRouter>
  );
}

export default App;
