package apperror

import "errors"

// ErrRateLimit is returned when a client exceeds the allowed request rate.
var ErrRateLimit = errors.New("rate limit reached")

// NotFoundError indicates a requested resource does not exist.
type NotFoundError struct{ Message string }

func (e NotFoundError) Error() string { return e.Message }

// ValidationError indicates invalid input from the client.
type ValidationError struct{ Message string }

func (e ValidationError) Error() string { return e.Message }

// ConflictError indicates a state conflict (e.g. duplicate key).
type ConflictError struct{ Message string }

func (e ConflictError) Error() string { return e.Message }

// UnauthorizedError indicates missing or invalid credentials.
type UnauthorizedError struct{ Message string }

func (e UnauthorizedError) Error() string { return e.Message }
