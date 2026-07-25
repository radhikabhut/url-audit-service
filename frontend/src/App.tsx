import React from 'react';
import { Header } from './components/Header';
import { AuditForm } from './components/AuditForm';
import { Loading } from './components/Loading';
import { ErrorAlert } from './components/ErrorAlert';
import { ResultCard } from './components/ResultCard';
import { Footer } from './components/Footer';
import { useAudit } from './hooks/useAudit';

const App: React.FC = () => {
  const { audit, isLoading, isError, isSuccess, result, error } = useAudit();

  const handleAuditSubmit = (url: string) => {
    audit(url);
  };

  return (
    <div className="flex flex-col min-h-screen">
      <main className="flex-grow container mx-auto px-4 py-8 flex flex-col items-center">
        <Header />
        
        <div className="w-full max-w-2xl mt-4">
          <AuditForm onSubmit={handleAuditSubmit} isLoading={isLoading} />

          {isLoading && (
            <div className="mt-8 animate-fadeIn">
              <Loading />
            </div>
          )}

          {isError && (
            <div className="mt-8">
              <ErrorAlert error={error} />
            </div>
          )}

          {isSuccess && result && !isLoading && (
            <ResultCard result={result} />
          )}
        </div>
      </main>
      
      <Footer />
    </div>
  );
};

export default App;
