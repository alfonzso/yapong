package main

import (
	"errors"
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
		time.Sleep(750 * time.Millisecond)
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
	TopLeft:     {-4, -4},
	Top:         {-4, 0},
	TopRight:    {-4, 4},
	MidLeft:     {0, -4},
	MidRight:    {0, 4},
	BottomLeft:  {4, -4},
	Bottom:      {4, 0},
	BottomRight: {4, 4},
}

func findDirection(p Points) DirectEnum {
	for idx, val := range directionMap {
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
	if x <= 0 || y <= 0 {
		fmt.Println("XY", x, y)
		return true
	}
	if x >= config.Screen.Height || y >= config.Screen.Width {
		fmt.Println("WI HEI", x, y, config.Screen.Width, config.Screen.Height)
		return true
	}
	return false
}

func caclulateDirection(x, y int, config Config, memo Memory) (Points, DirectEnum, error) {
	// fmt.Println(x, y, DirectHelper[memo.direction])
	// fmt.Println(x, y, memo.x, memo.y)
	dirX := x - memo.x
	dirY := y - memo.y
	// fmt.Println(dirX, dirY)
	// backDirX := dirX * -1
	// backDirY := dirY * -1
	possibleDirections := []Points{}
	possibleDirections = append(possibleDirections, Points{dirX * -1, dirY})
	possibleDirections = append(possibleDirections, Points{dirX, dirY * -1})
	possibleDirections = append(possibleDirections, Points{dirX * -1, dirY * -1})

	for _, val := range possibleDirections {
		nx, ny := val.x, val.y
		// fmt.Println(nx, ny)
		fmt.Println("----> ", nx, ny, nx+memo.x, ny+memo.y)
		if nx+memo.x >= 0 && ny+memo.y >= 0 && nx+memo.x <= config.Screen.Height && ny+memo.y <= config.Screen.Width {
			p := Points{nx + memo.x, ny + memo.y}
			// newDirName := findDirection(Points{nx, ny})
			newDirName := findDirection(val)
			fmt.Println("......", DirectHelper[newDirName])
			// return val, newDirName, nil
			return p, newDirName, nil
		}
	}
	return Points{}, 0, errors.New("Cant find good direction")
}

func BallAnimation(config Config, screenBuff *screenBuffer) {
	x := config.Screen.Height / 2
	y := config.Screen.Width / 2
	(*screenBuff)[0][0] = rune('x')

	gameMemory := Memory{Points{x, y}, (*screenBuff)[x][y], TopLeft}
	for true {

		// (*screenBuff)[gameMemory.x][gameMemory.y] = gameMemory.val
		// newX := x + directionMap[gameMemory.direction].x
		// newY := y + directionMap[gameMemory.direction].y
		bX := x
		bY := y
		x += directionMap[gameMemory.direction].x
		y += directionMap[gameMemory.direction].y

		// fmt.Println(x, y, gameMemory.direction)

		newDirection := gameMemory.direction

		// isBorder := checkBorders(x, y, config, *screenBuff)
		if isBorder := checkBorders(x, y, config); isBorder == true {
			_, newDir, _ := caclulateDirection(x, y, config, gameMemory)
			newDirection = newDir
			// fmt.Println("border ", isBorder, point)
			// x, y = point.x, point.y

			x = bX + directionMap[newDirection].x
			y = bY + directionMap[newDirection].y
		}
		gameMemory = Memory{Points{x, y}, (*screenBuff)[x][y], newDirection}
		(*screenBuff)[x][y] = ball
		// time.Sleep(1 * time.Second)
		time.Sleep(750 * time.Millisecond)
	}
}

func AnimateBall(config Config, screenBuff *screenBuffer) {
	x := config.Screen.Height / 2
	y := config.Screen.Width / 2
	x = 15
	y = 15
	// gameMemory := Memory{x, y, (*screenBuff)[x][y], BottomRight}
	gameMemory := Memory{Points{x, y}, (*screenBuff)[x][y], TopLeft}
	// (*screenBuff)[x][y] = ball
	// time.Sleep(1 * time.Second)
	for true {

		(*screenBuff)[gameMemory.x][gameMemory.y] = gameMemory.val
		// x += 1
		// y -= 1
		// x += direction[BottomRight].x
		// y += direction[BottomRight].y
		// x += directionMap[gameMemory.direction].x
		// y += directionMap[gameMemory.direction].y
		newX := x + directionMap[gameMemory.direction].x
		newY := y + directionMap[gameMemory.direction].y

		// if newX < 0 {
		// 	newX = x
		// }
		// if newY < 0 {
		// 	newY = y
		// }
		//
		// if newX < 1 && newY < 1 {
		// 	gameMemory.direction = BottomRight
		// }
		// fmt.Println(
		// 	len(*screenBuff), len((*screenBuff)[x]),
		// 	x, y,
		// 	config.Screen.Width,
		// 	config.Screen.Height,
		// 	newY >= config.Screen.Width-1,
		// )
		// xBuffLen := len(*screenBuff) - 1

		// if newX >= xBuffLen || (newX >= 0 && newY >= len((*screenBuff)[newX])-1) || newX < 2 || newY < 2 {
		if newX < 2 || newY < 2 || newX >= config.Screen.Height-1 || newY >= config.Screen.Width-1 {
			// if newX >= len(*screenBuff)-1 || newY >= len((*screenBuff)[newX])-1 {
			// newDirection := (gameMemory.direction + 3) % 9
			side := getSideName(newX, newY, config)
			newDirection := getDirection(side, gameMemory.direction)
			// fmt.Println("newdirrrrr", DirectHelper[gameMemory.direction], DirectHelper[newDirection], SideHelper[side])

			// x -= directionMap[gameMemory.direction].x
			// y -= directionMap[gameMemory.direction].y

			x += directionMap[newDirection].x
			y += directionMap[newDirection].y
			if x < 0 {
				x = 0
			}
			if y < 0 {
				y = 0
			}
			// fmt.Println(x, y)
			gameMemory = Memory{Points{x, y}, (*screenBuff)[x][y], newDirection}

		} else {
			// x += directionMap[gameMemory.direction].x
			// y += directionMap[gameMemory.direction].y
			x = newX
			y = newY
			// if x < 0 {
			// 	x = 0
			// }
			// if y < 0 {
			// 	y = 0
			// }
			gameMemory = Memory{Points{x, y}, (*screenBuff)[x][y], gameMemory.direction}
		}

		(*screenBuff)[x][y] = ball
		// time.Sleep(1 * time.Second)
		time.Sleep(750 * time.Millisecond)
		// time.Sleep(175 * time.Millisecond)
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
	// go AnimateBall(config, &screenBuff)
	go BallAnimation(config, &screenBuff)
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
