package model

import "strings"

func CleanReason(reason string) string {
	return strings.TrimSpace(strings.Join(strings.Fields(reason), " "))
}

func ValidStatus(status RecordStatus) bool {
	switch status {
	case StatusDraft, StatusReview, StatusApproved, StatusArchived, StatusRejected:
		return true
	default:
		return false
	}
}

func NextStatus(current RecordStatus, action string) (RecordStatus, bool) {
	switch action {
	case "submit":
		if current == StatusDraft || current == StatusRejected {
			return StatusReview, true
		}
	case "approve":
		if current == StatusReview {
			return StatusApproved, true
		}
	case "reject":
		if current == StatusReview {
			return StatusRejected, true
		}
	case "archive":
		if current == StatusApproved {
			return StatusArchived, true
		}
	}
	return current, false
}
