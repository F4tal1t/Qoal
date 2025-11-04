import React, { useState, useEffect } from 'react';
import anime from 'animejs/lib/anime.es.js';
import ProgressTracker from './ProgressTracker';
import GuestConversionTracker from './GuestConversionTracker';
import ConversionLimitBanner from './ConversionLimitBanner';
import './ConversionPage.css';
import axios from 'axios';

// Configure axios base URL
const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8000';
axios.defaults.baseURL = API_BASE_URL;
console.log('API Base URL:', API_BASE_URL);

const ConversionPage = ({ conversionType, onBack, user }) => {
  const [selectedFile, setSelectedFile] = useState(null);
  const [targetFormat, setTargetFormat] = useState('');
  const [jobId, setJobId] = useState(null);
  const [isUploading, setIsUploading] = useState(false);
  const [showRegistration, setShowRegistration] = useState(false);
  const [networkError, setNetworkError] = useState(false);

  // Test API connection on component mount
  useEffect(() => {
    const testConnection = async () => {
      try {
        console.log('Testing API connection...');
        const response = await axios.get('/api/conversions/test/');
        console.log('API connection successful:', response.data);
        setNetworkError(false);
      } catch (error) {
        console.error('API connection failed:', error.message);
        setNetworkError(true);
      }
    };
    testConnection();
  }, []);

  const isAuthenticated = !!user;
  
  const {
    guestStatus,
    updateAfterConversion,
    canConvert,
    remainingConversions
  } = GuestConversionTracker({
    onLimitReached: () => !isAuthenticated && setShowRegistration(true),
    onConversionUsed: (remaining) => {
      if (remaining <= 1 && !isAuthenticated) {
        setShowRegistration(true);
      }
    }
  });

  // Authenticated users can always convert
  const actualCanConvert = isAuthenticated || canConvert;
  const actualRemainingConversions = isAuthenticated ? 'unlimited' : remainingConversions;

  const formatOptions = {
    image: ['PNG', 'JPEG', 'WebP', 'HEIC', 'BMP', 'TIFF'],
    audio: ['MP3', 'WAV', 'FLAC', 'AAC', 'M4A', 'OGG'],
    video: ['MP4', 'AVI', 'MOV', 'WMV', 'WebM', 'MKV']
  };

  const handleFileSelect = (event) => {
    const file = event.target.files[0];
    setSelectedFile(file);
    
    // Trigger animation after state update
    setTimeout(() => {
      const fileInfo = document.querySelector('.file-info');
      if (fileInfo) {
        anime({
          targets: fileInfo,
          opacity: [0, 1],
          translateY: [20, 0],
          duration: 600,
          easing: 'easeOutCubic'
        });
      }
    }, 100);
  };

  const handleConvert = async () => {
    if (!selectedFile || !targetFormat) return;
    
    if (!actualCanConvert) {
      setShowRegistration(true);
      return;
    }
    
    setIsUploading(true);
    
    try {
      const formData = new FormData();
      formData.append('file', selectedFile);
      formData.append('target_format', targetFormat);
      formData.append('conversion_type', conversionType);
      
      console.log('Uploading file:', selectedFile.name, 'to format:', targetFormat);
      console.log('User authenticated:', isAuthenticated);
      
      const headers = {
        'Content-Type': 'multipart/form-data'
      };
      
      if (isAuthenticated) {
        headers['Authorization'] = `Bearer ${localStorage.getItem('access_token')}`;
      }
      
      const response = await axios.post('/api/conversions/upload/', formData, { headers });
      
      console.log('Upload response:', response.data);
      
      // Update conversion count for guests
      if (!isAuthenticated) {
        updateAfterConversion();
      }
      
      setJobId(response.data.job_id);
      setIsUploading(false);
      
      // Show success animation
      anime({
        targets: '.conversion-success',
        opacity: [0, 1],
        scale: [0.8, 1],
        duration: 600,
        easing: 'easeOutBack'
      });
      
    } catch (error) {
      console.error('Upload failed:', error.response?.data || error.message);
      if (error.response?.status === 429) {
        setShowRegistration(true);
      }
      setIsUploading(false);
    }
  };

  return (
    <div className="conversion-page">
      <div className="conversion-header">
        <button className="back-btn" onClick={onBack}>← Back</button>
        <h2>{conversionType.toUpperCase()} Conversion</h2>
      </div>

      {networkError && (
        <div style={{
          background: 'rgba(239, 68, 68, 0.1)',
          color: '#dc2626',
          padding: '12px',
          borderRadius: '8px',
          marginBottom: '16px',
          textAlign: 'center'
        }}>
          ⚠️ Cannot connect to backend API. Please check if the server is running.
        </div>
      )}

      {!isAuthenticated && (
        <ConversionLimitBanner
          remainingConversions={actualRemainingConversions}
          isAuthenticated={isAuthenticated}
          onRegisterClick={() => setShowRegistration(true)}
        />
      )}

      <div className="conversion-content">
        <div className="upload-section">
          <div className="file-drop-zone">
            <input
              type="file"
              id="file-input"
              onChange={handleFileSelect}
              accept={conversionType === 'image' ? 'image/*' : 
                     conversionType === 'audio' ? 'audio/*' : 'video/*'}
              hidden
            />
            <label htmlFor="file-input" className="drop-label">
              <span className="upload-icon">📁</span>
              <span>Click to select {conversionType} file</span>
            </label>
          </div>

          {selectedFile && (
            <div className="file-info">
              <p><strong>File:</strong> {selectedFile.name}</p>
              <p><strong>Size:</strong> {Math.round(selectedFile.size / 1024)} KB</p>
            </div>
          )}
        </div>

        <div className="format-section">
          <h3>Convert to:</h3>
          <div className="format-options">
            {formatOptions[conversionType]?.map(format => (
              <button
                key={format}
                className={`format-btn ${targetFormat === format ? 'active' : ''}`}
                onClick={() => setTargetFormat(format)}
              >
                {format}
              </button>
            ))}
          </div>
        </div>

        <button 
          className={`start-conversion-btn ${(!actualCanConvert || networkError) ? 'disabled' : ''}`}
          onClick={handleConvert}
          disabled={!selectedFile || !targetFormat || isUploading || !actualCanConvert || networkError}
        >
          {networkError ? 'Server Offline' :
           isUploading ? 'Uploading...' : 
           !actualCanConvert ? 'Register to Convert' : 
           'Start Conversion'}
        </button>

        {jobId && <ProgressTracker jobId={jobId} />}
        
        {/* Debug Panel */}
        <div style={{
          marginTop: '20px',
          padding: '10px',
          background: 'rgba(0,0,0,0.05)',
          borderRadius: '8px',
          fontSize: '12px',
          fontFamily: 'monospace'
        }}>
          <strong>Debug Info:</strong><br/>
          API URL: {API_BASE_URL}<br/>
          Network: {networkError ? '❌ Offline' : '✅ Online'}<br/>
          Auth: {isAuthenticated ? '✅ Logged In' : '❌ Guest'}<br/>
          Can Convert: {actualCanConvert ? '✅ Yes' : '❌ No'}<br/>
          {!isAuthenticated && `Remaining: ${actualRemainingConversions}`}
        </div>
        
        {showRegistration && (
          <div className="registration-prompt">
            <div className="prompt-content">
              <h3>Ready to unlock unlimited conversions?</h3>
              <p>Register for free to get unlimited daily conversions and access to premium features.</p>
              <div className="prompt-actions">
                <button className="register-btn-prompt">Register Free</button>
                <button 
                  className="maybe-later-btn"
                  onClick={() => setShowRegistration(false)}
                >
                  Maybe Later
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default ConversionPage;