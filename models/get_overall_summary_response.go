// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// GetOverallSummaryResponse get overall summary response
//
// swagger:model GetOverallSummaryResponse
type GetOverallSummaryResponse struct {

	// overall summary
	OverallSummary string `json:"overall_summary,omitempty"`
}

// Validate validates this get overall summary response
func (m *GetOverallSummaryResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this get overall summary response based on context it is used
func (m *GetOverallSummaryResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *GetOverallSummaryResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *GetOverallSummaryResponse) UnmarshalBinary(b []byte) error {
	var res GetOverallSummaryResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
