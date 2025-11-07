import React, { useState, useEffect } from 'react';
import StaggeredMenu from '../components/StaggeredMenu';

const Landing: React.FC = () => {
  const [conversionCount, setConversionCount] = useState(0);

  useEffect(() => {
    // Fetch conversion count from localStorage or API
    const savedCount = localStorage.getItem('conversionCount');
    if (savedCount) {
      setConversionCount(parseInt(savedCount));
    }
  }, []);

  const menuItems = [
    { label: 'Home', ariaLabel: 'Go to home page', link: '/' },
    { label: 'Conversions', ariaLabel: 'View conversions', link: '/conversions' },
    { label: 'Pricing', ariaLabel: 'View pricing', link: '/pricing' },
    { label: 'About', ariaLabel: 'Learn about us', link: '/about' },
    { label: 'Contact', ariaLabel: 'Get in touch', link: '/contact' }
  ];

  const socialItems = [
    { label: 'Twitter', link: 'https://twitter.com' },
    { label: 'GitHub', link: 'https://github.com' },
    { label: 'LinkedIn', link: 'https://linkedin.com' }
  ];

  return (
    <div className="min-h-screen">
      {/* Hero Section */}
      <section className="hero-section" style={{ height: '100vh', background: '#1a1a1a', position: 'relative' }}>
        {/* Header Logo */}
        <div style={{ 
          position: 'absolute', 
          top: '20px', 
          left: '20px', 
          zIndex: 10,
          display: 'flex',
          alignItems: 'center',
          gap: '10px'
        }}>
          <img src="/Qoalation.png" alt="Qoal Logo" style={{ height: '40px' }} />
          <img src="/QoalText.png" alt="Qoal Text" style={{ height: '30px' }} />
        </div>

        <StaggeredMenu
          position="right"
          items={menuItems}
          socialItems={socialItems}
          displaySocials={true}
          displayItemNumbering={true}
          menuButtonColor="#fff"
          openMenuButtonColor="#000"
          changeMenuColorOnOpen={true}
          colors={['#B19EEF', '#5227FF']}
          logoUrl="/Qoalation.png"
          accentColor="#ff6b6b"
          isFixed={false}
          onMenuOpen={() => console.log('Menu opened')}
          onMenuClose={() => console.log('Menu closed')}
        />
        
        <div className="hero-content">
          <h1 className="general-title">Convert Your Files</h1>
          <p>Fast, secure, and easy file conversion for all your needs</p>
          <button className="get-started-btn">Get Started</button>
        </div>
      </section>

      {/* Conversion Counter */}
      <section className="conversion-counter">
        <div className="counter-container">
          <h2>Files Converted</h2>
          <div className="counter-number">{conversionCount}</div>
          <p>Start converting your files now!</p>
        </div>
      </section>

      {/* Conversion Types Sections */}
      <section className="conversion-sections">
        <div className="section image-conversions">
          <h2>Image Conversions</h2>
          <p>Convert between JPG, PNG, WebP, GIF, BMP, TIFF</p>
          <a href="/convert/image">Convert Images</a>
        </div>

        <div className="section document-conversions">
          <h2>Document Conversions</h2>
          <p>Convert between PDF, DOCX, TXT, RTF, ODT</p>
          <a href="/convert/document">Convert Documents</a>
        </div>

        <div className="section audio-conversions">
          <h2>Audio Conversions</h2>
          <p>Convert between MP3, WAV, FLAC, M4A, OGG</p>
          <a href="/convert/audio">Convert Audio</a>
        </div>

        <div className="section video-conversions">
          <h2>Video Conversions</h2>
          <p>Convert between MP4, AVI, MOV, WebM, MKV</p>
          <a href="/convert/video">Convert Videos</a>
        </div>

        <div className="section archive-conversions">
          <h2>Archive Conversions</h2>
          <p>Convert between ZIP, RAR, 7Z, TAR, GZ</p>
          <a href="/convert/archive">Convert Archives</a>
        </div>
      </section>

      {/* Footer */}
      <footer className="footer">
        <div className="footer-content">
          <p>&copy; 2024 FileProcessor. All rights reserved.</p>
          <div className="footer-links">
            <a href="/privacy">Privacy Policy</a>
            <a href="/terms">Terms of Service</a>
          </div>
        </div>
      </footer>
    </div>
  );
};

export default Landing;