import React from 'react';

interface NavbarProps {
  transparent?: boolean;
}

const Navbar: React.FC<NavbarProps> = ({ transparent = true }) => {
  return (
    <nav style={{
      position: 'fixed',
      top: '1rem',
      left: '50%',
      transform: 'translateX(-50%)',
      zIndex: 100,
      backgroundColor: '#161b27',
      borderRadius: '1rem',
      padding: '0.75rem 2rem',
      fontFamily: 'inherit'
    }}>
      <div style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        gap: '2rem'
      }}>
        <a href="/" style={{
          color: '#fff',
          textDecoration: 'none',
          fontSize: '1.1rem',
          fontWeight: 500,
          padding: '0.5rem 1rem',
          borderRadius: '0.5rem',
          transition: 'all 0.3s ease'
        }}>Home</a>
        
        <a href="/convert" style={{
          color: '#fff',
          textDecoration: 'none',
          fontSize: '1.1rem',
          fontWeight: 500,
          padding: '0.5rem 1rem',
          borderRadius: '0.5rem',
          transition: 'all 0.3s ease'
        }}>Convert</a>
        
        <a href="/login" style={{
          color: '#161b27',
          backgroundColor: '#ffb947',
          textDecoration: 'none',
          fontSize: '1.1rem',
          fontWeight: 500,
          padding: '0.5rem 1rem',
          borderRadius: '0.5rem',
          transition: 'all 0.3s ease'
        }}>Login</a>
        
        <a href="/signup" style={{
          color: '#161b27',
          backgroundColor: '#ffb947',
          textDecoration: 'none',
          fontSize: '1.1rem',
          fontWeight: 500,
          padding: '0.5rem 1rem',
          borderRadius: '0.5rem',
          transition: 'all 0.3s ease'
        }}>Signup</a>
      </div>
      
      <style>{`
        nav a:hover {
          opacity: 0.8;
          transform: translateY(-1px);
        }
      `}</style>
    </nav>
  );
};

export default Navbar;