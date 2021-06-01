package main

import "fmt"

var resetfg = ccfg(255, 255, 255)
var resetbg = ccbg(0, 0, 0)

//colorcode_foregrown
func ccfg(r int, g int, b int) string {
	return "\u001b[38;2;" + fmt.Sprint(r) + ";" + fmt.Sprint(g) + ";" + fmt.Sprint(b) + "m"
}
func ccbg(r int, g int, b int) string {
	return "\u001b[48;2;" + fmt.Sprint(r) + ";" + fmt.Sprint(g) + ";" + fmt.Sprint(b) + "m"
}
