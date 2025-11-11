import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../services/api';
import { Upload } from '../components/animate-ui/icons/upload';
import { Paperclip } from '../components/animate-ui/icons/paperclip';
import { LoaderPinwheel } from '../components/animate-ui/icons/loader-pinwheel';
import { CircleCheckBig } from '../components/animate-ui/icons/circle-check-big';
import { MessageSquareWarning } from '../components/animate-ui/icons/message-square-warning';

type ConversionType = 'image' | 'document' | 'audio' | 'video' | 'archive';

const Convert: React.FC = () => {
  const [activeType, setActiveType] = useState<ConversionType>('image');
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [targetFormat, setTargetFormat] = useState('');
  const [isConverting, setIsConverting] = useState(false);
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setSelectedFile(e.target.files[0]);
    }
  };

  const handleConvert = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedFile || !targetFormat) return;

    setIsConverting(true);
    setError('');

    try {
      const response = await api.jobs.upload(selectedFile, targetFormat);
      navigate(`/status/${response.job_id}`);
    } catch (err: any) {
      setError(err.message || 'Conversion failed');
      setIsConverting(false);
    }
  };

  const formats: Record<ConversionType, string[]> = {
    image: ['jpg', 'png', 'webp', 'gif', 'bmp', 'tiff', 'svg'],
    document: ['pdf', 'docx', 'txt', 'xlsx', 'csv', 'rtf', 'odt'],
    audio: ['mp3', 'wav', 'flac', 'aac', 'ogg', 'm4a', 'wma'],
    video: ['mp4', 'avi', 'mov', 'wmv', 'flv', 'mkv', 'webm', 'm4v'],
    archive: ['zip', 'rar', '7z', 'tar', 'gz', 'bz2', 'xz']
  };

  const categoryIcons: Record<ConversionType, string> = {
    image: '/ImgIco.gif',
    document: '/DocIco.gif',
    audio: '/AudIco.gif',
    video: '/VidIco.gif',
    archive: '/ArcIco.gif'
  };

  return (
    <div className="convert-page">
      <div className="category-toggles">
        {(['image', 'document', 'audio', 'video', 'archive'] as ConversionType[]).map((cat) => (
          <button
            key={cat}
            className={`category-btn ${activeType === cat ? 'active' : ''}`}
            onClick={() => {
              setActiveType(cat);
              setSelectedFile(null);
              setTargetFormat('');
            }}
          >
            <img src={categoryIcons[cat]} alt={cat} />
            {cat.charAt(0).toUpperCase() + cat.slice(1)}
          </button>
        ))}
      </div>

      <div className="convert-container">
        <h1>Convert {activeType.charAt(0).toUpperCase() + activeType.slice(1)} Files</h1>

        <form onSubmit={handleConvert}>
          <div className="file-upload-area">
            <input
              type="file"
              id="file"
              onChange={handleFileChange}
              style={{ display: 'none' }}
            />
            <label htmlFor="file" className="upload-label">
              {selectedFile ? (
                <>
                  <Paperclip size={48} animateOnHover />
                  <p>{selectedFile.name}</p>
                  <span>{(selectedFile.size / 1024 / 1024).toFixed(2)} MB</span>
                </>
              ) : (
                <>
                  <Upload size={48} animateOnHover />
                  <p>Click to upload or drag and drop</p>
                </>
              )}
            </label>
          </div>

          <div className="form-group">
            <label htmlFor="targetFormat">Target Format</label>
            <select
              id="targetFormat"
              value={targetFormat}
              onChange={(e) => setTargetFormat(e.target.value)}
              required
            >
              <option value="">Select format</option>
              {formats[activeType].map((format) => (
                <option key={format} value={format}>
                  {format.toUpperCase()}
                </option>
              ))}
            </select>
          </div>

          {error && (
            <div className="error-message">
              <MessageSquareWarning size={20} />
              {error}
            </div>
          )}

          <button type="submit" disabled={isConverting || !selectedFile || !targetFormat}>
            {isConverting ? (
              <>
                <LoaderPinwheel size={20} animate loop />
                Converting...
              </>
            ) : (
              <>
                <CircleCheckBig size={20} animateOnHover />
                Convert File
              </>
            )}
          </button>
        </form>
      </div>
    </div>
  );
};

export default Convert;