# Tree Example

Example structure:
```
Tree
├── SubTree
│   └── LeafB
├── LeafA
└── LeafC
```

This file will be converted by `go-testmark` into golang [testify suite](https://github.com/stretchr/testify#suite-package).

The resulting file will include only setup and cleanup logic and also all dependent examples as tests. See more in section `Results`.

## Includes

- [Sub Tree](./SubTree)
- [Leaf A](./LeafA)
- [Leaf C](./LeafC)

## Run

The following command just creates a resource folder.

```bash
MY_TEST_DIR=resources 
echo "mkdir ${MY_TEST_DIR}"
```


## Cleanup

The following command just deletes the resource folder.

```bash
rm -rf ${MY_TEST_DIR}
```

# Results

The generated result of this example is:

```go
package tree

import (
	"os"
	"path/filepath"

	"github.com/networkservicemesh/gotestmd/pkg/suites/shell"
)

type Suite struct {
	shell.Suite
}

func (s *Suite) SetupSuite() {
	s.Suite.SetupSuite()

	dir := filepath.Join(os.Getenv("GOPATH"), "src", "/github.com/networkservicemesh/gotestmd/examples/Tree")
	r := s.Runner(dir)
	s.T().Cleanup(func() {
		r.Run(`rm -rf ${MY_TEST_DIR}`)

	})
	r.Run(`MY_TEST_DIR=resources 
echo "mkdir ${MY_TEST_DIR}"`)
}

func (s *Suite) TestLeafA() {
	dir := filepath.Join(os.Getenv("GOPATH"), "src", "/github.com/networkservicemesh/gotestmd/examples/Tree/LeafA")
	r := s.Runner(dir)

	r.Run(`echo "I'm leaf A"`)
}

func (s *Suite) TestLeafC() {
	dir := filepath.Join(os.Getenv("GOPATH"), "src", "/github.com/networkservicemesh/gotestmd/examples/Tree/LeafC")
	r := s.Runner(dir)

	r.Run(`echo "I'm leaf C"`)
}

```