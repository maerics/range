# range

Print numbers to stdout from range expression arguments.

## Usage
```
Usage: range <ranges> ... [flags]

Generate number ranges from mathematical expressions

Arguments:
  <ranges> ...    Range expressions to evaluate

Flags:
  -h, --help        Show context-sensitive help.
  -s, --step=1.0    Step size for range iteration

range: error: expected "<ranges> ..."
```

## Examples
```
$ range 3
0
1
2
$ range 2..5
2
3
4
5
$ range '[3,6)' --step=0.5
3.0
3.5
4.0
4.5
5.0
5.5
```
