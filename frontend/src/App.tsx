import { useState, useEffect } from 'react';
import { Routes, Route, Navigate, useNavigate, useLocation, NavLink } from 'react-router-dom';
import Training from './Training';
import LandingPage from './LandingPage';
import Stats from './Stats';
import Challenge from './Challenge';
import LogProblems from './LogProblems';
import RecentActivity from './RecentActivity';

function App() {
  const [handle, setHandle] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    const storedHandle = sessionStorage.getItem('cf_handle');
    setHandle(storedHandle);
    setLoading(false);
  }, []);

  useEffect(() => {
    if (!loading && !handle && location.pathname !== '/') {
      navigate('/');
    }
    if (!loading && handle && location.pathname === '/') {
      navigate('/training');
    }
  }, [handle, loading, navigate, location]);

  const handleLoginSuccess = (newHandle: string) => {
    sessionStorage.setItem('cf_handle', newHandle);
    setHandle(newHandle);
    navigate('/training');
  };

  const handleLogout = () => {
    sessionStorage.removeItem('cf_handle');
    setHandle(null);
    navigate('/');
  };

  if (loading) return null;

  return (
    <div className="h-screen w-screen overflow-hidden flex flex-col"
      style={{ background: 'linear-gradient(180deg, #0b1220 0%, #0e1628 100%)' }}
    >
      {handle && (
        <nav className="px-6 py-3 border-b border-slate-800 bg-slate-900/50 backdrop-blur-sm text-gray-200 font-mono text-sm tracking-wider flex justify-between items-center z-20">
          <div className="flex items-center">
            <div className="text-2xl font-extrabold tracking-tight bg-gradient-to-br from-white to-sky-500 bg-clip-text text-transparent pr-6 border-r border-slate-700">
              ASCENT
            </div>
            <div className="flex gap-6 ml-6 text-sm font-bold">
              <NavLink 
                to="/training" 
                className={({ isActive }) => 
                  `hover:text-sky-400 transition-colors ${isActive ? 'text-sky-400 font-bold' : 'text-slate-400'}`
                }
              >
                TRAINING
              </NavLink>
              <NavLink 
                to="/random" 
                className={({ isActive }) => 
                  `hover:text-sky-400 transition-colors ${isActive ? 'text-sky-400 font-bold' : 'text-slate-400'}`
                }
              >
                RANDOM
              </NavLink>
              <NavLink 
                to="/log" 
                className={({ isActive }) => 
                  `hover:text-sky-400 transition-colors ${isActive ? 'text-sky-400 font-bold' : 'text-slate-400'}`
                }
              >
                LOG
              </NavLink>
              <NavLink 
                to="/recent" 
                className={({ isActive }) => 
                  `hover:text-sky-400 transition-colors ${isActive ? 'text-sky-400 font-bold' : 'text-slate-400'}`
                }
              >
                RECENT
              </NavLink>
              <NavLink 
                to="/stats" 
                className={({ isActive }) => 
                  `hover:text-sky-400 transition-colors ${isActive ? 'text-sky-400 font-bold' : 'text-slate-400'}`
                }
              >
                STATS
              </NavLink>
            </div>
          </div>
          
          <div className="flex items-center gap-4">
            <span className="text-slate-400 text-xs">
              Logged in as <span className="text-white font-semibold">{handle}</span>
            </span>
            <button 
              onClick={handleLogout}
              className="text-xs text-red-400 hover:text-red-300 transition-colors border border-red-400/20 px-2 py-1 rounded hover:bg-red-400/10"
            >
              LOGOUT
            </button>
          </div>
        </nav>
      )}

      <main className="flex-1 overflow-hidden relative">
        <Routes>
          <Route path="/" element={<LandingPage onSuccess={handleLoginSuccess} />} />
          <Route path="/training" element={handle ? <Training /> : <Navigate to="/" />} />
          <Route path="/random" element={handle ? <Challenge /> : <Navigate to="/" />} />
          <Route path="/log" element={handle ? <LogProblems /> : <Navigate to="/" />} />
          <Route path="/recent" element={handle ? <RecentActivity /> : <Navigate to="/" />} />
          <Route path="/stats" element={handle ? <Stats /> : <Navigate to="/" />} />
          <Route path="*" element={<Navigate to="/" />} />
        </Routes>
      </main>
    </div>
  );
}

export default App;