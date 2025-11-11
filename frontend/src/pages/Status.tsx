import React, { useState, useEffect } from 'react';
import { api } from '../services/api';
import { LoaderPinwheel } from '../components/animate-ui/icons/loader-pinwheel';
import { CircleCheckBig } from '../components/animate-ui/icons/circle-check-big';
import { MessageSquareWarning } from '../components/animate-ui/icons/message-square-warning';
import { Download } from '../components/animate-ui/icons/download';

interface StatusProps {
  jobId?: string;
}

const Status: React.FC<StatusProps> = ({ jobId }) => {
  const [jobStatus, setJobStatus] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!jobId) return;

    const fetchStatus = async () => {
      try {
        const data = await api.jobs.getStatus(jobId);
        setJobStatus(data);
        setLoading(false);
      } catch (err: any) {
        setError(err.message);
        setLoading(false);
      }
    };

    fetchStatus();

    const interval = setInterval(() => {
      if (jobStatus?.status === 'pending' || jobStatus?.status === 'processing') {
        fetchStatus();
      }
    }, 2000);

    return () => clearInterval(interval);
  }, [jobId, jobStatus?.status]);

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'completed':
        return <CircleCheckBig size={48} animate />;
      case 'processing':
      case 'pending':
        return <LoaderPinwheel size={48} animate loop />;
      case 'failed':
        return <MessageSquareWarning size={48} animate />;
      default:
        return <LoaderPinwheel size={48} />;
    }
  };

  const handleDownload = async () => {
    if (jobId) {
      try {
        await api.jobs.download(jobId);
      } catch (err: any) {
        setError(err.message);
      }
    }
  };

  if (loading) {
    return (
      <div className="flex-center min-h-screen">
        <div className="loading-container">
          <LoaderPinwheel size={64} animate loop />
          <p>Loading status...</p>
        </div>
      </div>
    );
  }

  if (error || !jobStatus) {
    return (
      <div className="flex-center min-h-screen">
        <div className="error-container">
          <MessageSquareWarning size={64} animate />
          <p>{error || 'Job not found'}</p>
          <a href="/convert/image">Back to Convert</a>
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
            {getStatusIcon(jobStatus.status)}
            <h2>{jobStatus.status.charAt(0).toUpperCase() + jobStatus.status.slice(1)}</h2>
          </div>

          <div className="status-details">
            <div className="detail-row">
              <span className="label">Job ID:</span>
              <span className="value">{jobStatus.job_id}</span>
            </div>
            <div className="detail-row">
              <span className="label">File:</span>
              <span className="value">{jobStatus.original_filename}</span>
            </div>
            <div className="detail-row">
              <span className="label">Size:</span>
              <span className="value">{(jobStatus.file_size / 1024 / 1024).toFixed(2)} MB</span>
            </div>
            <div className="detail-row">
              <span className="label">Conversion:</span>
              <span className="value">{jobStatus.source_format.toUpperCase()} → {jobStatus.target_format.toUpperCase()}</span>
            </div>
            <div className="detail-row">
              <span className="label">Created:</span>
              <span className="value">{new Date(jobStatus.created_at).toLocaleString()}</span>
            </div>
            {jobStatus.error && (
              <div className="detail-row error">
                <span className="label">Error:</span>
                <span className="value">{jobStatus.error}</span>
              </div>
            )}
          </div>

          <div className="status-actions">
            {jobStatus.status === 'completed' && jobStatus.download_url && (
              <button onClick={handleDownload} className="download-btn">
                <Download size={20} animateOnHover />
                Download File
              </button>
            )}
            <a href="/convert/image" className="back-btn">
              New Conversion
            </a>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Status;