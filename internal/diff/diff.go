package diff

import (
	"fmt"
	"io"
	"log/slog"
	"reqcorder/internal/record"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// Generate a git-style diff between two resources (response, request, or template).
func DefaultDiff(w io.Writer, recordStorePath string, source string, target string, resource string) error {
	slog.Debug("Generating git-style diff", "recordStorePath", recordStorePath, "source", source, "target", target, "resource", resource)
	text1, text2, err := getTexts(recordStorePath, source, target, resource)
	if err != nil {
		slog.Error("Failed to get texts for diff", "error", err)
		return err
	}
	slog.Debug("Successfully retrieved texts for diff", "text1Length", len(text1), "text2Length", len(text2))
	err = gitStyleDiff(w, text1, text2, source, target)
	if err != nil {
		slog.Error("Failed to generate git-style diff", "error", err)
		return err
	}
	slog.Debug("Successfully generated git-style diff")
	return nil
}

// Generate an inline diff between two resources (response, request, or template).
func InlineDiff(w io.Writer, recordStorePath string, source string, target string, resource string) error {
	slog.Debug("Generating inline diff", "recordStorePath", recordStorePath, "source", source, "target", target, "resource", resource)
	text1, text2, err := getTexts(recordStorePath, source, target, resource)
	if err != nil {
		slog.Error("Failed to get texts for inline diff", "error", err)
		return err
	}
	slog.Debug("Successfully retrieved texts for inline diff", "text1Length", len(text1), "text2Length", len(text2))
	err = inlineDiff(w, text1, text2, source, target)
	if err != nil {
		slog.Error("Failed to generate inline diff", "error", err)
		return err
	}
	slog.Debug("Successfully generated inline diff")
	return nil
}

// Get response content by its ID for diff comparison.
func (d *DiffStore) getResponseByID(responseID string) (string, error) {
	slog.Debug("Getting response by ID for diff", "responseID", responseID, "recordStorePath", d.RecordStorePath)
	recordStore := &record.RecordStore{
		RecordStorePath: d.RecordStorePath,
		ResponseID:      responseID,
	}
	err := recordStore.GetResponseByID()
	if err != nil {
		slog.Error("Failed to get response by ID for diff", "error", err, "responseID", responseID)
		return "", err
	}
	slog.Debug("Successfully retrieved response for diff", "responseID", responseID)
	return string(recordStore.ResponseYaml), nil
}

// Get request content by its hash for diff comparison.
func (d *DiffStore) getRequestByHash(requestHash string) (string, error) {
	slog.Debug("Getting request by hash for diff", "requestHash", requestHash, "recordStorePath", d.RecordStorePath)
	recordStore := &record.RecordStore{
		RecordStorePath: d.RecordStorePath,
		RequestHash:     requestHash,
	}
	err := recordStore.GetRequestByHash()
	if err != nil {
		slog.Error("Failed to get request by hash for diff", "error", err, "requestHash", requestHash)
		return "", err
	}
	slog.Debug("Successfully retrieved request for diff", "requestHash", requestHash)
	return string(recordStore.RequestYaml), nil
}

// Get template content by its hash for diff comparison.
func (d *DiffStore) getTemplateByHash(templateHash string) (string, error) {
	slog.Debug("Getting template by hash for diff", "templateHash", templateHash, "recordStorePath", d.RecordStorePath)
	recordStore := &record.RecordStore{
		RecordStorePath: d.RecordStorePath,
		TemplateHash:    templateHash,
	}
	err := recordStore.GetTemplateByHash()
	if err != nil {
		slog.Error("Failed to get template by hash for diff", "error", err, "templateHash", templateHash)
		return "", err
	}
	slog.Debug("Successfully retrieved template for diff", "templateHash", templateHash)
	return string(recordStore.TemplateYaml), nil
}

// Get text content for two resources based on their type and identifiers.
func getTexts(recordStorePath string, source string, target string, resource string) (string, string, error) {
	slog.Debug("Getting texts for diff", "recordStorePath", recordStorePath, "source", source, "target", target, "resource", resource)
	diffStore := DiffStore{
		RecordStorePath: recordStorePath,
	}
	switch resource {
	case "response":
		slog.Debug("Getting response texts for diff", "sourceResponseID", source, "targetResponseID", target)
		text1, err := diffStore.getResponseByID(source)
		if err != nil {
			slog.Error("Failed to get source response for diff", "error", err, "responseID", source)
			return "", "", err
		}
		text2, err := diffStore.getResponseByID(target)
		if err != nil {
			slog.Error("Failed to get target response for diff", "error", err, "responseID", target)
			return "", "", err
		}
		slog.Debug("Successfully retrieved response texts for diff", "sourceLength", len(text1), "targetLength", len(text2))
		return text1, text2, nil
	case "request":
		slog.Debug("Getting request texts for diff", "sourceRequestHash", source, "targetRequestHash", target)
		text1, err := diffStore.getRequestByHash(source)
		if err != nil {
			slog.Error("Failed to get source request for diff", "error", err, "requestHash", source)
			return "", "", err
		}
		text2, err := diffStore.getRequestByHash(target)
		if err != nil {
			slog.Error("Failed to get target request for diff", "error", err, "requestHash", target)
			return "", "", err
		}
		slog.Debug("Successfully retrieved request texts for diff", "sourceLength", len(text1), "targetLength", len(text2))
		return text1, text2, nil
	case "template":
		slog.Debug("Getting template texts for diff", "sourceTemplateHash", source, "targetTemplateHash", target)
		text1, err := diffStore.getTemplateByHash(source)
		if err != nil {
			slog.Error("Failed to get source template for diff", "error", err, "templateHash", source)
			return "", "", err
		}
		text2, err := diffStore.getTemplateByHash(target)
		if err != nil {
			slog.Error("Failed to get target template for diff", "error", err, "templateHash", target)
			return "", "", err
		}
		slog.Debug("Successfully retrieved template texts for diff", "sourceLength", len(text1), "targetLength", len(text2))
		return text1, text2, nil
	default:
		slog.Error("Invalid diff type", "resource", resource)
		return "", "", fmt.Errorf("%w", ErrorInvalidDiffType)
	}
}

// Generate git-style diff output.
func gitStyleDiff(w io.Writer, text1 string, text2 string, filename1 string, filename2 string) error {
	lines1 := strings.Split(text1, "\n")
	lines2 := strings.Split(text2, "\n")
	dmp := diffmatchpatch.New()
	text1Lines := strings.Join(lines1, "\n")
	text2Lines := strings.Join(lines2, "\n")
	a, b, lineArray := dmp.DiffLinesToChars(text1Lines, text2Lines)
	diffs := dmp.DiffMain(a, b, false)
	diffs = dmp.DiffCharsToLines(diffs, lineArray)
	diffs = dmp.DiffCleanupSemantic(diffs)
	if _, err := fmt.Fprintf(w, "%s--- %s%s\n", ColorRed, filename1, ColorReset); err != nil {
		return fmt.Errorf("%w: %v", ErrorFailedToRenderDiff, err)
	}
	if _, err := fmt.Fprintf(w, "%s+++ %s%s\n\n", ColorGreen, filename2, ColorReset); err != nil {
		return fmt.Errorf("%w: %v", ErrorFailedToRenderDiff, err)
	}
	for _, diff := range diffs {
		lines := strings.Split(diff.Text, "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			switch diff.Type {
			case diffmatchpatch.DiffDelete:
				if _, err := fmt.Fprintf(w, "%s-%s%s\n", ColorRed, line, ColorReset); err != nil {
					return fmt.Errorf("%w: %v", ErrorFailedToRenderDiff, err)
				}
			case diffmatchpatch.DiffInsert:
				if _, err := fmt.Fprintf(w, "%s+%s%s\n", ColorGreen, line, ColorReset); err != nil {
					return fmt.Errorf("%w: %v", ErrorFailedToRenderDiff, err)
				}
			case diffmatchpatch.DiffEqual:
				if _, err := fmt.Fprintf(w, " %s\n", line); err != nil {
					return fmt.Errorf("%w: %v", ErrorFailedToRenderDiff, err)
				}
			}
		}
	}
	return nil
}

// Generate inline diff output with combined changes.
func inlineDiff(w io.Writer, text1 string, text2 string, filename1 string, filename2 string) error {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(text1, text2, false)
	dmp.DiffCleanupSemantic(diffs)
	if _, err := fmt.Fprintf(w, "%s%s%s\n", ColorGreen, filename2, ColorReset); err != nil {
		return fmt.Errorf("%w: %v", ErrorFailedToRenderDiff, err)
	}
	if _, err := fmt.Fprintf(w, "%s%s%s\n\n", ColorRed, filename1, ColorReset); err != nil {
		return fmt.Errorf("%w: %v", ErrorFailedToRenderDiff, err)
	}
	if _, err := fmt.Fprint(w, dmp.DiffPrettyText(diffs)); err != nil {
		return fmt.Errorf("%w: %v", ErrorFailedToRenderDiff, err)
	}
	return nil
}
