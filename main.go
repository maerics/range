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
	values, err := evaluateRange(expr, step)
	if err != nil {
		return err
	}
	for _, v := range values {
		fmt.Println(v)
	}
	return nil
}

func evaluateRange(expr string, step float64) ([]float64, error) {
	expr = strings.TrimSpace(expr)

	// Case 1: Single number (e.g., "10" -> 0-9)
	if matched, _ := regexp.MatchString(`^\d+$`, expr); matched {
		n, err := strconv.Atoi(expr)
		if err != nil {
			return nil, err
		}
		values := make([]float64, n)
		for i := 0; i < n; i++ {
			values[i] = float64(i)
		}
		return values, nil
	}

	// Case 2: Range with ".." (e.g., "2..4")
	if strings.Contains(expr, "..") {
		parts := strings.Split(expr, "..")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid range format: expected 'start..end'")
		}
		start, err := parseFloat(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid start value: %w", err)
		}
		end, err := parseFloat(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid end value: %w", err)
		}
		return generateRange(start, end, step, true, true), nil
	}

	// Case 3: Range with brackets (e.g., "[2.5,4.5]", "(2.5,4.5)", "[2.5,4.5)", etc.)
	// First trim outer quotes if present
	trimmedExpr := strings.Trim(expr, "'\"")
	// Check if it's a bracket notation after removing quotes
	if (strings.HasPrefix(trimmedExpr, "[") || strings.HasPrefix(trimmedExpr, "(")) &&
		(strings.HasSuffix(trimmedExpr, "]") || strings.HasSuffix(trimmedExpr, ")")) {
		// Determine inclusive/exclusive for start and end
		startInclusive := strings.HasPrefix(trimmedExpr, "[")
		endInclusive := strings.HasSuffix(trimmedExpr, "]")

		// Remove brackets and any remaining quotes
		trimmedExpr = strings.Trim(trimmedExpr, "[]()'\"")
		parts := strings.Split(trimmedExpr, ",")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid bracket range format: expected '[start,end]' or '(start,end)' or mixed")
		}
		start, err := parseFloat(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid start value: %w", err)
		}
		end, err := parseFloat(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, fmt.Errorf("invalid end value: %w", err)
		}
		return generateRange(start, end, step, startInclusive, endInclusive), nil
	}

	return nil, fmt.Errorf("unrecognized range expression format")
}

func parseFloat(s string) (float64, error) {
	s = strings.TrimSpace(s)
	return strconv.ParseFloat(s, 64)
}

func printRange(start, end, step float64) {
	values := generateRange(start, end, step, true, true)
	for _, v := range values {
		fmt.Println(v)
	}
}

func printRangeWithBounds(start, end, step float64, startInclusive, endInclusive bool) {
	values := generateRange(start, end, step, startInclusive, endInclusive)
	for _, v := range values {
		fmt.Println(v)
	}
}

func generateRange(start, end, step float64, startInclusive, endInclusive bool) []float64 {
	if step == 0 {
		step = 1.0
	}

	values := []float64{}

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
			values = append(values, val)
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
			values = append(values, val)
			val -= step
		}
	}

	return values
}
