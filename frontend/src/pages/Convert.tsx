import React, { useState, useEffect, useRef } from 'react';
import { api } from '../services/api';
import type { JobResponse, JobStatus } from '../services/api';
import { Upload } from '../components/animate-ui/icons/upload';
import { Paperclip } from '../components/animate-ui/icons/paperclip';
import { LoaderPinwheel } from '../components/animate-ui/icons/loader-pinwheel';
import { CircleCheckBig } from '../components/animate-ui/icons/circle-check-big';
import { MessageSquareWarning } from '../components/animate-ui/icons/message-square-warning';
import { Tabs, TabsList, TabsTrigger, TabsHighlight, TabsHighlightItem } from '../components/animate-ui/components/animate/tabs';

type ConversionType = 'image' | 'document' | 'audio' | 'video' | 'archive';

const Convert: React.FC = () => {
  const [activeType, setActiveType] = useState<ConversionType>('image');
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [filePreview, setFilePreview] = useState<string | null>(null);
  const [targetFormat, setTargetFormat] = useState('');
  const [isConverting, setIsConverting] = useState(false);
  const [error, setError] = useState('');
  const [isDragging, setIsDragging] = useState(false);
  const [conversionResult, setConversionResult] = useState<JobResponse | null>(null);
  const [jobStatus, setJobStatus] = useState<JobStatus | null>(null);
  const pollTimerRef = useRef<number | null>(null);
  
  const user = JSON.parse(localStorage.getItem('user') || '{}');
  const userName = user.name || 'User';

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setSelectedFile(e.target.files[0]);
      setFilePreview(null);
    }
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      const file = e.dataTransfer.files[0];
      const ext = '.' + file.name.split('.').pop()?.toLowerCase();
      
      const isValidType = (() => {
        switch (activeType) {
          case 'image':
            return file.type.startsWith('image/');
          case 'document':
            return ['.pdf', '.doc', '.docx', '.txt', '.xlsx', '.csv', '.rtf', '.odt'].includes(ext);
          case 'audio':
            return file.type.startsWith('audio/');
          case 'video':
            return file.type.startsWith('video/');
          case 'archive':
            return ['.zip', '.rar', '.7z', '.tar', '.gz', '.bz2', '.xz'].includes(ext);
          default:
            return false;
        }
      })();
      
      if (!isValidType) {
        setError(`Invalid file type for ${activeType} conversion`);
        return;
      }
      
      setError('');
      setSelectedFile(file);
      setFilePreview(null);
    }
  };

  useEffect(() => {
    return () => {
      if (pollTimerRef.current) {
        clearInterval(pollTimerRef.current);
      }
    };
  }, []);

  const handleConvert = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedFile || !targetFormat) return;

    setIsConverting(true);
    setError('');
    setConversionResult(null);
    setJobStatus(null);

    try {
      const response = await api.jobs.upload(selectedFile, targetFormat);
      console.log('Upload response:', response);
      setConversionResult(response);
      
      // Start polling immediately
      if (response.job_id) {
        startPolling(response.job_id);
      }
    } catch (err) {
      console.error('Upload error:', err);
      setError(err instanceof Error ? err.message : 'Conversion failed');
      setIsConverting(false);
    }
  };

  const startPolling = (jobId: string) => {
    if (pollTimerRef.current) {
      clearInterval(pollTimerRef.current);
    }
    
    const poll = () => {
      api.jobs.getStatus(jobId)
        .then(status => {
          console.log('Status:', status.status);
          setJobStatus(prev => {
            if (prev?.status !== status.status) {
              return status;
            }
            return prev;
          });
          
          if (status.status === 'completed' || status.status === 'failed') {
            setIsConverting(false);
            if (pollTimerRef.current) {
              clearInterval(pollTimerRef.current);
              pollTimerRef.current = null;
            }
          }
        })
        .catch(err => console.error('Poll error:', err));
    };

    poll();
    pollTimerRef.current = window.setInterval(poll, 2000);
    console.log('Polling started, timer ID:', pollTimerRef.current);
  };

  const handleDownload = async () => {
    if (jobStatus?.job_id && jobStatus.status === 'completed') {
      try {
        const filename = `converted_${jobStatus.original_filename.split('.')[0]}.${jobStatus.target_format}`;
        await api.jobs.download(jobStatus.job_id, filename);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Download failed');
      }
    }
  };

  const handleReset = () => {
    if (pollTimerRef.current) {
      clearInterval(pollTimerRef.current);
      pollTimerRef.current = null;
    }
    setSelectedFile(null);
    setFilePreview(null);
    setTargetFormat('');
    setConversionResult(null);
    setJobStatus(null);
    setError('');
    setIsConverting(false);
  };

  const formats: Record<ConversionType, string[]> = {
    image: ['jpg', 'png', 'webp', 'gif', 'bmp', 'tiff', 'svg'],
    document: ['pdf', 'docx', 'txt', 'xlsx', 'csv', 'rtf', 'odt'],
    audio: ['mp3', 'wav', 'flac', 'aac', 'ogg', 'm4a', 'wma'],
    video: ['mp4', 'avi', 'mov', 'wmv', 'flv', 'mkv', 'webm', 'm4v'],
    archive: ['zip', 'rar', '7z', 'tar', 'gz', 'bz2', 'xz']
  };

  const acceptedTypes: Record<ConversionType, string> = {
    image: 'image/*',
    document: '.pdf,.doc,.docx,.txt,.xlsx,.csv,.rtf,.odt',
    audio: 'audio/*',
    video: 'video/*',
    archive: '.zip,.rar,.7z,.tar,.gz,.bz2,.xz'
  };

  const categoryIcons: Record<ConversionType, string> = {
    image: '/ImgIco.gif',
    document: '/DocIco.gif',
    audio: '/AudIco.gif',
    video: '/VidIco.gif',
    archive: '/ArcIco.gif'
  };

  return (
    <div className="fixed inset-0 flex items-center justify-center">
      <div className="halftone-bg fixed inset-0 z-0">
        <div className="halftone-noise" />
        <div className="absolute inset-0 opacity-25 mix-blend-overlay" style={{
          backgroundImage: 'url(/BayerDithering.png)',
          backgroundRepeat: 'repeat',
          backgroundSize: '20px 20px'
        }} />
      </div>
      
      <div className="relative z-10 w-full max-w-4xl mx-auto p-8 rounded-lg" style={{
        background: 'rgba(255, 255, 255, 0.05)',
        border: '1px solid rgba(255, 255, 255, 0.1)',
        backdropFilter: 'blur(10px)'
      }}>
        <div className="mb-6">
          <h2 className="text-xl font-display" style={{ color: 'var(--color-text)' }}>Welcome, {userName}!</h2>
        </div>
        
        <div className="flex gap-6">
          <Tabs defaultValue="image" className="flex-shrink-0" value={activeType} onValueChange={(val: string) => {
            setActiveType(val as ConversionType);
            setSelectedFile(null);
            setFilePreview(null);
            setTargetFormat('');
          }}>
            <TabsHighlight className="absolute z-0 inset-0 border border-[#E08A00] rounded-md shadow-sm bg-[#E08A00]">
              <TabsList className="flex-col h-full w-32 gap-2 py-4">
                {(['image', 'document', 'audio', 'video', 'archive'] as ConversionType[]).map((cat) => (
                  <TabsHighlightItem key={cat} value={cat}>
                    <TabsTrigger value={cat} className="w-full flex-col h-20 gap-2">
                      <img src={categoryIcons[cat]} alt={cat} style={{ width: '48px', height: '48px' }} />
                      <span className="text-sm font-medium font-display">{cat.charAt(0).toUpperCase() + cat.slice(1)}</span>
                    </TabsTrigger>
                  </TabsHighlightItem>
                ))}
              </TabsList>
            </TabsHighlight>
          </Tabs>

          <div className="flex-1 flex flex-col">
            <h1 className="text-2xl font-bold mb-6" style={{ color: 'var(--color-text)' }}>Convert {activeType.charAt(0).toUpperCase() + activeType.slice(1)} Files</h1>

            <form onSubmit={handleConvert} className="flex-1 flex flex-col">
              <div 
                className="p-6 rounded-lg border border-dashed" 
                style={{
                  background: isDragging ? 'rgba(255, 255, 255, 0.1)' : 'rgba(255, 255, 255, 0.05)',
                  borderColor: isDragging ? '#E08A00' : 'rgba(255, 255, 255, 0.2)'
                }}
                onDragOver={handleDragOver}
                onDragLeave={handleDragLeave}
                onDrop={handleDrop}
              >
                <input
                  type="file"
                  id="file"
                  onChange={handleFileChange}
                  accept={acceptedTypes[activeType]}
                  style={{ display: 'none' }}
                />
                <label htmlFor="file" className="cursor-pointer flex flex-col items-center gap-3">
                  {selectedFile ? (
                    <>
                      {filePreview && (
                        <img src={filePreview} alt="Preview" className="max-w-full max-h-48 rounded-md mb-2" />
                      )}
                      <Paperclip size={32} animateOnHover className="text-white" />
                      <p style={{ color: 'var(--color-text)' }}>{selectedFile.name}</p>
                      <span style={{ color: 'var(--color-text)', opacity: 0.7 }}>{(selectedFile.size / 1024 / 1024).toFixed(2)} MB</span>
                    </>
                  ) : (
                    <>
                      <Upload size={48} animateOnHover className="text-white" />
                      <p style={{ color: 'var(--color-text)' }}>Click to upload or drag and drop</p>
                    </>
                  )}
                </label>
              </div>

              <div className="mt-4">
                <label htmlFor="targetFormat" className="block text-sm mb-2" style={{ color: 'var(--color-text)' }}>Target Format</label>
                <select
                  id="targetFormat"
                  value={targetFormat}
                  onChange={(e) => setTargetFormat(e.target.value)}
                  required
                  className="w-full px-4 py-2 rounded-md border text-[#cae2e2]"
                  style={{
                    background: 'rgba(255, 255, 255, 0.05)',
                    border: '1px solid rgba(255, 255, 255, 0.1)'
                  }}
                >
                  <option value="" className="bg-[#161B27] text-[#cae2e2]">Select format</option>
                  {formats[activeType].map((format) => (
                    <option key={format} value={format} className="bg-[#161B27] text-[#cae2e2]">
                      {format.toUpperCase()}
                    </option>
                  ))}
                </select>
              </div>

              <div className="flex-1" />

              {error && (
                <div className="flex items-center gap-2 text-red-500 text-sm mb-4">
                  <MessageSquareWarning size={20} />
                  {error}
                </div>
              )}

              {conversionResult && (
                <div className="mt-6 p-6 rounded-lg" style={{
                  background: jobStatus?.status === 'completed' ? 'rgba(0, 255, 0, 0.1)' : 'rgba(255, 165, 0, 0.1)',
                  border: jobStatus?.status === 'completed' ? '1px solid rgba(0, 255, 0, 0.3)' : '1px solid rgba(255, 165, 0, 0.3)'
                }}>
                  {jobStatus?.status === 'failed' ? (
                    <>
                      <div className="flex items-center gap-3 mb-4">
                        <MessageSquareWarning size={32} className="text-red-500" />
                        <h3 className="text-lg font-semibold" style={{ color: 'var(--color-text)' }}>Conversion Failed</h3>
                      </div>
                      <p className="text-red-500 mb-4">{jobStatus.error || 'Unknown error'}</p>
                      <button type="button" onClick={handleReset} className="w-full py-2 rounded-md font-medium text-white" style={{ background: 'rgba(255, 255, 255, 0.1)' }}>Try Again</button>
                    </>
                  ) : jobStatus?.status === 'completed' ? (
                    <>
                      <div className="flex items-center gap-3 mb-4">
                        <CircleCheckBig size={32} animate className="text-green-500" />
                        <h3 className="text-lg font-semibold" style={{ color: 'var(--color-text)' }}>Conversion Complete!</h3>
                      </div>
                      <div className="space-y-2 mb-4" style={{ color: 'var(--color-text)', opacity: 0.8 }}>
                        <p>File: {jobStatus.original_filename}</p>
                        <p>Format: {jobStatus.source_format.toUpperCase()} → {jobStatus.target_format.toUpperCase()}</p>
                      </div>
                      <div className="flex gap-3">
                        <button type="button" onClick={handleDownload} className="flex-1 py-2 rounded-md font-medium bg-[#E08A00] text-white">Download</button>
                        <button type="button" onClick={handleReset} className="flex-1 py-2 rounded-md font-medium text-white" style={{ background: 'rgba(255, 255, 255, 0.1)' }}>Convert Another</button>
                      </div>
                    </>
                  ) : (
                    <>
                      <div className="flex items-center gap-3 mb-4">
                        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-orange-500"></div>
                        <h3 className="text-lg font-semibold" style={{ color: 'var(--color-text)' }}>Processing...</h3>
                      </div>
                      <div className="space-y-2" style={{ color: 'var(--color-text)', opacity: 0.8 }}>
                        <p>File: {conversionResult.original_filename}</p>
                        <p>Format: {conversionResult.source_format.toUpperCase()} → {conversionResult.target_format.toUpperCase()}</p>
                        <p>Status: {jobStatus?.status || 'pending'}</p>
                      </div>
                    </>
                  )}
                </div>
              )}
              
              {!conversionResult && (
                <div className="flex justify-center mt-6">
                  <button
                    type="submit"
                    disabled={isConverting || !selectedFile || !targetFormat}
                    className="max-w-md w-full py-3 rounded-md font-medium transition-colors bg-[#E08A00] text-white flex items-center justify-center gap-2"
                    style={{ opacity: (isConverting || !selectedFile || !targetFormat) ? 0.5 : 1 }}
                  >
                    {isConverting ? (
                      <>
                        <LoaderPinwheel size={20} loop/>
                        Converting...
                      </>
                    ) : (
                      <>
                        <CircleCheckBig size={20} animateOnHover />
                        Convert File
                      </>
                    )}
                  </button>
                </div>
              )}
            </form>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Convert;