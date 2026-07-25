import React from 'react';
import { AlertCircle, ShieldAlert } from 'lucide-react';
import { APIError } from '../types';

interface ErrorAlertProps {
  error: APIError | null | unknown;
}

export const ErrorAlert: React.FC<ErrorAlertProps> = ({ error }) => {
  if (!error) return null;

  const apiError = error as APIError;
  const isStructured = apiError?.error?.code !== undefined;

  const errorCode = isStructured ? apiError.error.code : 'UNKNOWN_ERROR';
  const errorMessage = isStructured ? apiError.error.message : 'An unexpected error occurred. Please try again.';
  const requestId = apiError?.requestId;

  return (
    <div className="bg-red-50 border border-red-200 rounded-xl p-5 max-w-xl mx-auto my-6 shadow-sm">
      <div className="flex items-start gap-4">
        <div className="p-2 rounded-lg bg-red-100 text-red-600">
          {errorCode === 'INVALID_URL' ? (
            <ShieldAlert className="w-5.5 h-5.5" />
          ) : (
            <AlertCircle className="w-5.5 h-5.5" />
          )}
        </div>
        <div className="flex-1">
          <h3 className="text-sm font-bold text-red-800 tracking-wide uppercase">
            Audit Failed: {errorCode.replace(/_/g, ' ')}
          </h3>
          <p className="text-sm text-red-700 mt-1 leading-relaxed">
            {errorMessage}
          </p>
          {requestId && (
            <div className="flex items-center gap-1.5 mt-3 text-[10px] font-mono text-slate-500 bg-red-100/50 py-1 px-2.5 rounded border border-red-200/40 w-fit">
              <span>Request ID:</span>
              <span className="select-all text-slate-600">{requestId}</span>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
