import { Dashboard } from './Dashboard'
import './index.css'

function App() {
  return (
    <div className="min-h-screen bg-gray-100">
      <nav className="bg-indigo-600 shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex">
              <div className="flex-shrink-0 flex items-center">
                <span className="text-white text-xl font-bold">EdgeSynth</span>
              </div>
            </div>
          </div>
        </div>
      </nav>
      <Dashboard />
    </div>
  )
}

export default App
