package parser

type ClosureCheck struct {
	vars     map[string]bool
	closures []string
}

func NewClosureCheck() *ClosureCheck {
	return &ClosureCheck{
		vars:     map[string]bool{},
		closures: nil,
	}
}

func (c *ClosureCheck) AddVar(label string) {
	c.vars[label] = true
}
func (c *ClosureCheck) Visit(label string) bool {
	var closure bool
	if c.vars[label] == false {
		c.closures = append(c.closures, label)
		closure = true
	}
	c.AddVar(label)
	return closure
}
