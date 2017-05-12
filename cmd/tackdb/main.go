package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mtso/tack"
)

func main() {
	handle := tack.CreateDb().GetCommands()

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("> ")
	for scanner.Scan() {
		args := strings.Split(scanner.Text(), " ")

		if cmd, ok := handle[strings.ToUpper(args[0])]; !ok {
			fmt.Println("UNRECOGNIZED COMMAND")
		} else {
			if resp, err := cmd(args[1:]...); err == tack.ErrEnd {
				break
			} else if err != nil {
				fmt.Println(err)
			} else if resp != "" {
				fmt.Println(resp)
			}
		}
		fmt.Printf("> ")
	}
}
