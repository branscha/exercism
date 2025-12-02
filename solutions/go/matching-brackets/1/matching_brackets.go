package brackets

import (
    "errors"
)

const (
    PARENS int = iota
    BRACKETS 
    BRACES
)

type stack []int 

func newStack() stack {
    return []int{}
}

var ErrStack = errors.New("Stack eror")

func (s *stack) size() int {
    return len(*s)
}

func (s *stack) push(v int) {
    *s = append(*s, v)
}

func (s *stack) peek() (int, error){
    if s.size() <= 0 {
        return 0, ErrStack
    }
    return (*s)[len(*s)-1], nil
}

func (s *stack) pop() (int, error) {
    val, err := s.peek()
    if err != nil {
        return 0, err
    }
    *s = (*s)[:len(*s)-1]
    return val, nil
}

func Bracket(input string) bool {
    s := newStack()
   for _, c := range input {
       switch c {
           case '(':
        	   s.push(PARENS)
           case '[':
              s.push(BRACKETS)
           case '{':
              s.push(BRACES)
           case ')':
              if err := match(&s, PARENS); err != nil {
                  return false
              }
           case '}':
              if err := match(&s, BRACES); err != nil {
                  return false
              }
           case ']': 
              if err := match(&s, BRACKETS); err != nil {
                  return false
              } 
       }
   }
   return s.size() == 0
}

func match(s *stack, t int) error {
    v, err := s.peek()
    if err != nil {
        return err
    }
    if v != t {
        return ErrStack
    } 
    s.pop()
    return nil
}
