# Sub Tree

This example will be generated into a test of a Suite.

## Includes

[Leaf B](./LeafB)

## Run

```bash
echo "I'm sub tree"
```

## Cleanup

```bash
echo "Sub tree is done"
```

# Results

The result of generating a suite from this file will be:

```go
package subtree

import (
	"os"
	"path/filepath"

	"github.com/networkservicemesh/gotestmd/pkg/suites/shell"
)

type Suite struct {
	shell.Suite
}

func (s *Suite) SetupSuite() {
	dir := filepath.Join(os.Getenv("GOPATH"), "/github.com/networkservicemesh/gotestmd/examples/Tree/SubTree")
	r := s.Runner(dir)

	s.T().Cleanup(func() {
		r.Run(`echo "Sub tree is done"`)
	})

	r.Run(`echo "I'm sub tree"`)
}

func (s *Suite) TestLeafB() {
	dir := filepath.Join(os.Getenv("GOPATH"), "/Users/user/Projects/NSM/gotestmd/examples/Tree/SubTree/LeafB")
	r := s.Runner(dir)

	r.Run(`echo "I'm leaf B"`)
}
```