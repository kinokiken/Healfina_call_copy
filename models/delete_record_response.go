// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// DeleteRecordResponse delete record response
//
// swagger:model DeleteRecordResponse
type DeleteRecordResponse struct {

	// status
	Status string `json:"status,omitempty"`
}

// Validate validates this delete record response
func (m *DeleteRecordResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this delete record response based on context it is used
func (m *DeleteRecordResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *DeleteRecordResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DeleteRecordResponse) UnmarshalBinary(b []byte) error {
	var res DeleteRecordResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
