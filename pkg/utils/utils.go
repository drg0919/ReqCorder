package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"

	"github.com/goccy/go-yaml"
)

// Create a preview of content for displaying on the terminal.
func CreatePreview(body string) string {
	if len(body) > 80 {
		return body[:76] + "..."
	}
	return body
}

// Convert any object to its JSON representation.
func Prettify(v any) (string, error) {
	temp, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrorFailedToMarshalJSON, err)
	}
	result := string(temp)
	if quoted, err := strconv.Unquote(result); err == nil {
		return quoted, nil
	}
	return result, nil
}

// Convert a byte slice to YAML content using the `github.com/goccy/go-yaml` library.
func ConvertToYAML(r any) ([]byte, error) {
	obj, err := yaml.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrorFailedToMarshalYAML, err)
	}
	return obj, nil
}

// Read a YAML file into an object using the `github.com/goccy/go-yaml` library.
func ReadYAMLFile(file string, sourceObject any) error {
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("%w %q: %v", ErrorFailedToReadFile, file, err)
	}
	err = yaml.Unmarshal([]byte(yamlFile), sourceObject)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrorFailedToUnmarshalYAML, err)
	}
	return nil
}

// Read a file and return its content.
func ReadFile(file string) ([]byte, error) {
	fileContent, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("%w %q: %v", ErrorFailedToReadFile, file, err)
	}
	return fileContent, nil
}

// Calculate the MD5 hash of the given content.
func CalculateMD5Hash(content []byte) string {
	sum := md5.Sum(content)
	hash := hex.EncodeToString(sum[:])
	return hash
}

// Calculate the MD5 hash from the file's content.
func CalculateMD5HashFromFile(file string) (string, error) {
	fileContent, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("%w %q: %v", ErrorFailedToReadFile, file, err)
	}
	return CalculateMD5Hash(fileContent), nil
}

// Ensure the specified path exists.
func EnsureDir(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("%w %q: %v", ErrorFailedToCreateDirectory, path, err)
	}
	return nil
}

// Print the given error to the specified writer, prefixing with "error: ". Panics on write failure.
func PrintError(w io.Writer, err error) {
	_, e := fmt.Fprintf(w, "error: %v\n", err)
	if e != nil {
		panic(fmt.Errorf("failed to write to writer: %w", e))
	}
}

// Print to the specified writer using fmt.Fprint. Panics on write failure.
func Fprint(w io.Writer, a ...any) {
	_, err := fmt.Fprint(w, a...)
	if err != nil {
		panic(fmt.Errorf("failed to write to writer: %w", err))
	}
}

// Print to the specified writer using fmt.Fprintf with the given format. Panics on write failure.
func Fprintf(w io.Writer, format string, a ...any) {
	_, err := fmt.Fprintf(w, format, a...)
	if err != nil {
		panic(fmt.Errorf("failed to write to writer: %w", err))
	}
}

// Print to the specified writer using fmt.Fprintln, appending a newline. Panics on write failure.
func Fprintln(w io.Writer, a ...any) {
	_, err := fmt.Fprintln(w, a...)
	if err != nil {
		panic(fmt.Errorf("failed to write to writer: %w", err))
	}
}

// Helper function to log error types.
func LogError(err error) slog.Value {
	return slog.GroupValue(
		slog.String("error", err.Error()),
	)
}
