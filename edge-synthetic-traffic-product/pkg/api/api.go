package api

const (
	// BasePath is the base path for all API routes.
	BasePath = "/api/v1"

	// Agent routes
	RouteRegisterAgent = BasePath + "/agents/register"
	RouteHeartbeat     = BasePath + "/agents/heartbeat"

	// Task routes
	RouteGetTasks      = BasePath + "/tasks"
	RouteSubmitResults = BasePath + "/results"
)

// RegisterRequest is sent by the agent when it starts up.
type RegisterRequest struct {
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
}

// RegisterResponse contains the assigned Agent ID.
type RegisterResponse struct {
	AgentID string `json:"agentId"`
}

// HeartbeatRequest is sent periodically by the agent.
type HeartbeatRequest struct {
	AgentID string `json:"agentId"`
}
