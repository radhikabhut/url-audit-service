import React from 'react';
import { Search, Loader2 } from 'lucide-react';

interface AuditButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  isLoading: boolean;
}

export const AuditButton: React.FC<AuditButtonProps> = ({ isLoading, disabled, ...props }) => {
  return (
    <button
      type="submit"
      disabled={isLoading || disabled}
      className={`relative flex items-center justify-center gap-2 px-6 py-3.5 rounded-xl text-sm font-semibold tracking-wide text-white transition-all duration-200 select-none
        ${isLoading || disabled 
          ? 'bg-indigo-600/50 cursor-not-allowed border-indigo-500/20' 
          : 'bg-indigo-600 hover:bg-indigo-500 active:scale-[0.98] hover:shadow-lg hover:shadow-indigo-500/10 active:shadow-none border-indigo-500/40'
        } border`}
      {...props}
    >
      {isLoading ? (
        <>
          <Loader2 className="w-4 h-4 animate-spin text-white/80" />
          <span>Auditing...</span>
        </>
      ) : (
        <>
          <Search className="w-4 h-4 text-white/80" />
          <span>Audit URL</span>
        </>
      )}
    </button>
  );
};
