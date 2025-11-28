package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/alecthomas/kong"
)

type CLI struct {
	Step   float64  `short:"s" help:"Step size for range iteration" default:"1.0"`
	Ranges []string `arg:"" help:"Range expressions to evaluate"`
}

func main() {
	var cli CLI
	ctx := kong.Parse(&cli,
		kong.Description("Generate number ranges from mathematical expressions"),
		kong.UsageOnError(),
	)

	ctx.FatalIfErrorf(processRanges(&cli))
}

func processRanges(cli *CLI) error {
	for _, expr := range cli.Ranges {
		if err := processRange(expr, cli.Step); err != nil {
			return fmt.Errorf("error processing range %q: %w", expr, err)
		}
	}
	return nil
}

func processRange(expr string, step float64) error {
	expr = strings.TrimSpace(expr)

	// Case 1: Single number (e.g., "10" -> 0-9)
	if matched, _ := regexp.MatchString(`^\d+$`, expr); matched {
		n, err := strconv.Atoi(expr)
		if err != nil {
			return err
		}
		for i := 0; i < n; i++ {
			fmt.Println(i)
		}
		return nil
	}

	// Case 2: Range with ".." (e.g., "2..4")
	if strings.Contains(expr, "..") {
		parts := strings.Split(expr, "..")
		if len(parts) != 2 {
			return fmt.Errorf("invalid range format: expected 'start..end'")
		}
		start, err := parseFloat(parts[0])
		if err != nil {
			return fmt.Errorf("invalid start value: %w", err)
		}
		end, err := parseFloat(parts[1])
		if err != nil {
			return fmt.Errorf("invalid end value: %w", err)
		}
		printRange(start, end, step)
		return nil
	}

	// Case 3: Range with brackets (e.g., "[2.5,4.5]", "(2.5,4.5)", "[2.5,4.5)", etc.)
	if (strings.HasPrefix(expr, "[") || strings.HasPrefix(expr, "(")) &&
		(strings.HasSuffix(expr, "]") || strings.HasSuffix(expr, ")")) {
		// Determine inclusive/exclusive for start and end
		startInclusive := strings.HasPrefix(expr, "[")
		endInclusive := strings.HasSuffix(expr, "]")

		// Remove brackets and quotes if present
		expr = strings.Trim(expr, "[]()'\"")
		parts := strings.Split(expr, ",")
		if len(parts) != 2 {
			return fmt.Errorf("invalid bracket range format: expected '[start,end]' or '(start,end)' or mixed")
		}
		start, err := parseFloat(strings.TrimSpace(parts[0]))
		if err != nil {
			return fmt.Errorf("invalid start value: %w", err)
		}
		end, err := parseFloat(strings.TrimSpace(parts[1]))
		if err != nil {
			return fmt.Errorf("invalid end value: %w", err)
		}
		printRangeWithBounds(start, end, step, startInclusive, endInclusive)
		return nil
	}

	return fmt.Errorf("unrecognized range expression format")
}

func parseFloat(s string) (float64, error) {
	s = strings.TrimSpace(s)
	return strconv.ParseFloat(s, 64)
}

func printRange(start, end, step float64) {
	printRangeWithBounds(start, end, step, true, true)
}

func printRangeWithBounds(start, end, step float64, startInclusive, endInclusive bool) {
	if step == 0 {
		step = 1.0
	}

	if start <= end {
		// Ascending range
		val := start
		if !startInclusive {
			val += step
		}
		for {
			if endInclusive {
				if val > end {
					break
				}
			} else {
				if val >= end {
					break
				}
			}
			fmt.Println(val)
			val += step
		}
	} else {
		// Descending range
		val := start
		if !startInclusive {
			val -= step
		}
		for {
			if endInclusive {
				if val < end {
					break
				}
			} else {
				if val <= end {
					break
				}
			}
			fmt.Println(val)
			val -= step
		}
	}
}
