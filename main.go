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
	DefaultDE DirectEnum = iota
	TopLeft
	Top
	TopRight
	MidLeft
	Center
	MidRight
	BottomLeft
	Bottom
	BottomRight
)

var DirectHelper = map[DirectEnum]string{
	DefaultDE:   "defa",
	TopLeft:     "TopLeft",
	Top:         "Top",
	TopRight:    "TopRight",
	MidLeft:     "MidLeft",
	Center:      "Center",
	MidRight:    "MidRight",
	BottomLeft:  "BottomLeft",
	Bottom:      "Bottom",
	BottomRight: "BottomRight",
}

type SideEnum int

const (
	DefaultSE SideEnum = iota
	TopSide
	LeftSide
	RightSide
	BottomSide
)

var SideHelper = map[SideEnum]string{
	DefaultSE:  "defa",
	TopSide:    "TopSide",
	LeftSide:   "LeftSide",
	RightSide:  "RightSide",
	BottomSide: "BottomSide",
}

type DirectionXY struct {
	x int
	y int
}

var directionMap = map[DirectEnum]DirectionXY{
	TopLeft:     {-2, -2},
	Top:         {-2, 0},
	TopRight:    {-2, 2},
	MidLeft:     {0, -2},
	MidRight:    {0, 2},
	BottomLeft:  {2, -2},
	Bottom:      {2, 0},
	BottomRight: {2, 2},
}

func getSideName(x, y int, config Config) SideEnum {
	if x <= 1 {
		return TopSide
	}
	if x >= config.Screen.Height-1 {
		return BottomSide
	}
	if y <= 1 {
		return LeftSide
	}
	if y >= config.Screen.Width-1 {
		return RightSide
	}
	return DefaultSE
}

// var directionMap = map[

func getDirection(side SideEnum, direction DirectEnum) DirectEnum {
	if side == TopSide && direction == TopRight {
		return BottomRight
	}
	if side == RightSide && direction == BottomRight {
		return BottomLeft
	}
	if side == BottomSide && direction == BottomLeft {
		return TopLeft
	}
	if side == TopSide && direction == TopLeft {
		return BottomLeft
	}
	if side == LeftSide && direction == BottomLeft {
		return BottomRight
	}
	if side == BottomSide && direction == BottomRight {
		return TopRight
	}
	return DefaultDE
}

func AnimateBall(config Config, screenBuff *screenBuffer) {
	x := config.Screen.Height / 2
	y := config.Screen.Width / 2
	// x := 31
	// y := 82
	gameMemory := Memory{x, y, (*screenBuff)[x][y], BottomRight}
	// (*screenBuff)[x][y] = ball
	// time.Sleep(1 * time.Second)
	for true {

		(*screenBuff)[gameMemory.x][gameMemory.y] = gameMemory.val
		// x += 1
		// y -= 1
		// x += direction[BottomRight].x
		// y += direction[BottomRight].y
		x += directionMap[gameMemory.direction].x
		y += directionMap[gameMemory.direction].y

		gameMemory = Memory{x, y, (*screenBuff)[x][y], gameMemory.direction}
		fmt.Println(
			len(*screenBuff), len((*screenBuff)[x]),
			x, y,
		)

		if x >= len(*screenBuff)-1 || y >= len((*screenBuff)[x])-1 || x == 0 || y == 0 {
			// newDirection := (gameMemory.direction + 3) % 9
			side := getSideName(x, y, config)
			newDirection := getDirection(side, gameMemory.direction)
			fmt.Println("newdirrrrr", DirectHelper[gameMemory.direction], DirectHelper[newDirection], SideHelper[side])

			x -= directionMap[gameMemory.direction].x
			y -= directionMap[gameMemory.direction].y

			x += directionMap[newDirection].x
			y += directionMap[newDirection].y

			gameMemory = Memory{x, y, (*screenBuff)[x][y], newDirection}

		}
		(*screenBuff)[x][y] = ball
		// time.Sleep(1 * time.Second)
		// time.Sleep(250 * time.Millisecond)
		time.Sleep(125 * time.Millisecond)
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
	time.Sleep(25 * time.Minute)
}
