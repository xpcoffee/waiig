package main

import (
	"fmt"
	"monkey/grapher"
	"monkey/repl"
	"os"
	"os/user"
)

func main() {
	runRepl()
}

func runRepl() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello, %s! Welcome to the Monkey programming language!\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}

func graphAst() {
	input := `
	   let hello = fn(x,y) {
	       fn(z) {
	           x + y + z;
	       }
	   }; hello(1,2)(3);
	   `
	graph := grapher.New(input).GetDot()
	fmt.Println(graph)
}
