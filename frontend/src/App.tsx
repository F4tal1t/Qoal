import React, { Suspense } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Landing from './pages/Landing';
import Auth from './pages/Auth';
import Convert from './pages/Convert';
import Layout from './components/Layout';
import Loader from './components/Loader';

function App() {
  return (
    <Router>
      <Suspense fallback={<Loader />}>
        <Routes>
          <Route path="/" element={<Landing />} />
          <Route path="/auth" element={<Auth />} />
          <Route path="/convert" element={<Layout><Convert /></Layout>} />
        </Routes>
      </Suspense>
    </Router>
  );
}

export default App;