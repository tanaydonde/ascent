import { useState } from 'react';

const LogProblems = () => {
  const [problemId, setProblemId] = useState('');
  const [time, setTime] = useState('');
  
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isSyncing, setIsSyncing] = useState(false);
  
  const [message, setMessage] = useState<{ text: string, type: 'success' | 'error' } | null>(null);

  const handle = sessionStorage.getItem('cf_handle');

  const handleManualSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!handle || !problemId) return;

    setIsSubmitting(true);
    setMessage(null);

    try {
      const res = await fetch(`https://ascent-backend-842l.onrender.com/api/submit/${handle}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          problem_id: problemId.trim().toUpperCase(),
          time_spent_minutes: time ? parseInt(time) : 0
        })
      });

      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || "Failed to log problem");
      }

      setMessage({ text: `Successfully logged ${problemId}!`, type: 'success' });
      setProblemId('');
      setTime('');
    } catch (err: any) {
        let msg = err.message;
        if (msg.includes("not solved")) msg = "Codeforces says this isn't solved yet. Did you submit it?";
        setMessage({ text: msg, type: 'error' });
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleSync = async () => {
    if (!handle) return;

    setIsSyncing(true);
    setMessage(null);

    try {
      const res = await fetch(`https://ascent-backend-842l.onrender.com/api/sync/${handle}`, {
        method: 'POST'
      });

      if (!res.ok) throw new Error("Sync failed");

      setMessage({ 
        text: `Sync complete!`, 
        type: 'success' 
      });
    } catch (err: any) {
      setMessage({ text: "Failed to sync with Codeforces.", type: 'error' });
    } finally {
      setIsSyncing(false);
    }
  };

  return (
    <div className="h-full w-full flex items-center justify-center p-6">
      <div className="w-full max-w-4xl grid grid-cols-1 md:grid-cols-2 gap-8">
        
        <div className="bg-slate-900/50 backdrop-blur-md border border-slate-700 rounded-2xl p-8 flex flex-col">
          <div className="mb-6">
            <h2 className="text-2xl font-bold text-white flex items-center gap-3">
              Manual Log
              <span className="text-[10px] bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 px-2 py-0.5 rounded uppercase tracking-wide">
                Recommended
              </span>
            </h2>
            <p className="text-slate-400 text-sm mt-2">
              Manually entering the time taken helps calculate your rating mastery much more accurately than auto-syncing.
            </p>
          </div>

          <form onSubmit={handleManualSubmit} className="space-y-4 flex-1">
            <div>
              <label className="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1">
                Problem ID
              </label>
              <input 
                type="text" 
                value={problemId}
                onChange={(e) => setProblemId(e.target.value)}
                placeholder="e.g. 158B"
                className="w-full bg-slate-800 border border-slate-700 rounded-lg px-4 py-3 text-gray-200 focus:ring-2 focus:ring-sky-500 outline-none uppercase placeholder:normal-case"
              />
            </div>

            <div>
              <label className="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1">
                Time Taken (Minutes)
              </label>
              <input 
                type="number" 
                value={time}
                onChange={(e) => setTime(e.target.value)}
                placeholder="Optional but recommended"
                className="w-full bg-slate-800 border border-slate-700 rounded-lg px-4 py-3 text-gray-200 focus:ring-2 focus:ring-sky-500 outline-none"
              />
            </div>

            <button 
              type="submit"
              disabled={isSubmitting || !problemId}
              className="w-full py-3 bg-sky-600 hover:bg-sky-500 disabled:bg-slate-800 disabled:text-slate-600 text-white font-bold rounded-lg transition-all shadow-lg shadow-sky-900/20 mt-4"
            >
              {isSubmitting ? "Verifying..." : "Log Problem"}
            </button>
          </form>
        </div>

        <div className="bg-slate-900/50 backdrop-blur-md border border-slate-700 rounded-2xl p-8 flex flex-col justify-between">
          <div>
            <h2 className="text-2xl font-bold text-white mb-2">Auto-Sync</h2>
            <p className="text-slate-400 text-sm">
              One-click sync with Codeforces. Useful if you've solved many problems recently and just want to update your list.
            </p>
            
            <div className="mt-6 p-4 bg-amber-500/5 border border-amber-500/10 rounded-lg">
              <p className="text-amber-500/80 text-xs leading-relaxed">
                <strong>Note:</strong> Codeforces doesn't report how long you took to solve a problem. We apply a default estimate for synced problems, which makes your mastery rating less precise than manual logging.
              </p>
            </div>
          </div>

          <div className="mt-8">
            <button 
              onClick={handleSync}
              disabled={isSyncing}
              className="w-full py-3 bg-slate-800 hover:bg-slate-700 border border-slate-600 text-slate-300 font-bold rounded-lg transition-all flex items-center justify-center gap-2"
            >
              {isSyncing ? (
                <>
                  <span className="w-4 h-4 border-2 border-slate-400 border-t-transparent rounded-full animate-spin"></span>
                  Syncing...
                </>
              ) : (
                "Sync with Codeforces"
              )}
            </button>
          </div>
        </div>

      </div>

      {message && (
        <div className={`fixed bottom-8 left-1/2 -translate-x-1/2 px-6 py-3 rounded-xl border shadow-2xl backdrop-blur-xl animate-in fade-in slide-in-from-bottom-4 ${
          message.type === 'success' 
            ? 'bg-emerald-900/80 border-emerald-500/50 text-emerald-100' 
            : 'bg-rose-900/80 border-rose-500/50 text-rose-100'
        }`}>
          {message.text}
        </div>
      )}
    </div>
  );
};

export default LogProblems;