package model

// Registry represents the full genome.json file structure.
type Registry struct {
	ProjectName     string           `json:"project_name"`
	RootPath        string           `json:"root_path"`
	LastMutatedUUID string           `json:"last_mutated_uuid,omitempty"`
	Nodes           map[string]*Node `json:"nodes"` // Using *Node pointers for efficiency
}

// NewRegistry initializes a fresh birth certificate for the project.
func NewRegistry(project, targetDir string) *Registry {
	return &Registry{
		ProjectName: project,
		RootPath:    targetDir,
		Nodes:       make(map[string]*Node),
	}
}

// SurgeryRequest represents a mutation instruction for the SAAYN agent.
type SurgeryRequest struct {
	TargetUUID  string
	NewLogic    string
	Registry    *Registry
	ProjectRoot string
}
