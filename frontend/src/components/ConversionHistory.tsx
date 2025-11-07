import React, { useState, useEffect } from 'react';

interface ConversionJob {
  id: string;
  type: string;
  status: string;
  input_file: string;
  output_file?: string;
  created_at: string;
}

const ConversionHistory: React.FC = () => {
  const [conversions, setConversions] = useState<ConversionJob[]>([]);
  const [showHistory, setShowHistory] = useState(false);

  useEffect(() => {
    // Load conversion history from localStorage
    const savedConversions = localStorage.getItem('conversionHistory');
    if (savedConversions) {
      setConversions(JSON.parse(savedConversions));
    }
  }, []);

  if (conversions.length === 0) {
    return null;
  }

  return (
    <div className="conversion-history">
      <button 
        onClick={() => setShowHistory(!showHistory)}
        className="history-toggle"
      >
        {showHistory ? 'Hide' : 'Show'} Recent Conversions ({conversions.length})
      </button>
      
      {showHistory && (
        <div className="history-list">
          {conversions.map((conversion) => (
            <div key={conversion.id} className="history-item">
              <div className="history-info">
                <span className="file-name">{conversion.input_file}</span>
                <span className={`status ${conversion.status}`}>
                  {conversion.status}
                </span>
              </div>
              {conversion.status === 'completed' && conversion.output_file && (
                <a href={`/download/${conversion.id}`} className="download-link">
                  Download
                </a>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default ConversionHistory;