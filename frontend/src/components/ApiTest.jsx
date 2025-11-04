import React, { useState, useEffect } from 'react';
import axios from 'axios';

axios.defaults.baseURL = process.env.REACT_APP_API_URL || 'http://localhost:8000';

const ApiTest = () => {
  const [testResults, setTestResults] = useState({});

  useEffect(() => {
    testApis();
  }, []);

  const testApis = async () => {
    const tests = [
      { name: 'Backend Connection', url: '/api/conversions/test/' },
      { name: 'Guest Status', url: '/api/conversions/guest-status/' },
      { name: 'Auth Register', url: '/api/auth/register/', method: 'POST', data: { email: 'test@test.com', password: 'test123', username: 'testuser' } }
    ];

    const results = {};
    
    for (const test of tests) {
      try {
        const response = test.method === 'POST' 
          ? await axios.post(test.url, test.data)
          : await axios.get(test.url);
        
        results[test.name] = { 
          status: 'SUCCESS', 
          data: response.data,
          statusCode: response.status 
        };
      } catch (error) {
        results[test.name] = { 
          status: 'ERROR', 
          error: error.response?.data || error.message,
          statusCode: error.response?.status || 'Network Error'
        };
      }
    }
    
    setTestResults(results);
  };

  return (
    <div style={{ padding: '20px', fontFamily: 'monospace', fontSize: '12px' }}>
      <h3>API Connection Test</h3>
      <button onClick={testApis} style={{ marginBottom: '20px' }}>Retest APIs</button>
      
      {Object.entries(testResults).map(([name, result]) => (
        <div key={name} style={{ 
          marginBottom: '15px', 
          padding: '10px', 
          border: '1px solid #ccc',
          backgroundColor: result.status === 'SUCCESS' ? '#d4edda' : '#f8d7da'
        }}>
          <strong>{name}</strong> - {result.status} ({result.statusCode})
          <pre style={{ fontSize: '10px', marginTop: '5px' }}>
            {JSON.stringify(result.data || result.error, null, 2)}
          </pre>
        </div>
      ))}
    </div>
  );
};

export default ApiTest;