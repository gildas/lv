package kubectl

type Selector struct {
	Name    string   `json:"name" yaml:"name"`
	Aliases []string `json:"aliases" yaml:"aliases"`
	Label   string   `json:"label" yaml:"label"`
	Usage   string   `json:"usage" yaml:"usage"`
	Charts  []string `json:"charts" yaml:"charts"`
}
