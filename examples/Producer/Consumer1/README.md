# Consumer 1

## Requires

- [Producer](../)

## Run

```bash
echo "I'm the first consumer"
```

# Results
The result of generating a suite is:
```go
// Code generated by gotestmd DO NOT EDIT.
package consumer1

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
	r := s.Runner("examples/Producer/Consumer1")
	r.Run(`echo "I'm the first consumer"`)
}
func (s *Suite) Test() {}
```
Note: the result has not producer setup/teardown logic because this Consumer is used by [Consumer3](../Consumer3) that contains required dependency.