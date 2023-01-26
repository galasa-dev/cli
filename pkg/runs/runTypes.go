/*
 * Copyright contributors to the Galasa project
 */
package runs

type TestRun struct {
	Name      string            `yaml:"name" json:"name"`
	Bundle    string            `yaml:"bundle" json:"bundle"`
	Class     string            `yaml:"class" json:"class"`
	Stream    string            `yaml:"stream" json:"stream"`
	Status    string            `yaml:"status" json:"status"`
	Result    string            `yaml:"result" json:"result"`
	Overrides map[string]string `yaml:"overrides" json:"overrides"`
	Tests     []TestMethod      `yaml:"tests" json:"tests"`
}

type TestMethod struct {
	Method string `yaml:"name" json:"name"`
	Result string `yaml:"result" json:"result"`
}

func DeepClone(original map[string]*TestRun) map[string]*TestRun {
	new := make(map[string]*TestRun)
	for k, v := range original {
		new[k] = v
	}

	return new
}
