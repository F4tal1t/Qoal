import React, { useState, useEffect } from 'react';
import anime from 'animejs/lib/anime.es.js';
import './ProgressTracker.css';

const ProgressTracker = ({ jobId }) => {
  const [jobStatus, setJobStatus] = useState(null);
  const [progress, setProgress] = useState(0);

  useEffect(() => {
    if (jobId) {
      pollJobStatus();
    }
  }, [jobId]);

  useEffect(() => {
    const progressBar = document.querySelector('.progress-fill');
    
    if (progressBar) {
      // Animate progress bar
      anime({
        targets: progressBar,
        width: `${progress}%`,
        duration: 800,
        easing: 'easeOutCubic'
      });
    }
  }, [progress]);

  const pollJobStatus = async () => {
    try {
      const response = await fetch(`${process.env.REACT_APP_API_URL || 'http://localhost:8000'}/api/conversions/job/${jobId}/status/`);
      const data = await response.json();
      
      setJobStatus(data);
      setProgress(data.progress || 0);
      
      // Continue polling if not completed or failed
      if (data.status === 'processing' || data.status === 'pending') {
        setTimeout(pollJobStatus, 2000);
      }
    } catch (error) {
      console.error('Error polling job status:', error);
      // Retry after error
      setTimeout(pollJobStatus, 5000);
    }
  };

  if (!jobStatus) {
    return (
      <div className="progress-tracker loading">
        <div className="loading-spinner"></div>
        <p>Loading job status...</p>
      </div>
    );
  }

  return (
    <div className="progress-tracker">
      <div className="conversion-info">
        <h3>{jobStatus.file_info?.original_filename}</h3>
        <p className="format-info">
          {jobStatus.file_info?.source_format?.toUpperCase()} → {jobStatus.file_info?.target_format?.toUpperCase()}
        </p>
      </div>
      
      <div className="progress-container">
        <div className="progress-bar">
          <div className="progress-fill" style={{width: '0%'}}></div>
          <div className="progress-dither"></div>
        </div>
        <span className="progress-text">{progress}%</span>
      </div>
      
      <div className="status-info">
        <div className="status-badge" data-status={jobStatus.status}>
          {jobStatus.current_stage || jobStatus.status}
        </div>
        {jobStatus.estimated_time_remaining && (
          <p className="eta">ETA: {jobStatus.estimated_time_remaining}s</p>
        )}
      </div>

      {jobStatus.status === 'completed' && jobStatus.download_url && (
        <div className="download-section">
          <a 
            href={jobStatus.download_url} 
            className="download-btn"
            download
          >
            Download Result
          </a>
          <p className="file-size">
            Size: {Math.round(jobStatus.output_file_size / 1024)} KB
          </p>
        </div>
      )}

      {jobStatus.status === 'failed' && (
        <div className="error-section">
          <p className="error-message">
            Conversion failed. Please try again.
          </p>
        </div>
      )}
    </div>
  );
};

export default ProgressTracker;