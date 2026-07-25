import axios from 'axios';
import { AuditResult } from './types';

const apiBaseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

export const apiClient = axios.create({
  baseURL: apiBaseUrl,
  headers: {
    'Content-Type': 'application/json',
  },
});

export const auditURL = async (url: string): Promise<AuditResult> => {
  const response = await apiClient.post<AuditResult>('/api/v1/audit', { url });
  return response.data;
};
