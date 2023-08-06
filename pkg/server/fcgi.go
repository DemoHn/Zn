package server

///// copied from net/http/fcgi/fcgi.go

// Package fcgi implements the FastCGI protocol.
//
// See https://fast-cgi.github.io/ for an unofficial mirror of the
// original documentation.
//
// Currently only the responder role is supported.

// This file defines the raw protocol and some utilities used by the child and
// the host.

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/cgi"
	"os"
	"strings"
	"sync"
	"time"
)

// recType is a record type, as defined by
// https://web.archive.org/web/20150420080736/http://www.fastcgi.com/drupal/node/6?q=node/22#S8
type recType uint8

const (
	typeBeginRequest    recType = 1
	typeAbortRequest    recType = 2
	typeEndRequest      recType = 3
	typeParams          recType = 4
	typeStdin           recType = 5
	typeStdout          recType = 6
	typeStderr          recType = 7
	typeData            recType = 8
	typeGetValues       recType = 9
	typeGetValuesResult recType = 10
	typeUnknownType     recType = 11
)

// keep the connection between web-server and responder open after request
const flagKeepConn = 1

const (
	maxWrite = 65535 // maximum record body
	maxPad   = 255
)

const (
	roleResponder = iota + 1 // only Responders are implemented.
	roleAuthorizer
	roleFilter
)

const (
	statusRequestComplete = iota
	statusCantMultiplex
	statusOverloaded
	statusUnknownRole
)

type header struct {
	Version       uint8
	Type          recType
	Id            uint16
	ContentLength uint16
	PaddingLength uint8
	Reserved      uint8
}

type beginRequest struct {
	role     uint16
	flags    uint8
	reserved [5]uint8
}

func (br *beginRequest) read(content []byte) error {
	if len(content) != 8 {
		return errors.New("fcgi: invalid begin request record")
	}
	br.role = binary.BigEndian.Uint16(content)
	br.flags = content[2]
	return nil
}

// for padding so we don't have to allocate all the time
// not synchronized because we don't care what the contents are
var pad [maxPad]byte

func (h *header) init(recType recType, reqId uint16, contentLength int) {
	h.Version = 1
	h.Type = recType
	h.Id = reqId
	h.ContentLength = uint16(contentLength)
	h.PaddingLength = uint8(-contentLength & 7)
}

// conn sends records over rwc
type conn struct {
	mutex sync.Mutex
	rwc   io.ReadWriteCloser

	// to avoid allocations
	buf bytes.Buffer
	h   header
}

func newConn(rwc io.ReadWriteCloser) *conn {
	return &conn{rwc: rwc}
}

func (c *conn) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.rwc.Close()
}

type record struct {
	h   header
	buf [maxWrite + maxPad]byte
}

func (rec *record) read(r io.Reader) (err error) {
	if err = binary.Read(r, binary.BigEndian, &rec.h); err != nil {
		return err
	}
	if rec.h.Version != 1 {
		return errors.New("fcgi: invalid header version")
	}
	n := int(rec.h.ContentLength) + int(rec.h.PaddingLength)
	if _, err = io.ReadFull(r, rec.buf[:n]); err != nil {
		return err
	}
	return nil
}

func (r *record) content() []byte {
	return r.buf[:r.h.ContentLength]
}

// writeRecord writes and sends a single record.
func (c *conn) writeRecord(recType recType, reqId uint16, b []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.buf.Reset()
	c.h.init(recType, reqId, len(b))
	if err := binary.Write(&c.buf, binary.BigEndian, c.h); err != nil {
		return err
	}
	if _, err := c.buf.Write(b); err != nil {
		return err
	}
	if _, err := c.buf.Write(pad[:c.h.PaddingLength]); err != nil {
		return err
	}
	_, err := c.rwc.Write(c.buf.Bytes())
	return err
}

func (c *conn) writeEndRequest(reqId uint16, appStatus int, protocolStatus uint8) error {
	b := make([]byte, 8)
	binary.BigEndian.PutUint32(b, uint32(appStatus))
	b[4] = protocolStatus
	return c.writeRecord(typeEndRequest, reqId, b)
}

func (c *conn) writePairs(recType recType, reqId uint16, pairs map[string]string) error {
	w := newWriter(c, recType, reqId)
	b := make([]byte, 8)
	for k, v := range pairs {
		n := encodeSize(b, uint32(len(k)))
		n += encodeSize(b[n:], uint32(len(v)))
		if _, err := w.Write(b[:n]); err != nil {
			return err
		}
		if _, err := w.WriteString(k); err != nil {
			return err
		}
		if _, err := w.WriteString(v); err != nil {
			return err
		}
	}
	w.Close()
	return nil
}

func readSize(s []byte) (uint32, int) {
	if len(s) == 0 {
		return 0, 0
	}
	size, n := uint32(s[0]), 1
	if size&(1<<7) != 0 {
		if len(s) < 4 {
			return 0, 0
		}
		n = 4
		size = binary.BigEndian.Uint32(s)
		size &^= 1 << 31
	}
	return size, n
}

func readString(s []byte, size uint32) string {
	if size > uint32(len(s)) {
		return ""
	}
	return string(s[:size])
}

func encodeSize(b []byte, size uint32) int {
	if size > 127 {
		size |= 1 << 31
		binary.BigEndian.PutUint32(b, size)
		return 4
	}
	b[0] = byte(size)
	return 1
}

// bufWriter encapsulates bufio.Writer but also closes the underlying stream when
// Closed.
type bufWriter struct {
	closer io.Closer
	*bufio.Writer
}

func (w *bufWriter) Close() error {
	if err := w.Writer.Flush(); err != nil {
		w.closer.Close()
		return err
	}
	return w.closer.Close()
}

func newWriter(c *conn, recType recType, reqId uint16) *bufWriter {
	s := &streamWriter{c: c, recType: recType, reqId: reqId}
	w := bufio.NewWriterSize(s, maxWrite)
	return &bufWriter{s, w}
}

// streamWriter abstracts out the separation of a stream into discrete records.
// It only writes maxWrite bytes at a time.
type streamWriter struct {
	c       *conn
	recType recType
	reqId   uint16
}

func (w *streamWriter) Write(p []byte) (int, error) {
	nn := 0
	for len(p) > 0 {
		n := len(p)
		if n > maxWrite {
			n = maxWrite
		}
		if err := w.c.writeRecord(w.recType, w.reqId, p[:n]); err != nil {
			return nn, err
		}
		nn += n
		p = p[n:]
	}
	return nn, nil
}

func (w *streamWriter) Close() error {
	// send empty record to close the stream
	return w.c.writeRecord(w.recType, w.reqId, nil)
}

///// copied from net/http/fcgi/child.go

// This file implements FastCGI from the perspective of a child process.

// request holds the state for an in-progress request. As soon as it's complete,
// it's converted to an http.Request.
type request struct {
	pw        *io.PipeWriter
	reqId     uint16
	params    map[string]string
	buf       [1024]byte
	rawParams []byte
	keepConn  bool
}

// envVarsContextKey uniquely identifies a mapping of CGI
// environment variables to their values in a request context
type envVarsContextKey struct{}

func newRequest(reqId uint16, flags uint8) *request {
	r := &request{
		reqId:    reqId,
		params:   map[string]string{},
		keepConn: flags&flagKeepConn != 0,
	}
	r.rawParams = r.buf[:0]
	return r
}

// parseParams reads an encoded []byte into Params.
func (r *request) parseParams() {
	text := r.rawParams
	r.rawParams = nil
	for len(text) > 0 {
		keyLen, n := readSize(text)
		if n == 0 {
			return
		}
		text = text[n:]
		valLen, n := readSize(text)
		if n == 0 {
			return
		}
		text = text[n:]
		if int(keyLen)+int(valLen) > len(text) {
			return
		}
		key := readString(text, keyLen)
		text = text[keyLen:]
		val := readString(text, valLen)
		text = text[valLen:]
		r.params[key] = val
	}
}

// response implements http.ResponseWriter.
type response struct {
	req            *request
	header         http.Header
	code           int
	wroteHeader    bool
	wroteCGIHeader bool
	w              *bufWriter
}

func newResponse(c *child, req *request) *response {
	return &response{
		req:    req,
		header: http.Header{},
		w:      newWriter(c.conn, typeStdout, req.reqId),
	}
}

func (r *response) Header() http.Header {
	return r.header
}

func (r *response) Write(p []byte) (n int, err error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	if !r.wroteCGIHeader {
		r.writeCGIHeader(p)
	}
	return r.w.Write(p)
}

func (r *response) WriteHeader(code int) {
	if r.wroteHeader {
		return
	}
	r.wroteHeader = true
	r.code = code
	if code == http.StatusNotModified {
		// Must not have body.
		r.header.Del("Content-Type")
		r.header.Del("Content-Length")
		r.header.Del("Transfer-Encoding")
	}
	if r.header.Get("Date") == "" {
		r.header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
	}
}

// writeCGIHeader finalizes the header sent to the client and writes it to the output.
// p is not written by writeHeader, but is the first chunk of the body
// that will be written. It is sniffed for a Content-Type if none is
// set explicitly.
func (r *response) writeCGIHeader(p []byte) {
	if r.wroteCGIHeader {
		return
	}
	r.wroteCGIHeader = true
	fmt.Fprintf(r.w, "Status: %d %s\r\n", r.code, http.StatusText(r.code))
	if _, hasType := r.header["Content-Type"]; r.code != http.StatusNotModified && !hasType {
		r.header.Set("Content-Type", http.DetectContentType(p))
	}
	r.header.Write(r.w)
	r.w.WriteString("\r\n")
	r.w.Flush()
}

func (r *response) Flush() {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	r.w.Flush()
}

func (r *response) Close() error {
	r.Flush()
	return r.w.Close()
}

type child struct {
	conn    *conn
	handler http.Handler

	mu       sync.Mutex          // protects requests:
	requests map[uint16]*request // keyed by request ID
}

func newChild(rwc io.ReadWriteCloser, handler http.Handler) *child {
	return &child{
		conn:     newConn(rwc),
		handler:  handler,
		requests: make(map[uint16]*request),
	}
}

func (c *child) serve() {
	defer c.conn.Close()
	defer c.cleanUp()
	var rec record
	for {
		if err := rec.read(c.conn.rwc); err != nil {
			return
		}
		if err := c.handleRecord(&rec); err != nil {
			return
		}
	}
}

var errCloseConn = errors.New("fcgi: connection should be closed")

var emptyBody = io.NopCloser(strings.NewReader(""))

// ErrRequestAborted is returned by Read when a handler attempts to read the
// body of a request that has been aborted by the web server.
var ErrRequestAborted = errors.New("fcgi: request aborted by web server")

// ErrConnClosed is returned by Read when a handler attempts to read the body of
// a request after the connection to the web server has been closed.
var ErrConnClosed = errors.New("fcgi: connection to web server closed")

func (c *child) handleRecord(rec *record) error {
	c.mu.Lock()
	req, ok := c.requests[rec.h.Id]
	c.mu.Unlock()
	if !ok && rec.h.Type != typeBeginRequest && rec.h.Type != typeGetValues {
		// The spec says to ignore unknown request IDs.
		return nil
	}

	switch rec.h.Type {
	case typeBeginRequest:
		if req != nil {
			// The server is trying to begin a request with the same ID
			// as an in-progress request. This is an error.
			return errors.New("fcgi: received ID that is already in-flight")
		}

		var br beginRequest
		if err := br.read(rec.content()); err != nil {
			return err
		}
		if br.role != roleResponder {
			c.conn.writeEndRequest(rec.h.Id, 0, statusUnknownRole)
			return nil
		}
		req = newRequest(rec.h.Id, br.flags)
		c.mu.Lock()
		c.requests[rec.h.Id] = req
		c.mu.Unlock()
		return nil
	case typeParams:
		// NOTE(eds): Technically a key-value pair can straddle the boundary
		// between two packets. We buffer until we've received all parameters.
		if len(rec.content()) > 0 {
			req.rawParams = append(req.rawParams, rec.content()...)
			return nil
		}
		req.parseParams()
		return nil
	case typeStdin:
		content := rec.content()
		if req.pw == nil {
			var body io.ReadCloser
			if len(content) > 0 {
				// body could be an io.LimitReader, but it shouldn't matter
				// as long as both sides are behaving.
				body, req.pw = io.Pipe()
			} else {
				body = emptyBody
			}
			go c.serveRequest(req, body)
		}
		if len(content) > 0 {
			// TODO(eds): This blocks until the handler reads from the pipe.
			// If the handler takes a long time, it might be a problem.
			req.pw.Write(content)
		} else if req.pw != nil {
			req.pw.Close()
		}
		return nil
	case typeGetValues:
		values := map[string]string{"FCGI_MPXS_CONNS": "1"}
		c.conn.writePairs(typeGetValuesResult, 0, values)
		return nil
	case typeData:
		// If the filter role is implemented, read the data stream here.
		return nil
	case typeAbortRequest:
		c.mu.Lock()
		delete(c.requests, rec.h.Id)
		c.mu.Unlock()
		c.conn.writeEndRequest(rec.h.Id, 0, statusRequestComplete)
		if req.pw != nil {
			req.pw.CloseWithError(ErrRequestAborted)
		}
		if !req.keepConn {
			// connection will close upon return
			return errCloseConn
		}
		return nil
	default:
		b := make([]byte, 8)
		b[0] = byte(rec.h.Type)
		c.conn.writeRecord(typeUnknownType, 0, b)
		return nil
	}
}

// filterOutUsedEnvVars returns a new map of env vars without the
// variables in the given envVars map that are read for creating each http.Request
func filterOutUsedEnvVars(envVars map[string]string) map[string]string {
	withoutUsedEnvVars := make(map[string]string)
	for k, v := range envVars {
		if addFastCGIEnvToContext(k) {
			withoutUsedEnvVars[k] = v
		}
	}
	return withoutUsedEnvVars
}

func (c *child) serveRequest(req *request, body io.ReadCloser) {
	r := newResponse(c, req)
	httpReq, err := cgi.RequestFromMap(req.params)
	if err != nil {
		// there was an error reading the request
		r.WriteHeader(http.StatusInternalServerError)
		c.conn.writeRecord(typeStderr, req.reqId, []byte(err.Error()))
	} else {
		httpReq.Body = body
		withoutUsedEnvVars := filterOutUsedEnvVars(req.params)
		envVarCtx := context.WithValue(httpReq.Context(), envVarsContextKey{}, withoutUsedEnvVars)
		httpReq = httpReq.WithContext(envVarCtx)
		c.handler.ServeHTTP(r, httpReq)
	}
	// Make sure we serve something even if nothing was written to r
	r.Write(nil)
	r.Close()
	c.mu.Lock()
	delete(c.requests, req.reqId)
	c.mu.Unlock()
	c.conn.writeEndRequest(req.reqId, 0, statusRequestComplete)

	// Consume the entire body, so the host isn't still writing to
	// us when we close the socket below in the !keepConn case,
	// otherwise we'd send a RST. (golang.org/issue/4183)
	// TODO(bradfitz): also bound this copy in time. Or send
	// some sort of abort request to the host, so the host
	// can properly cut off the client sending all the data.
	// For now just bound it a little and
	io.CopyN(io.Discard, body, 100<<20)
	body.Close()

	if !req.keepConn {
		c.conn.Close()
	}
}

func (c *child) cleanUp() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, req := range c.requests {
		if req.pw != nil {
			// race with call to Close in c.serveRequest doesn't matter because
			// Pipe(Reader|Writer).Close are idempotent
			req.pw.CloseWithError(ErrConnClosed)
		}
	}
}

// addFastCGIEnvToContext reports whether to include the FastCGI environment variable s
// in the http.Request.Context, accessible via ProcessEnv.
func addFastCGIEnvToContext(s string) bool {
	// Exclude things supported by net/http natively:
	switch s {
	case "CONTENT_LENGTH", "CONTENT_TYPE", "HTTPS",
		"PATH_INFO", "QUERY_STRING", "REMOTE_ADDR",
		"REMOTE_HOST", "REMOTE_PORT", "REQUEST_METHOD",
		"REQUEST_URI", "SCRIPT_NAME", "SERVER_PROTOCOL":
		return false
	}
	if strings.HasPrefix(s, "HTTP_") {
		return false
	}
	// Explicitly include FastCGI-specific things.
	// This list is redundant with the default "return true" below.
	// Consider this documentation of the sorts of things we expect
	// to maybe see.
	switch s {
	case "REMOTE_USER":
		return true
	}
	// Unknown, so include it to be safe.
	return true
}

//////////////////////////////
/////  PUBLIC FUNCTIONS //////
//////////////////////////////

// Serve accepts incoming FastCGI connections on the listener l, creating a new
// goroutine for each. The goroutine reads requests and then calls handler
// to reply to them.
// If l is nil, Serve accepts connections from os.Stdin.
// If handler is nil, http.DefaultServeMux is used.
func Serve(l net.Listener, handler http.Handler) error {
	if l == nil {
		var err error
		l, err = net.FileListener(os.Stdin)
		if err != nil {
			return err
		}
		defer l.Close()
	}
	if handler == nil {
		handler = http.DefaultServeMux
	}
	for {
		rw, err := l.Accept()
		if err != nil {
			return err
		}
		c := newChild(rw, handler)
		go c.serve()
	}
}

// ProcessEnv returns FastCGI environment variables associated with the request r
// for which no effort was made to be included in the request itself - the data
// is hidden in the request's context. As an example, if REMOTE_USER is set for a
// request, it will not be found anywhere in r, but it will be included in
// ProcessEnv's response (via r's context).
func ProcessFCGIEnv(r *http.Request) map[string]string {
	env, _ := r.Context().Value(envVarsContextKey{}).(map[string]string)
	return env
}

func AcceptFCGIRequest(conn net.Conn, handler http.Handler) {
	c := newChild(conn, handler)
	c.serve()
}
