# Leaf A

This example will be generated into a test of the _Tree_ suite.

## Run

```bash
echo "I'm leaf A"
```

# Results

The generated result of this example is:

```go
func (s *Suite) TestLeafA() {
	dir := filepath.Join(os.Getenv("GOPATH"), "src", "/github.com/networkservicemesh/gotestmd/examples/Tree/LeafA")
    r := s.Runner(dir)    
	
	r.Run(`echo "I'm leaf A"`)
}
```