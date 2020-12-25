# Leaf C

This example will be generated into a test of the _Tree_ suite.


## Run

```bash
echo "I'm leaf C"
```

# Results

The result of generating a suite from this file will be:

```go
func (s *Suite) TestLeafC() {
	dir := filepath.Join(os.Getenv("GOPATH"), "src", "/github.com/networkservicemesh/gotestmd/examples/Tree/LeafC")
    r := s.Runner(dir)    
	r.Run(`echo "I'm leaf C"`)
}
```