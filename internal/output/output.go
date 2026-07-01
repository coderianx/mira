package output

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

type Style int

const (
	ResetStyle Style = 0
	BoldStyle  Style = 1
	DimStyle   Style = 2
)

type Color int

const (
	Black   Color = 30
	Red     Color = 31
	Green   Color = 32
	Yellow  Color = 33
	Blue    Color = 34
	Magenta Color = 35
	Cyan    Color = 36
	White   Color = 37
)

var w io.Writer = os.Stdout

func SetOutput(out io.Writer) {
	w = out
}

func ansi(color Color, style ...Style) string {
	s := make([]string, 0, 2)
	for _, st := range style {
		s = append(s, fmt.Sprintf("%d", st))
	}
	s = append(s, fmt.Sprintf("%d", color))
	return fmt.Sprintf("\033[%sm", strings.Join(s, ";"))
}

const reset = "\033[0m"

func cs(s string, color Color, style ...Style) string {
	return ansi(color, style...) + s + reset
}

func Title(s string) {
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, cs(" "+s+" ", Cyan, BoldStyle))
}

func Section(s string) {
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, cs(s, Cyan, BoldStyle))
}

func Subtitle(s string) {
	fmt.Fprintln(w, "  "+cs(s, White, BoldStyle))
}

func Info(msg string) {
	fmt.Fprintln(w, cs("  ◇ ", Blue)+msg)
}

func Infof(format string, args ...any) {
	Info(fmt.Sprintf(format, args...))
}

func Success(msg string) {
	fmt.Fprintln(w, cs("  ✓ ", Green)+msg)
}

func Successf(format string, args ...any) {
	Success(fmt.Sprintf(format, args...))
}

func Warning(msg string) {
	fmt.Fprintln(w, cs("  ⚠ ", Yellow)+msg)
}

func Warningf(format string, args ...any) {
	Warning(fmt.Sprintf(format, args...))
}

func Error(msg string) {
	fmt.Fprintln(w, cs("  ✗ ", Red)+msg)
}

func Errorf(format string, args ...any) {
	Error(fmt.Sprintf(format, args...))
}

func Dim(format string, args ...any) {
	s := fmt.Sprintf(format, args...)
	fmt.Fprintln(w, cs("  "+s, White, DimStyle))
}

func KeyValue(key, value string) {
	fmt.Fprintf(w, "  %s  %s\n", cs(pad(key, 16)+":", White, BoldStyle), value)
}

func KeyValuef(key, format string, args ...any) {
	KeyValue(key, fmt.Sprintf(format, args...))
}

func pad(s string, n int) string {
	l := utf8.RuneCountInString(s)
	if l >= n {
		return s
	}
	return s + strings.Repeat(" ", n-l)
}

func Badge(text string, color Color) string {
	return ansi(color, BoldStyle) + " " + text + " " + reset
}

func Step(current, total int, msg string) {
	fmt.Fprintf(w, "  %s  %s\n",
		cs(fmt.Sprintf("[%d/%d]", current, total), Cyan, BoldStyle),
		msg,
	)
}

func Header(repo string, version string) {
	const width = 60
	title := " mira "
	if version != "" {
		title += version + " "
	}
	if repo != "" {
		title += "· " + repo + " "
	}

	fmt.Fprintln(w, "")
	fmt.Fprintln(w, cs("┌"+strings.Repeat("─", width-2)+"┐", Cyan, DimStyle))
	fmt.Fprintln(w, cs("│", Cyan, DimStyle)+cs(padCenter(title, width-2), Cyan, BoldStyle)+cs("│", Cyan, DimStyle))
	fmt.Fprintln(w, cs("└"+strings.Repeat("─", width-2)+"┘", Cyan, DimStyle))
	fmt.Fprintln(w, "")
}

func padCenter(s string, n int) string {
	l := utf8.RuneCountInString(s)
	if l >= n {
		return s
	}
	left := (n - l) / 2
	right := n - l - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}

func Box(text string, color Color) {
	const width = 60
	lines := strings.Split(text, "\n")
	fmt.Fprintln(w, cs("┌"+strings.Repeat("─", width-2)+"┐", color, DimStyle))
	for _, line := range lines {
		fmt.Fprintln(w, cs("│", color, DimStyle)+padCenter(line, width-2)+cs("│", color, DimStyle))
	}
	fmt.Fprintln(w, cs("└"+strings.Repeat("─", width-2)+"┘", color, DimStyle))
}

func Table(headers []string, rows [][]string) {
	if len(rows) == 0 {
		return
	}

	cols := len(headers)
	widths := make([]int, cols)
	for i, h := range headers {
		widths[i] = utf8.RuneCountInString(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < cols {
				if l := utf8.RuneCountInString(cell); l > widths[i] {
					widths[i] = l
				}
			}
		}
	}

	sep := "─"
	fmt.Fprint(w, cs("  ┌", Cyan, DimStyle))
	for i, wd := range widths {
		if i > 0 {
			fmt.Fprint(w, cs("┬", Cyan, DimStyle))
		}
		fmt.Fprint(w, cs(strings.Repeat(sep, wd+2), Cyan, DimStyle))
	}
	fmt.Fprintln(w, cs("┐", Cyan, DimStyle))

	fmt.Fprint(w, cs("  │", Cyan, DimStyle))
	for i, h := range headers {
		p := widths[i] - utf8.RuneCountInString(h)
		fmt.Fprintf(w, " %s %s", h+strings.Repeat(" ", p), cs("│", Cyan, DimStyle))
	}
	fmt.Fprintln(w)

	fmt.Fprint(w, cs("  ├", Cyan, DimStyle))
	for i, wd := range widths {
		if i > 0 {
			fmt.Fprint(w, cs("┼", Cyan, DimStyle))
		}
		fmt.Fprint(w, cs(strings.Repeat(sep, wd+2), Cyan, DimStyle))
	}
	fmt.Fprintln(w, cs("┤", Cyan, DimStyle))

	for _, row := range rows {
		fmt.Fprint(w, cs("  │", Cyan, DimStyle))
		for i, cell := range row {
			if i >= cols {
				break
			}
			p := widths[i] - utf8.RuneCountInString(cell)
			fmt.Fprintf(w, " %s %s", cell+strings.Repeat(" ", p), cs("│", Cyan, DimStyle))
		}
		fmt.Fprintln(w)
	}

	fmt.Fprint(w, cs("  └", Cyan, DimStyle))
	for i, wd := range widths {
		if i > 0 {
			fmt.Fprint(w, cs("┴", Cyan, DimStyle))
		}
		fmt.Fprint(w, cs(strings.Repeat(sep, wd+2), Cyan, DimStyle))
	}
	fmt.Fprintln(w, cs("┘", Cyan, DimStyle))
}

func Fatal(err error) {
	fmt.Fprintln(w, "")
	fmt.Fprintf(w, "  %s %s\n", cs("✗", Red, BoldStyle), cs(err.Error(), Red))
	fmt.Fprintln(w, "")
	os.Exit(1)
}
