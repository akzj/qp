package qp

type closureCheck struct {
	vars     map[string]bool
	closures []string
}

func newClosureCheck() *closureCheck {
	return &closureCheck{
		vars:     map[string]bool{},
		closures: nil,
	}
}

func (c *closureCheck) addVar(label string) {
	c.vars[label] = true
}
func (c *closureCheck) visit(label string) bool {
	var closure bool
	if c.vars[label] == false {
		c.closures = append(c.closures, label)
		closure = true
	}
	c.addVar(label)
	return closure
}
