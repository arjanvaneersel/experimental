package sum_test

import (
	"fmt"

	"github.com/arjanvaneersel/experimental/testing/sum"
)

func ExampleInts() {
	s := sum.Ints(1, 2, 3, 4, 5)
	fmt.Println("The sum of one to five is", s)
	// Output:
	// The sum of one to five is 15
}
