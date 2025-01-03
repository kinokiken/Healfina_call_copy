// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// GenericResponse generic response
//
// swagger:model GenericResponse
type GenericResponse struct {

	// status
	Status string `json:"status,omitempty"`
}

// Validate validates this generic response
func (m *GenericResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this generic response based on context it is used
func (m *GenericResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *GenericResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *GenericResponse) UnmarshalBinary(b []byte) error {
	var res GenericResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
