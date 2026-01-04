import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../services/api';
import { UserRound } from '../components/animate-ui/icons/user-round';
import { LockKeyhole } from '../components/animate-ui/icons/lock-keyhole';
import { Tabs, TabsList, TabsHighlight, TabsHighlightItem, TabsTrigger, TabsContents, TabsContent } from '../components/animate-ui/components/animate/tabs';

const Auth: React.FC = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [name, setName] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const navigate = useNavigate();

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSuccess('');
    setLoading(true);
    try {
      const response = await api.auth.login(email, password);
      localStorage.setItem('token', response.token);
      localStorage.setItem('user', JSON.stringify(response.user));
      setSuccess('Login successful! Redirecting...');
      setTimeout(() => {
        window.location.href = '/convert';
      }, 1500);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Login failed');
    } finally {
      setLoading(false);
    }
  };

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSuccess('');
    setLoading(true);
    try {
      const response = await api.auth.register(email, password, name);
      localStorage.setItem('token', response.token);
      localStorage.setItem('user', JSON.stringify(response.user));
      setSuccess('Registration successful! Redirecting...');
      setTimeout(() => {
        window.location.href = '/convert';
      }, 1500);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Registration failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen relative">
      <div className="fixed top-0 left-0 right-0 z-20 text-center py-3 px-4" style={{
        background: 'rgba(255, 185, 71, 0.95)',
        color: '#161B27',
        fontSize: 'clamp(0.75rem, 2vw, 0.875rem)',
        fontWeight: '500'
      }}>
        Backend currently unavailable due to expired Render Database free access (Dec 14, 2025) (<a href="https://github.com/F4tal1t/Qoal" target="_blank" rel="noopener noreferrer" style={{ color: '#161B27' }}> Project Link</a>)
      </div>
      
      <div className="min-h-screen flex items-center justify-center relative p-4" style={{ paddingTop: '4rem' }}>
        <div className="halftone-bg fixed inset-0 z-0">
          <div className="halftone-noise" />
          <div className="absolute inset-0 opacity-25 mix-blend-overlay" style={{
            backgroundImage: 'url(/BayerDithering.png)',
            backgroundRepeat: 'repeat',
            backgroundSize: '20px 20px'
          }} />
        </div>
        
        <div className="relative z-10 w-full max-w-md rounded-lg" style={{
          background: 'rgba(255, 255, 255, 0.05)',
          border: '1px solid rgba(255, 255, 255, 0.1)',
          backdropFilter: 'blur(10px)',
          padding: 'clamp(1rem, 4vw, 2rem)'
        }}>
          <div className="flex justify-center" style={{ marginBottom: 'clamp(1rem, 4vw, 2rem)' }}>
            <img src="/Qoalation.png" alt="Qoal" style={{ height: 'clamp(32px, 8vw, 48px)' }} />
          </div>
          
          <Tabs defaultValue="login">
          <TabsHighlight>
            <TabsList style={{ display: 'flex', justifyContent: 'center', gap: 'clamp(0.5rem, 2vw, 1rem)', marginBottom: 'clamp(1rem, 3vw, 2rem)' }}>
              <TabsHighlightItem value="login">
                <TabsTrigger value="login" >Login</TabsTrigger>
              </TabsHighlightItem>
              <TabsHighlightItem value="register">
                <TabsTrigger value="register" >Register</TabsTrigger>
              </TabsHighlightItem>
            </TabsList>
          </TabsHighlight>
          
          <TabsContents>
            <TabsContent value="login">
              <form onSubmit={handleLogin} style={{ display: 'flex', flexDirection: 'column', gap: 'clamp(0.75rem, 2vw, 1rem)', marginTop: 'clamp(1rem, 3vw, 1.5rem)' }}>
                <div>
                  <label className="flex items-center gap-2 mb-2" style={{ color: 'var(--color-text)', fontSize: 'clamp(0.75rem, 2vw, 0.875rem)' }}>
                    <UserRound size={18} animateOnHover />
                    Email
                  </label>
                  <input
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required
                    className="w-full rounded-md border"
                    style={{
                      background: 'rgba(255, 255, 255, 0.05)',
                      border: '1px solid rgba(255, 255, 255, 0.1)',
                      color: 'var(--color-text)',
                      padding: 'clamp(0.5rem, 2vw, 0.75rem) clamp(0.75rem, 3vw, 1rem)',
                      fontSize: 'clamp(0.875rem, 2vw, 1rem)'
                    }}
                  />
                </div>
                
                <div>
                  <label className="flex items-center gap-2 mb-2" style={{ color: 'var(--color-text)', fontSize: 'clamp(0.75rem, 2vw, 0.875rem)' }}>
                    <LockKeyhole size={18} animateOnHover />
                    Password
                  </label>
                  <input
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                    className="w-full rounded-md border"
                    style={{
                      background: 'rgba(255, 255, 255, 0.05)',
                      border: '1px solid rgba(255, 255, 255, 0.1)',
                      color: 'var(--color-text)',
                      padding: 'clamp(0.5rem, 2vw, 0.75rem) clamp(0.75rem, 3vw, 1rem)',
                      fontSize: 'clamp(0.875rem, 2vw, 1rem)'
                    }}
                  />
                </div>
                
                {error && <div className="text-red-500" style={{ fontSize: 'clamp(0.75rem, 2vw, 0.875rem)' }}>{error}</div>}
                {success && <div className="text-green-500" style={{ fontSize: 'clamp(0.75rem, 2vw, 0.875rem)' }}>{success}</div>}
                
                <button
                  type="submit"
                  disabled={loading}
                  className="w-full rounded-md font-medium transition-colors bg-[#ffb947] text-[#161B27]"
                  style={{ 
                    opacity: loading ? 0.7 : 1,
                    padding: 'clamp(0.75rem, 3vw, 1rem)',
                    fontSize: 'clamp(0.875rem, 2vw, 1rem)',
                    marginTop: 'clamp(0.5rem, 2vw, 1rem)'
                  }}
                >
                  {loading ? 'Processing...' : 'Login'}
                </button>
              </form>
            </TabsContent>
            
            <TabsContent value="register">
              <form onSubmit={handleRegister} style={{ display: 'flex', flexDirection: 'column', gap: 'clamp(0.75rem, 2vw, 1rem)', marginTop: 'clamp(1rem, 3vw, 1.5rem)' }}>
                <div>
                  <label className="flex items-center gap-2 mb-2" style={{ color: 'var(--color-text)', fontSize: 'clamp(0.75rem, 2vw, 0.875rem)' }}>
                    <UserRound size={18} animateOnHover />
                    Name
                  </label>
                  <input
                    type="text"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    required
                    className="w-full rounded-md border"
                    style={{
                      background: 'rgba(255, 255, 255, 0.05)',
                      border: '1px solid rgba(255, 255, 255, 0.1)',
                      color: 'var(--color-text)',
                      padding: 'clamp(0.5rem, 2vw, 0.75rem) clamp(0.75rem, 3vw, 1rem)',
                      fontSize: 'clamp(0.875rem, 2vw, 1rem)'
                    }}
                  />
                </div>
                
                <div>
                  <label className="flex items-center gap-2 mb-2" style={{ color: 'var(--color-text)', fontSize: 'clamp(0.75rem, 2vw, 0.875rem)' }}>
                    <UserRound size={18} animateOnHover />
                    Email
                  </label>
                  <input
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required
                    className="w-full rounded-md border"
                    style={{
                      background: 'rgba(255, 255, 255, 0.05)',
                      border: '1px solid rgba(255, 255, 255, 0.1)',
                      color: 'var(--color-text)',
                      padding: 'clamp(0.5rem, 2vw, 0.75rem) clamp(0.75rem, 3vw, 1rem)',
                      fontSize: 'clamp(0.875rem, 2vw, 1rem)'
                    }}
                  />
                </div>
                
                <div>
                  <label className="flex items-center gap-2 mb-2" style={{ color: 'var(--color-text)', fontSize: 'clamp(0.75rem, 2vw, 0.875rem)' }}>
                    <LockKeyhole size={18} animateOnHover />
                    Password
                  </label>
                  <input
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                    className="w-full rounded-md border"
                    style={{
                      background: 'rgba(255, 255, 255, 0.05)',
                      border: '1px solid rgba(255, 255, 255, 0.1)',
                      color: 'var(--color-text)',
                      padding: 'clamp(0.5rem, 2vw, 0.75rem) clamp(0.75rem, 3vw, 1rem)',
                      fontSize: 'clamp(0.875rem, 2vw, 1rem)'
                    }}
                  />
                </div>
                
                {error && <div className="text-red-500" style={{ fontSize: 'clamp(0.75rem, 2vw, 0.875rem)' }}>{error}</div>}
                {success && <div className="text-green-500" style={{ fontSize: 'clamp(0.75rem, 2vw, 0.875rem)' }}>{success}</div>}
                
                <button
                  type="submit"
                  disabled={loading}
                  className="w-full rounded-md font-medium transition-colors bg-[#ffb947] text-[#161B27]"
                  style={{ 
                    opacity: loading ? 0.7 : 1,
                    padding: 'clamp(0.75rem, 3vw, 1rem)',
                    fontSize: 'clamp(0.875rem, 2vw, 1rem)',
                    marginTop: 'clamp(0.5rem, 2vw, 1rem)'
                  }}
                >
                  {loading ? 'Processing...' : 'Register'}
                </button>
              </form>
            </TabsContent>
          </TabsContents>
        </Tabs>
      </div>
    </div>
  );
};

export default Auth;
