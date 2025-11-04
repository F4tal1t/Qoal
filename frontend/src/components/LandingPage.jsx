import React, { useState, useEffect } from 'react';
import anime from 'animejs/lib/anime.es.js';
import './LandingPage.css';
import backgroundVideo from '../background-video.mp4';
import ConversionPage from './ConversionPage';
import AuthPage from './AuthPage';
import ApiTest from './ApiTest';

const LandingPage = () => {
  const [currentPage, setCurrentPage] = useState('landing');
  const [selectedConversionType, setSelectedConversionType] = useState(null);
  const [user, setUser] = useState(null);

  useEffect(() => {
    // Check if user is logged in
    const savedUser = localStorage.getItem('user');
    if (savedUser) {
      setUser(JSON.parse(savedUser));
    }
  }, []);

  useEffect(() => {
    // Animate title on load
    anime({
      targets: '.title',
      translateY: [-50, 0],
      opacity: [0, 1],
      duration: 1200,
      easing: 'easeOutCubic',
      delay: 300
    });

    // Animate tagline
    anime({
      targets: '.tagline',
      translateY: [30, 0],
      opacity: [0, 1],
      duration: 1000,
      easing: 'easeOutCubic',
      delay: 600
    });

    // Animate project description
    anime({
      targets: '.project-description',
      translateY: [20, 0],
      opacity: [0, 1],
      duration: 1000,
      easing: 'easeOutCubic',
      delay: 800
    });

    // Animate conversion bar
    anime({
      targets: '.conversion-bar',
      scale: [0.9, 1],
      opacity: [0, 1],
      duration: 800,
      easing: 'easeOutBack',
      delay: 1000
    });

    // Add click handlers for option buttons
    document.querySelectorAll('.option-btn').forEach(btn => {
      btn.addEventListener('click', () => {
        document.querySelectorAll('.option-btn').forEach(b => b.classList.remove('active'));
        btn.classList.add('active');
        setSelectedConversionType(btn.dataset.type);
      });
    });

    // Folder animations removed - folder is in video
  }, []);

  const handleAuthSuccess = (userData) => {
    setUser(userData);
    setCurrentPage('landing');
  };

  const handleLogout = () => {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('user');
    setUser(null);
  };

  if (currentPage === 'conversion') {
    return (
      <ConversionPage 
        conversionType={selectedConversionType}
        onBack={() => setCurrentPage('landing')}
        user={user}
      />
    );
  }

  if (currentPage === 'auth') {
    return (
      <AuthPage 
        onBack={() => setCurrentPage('landing')}
        onAuthSuccess={handleAuthSuccess}
      />
    );
  }

  if (currentPage === 'test') {
    return (
      <div>
        <button onClick={() => setCurrentPage('landing')} style={{ margin: '20px' }}>← Back</button>
        <ApiTest />
      </div>
    );
  }

  return (
    <div className="landing-container">
      {/* Background Video */}
      <video 
        className="bg-video" 
        autoPlay 
        muted 
        loop 
        playsInline
        src={backgroundVideo}
        onLoadedData={() => console.log('Video loaded successfully')}
        onError={(e) => console.error('Video failed to load:', e)}
      />

      {/* Dithered Overlay */}
      <div className="dither-overlay"></div>
      <div className="moving-shapes"></div>

      {/* Navigation */}
      <nav className="nav-bar">
        <div className="nav-home" onClick={() => setCurrentPage('test')} style={{ cursor: 'pointer' }}>
          <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
            <path d="M10 20v-6h4v6h5v-8h3L12 3 2 12h3v8z"/>
          </svg>
        </div>
        {user ? (
          <div className="user-menu">
            <span className="user-name">Hi, {user.username}</span>
            <button className="logout-btn" onClick={handleLogout}>Logout</button>
          </div>
        ) : (
          <button 
            className="register-btn"
            onClick={() => setCurrentPage('auth')}
          >
            Register / Login
          </button>
        )}
      </nav>

      {/* Main Content */}
      <main className="main-content">
        <div className="content-left">
          <h1 className="title">Qoal</h1>
          <p className="tagline">
            Cloud File Processing<br />
            Made Simple
          </p>
          <p className="project-description">
            Transform your media files with professional-grade quality. 
            Convert images, process audio, compress videos, and handle documents 
            with lightning-fast cloud processing.
          </p>
          <div className="conversion-bar">
            <div className="conversion-options">
              <button className="option-btn active" data-type="image">
                <span className="option-icon">🖼️</span>
                <span className="option-text">Image</span>
              </button>
              <button className="option-btn" data-type="audio">
                <span className="option-icon">🎵</span>
                <span className="option-text">Audio</span>
              </button>
              <button className="option-btn" data-type="video">
                <span className="option-icon">🎬</span>
                <span className="option-text">Video</span>
              </button>
            </div>
            <button 
              className="convert-now-btn"
              onClick={() => {
                if (selectedConversionType) {
                  setCurrentPage('conversion');
                }
              }}
            >
              Convert Now
            </button>
          </div>
        </div>

        <div className="content-right">
          {/* Folder removed - present in video */}
        </div>
      </main>


    </div>
  );
};

export default LandingPage;