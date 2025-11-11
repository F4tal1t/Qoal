import React from 'react';

interface NavbarProps {
  transparent?: boolean;
}

const Navbar: React.FC<NavbarProps> = ({ transparent = true }) => {
  return (
    <nav style={{
      position: 'fixed',
      top: 0,
      left: 0,
      width: '100%',
      zIndex: 100,
      padding: '1.5rem 2rem',
      backgroundColor: transparent ? 'transparent' : 'rgba(255, 255, 255, 0.95)',
      backdropFilter: transparent ? 'none' : 'blur(10px)',
      fontFamily: 'inherit'
    }}>
      <div style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        gap: '2rem'
      }}>
        <a href="/" style={{
          color: transparent ? '#fff' : '#333',
          textDecoration: 'none',
          fontSize: '1.1rem',
          fontWeight: 500,
          padding: '0.5rem 1rem',
          borderRadius: '0.5rem',
          transition: 'all 0.3s ease'
        }}>Home</a>
        
        {/* Convert Dropdown */}
        <div className="nav-dropdown" style={{
          position: 'relative'
        }}>
          <a href="#" style={{
            color: transparent ? '#fff' : '#333',
            textDecoration: 'none',
            fontSize: '1.1rem',
            fontWeight: 500,
            padding: '0.5rem 1rem',
            borderRadius: '0.5rem',
            transition: 'all 0.3s ease',
            display: 'flex',
            alignItems: 'center',
            gap: '0.3rem'
          }}>
            Convert
            <span style={{ fontSize: '0.7rem', marginTop: '2px' }}>▼</span>
          </a>
          
          <div className="dropdown-content" style={{
            position: 'absolute',
            top: '100%',
            left: '50%',
            transform: 'translateX(-50%)',
            backgroundColor: 'rgba(255, 255, 255, 0.95)',
            backdropFilter: 'blur(10px)',
            borderRadius: '0.5rem',
            padding: '0.5rem',
            minWidth: '200px',
            boxShadow: '0 4px 20px rgba(0, 0, 0, 0.1)',
            opacity: 0,
            visibility: 'hidden',
            transition: 'all 0.3s ease',
            marginTop: '0.5rem',
            border: '1px solid rgba(255, 255, 255, 0.2)'
          }}>
            <a href="/convert/image" style={{
              display: 'block',
              padding: '0.75rem 1rem',
              color: '#333',
              textDecoration: 'none',
              borderRadius: '0.25rem',
              transition: 'all 0.3s ease',
              fontSize: '0.95rem',
              fontWeight: 500
            }}>Image Conversions</a>
            <a href="/convert/video" style={{
              display: 'block',
              padding: '0.75rem 1rem',
              color: '#333',
              textDecoration: 'none',
              borderRadius: '0.25rem',
              transition: 'all 0.3s ease',
              fontSize: '0.95rem',
              fontWeight: 500
            }}>Video Conversions</a>
            <a href="/convert/audio" style={{
              display: 'block',
              padding: '0.75rem 1rem',
              color: '#333',
              textDecoration: 'none',
              borderRadius: '0.25rem',
              transition: 'all 0.3s ease',
              fontSize: '0.95rem',
              fontWeight: 500
            }}>Audio Conversions</a>
            <a href="/convert/document" style={{
              display: 'block',
              padding: '0.75rem 1rem',
              color: '#333',
              textDecoration: 'none',
              borderRadius: '0.25rem',
              transition: 'all 0.3s ease',
              fontSize: '0.95rem',
              fontWeight: 500
            }}>Document Conversions</a>
            <a href="/convert/archive" style={{
              display: 'block',
              padding: '0.75rem 1rem',
              color: '#333',
              textDecoration: 'none',
              borderRadius: '0.25rem',
              transition: 'all 0.3s ease',
              fontSize: '0.95rem',
              fontWeight: 500
            }}>Archive Conversions</a>
          </div>
        </div>
        
        <a href="/about" style={{
          color: transparent ? '#fff' : '#333',
          textDecoration: 'none',
          fontSize: '1.1rem',
          fontWeight: 500,
          padding: '0.5rem 1rem',
          borderRadius: '0.5rem',
          transition: 'all 0.3s ease'
        }}>About</a>
        
        <a href="/contact" style={{
          color: transparent ? '#fff' : '#333',
          textDecoration: 'none',
          fontSize: '1.1rem',
          fontWeight: 500,
          padding: '0.5rem 1rem',
          borderRadius: '0.5rem',
          transition: 'all 0.3s ease'
        }}>Contact</a>
      </div>
      
      <style>{`
        .nav-dropdown:hover .dropdown-content {
          opacity: 1;
          visibility: visible;
          transform: translateX(-50%) translateY(0);
        }
        
        .nav-dropdown .dropdown-content a:hover {
          background-color: rgba(255, 120, 90, 0.1);
          color: #ff785a;
        }
        
        nav a:hover {
          background-color: rgba(255, 255, 255, 0.1);
          transform: translateY(-1px);
        }
        
        .nav-dropdown:hover > a {
          background-color: rgba(255, 255, 255, 0.1);
        }
      `}</style>
    </nav>
  );
};

export default Navbar;