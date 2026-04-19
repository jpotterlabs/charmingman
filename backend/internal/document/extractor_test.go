package document

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractText_PlainText(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "test.txt")
	err := os.WriteFile(filePath, []byte("Hello, World!"), 0644)
	assert.NoError(t, err)

	text, err := ExtractText(filePath)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World!", text)
}

func TestExtractText_Markdown(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "readme.md")
	err := os.WriteFile(filePath, []byte("# Hello\n\nWorld!"), 0644)
	assert.NoError(t, err)

	text, err := ExtractText(filePath)
	assert.NoError(t, err)
	assert.Equal(t, "# Hello\n\nWorld!", text)
}

func TestExtractText_UnsupportedExtension(t *testing.T) {
	text, err := ExtractText("test.xyz")
	assert.Error(t, err)
	assert.Empty(t, text)
	assert.Contains(t, err.Error(), "unsupported file extension")
}

func TestExtractText_FileNotFound(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "nonexistent.txt")

	text, err := ExtractText(filePath)
	assert.Error(t, err)
	assert.Empty(t, text)
	assert.True(t, os.IsNotExist(err))
}
