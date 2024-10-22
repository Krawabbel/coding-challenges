package calculator

import (
	"fmt"
	"math"
	"strconv"
)

const (
	TOKEN_NUMBER = iota
	TOKEN_PLUS
	TOKEN_MINUS
	TOKEN_TIMES
	TOKEN_SLASH
	TOKEN_PAREN_OPEN
	TOKEN_PAREN_CLOSED
	TOKEN_CIRCUMFLEX
)

func Calculate(expr string) (float64, error) {

	tokens, err := tokenize(expr)
	if err != nil {
		return 0, err
	}

	result, err := evalSum(&tokens)
	if err != nil {
		return 0, err
	}

	if len(tokens) > 0 {
		return 0, fmt.Errorf("ill-formed expression")
	}

	return result, nil
}

type token struct {
	spec   int
	number float64
}

func isNumeric(b byte) bool {
	return (b >= '0' && b <= '9') || (b == '.')
}

func tokenize(expr string) ([]token, error) {

	ts := make([]token, 0)

	for ptr := 0; ptr < len(expr); ptr++ {
		b := expr[ptr]
		switch {
		case b == '+':
			ts = append(ts, token{spec: TOKEN_PLUS})
		case b == '-':
			ts = append(ts, token{spec: TOKEN_MINUS})
		case b == '*' || b == 'x':
			ts = append(ts, token{spec: TOKEN_TIMES})
		case b == '/':
			ts = append(ts, token{spec: TOKEN_SLASH})
		case b == '(':
			ts = append(ts, token{spec: TOKEN_PAREN_OPEN})
		case b == ')':
			ts = append(ts, token{spec: TOKEN_PAREN_CLOSED})
		case b == '^':
			ts = append(ts, token{spec: TOKEN_CIRCUMFLEX})
		case isNumeric(b):
			end := ptr + 1
			for end < len(expr) {
				if !isNumeric(expr[end]) {
					break
				}
				end++
			}
			literal := expr[ptr:end]
			number, err := strconv.ParseFloat(literal, 64)
			if err != nil {
				return nil, err
			}
			ts = append(ts, token{spec: TOKEN_NUMBER, number: number})
			ptr = end - 1
		case b == ' ':
			// ignore whitespace
		default:
			return nil, fmt.Errorf("ill-formed expression: unexpected character '%v'", b)
		}
	}

	return ts, nil
}

func match(ts *[]token, specs ...int) (token, bool) {
	if len(*ts) == 0 {
		return token{}, false
	}

	t := (*ts)[0]
	for _, spec := range specs {
		if t.spec == spec {
			*ts = (*ts)[1:]
			return t, true
		}
	}

	return token{}, false
}

func evalExpr(ts *[]token) (float64, error) {
	return evalSum(ts)
}

func evalSum(ts *[]token) (float64, error) {

	result, err := evalTerm(ts)
	if err != nil {
		return 0, err
	}

	for {
		op, matched := match(ts, TOKEN_PLUS, TOKEN_MINUS)
		if !matched {
			return result, nil
		}

		right, err := evalTerm(ts)
		if err != nil {
			return 0, err
		}

		switch op.spec {
		case TOKEN_PLUS:
			result += right
		case TOKEN_MINUS:
			result -= right
		default:
			return 0, fmt.Errorf("unreachable")
		}

	}
}

func evalTerm(ts *[]token) (float64, error) {

	result, err := evalPower(ts)
	if err != nil {
		return 0, err
	}

	for {
		op, matched := match(ts, TOKEN_TIMES, TOKEN_SLASH)

		if !matched {
			return result, nil
		}

		right, err := evalPower(ts)
		if err != nil {
			return 0, err
		}

		switch op.spec {
		case TOKEN_TIMES:
			result *= right
		case TOKEN_SLASH:
			result /= right
		default:
			return 0, fmt.Errorf("unreachable")
		}
	}
}

func evalPower(ts *[]token) (float64, error) {

	result, err := evalGroup(ts)
	if err != nil {
		return 0, err
	}

	for {
		_, matched := match(ts, TOKEN_CIRCUMFLEX)

		if !matched {
			return result, nil
		}

		right, err := evalGroup(ts)
		if err != nil {
			return 0, err
		}

		result = math.Pow(result, right)
	}
}

func evalGroup(ts *[]token) (float64, error) {

	if _, matched := match(ts, TOKEN_PAREN_OPEN); !matched {
		return evalUnary(ts)
	}

	expr, err := evalExpr(ts)
	if err != nil {
		return 0, err
	}

	if _, matched := match(ts, TOKEN_PAREN_CLOSED); !matched {
		return 0, fmt.Errorf("ill-formed expression: missing closing parenthesis")
	}

	return expr, nil
}

func evalUnary(ts *[]token) (float64, error) {

	unary, matched := match(ts, TOKEN_PLUS, TOKEN_MINUS)

	right, err := evalNumber(ts)
	if err != nil {
		return 0, err
	}

	if matched {
		switch unary.spec {
		case TOKEN_MINUS:
			return -right, nil
		case TOKEN_PLUS:
			return right, nil
		default:
		}
		return 0, fmt.Errorf("unreachable")
	}

	return right, nil
}

func evalNumber(ts *[]token) (float64, error) {
	t, matched := match(ts, TOKEN_NUMBER)
	if !matched {
		return 0, fmt.Errorf("ill-formed expression: missing number")
	}
	return t.number, nil
}
