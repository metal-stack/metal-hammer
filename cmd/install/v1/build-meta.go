package v1

type BuildMeta struct {
	Version  string `json:"buildVersion" yaml:"buildVersion"`
	Date     string `json:"buildDate" yaml:"buildDate"`
	SHA      string `json:"buildSHA" yaml:"buildSHA"`
	Revision string `json:"buildRevision" yaml:"buildRevision"`

	IgnitionVersion string `json:"ignitionVersion" yaml:"ignitionVersion"`
}
