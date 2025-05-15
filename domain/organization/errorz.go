package organization

import "clean-architecture/pkg/errorz"

var (
	// ErrOrganizationNotFound is returned when an organization is not found
	ErrOrganizationNotFound = errorz.ErrNotFound.JoinError("organization not found")

	// ErrInvalidOrganizationData is returned when invalid data is provided
	ErrInvalidOrganizationData = errorz.ErrBadRequest.JoinError("invalid organization data")
)
