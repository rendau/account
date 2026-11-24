package core

import "strings"

// maskPhone hides a phone for logs: the first 4 and last 2 digits stay,
// 77011234567 → 7701*****67.
func maskPhone(phone string) string {
	const keepHead, keepTail = 4, 2

	if len(phone) <= keepHead+keepTail {
		return strings.Repeat("*", len(phone))
	}

	return phone[:keepHead] + strings.Repeat("*", len(phone)-keepHead-keepTail) + phone[len(phone)-keepTail:]
}
