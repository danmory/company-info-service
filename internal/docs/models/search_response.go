// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// SearchResponse search response
//
// swagger:model SearchResponse
type SearchResponse struct {

	// ceo
	Ceo string `json:"ceo,omitempty"`

	// inn
	Inn string `json:"inn,omitempty"`

	// name
	Name string `json:"name,omitempty"`
}

// Validate validates this search response
func (m *SearchResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this search response based on context it is used
func (m *SearchResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *SearchResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SearchResponse) UnmarshalBinary(b []byte) error {
	var res SearchResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
