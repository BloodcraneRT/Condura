import { Dashboard } from './Dashboard'
import './index.css'

function App() {
  return (
    <div className="min-h-screen bg-[#f9f9f9]">
      <nav className="bg-black text-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <div className="flex-shrink-0 flex items-center gap-2">
                <div className="w-4 h-4 bg-white rounded-sm"></div>
                <span className="text-white text-xl font-semibold tracking-tight">EdgeSynth</span>
              </div>
            </div>
            <div className="flex items-center">
              <span className="text-xs uppercase tracking-wider font-semibold text-gray-400">Control Plane</span>
            </div>
          </div>
        </div>
      </nav>
      <Dashboard />
    </div>
  )
}

export default App
