import{ Suspense } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Landing from './pages/Landing';
import Auth from './pages/Auth';
import Convert from './pages/Convert';
import Layout from './components/Layout';
import Loader from './components/Loader';
import ProtectedRoute from './components/ProtectedRoute';

function App() {
  return (
    <Router>
      <Suspense fallback={<Loader />}>
        <Routes>
          <Route path="/" element={<Landing />} />
          <Route path="/auth" element={<Auth />} />
          <Route path="/convert" element={
            <ProtectedRoute>
              <Layout><Convert /></Layout>
            </ProtectedRoute>
          } />
        </Routes>
      </Suspense>
    </Router>
  );
}

export default App;