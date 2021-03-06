// Code generated by go-swagger; DO NOT EDIT.

package gen

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// DataRequestV1 DataRequest
//
// swagger:model dataRequestV1
type DataRequestV1 struct {

	// id
	ID string `json:"id,omitempty"`
}

// Validate validates this data request v1
func (m *DataRequestV1) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this data request v1 based on context it is used
func (m *DataRequestV1) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *DataRequestV1) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DataRequestV1) UnmarshalBinary(b []byte) error {
	var res DataRequestV1
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
