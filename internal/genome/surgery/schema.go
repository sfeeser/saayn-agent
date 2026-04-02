package surgery

const (
	SurgeryPlanVersion  = "v1"
	SurgeryPatchVersion = "v1"
)

const (
	ActionReplaceNode = "replace_node"
)

const (
	ReasonPrimaryTarget = "primary_target"
	ReasonKnownCaller   = "known_caller"
	ReasonSameFile      = "same_file"
	ReasonSamePackage   = "same_package"
)

const (
	RiskLow          = "LOW"
	RiskMedium       = "MEDIUM"
	RiskHigh         = "HIGH"
	RiskCatastrophic = "CATASTROPHIC"
)

// TargetAnchor provides stable identifiers so downstream tools don't lose the node.
type TargetAnchor struct {
	UUID      string `yaml:"uuid" json:"uuid"`
	PublicID  string `yaml:"public_id" json:"public_id"`
	FilePath  string `yaml:"file_path" json:"file_path"`
	NodeType  string `yaml:"node_type" json:"node_type"`
	LogicHash string `yaml:"logic_hash,omitempty" json:"logic_hash,omitempty"`
}

// SurgeryPlan is the staging document for code modification.
type SurgeryPlan struct {
	Version        string          `yaml:"version"`
	Intent         string          `yaml:"intent"`
	PlannerModel   string          `yaml:"planner_model"`
	Target         TargetAnchor    `yaml:"target"`
	Context        []ContextNode   `yaml:"context"`
	ImpactAnalysis *ImpactAnalysis `yaml:"impact_analysis,omitempty"`
}

// ContextNode holds the raw code the LLM needs to read.
type ContextNode struct {
	TargetAnchor `yaml:",inline"`
	Reason       string `yaml:"reason"`
	ReasonDetail string `yaml:"reason_detail,omitempty"`
	SourceCode   string `yaml:"source_code"`
}

// ImpactAnalysis represents the 4-Pillar Risk Heuristic.
type ImpactAnalysis struct {
	RiskLevel     string      `yaml:"risk_level"`
	RiskReasoning string      `yaml:"risk_reasoning"`
	RiskFactors   RiskFactors `yaml:"risk_factors"`
}

type RiskFactors struct {
	Visibility        string `yaml:"visibility"` // "private", "package", "public_api", "external_boundary"
	BoundaryCrossings int    `yaml:"boundary_crossings"`
	TransitiveDepth   int    `yaml:"transitive_depth"`
	TotalBlastRadius  int    `yaml:"total_blast_radius"`
}

// SurgeryPatch is the constrained, machine-checkable instruction set for the AST Surgeon.
type SurgeryPatch struct {
	Version    string         `yaml:"version" json:"version"`
	TargetNode TargetAnchor   `yaml:"target_node" json:"target_node"`
	Action     string         `yaml:"action" json:"action"`
	NewCode    string         `yaml:"new_code" json:"new_code"`
	Rationale  string         `yaml:"rationale" json:"rationale"`
	Metadata   *PatchMetadata `yaml:"metadata,omitempty" json:"metadata,omitempty"`
}

type PatchMetadata struct {
	Status        string          `yaml:"status" json:"status"`
	Attempts      int             `yaml:"attempts" json:"attempts"`
	ReviewHistory []ReviewAttempt `yaml:"review_history,omitempty" json:"review_history,omitempty"`
}

type ReviewAttempt struct {
	Attempt   int      `yaml:"attempt" json:"attempt"`
	Stage     string   `yaml:"stage" json:"stage"`
	Errors    []string `yaml:"errors,omitempty" json:"errors,omitempty"`
	AutoFixed bool     `yaml:"auto_fixed,omitempty" json:"auto_fixed,omitempty"`
}
