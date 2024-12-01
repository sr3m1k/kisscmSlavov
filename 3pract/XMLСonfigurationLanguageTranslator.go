package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type ConfigTranslator struct {
	constants map[string]interface{}
}

func NewConfigTranslator() *ConfigTranslator {
	return &ConfigTranslator{
		constants: make(map[string]interface{}),
	}
}

func (ct *ConfigTranslator) translateXMLToConfig(input io.Reader) error {
	decoder := xml.NewDecoder(input)
	var result strings.Builder

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("XML parsing error: %v", err)
		}

		switch t := token.(type) {
		case xml.StartElement:
			result.WriteString(translateElement(t, decoder))
		}
	}

	fmt.Println(result.String())
	return nil
}

func translateElement(start xml.StartElement, decoder *xml.Decoder) string {
	var value string
	var _ error

	// Обработка текстового содержимого элемента
	next, _ := decoder.Token()
	if charData, ok := next.(xml.CharData); ok {
		value = string(charData)
	}

	switch start.Name.Local {
	case "number":
		return value
	case "string":
		return fmt.Sprintf("[[%s]]", value)
	case "array":
		return translateArray(start, decoder)
	case "constant":
		return translateConstant(start, value)
	case "compute":
		return translateCompute(value)
	default:
		return ""
	}
}

func translateArray(start xml.StartElement, decoder *xml.Decoder) string {
	var items []string

	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}

		switch t := token.(type) {
		case xml.EndElement:
			if t.Name.Local == "array" {
				return "'(" + strings.Join(items, " ") + ")"
			}
		case xml.StartElement:
			items = append(items, translateElement(t, decoder))
		}
	}

	return ""
}

func translateConstant(start xml.StartElement, value string) string {
	var name string
	for _, attr := range start.Attr {
		if attr.Name.Local == "name" {
			name = attr.Value
			break
		}
	}
	return fmt.Sprintf("%s -> %s", value, name)
}

func translateCompute(expr string) string {
	tokens := strings.Fields(expr)
	stack := []int{}

	for _, token := range tokens {
		switch token {
		case "+":
			if len(stack) < 2 {
				fmt.Fprintf(os.Stderr, "Invalid compute expression: %s\n", expr)
				return ""
			}
			b, a := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			stack = append(stack, a+b)
		case "ord()":
			if len(stack) < 1 {
				fmt.Fprintf(os.Stderr, "Invalid compute expression: %s\n", expr)
				return ""
			}
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			stack = append(stack, int(rune(top)))
		default:
			val, err := strconv.Atoi(token)
			if err == nil {
				stack = append(stack, val)
			}
		}
	}

	if len(stack) != 1 {
		fmt.Fprintf(os.Stderr, "Invalid compute expression: %s\n", expr)
		return ""
	}

	return fmt.Sprintf("$%d$", stack[0])
}

func isValidName(name string) bool {
	if len(name) == 0 {
		return false
	}

	// Первый символ - только буква в нижнем регистре
	if !unicode.IsLower(rune(name[0])) {
		return false
	}

	// Остальные символы - буквы, цифры или underscore
	for _, ch := range name[1:] {
		if !unicode.IsLower(ch) && !unicode.IsDigit(ch) && ch != '_' {
			return false
		}
	}

	return true
}

func main() {
	translator := NewConfigTranslator()

	err := translator.translateXMLToConfig(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Translation error: %v\n", err)
		os.Exit(1)
	}
}
