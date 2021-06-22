package yml

type YML struct {
	Name        string              `yaml:"name,omitempty" json:"name"`
	Version     string              `yaml:"version,omitempty" json:"version"`
	DisplayName string              `yaml:"displayName,omitempty" json:"displayName"`
	Triggers    map[string]*Trigger `yaml:"triggers,omitempty" json:"triggers"`
	Variables   map[string]string   `yaml:"variables,omitempty" json:"variables"`
	Stages      []*Stage            `yaml:"stages,omitempty" json:"stages"`
}

type Trigger struct {
	AutoCancel     bool       `yaml:"autoCancel,omitempty" json:"autoCancel,omitempty"`
	Timeout        string     `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	Branches       *Condition `yaml:"branches,omitempty" json:"branches,omitempty"`
	Tags           *Condition `yaml:"tags,omitempty" json:"tags,omitempty"`
	Paths          *Condition `yaml:"paths,omitempty" json:"paths,omitempty"`
	Notes          *Condition `yaml:"notes,omitempty" json:"notes,omitempty"`
	CommitMessages *Condition `yaml:"commitMessages,omitempty" json:"commitMessages,omitempty"`
}

type Condition struct {
	Include []string `yaml:"include,omitempty" json:"include,omitempty"`
	Exclude []string `yaml:"exclude,omitempty" json:"exclude,omitempty"`
}

type Stage struct {
	Stage       string  `yaml:"stage" json:"stage"`
	Name        string  `yaml:"name,omitempty" json:"name"`
	DisplayName string  `yaml:"displayName,omitempty" json:"displayName"`
	Steps       []*Step `yaml:"steps,omitempty" json:"steps"`
}

type Step struct {
	Step            string            `yaml:"step" json:"step"`
	DisplayName     string            `yaml:"displayName,omitempty" json:"displayName"`
	Name            string            `yaml:"name,omitempty" json:"name"`
	Environments    map[string]string `yaml:"environments,omitempty" json:"environments"`
	Commands        []interface{}     `yaml:"commands,omitempty" json:"commands"`
	DependsOn       []string          `yaml:"dependsOn,omitempty" json:"dependsOn"`
	Image           string            `yaml:"image,omitempty" json:"image"`
	Artifacts       []*Artifact       `yaml:"artifacts,omitempty" json:"artifacts"`
	DependArtifacts []*DependArtifact `yaml:"dependArtifacts,omitempty" json:"dependArtifacts"`
}

type Artifact struct {
	Name       string `yaml:"name,omitempty" json:"name"`
	Scope      string `yaml:"scope,omitempty" json:"scope"`
	Path       string `yaml:"path,omitempty" json:"path"`
	Repository string `yaml:"repository,omitempty" json:"repository"`
	Value      string `yaml:"value,omitempty" json:"value"`
}

type DependArtifact struct {
	BuildName string `yaml:"buildName,omitempty" json:"buildName"`
	StageName string `yaml:"stageName,omitempty" json:"stageName"`
	StepName  string `yaml:"stepName,omitempty" json:"stepName"`

	Type       string `yaml:"type,omitempty" json:"type"`
	Repository string `yaml:"repository,omitempty" json:"repository"`
	Name       string `yaml:"name,omitempty" json:"name"`
	Target     string `yaml:"target,omitempty" json:"target"`
	IsForce    bool   `yaml:"isForce,omitempty" json:"isForce"`

	SourceStage string `yaml:"sourceStage,omitempty" json:"sourceStage"`
	SourceStep  string `yaml:"sourceStep,omitempty" json:"sourceStep"`
}
