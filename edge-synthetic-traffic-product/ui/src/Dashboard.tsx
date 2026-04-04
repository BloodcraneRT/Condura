import React, { useEffect, useState } from 'react';

interface Agent {
  id: string;
  hostname: string;
  os: string;
  ip: string;
  status: string;
  lastSeen: string;
}

interface Result {
  id: string;
  taskId: string;
  agentId: string;
  timestamp: string;
  success: boolean;
  latencyMs: number;
  errorMsg: string;
}

interface Source {
  id: string;
  name: string;
  target: string;
  isDefault: boolean;
}

const API_URL = 'http://localhost:8080/api/v1';

export const Dashboard: React.FC = () => {
  const [agents, setAgents] = useState<Agent[]>([]);
  const [results, setResults] = useState<Result[]>([]);
  const [sources, setSources] = useState<Source[]>([]);

  // Task form state
  const [selectedAgent, setSelectedAgent] = useState('');
  const [selectedType, setSelectedType] = useState('ping');
  const [selectedSource, setSelectedSource] = useState('');

  // Source form state
  const [newSourceName, setNewSourceName] = useState('');
  const [newSourceTarget, setNewSourceTarget] = useState('');

  const fetchData = () => {
    fetch(`${API_URL}/ui/agents`)
      .then(res => res.json())
      .then(data => setAgents(data))
      .catch(console.error);

    fetch(`${API_URL}/ui/results`)
      .then(res => res.json())
      .then(data => setResults(data))
      .catch(console.error);

    fetch(`${API_URL}/ui/sources`)
      .then(res => res.json())
      .then(data => setSources(data))
      .catch(console.error);
  };

  useEffect(() => {
    fetchData();
    const interval = setInterval(fetchData, 5000);
    return () => clearInterval(interval);
  }, []);

  const handleCreateTask = (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedAgent || !selectedSource) return;

    fetch(`${API_URL}/ui/tasks`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        agentId: selectedAgent,
        type: selectedType,
        target: selectedSource,
        interval: 10,
        enabled: true
      })
    }).then(() => {
      alert("Task scheduled successfully!");
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
    <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <div className="px-4 py-6 sm:px-0">
        <h1 className="text-3xl font-bold text-gray-900 mb-8">EdgeSynth Dashboard</h1>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mb-12">
          {/* Create Task Form */}
          <div className="bg-white shadow sm:rounded-md p-6">
            <h2 className="text-xl font-semibold mb-4 text-gray-700">Schedule Synthetic Task</h2>
            <form onSubmit={handleCreateTask} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700">Agent</label>
                <select
                  className="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md"
                  value={selectedAgent}
                  onChange={(e) => setSelectedAgent(e.target.value)}
                  required
                >
                  <option value="" disabled>Select an agent...</option>
                  {agents.map(a => <option key={a.id} value={a.id}>{a.hostname} ({a.ip})</option>)}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700">Test Type</label>
                <select
                  className="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md"
                  value={selectedType}
                  onChange={(e) => setSelectedType(e.target.value)}
                >
                  <option value="ping">Ping</option>
                  <option value="http">HTTP</option>
                  <option value="tcp">TCP</option>
                  <option value="udp">UDP</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700">Source Target</label>
                <select
                  className="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md"
                  value={selectedSource}
                  onChange={(e) => setSelectedSource(e.target.value)}
                  required
                >
                  <option value="" disabled>Select a source target...</option>
                  {sources.map(s => <option key={s.id} value={s.target}>{s.name} ({s.target})</option>)}
                </select>
              </div>
              <button type="submit" className="inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
                Schedule Task
              </button>
            </form>
          </div>

          {/* Add Source Form */}
          <div className="bg-white shadow sm:rounded-md p-6">
            <h2 className="text-xl font-semibold mb-4 text-gray-700">Add New Source</h2>
            <form onSubmit={handleCreateSource} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700">Name</label>
                <input
                  type="text"
                  required
                  placeholder="e.g. My Custom App"
                  className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                  value={newSourceName}
                  onChange={(e) => setNewSourceName(e.target.value)}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700">Target URL/IP</label>
                <input
                  type="text"
                  required
                  placeholder="e.g. https://my-app.local"
                  className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                  value={newSourceTarget}
                  onChange={(e) => setNewSourceTarget(e.target.value)}
                />
              </div>
              <button type="submit" className="inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500">
                Add Source
              </button>
            </form>
          </div>
        </div>

        <div className="mb-12">
          <h2 className="text-xl font-semibold mb-4 text-gray-700">Connected Agents</h2>
          <div className="bg-white shadow overflow-hidden sm:rounded-md">
            <ul className="divide-y divide-gray-200">
              {agents.length === 0 ? (
                <li className="px-4 py-4 sm:px-6 text-gray-500">No agents connected.</li>
              ) : (
                agents.map(agent => (
                  <li key={agent.id} className="px-4 py-4 sm:px-6 flex items-center justify-between">
                    <div>
                      <p className="text-sm font-medium text-indigo-600 truncate">{agent.hostname}</p>
                      <p className="mt-2 flex items-center text-sm text-gray-500">
                        {agent.os} - {agent.ip}
                      </p>
                    </div>
                    <div className="ml-2 flex-shrink-0 flex flex-col items-end">
                      <p className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">
                        {agent.status}
                      </p>
                      <p className="mt-2 text-sm text-gray-500">
                        Last seen: {new Date(agent.lastSeen).toLocaleTimeString()}
                      </p>
                    </div>
                  </li>
                ))
              )}
            </ul>
          </div>
        </div>

        <div>
          <h2 className="text-xl font-semibold mb-4 text-gray-700">Recent Test Results</h2>
          <div className="bg-white shadow overflow-hidden sm:rounded-md">
            <ul className="divide-y divide-gray-200">
              {results.length === 0 ? (
                <li className="px-4 py-4 sm:px-6 text-gray-500">No results yet.</li>
              ) : (
                results.map(res => (
                  <li key={res.id} className="px-4 py-4 sm:px-6">
                    <div className="flex items-center justify-between">
                      <p className="text-sm font-medium text-gray-900 truncate">
                        Task: {res.taskId} <span className="text-gray-500 text-xs ml-2">(Agent: {res.agentId.substring(0,8)}...)</span>
                      </p>
                      <div className="ml-2 flex-shrink-0 flex">
                        {res.success ? (
                          <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">Success</span>
                        ) : (
                          <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-red-100 text-red-800">Failed</span>
                        )}
                      </div>
                    </div>
                    <div className="mt-2 sm:flex sm:justify-between">
                      <div className="sm:flex">
                        <p className="flex items-center text-sm text-gray-500">
                          Latency: {res.latencyMs.toFixed(2)} ms
                        </p>
                      </div>
                      <div className="mt-2 flex items-center text-sm text-gray-500 sm:mt-0">
                        <p>{new Date(res.timestamp).toLocaleString()}</p>
                      </div>
                    </div>
                    {!res.success && (
                      <p className="mt-1 text-xs text-red-500">{res.errorMsg}</p>
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
