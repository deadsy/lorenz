package main

import (
	"fmt"

	"github.com/deadsy/lorenz/lorenz"
)

func main() {
	w := lorenz.New(1.0, 0.2)
	dt := 0.00001
	for {
		w.Run(dt)
		fmt.Printf("%s\n", w)
	}
}
