package traefik_add_trace_id

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

type traceIdFormat string

const (
	formatUUID traceIdFormat = "uuid"
	formatHex  traceIdFormat = "hex"
)

const defaultTraceID = "X-Trace-Id"
const defaultTraceIDFormat = formatUUID

// Config the plugin configuration.
type Config struct {
	HeaderPrefix string        `json:"headerPrefix"`
	HeaderName   string        `json:"headerName"`
	Format       traceIdFormat `json:"format"`
	Verbose      bool          `json:"verbose"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		HeaderPrefix: "",
		HeaderName:   defaultTraceID,
		Format:       defaultTraceIDFormat,
		Verbose:      false,
	}
}

// TraceIDHeader header if it's missing
type TraceIDHeader struct {
	headerName   string
	headerPrefix string
	format       traceIdFormat
	name         string
	next         http.Handler
	verbose      bool
}

// New created a new TraceIDHeader plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	if config == nil {
		return nil, fmt.Errorf("config can not be nil")
	}

	tIDHdr := &TraceIDHeader{
		next:    next,
		name:    name,
		verbose: config.Verbose,
	}

	if config.HeaderName == "" {
		tIDHdr.headerName = defaultTraceID
	} else {
		tIDHdr.headerName = config.HeaderName
	}

	if config.Format == formatUUID || config.Format == formatHex {
		tIDHdr.format = config.Format
	} else {
		tIDHdr.format = defaultTraceIDFormat
	}

	tIDHdr.headerPrefix = config.HeaderPrefix

	return tIDHdr, nil
}

func (t *TraceIDHeader) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	headerArr := req.Header[t.headerName]

	// continue with generation only if header is missing or empty
	if len(headerArr) <= 0 || headerArr[0] == "" {

		uuid := newUUID()
		traceId := ""

		if t.format == formatUUID {
			traceId = uuid.String()
		} else if t.format == formatHex {
			traceId = uuid.HexString()
		}

		randomUUID := fmt.Sprintf("%s%s", t.headerPrefix, traceId)

		req.Header.Set(t.headerName, randomUUID)

		if t.verbose {
			log.Println(req.Header[t.headerName][0])
		}
	}

	t.next.ServeHTTP(rw, req)
}
