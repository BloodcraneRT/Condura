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
  ttfbMs: number;
  dnsTimeMs: number;
  connectTimeMs: number;
  hops: number;
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
  const [taskConfig, setTaskConfig] = useState(''); // Extra config like interface for PCAP

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
        config: taskConfig,
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
    <div className="flex flex-col gap-10">
      <div className="flex flex-col gap-8">
        <h1 className="text-4xl font-extrabold tracking-tight text-[#222222]">Overview</h1>

        {/* Metrics Overview */}
        {metrics && (
          <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
            <div className="bg-white shadow-sm border border-gray-100 rounded-2xl p-6 flex flex-col justify-center">
              <p className="text-xs font-bold text-gray-500 uppercase tracking-widest mb-1">Total Executions</p>
              <p className="text-4xl font-extrabold text-[#222222]">{metrics.totalTests}</p>
            </div>
            <div className="bg-white shadow-sm border border-gray-100 rounded-2xl p-6 flex flex-col justify-center">
              <p className="text-xs font-bold text-gray-500 uppercase tracking-widest mb-1">Health Score</p>
              <p className={`text-4xl font-extrabold ${metrics.totalTests > 0 && metrics.totalSuccesses / metrics.totalTests > 0.9 ? 'text-emerald-500' : 'text-rose-500'}`}>
                {metrics.totalTests > 0 ? Math.round((metrics.totalSuccesses / metrics.totalTests) * 100) : 0}%
              </p>
            </div>
            <div className="bg-white shadow-sm border border-gray-100 rounded-2xl p-6 flex flex-col justify-center">
              <p className="text-xs font-bold text-gray-500 uppercase tracking-widest mb-1">Avg Latency</p>
              <p className="text-4xl font-extrabold text-[#222222]">{metrics.avgLatencyMs.toFixed(1)}<span className="text-xl font-medium text-gray-400 ml-1">ms</span></p>
            </div>
            <div className="bg-white shadow-sm border border-gray-100 rounded-2xl p-6 flex flex-col justify-center">
              <p className="text-xs font-bold text-gray-500 uppercase tracking-widest mb-1">Data Throughput</p>
              <p className="text-4xl font-extrabold text-[#222222]">
                {((metrics.totalBytesRecv + metrics.totalBytesSent) / 1024 / 1024).toFixed(1)}<span className="text-xl font-medium text-gray-400 ml-1">MB</span>
              </p>
            </div>
          </div>
        )}

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Create Task Form */}
          <div className="bg-white shadow-md border border-gray-100 rounded-3xl p-8 transition-shadow duration-300">
            <h2 className="text-xl font-bold mb-6 text-[#222222] tracking-tight">Schedule New Test</h2>
            <form onSubmit={handleCreateTask} className="space-y-6">
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">Test Type</label>
                <select
                  className="block w-full border-gray-200 bg-white py-3 px-4 focus:outline-none focus:ring-2 focus:ring-rose-500 focus:border-rose-500 sm:text-sm rounded-xl border shadow-sm transition-shadow"
                  value={selectedType}
                  onChange={(e) => setSelectedType(e.target.value)}
                >
                  <option value="ping">Ping</option>
                  <option value="http">HTTP (Web Load)</option>
                  <option value="tcp">TCP</option>
                  <option value="udp">UDP</option>
                  <option value="dns">DNS Resolution</option>
                  <option value="traceroute">Traceroute (Path)</option>
                  <option value="download">Data Download</option>
                  <option value="upload">Data Upload</option>
                  <option value="pcap_replay">PCAP Replay</option>
                  <option value="ostinato">Traffic Stream (Ostinato)</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">Source Target</label>
                <select
                  className="block w-full border-gray-200 bg-white py-3 px-4 focus:outline-none focus:ring-2 focus:ring-rose-500 focus:border-rose-500 sm:text-sm rounded-xl border shadow-sm transition-shadow"
                  value={selectedSource}
                  onChange={(e) => setSelectedSource(e.target.value)}
                  required
                >
                  <option value="" disabled>Select a target...</option>
                  {sources.map(s => <option key={s.id} value={s.target}>{s.name} ({s.target})</option>)}
                </select>
              </div>
              {selectedType === 'pcap_replay' && (
                <div>
                  <label className="block text-sm font-semibold text-gray-700 mb-2">Network Interface (e.g., eth0)</label>
                  <input
                    type="text"
                    placeholder="eth0"
                    className="block w-full border-gray-200 bg-white py-3 px-4 focus:outline-none focus:ring-2 focus:ring-rose-500 focus:border-rose-500 sm:text-sm rounded-xl border shadow-sm transition-shadow"
                    value={taskConfig}
                    onChange={(e) => setTaskConfig(e.target.value)}
                  />
                </div>
              )}
              {selectedType === 'ostinato' && (
                <div>
                  <label className="block text-sm font-semibold text-gray-700 mb-2">Ostinato Stream Config (JSON)</label>
                  <textarea
                    placeholder={'{\n  "protocol": "udp",\n  "packetSize": 512,\n  "pps": 100,\n  "duration": 5\n}'}
                    className="block w-full border-gray-200 bg-white py-3 px-4 focus:outline-none focus:ring-2 focus:ring-rose-500 focus:border-rose-500 sm:text-sm rounded-xl border shadow-sm transition-shadow"
                    rows={5}
                    value={taskConfig}
                    onChange={(e) => setTaskConfig(e.target.value)}
                  />
                </div>
              )}
              <button type="submit" className="w-full mt-4 justify-center py-3.5 px-4 text-sm font-bold rounded-xl text-white bg-gradient-to-r from-rose-500 to-rose-600 hover:from-rose-600 hover:to-rose-700 shadow-md hover:shadow-lg transition-all focus:ring-2 focus:ring-offset-2 focus:ring-rose-500">
                Run Task
              </button>
            </form>
          </div>

          {/* Add Source Form */}
          <div className="bg-white shadow-md border border-gray-100 rounded-3xl p-8 transition-shadow duration-300">
            <h2 className="text-xl font-bold mb-6 text-[#222222] tracking-tight">New Target Source</h2>
            <form onSubmit={handleCreateSource} className="space-y-6">
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">Name</label>
                <input
                  type="text"
                  required
                  placeholder="e.g. Production API"
                  className="block w-full border-gray-200 bg-white py-3 px-4 focus:outline-none focus:ring-2 focus:ring-rose-500 focus:border-rose-500 sm:text-sm rounded-xl border shadow-sm transition-shadow"
                  value={newSourceName}
                  onChange={(e) => setNewSourceName(e.target.value)}
                />
              </div>
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">Target URL or IP</label>
                <input
                  type="text"
                  required
                  placeholder="e.g. https://api.example.com"
                  className="block w-full border-gray-200 bg-white py-3 px-4 focus:outline-none focus:ring-2 focus:ring-rose-500 focus:border-rose-500 sm:text-sm rounded-xl border shadow-sm transition-shadow"
                  value={newSourceTarget}
                  onChange={(e) => setNewSourceTarget(e.target.value)}
                />
              </div>
              <button type="submit" className="w-full mt-4 justify-center py-3.5 px-4 text-sm font-bold rounded-xl text-gray-700 bg-white border-2 border-gray-200 hover:border-gray-300 hover:bg-gray-50 shadow-sm transition-all">
                Save Target
              </button>
            </form>
          </div>
        </div>

        <div className="mt-6 mb-12">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-2xl font-bold tracking-tight text-[#222222]">Active Tasks</h2>
          </div>
          <div className="bg-white shadow-sm border border-gray-100 rounded-3xl overflow-hidden">
            <ul className="divide-y divide-gray-50">
              {tasks.length === 0 ? (
                <li className="px-8 py-10 text-center text-sm font-medium text-gray-400">No tasks currently scheduled.</li>
              ) : (
                tasks.map(task => (
                  <li key={task.id} className="px-8 py-6 hover:bg-gray-50 transition-colors flex justify-between items-center group">
                    <div className="flex items-center gap-4">
                      <div className="w-12 h-12 bg-rose-50 rounded-2xl flex items-center justify-center border border-rose-100">
                        <svg className="w-5 h-5 text-rose-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" /></svg>
                      </div>
                      <div>
                        <p className="text-sm font-bold text-[#222222] uppercase tracking-wide">{task.type}</p>
                        <p className="text-xs font-medium text-gray-500 mt-1.5">{task.target}</p>
                        <p className="text-[11px] text-gray-400 mt-1 font-medium">Runs every {task.interval}s</p>
                      </div>
                    </div>
                    <button
                      onClick={() => handleDeleteTask(task.id)}
                      className="opacity-0 group-hover:opacity-100 px-4 py-2 bg-red-50 text-red-600 text-xs font-bold rounded-xl hover:bg-red-100 transition-all"
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
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-2xl font-bold tracking-tight text-[#222222]">Recent Traffic Logs</h2>
          </div>
          <div className="bg-white shadow-sm border border-gray-100 rounded-3xl overflow-hidden">
            <ul className="divide-y divide-gray-50">
              {results.length === 0 ? (
                <li className="px-8 py-10 text-center text-sm font-medium text-gray-400">No tests executed yet.</li>
              ) : (
                results.map(res => (
                  <li key={res.id} className="px-8 py-6 hover:bg-gray-50 transition-colors">
                    <div className="flex items-center justify-between mb-4">
                      <div className="flex items-center gap-3">
                        {res.success ? (
                          <div className="w-2 h-2 rounded-full bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.5)]"></div>
                        ) : (
                          <div className="w-2 h-2 rounded-full bg-rose-500 shadow-[0_0_8px_rgba(244,63,94,0.5)]"></div>
                        )}
                        <p className="text-sm font-bold text-[#222222]">
                          Task <span className="text-gray-400 font-medium ml-1">#{res.taskId.substring(0,8)}</span>
                        </p>
                      </div>
                      <span className="text-xs font-semibold text-gray-400 bg-gray-50 px-3 py-1 rounded-lg border border-gray-100">{new Date(res.timestamp).toLocaleString()}</span>
                    </div>
                    <div className="flex items-center gap-6 flex-wrap">
                      <div className="flex flex-col">
                        <span className="text-[10px] uppercase font-bold text-gray-400 tracking-wider mb-1">Latency</span>
                        <span className="text-sm font-bold text-[#222222]">{res.latencyMs.toFixed(1)}ms</span>
                      </div>
                      {res.ttfbMs > 0 && (
                        <div className="flex flex-col">
                          <span className="text-[10px] uppercase font-bold text-gray-400 tracking-wider mb-1">TTFB</span>
                          <span className="text-sm font-bold text-[#222222]">{res.ttfbMs.toFixed(1)}ms</span>
                        </div>
                      )}
                      {res.dnsTimeMs > 0 && (
                        <div className="flex flex-col">
                          <span className="text-[10px] uppercase font-bold text-gray-400 tracking-wider mb-1">DNS</span>
                          <span className="text-sm font-bold text-[#222222]">{res.dnsTimeMs.toFixed(1)}ms</span>
                        </div>
                      )}
                      {res.connectTimeMs > 0 && (
                        <div className="flex flex-col">
                          <span className="text-[10px] uppercase font-bold text-gray-400 tracking-wider mb-1">Connect</span>
                          <span className="text-sm font-bold text-[#222222]">{res.connectTimeMs.toFixed(1)}ms</span>
                        </div>
                      )}
                      {res.bytesRecv > 0 && (
                        <div className="flex flex-col">
                          <span className="text-[10px] uppercase font-bold text-gray-400 tracking-wider mb-1">Downloaded</span>
                          <span className="text-sm font-bold text-[#222222]">{(res.bytesRecv / 1024 / 1024).toFixed(2)} MB</span>
                        </div>
                      )}
                      {res.bytesSent > 0 && (
                        <div className="flex flex-col">
                          <span className="text-[10px] uppercase font-bold text-gray-400 tracking-wider mb-1">Uploaded</span>
                          <span className="text-sm font-bold text-[#222222]">{(res.bytesSent / 1024 / 1024).toFixed(2)} MB</span>
                        </div>
                      )}
                      {res.hops > 0 && (
                        <div className="flex flex-col">
                          <span className="text-[10px] uppercase font-bold text-gray-400 tracking-wider mb-1">Hops</span>
                          <span className="text-sm font-bold text-[#222222]">{res.hops}</span>
                        </div>
                      )}
                    </div>
                    {!res.success && (
                      <p className="mt-4 text-xs font-medium text-rose-600 bg-rose-50 p-3 rounded-xl border border-rose-100 whitespace-pre-wrap">{res.errorMsg}</p>
                    )}
                    {res.success && res.errorMsg && (
                      <div className="mt-4 text-xs font-mono text-gray-600 bg-gray-50 p-4 rounded-xl border border-gray-100 overflow-auto max-h-40">
                        <pre className="text-[10px] leading-relaxed">{res.errorMsg}</pre>
                      </div>
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
