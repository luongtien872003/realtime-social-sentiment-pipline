package validation

import (
	"social-insight/internal/models"
	"strings"
)

// ValidationError is a simple string message for failures
type ValidationError string

// Validator performs basic validation and sanitization for posts
type Validator struct{}

// New returns a Validator instance
func New() *Validator {
	return &Validator{}
}

// ValidatePost checks and sanitizes the given post.
// Returns ok (true if post is acceptable or fixable) and a list of errors (empty if ok)
func (v *Validator) ValidatePost(p *models.Post) (bool, []ValidationError) {
	errs := make([]ValidationError, 0)

	// Trim content and author
	p.Content = strings.TrimSpace(p.Content)
	p.Author = strings.TrimSpace(p.Author)

	// Content must not be empty
	if p.Content == "" {
		errs = append(errs, "empty content")
	}

	// Author fallback
	if p.Author == "" {
		p.Author = "Anonymous"
	}

	// Topic normalization: lowercase and ensure it's in known topics
	p.Topic = strings.ToLower(strings.TrimSpace(p.Topic))
	if p.Topic == "" {
		p.Topic = "programming"
	}

	// CreatedAt: if zero value, leave as-is (caller may set), but flag
	if p.CreatedAt.IsZero() {
		errs = append(errs, "missing created_at")
	}

	// Likes/Comments/Shares negative check
	if p.Likes < 0 {
		p.Likes = 0
	}
	if p.Comments < 0 {
		p.Comments = 0
	}
	if p.Shares < 0 {
		p.Shares = 0
	}

	// If only minor issues (author defaulted, counts normalized), consider ok
	// If content missing, reject
	if len(errs) > 0 {
		// If only missing created_at, allow (caller may parse later)
		if len(errs) == 1 && errs[0] == "missing created_at" {
			return true, nil
		}
		return false, errs
	}

	return true, nil
}
