package acase

type Task struct {
	Name     string    `yaml:"Name"`
	Enable   bool      `yaml:"Enable"`
	Routines []Routine `yaml:"Routines"`
}
