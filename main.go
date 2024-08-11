package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"golang.org/x/term"
)

var block = rune('█')
var half = rune('¦')
var ball = rune('■')

var config = Config{}
var speedMS = 100

type screenLine = []rune
type screenBuffer = []screenLine
type Points struct {
	x int
	y int
}
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
	c.Screen.Width = width - 1
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
				tmpArr = append(tmpArr, ' ')
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
		// time.Sleep(1 * time.Second)
		time.Sleep(time.Duration(speedMS) * time.Millisecond)
	}
}

type Memory struct {
	// x         int
	// y         int
	Points
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

var directionMap = map[DirectEnum]Points{
	TopLeft:     {-1, -1},
	Top:         {-1, 0},
	TopRight:    {-1, 1},
	MidLeft:     {0, -1},
	MidRight:    {0, 1},
	BottomLeft:  {1, -1},
	Bottom:      {1, 0},
	BottomRight: {1, 1},
}

func findDirection(p Points) DirectEnum {
	// isXNegative := p.x <= 0
	// isYNegative := p.y <= 0
	for idx, val := range directionMap {
		// isValXNeg := val.x <= 0
		// isValYNeg := val.y <= 0
		// if isXNegative == isValXNeg && isYNegative == isValYNeg {
		if p == val {
			return idx
		}
	}
	return 0
}

// TopLeft:     {-4, -4},
// TopRight:    {-4, 4},
// BottomLeft:  {4, -4},
// BottomRight: {4, 4},

func getSideName(x, y int, config Config) SideEnum {
	if x <= 2 {
		return TopSide
	}
	if x >= config.Screen.Height-1 {
		return BottomSide
	}
	if y <= 2 {
		return LeftSide
	}
	if y >= config.Screen.Width-1 {
		return RightSide
	}
	return DefaultSE
}

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
	if side == LeftSide && direction == TopLeft {
		return TopRight
	}
	if side == RightSide && direction == TopRight {
		return TopLeft
	}
	if side == TopSide && direction == BottomLeft {
		return TopLeft
	}
	fmt.Println(SideHelper[side], DirectHelper[direction])
	return DefaultDE
}

func checkBorders(x, y int, config Config) bool {
	if x < 0 || y < 0 {
		// fmt.Println("XY", x, y)
		return true
	}
	if x > config.Screen.Height-1 || y > config.Screen.Width-1 {
		// fmt.Println("WI HEI", x, y, config.Screen.Width, config.Screen.Height)
		return true
	}
	return false
}

func caclulateDirection(x, y int, config Config, memo Memory) (Points, DirectEnum, error) {
	dirX := x - memo.x
	dirY := y - memo.y
	possibleDirections := []Points{}
	possibleDirections = append(possibleDirections, Points{dirX * -1, dirY})
	possibleDirections = append(possibleDirections, Points{dirX, dirY * -1})
	possibleDirections = append(possibleDirections, Points{dirX * -1, dirY * -1})

	for _, val := range possibleDirections {
		nx, ny := val.x, val.y
		if nx+memo.x > 0 && ny+memo.y > 0 && nx+memo.x < config.Screen.Height-1 && ny+memo.y < config.Screen.Width-1 {
			p := Points{nx + memo.x, ny + memo.y}
			newDirName := findDirection(Points{nx, ny / 10})
			if newDirName == 0 {
				fmt.Println(val, memo, p)
				os.Exit(1)
			}
			return p, newDirName, nil
		}
	}
	return Points{}, 0, errors.New("Cant find good direction")
}

func DrawSideBalls(p Points, config Config, screenBuff *screenBuffer, memo *Memory) {

	if p.y < 0 || p.y > config.Screen.Width-1 {
		newY := config.Screen.Width - 1
		if p.y < 0 {
			newY = 0
		}
		*memo = Memory{p, PointsToScreenBuff(Points{p.x, newY}, *screenBuff), memo.direction}
		(*screenBuff)[p.x][newY] = ball
		time.Sleep(time.Duration(speedMS) * time.Millisecond)
		(*screenBuff)[memo.x][newY] = memo.val
	}
}

func BallAnimation(p Points, screenBuff *screenBuffer, memo *Memory) {

	for true {

		(*screenBuff)[memo.x][memo.y] = memo.val

		p.x += directionMap[memo.direction].x
		p.y += directionMap[memo.direction].y * 10

		// fmt.Println(x, y, memo.direction)

		newDirection := memo.direction

		if isBorder := checkBorders(p.x, p.y, config); isBorder == true {
			point, newDir, _ := caclulateDirection(p.x, p.y, config, *memo)
			newDirection = newDir
			DrawSideBalls(p, config, screenBuff, memo)
			p.x = point.x
			p.y = point.y
			*memo = Memory{p, PointsToScreenBuff(p, *screenBuff), newDirection}
		} else {
			// *memo = Memory{p, (*screenBuff)[p.x][p.y], memo.direction}
			*memo = Memory{p, PointsToScreenBuff(p, *screenBuff), memo.direction}
		}
		(*screenBuff)[p.x][p.y] = ball

		time.Sleep(time.Duration(speedMS) * time.Millisecond)
	}
}

func drawPlayerBlock(screenBuff *screenBuffer) {
	for x := 0; x < 10; x++ {
		for y := 0; y < 2; y++ {
			(*screenBuff)[x+10][y] = block
		}
	}
}

func PointsToScreenBuff(p Points, screenBuff screenBuffer) rune {
	return screenBuff[p.x][p.y]
}

func main() {
	var screenBuff = screenBuffer{}
	config = NewConfig()
	GetSize(&config)

	InitLevel(&screenBuff)
	PrintLevel(screenBuff)
	drawPlayerBlock(&screenBuff)

	p := Points{config.Screen.Height / 2, config.Screen.Width / 2}
	gameMemory := Memory{p, PointsToScreenBuff(p, screenBuff), TopLeft}

	go Animation(&screenBuff)

	go BallAnimation(p, &screenBuff, &gameMemory)
	time.Sleep(25 * time.Minute)
}
