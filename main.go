package main

import (
	"fmt"

	"golang.org/x/term"
)

var block = rune('█')
var half = rune('¦')
var config = Config{}

type screenLine = []rune
type screenBuffer = []screenLine

var screenBuff = screenBuffer{}

type Screen struct {
	Width  int
	Height int
}
type Config struct {
	Screen Screen
}

func NewConfig() Config {
	return Config{
		Screen: Screen{0, 0},
	}
}
func GetSize(c *Config) {
	// if term.IsTerminal(0) {
	// 	println("in a term")
	// } else {
	// 	println("not in a term")
	// }
	width, height, err := term.GetSize(0)
	if err != nil {
		return
	}
	c.Screen.Width = width - 5
	c.Screen.Height = height
	// println("width:", width, "height:", height)
}
func main() {
	// fmt.Println("ekekekk")
	config = NewConfig()
	GetSize(&config)
	fmt.Println(config.Screen.Width)
	fmt.Println(config.Screen.Height)
	// s := make([]int, 0, 10)
	// for(unsigned int ch = 128; ch < 256; ch++)
	//   {
	//       printf("%d = %c\t\t", ch, ch);
	//   }
	fmt.Println(block)
	halfBlockPlace := config.Screen.Width / 2

	for h := 0; h < config.Screen.Height/2; h++ {
		// fmt.Print(h)
		tmpArr := []rune{}
		for w := 0; w < config.Screen.Width; w++ {
			if w == halfBlockPlace {
				// fmt.Print(half)
				tmpArr = append(tmpArr, half)
			} else {
				// fmt.Print(" ")
				tmpArr = append(tmpArr, ' ')
			}
		}
		// fmt.Println()
		screenBuff = append(screenBuff, tmpArr)
	}
	// fmt.Println(screenBuff)

	for _, scrn := range screenBuff {
		for _, row := range scrn {
			fmt.Print(string(row))
		}
		fmt.Println()
	}
	// for ch := 128; ch < 256; ch++ {
	// 	fmt.Printf("%d = %c\t\t", ch, ch)
	// 	if ch%7 == 0 {
	// 		fmt.Println()
	// 	}
	// }
}
