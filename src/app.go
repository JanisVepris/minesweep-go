package src

import (
	"fmt"
	"minesweep-go/src/field"
	term "minesweep-go/src/util"
	"os"
	"strings"
)

const (
	maxFieldWidth     = 30
	minFieldWidth     = 5
	maxFieldHeight    = 30
	minFieldHeight    = 5
	maxMinePercentage = 80
	minMinePercentage = 1
)

func Run() {
	width, _ := GetWidth()
	height, _ := GetHeight()
	minePercentage, _ := GetMinePercentage()

	gameField := field.NewField(width, height, minePercentage)

	GameLoop(gameField)
}

func GameLoop(gameField *field.Field) {
	// Disable input buffering
	if err := term.SetTerminalToRaw(); err != nil {
		fmt.Println("Error setting terminal to raw mode:", err)
		return
	}
	defer term.RestoreTerminal()

	end := false
	lost := false
	won := false

	PrintHeader(GetHeader(gameField))
	gameField.Print()

	for !end {
		// Read a single character
		char, _, err := term.GetCharStr()
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}

		cmd := strings.ToLower(char)

		switch cmd {
		case "q":
			end = true
		case "w":
			gameField.MoveCursor(field.UP)
		case "a":
			gameField.MoveCursor(field.LEFT)
		case "s":
			gameField.MoveCursor(field.DOWN)
		case "d":
			gameField.MoveCursor(field.RIGHT)
		case "f":
			gameField.PlaceFlagAtCursor()
		case " ":
			cell := gameField.OpenCellAtCursor()
			if cell != nil && cell.IsMine {
				gameField.OpenAllMines()

				lost = true
				end = true
			}

			if gameField.IsComplete() {
				won = true
				end = true
			}
		}

		header := GetHeader(gameField)
		term.Clear(gameField.Height+len(header)+2, gameField.Width*3+2)
		PrintHeader(header)
		gameField.Print()
	}

	if lost {
		fmt.Println(term.Colorize("\nYou lost :(", term.Red))
	}

	if won {
		fmt.Println(term.Colorize("\nYou WON :)!", term.Green))
	}
}

func PrintHeader(lines []string) {
	for _, line := range lines {
		fmt.Print(line)
	}
}

func GetHeader(gameField *field.Field) []string {
	return []string{
		fmt.Sprintf("============================\n"),
		fmt.Sprintf("Field Width: %d\n", gameField.Width),
		fmt.Sprintf("Field Height: %d\n", gameField.Height),
		fmt.Sprintf("Mine count: %d\n", gameField.MineCount),
		fmt.Sprintf("Flag count: %d\n", gameField.GetFlagCount()),

		fmt.Sprint("Controls:\n"),
		fmt.Sprint("  W    \t\tMove cursor up\n"),
		fmt.Sprint("  A    \t\tMove cursor left\n"),
		fmt.Sprint("  S    \t\tMove cursor down\n"),
		fmt.Sprint("  D    \t\tMove cursor right\n"),
		fmt.Sprint("  F    \t\tPlace a flag\n"),
		fmt.Sprint("  SPACE\t\tOpen cell\n"),
		fmt.Sprint("  Q    \t\tQuit\n"),
		fmt.Sprint("============================\n"),
		fmt.Sprint("\n"),
	}
}

func GetWidth() (width int, promptLen int) {
	prompt := fmt.Sprintf("Enter field width (%d-%d): ", minFieldWidth, maxFieldWidth)

	width = term.AskForInt(prompt)

	if width < minFieldWidth || width > maxFieldWidth {
		fmt.Printf("Invalid width: %d\n", width)
		os.Exit(0)
	}

	return width, len(prompt)
}

func GetHeight() (height int, promptLen int) {
	prompt := fmt.Sprintf("Enter field height (%d-%d): ", minFieldHeight, maxFieldHeight)

	height = term.AskForInt(prompt)

	if height < minFieldHeight || height > maxFieldHeight {
		fmt.Printf("Invalid height: %d\n", height)
		os.Exit(0)
	}

	return height, len(prompt)
}

func GetMinePercentage() (perc int, promptLen int) {
	prompt := fmt.Sprintf("Enter mine percentage (%d-%d): ", minMinePercentage, maxMinePercentage)

	perc = term.AskForInt(prompt)

	if perc < minMinePercentage || perc > maxMinePercentage {
		fmt.Printf("Invalid mine percentage: %d\n", perc)
		os.Exit(0)
	}

	return perc, len(prompt)
}
