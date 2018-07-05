// run tests: go test -v
// run selected: go test -v -run="Invalid"
// documentation server
// godoc -http=":7799" -goroot=$GOPATH
// from console
// godoc github.com/kdlug/mail/email
//
// Code Coverage
// go test -cover
// Write coverage profile to a file (-coverprofile flag automatically sets -cover to enable coverage analysis):
// go test -coverprofile=coverage.out
//
// Analyze results
// go tool cover -func=coverage.out
// go tool cover -html=coverage.out
//
// Heat maps
// go test -covermode=count -coverprofile=count.out fmt
// go tool cover -func=count.out
// go tool cover -html=count.out

package email

import (
	"fmt"
	"testing"
)

func ExampleValidateEmail() {
	res, err := ValidateEmail("john.doe@gmail.com")
	fmt.Println(res)
	fmt.Println(err)
	// Output:
	// true
	// <nil>
}

func TestValidateInvalidEmail(t *testing.T) {
	email := "john.doe"

	if res, _ := ValidateEmail(email); res == true {
		t.Errorf("ValidateEmail(%q) = %v", email, true)
	}
}
func TestValidateNonExistentEmail(t *testing.T) {
	email := "john.doe@gmail"

	if res, _ := ValidateEmail(email); res == true {
		t.Errorf("ValidateEmail(%q) = %v", email, false)
	}
}

// Table driven test
func TestValidateEmail(t *testing.T) {
	var tests = []struct {
		input string
		want  bool
	}{
		{"john.doe@gmail.com", true},
		{"john.doe", false},
		{"", false},
		{"john.doe@gmail", false},
	}

	for _, test := range tests {
		if got, _ := ValidateEmail(test.input); got != test.want {
			t.Errorf("validateEmail(%q) = %v", test.input, test.want)
		}
	}

}

func BenchmarkValidateEmailEmpty(b *testing.B) {
	// run function b.N times
	for n := 0; n < b.N; n++ {
		ValidateEmail("")
	}
}

func BenchmarkValidateEmailValid(b *testing.B) {
	// run function b.N times
	for n := 0; n < b.N; n++ {
		ValidateEmail("jogn.doe@gmail.com")
	}
}

// Benchmark functions start with Benchmark
// Each benchmark must execute the code under test b.N time
// Each benchmark is run for a minimum of 1 second by default.
// If the second has not elapsed when the Benchmark function returns, the value of b.N is increased in the sequence 1, 2, 5, 10, 20, 50, … and the function run again.
// go test -bench=.
