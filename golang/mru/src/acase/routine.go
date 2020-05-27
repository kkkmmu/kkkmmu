package acase

type Routine struct {
	Device string   `yaml:"Device"`
	Type   string   `yaml:"Type"`
	Delay  int      `yaml:"Delay"`
	API    string   `yaml:"API"`
	Params []string `yaml:"Params"`
	Expect string   `yaml:"Expect"`
}
