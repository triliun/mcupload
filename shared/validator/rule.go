package validator

import (
	"context"
	"net/http"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/jmoiron/sqlx"
)

type DBGetterInterface interface {
	GetDB(ctx context.Context) *sqlx.DB
}

type Rule struct {
	DBGetter DBGetterInterface
}

func NewRule(dbGetter DBGetterInterface) *Rule {
	return &Rule{DBGetter: dbGetter}
}

func (r *Rule) UsernameRules(isRequired bool, isUsernameExist bool, request *http.Request) []validation.Rule {
	rules := []validation.Rule{}

	if isRequired {
		rules = append(rules,
			validation.Required.Error("Username is required"),
			validation.Length(3, 30).Error("Username must be between 3-30 characters"),
			validation.Match(regexp.MustCompile(`^[a-z]`)).Error("Username must start with lowercase letter"),
			validation.Match(regexp.MustCompile(`^[a-z0-9_]*$`)).Error("Only lowercase letters, numbers, and underscores allowed"),
			validation.By(checkReservedUsername),
		)
	}

	if isUsernameExist && isRequired {
		fn := validation.By(func(value any) error {
			error := validation.NewError("validation_unique_username", "Username is already taken")

			if request == nil {
				return error
			}

			db := r.DBGetter.GetDB(request.Context())
			if db == nil {
				return error
			}

			var exists bool
			query := "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)"
			err := db.GetContext(request.Context(), &exists, query, value)
			if err != nil || exists {
				return error
			}

			return nil
		})

		rules = append(rules, fn)
	}

	return rules
}

func (r *Rule) EmailRules(isRequired bool, isEmailExist bool, request *http.Request) []validation.Rule {
	rules := []validation.Rule{}

	if isRequired {
		rules = append(rules,
			validation.Required.Error("Email is required"),
			is.EmailFormat.Error("Invalid email format"),
		)
	}

	if isEmailExist && isRequired {
		fn := validation.By(func(value any) error {
			error := validation.NewError("validation_unique_email", "Email is already taken")

			if request == nil {
				return error
			}

			db := r.DBGetter.GetDB(request.Context())
			if db == nil {
				return error
			}

			var exists bool
			query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
			err := db.GetContext(request.Context(), &exists, query, value)
			if err != nil || exists {
				return error
			}

			return nil
		})

		rules = append(rules, fn)
	}

	return rules
}

func (r *Rule) AgeRules(isRequired bool) []validation.Rule {
	rules := []validation.Rule{
		validation.When(isRequired,
			validation.Required.Error("Age is required"),
			validation.Min(13).Error("Must be at least 13 years old"),
			validation.Max(120).Error("Age seems invalid"),
		),
	}

	return rules
}

func (r *Rule) PasswordRules(isRequired bool) []validation.Rule {
	rules := []validation.Rule{
		validation.When(isRequired,
			validation.Required.Error("Password is required"),
			validation.Length(8, 0).Error("Password must be at least 8 characters"),
		),
	}

	return rules
}

func (r *Rule) UUIDRules(isRequired bool) []validation.Rule {
	rules := []validation.Rule{
		validation.When(isRequired,
			validation.Required.Error("ID is required"),
			validation.By(validateUUID),
		),
	}

	return rules
}

func (r *Rule) EditionRules(isRequired bool) []validation.Rule {
	rules := []validation.Rule{
		validation.When(isRequired,
			validation.Required.Error("Edition is required"),
			validation.In("Bedrock", "Java").Error("Edition must be Bedrock or Java"),
		),
	}

	return rules
}

func (r *Rule) VersionRules(isRequired bool) []validation.Rule {
	rules := []validation.Rule{
		validation.When(isRequired,
			validation.Required.Error("Version is required"),
			validation.By(validateVersion),
		),
	}

	return rules
}

func (r *Rule) SlugRules(isRequired bool, isSlugExist bool, request *http.Request) []validation.Rule {
	rules := []validation.Rule{}

	if isRequired {
		rules = append(rules, validation.Required.Error("Slug is required"),
			validation.Length(3, 50).Error("Slug must be between 3-50 characters"),
			validation.By(validateSlug))
	}

	if isSlugExist && isRequired {
		fn := validation.By(func(value any) error {
			error := validation.NewError("validation_unique_slug", "Slug is already taken")

			if request == nil {
				return error
			}

			db := r.DBGetter.GetDB(request.Context())
			if db == nil {
				return error
			}

			var exists bool
			query := "SELECT EXISTS(SELECT 1 FROM resource_packs WHERE slug = $1)"
			err := db.GetContext(request.Context(), &exists, query, value)
			if err != nil || exists {
				return error
			}

			return nil
		})

		rules = append(rules, fn)
	}

	return rules
}

func (r *Rule) TitleRules(isRequired bool) []validation.Rule {
	rules := []validation.Rule{
		validation.When(isRequired,
			validation.Required.Error("Title is required"),
			validation.Length(5, 100).Error("Title must be between 5-100 characters"),
			validation.By(validateTitle),
		),
	}

	return rules
}
