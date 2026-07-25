import React from 'react';
import { Activity } from 'lucide-react';

export const Header: React.FC = () => {
  return (
    <header className="flex flex-col items-center justify-center text-center py-10">
      <div className="flex items-center gap-3 mb-2.5">
        <div className="flex items-center justify-center w-11 h-11 rounded-xl bg-indigo-600 shadow-md shadow-indigo-600/10">
          <Activity className="w-5.5 h-5.5 text-white animate-pulse" />
        </div>
        <h1 className="text-3xl font-extrabold tracking-tight text-slate-900">
          Page Pulse
        </h1>
      </div>
      <p className="text-xs font-bold tracking-wider uppercase text-indigo-600 mb-1">
        Production-Grade URL Audit Service
      </p>
      <p className="text-sm text-slate-500 max-w-md">
        Analyze target site availability, load latency, SSL certificates, headers, page titles, and caching status in real-time.
      </p>
    </header>
  );
};
