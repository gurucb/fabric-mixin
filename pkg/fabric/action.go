package fabric

import (
	"get.porter.sh/porter/pkg/exec/builder"
)

var _ builder.ExecutableAction = Action{}
var _ builder.BuildableAction = Action{}

type Action struct {
	Name  string
	Steps []Step // using UnmarshalYAML so that we don't need a custom type per action
	// RuntimeConfig runtime.RuntimeConfig
}

// MakeSteps builds a slice of Steps for data to be unmarshaled into.
func (a Action) MakeSteps() interface{} {
	return &[]Step{}
}

// UnmarshalYAML takes any yaml in this form
// ACTION:
// - aws: ...
// and puts the steps into the Action.Steps field
func (a *Action) UnmarshalYAML(unmarshal func(interface{}) error) error {

	results, err := builder.UnmarshalAction(unmarshal, a)
	if err != nil {
		return err
	}

	for actionName, action := range results {
		a.Name = actionName
		for _, result := range action {
			step := result.(*[]Step)
			a.Steps = append(a.Steps, *step...)
		}
		break // There is only 1 action
	}
	return nil
}

func (a Action) GetSteps() []builder.ExecutableStep {
	steps := make([]builder.ExecutableStep, len(a.Steps))
	for i := range a.Steps {
		steps[i] = a.Steps[i]
	}

	return steps
}

var _ builder.ExecutableStep = Step{}
var _ builder.StepWithOutputs = Step{}
var _ builder.SuppressesOutput = Step{}

type Step struct {
	Instruction `yaml:"fabric"`
}

type Instruction struct {
	// Description  string      `yaml:"description",omitempty`
	Licenses          interface{}   `yaml:"license"`
	Dependencies      interface{}   `yaml:"dependencies"`
	SupportedRegions  []string      `yaml:"supportedRegions"`
	TargetEnvironment string        `yaml:"targetEnvironment"`
	PackageId         string        `yaml:"packageId"`
	Arguments         []string      `yaml:"arguments,omitempty"`
	Flags             builder.Flags `yaml:"flags,omitempty"`
	// Outputs        []Output      `yaml:"outputs,omitempty"`
	// SuppressOutput bool          `yaml:"suppress-output,omitempty"`
}

type License struct {
	SKUList []Skus `yaml:"license"`
}
type Skus struct {
	SKUs     []string `yaml:"skus",omitempty`
	Operator string   `yaml:"operator",omitempty`
}
type Depends struct {
	DependencyList []Dependency `yaml:"dependencies"`
}
type Dependency struct {
	Type                        string `yaml:"type"`
	Query                       string `yaml:"query"`
	dependencyCheckCondition    string `yaml:"dependencyCheckCondition"`
	dependencyCheckResultAction string `yaml:"dependencyCheckResultAction"`
}

func (s Step) GetCommand() string {
	return "CompositeSolution"
}

func (s Step) GetWorkingDir() string {
	return ""
}

func (s Step) GetArguments() []string {
	// args := make([]string, 0, len(s.Arguments)+2)

	// // Specify the Service and Operation
	// args = append(args, s.Service)
	// args = append(args, s.Operation)

	// // Append the positional arguments
	// args = append(args, s.Arguments...)

	// return args
	return nil
}

func (s Step) GetFlags() builder.Flags {
	// Always request json formatted output
	// return append(s.Flags, builder.NewFlag("output", "json"))
	return s.Flags
}

func (s Step) GetOutputs() []builder.Output {
	// outputs := make([]builder.Output, len(s.Outputs))
	// for i := range s.Outputs {
	// 	outputs[i] = s.Outputs[i]
	// }
	// return outputs
	return nil
}

func (s Step) SuppressesOutput() bool {
	// return s.SuppressOutput
	return false
}

var _ builder.OutputJsonPath = Output{}

type Output struct {
	Name     string `yaml:"name"`
	JsonPath string `yaml:"jsonPath"`
}

func (o Output) GetName() string {
	return o.Name
}

func (o Output) GetJsonPath() string {
	return o.JsonPath
}
