const API_BASE = 'http://localhost:8000/api';

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

const getToken = () => localStorage.getItem('token');

export const api = {
  auth: {
    login: async (email: string, password: string): Promise<AuthResponse> => {
      const res = await fetch(`${API_BASE}/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
      });
      if (!res.ok) throw new Error('Invalid credentials');
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
