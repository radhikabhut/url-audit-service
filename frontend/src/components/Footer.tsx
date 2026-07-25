import React from 'react';

export const Footer: React.FC = () => {
  return (
    <footer className="text-center py-8 mt-auto border-t border-slate-150 text-xs text-slate-400">
      <p>© {new Date().getFullYear()} Page Pulse URL Audit Service. All rights reserved.</p>
      <p className="mt-1 text-[10px] text-slate-350 font-mono">Build v1.0.0 (Go + React + Vite + TS)</p>
      <p className="mt-2">
        <a 
          href="https://digitalheroesco.com" 
          target="_blank" 
          rel="noopener noreferrer" 
          className="text-indigo-600 hover:underline hover:text-indigo-500 font-medium transition-colors"
        >
          Built for Digital Heroes Training Task
        </a>
      </p>
    </footer>
  );
};
