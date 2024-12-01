package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestTranslateNumber(t *testing.T) {
	input := `<number>42</number>`
	expected := "42"
	result := testTranslate(input)
	if strings.TrimSpace(result) != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestTranslateString(t *testing.T) {
	input := `<string>Hello World</string>`
	expected := "[[Hello World]]"
	result := testTranslate(input)
	if strings.TrimSpace(result) != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestTranslateArray(t *testing.T) {
	input := `<array><number>1</number><number>2</number><number>3</number></array>`
	expected := "'(1 2 3)"
	result := testTranslate(input)
	if strings.TrimSpace(result) != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestTranslateConstant(t *testing.T) {
	input := `<constant name="x">42</constant>`
	expected := "42 -> x"
	result := testTranslate(input)
	if strings.TrimSpace(result) != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestTranslateCompute(t *testing.T) {
	input := `<compute>1 2 +</compute>`
	expected := "$3$"
	result := testTranslate(input)
	if strings.TrimSpace(result) != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestTranslateNestedStructures(t *testing.T) {
	input := `
	<array>
		<number>1</number>
		<array>
			<number>2</number>
			<number>3</number>
		</array>
		<string>test</string>
	</array>`
	expected := "'(1 '(2 3) [[test]])"
	result := testTranslate(input)
	if strings.TrimSpace(result) != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func testTranslate(input string) string {
	inputReader := strings.NewReader(input)
	var outputBuffer bytes.Buffer
	oldStdout := os.Stdout
	//os.Stdout = &outputBuffer

	translator := NewConfigTranslator()
	translator.translateXMLToConfig(inputReader)

	os.Stdout = oldStdout
	return outputBuffer.String()
}

func TestIsValidName(t *testing.T) {
	testCases := []struct {
		name     string
		expected bool
	}{
		{"x", true},
		{"abc123", true},
		{"x_y_z", true},
		{"X", false},
		{"1abc", false},
		{"abc!", false},
	}

	for _, tc := range testCases {
		result := isValidName(tc.name)
		if result != tc.expected {
			t.Errorf("Name '%s': expected %v, got %v", tc.name, tc.expected, result)
		}
	}
}
