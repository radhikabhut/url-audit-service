import React, { useState } from 'react';
import { Globe, AlertTriangle } from 'lucide-react';
import { AuditButton } from './AuditButton';

interface AuditFormProps {
  onSubmit: (url: string) => void;
  isLoading: boolean;
}

export const AuditForm: React.FC<AuditFormProps> = ({ onSubmit, isLoading }) => {
  const [url, setUrl] = useState('');
  const [validationError, setValidationError] = useState('');

  const validate = (val: string): boolean => {
    if (!val.trim()) {
      setValidationError('URL cannot be empty');
      return false;
    }

    try {
      const parsed = new URL(val);
      if (parsed.protocol !== 'http:' && parsed.protocol !== 'https:') {
        setValidationError('Only HTTP and HTTPS protocols are allowed');
        return false;
      }
    } catch {
      setValidationError('Please enter a valid, well-formed URL (e.g., https://example.com)');
      return false;
    }

    setValidationError('');
    return true;
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = e.target.value;
    setUrl(val);
    if (validationError) {
      try {
        const parsed = new URL(val);
        if (parsed.protocol === 'http:' || parsed.protocol === 'https:') {
          setValidationError('');
        }
      } catch {
        // keep error
      }
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (validate(url)) {
      onSubmit(url.trim());
    }
  };

  return (
    <form onSubmit={handleSubmit} className="w-full max-w-xl mx-auto">
      <div className="saas-card rounded-2xl p-6 shadow-sm shadow-slate-200/80 border border-slate-200/60 bg-white relative overflow-hidden group">
        <div className="absolute top-0 left-0 right-0 h-[2px] bg-indigo-600 opacity-85"></div>
        <div className="flex flex-col gap-4">
          <label htmlFor="url" className="text-xs font-semibold text-slate-400 tracking-wider uppercase">
            Target Website URL
          </label>
          <div className="relative flex items-center">
            <Globe className="absolute left-4 w-5 h-5 text-slate-400 group-focus-within:text-indigo-600 transition-colors" />
            <input
              type="text"
              id="url"
              name="url"
              value={url}
              onChange={handleInputChange}
              placeholder="https://example.com"
              disabled={isLoading}
              className={`w-full bg-white border ${
                validationError ? 'border-red-400 focus:border-red-500' : 'border-slate-200 focus:border-indigo-600'
              } rounded-xl py-3.5 pl-12 pr-4 text-sm text-slate-800 placeholder-slate-400 focus:outline-none focus:ring-4 ${
                validationError ? 'focus:ring-red-500/5' : 'focus:ring-indigo-600/5'
              } transition-all duration-200`}
            />
          </div>

          {validationError && (
            <div className="flex items-center gap-2 text-xs text-red-600 mt-1 bg-red-50 border border-red-200 py-2 px-3 rounded-lg animate-fadeIn">
              <AlertTriangle className="w-3.5 h-3.5 shrink-0" />
              <span>{validationError}</span>
            </div>
          )}

          <div className="flex justify-end mt-2">
            <AuditButton isLoading={isLoading} disabled={!url} />
          </div>
        </div>
      </div>
    </form>
  );
};
