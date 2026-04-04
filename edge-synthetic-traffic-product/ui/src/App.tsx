import { Dashboard } from './Dashboard'
import './index.css'

function App() {
  return (
    <div className="min-h-screen bg-[#f7f7f9] text-[#222222]">
      <nav className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-6 sm:px-8 lg:px-12">
          <div className="flex justify-between h-[72px]">
            <div className="flex items-center">
              <div className="flex-shrink-0 flex items-center gap-3">
                <div className="w-8 h-8 rounded-full bg-gradient-to-tr from-[#FF385C] to-[#E61E4D] flex items-center justify-center shadow-sm">
                  <svg className="w-4 h-4 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={3}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
                  </svg>
                </div>
                <span className="text-[#222222] text-xl font-extrabold tracking-tight">EdgeSynth</span>
              </div>
            </div>
            <div className="flex items-center gap-4">
              <span className="px-4 py-1.5 rounded-full bg-gray-100 text-xs font-bold text-gray-600 tracking-wide">
                Agent Node
              </span>
            </div>
          </div>
        </div>
      </nav>
      <main className="max-w-7xl mx-auto px-6 sm:px-8 lg:px-12 py-10">
        <Dashboard />
      </main>
    </div>
  )
}

export default App
