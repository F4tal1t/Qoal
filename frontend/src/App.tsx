import React from 'react';
import { BrowserRouter as Router, Routes, Route, useParams } from 'react-router-dom';
import Landing from './pages/Landing';
import Login from './pages/Login';
import Register from './pages/Register';

import Convert from './pages/Convert';
import Status from './pages/Status';
import Profile from './pages/Profile';
import Layout from './components/Layout';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Landing />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/*" element={<AppLayout />} />
      </Routes>
    </Router>
  );
}

// AppLayout wraps pages with navigation (except landing page)
const AppLayout: React.FC = () => {
  return (
    <Layout>
      <Routes>
        <Route path="/convert/:type" element={<ConvertWrapper />} />
        <Route path="/status/:id" element={<StatusWrapper />} />
        <Route path="/profile" element={<Profile />} />
      </Routes>
    </Layout>
  );
};

// Wrapper components to handle route parameters
const ConvertWrapper: React.FC = () => {
  const { type } = useParams<{ type: string }>();
  return <Convert type={type as any} />;
};

const StatusWrapper: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  return <Status jobId={id} />;
};

export default App;