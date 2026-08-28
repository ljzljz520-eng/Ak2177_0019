package model

import "fmt"

func RecordID(warehouse, cycle string) string {
	return fmt.Sprintf("rec-%s-%s", normalizeKey(warehouse), normalizeKey(cycle))
}

func EventID(recordID string, sequence int) string {
	return fmt.Sprintf("evt-%s-%04d", normalizeKey(recordID), sequence)
}

func WorkflowID(recordID, name string) string {
	return fmt.Sprintf("wf-%s-%s", normalizeKey(recordID), normalizeKey(name))
}

func AttachmentID(recordID, name string) string {
	return fmt.Sprintf("att-%s-%s", normalizeKey(recordID), normalizeKey(name))
}

func normalizeKey(value string) string {
	result := make([]rune, 0, len(value))
	for _, ch := range value {
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') {
			result = append(result, ch)
		} else {
			result = append(result, '-')
		}
	}
	return string(result)
}
