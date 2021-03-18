# Requires

[This example links to its parent](../../Bidirecitonal)

## Run

```bash
echo Running example1...
```

## Cleanup

```bash
echo Terminating example1...
```

# Results

For this example gotestmd generates this:

```go
func (s *Suite) TestExample1() {
	r := s.Runner("examples/Bidirecitonal/Example1")
	s.T().Cleanup(func() {
		r.Run(`echo Terminating example1...`)
	})
	r.Run(`echo Running example1...`)
}
```