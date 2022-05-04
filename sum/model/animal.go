// Package model is needed because easyjson cannot use the package main.
// See https://github.com/mailru/easyjson/issues/236
package model

// Animal was generated with easyjson `model/animal.go`
//easyjson:json
type Animal struct {
	Name  string
	Order string
}
