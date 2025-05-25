// github.com/DauletBai/shanraq.org/http/context.go
package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5" 
	"github.com/DauletBai/shanraq.org/core/kernel"
	"github.com/DauletBai/shanraq.org/core/logger"
)

type HandlerFunc func(*Context)
// Context wraps the standard http.ResponseWriter and *http.Request,
// and provides access to framework services and helper methods.
type Context struct {
	BaseResponseWriter http.ResponseWriter
	BaseRequest        *http.Request
	kernel             *kernel.Kernel
	params             map[string]string // For storing URL parameters, form data, etc. after parsing
}

// NewContext creates a new Shanraq HTTP Context.
// This is typically called by the adapter that integrates with the chosen router.
func NewContext(w http.ResponseWriter, r *http.Request, k *kernel.Kernel) *Context {
	return &Context{
		BaseResponseWriter: w,
		BaseRequest:        r,
		kernel:             k,
		params:             make(map[string]string), // Initialize params map
	}
}

// Kernel returns the application kernel.
func (c *Context) Kernel() *kernel.Kernel {
	return c.kernel
}

// Logger returns the application logger from the kernel.
func (c *Context) Logger() logger.Interface {
	return c.kernel.Logger()
}

// Request returns the underlying *http.Request.
func (c *Context) Request() *http.Request {
	return c.BaseRequest
}

// Response returns the underlying http.ResponseWriter.
func (c *Context) Response() http.ResponseWriter {
	return c.BaseResponseWriter
}

// SetHeader sets a header on the response.
func (c *Context) SetHeader(key, value string) {
	c.BaseResponseWriter.Header().Set(key, value)
}

// ---- URL Parameters ----

// URLParam retrieves a URL parameter by name.
// This relies on the router (e.g., chi) to have parsed these into the request context.
func (c *Context) URLParam(name string) string {
	return chi.URLParam(c.BaseRequest, name)
}

// URLParamInt retrieves a URL parameter by name and converts it to an int.
// Returns 0 and an error if the parameter is not found or not a valid integer.
func (c *Context) URLParamInt(name string) (int, error) {
	valStr := c.URLParam(name)
	if valStr == "" {
		return 0, NewHTTPError(http.StatusBadRequest, "missing URL parameter: "+name)
	}
	valInt, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, NewHTTPError(http.StatusBadRequest, "invalid URL parameter format for '"+name+"': expected integer")
	}
	return valInt, nil
}

// ---- Query Parameters ----

// QueryParam retrieves a query parameter from the URL.
func (c *Context) QueryParam(name string) string {
	return c.BaseRequest.URL.Query().Get(name)
}

// QueryParamInt retrieves a query parameter and converts it to an int.
func (c *Context) QueryParamInt(name string, defaultValue ...int) (int, error) {
	valStr := c.QueryParam(name)
	if valStr == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0], nil
		}
		return 0, NewHTTPError(http.StatusBadRequest, "missing query parameter: "+name)
	}
	valInt, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, NewHTTPError(http.StatusBadRequest, "invalid query parameter format for '"+name+"': expected integer")
	}
	return valInt, nil
}

// QueryParams retrieves all query parameters.
func (c *Context) QueryParams() map[string][]string {
	return c.BaseRequest.URL.Query()
}

// ---- Form Data ----

// FormValue retrieves a form value by name (from POST/PUT form data or URL query).
// It calls ParseMultipartForm or ParseForm if necessary.
func (c *Context) FormValue(name string) string {
	// Ensure form is parsed. Max 32MB form data.
	if c.BaseRequest.Form == nil {
		// This value (32 << 20) is the default used by http.Request.ParseMultipartForm
		if err := c.BaseRequest.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			// Log error but don't necessarily stop, ParseForm might still work
			c.Logger().Error("Error parsing multipart form, falling back to ParseForm", "error", err)
		}
		if err := c.BaseRequest.ParseForm(); err != nil {
			c.Logger().Error("Error parsing form", "error", err)
			return "" // Return empty if form parsing fails
		}
	}
	return c.BaseRequest.FormValue(name)
}

// ---- JSON Handling ----

// BindJSON binds the request body as JSON into the provided structure.
func (c *Context) BindJSON(obj interface{}) error {
	if c.BaseRequest.Body == nil {
		return NewHTTPError(http.StatusBadRequest, "request body is empty")
	}
	defer c.BaseRequest.Body.Close()
	decoder := json.NewDecoder(c.BaseRequest.Body)
	// Consider adding options like decoder.DisallowUnknownFields() based on config
	if err := decoder.Decode(obj); err != nil {
		return NewHTTPError(http.StatusBadRequest, "failed to bind JSON: "+err.Error())
	}
	return nil
}

// JSON sends a JSON response with the given status code and data.
func (c *Context) JSON(statusCode int, data interface{}) error {
	c.SetHeader("Content-Type", "application/json; charset=UTF-8")
	c.BaseResponseWriter.WriteHeader(statusCode)
	encoder := json.NewEncoder(c.BaseResponseWriter)
	if err := encoder.Encode(data); err != nil {
		// If encoding fails, it's hard to recover the response, log it.
		c.Logger().Error("Failed to encode JSON response", "error", err, "statusCode", statusCode)
		return err // Return error so handler is aware
	}
	return nil
}

// ---- Error Handling ----

// Error sends a JSON error response.
// It uses the HTTPError type if the error is of that type,
// otherwise, it defaults to a 500 internal server error.
func (c *Context) Error(err error) {
	httpErr, ok := err.(*HTTPError)
	if !ok {
		// For non-HTTPError types, wrap it as a generic internal server error
		// Log the original error for debugging purposes
		c.Logger().Error("Internal server error", "original_error", err.Error())
		httpErr = NewHTTPError(http.StatusInternalServerError, "An unexpected error occurred.")
	} else {
		// If it's a client error (4xx), log it as INFO or WARN.
		// If it's a server error (5xx), log it as ERROR.
		if httpErr.Code >= 400 && httpErr.Code < 500 {
			c.Logger().Info("Client error", "status", httpErr.Code, "message", httpErr.Message, "details", httpErr.Details)
		} else if httpErr.Code >= 500 {
			c.Logger().Error("Server error", "status", httpErr.Code, "message", httpErr.Message, "details", httpErr.Details)
		}
	}

	// Attempt to send JSON error response
	if jsonErr := c.JSON(httpErr.Code, httpErr); jsonErr != nil {
		// Fallback if JSON encoding itself fails
		http.Error(c.BaseResponseWriter, `{"error":"Failed to marshal error response"}`, http.StatusInternalServerError)
	}
}

// ---- Other Response Types (Placeholders, can be expanded) ----

// String sends a plain text response.
func (c *Context) String(statusCode int, format string, args ...interface{}) error {
	c.SetHeader("Content-Type", "text/plain; charset=UTF-8")
	c.BaseResponseWriter.WriteHeader(statusCode)
	if _, err := c.BaseResponseWriter.Write([]byte(http.StatusText(statusCode) + "\n" + format)); err != nil {
		c.Logger().Error("Failed to write string response", "error", err)
		return err
	}
	// Note: fmt.Fprintf(c.BaseResponseWriter, format, args...) might be more conventional
	return nil
}

// HTML sends an HTML response. (Requires a template engine integration later)
func (c *Context) HTML(statusCode int, templateName string, data interface{}) error {
	// Placeholder for template rendering logic
	// e.g., err := c.kernel.TemplateRenderer().Render(c.BaseResponseWriter, templateName, data)
	// For now, just a simple message
	c.SetHeader("Content-Type", "text/html; charset=UTF-8")
	c.BaseResponseWriter.WriteHeader(statusCode)
	mockHTML := `<h1>Template: ` + templateName + `</h1><pre>` + MustMarshal(data) + `</pre>`
	if _, err := c.BaseResponseWriter.Write([]byte(mockHTML)); err != nil {
		c.Logger().Error("Failed to write HTML response", "error", err)
		return err
	}
	return nil
}

// MustMarshal is a helper for debugging, careful with sensitive data in production.
func MustMarshal(data interface{}) string {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "Error marshalling data: " + err.Error()
	}
	return string(b)
}