const API_BASE = 'https://qoal-backend.onrender.com/api';

interface AuthResponse {
  token: string;
  user: { id: number; email: string; name: string };
  expires_at: string;
}

export interface JobResponse {
  success: boolean;
  message: string;
  job_id: string;
  status: string;
  original_filename: string;
  file_size: number;
  source_format: string;
  target_format: string;
  created_at: string;
}

export interface JobStatus {
  job_id: string;
  status: string;
  original_filename: string;
  file_size: number;
  source_format: string;
  target_format: string;
  error?: string;
  download_url?: string;
}

const getToken = () => localStorage.getItem('token');

export const api = {
  auth: {
    register: async (email: string, password: string, name: string): Promise<AuthResponse> => {
      const res = await fetch(`${API_BASE}/auth/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password, name }),
      });
      if (!res.ok) throw new Error(await res.text());
      return res.json();
    },

    login: async (email: string, password: string): Promise<AuthResponse> => {
      const res = await fetch(`${API_BASE}/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
      });
      if (!res.ok) throw new Error('Invalid credentials');
      return res.json();
    },

    getProfile: async () => {
      const res = await fetch(`${API_BASE}/auth/profile`, {
        headers: { Authorization: `Bearer ${getToken()}` },
      });
      if (!res.ok) throw new Error('Unauthorized');
      return res.json();
    },
  },

  jobs: {
    upload: async (file: File, targetFormat: string, qualityPreset?: string): Promise<JobResponse> => {
      const formData = new FormData();
      formData.append('file', file);
      formData.append('target_format', targetFormat);
      if (qualityPreset) formData.append('quality_preset', qualityPreset);

      const res = await fetch(`${API_BASE}/upload`, {
        method: 'POST',
        headers: { Authorization: `Bearer ${getToken()}` },
        body: formData,
      });
      if (!res.ok) throw new Error(await res.text());
      return res.json();
    },

    getStatus: async (jobId: string): Promise<JobStatus> => {
      const res = await fetch(`${API_BASE}/jobs/${jobId}`, {
        headers: { Authorization: `Bearer ${getToken()}` },
      });
      if (!res.ok) throw new Error('Failed to get job status');
      return res.json();
    },

    download: async (jobId: string, filename: string) => {
      const res = await fetch(`${API_BASE}/download/${jobId}`, {
        headers: { Authorization: `Bearer ${getToken()}` },
      });
      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || 'Download failed');
      }
      
      const blob = await res.blob();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = filename;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      window.URL.revokeObjectURL(url);
    },
  },
};
