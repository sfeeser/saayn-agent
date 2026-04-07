package model

// Specbook represents the Genotype: the desired blueprint for a project.
type Specbook struct {
	Project string     `yaml:"project"`
	Version string     `yaml:"version"`
	Nodes   []SpecNode `yaml:"nodes"`
}

// SpecNode defines the requirements for a single unit of code (Function, Struct, etc.)
type SpecNode struct {
	Name      string   `yaml:"name"`       // Unique ID in the spec (e.g. "Task")
	Path      string   `yaml:"path"`       // Destination file (e.g. "internal/models/task.go")
	Type      string   `yaml:"type"`       // struct, function, method
	DependsOn []string `yaml:"depends_on"` // Names of nodes that must exist first
	Logic     string   `yaml:"logic"`      // The "Soul": Natural language intent

	// Optional fields for structured generation
	Receiver string  `yaml:"receiver,omitempty"` // For methods (e.g. "*Calculator")
	Fields   []Field `yaml:"fields,omitempty"`   // For structs
}

type Field struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}
