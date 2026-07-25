import React from 'react';
import { 
  CheckCircle2, 
  XCircle, 
  Clock, 
  FileText, 
  Database, 
  FileCheck, 
  Calendar, 
  Layers,
  Globe2 
} from 'lucide-react';
import { AuditResult } from '../types';

interface ResultCardProps {
  result: AuditResult | null;
}

export const ResultCard: React.FC<ResultCardProps> = ({ result }) => {
  if (!result) return null;

  const getStatusColor = (code: number) => {
    if (code >= 200 && code < 300) return 'text-emerald-700 bg-emerald-50 border-emerald-200';
    if (code >= 300 && code < 400) return 'text-indigo-700 bg-indigo-50 border-indigo-200';
    return 'text-rose-700 bg-rose-50 border-rose-200';
  };

  const getResponseTimeColor = (ms: number) => {
    if (ms < 200) return 'text-emerald-600';
    if (ms < 500) return 'text-amber-600';
    return 'text-rose-600';
  };

  const formatBytes = (bytes: number) => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  return (
    <div className="w-full max-w-2xl mx-auto my-8 space-y-5 animate-fadeIn">
      {/* Overview Card */}
      <div className="saas-card rounded-2xl p-6 shadow-sm shadow-slate-200/80 border border-slate-200/60 bg-white relative overflow-hidden">
        <div className="absolute top-0 left-0 right-0 h-[2px] bg-gradient-to-r from-emerald-500 via-teal-500 to-indigo-500 opacity-80"></div>
        
        <div className="flex flex-col gap-4">
          <div className="flex flex-wrap items-center justify-between gap-3 border-b border-slate-100 pb-4">
            <div className="flex items-center gap-2">
              <Globe2 className="w-4.5 h-4.5 text-indigo-600" />
              <span className="text-sm font-semibold text-slate-800 break-all select-all">
                {result.url}
              </span>
            </div>
            
            <div className="flex items-center gap-2">
              {result.reachable ? (
                <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-emerald-50 text-emerald-700 border border-emerald-200">
                  <CheckCircle2 className="w-3.5 h-3.5" />
                  Reachable
                </span>
              ) : (
                <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-rose-50 text-rose-700 border border-rose-200">
                  <XCircle className="w-3.5 h-3.5" />
                  Unreachable
                </span>
              )}

              {result.cached && (
                <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-indigo-50 text-indigo-700 border border-indigo-200">
                  <Layers className="w-3.5 h-3.5" />
                  Cached
                </span>
              )}
            </div>
          </div>

          <div className="py-1">
            <h2 className="text-xs font-semibold text-slate-400 tracking-wider uppercase mb-1">
              Page Title
            </h2>
            <p className="text-xl font-bold text-slate-850 leading-tight">
              {result.title || <span className="text-slate-400 italic font-normal text-sm">No page title extracted</span>}
            </p>
          </div>
        </div>
      </div>

      {/* Grid of detailed specs */}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        {/* Metric 1: Status Code */}
        <div className="saas-card rounded-xl p-5 flex items-center justify-between bg-white border border-slate-100 shadow-sm shadow-slate-200/50">
          <div className="space-y-1">
            <span className="text-xs font-semibold text-slate-400 tracking-wider uppercase block">HTTP Status</span>
            <span className="text-2xl font-extrabold text-slate-800">{result.statusCode || 'N/A'}</span>
          </div>
          <div className={`px-3 py-1 rounded-lg border text-xs font-bold ${getStatusColor(result.statusCode)}`}>
            {result.statusCode === 200 ? 'OK' : result.statusCode ? `HTTP ${result.statusCode}` : 'FAILED'}
          </div>
        </div>

        {/* Metric 2: Latency */}
        <div className="saas-card rounded-xl p-5 flex items-center justify-between bg-white border border-slate-100 shadow-sm shadow-slate-200/50">
          <div className="space-y-1">
            <span className="text-xs font-semibold text-slate-400 tracking-wider uppercase block">Response Time</span>
            <span className={`text-2xl font-extrabold ${getResponseTimeColor(result.responseTimeMs)}`}>
              {result.responseTimeMs} <span className="text-sm font-semibold text-slate-400">ms</span>
            </span>
          </div>
          <div className="p-2.5 rounded-lg bg-slate-50 text-slate-500 border border-slate-100">
            <Clock className="w-5 h-5" />
          </div>
        </div>

        {/* Metric 3: Content Type */}
        <div className="saas-card rounded-xl p-5 flex items-center justify-between bg-white border border-slate-100 shadow-sm shadow-slate-200/50">
          <div className="space-y-1">
            <span className="text-xs font-semibold text-slate-400 tracking-wider uppercase block">Content Type</span>
            <span className="text-sm font-bold text-slate-700 truncate max-w-[200px] block" title={result.contentType}>
              {result.contentType || 'N/A'}
            </span>
          </div>
          <div className="p-2.5 rounded-lg bg-slate-50 text-slate-500 border border-slate-100">
            <FileText className="w-5 h-5" />
          </div>
        </div>

        {/* Metric 4: Content Length */}
        <div className="saas-card rounded-xl p-5 flex items-center justify-between bg-white border border-slate-100 shadow-sm shadow-slate-200/50">
          <div className="space-y-1">
            <span className="text-xs font-semibold text-slate-400 tracking-wider uppercase block">Content Size</span>
            <span className="text-2xl font-extrabold text-slate-800">
              {result.contentLength !== undefined && result.contentLength >= 0 ? formatBytes(result.contentLength) : 'N/A'}
            </span>
          </div>
          <div className="p-2.5 rounded-lg bg-slate-50 text-slate-500 border border-slate-100">
            <Database className="w-5 h-5" />
          </div>
        </div>

        {/* Metric 5: Cached */}
        <div className="saas-card rounded-xl p-5 flex items-center justify-between bg-white border border-slate-100 shadow-sm shadow-slate-200/50">
          <div className="space-y-1">
            <span className="text-xs font-semibold text-slate-400 tracking-wider uppercase block">Served From Cache</span>
            <span className="text-2xl font-extrabold text-slate-850">{result.cached ? 'Yes' : 'No'}</span>
          </div>
          <div className={`p-2.5 rounded-lg border ${result.cached ? 'bg-indigo-50 text-indigo-600 border-indigo-100' : 'bg-slate-50 text-slate-400 border-slate-100'}`}>
            <FileCheck className="w-5 h-5" />
          </div>
        </div>

        {/* Metric 6: Checked At */}
        <div className="saas-card rounded-xl p-5 flex items-center justify-between bg-white border border-slate-100 shadow-sm shadow-slate-200/50">
          <div className="space-y-1">
            <span className="text-xs font-semibold text-slate-400 tracking-wider uppercase block">Audited At</span>
            <span className="text-xs font-bold text-slate-600 block" title={result.checkedAt}>
              {result.checkedAt ? new Date(result.checkedAt).toLocaleString() : 'N/A'}
            </span>
          </div>
          <div className="p-2.5 rounded-lg bg-slate-50 text-slate-500 border border-slate-100">
            <Calendar className="w-5 h-5" />
          </div>
        </div>
      </div>
    </div>
  );
};
