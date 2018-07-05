package randgen

import "fmt"

func ExampleGenerateStr() {
	res := GenerateStr(0, "[TEST-", "]")
	fmt.Println(res)
	// Output:
	// [TEST-]
}
