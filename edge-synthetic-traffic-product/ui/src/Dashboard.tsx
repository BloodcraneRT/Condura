import React, { useEffect, useState } from 'react';

interface Result {
  id: string;
  taskId: string;
  timestamp: string;
  success: boolean;
  latencyMs: number;
  errorMsg: string;
  bytesSent: number;
  bytesRecv: number;
}

interface Source {
  id: string;
  name: string;
  target: string;
  isDefault: boolean;
}

interface Metrics {
  totalTests: number;
  totalSuccesses: number;
  totalFailures: number;
  avgLatencyMs: number;
  totalBytesSent: number;
  totalBytesRecv: number;
}

const API_URL = '/api/v1';

export const Dashboard: React.FC = () => {
  const [results, setResults] = useState<Result[]>([]);
  const [sources, setSources] = useState<Source[]>([]);
  const [tasks, setTasks] = useState<any[]>([]);
  const [metrics, setMetrics] = useState<Metrics | null>(null);

  // Task form state
  const [selectedType, setSelectedType] = useState('ping');
  const [selectedSource, setSelectedSource] = useState('');

  // Source form state
  const [newSourceName, setNewSourceName] = useState('');
  const [newSourceTarget, setNewSourceTarget] = useState('');

  const fetchData = () => {
    fetch(`${API_URL}/ui/tasks`)
      .then(res => res.json())
      .then(data => setTasks(data))
      .catch(console.error);

    fetch(`${API_URL}/ui/results`)
      .then(res => res.json())
      .then(data => setResults(data))
      .catch(console.error);

    fetch(`${API_URL}/ui/sources`)
      .then(res => res.json())
      .then(data => setSources(data))
      .catch(console.error);

    fetch(`${API_URL}/ui/metrics`)
      .then(res => res.json())
      .then(data => setMetrics(data))
      .catch(console.error);
  };

  useEffect(() => {
    fetchData();
    const interval = setInterval(fetchData, 5000);
    return () => clearInterval(interval);
  }, []);

  const handleCreateTask = (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedSource) return;

    fetch(`${API_URL}/ui/tasks`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        type: selectedType,
        target: selectedSource,
        interval: 10,
        enabled: true
      })
    }).then(() => {
      alert("Task scheduled successfully!");
      fetchData();
    }).catch(console.error);
  };

  const handleDeleteTask = (id: string) => {
    fetch(`${API_URL}/ui/tasks?id=${id}`, {
      method: 'DELETE'
    }).then(() => {
      fetchData();
    }).catch(console.error);
  };

  const handleCreateSource = (e: React.FormEvent) => {
    e.preventDefault();
    if (!newSourceName || !newSourceTarget) return;

    fetch(`${API_URL}/ui/sources`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        name: newSourceName,
        target: newSourceTarget,
        isDefault: false
      })
    }).then(() => {
      setNewSourceName('');
      setNewSourceTarget('');
      fetchData();
      alert("Source added successfully!");
    }).catch(console.error);
  };

  return (
    <div className="max-w-7xl mx-auto py-10 sm:px-6 lg:px-8">
      <div className="px-4 sm:px-0">
        <h1 className="text-4xl font-bold tracking-tight text-black mb-10">Dashboard</h1>

        {/* Metrics Overview */}
        {metrics && (
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-10">
            <div className="bg-white border border-gray-200 rounded-xl p-6 flex flex-col items-center justify-center">
              <p className="text-sm font-semibold text-gray-500 uppercase tracking-wider mb-2">Total Tests</p>
              <p className="text-3xl font-black text-black">{metrics.totalTests}</p>
            </div>
            <div className="bg-white border border-gray-200 rounded-xl p-6 flex flex-col items-center justify-center">
              <p className="text-sm font-semibold text-gray-500 uppercase tracking-wider mb-2">Success Rate</p>
              <p className={`text-3xl font-black ${metrics.totalTests > 0 && metrics.totalSuccesses / metrics.totalTests > 0.9 ? 'text-green-600' : 'text-red-600'}`}>
                {metrics.totalTests > 0 ? Math.round((metrics.totalSuccesses / metrics.totalTests) * 100) : 0}%
              </p>
            </div>
            <div className="bg-white border border-gray-200 rounded-xl p-6 flex flex-col items-center justify-center">
              <p className="text-sm font-semibold text-gray-500 uppercase tracking-wider mb-2">Avg Latency</p>
              <p className="text-3xl font-black text-black">{metrics.avgLatencyMs.toFixed(1)}ms</p>
            </div>
            <div className="bg-white border border-gray-200 rounded-xl p-6 flex flex-col items-center justify-center text-center">
              <p className="text-sm font-semibold text-gray-500 uppercase tracking-wider mb-2">Data Transferred</p>
              <p className="text-3xl font-black text-black">
                {((metrics.totalBytesRecv + metrics.totalBytesSent) / 1024 / 1024).toFixed(1)} <span className="text-xl">MB</span>
              </p>
            </div>
          </div>
        )}

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-12">
          {/* Create Task Form */}
          <div className="bg-white border border-gray-200 rounded-xl p-8">
            <h2 className="text-lg font-semibold mb-6 text-black tracking-tight">Schedule Task</h2>
            <form onSubmit={handleCreateTask} className="space-y-5">
              <div>
                <label className="block text-sm font-medium text-gray-900 mb-2">Test Type</label>
                <select
                  className="block w-full border-gray-300 bg-gray-50 py-3 px-4 focus:outline-none focus:ring-0 focus:border-black sm:text-sm rounded-lg border"
                  value={selectedType}
                  onChange={(e) => setSelectedType(e.target.value)}
                >
                  <option value="ping">Ping</option>
                  <option value="http">HTTP</option>
                  <option value="tcp">TCP</option>
                  <option value="udp">UDP</option>
                  <option value="download">Data Download</option>
                  <option value="upload">Data Upload</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-900 mb-2">Source Target</label>
                <select
                  className="block w-full border-gray-300 bg-gray-50 py-3 px-4 focus:outline-none focus:ring-0 focus:border-black sm:text-sm rounded-lg border"
                  value={selectedSource}
                  onChange={(e) => setSelectedSource(e.target.value)}
                  required
                >
                  <option value="" disabled>Select a target...</option>
                  {sources.map(s => <option key={s.id} value={s.target}>{s.name} ({s.target})</option>)}
                </select>
              </div>
              <button type="submit" className="w-full mt-2 justify-center py-3 px-4 text-sm font-semibold rounded-lg text-white bg-black hover:bg-gray-900 transition-colors">
                Run Task
              </button>
            </form>
          </div>

          {/* Add Source Form */}
          <div className="bg-white border border-gray-200 rounded-xl p-8">
            <h2 className="text-lg font-semibold mb-6 text-black tracking-tight">New Target Source</h2>
            <form onSubmit={handleCreateSource} className="space-y-5">
              <div>
                <label className="block text-sm font-medium text-gray-900 mb-2">Name</label>
                <input
                  type="text"
                  required
                  placeholder="e.g. Production API"
                  className="block w-full border-gray-300 bg-gray-50 py-3 px-4 focus:outline-none focus:ring-0 focus:border-black sm:text-sm rounded-lg border"
                  value={newSourceName}
                  onChange={(e) => setNewSourceName(e.target.value)}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-900 mb-2">Target URL or IP</label>
                <input
                  type="text"
                  required
                  placeholder="e.g. https://api.example.com"
                  className="block w-full border-gray-300 bg-gray-50 py-3 px-4 focus:outline-none focus:ring-0 focus:border-black sm:text-sm rounded-lg border"
                  value={newSourceTarget}
                  onChange={(e) => setNewSourceTarget(e.target.value)}
                />
              </div>
              <button type="submit" className="w-full mt-2 justify-center py-3 px-4 text-sm font-semibold rounded-lg text-black bg-gray-100 hover:bg-gray-200 transition-colors">
                Save Target
              </button>
            </form>
          </div>
        </div>

        <div className="mb-12">
          <h2 className="text-xl font-bold tracking-tight text-black mb-6">Active Tasks</h2>
          <div className="bg-white border border-gray-200 rounded-xl overflow-hidden">
            <ul className="divide-y divide-gray-100">
              {tasks.length === 0 ? (
                <li className="px-6 py-8 text-center text-sm text-gray-500">No tasks currently scheduled.</li>
              ) : (
                tasks.map(task => (
                  <li key={task.id} className="px-6 py-5 hover:bg-gray-50 transition-colors flex justify-between items-center">
                    <div>
                      <p className="text-sm font-bold text-gray-900 uppercase tracking-wide">{task.type}</p>
                      <p className="text-xs text-gray-500 mt-1">{task.target}</p>
                      <p className="text-[10px] text-gray-400 mt-1">Runs every {task.interval}s • ID: {task.id.substring(0,8)}</p>
                    </div>
                    <button
                      onClick={() => handleDeleteTask(task.id)}
                      className="px-3 py-1 bg-red-50 text-red-600 text-xs font-bold rounded-md hover:bg-red-100 transition-colors"
                    >
                      Delete
                    </button>
                  </li>
                ))
              )}
            </ul>
          </div>
        </div>

        <div>
          <h2 className="text-xl font-bold tracking-tight text-black mb-6">Recent Traffic Results</h2>
          <div className="bg-white border border-gray-200 rounded-xl overflow-hidden">
            <ul className="divide-y divide-gray-100">
              {results.length === 0 ? (
                <li className="px-6 py-8 text-center text-sm text-gray-500">No tests executed yet.</li>
              ) : (
                results.map(res => (
                  <li key={res.id} className="px-6 py-5 hover:bg-gray-50 transition-colors">
                    <div className="flex items-center justify-between mb-2">
                      <p className="text-sm font-semibold text-gray-900">
                        Task <span className="text-gray-500 font-normal">{res.taskId.substring(0,8)}</span>
                      </p>
                      <div>
                        {res.success ? (
                          <span className="px-3 py-1 text-[10px] uppercase tracking-wider font-bold rounded-md bg-green-50 text-green-700">Success</span>
                        ) : (
                          <span className="px-3 py-1 text-[10px] uppercase tracking-wider font-bold rounded-md bg-red-50 text-red-700">Failed</span>
                        )}
                      </div>
                    </div>
                    <div className="flex items-center justify-between text-xs text-gray-500">
                      <div className="flex items-center gap-4">
                        <span>Latency: <strong className="text-gray-900 font-semibold">{res.latencyMs.toFixed(2)}ms</strong></span>
                        {res.bytesRecv > 0 && <span>Downloaded: <strong className="text-gray-900 font-semibold">{(res.bytesRecv / 1024 / 1024).toFixed(2)} MB</strong></span>}
                        {res.bytesSent > 0 && <span>Uploaded: <strong className="text-gray-900 font-semibold">{(res.bytesSent / 1024 / 1024).toFixed(2)} MB</strong></span>}
                      </div>
                      <span>{new Date(res.timestamp).toLocaleString()}</span>
                    </div>
                    {!res.success && (
                      <p className="mt-3 text-xs text-red-600 bg-red-50 p-2 rounded-md">{res.errorMsg}</p>
                    )}
                  </li>
                ))
              )}
            </ul>
          </div>
        </div>
      </div>
    </div>
  );
};
