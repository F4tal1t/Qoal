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
  const navigate = useNavigate();

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const response = await api.auth.login(email, password);
      localStorage.setItem('token', response.token);
      localStorage.setItem('user', JSON.stringify(response.user));
      navigate('/convert/image');
    } catch (err: any) {
      setError(err.message || 'Login failed');
    } finally {
      setLoading(false);
    }
  };

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const response = await api.auth.register(email, password, name);
      localStorage.setItem('token', response.token);
      localStorage.setItem('user', JSON.stringify(response.user));
      navigate('/convert/image');
    } catch (err: any) {
      setError(err.message || 'Registration failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center relative">
      <div className="halftone-bg fixed inset-0 z-0">
        <div className="halftone-noise" />
        <div className="absolute inset-0 opacity-25 mix-blend-overlay" style={{
          backgroundImage: 'url(/BayerDithering.png)',
          backgroundRepeat: 'repeat',
          backgroundSize: '20px 20px'
        }} />
      </div>
      
      <div className="relative z-10 w-full max-w-md mx-4 p-8 rounded-lg" style={{
        background: 'rgba(255, 255, 255, 0.05)',
        border: '1px solid rgba(255, 255, 255, 0.1)',
        backdropFilter: 'blur(10px)'
      }}>
        <div className="flex justify-center mb-8">
          <img src="/Qoalation.png" alt="Qoal" style={{ height: '48px' }} />
        </div>
        
        <Tabs defaultValue="login">
          <TabsHighlight>
            <TabsList style={{ display: 'flex', justifyContent: 'center', gap: '1rem', marginBottom: '2rem' }}>
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
              <form onSubmit={handleLogin} className="space-y-4 mt-6">
                <div>
                  <label className="flex items-center gap-2 text-sm mb-2" style={{ color: 'var(--color-text)' }}>
                    <UserRound size={18} animateOnHover />
                    Email
                  </label>
                  <input
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required
                    className="w-full px-4 py-2 rounded-md border"
                    style={{
                      background: 'rgba(255, 255, 255, 0.05)',
                      border: '1px solid rgba(255, 255, 255, 0.1)',
                      color: 'var(--color-text)'
                    }}
                  />
                </div>
                
                <div>
                  <label className="flex items-center gap-2 text-sm mb-2" style={{ color: 'var(--color-text)' }}>
                    <LockKeyhole size={18} animateOnHover />
                    Password
                  </label>
                  <input
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                    className="w-full px-4 py-2 rounded-md border"
                    style={{
                      background: 'rgba(255, 255, 255, 0.05)',
                      border: '1px solid rgba(255, 255, 255, 0.1)',
                      color: 'var(--color-text)'
                    }}
                  />
                </div>
                
                {error && <div className="text-red-500 text-sm">{error}</div>}
                
                <button
                  type="submit"
                  disabled={loading}
                  className="w-full py-3 rounded-md font-medium transition-colors bg-[#ffb947] text-[#161B27] mt-4"
                  style={{ opacity: loading ? 0.7 : 1 }}
                >
                  {loading ? 'Processing...' : 'Login'}
                </button>
              </form>
            </TabsContent>
            
            <TabsContent value="register">
              <form onSubmit={handleRegister} className="space-y-4 mt-6">
                <div>
                  <label className="flex items-center gap-2 text-sm mb-2" style={{ color: 'var(--color-text)' }}>
                    <UserRound size={18} animateOnHover />
                    Name
                  </label>
                  <input
                    type="text"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    required
                    className="w-full px-4 py-2 rounded-md border"
                    style={{
                      background: 'rgba(255, 255, 255, 0.05)',
                      border: '1px solid rgba(255, 255, 255, 0.1)',
                      color: 'var(--color-text)'
                    }}
                  />
                </div>
                
                <div>
                  <label className="flex items-center gap-2 text-sm mb-2" style={{ color: 'var(--color-text)' }}>
                    <UserRound size={18} animateOnHover />
                    Email
                  </label>
                  <input
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required
                    className="w-full px-4 py-2 rounded-md border"
                    style={{
                      background: 'rgba(255, 255, 255, 0.05)',
                      border: '1px solid rgba(255, 255, 255, 0.1)',
                      color: 'var(--color-text)'
                    }}
                  />
                </div>
                
                <div>
                  <label className="flex items-center gap-2 text-sm mb-2" style={{ color: 'var(--color-text)' }}>
                    <LockKeyhole size={18} animateOnHover />
                    Password
                  </label>
                  <input
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                    className="w-full px-4 py-2 rounded-md border"
                    style={{
                      background: 'rgba(255, 255, 255, 0.05)',
                      border: '1px solid rgba(255, 255, 255, 0.1)',
                      color: 'var(--color-text)'
                    }}
                  />
                </div>
                
                {error && <div className="text-red-500 text-sm">{error}</div>}
                
                <button
                  type="submit"
                  disabled={loading}
                  className="w-full py-3 rounded-md font-medium transition-colors bg-[#ffb947] text-[#161B27] mt-4"
                  style={{ opacity: loading ? 0.7 : 1 }}
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
