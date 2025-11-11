import React from 'react';

const Loader: React.FC = () => {
  return (
    <div style={{
      position: 'fixed',
      top: 0,
      left: 0,
      width: '100%',
      height: '100%',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      backgroundColor: '#0a0e17'
    }}>
      <div className="halftone-bg" style={{ 
        position: 'fixed', 
        top: 0, 
        left: 0, 
        width: '100%', 
        height: '100%', 
        zIndex: 0
      }}>
        <div className="halftone-noise" />
        <div style={{
          position: 'absolute',
          top: 0,
          left: 0,
          width: '100%',
          height: '100%',
          backgroundImage: 'url(/BayerDithering.png)',
          backgroundRepeat: 'repeat',
          backgroundSize: '20px 20px',
          opacity: 0.25,
          mixBlendMode: 'overlay'
        }} />
      </div>
      <div style={{
        position: 'relative',
        zIndex: 1,
        fontSize: '2rem',
        color: '#ffb947',
        fontWeight: 600
      }}>Loading...</div>
    </div>
  );
};

export default Loader;
