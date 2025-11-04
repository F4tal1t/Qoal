import React from 'react';

const ConversionLimitBanner = ({ remainingConversions, isAuthenticated, onRegisterClick }) => {
  if (isAuthenticated) return null;

  const getBannerStyle = () => {
    const baseStyle = {
      background: 'rgba(255, 255, 255, 0.9)',
      borderLeft: '4px solid #D4AF37',
      borderRadius: '8px',
      padding: '16px',
      marginBottom: '20px',
      backdropFilter: 'blur(10px)',
      fontFamily: 'dotemp-demo, monospace'
    };
    
    if (remainingConversions === 0) {
      return { ...baseStyle, borderLeftColor: '#ef4444', background: 'rgba(254, 242, 242, 0.9)' };
    } else if (remainingConversions === 1) {
      return { ...baseStyle, borderLeftColor: '#f59e0b', background: 'rgba(255, 251, 235, 0.9)' };
    }
    return baseStyle;
  };

  const getButtonStyle = () => {
    const baseStyle = {
      background: '#D4AF37',
      color: 'white',
      border: 'none',
      padding: '8px 16px',
      borderRadius: '6px',
      fontFamily: 'dotemp-demo, monospace',
      fontSize: '0.85rem',
      fontWeight: '600',
      cursor: 'pointer',
      whiteSpace: 'nowrap'
    };
    
    return remainingConversions === 0 
      ? { ...baseStyle, background: '#ef4444' }
      : baseStyle;
  };

  const getMessage = () => {
    if (remainingConversions === 0) {
      return 'Free conversions used up! Register for unlimited conversions.';
    } else if (remainingConversions === 1) {
      return 'Last free conversion! Register to continue converting files.';
    }
    return `${remainingConversions} free conversions remaining today.`;
  };

  return (
    <div style={getBannerStyle()}>
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: '16px' }}>
        <div style={{ flex: 1 }}>
          <p style={{ fontSize: '0.9rem', fontWeight: '600', color: '#333', marginBottom: '4px' }}>
            {getMessage()}
          </p>
          {remainingConversions > 0 && (
            <p style={{ fontSize: '0.8rem', color: '#666' }}>
              Register now for unlimited daily conversions and premium features.
            </p>
          )}
        </div>
        <button
          onClick={onRegisterClick}
          style={getButtonStyle()}
        >
          {remainingConversions === 0 ? 'Register Now' : 'Sign Up Free'}
        </button>
      </div>
    </div>
  );
};

export default ConversionLimitBanner;