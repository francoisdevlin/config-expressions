# Patterns Reference
The following patterns are available in a label, with the high precedence matches towards the top

1. An exact match
1. An enum match
1. A regex match
1. A wildcard match
1. A deep wildcard match

## Exact Match
An exact match is simply a label.  It looks like this

    a.b.c

## An Enum Match
An enum match is delimited by commas.  It looks like this

	a1,a2.b.c

This will match `a1.b.c` and `a2.b.c`

## A Regex Match
A regex match is surrounded by slashes.  Using the `.` is NOT supported.  Instead, `\w` is your friend  

	/a\w*/.b.c

This will match `a.b.c`, `apple.b.c`, and `aardvark.b.c` 

## Wildcard Match
A wildcard match is specified like this

	*.b.c

This will match `a.b.c`, `b.b.c`, and any three element path that ends with `b.c`

## Deep Wildcard Match
A deep wildcard match is specified like this

	**.c

This will match anything ending with `.c`, such as `a.b.c`, `b.c`, or even `c`.  The pattern is greedy, so it will also match `a.b.c.c`

# Learning More

There are collection of directories that include the actual source for this document,and they act as a test suite for the tool.  Each directory will contain:

* A config file, conf.jsonw
* A set of expected keys, in passing-keys.csv
* A human readable explaination of the test case, in input.md

You can play with these yourself, or add items to them.  Run the `tests.sh` command to regenerate this documentation.
