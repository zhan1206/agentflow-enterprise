import React from 'react';
import { ConfigProvider } from 'antd';
import { Provider } from 'react-redux';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import zhCN from 'antd/locale/zh_CN';

import { store } from './store';
import { MainLayout } from './components/layout/MainLayout';
import { Dashboard } from './pages/Dashboard';
import { Agents } from './pages/Agents';
import { Tasks } from './pages/Tasks';
import { Workflows } from './pages/Workflows';
import { Observability } from './pages/Observability';
import { Settings } from './pages/Settings';

import './App.css';

const App: React.FC = () => {
  return (
    <ConfigProvider locale={zhCN}>
      <Provider store={store}>
        <BrowserRouter>
          <MainLayout>
            <Routes>
              <Route path="/" element={<Navigate to="/dashboard" replace />} />
              <Route path="/dashboard" element={<Dashboard />} />
              <Route path="/agents/*" element={<Agents />} />
              <Route path="/tasks/*" element={<Tasks />} />
              <Route path="/workflows/*" element={<Workflows />} />
              <Route path="/observability/*" element={<Observability />} />
              <Route path="/settings/*" element={<Settings />} />
            </Routes>
          </MainLayout>
        </BrowserRouter>
      </Provider>
    </ConfigProvider>
  );
};

export default App;
