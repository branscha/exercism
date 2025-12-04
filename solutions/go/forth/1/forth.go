package forth

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var ErrForth = errors.New("forth error")

var intToken *regexp.Regexp
var wordToken *regexp.Regexp

func init() {
	intToken = regexp.MustCompile("-?[0-9]+")
	wordToken = regexp.MustCompile("[[:graph:]]*")
}

type forthStack []int

func (s *forthStack) len() int {
	return len(*s)
}

func (s *forthStack) peek() (int, error) {
	l := s.len()
	if l <= 0 {
		return 0, ErrForth
	}
	return (*s)[l-1], nil
}

func (s *forthStack) pop() (int, error) {
	val, err := s.peek()
	if err != nil {
		return 0, err
	}
	(*s) = (*s)[:s.len()-1]
	return val, nil
}

func (s *forthStack) pop2() (int, int, error) {
	fst, err := s.pop()
	if err != nil {
		return 0, 0, ErrForth
	}

	snd, err := s.pop()
	if err != nil {
		return 0, 0, ErrForth
	}
	return fst, snd, nil
}

func (s *forthStack) binop(op func(int, int) (int, error)) error {
	fst, snd, err := s.pop2()
	if err != nil {
		return err
	}
	val, err := op(fst, snd)
	if err != nil {
		return ErrForth
	}
	(*s) = append(*s, val)
	return nil
}

func (s *forthStack) plus() error {
	return s.binop(func(a int, b int) (int, error) {
		return a + b, nil
	})
}

func (s *forthStack) minus() error {
	return s.binop(func(a int, b int) (int, error) {
		return b - a, nil
	})
}

func (s *forthStack) mult() error {
	return s.binop(func(a int, b int) (int, error) {
		return a * b, nil
	})
}

func (s *forthStack) div() error {
	return s.binop(func(a int, b int) (int, error) {
		if a == 0 {
			return 0, ErrForth
		}
		return b / a, nil
	})
}

func (s *forthStack) dup() error {
	val, err := s.peek()
	if err != nil {
		return err
	}
	(*s) = append((*s), val)
	return nil
}

func (s *forthStack) drop() error {
	l := s.len()
	if l <= 0 {
		return ErrForth
	}
	(*s) = (*s)[:l-1]
	return nil
}

func (s *forthStack) swap() error {
	fst, snd, err := s.pop2()
	if err != nil {
		return err
	}
	(*s) = append((*s), fst, snd)
	return nil
}

func (s *forthStack) over() error {
	fst, snd, err := s.pop2()
	if err != nil {
		return err
	}
	(*s) = append((*s), snd, fst, snd)
	return nil
}

func (s *forthStack) pint(expr string) error {
	val, err := strconv.Atoi(expr)
	if err != nil {
		return err
	}
	(*s) = append((*s), val)
	return nil
}

type dictionary map[string][]string

func (d *dictionary) lookup(name string) ([]string, bool) {
	block, ok := (*d)[name]
	return block, ok
}

func (d *dictionary) add(name string, block []string) {
	var translated []string
	for _, expr := range block {
		l, ok := d.lookup(expr)
		if ok {
			translated = append(translated, l...)
		} else {
			translated = append(translated, expr)
		}
	}
	(*d)[name] = translated
}

func Forth(input []string) ([]int, error) {
	fmt.Printf("==> forth input: %v \n", input)
	stack := forthStack{}
	words := dictionary{}
	for _, line := range input {
		exprLst := strings.Fields(line)
		err := forth2(exprLst, &stack, &words)
		if err != nil {
			return nil, err
		}
	}
	return stack, nil
}

func forth2(input []string, stack *forthStack, words *dictionary) error {
loop:
	for i := 0; i < len(input); i++ {
		var err error
		expr := strings.ToUpper(input[i])

		// The word dictionary has precedence over built-in words.
		if block, ok := (*words)[expr]; ok {
			// Recursion step on definition body
			err = forth2(block, stack, words)
		} else {
			switch expr {
			case "+":
				err = stack.plus()
			case "*":
				err = stack.mult()
			case "/":
				err = stack.div()
			case "-":
				err = stack.minus()
			case "DUP":
				err = stack.dup()
			case "DROP":
				err = stack.drop()
			case "SWAP":
				err = stack.swap()
			case "OVER":
				err = stack.over()
			case ":":
				def := []string{}
				for i++; i < len(input); i++ {
					if input[i] == ";" {
						if len(def) <= 0 {
							err = ErrForth
						} else {
							if wordToken.MatchString(def[0]) && !intToken.MatchString(def[0]) {
								// Definition name has correct syntax, we can add the definition.
								words.add(def[0], def[1:])
								continue loop
							} else {
								// Definition name is not valid.
								err = ErrForth
								break
							}
						}
					} else {
						// Keep gobbling the definition.
						def = append(def, strings.ToUpper(input[i]))
					}
				} // definition lop
				err = ErrForth
			default:
				if intToken.MatchString(expr) {
					err = stack.pint(expr)
				} else {
					err = ErrForth
				}
			} // switch
		}

		if err != nil {
			return err
		}

	} // input loop
	return nil
}
