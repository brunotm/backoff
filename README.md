# backoff
### backoff provides two simple backoff and jitter implementations for retrying operations
### As described in:
 - https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/

### Example:
```go
package main

import (
	"fmt"
	"time"

	"github.com/brunotm/backoff"
)

func main() {
	count := 0
	err := backoff.Retry(
		5,                      // attempts
		100*time.Millisecond,   // min
		3*time.Second,          // max
		func() error {
			count++
			fmt.Println("Count: ", count)
			if count == 5 {
				return nil
			}
			return fmt.Errorf("op error")
		})
	fmt.Println(err)
}
```

Written by Bruno Moura <brunotm@gmail.com>
