## Description

This example shows that examples can have not only Tree or Producer-Consumer structure,  but also they can be bidirectional. 
This is useful to navigating from one example to another to be in context.

# Includes

- [This examples links to](./Example1) 

# Result


```go
// Code generated by gotestmd DO NOT EDIT.
package bidirecitonal

import (
	"github.com/networkservicemesh/gotestmd/pkg/suites/shell"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	shell.Suite
}

func (s *Suite) SetupSuite() {
	parents := []interface{}{&s.Suite}
	for _, p := range parents {
		if v, ok := p.(suite.TestingSuite); ok {
			v.SetT(s.T())
		}
		if v, ok := p.(suite.SetupAllSuite); ok {
			v.SetupSuite()
		}
	}
}
func (s *Suite) TestExample1() {
	r := s.Runner("examples/Bidirecitonal/Example1")
	s.T().Cleanup(func() {
		r.Run(`echo Terminating example1...`)
	})
	r.Run(`echo Running example1...`)
}

```