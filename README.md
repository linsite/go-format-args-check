# go-format-args-check

The Go linter `go-format-args-check` checks that printf-like functions are named with `f` at the end.

For example, `Printf` should have exactly 3 args, but got 2. You'll get a error.

```go
package main

import "log"

func main() {
	log.Printf("name: %s, age: %d, test: %s", "test", 3)
}
```


```
main.go:6:2: formatting function 'Printf' args should match % count
```
