package dolphin

// HTTP Request methods.
const (
	// MethodGet is the HTTP GET method.
	MethodGet = "GET"
	// MethodDelete is the HTTP DELETE method.
	MethodDelete = "DELETE"
	// MethodHead is the HTTP HEAD method.
	MethodHead = "HEAD"
	// MethodOptions is the HTTP OPTIONS method.
	MethodOptions = "OPTIONS"
	// MethodPatch is the HTTP PATCH method.
	MethodPatch = "PATCH"
	// MethodPost is the HTTP POST method.
	MethodPost = "POST"
	// MethodPut is the HTTP PUT method.
	MethodPut = "PUT"
)

// HTTP header names.
const (
	// HeaderContentType is the "Content-Type" header.
	HeaderContentType = "Content-Type"
	// HeaderLocation is the "Location" header.
	HeaderLocation = "Location"
	// HeaderReferrer is the "Referer" header.
	HeaderReferrer = "Referrer"
	// HeaderSetCookie is the "Set-Cookie" header.
	HeaderSetCookie = "Set-Cookie"
)

// MIME types for HTTP content type.
const (
	// MIMETypeJSON indicates the "application/json" MIME type.
	MIMETypeJSON = "application/json"
	// MIMETypeHTML indicates the "text/html" MIME type.
	MIMETypeHTML = "text/html"
	// MIMETypeText indicates the "text/plain" MIME type.
	MIMETypeText = "text/plain"
)
