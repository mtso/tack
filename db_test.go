package tack

import (
	"strings"
	"testing"
)

func TestDb(t *testing.T) {
	testcases := []struct {
		in   string
		want []interface{}
	}{
		{
			"SET ex 10\nGET ex\nUNSET ex\nGET ex\nEND",
			[]interface{}{nil, "10", nil, ErrNotFound, ErrEnd},
		},
		{
			"SET a 10\nSET b 10\nNUMEQUALTO 10\nNUMEQUALTO 20\nSET b 30\nNUMEQUALTO 10\nEND",
			[]interface{}{nil, nil, 2, 0, nil, 1, ErrEnd},
		},
		{
			"BEGIN\nSET a 10\nGET a\nBEGIN\nSET a 20\nGET a\nROLLBACK\nGET a\nROLLBACK\nGET a\nEND",
			[]interface{}{nil, nil, "10", nil, nil, "20", nil, "10", nil, ErrNotFound, ErrEnd},
		},
		{
			"BEGIN\nSET a 30\nBEGIN\nSET a 40\nCOMMIT\nGET a\nROLLBACK\nCOMMIT\nEND",
			[]interface{}{nil, nil, nil, nil, nil, "40", ErrNoTransaction, ErrNoTransaction, ErrEnd},
		},
		{
			"SET a 50\nBEGIN\nGET a\nSET a 60\nBEGIN\nUNSET a\nGET a\nROLLBACK\nGET a\nCOMMIT\nGET a\nEND",
			[]interface{}{nil, nil, "50", nil, nil, nil, ErrNotFound, nil, "60", nil, "60", ErrEnd},
		},
		{
			"SET a 10\nBEGIN\nNUMEQUALTO 10\nBEGIN\nUNSET a\nNUMEQUALTO 10\nROLLBACK\nNUMEQUALTO 10\nCOMMIT\nEND",
			[]interface{}{nil, nil, 1, nil, nil, 0, nil, 1, nil, ErrEnd},
		},
	}

	for _, testcase := range testcases {
		handle := MakeHandler()
		inputs := strings.Split(testcase.in, "\n")

		for i, input := range inputs {
			raw := strings.Split(input, " ")
			cmd := handle[raw[0]]
			arg := convertArgs(raw[1:])
			got := cmd(arg...)

			if got != testcase.want[i] {
				t.Errorf("Expected %q => %q, but got %q", input, got, testcase.want[i])
			}
		}
	}
}

func BenchmarkSet(b *testing.B) {
	handle := MakeHandler()
	set := handle["SET"]
	for n := 0; n < b.N; n++ {
		set(string(n+n), n)
	}
}

// var getHandler = MakeHandler()
// var set = handle["SET"]
// for n := 0; n < 1000; n++ {
// 	set(string(n + n), n)
// }
// func BenchmarkGet(b *testing.B) {
// 	// handle := MakeHandler()
// 	// set := handle["SET"]
// 	// get := handle["GET"]
// 	for n := 0; n < b.N; n++ {
// 		// set(string(n + n), n)
// 		get(string(n % 1100))
// 	}
// }

func convertArgs(input []string) (args []interface{}) {
	args = make([]interface{}, len(input))
	for i := range args {
		args[i] = interface{}(input[i])
	}
	return
}
