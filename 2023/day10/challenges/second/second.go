package challenge2

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"slices"
	"strings"
)

type Pipe struct {
	Character string
	Right     bool
	Left      bool
	Top       bool
	Down      bool
	x         int
	y         int
}

func Run() {
	input, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	pipes := getPipes(input)

	var startPipe Pipe

	for _, row := range pipes {
		for _, pipe := range row {
			if pipe.Character == "S" {
				startPipe = pipe
			}
		}
	}

	directions := map[string]map[string]int{
		"Right": {
			"x": 1,
		},
		"Left": {
			"x": -1,
		},
		"Top": {
			"y": -1,
		},
		"Down": {
			"y": 1,
		},
	}

	steps := make([]int, 4)
	var loops [][][]int
	for range directions {
		loops = append(loops, [][]int{{startPipe.y, startPipe.x}})
	}

	index := 0

	for _direction, coordinates := range directions {
		direction := _direction
		for _coordinate, _coordinateValue := range coordinates {
			coordinate := _coordinate
			coordinateValue := _coordinateValue
			pipe := startPipe
			for {
				matchDirection, nextDirection, nextPipe := walk(pipes, direction, coordinate, coordinateValue, pipe)

				if !matchDirection {
					index++
					break
				}

				loops[index] = append(loops[index], []int{nextPipe.y, nextPipe.x})

				for key, value := range directions[nextDirection] {
					coordinate = key
					coordinateValue = value
				}

				pipe = nextPipe
				direction = nextDirection
				steps[index]++
			}
		}
	}

	var mainLoop [][]int

	for _, v := range loops {
		if len(v) > len(mainLoop) {
			mainLoop = v
		}
	}

	mainLoopMap := make(map[string]bool)

	for _, v := range mainLoop {
		mainLoopMap[fmt.Sprintf("%d,%d", v[0], v[1])] = true
	}

	truthyS := getTruthySValue(directions, startPipe, pipes)

	reg := regexp.MustCompile(`[S]+`)
	rows := strings.Split(strings.TrimSpace(reg.ReplaceAllString(string(input), truthyS)), "\n")

	var mutableRows [][]rune

	for i := range rows {
		mutableRows = append(mutableRows, []rune(rows[i]))
	}

	for y, row := range mutableRows {
		for x := range row {
			coordinate := fmt.Sprintf("%d,%d", y, x)
			if ok := mainLoopMap[coordinate]; !ok {
				mutableRows[y][x] = rune('.')
			}
		}
	}

	for i := range mutableRows {
		rows[i] = string(mutableRows[i])
	}

	outsideLoopMap := make(map[string]bool)

	for y, row := range rows {
		within := false
		up := false
		for x, _ch := range row {
			ch := string(_ch)

			if ch == "|" {
				within = !within
			}
			if ch == "F" || ch == "L" {
				up = ch == "L"
			}
			if ch == "7" || ch == "J" {
				if up && ch != "J" {
					within = !within
				}

				if !up && ch != "7" {
					within = !within
				}
				up = false
			}

			if !within {
				outsideLoopMap[fmt.Sprintf("%d,%d", y, x)] = true
			}
		}
	}

	area := len(rows) * len(rows[0])

	res := area - unionLen(outsideLoopMap, mainLoopMap)

	fmt.Printf("Inside loop: %d\n", res)
}

func unionLen(set1, set2 map[string]bool) int {
	union := make(map[string]bool)

	for k := range set1 {
		union[k] = true
	}
	for k := range set2 {
		union[k] = true
	}

	return len(union)
}

func walk(rows [][]Pipe, direction, coordinate string, coordinateValue int, pipe Pipe) (bool, string, Pipe) {
	x, y := pipe.x, pipe.y

	if coordinate == "x" {
		x += coordinateValue
	} else {
		y += coordinateValue
	}

	if x < 0 || y < 0 || x > len(rows[0])-1 || y > len(rows)-1 {
		return false, direction, pipe
	}
	nextPipe := rows[y][x]
	if nextPipe.Character == "." || nextPipe.Character == "S" {
		return false, direction, pipe
	}
	// fmt.Printf("Direction: %s\nPipe: %+v\nNextPipe: %+v\nX: %d\nY: %d\n", direction, pipe, nextPipe, x, y)

	nextDirection := direction

	directions := []string{
		"Right",
		"Left",
		"Top",
		"Down",
	}

	reverseDirection := getReverseDirection(direction)

	matchDirection := reflect.ValueOf(nextPipe).FieldByName(reverseDirection)

	if !matchDirection.Bool() {
		return false, nextDirection, nextPipe
	}

	v := reflect.ValueOf(nextPipe)
	typ := v.Type()

	for i := 0; i < v.NumField(); i++ {
		key := typ.Field(i).Name
		if slices.Index(directions, key) == -1 || key == reverseDirection {
			continue
		}
		value := v.Field(i)
		if value.Bool() {
			nextDirection = key
			break
		}
	}

	return matchDirection.Bool(), nextDirection, nextPipe
}

func getPipes(input []byte) [][]Pipe {
	ground := strings.Split(strings.TrimSpace(string(input)), "\n")

	rows := make([][]Pipe, len(ground))

	for y, row := range ground {
		for x, pipeRune := range row {
			pipe := string(pipeRune)
			switch pipe {
			case "|":
				rows[y] = append(rows[y], Pipe{
					Character: pipe,
					Right:     false,
					Left:      false,
					Top:       true,
					Down:      true,
					x:         x,
					y:         y,
				})
			case "-":
				rows[y] = append(rows[y], Pipe{
					Character: pipe,
					Right:     true,
					Left:      true,
					Top:       false,
					Down:      false,
					x:         x,
					y:         y,
				})
			case "J":
				rows[y] = append(rows[y], Pipe{
					Character: pipe,
					Right:     false,
					Left:      true,
					Top:       true,
					Down:      false,
					x:         x,
					y:         y,
				})
			case "L":
				rows[y] = append(rows[y], Pipe{
					Character: pipe,
					Right:     true,
					Left:      false,
					Top:       true,
					Down:      false,
					x:         x,
					y:         y,
				})
			case "7":
				rows[y] = append(rows[y], Pipe{
					Character: pipe,
					Right:     false,
					Left:      true,
					Top:       false,
					Down:      true,
					x:         x,
					y:         y,
				})
			case "F":
				rows[y] = append(rows[y], Pipe{
					Character: pipe,
					Right:     true,
					Left:      false,
					Top:       false,
					Down:      true,
					x:         x,
					y:         y,
				})
			case ".", "S":
				rows[y] = append(rows[y], Pipe{
					Character: pipe,
					Right:     false,
					Left:      false,
					Top:       false,
					Down:      false,
					x:         x,
					y:         y,
				})
			default:
				err := fmt.Sprintf("Unkown pipe error at position: [%d][%d]", y, x)
				panic(err)
			}
		}
	}

	return rows
}

func getReverseDirection(direction string) string {
	switch direction {
	case "Right":
		return "Left"
	case "Left":
		return "Right"
	case "Down":
		return "Top"
	case "Top":
		return "Down"
	default:
		panic("invalid direction")
	}
}

func getTruthySValue(directions map[string]map[string]int, startPipe Pipe, pipes [][]Pipe) string {
	var maybeS []string

	for direction, dValue := range directions {
		for coordinate, cValue := range dValue {
			x, y := startPipe.x, startPipe.y
			if coordinate == "x" {
				x += cValue
			} else {
				y += cValue
			}

			if y < 0 || y > len(pipes) {
				continue
			}

			if x < 0 || x > len(pipes[0]) {
				continue
			}

			reverseDirection := getReverseDirection(direction)
			nextPipe := pipes[y][x]
			matchDirection := reflect.ValueOf(nextPipe).FieldByName(reverseDirection)

			if matchDirection.Bool() {
				maybeS = append(maybeS, direction)
			}
		}
	}

	var truthyS string

	switch strings.Join(maybeS, "") {
	case "TopDown", "DownTop":
		truthyS = "|"
	case "LeftRight", "RightLeft":
		truthyS = "-"
	case "DownRight", "RightDown":
		truthyS = "F"
	case "DownLeft", "LeftDown":
		truthyS = "7"
	case "TopLeft", "LeftTop":
		truthyS = "J"
	case "TopRight", "RightTop":
		truthyS = "L"
	}

	return truthyS
}
