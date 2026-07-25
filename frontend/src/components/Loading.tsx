import React from 'react';

export const Loading: React.FC = () => {
  return (
    <div className="flex flex-col items-center justify-center py-12 px-4">
      <div className="relative w-14 h-14 mb-4">
        <div className="absolute inset-0 rounded-full border-4 border-indigo-100"></div>
        <div className="absolute inset-0 rounded-full border-4 border-t-indigo-600 border-r-indigo-400 border-b-indigo-200 border-l-transparent animate-spin"></div>
      </div>
      <p className="text-sm font-semibold text-slate-700 animate-pulse">
        Auditing URL...
      </p>
      <p className="text-xs text-slate-400 mt-1">
        Preventing SSRF, resolving DNS, and performing GET checks
      </p>
    </div>
  );
};
