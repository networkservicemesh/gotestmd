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
package consumer1

import (
	"os"
	"path/filepath"

	"github.com/stretchr/testify/suite"

	"github.com/networkservicemesh/gotestmd/pkg/suites/shell"
	"github.com/networkservicemesh/gotestmd/test-examples/producer"
)

type Suite struct {
	shell.Suite
	producerSuite producer.Suite
}

func (s *Suite) SetupSuite() {
	suite.Run(s.T(), &s.producerSuite)
	dir := filepath.Join(os.Getenv("GOPATH"), "src", "/github.com/networkservicemesh/gotestmd/examples/Producer/Consumer1")
	r := s.Runner(dir)
	r.Run(`echo "I'm the first consumer"`)
}
func (s *Suite) Test() {
}
```