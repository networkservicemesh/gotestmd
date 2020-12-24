# Leaf B

This example will be generated into a test of the _SubTree_ suite.

## Run

```bash
echo "I'm leaf B"
```

# Results

The result of generating a suite from this file will be:

```go
func (s *Suite) TestLeafB() {
	dir := filepath.Join(os.Getenv("GOPATH"), "src", "/github.com/networkservicemesh/gotestmd/examples/Tree/SubTree/LeafB")
	r := s.Runner(dir)
	r.Run(`echo "I'm leaf B"`)
}

```