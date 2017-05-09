package main

import (
	"bufio"
	"os"
	"github.com/mtso/tack"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		tack.Hello(string(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		tack.Hello(err.Error())
	}
}

// import (
// 	"bufio"
// 	"os"
// 	"io"
// 	"fmt"
// )

// func main() {
// 	reader := bufio.NewReader(os.Stdin)

// 	for {
// 		s, err := reader.ReadString('\n')
// 		if err == io.EOF {
// 			return
// 		}
// 		fmt.Println(s)
// 	}
// }

// import (
// 	"fmt"
// 	"github.com/mtso/tack"
// )

// func main() {
// 	n := tack.Node{
// 		Name: "foo",
// 		Value: 4,
// 	}
// 	n.Child[0] = &n
// 	n.Child[1] = &n

// 	fmt.Println(n)
// }