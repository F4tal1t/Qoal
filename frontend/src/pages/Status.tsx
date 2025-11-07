import React, { useState, useEffect } from 'react';

interface StatusProps {
  jobId?: string;
}

const Status: React.FC<StatusProps> = ({ jobId }) => {
  const [jobStatus] = useState({
    id: jobId || '1',
    status: 'processing',
    type: 'image',
    input_file: 'input.jpg',
    output_file: 'output.png',
    progress: 75,
    created_at: '2024-01-01T12:00:00Z',
    completed_at: null
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Fetch job status from API
    const fetchStatus = async () => {
      try {
        // const response = await fetch(`/api/status/${jobId}`);
        // const data = await response.json();
        // setJobStatus(data);
        
        // Mock data
        setTimeout(() => {
          setLoading(false);
        }, 1000);
      } catch (error) {
        console.error('Error fetching status:', error);
        setLoading(false);
      }
    };

    fetchStatus();

    // Poll for status updates
    const interval = setInterval(() => {
      fetchStatus();
    }, 5000);

    return () => clearInterval(interval);
  }, [jobId]);

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed':
        return 'text-green-600';
      case 'processing':
        return 'text-blue-600';
      case 'failed':
        return 'text-red-600';
      default:
        return 'text-gray-600';
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'completed':
        return '✅';
      case 'processing':
        return '⏳';
      case 'failed':
        return '❌';
      default:
        return '❓';
    }
  };

  if (loading) {
    return (
      <div className="flex-center min-h-screen">
        <div className="loading-container">
          <div className="spinner"></div>
          <p>Loading status...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex-center min-h-screen">
      <div className="status-container">
        <h1>Conversion Status</h1>
        
        <div className="status-card">
          <div className="status-header">
            <span className="status-icon">{getStatusIcon(jobStatus.status)}</span>
            <span className={`status-text ${getStatusColor(jobStatus.status)}`}>
              {jobStatus.status.charAt(0).toUpperCase() + jobStatus.status.slice(1)}
            </span>
          </div>

          <div className="status-details">
            <div className="detail-row">
              <span className="label">Job ID:</span>
              <span className="value">{jobStatus.id}</span>
            </div>
            <div className="detail-row">
              <span className="label">Type:</span>
              <span className="value">{jobStatus.type.toUpperCase()}</span>
            </div>
            <div className="detail-row">
              <span className="label">Input File:</span>
              <span className="value">{jobStatus.input_file}</span>
            </div>
            {jobStatus.output_file && (
              <div className="detail-row">
                <span className="label">Output File:</span>
                <span className="value">{jobStatus.output_file}</span>
              </div>
            )}
            <div className="detail-row">
              <span className="label">Created:</span>
              <span className="value">{new Date(jobStatus.created_at).toLocaleString()}</span>
            </div>
          </div>

          {jobStatus.status === 'processing' && (
            <div className="progress-bar">
              <div className="progress-fill" style={{ width: `${jobStatus.progress}%` }}></div>
            </div>
          )}

          <div className="status-actions">
            {jobStatus.status === 'completed' && jobStatus.output_file && (
              <a href={`/download/${jobStatus.id}`} className="download-btn">
                Download File
              </a>
            )}
            <a href="/" className="back-btn">
              Back to Home
            </a>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Status;