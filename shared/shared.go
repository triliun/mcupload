package shared

import (
	"github.com/triliun/mcupload/backend/shared/crypto"
	"github.com/triliun/mcupload/backend/shared/request/binder"
	"github.com/triliun/mcupload/backend/shared/request/context"
	"github.com/triliun/mcupload/backend/shared/request/parser"
	"github.com/triliun/mcupload/backend/shared/response/http"
	api "github.com/triliun/mcupload/backend/shared/types/http"
	"github.com/triliun/mcupload/backend/shared/types/sql"
	"github.com/triliun/mcupload/backend/shared/validator"
)

var (
	// crypto

	Crypto = &crypto.Crypto{}

	// request

	Binder  = binder.NewBinder()
	Context = &context.Context{}
	Param   = &parser.Param{}
	Query   = &parser.Query{}

	// response

	HTTP = http.HTTP{}

	// types DTO

	SQL = &sql.SQL{}

	// validator

	Validator = &validator.Validator{Rule: *validator.NewRule(Context)}
)

//  types

type (
	APIResponse           api.APIResponse
	InfiniteScrollRequest api.InfiniteScrollRequest
)
