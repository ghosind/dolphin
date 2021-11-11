package dolphin

import "errors"

// ErrNoTLSCert is returned by the TLS server when no certificate file is provided.
var ErrNoTLSCert = errors.New("no TLS certificate file")

// ErrNoTLSKey is returned by the TLS server when no private key file is provided.
var ErrNoTLSKey = errors.New("no TLS key file")

// ErrInvalidStatusCode is returned by the response status code is not valid.
var ErrInvalidStatusCode = errors.New("invalid status code")
