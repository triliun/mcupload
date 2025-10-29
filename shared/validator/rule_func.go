package validator

import (
	"regexp"
	"slices"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func checkReservedUsername(value any) error {
	error := validation.NewError("validation_username_reserved", "Username not allowed")

	username, ok := value.(string)
	if !ok {
		return error
	}

	reserved := []string{"admin", "ceo", "root", "system", "user"}
	if slices.Contains(reserved, strings.ToLower(username)) {
		return error
	}

	return nil
}

func validateUUID(value any) error {
	_, ok := value.(uuid.UUID)
	if !ok {
		return validation.NewError("validation_uuid", "Invalid id format")
	}

	return nil
}

func validateVersion(value any) error {
	error := validation.NewError("validation_version", "Version format not supported")

	version, ok := value.(string)
	if !ok {
		return error
	}

	// // For release versions: 1.16.5, 1.19, 1.20.2
	// releaseVersionRegex := regexp.MustCompile(`^\d+\.\d+(\.\d+)?$`)

	// // For snapshot versions: 23w14a, 1.18-pre1, 20w45a
	// snapshotVersionRegex := regexp.MustCompile(`^(\d+w\d+[a-z]|\d+\.\d+(-pre\d+)?|[a-z]+\d+\.\d+(\.\d+)?)$`)

	// // Accept release version OR snapshot version
	// if releaseVersionRegex.MatchString(version) || snapshotVersionRegex.MatchString(version) {
	// 	return nil
	// }

	// Simplified and more accurate version validation
	versionRegex := regexp.MustCompile(`^(\d+\.\d+(\.\d+)?|\d+w\d+[a-z]|\d+\.\d+-pre\d+|[a-z]+\d+\.\d+(\.\d+)?)$`)

	if versionRegex.MatchString(version) {
		return nil
	}

	return error
}

func validateSlug(value any) error {
	error := validation.NewError("validation_slug", "Only lowercase letters, numbers, and hyphens allowed")

	slug, ok := value.(string)
	if !ok {
		return error
	}

	// Check for leading/trailing hyphens
	if strings.HasPrefix(slug, "-") || strings.HasSuffix(slug, "-") {
		return error
	}

	// Check for consecutive hyphens (optional)
	if strings.Contains(slug, "--") {
		return error
	}

	slugRegex := regexp.MustCompile(`^[a-z0-9-]+$`)

	if slugRegex.MatchString(slug) {
		return nil
	}

	return error
}

func validateTitle(value any) error {
	error := validation.NewError("validation_title", "Must start with capital letter and can contain letters, numbers, spaces, and basic punctuation")

	slug, ok := value.(string)
	if !ok {
		return error
	}
	titleRegex := regexp.MustCompile(`^[\p{L}\p{N}\s,.!?:;'"()&\[\]-]+$`)
	// titleRegex := regexp.MustCompile(`^[A-Z][A-Za-z0-9\s,.!?:;'\"()&\[\]-]*$`)

	if titleRegex.MatchString(slug) {
		return nil
	}

	return error
}
