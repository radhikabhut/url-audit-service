import { useMutation } from '@tanstack/react-query';
import axios from 'axios';
import { auditURL } from '../api';
import { AuditResult, APIError } from '../types';

export interface UseAuditResult {
  audit: (url: string) => void;
  isLoading: boolean;
  isError: boolean;
  isSuccess: boolean;
  result: AuditResult | null;
  error: APIError | Error | null;
  reset: () => void;
}

export const useAudit = (): UseAuditResult => {
  const mutation = useMutation<AuditResult, any, string>({
    mutationFn: auditURL,
  });

  const getCleanError = (err: any): APIError | Error | null => {
    if (!err) return null;
    if (axios.isAxiosError(err) && err.response?.data) {
      return err.response.data as APIError;
    }
    return err as Error;
  };

  return {
    audit: mutation.mutate,
    isLoading: mutation.isPending,
    isError: mutation.isError,
    isSuccess: mutation.isSuccess,
    result: mutation.data || null,
    error: getCleanError(mutation.error),
    reset: mutation.reset,
  };
};
