package upload_pack

import (
	"mime/multipart"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/triliun/mcupload/backend/shared"
)

// Representaions table pack_files
type File struct {
	ID               uint64    `db:"id"`
	UserID           uint64    `db:"user_id"` // reference to id in table users
	PackID           uint32    `db:"pack_id"` // reference to id in table resource_pack
	FileName         string    `db:"filename"`
	Size             int64     `db:"size"`
	Extention        string    `db:"extention"`
	Path             string    `db:"path"`
	VersionFile      string    `db:"version_file"`
	VersionMinecraft uuid.UUID `db:"version_minecraft"` // reference to id in table version_minecraft
	Hash             string    `db:"hash"`
	CreatedAt        time.Time `db:"created_at"`
}

type CreateRequest struct {
	File             *multipart.FileHeader `form:"file" validate:"requrired,file"`
	VersionFile      string                `form:"version_file" validate:"requrired"`
	VersionMinecraft uuid.UUID             `form:"version_minecraft" validate:"requrired,number,version_minecraft"`
	Hash             string                `form:"hash_content" validate:"requrired,sha256"`
	userID           uint64
	packID           uint32
}

type CreateResponse struct {
	FileName  string    `json:"filename"`
	Hash      string    `json:"hash_content"`
	CreatedAt time.Time `json:"created_at"`
}

type GetResponse struct {
	FileName         string    `db:"filename" json:"filename"`
	Size             int64     `db:"size" json:"size"`
	Extention        string    `db:"extention" json:"extention"`
	Path             string    `db:"path" json:"path"`
	VersionFile      string    `db:"version_file" json:"version_file"`
	VersionMinecraft uuid.UUID `db:"version_minecraft" json:"version_minecraft"`
	Edition          string    `db:"edition" json:"edition"`
	Hash             string    `db:"hash" json:"hash_content"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
}

type GetMyFileResponse struct {
	ID               uint64    `db:"id" json:"id"`
	FileName         string    `db:"filename" json:"filename"`
	Size             int64     `db:"size" json:"size"`
	Extention        string    `db:"extention" json:"extention"`
	VersionFile      string    `db:"version_file" json:"version_file"`
	VersionMinecraft uuid.UUID `db:"version_minecraft" json:"version_minecraft"`
	Edition          string    `db:"edition" json:"edition"`
	Hash             string    `db:"hash" json:"hash_content"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
}

func (r CreateRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.File,
			validation.Required.Error("File is required"),
		),
		validation.Field(&r.VersionFile,
			shared.Validator.Rule.VersionRules(true)...,
		),
		validation.Field(&r.VersionMinecraft,
			validation.Required.Error("Minecraft version is required"),
			is.UUID.Error("Invalid id format"),
		),
		validation.Field(&r.Hash),
	)
}
