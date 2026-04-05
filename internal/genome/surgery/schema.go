package surgery

const (
	SurgeryPlanVersion  = "v1"
	SurgeryPatchVersion = "v1"
)

// Action constants define what the AST Surgeon should do with a node.
const (
	ActionReplaceNode = "replace_node"
	ActionNone        = "none"
)

// Role constants define the node's relationship to the surgery intent.
const (
	RolePrimaryTarget = "primary_target"
	RoleBlastRadius   = "blast_radius"
)

// Status constants define the outcome of the AI's review for a specific node.
const (
	StatusApproved = "approved"
	StatusNoChange = "no_change"
)

// Reason constants for context gathering.
const (
	ReasonPrimaryTarget = "primary_target"
	ReasonKnownCaller   = "known_caller"
	ReasonSameFile      = "same_file"
	ReasonSamePackage   = "same_package"
)

// Risk levels for impact analysis.
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

// ContextNode holds the raw code the LLM needs to read for the blast radius.
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

// SurgeryPatchSet is the V5 multi-file instruction set for the AST Surgeon.
type SurgeryPatchSet struct {
	Version string          `yaml:"version" json:"version"`
	Records []PatchRecord   `yaml:"records" json:"records"`
	Review  *PatchSetReview `yaml:"review,omitempty" json:"review,omitempty"`
}

// PatchSetReview holds the audit trail for the entire planning session.
type PatchSetReview struct {
	Status   string          `yaml:"status" json:"status"`
	Attempts int             `yaml:"attempts" json:"attempts"`
	History  []ReviewAttempt `yaml:"history,omitempty" json:"history,omitempty"`
}

// PatchRecord represents a single file/node modification decision within a multi-file surgery.
type PatchRecord struct {
	TargetNode TargetAnchor `yaml:"target_node" json:"target_node"`
	Role       string       `yaml:"role" json:"role"`     // "primary_target" | "blast_radius"
	Status     string       `yaml:"status" json:"status"` // "approved" | "no_change"
	Action     string       `yaml:"action" json:"action"` // "replace_node" | "none"
	NewCode    string       `yaml:"new_code,omitempty" json:"new_code,omitempty"`
	Rationale  string       `yaml:"rationale" json:"rationale"`
}

// ReviewAttempt tracks individual AST/compiler validation cycles during the planning loop.
type ReviewAttempt struct {
	Attempt   int      `yaml:"attempt" json:"attempt"`
	Stage     string   `yaml:"stage" json:"stage"`
	Errors    []string `yaml:"errors,omitempty" json:"errors,omitempty"`
	AutoFixed bool     `yaml:"auto_fixed,omitempty" json:"auto_fixed,omitempty"`
}
