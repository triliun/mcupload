package validator

import "fmt"

// User-friendly validation error messages
func messageForTag(tag string, param string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "number":
		return "Invalid number format"
	case "uuid":
		return "Invalid ID format"
	case "min":
		return fmt.Sprintf("Must be at least %s characters", param)
	case "max":
		return fmt.Sprintf("Cannot exceed %s characters", param)
	case "eqfield":
		return fmt.Sprintf("Must match %s field", param)
	case "username":
		return "Only letters, numbers, and underscores allowed"
	case "unique_username":
		return "Username is already taken"
	case "unique_email":
		return "Email is already registered"
	case "slug":
		return "Only lowercase letters, numbers, and hyphens allowed"
	case "title":
		return "Must start with capital letter and can contain letters, numbers, spaces, and basic punctuation"
	case "minecraft_version":
		return "Invalid version format"
	default:
		return "Invalid value"
	}
}
