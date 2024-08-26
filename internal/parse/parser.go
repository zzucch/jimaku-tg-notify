package parse

import (
	"fmt"
	"strconv"
	"unicode"
)

func parser(tokens []token) (Date, error) {
	iter := tokens

	expect := func(expected tokenType) error {
		if len(iter) == 0 || iter[0].Type != expected {
			return fmt.Errorf("expected %v, got %v", expected, iter[0].Type)
		}
		iter = iter[1:]
		return nil
	}

	getNumber := func() (int, error) {
		if len(iter) == 0 || iter[0].Type != Num {
			return 0, fmt.Errorf("expected number")
		}
		value := iter[0].Value
		iter = iter[1:]
		return value, nil
	}

	if err := expect(Start); err != nil {
		return Date{}, err
	}
	year, err := getNumber()
	if err != nil {
		return Date{}, err
	}
	if err := expect(Minus); err != nil {
		return Date{}, err
	}
	month, err := getNumber()
	if err != nil {
		return Date{}, err
	}
	if err := expect(Minus); err != nil {
		return Date{}, err
	}
	day, err := getNumber()
	if err != nil {
		return Date{}, err
	}
	hour, err := getNumber()
	if err != nil {
		return Date{}, err
	}
	if err := expect(Colon); err != nil {
		return Date{}, err
	}
	minute, err := getNumber()
	if err != nil {
		return Date{}, err
	}
	if err := expect(Colon); err != nil {
		return Date{}, err
	}
	second, err := getNumber()
	if err != nil {
		return Date{}, err
	}
	if err := offsetSignParser(&iter); err != nil {
		return Date{}, err
	}
	if _, err := getNumber(); err != nil {
		return Date{}, err
	}
	if err := expect(Colon); err != nil {
		return Date{}, err
	}
	if _, err := getNumber(); err != nil {
		return Date{}, err
	}

	return Date{Year: year, Month: month, Day: day, Hour: hour, Minute: minute, Second: second}, nil
}

func offsetSignParser(tokens *[]token) error {
	if len(*tokens) == 0 {
		return fmt.Errorf("expected + or -")
	}
	switch (*tokens)[0].Type {
	case Plus, Minus:
		*tokens = (*tokens)[1:]
		return nil
	default:
		return fmt.Errorf("expected + or -")
	}
}

func parseInt(input string) (int, int) {
	end := 0
	for end < len(input) && unicode.IsDigit(rune(input[end])) {
		end++
	}

	if end == 0 {
		return 0, 0
	}

	num, err := strconv.Atoi(input[:end])
	if err != nil {
		return 0, 0
	}

	return num, end
}
