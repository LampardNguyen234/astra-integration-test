package framework

import (
	"fmt"
	"strings"
	"testing"
)

func TestNode(t *testing.T) {
	n := Describe("Test a module",
		Before(func() {
			fmt.Println("he;;p")
		}),
		Context("00",
			It("001", func() {
				fmt.Println("001")
			}),
			It("002", func() {
				panic("002")
			}),
		),
	)
	n.Run()
	r := n.Logs()
	fmt.Println(strings.Join(r, "\n"))
}
