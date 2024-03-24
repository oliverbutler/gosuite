package design

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func GetBorderColor(selected bool) lipgloss.TerminalColor {
	if selected {
		return lipgloss.Color("50")
	}
	return lipgloss.Color("255")
}

// stripANSI removes ANSI escape sequences from a string.
func stripANSI(str string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(str, "")
}

func addTextToBorder(content string, index int, text string, selected bool) string {
	lines := strings.Split(content, "\n")
	if len(lines) < 1 {
		return content // If there's no content, just return it unchanged.
	}

	// Render the insertion text with optional bold styling.
	insertionText := lipgloss.NewStyle().
		Bold(selected).
		Foreground(GetBorderColor(selected)).
		Render("[" + strconv.Itoa(index) + "] " + text)

	// Calculate the visible length of the insertionText by stripping ANSI codes.
	visibleInsertionLength := len(stripANSI(insertionText))

	magicNumber := 6

	// Process the first line to find the third visible character position.
	strippedFirstLine := stripANSI(lines[0])
	// Ensuring not to exceed the length of the stripped line.
	if len(strippedFirstLine) < magicNumber {
		return content // Not enough length for insertion.
	}

	// The third visible character position in the stripped content.
	// Adjust 'insertAt' dynamically based on your requirements.
	insertAt := magicNumber + len([]rune(strippedFirstLine)[:magicNumber])

	// Convert the original first line (with ANSI codes) to runes.
	runes := []rune(lines[0])
	// Calculate the cut index; adjust dynamically if necessary.
	cutIndex := insertAt + visibleInsertionLength
	if cutIndex > len(runes) {
		cutIndex = len(runes) // Ensure cutIndex does not exceed the line length.
	}

	beforeInsertion := string(runes[:insertAt])
	afterInsertion := string(runes[cutIndex:])

	// Reassemble the modified top border.
	lines[0] = beforeInsertion + insertionText + afterInsertion

	return strings.Join(lines, "\n")
}

func CreatePane(
	index int,
	title string,
	selected bool,
	width int,
	height int,
	content string,
) string {
	styledContent := lipgloss.NewStyle().
		Padding(1, 1).
		Border(lipgloss.RoundedBorder()).
		MarginRight(1).
		Width(width).
		Height(height).BorderForeground(GetBorderColor(selected)).Render(content)

	return addTextToBorder(styledContent, index, title, selected)
}
