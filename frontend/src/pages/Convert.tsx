import React, { useState } from 'react';
import ConversionHistory from '../components/ConversionHistory';

interface ConvertProps {
  type: 'image' | 'document' | 'audio' | 'video' | 'archive';
}

const Convert: React.FC<ConvertProps> = ({ type }) => {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [targetFormat, setTargetFormat] = useState('');
  const [isConverting, setIsConverting] = useState(false);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setSelectedFile(e.target.files[0]);
    }
  };

  const handleConvert = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedFile || !targetFormat) return;

    setIsConverting(true);
    
    // Handle conversion logic here
    const formData = new FormData();
    formData.append('file', selectedFile);
    formData.append('targetFormat', targetFormat);
    formData.append('type', type);

    try {
      // const response = await fetch('/api/process', {
      //   method: 'POST',
      //   body: formData,
      // });
      // const data = await response.json();
      // console.log('Conversion started:', data);
      
      // Mock conversion
      setTimeout(() => {
        setIsConverting(false);
      
      // Increment conversion counter
      const currentCount = parseInt(localStorage.getItem('conversionCount') || '0');
      const newCount = currentCount + 1;
      localStorage.setItem('conversionCount', newCount.toString());
      
      // Create conversion job
      const job = {
        id: Date.now().toString(),
        type: type,
        status: 'completed',
        input_file: selectedFile.name,
        output_file: selectedFile.name.replace(/\.[^/.]+$/, '') + '.' + targetFormat,
        created_at: new Date().toISOString()
      };
      
      // Add to conversion history
      const history = JSON.parse(localStorage.getItem('conversionHistory') || '[]');
      const updatedHistory = [job, ...history].slice(0, 10);
      localStorage.setItem('conversionHistory', JSON.stringify(updatedHistory));
      
      alert('Conversion completed successfully!');
      
      // Navigate to status page
      window.location.href = `/status/${job.id}`;
      }, 2000);
    } catch (error) {
      console.error('Conversion error:', error);
      setIsConverting(false);
    }
  };

  const getFormats = () => {
    const formats = {
      image: ['JPG', 'PNG', 'WebP', 'GIF', 'BMP', 'TIFF'],
      document: ['PDF', 'DOCX', 'TXT', 'RTF', 'ODT'],
      audio: ['MP3', 'WAV', 'FLAC', 'M4A', 'OGG'],
      video: ['MP4', 'AVI', 'MOV', 'WebM', 'MKV'],
      archive: ['ZIP', 'RAR', '7Z', 'TAR', 'GZ']
    };
    return formats[type] || [];
  };

  return (
    <div className="flex-center min-h-screen">
      <div className="convert-container">
        <h1>Convert {type.charAt(0).toUpperCase() + type.slice(1)} Files</h1>
        
        <form onSubmit={handleConvert}>
          <div className="form-group">
            <label htmlFor="file">Select File</label>
            <input
              type="file"
              id="file"
              onChange={handleFileChange}
              accept={`.${type === 'archive' ? 'zip,rar,7z,tar,gz' : type === 'document' ? 'pdf,docx,txt,rtf,odt' : type === 'audio' ? 'audio/*' : type === 'video' ? 'video/*' : 'image/*'}`}
              required
            />
            {selectedFile && (
              <p className="file-info">Selected: {selectedFile.name}</p>
            )}
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
              {getFormats().map((format) => (
                <option key={format} value={format.toLowerCase()}>
                  {format}
                </option>
              ))}
            </select>
          </div>

          <button type="submit" disabled={isConverting}>
            {isConverting ? 'Converting...' : 'Convert File'}
          </button>
        </form>

        <div className="conversion-info">
          <h3>Supported Formats</h3>
          <p>Input: {getFormats().join(', ')}</p>
          <p>Output: {getFormats().join(', ')}</p>
        </div>
        
        <ConversionHistory />
        
        <div className="back-link">
          <a href="/">Back to Home</a>
        </div>
      </div>
    </div>
  );
};

export default Convert;