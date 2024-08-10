package main

import (
	"fmt"
	"time"

	"golang.org/x/term"
)

var block = rune('█')
var half = rune('¦')
var ball = rune('■')

var config = Config{}

type screenLine = []rune
type screenBuffer = []screenLine

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
	width, height, err := term.GetSize(0)
	if err != nil {
		return
	}
	c.Screen.Width = width - 5
	c.Screen.Height = height / 2
}

func InitLevel(screenBuff *screenBuffer) {
	halfBlockPlace := config.Screen.Width / 2

	for h := 0; h < config.Screen.Height; h++ {
		tmpArr := []rune{}
		for w := 0; w < config.Screen.Width; w++ {
			if w == halfBlockPlace {
				tmpArr = append(tmpArr, half)
			} else {
				tmpArr = append(tmpArr, '.')
			}
		}
		*screenBuff = append(*screenBuff, tmpArr)
	}
}

func ClearLevel(screenBuff screenBuffer) {
	for _, scrn := range screenBuff {
		for range scrn {
			fmt.Print(" ")
		}
		fmt.Println()
	}
}
func PrintLevel(screenBuff screenBuffer) {
	for _, scrn := range screenBuff {
		for _, row := range scrn {
			fmt.Print(string(row))
		}
		fmt.Println()
	}
}

func Animation(screenBuff *screenBuffer) {
	for true {
		// ClearLevel(*screenBuff)
		fmt.Print("\033[H")
		fmt.Print("\033[2J")
		PrintLevel(*screenBuff)
		time.Sleep(1 * time.Second)
	}
}

type Memory struct {
	x         int
	y         int
	val       rune
	direction DirectEnum
}

type DirectEnum int

const (
	TopLeft DirectEnum = iota
	Top
	TopRight
	MidLeft
	MidRight
	BottomLeft
	Bottom
	BottomRight
)

type DirectionXY struct {
	x int
	y int
}

var direction = map[DirectEnum]DirectionXY{
	TopLeft:     {-1, -1},
	Top:         {-1, 0},
	TopRight:    {-1, 1},
	MidLeft:     {0, -1},
	MidRight:    {0, 1},
	BottomLeft:  {1, -1},
	Bottom:      {1, 0},
	BottomRight: {1, 1},
}

func AnimateBall(config Config, screenBuff *screenBuffer) {
	x := config.Screen.Height / 2
	y := config.Screen.Width / 2
	gameMemory := Memory{x, y, (*screenBuff)[x][y], BottomRight}
	// (*screenBuff)[x][y] = ball
	// time.Sleep(1 * time.Second)
	for true {

		(*screenBuff)[gameMemory.x][gameMemory.y] = gameMemory.val
		// x += 1
		// y -= 1
		// x += direction[BottomRight].x
		// y += direction[BottomRight].y
		x += direction[gameMemory.direction].x
		y += direction[gameMemory.direction].y

		gameMemory = Memory{x, y, (*screenBuff)[x][y], gameMemory.direction}
		if len(*screenBuff) >= x || len((*screenBuff)[x]) >= y {
			newDirection := gameMemory.direction + 2%7

			x -= direction[gameMemory.direction].x
			y -= direction[gameMemory.direction].y

			x += direction[newDirection].x
			y += direction[newDirection].y

			gameMemory = Memory{x, y, (*screenBuff)[x][y], newDirection}

		}
		(*screenBuff)[x][y] = ball
		time.Sleep(1 * time.Second)
	}
}

func main() {
	var screenBuff = screenBuffer{}
	config = NewConfig()
	GetSize(&config)

	InitLevel(&screenBuff)
	PrintLevel(screenBuff)

	go Animation(&screenBuff)

	// tmpRune := BeforeStep{}
	go AnimateBall(config, &screenBuff)
	// for x, scrn := range screenBuff {
	// 	for y := range scrn {
	// 		// fmt.Print(string(row))
	//
	// 		if x == y {
	// 			if (BeforeStep{}) == tmpRune {
	// 				tmpRune = BeforeStep{x, y, screenBuff[x][y]}
	// 			} else {
	// 				screenBuff[tmpRune.x][tmpRune.y] = tmpRune.val
	// 				tmpRune = BeforeStep{x, y, screenBuff[x][y]}
	// 			}
	// 			screenBuff[x][y] = ball
	// 		}
	// 		// else{
	// 		// 	scree
	// 		// }
	//
	// 	}
	// 	time.Sleep(1 * time.Second)
	// }
	time.Sleep(25 * time.Second)
}
