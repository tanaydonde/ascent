import { useState, useEffect } from 'react';

interface Problem {
  id: string;
  name: string;
  rating: number;
  tags: string[];
  link?: string;
}

interface SubmitModalProps {
  isOpen: boolean;
  onClose: () => void;
  problem: Problem | null;
  onSubmit: (timeMinutes: number) => void;
  isSubmitting: boolean;
  isSuccess: boolean;
  error: string | null;
}

const SubmitModal = ({ isOpen, onClose, problem, onSubmit, isSubmitting, isSuccess, error }: SubmitModalProps) => {
  const [time, setTime] = useState<string>('');

  if (!isOpen || !problem) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/60 backdrop-blur-sm">
      <div className="bg-slate-900 border border-slate-700 rounded-xl shadow-2xl w-full max-w-md p-6 relative animate-in fade-in zoom-in duration-200">
        
        {isSuccess ? (
          <div className="flex flex-col items-center justify-center py-6 text-center animate-in fade-in slide-in-from-bottom-2 duration-300">
            <div className="w-16 h-16 bg-emerald-500/20 rounded-full flex items-center justify-center mb-4 ring-1 ring-emerald-500/50 shadow-[0_0_20px_rgba(16,185,129,0.2)]">
              <svg className="w-8 h-8 text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={3}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <h3 className="text-xl font-bold text-gray-100 mb-1">Great Job!</h3>
            <p className="text-slate-400 text-sm">
              Challenge problem complete. Mastery updated.
            </p>
          </div>
        ) : (
          <>
            <h3 className="text-xl font-bold text-gray-100 mb-1">Verify Solution</h3>
            <p className="text-slate-400 text-sm mb-6">
              Did you solve <span className="text-sky-400 font-mono">{problem.id}</span> on Codeforces?
            </p>

            <div className="space-y-3 mb-6">
              <label className="block text-xs font-semibold text-slate-500 uppercase tracking-wider">
                Time Taken (Optional)
              </label>
              <div className="relative">
                <input 
                  type="number" 
                  value={time}
                  onChange={(e) => setTime(e.target.value)}
                  placeholder="0"
                  className="w-full bg-slate-800 border border-slate-700 rounded-lg px-4 py-3 text-gray-200 focus:ring-2 focus:ring-sky-500 focus:border-transparent outline-none transition-all placeholder:text-slate-600 [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
                />
                <span className="absolute right-4 top-3.5 text-slate-500 text-sm">min</span>
              </div>
            </div>

            {error && (
              <div className="mb-4 p-3 bg-rose-500/10 border border-rose-500/20 rounded-lg text-rose-400 text-xs">
                {error}
              </div>
            )}

            <div className="flex gap-3 mt-2">
              <button 
                onClick={onClose}
                className="flex-1 px-4 py-2.5 rounded-lg bg-slate-800 text-slate-300 font-medium hover:bg-slate-700 transition-colors text-sm"
                disabled={isSubmitting}
              >
                Cancel
              </button>
              <button 
                onClick={() => onSubmit(time === '' ? 0 : parseInt(time))}
                disabled={isSubmitting}
                className="flex-1 px-4 py-2.5 rounded-lg bg-sky-500 text-white font-bold hover:bg-sky-400 transition-colors text-sm flex justify-center items-center gap-2 shadow-lg shadow-sky-500/20"
              >
                {isSubmitting ? (
                  <span className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"/>
                ) : (
                  <>Verify & Submit</>
                )}
              </button>
            </div>
          </>
        )}
      </div>
    </div>
  );
};

const getLink = (id: string) => {
  const match = id.match(/^(\d+)(.+)$/);
  if (match) {
    return `https://codeforces.com/problemset/problem/${match[1]}/${match[2]}`;
  }
  return `https://codeforces.com/problemset/problem/${id}`;
};

const Challenge = () => {
  const [problem, setProblem] = useState<Problem | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const [modalOpen, setModalOpen] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [isSuccess, setIsSuccess] = useState(false);

  const handle = sessionStorage.getItem('cf_handle');

  useEffect(() => {
    const fetchDaily = async () => {
      if (!handle) return;
      
      try {
        setLoading(true);
        const res = await fetch(`https://ascent-backend-842l.onrender.com/api/daily?handle=${handle}`);
        if (!res.ok) throw new Error('Failed to fetch challenge problem');
        
        const data = await res.json();
        setProblem({
            ...data,
            link: getLink(data.id)
        });
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchDaily();
  }, [handle]);

  const handleSubmitVerification = async (timeMinutes: number) => {
    if (!problem || !handle) return;
    setIsSubmitting(true);
    setSubmitError(null);

    try {
      const res = await fetch(`https://ascent-backend-842l.onrender.com/api/submit/${handle}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          problem_id: problem.id,
          time_spent_minutes: timeMinutes
        })
      });

      if (!res.ok) {
        const errMsg = await res.text();
        throw new Error(errMsg || "Verification failed");
      }
      
      setIsSubmitting(false);
      setIsSuccess(true);

      setTimeout(() => {
        setModalOpen(false);
        setIsSuccess(false);
      }, 2000);
      
    } catch (err: any) {
      let msg = err.message;
      if (msg.includes("not solved")) msg = "Codeforces says you haven't solved this yet! Please wait a minute if you just submitted.";
      else if (msg.includes("already solved")) msg = "You've already tracked this problem! Good job.";
      
      setSubmitError(msg);
      setIsSubmitting(false);
    }
  };

  if (loading) {
    return (
      <div className="h-full flex items-center justify-center text-slate-500 font-mono animate-pulse">
        LOADING...
      </div>
    );
  }

  if (error) {
    return (
      <div className="h-full flex flex-col items-center justify-center gap-4">
        <div className="text-red-400 font-mono">SYSTEM FAILURE: {error}</div>
        <button 
          onClick={() => window.location.reload()}
          className="px-4 py-2 bg-red-500/10 border border-red-500/50 text-red-400 rounded hover:bg-red-500/20 transition-all"
        >
          RETRY SEQUENCE
        </button>
      </div>
    );
  }

  return (
    <>
      <div className="h-full w-full flex items-center justify-center p-6">
        <div className="max-w-2xl w-full bg-slate-900/50 backdrop-blur-md border border-slate-700 rounded-2xl p-8 relative overflow-hidden group">
          
          <div className="absolute top-0 right-0 w-64 h-64 bg-sky-500/5 rounded-full blur-3xl -translate-y-1/2 translate-x-1/2 group-hover:bg-sky-500/10 transition-all duration-700" />

          <div className="relative z-10">
            <div className="flex justify-between items-start mb-6">
              <div>
                <h2 className="text-2xl font-bold text-white mb-1">Random Challenge Problem</h2>
                <p className="text-slate-400 text-sm">
                  System selected challenge based on your recent performance.
                </p>
              </div>
            </div>

            {problem && (
              <div className="bg-slate-800/50 rounded-xl p-6 border border-slate-700/50">
                <div className="flex justify-between items-center mb-4">
                  <span className="text-3xl font-bold text-white tracking-tight">
                    {problem.id}
                  </span>
                  <span className={`text-lg font-mono font-bold ${
                    problem.rating >= 2000 ? 'text-red-400' :
                    problem.rating >= 1600 ? 'text-blue-400' :
                    problem.rating >= 1200 ? 'text-cyan-400' : 'text-green-400'
                  }`}>
                    {problem.rating}
                  </span>
                </div>
                
                <h3 className="text-xl text-slate-200 mb-6 font-medium">
                  {problem.name}
                </h3>

                <div className="flex gap-4">
                  <a 
                    href={problem.link} 
                    target="_blank" 
                    rel="noreferrer"
                    className="flex-1 text-center py-4 bg-sky-600 hover:bg-sky-500 text-white font-bold rounded-lg transition-all shadow-lg shadow-sky-900/20"
                  >
                    SOLVE ON CODEFORCES
                  </a>
                  
                  <button 
                    onClick={() => {
                        setSubmitError(null);
                        setModalOpen(true);
                    }}
                    className="flex-1 py-4 bg-emerald-600/10 hover:bg-emerald-600/20 border border-emerald-500/50 text-emerald-400 font-bold rounded-lg transition-all"
                  >
                    VERIFY SOLUTION
                  </button>
                </div>

              </div>
            )}
          </div>
        </div>
      </div>

      <SubmitModal 
        isOpen={modalOpen}
        onClose={() => setModalOpen(false)}
        problem={problem}
        onSubmit={handleSubmitVerification}
        isSubmitting={isSubmitting}
        isSuccess={isSuccess}
        error={submitError}
      />
    </>
  );
};

export default Challenge;