package documents

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecursiveCharacterTextSplitter_SplitText(t *testing.T) {
	splitter := &RecursiveCharacterTextSplitter{
		Separators:       []string{"\n\n", "\n", " "},
		KeepSeparator:    false,
		IsSeparatorRegex: false,
		ChunkSize:        10,
		OverlapSize:      0,
		LengthFunction:   func(s string) int { return len(s) },
	}

	text := "Hello\nWorld\n\nThis is a test\nof the text splitter."
	expected := []string{"Hello", "World", "This", "is", "a", "test", "of", "the", "text", "splitter."}

	result := splitter.SplitText(text)
	assert.Equal(t, expected, result, "The text should be split correctly without overlap.")
}

func TestFromLanguage(t *testing.T) {
	testCases := []struct {
		language Language
		wantErr  bool
	}{
		{PYTHON, false},
		{GO, false},
		{HTML, false},
		{JS, false},
		{TS, false},
		{MARKDOWN, false},
		{JSON, false},
		{Language("INVALID"), true},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(string(tc.language), func(t *testing.T) {
			splitter, err := FromLanguage(tc.language)
			if tc.wantErr {
				require.Error(t, err, "Expected an error for invalid language")
				assert.Nil(t, splitter, "Splitter should be nil for invalid language")
			} else {
				require.NoError(t, err, "Did not expect an error for valid language")
				assert.NotNil(t, splitter, "Splitter should not be nil for valid language")
			}
		})
	}
}

func TestGetSeparatorsForLanguage(t *testing.T) {
	testCases := []struct {
		language Language
		wantErr  bool
	}{
		{PYTHON, false},
		{GO, false},
		{HTML, false},
		{JS, false},
		{TS, false},
		{MARKDOWN, false},
		{JSON, false},
		{Language("INVALID"), true},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(string(tc.language), func(t *testing.T) {
			separators, err := GetSeparatorsForLanguage(tc.language)
			if tc.wantErr {
				require.Error(t, err, "Expected an error for invalid language")
				assert.Nil(t, separators, "Separators should be nil for invalid language")
			} else {
				require.NoError(t, err, "Did not expect an error for valid language")
				assert.NotEmpty(t, separators, "Separators should not be empty for valid language")
			}
		})
	}
}

func TestSplitTextByCount(t *testing.T) {
	text := "Hello World! This is a test of the text splitter."
	expected := []string{"Hello Wo", "rld! Thi", "s is a t", "est of t", "he text ", "splitter", "."}

	result := SplitTextByCount(text, 8)
	assert.Equal(t, expected, result, "The text should be split correctly by count.")
}
