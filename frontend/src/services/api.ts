const API_BASE = 'http://localhost:8000/api';

interface AuthResponse {
  token: string;
  user: { id: number; email: string; name: string };
  expires_at: string;
}

interface JobResponse {
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

interface JobStatus {
  job_id: string;
  status: string;
  original_filename: string;
  file_size: number;
  source_format: string;
  target_format: string;
  error: string;
  created_at: string;
  updated_at: string;
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
      if (!res.ok) throw new Error('Job not found');
      return res.json();
    },

    list: async (page = 1, limit = 10) => {
      const res = await fetch(`${API_BASE}/jobs?page=${page}&limit=${limit}`, {
        headers: { Authorization: `Bearer ${getToken()}` },
      });
      if (!res.ok) throw new Error('Failed to fetch jobs');
      return res.json();
    },

    download: async (jobId: string) => {
      const res = await fetch(`${API_BASE}/download/${jobId}`, {
        headers: { Authorization: `Bearer ${getToken()}` },
      });
      if (!res.ok) throw new Error('Download failed');
      const blob = await res.blob();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `converted_${jobId}`;
      a.click();
      window.URL.revokeObjectURL(url);
    },
  },
};
