package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mtso/tack"
)

func main() {
	handle := tack.MakeHandler()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		str := strings.Split(scanner.Text(), " ")

		if cmd, ok := handle[strings.ToUpper(str[0])]; !ok {
			fmt.Println("UNRECOGNIZED COMMAND")
		} else {
			args := convertArgs(str[1:])

			if resp := cmd(args...); resp == tack.ErrEnd {
				break
			} else if resp != nil {
				fmt.Println(resp)
			}
		}
	}
}

func convertArgs(input []string) (args []interface{}) {
	args = make([]interface{}, len(input))
	for i := range args {
		args[i] = interface{}(input[i])
	}
	return
}
