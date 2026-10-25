package field

import (
	"fmt"
	"math/rand"
	term "minesweep-go/src/util"
	"strconv"
	"strings"
)

type Direction int

const (
	LEFT  Direction = iota
	RIGHT Direction = iota
	UP    Direction = iota
	DOWN  Direction = iota
)

type Cursor struct {
	X, Y int
}

type Cell struct {
	X, Y    int
	IsMine  bool
	IsOpen  bool
	HasFlag bool
}

type Field struct {
	Initialized bool
	Width       int
	Height      int
	MineCount   int
	Cells       [][]*Cell
	Cur         *Cursor
}

func (field *Field) IsComplete() bool {
	for y := 0; y < field.Height; y++ {
		for x := 0; x < field.Width; x++ {
			cell := field.Cells[y][x]

			if !cell.IsMine && !cell.IsOpen {
				return false
			}
		}
	}

	return true
}

func (field *Field) InitMines() {
	for i := 0; i < field.MineCount; i++ {
		x := rand.Intn(field.Width)
		y := rand.Intn(field.Height)

		cell := field.Cells[y][x]
		if cell.IsMine || field.IsCellAtCursor(cell) {
			i--
			continue
		}

		field.Cells[y][x].IsMine = true
	}

	field.Initialized = true
}

func (field *Field) OpenCellAtCursor() *Cell {
	cell := field.Cells[field.Cur.Y][field.Cur.X]

	if cell.IsOpen || cell.HasFlag {
		return nil
	}

	field.OpenCell(field.Cells[field.Cur.Y][field.Cur.X])

	return cell
}

func (field *Field) GetFlagCount() int {
	count := 0

	for y := 0; y < field.Height; y++ {
		for x := 0; x < field.Width; x++ {
			if field.Cells[y][x].HasFlag {
				count++
			}
		}
	}

	return count
}

func (field *Field) PlaceFlagAtCursor() {
	cell := field.Cells[field.Cur.Y][field.Cur.X]

	if cell.IsOpen {
		return
	}

	cell.HasFlag = !cell.HasFlag
}

func (field *Field) MoveCursor(dir Direction) bool {
	switch dir {
	case LEFT:
		if field.Cur.X > 0 {
			field.Cur.X--
		}
	case RIGHT:
		if field.Cur.X < field.Width-1 {
			field.Cur.X++
		}
	case UP:
		if field.Cur.Y > 0 {
			field.Cur.Y--
		}
	case DOWN:
		if field.Cur.Y < field.Height-1 {
			field.Cur.Y++
		}
	}

	return true
}

func (field *Field) OpenAllMines() {
	for y := 0; y < field.Height; y++ {
		for x := 0; x < field.Width; x++ {
			if field.Cells[y][x].IsMine {
				field.Cells[y][x].IsOpen = true
			}
		}
	}
}

func (field *Field) IsCellAtCursor(cell *Cell) bool {
	return field.Cur.X == cell.X && field.Cur.Y == cell.Y
}

func (field *Field) OpenCell(cell *Cell) {
	if !field.Initialized {
		field.InitMines()
	}

	if cell.HasFlag {
		return
	}

	cell.IsOpen = true

	if cell.IsMine || field.CountAdjacentMines(cell) > 0 {
		return
	}

	adjacentCells := field.GetAdjacentCells(cell)

	for _, adjacentCell := range adjacentCells {
		if adjacentCell.IsOpen || adjacentCell.IsMine {
			continue
		}

		field.OpenCell(adjacentCell)
	}
}

func (field *Field) FindAdjacentEmptyCells(cell *Cell) []*Cell {
	emptyCells := []*Cell{}

	return emptyCells
}

func (field *Field) GetAdjacentCells(cell *Cell) []*Cell {
	adjacentCells := []*Cell{}

	offsets := [][]int{
		{-1, 0},  // left
		{1, 0},   // right
		{0, -1},  // top
		{0, 1},   // bottom
		{-1, -1}, // top left
		{1, -1},  // top right
		{-1, 1},  // bottom left
		{1, 1},   // bottom right
	}

	for _, offset := range offsets {
		x := cell.X + offset[0]
		y := cell.Y + offset[1]

		if x < 0 || x >= field.Width || y < 0 || y >= field.Height {
			continue // out of field bounds
		}

		adjacentCells = append(adjacentCells, field.Cells[y][x])
	}

	return adjacentCells
}

func (field *Field) Print() {

	/*
	   ╔ ═ ╗

	   ║   ║

	   ╚ ═ ╝


	*/

	fmt.Print("╔", strings.Repeat("═", field.Width*3), "╗\n")

	for y := 0; y < field.Height; y++ {

		fmt.Print("║")
		for x := 0; x < field.Width; x++ {
			cell := field.Cells[y][x]

			prefix := " "
			suffix := " "
			content := term.Colorize("■", term.DarkGray)
			// content := "."

			if field.IsCellAtCursor(cell) {
				prefix = term.Colorize("[", term.Cyan)
				suffix = term.Colorize("]", term.Cyan)
			}

			if cell.IsOpen && cell.IsMine {
				content = term.Colorize("@", term.Red)
			} else if cell.IsOpen {
				adjMineCount := field.CountAdjacentMines(cell)

				if adjMineCount == 0 {
					content = term.Colorize(" ", term.DarkGray)
				} else {
					content = colorizeAdjMines(adjMineCount)
				}
			} else if cell.HasFlag {
				content = term.Colorize("¶", term.Red)
			} else if field.IsCellAtCursor(cell) {
				content = term.Colorize(content, term.Cyan)
			}

			fmt.Printf("%s%s%s", prefix, content, suffix)
		}
		fmt.Print("║\n")
	}

	fmt.Print("╚", strings.Repeat("═", field.Width*3), "╝\n")
}

func (field *Field) CountAdjacentMines(cell *Cell) int {
	count := 0

	adjacentCells := field.GetAdjacentCells(cell)

	for _, adjacentCell := range adjacentCells {
		if adjacentCell.IsMine {
			count++
		}
	}

	return count
}

func NewField(width int, height int, minePercentage int) *Field {
	field := &Field{
		Initialized: false,
		Width:       width,
		Height:      height,
		MineCount:   width * height * minePercentage / 100,
		Cur:         &Cursor{X: 0, Y: 0},
	}

	field.Cells = make([][]*Cell, field.Height)

	for y := 0; y < field.Height; y++ {
		field.Cells[y] = make([]*Cell, field.Width)

		for x := 0; x < field.Width; x++ {
			field.Cells[y][x] = &Cell{X: x, Y: y, IsMine: false, IsOpen: false}
		}

	}

	return field
}

func colorizeAdjMines(adjMines int) string {
	countStr := strconv.Itoa(adjMines)

	if adjMines == 1 {
		return term.Colorize(countStr, term.Blue)
	}

	if adjMines == 2 {
		return term.Colorize(countStr, term.Yellow)
	}

	return term.Colorize(countStr, term.Orange)
}
