package repl

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    " Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    "Golang",
			expected: []string{"golang"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "     ",
			expected: []string{},
		},
		{
			input:    "  Go   es    GENIAL  ",
			expected: []string{"go", "es", "genial"},
		},
		{
			input:    "\tHola \n mundo\r",
			expected: []string{"hola", "mundo"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Not sufficent cases")
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Word and Expected doesn't match")
			}
		}
	}
}
