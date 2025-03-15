package node

type DeploymentConfig struct {
	Workers map[string]struct {
		Type   string `yaml:"type"`
		Config any    `yaml:"config"`
	}
	StartOrder []string `yaml:"start_order"`
	StopOrder  []string `yaml:"stop_order"`
}
