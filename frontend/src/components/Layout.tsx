import React from 'react';
import { Link, useLocation } from 'react-router-dom';

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const location = useLocation();
  const isLandingPage = location.pathname === '/';

  return (
    <div className="app-container">
      {!isLandingPage && (
        <nav className="app-navbar">
          <div className="nav-container">
            <Link to="/" className="logo">FileProcessor</Link>
            <div className="nav-links">
              <Link to="/convert/image">Convert Image</Link>
              <Link to="/convert/document">Convert Document</Link>
              <Link to="/convert/audio">Convert Audio</Link>
              <Link to="/convert/video">Convert Video</Link>
              <Link to="/convert/archive">Convert Archive</Link>
            </div>
            <div className="nav-actions">
              <Link to="/profile">Profile</Link>
              <button className="logout-btn">Logout</button>
            </div>
          </div>
        </nav>
      )}
      
      <main className="main-content">
        {children}
      </main>
    </div>
  );
};

export default Layout;