/*
Package backoff provides two simple backoff and jitter implementations for retrying operations
as described in https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/

Example:

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

*/
package backoff

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Retry the given function n times jittering between max and min time.Duration
func Retry(attempts int64, min, max time.Duration, f func() error) (err error) {
	return retry(false, attempts, max, min, f)
}

// RetryDec like Retry but increase the jitter based on the last jitter value
func RetryDec(attempts int64, min, max time.Duration, f func() error) (err error) {
	return retry(true, attempts, max, min, f)
}

func retry(dec bool, attempts int64, min, max time.Duration, f func() error) (err error) {

	var j int64
	mx := int64(max)
	mn := int64(min)

	for {

		// Quit on a successful attempt
		if err = f(); err == nil {
			break
		}

		// Decrease attempt count and quit if this was the last one
		if attempts--; attempts == 0 {
			break
		}

		j = jitter(dec, j, mx, mn, attempts)
		time.Sleep(time.Duration(j))

	}

	return err
}

func jitter(dec bool, current int64, mx, mn, count int64) (j int64) {
	if dec {
		if current == 0 {
			current = mn
		}
		j = min(mx, rand.Int63n(3*current-mn)+mn)
	} else {
		j = rand.Int63n(min(mx, mn*pow(2, count)))
	}

	if j < mn {
		return mn
	}
	return j

}

// pow for int64
func pow(a, b int64) int64 {
	var p int64 = 1
	for b > 0 {
		if b&1 != 0 {
			p *= a
		}
		b >>= 1
		a *= a
	}
	return p
}

// min for int64
func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
