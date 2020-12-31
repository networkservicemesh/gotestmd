# Hello World Example

This file will be converted by `gotestmd` into go [testify suite](https://github.com/stretchr/testify#suite-package).

The file will include only one test because this example doesn't depend on others and others don't depend on this example.

## Run

```bash
echo "Hello world!"
```

# Results

The result of generating a suite is:
```go
package helloworld

import (
	"os"
	"path/filepath"

	"github.com/networkservicemesh/gotestmd/pkg/suites/shell"
)

type Suite struct {
	shell.Suite
}

func (s *Suite) SetupSuite() {
	dir := filepath.Join(os.Getenv("GOPATH"), "/github.com/networkservicemesh/gotestmd/examples/HelloWorld")
	r := s.Runner(dir)
	r.Run(`echo "Hello world!"`)
}
func (s *Suite) Test() {}
```