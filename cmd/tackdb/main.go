package main

import (
	"bufio"
	"os"
	"github.com/mtso/tack"

	"fmt"
	// "strings"
)

func main() {
	db := tack.CreateDb()
	scanner := bufio.NewScanner(os.Stdin)

	db.Set("hello", 4)

	fmt.Println(db.Get("hello"))
	fmt.Println(db.NumEqualTo(4))
	fmt.Println(db.NumEqualTo(5))

	// for scanner.Scan() {
	// 	str := strings.Split(scanner.Text(), " ")
	// 	db.Set(str[1], str[2])
	// }

	if err := scanner.Err(); err != nil {
		
	}
}
