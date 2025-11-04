import React, { useState, useEffect } from 'react';
import axios from 'axios';
import anime from 'animejs/lib/anime.es.js';
import './AuthPage.css';

axios.defaults.baseURL = process.env.REACT_APP_API_URL || 'http://localhost:8000';

const AuthPage = ({ onBack, onAuthSuccess }) => {
  const [isLogin, setIsLogin] = useState(false); // Start with Register
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    username: ''
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleToggle = (newIsLogin) => {
    if (newIsLogin !== isLogin) {
      // Animate form transition
      anime({
        targets: '.auth-form',
        translateX: [newIsLogin ? -20 : 20, 0],
        opacity: [0.7, 1],
        duration: 400,
        easing: 'easeOutCubic'
      });
      
      setIsLogin(newIsLogin);
      setError('');
      setFormData({ email: '', password: '', username: '' });
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      const endpoint = isLogin ? '/api/auth/login/' : '/api/auth/register/';
      const payload = isLogin 
        ? { email: formData.email, password: formData.password }
        : { email: formData.email, password: formData.password, username: formData.username };

      console.log('Sending request to:', endpoint, payload);
      const response = await axios.post(endpoint, payload);
      
      console.log('Auth response:', response.data);
      if (response.data.tokens) {
        localStorage.setItem('access_token', response.data.tokens.access);
        localStorage.setItem('refresh_token', response.data.tokens.refresh);
        localStorage.setItem('user', JSON.stringify(response.data.user));
        onAuthSuccess(response.data.user);
      }
    } catch (error) {
      console.error('Auth error:', error.response?.data || error.message);
      setError(error.response?.data?.error || error.response?.data?.message || 'Authentication failed');
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (e) => {
    setFormData(prev => ({
      ...prev,
      [e.target.name]: e.target.value
    }));
  };

  return (
    <div className="auth-page">
      <div className="auth-container">
        <button className="back-btn" onClick={onBack}>← Back</button>
        
        <div className="auth-form-container">
          <div className="auth-toggle">
            <button 
              className={!isLogin ? 'active' : ''}
              onClick={() => handleToggle(false)}
            >
              Register
            </button>
            <button 
              className={isLogin ? 'active' : ''}
              onClick={() => handleToggle(true)}
            >
              Login
            </button>
          </div>

          <form onSubmit={handleSubmit} className="auth-form">
            <h2>{isLogin ? 'Welcome Back' : 'Create Account'}</h2>
            
            {error && <div className="error-message">{error}</div>}
            
            <input
              type="email"
              name="email"
              placeholder="Email"
              value={formData.email}
              onChange={handleChange}
              required
            />
            
            <input
              type="password"
              name="password"
              placeholder="Password"
              value={formData.password}
              onChange={handleChange}
              required
            />
            
            {!isLogin && (
              <input
                type="text"
                name="username"
                placeholder="Username"
                value={formData.username}
                onChange={handleChange}
                required
              />
            )}
            
            <button type="submit" disabled={loading} className="submit-btn">
              {loading ? 'Processing...' : (isLogin ? 'Sign In' : 'Create Account')}
            </button>
          </form>
        </div>
      </div>
    </div>
  );
};

export default AuthPage;