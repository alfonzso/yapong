package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/MarinX/keylogger"
	"github.com/eiannone/keyboard"
	"github.com/sirupsen/logrus"
	"golang.org/x/term"
)

var block = rune('█')
var half = rune('¦')
var ball = rune('■')

var config = Config{}
var speedMS = 150

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
	c.Screen.Width = (width - 1) // / 2
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
				ch := ' '
				if h == 0 || h == config.Screen.Height-1 {
					ch = '='
				}
				tmpArr = append(tmpArr, ch)
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
func PrintLevel(screenBuff *screenBuffer) {
	for _, scrn := range *screenBuff {
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
		PrintLevel(screenBuff)
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

// def line_intersection(line1, line2):
//     xdiff = (line1[0][0] - line1[1][0], line2[0][0] - line2[1][0])
//     ydiff = (line1[0][1] - line1[1][1], line2[0][1] - line2[1][1])
//
//     def det(a, b):
//         return a[0] * b[1] - a[1] * b[0]
//
//     div = det(xdiff, ydiff)
//     if div == 0:
//        raise Exception('lines do not intersect')
//
//     d = (det(*line1), det(*line2))
//     x = det(d, xdiff) / div
//     y = det(d, ydiff) / div
//     return x, y

func isLinesIntersect(A []Points, B []Points) bool {
	xdiff := Points{A[0].x - A[1].x, B[0].x - B[1].x}
	ydiff := Points{A[0].y - A[1].y, B[0].y - B[1].y}

	det := func(a Points, b Points) int {
		return a.x*b.y - a.y - b.x
	}
	div := det(xdiff, ydiff)
	if div == 0 {
		return false
	}
	return true
	// d := Points{det(A[0], A[1]), det(B[0], B[1])}
	// x := det(d, xdiff) / div
	// y := det(d, ydiff) / div
	// return x, y
}

func direction(p, q, r Points) int {
	return (q.y-p.y)*(r.x-q.x) - (q.x-p.x)*(r.y-q.y)
}

func areCollinearAndOverlapping(a1, b1, a2, b2 Points) bool {
	// # Check if the line segments are collinear
	if direction(a1, b1, a2) == 0 {
		// # Check if the line segments overlap
		if a2.x <= max(a1.x, b1.x) && a2.x >= min(a1.x, b1.x) && a2.y <= max(a1.y, b1.y) && a2.y >= min(a1.y, b1.y) {
			return true
		}
	}
	return false
}

func isintersect(a1, b1, a2, b2 Points) bool {
	// Compute the directions of the four line segments
	d1 := direction(a1, b1, a2)
	d2 := direction(a1, b1, b2)
	d3 := direction(a2, b2, a1)
	d4 := direction(a2, b2, b1)

	// Check if the two line segments intersect
	if ((d1 > 0 && d2 < 0) || (d1 < 0 && d2 > 0)) && ((d3 > 0 && d4 < 0) || (d3 < 0 && d4 > 0)) {
		return true
	}
	// Check if the line segments are collinear && overlapping
	if areCollinearAndOverlapping(a1, b1, a2, b2) || areCollinearAndOverlapping(a2, b2, a1, b1) {
		return true
	}
	return false
}

func BallAnimation(p Points, screenBuff *screenBuffer, memo *Memory, pAX, pAY, pBX, pBY *Points) {

	for true {

		(*screenBuff)[memo.x][memo.y] = memo.val

		p.x += directionMap[memo.direction].x
		p.y += directionMap[memo.direction].y * 10

		// fmt.Println(x, y, memo.direction)

		newDirection := memo.direction

		// x -> 20 -- 30
		// y -> 1  -- 2
		// playaLine := []Points{{15, 1}, {25, 1}}
		ballLine := []Points{{memo.x, memo.y}, {p.x, p.y}}

		AInter := isintersect(*pAX, *pAY, Points{memo.x, memo.y}, Points{p.x, p.y})
		BInter := isintersect(*pBX, *pBY, Points{memo.x, memo.y}, Points{p.x, p.y})
		// res := isLinesIntersect(ballLine, playaLine)
		fmt.Println("..........................", pBX, pBY, ballLine)
		if AInter || BInter {
			// os.Exit(1)
			time.Sleep(5 * time.Second)
		}

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

func KeyBoardHandler() {
	// find keyboard device, does not require a root permission
	keyboard := keylogger.FindKeyboardDevice()

	logrus.Println("Found a keyboard at", keyboard)
	// init keylogger with keyboard
	k, err := keylogger.New(keyboard)
	if err != nil {
		logrus.Error(err)
		return
	}

	defer k.Close()

	logrus.Println("reading ... ...............")
	events := k.Read()
	logrus.Println("reading 1... ...............", events)

	// range of events
	for e := range events {
		logrus.Println(e)
		switch e.Type {
		// EvKey is used to describe state changes of keyboards, buttons, or other key-like devices.
		// check the input_event.go for more events
		case keylogger.EvKey:

			// if the state of key is pressed
			if e.KeyPress() {
				logrus.Println("[event] press key ", e.KeyString())
			}

			// if the state of key is released
			if e.KeyRelease() {
				logrus.Println("[event] release key ", e.KeyString())
			}

			break
		}
	}
	logrus.Println("exited ...............")
}

func cleanupPlayer(screenBuff *screenBuffer, playa int) {
	xLen := len(*screenBuff)
	for i := 0; i < xLen; i++ {
		for y := playa; y < playa+2; y++ {
			(*screenBuff)[i][y] = rune(' ')
		}
	}
}

func drawPlayerBlock(screenBuff *screenBuffer, ws int, isAPlayer bool) (Points, Points) {
	pFirstIndex := ws
	pLen := 5

	player := 0
	if !isAPlayer {
		player = len((*screenBuff)[0]) - 2
	}
	cleanupPlayer(screenBuff, player)
	for x := pFirstIndex; x < pFirstIndex+pLen; x++ {
		for y := player; y < player+2; y++ {
			(*screenBuff)[x][y] = block
		}
	}

	return Points{pFirstIndex, player + 2}, Points{pFirstIndex + pLen, player + 2}
}

func PointsToScreenBuff(p Points, screenBuff screenBuffer) rune {
	return screenBuff[p.x][p.y]
}

var reader = bufio.NewReader(os.Stdin)

func readKey(input chan rune) {
	for true {
		char, _, err := keyboard.GetSingleKey()
		if err != nil {
			panic(err)
		}
		if char == rune('q') {
			os.Exit(1)
		}
		if char == rune('w') || char == rune('s') || char == rune('i') || char == rune('k') {
			input <- char
		}
		if char == rune('r') {
			fmt.Print("\033[H")
			fmt.Print("\033[2J")
		}
		// fmt.Printf("You pressed: %q\r\n", char)
	}
	// scanner := bufio.NewScanner(os.Stdin)
	// logrus.Println("kekek")
	// for scanner.Scan() {
	// 	fmt.Printf("--->>> %s", scanner.Text())
	// }
	// return
	// for true {
	// 	char, _, err := reader.ReadRune()
	// 	if err != nil {
	// 		logrus.Fatal(err)
	// 	}
	// 	input <- char
	// }
}

func main() {
	var screenBuff = screenBuffer{}
	ws, ik := 10, 10

	config = NewConfig()
	GetSize(&config)

	InitLevel(&screenBuff)
	PrintLevel(&screenBuff)
	pAPozX, pAPozY := drawPlayerBlock(&screenBuff, ws, true)
	pBPozX, pBPozY := drawPlayerBlock(&screenBuff, ik, false)

	p := Points{config.Screen.Height / 2, config.Screen.Width / 2}
	gameMemory := Memory{p, PointsToScreenBuff(p, screenBuff), TopLeft}

	go Animation(&screenBuff)

	go BallAnimation(p, &screenBuff, &gameMemory, &pAPozX, &pAPozY, &pBPozX, &pBPozY)
	// logrus.Println(pPozA, pPozB, gameMemory)

	// go KeyBoardHandler()
	input := make(chan rune, 1)
	// fmt.Println("Checking keyboard input...")
	go readKey(input)
	for true {
		select {
		case i := <-input:
			if i == rune('w') && ws > 0 {
				ws += -2
			}
			if i == rune('s') && ws < config.Screen.Height-5 {
				ws += 2
			}
			if i == rune('i') && ik > 0 {
				ik += -2
			}
			if i == rune('k') && ik < config.Screen.Height-5 {
				ik += 2
			}
			pAPozX, pAPozY = drawPlayerBlock(&screenBuff, ws, true)
			pBPozX, pBPozY = drawPlayerBlock(&screenBuff, ik, false)
		}
	}
	time.Sleep(25 * time.Minute)
}
