# Consumer 2

## Requires

- [Producer](../)

## Run

```bash
echo "I'm the second consumer"
```

# Results

The result of generating a suite is:
```go
package consumer2

import (
	"github.com/networkservicemesh/gotestmd/pkg/suites/shell"
	"github.com/networkservicemesh/gotestmd/test-examples/producer"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	shell.Suite
	producerSuite producer.Suite
}

func (s *Suite) SetupSuite() {
	suite.Run(s.T(), &s.producerSuite)
	r := s.Runner("examples/Producer/Consumer2")
	r.Run(`echo "I'm the second consumer"`)
}
func (s *Suite) Test() {}
```