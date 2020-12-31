# Consumer 3

# Requires

- [Producer](../)

## Includes

- [Consumer 1](../Consumer1)

## Run

```bash
echo "I'm the third consumer"
```

# Results

The result of generating a suite is:
```go
package consumer3

import (
	"os"
	"path/filepath"

	"github.com/networkservicemesh/gotestmd/pkg/suites/shell"
	"github.com/networkservicemesh/gotestmd/test-examples/producer"
	"github.com/networkservicemesh/gotestmd/test-examples/producer/consumer1"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	shell.Suite
	producerSuite  producer.Suite
	consumer1Suite consumer1.Suite
}

func (s *Suite) SetupSuite() {
	suite.Run(s.T(), &s.producerSuite)
	suite.Run(s.T(), &s.consumer1Suite)

	dir := filepath.Join(os.Getenv("GOPATH"), "/github.com/networkservicemesh/gotestmd/examples/Producer/Consumer3")
	r := s.Runner(dir)

	r.Run(`echo "I'm the third consumer"`)
}

func (s *Suite) Test() {}
```