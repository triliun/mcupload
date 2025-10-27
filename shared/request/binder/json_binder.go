package binder

import (
	"io"
	"net/http"

	"github.com/bytedance/sonic"
)

type Binder struct {
	API sonic.API
}

func NewBinder() *Binder {
	api := sonic.Config{
		EscapeHTML: false,
		UseNumber:  false,
	}.Froze()

	return &Binder{API: api}
}

// BindJSON binds the request to provide struct using sonic
func (b *Binder) BindJSON(r *http.Request, object any) error {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = b.API.Unmarshal(body, object)

	return err
}

// BindWJSON binds the response to provide struct using sonic
func (b *Binder) BindWJSON(w *http.Response, object any) error {
	defer w.Body.Close()
	body, err := io.ReadAll(w.Body)
	if err != nil {
		return err
	}

	err = b.API.Unmarshal(body, object)

	return err
}
