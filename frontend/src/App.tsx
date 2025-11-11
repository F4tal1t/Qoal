import React, { Suspense } from 'react';
import { BrowserRouter as Router, Routes, Route, useParams } from 'react-router-dom';
import Landing from './pages/Landing';
import Auth from './pages/Auth';
import Convert from './pages/Convert';
import Status from './pages/Status';
import Profile from './pages/Profile';
import Layout from './components/Layout';
import Loader from './components/Loader';

function App() {
  return (
    <Router>
      <Suspense fallback={<Loader />}>
        <Routes>
          <Route path="/" element={<Landing />} />
          <Route path="/auth" element={<Auth />} />
          <Route path="/*" element={<AppLayout />} />
        </Routes>
      </Suspense>
    </Router>
  );
}

// AppLayout wraps pages with navigation (except landing page)
const AppLayout: React.FC = () => {
  return (
    <Layout>
      <Routes>
        <Route path="/convert" element={<Convert />} />
        <Route path="/status/:id" element={<StatusWrapper />} />
        <Route path="/profile" element={<Profile />} />
      </Routes>
    </Layout>
  );
};

const StatusWrapper: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  return <Status jobId={id} />;
};

export default App;