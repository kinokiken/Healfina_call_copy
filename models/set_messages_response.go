// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// SetMessagesResponse set messages response
//
// swagger:model SetMessagesResponse
type SetMessagesResponse struct {

	// status
	Status string `json:"status,omitempty"`
}

// Validate validates this set messages response
func (m *SetMessagesResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this set messages response based on context it is used
func (m *SetMessagesResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *SetMessagesResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SetMessagesResponse) UnmarshalBinary(b []byte) error {
	var res SetMessagesResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}