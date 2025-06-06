package email

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

// WriteTo writes the entire email message to the provided io.Writer in MIME format.
// It implements the io.WriterTo interface. Returns the number of bytes written and any error encountered.
func (m *Message) WriteTo(w io.Writer) (int64, error) {
	mw := &messageWriter{w: w}
	mw.writeMessage(m)
	return mw.n, mw.err
}

// writeMessage writes the provided Message to the underlying writer in MIME format.
func (w *messageWriter) writeMessage(m *Message) {
	if _, ok := m.header["Mime-Version"]; !ok {
		w.writeString("Mime-Version: 1.0\r\n")
	}
	if _, ok := m.header["Date"]; !ok {
		w.writeHeader("Date", m.FormatDate(now()))
	}
	w.writeHeaders(m.header)

	if m.hasMixedPart() {
		w.openMultipart("mixed")
	}

	if m.hasAlternativePart() {
		w.openMultipart("alternative")
	}
	for _, part := range m.parts {
		w.writePart(part, m.charset)
	}
	if m.hasAlternativePart() {
		w.closeMultipart()
	}

	w.addFiles(m.attachments, true)
	if m.hasMixedPart() {
		w.closeMultipart()
	}
}

func (m *Message) hasMixedPart() bool {
	return (len(m.parts) > 0 && len(m.attachments) > 0) || len(m.attachments) > 1
}

func (m *Message) hasAlternativePart() bool {
	return len(m.parts) > 1
}

type messageWriter struct {
	w          io.Writer
	n          int64
	writers    [3]*multipart.Writer
	partWriter io.Writer
	depth      uint8
	err        error
}

// openMultipart opens a new multipart section with the specified MIME type.
func (w *messageWriter) openMultipart(mimeType string) {
	mw := multipart.NewWriter(w)
	contentType := "multipart/" + mimeType + ";\r\n boundary=" + mw.Boundary()
	w.writers[w.depth] = mw

	if w.depth == 0 {
		w.writeHeader("Content-Type", contentType)
		w.writeString("\r\n")
	} else {
		w.createPart(map[string][]string{
			"Content-Type": {contentType},
		})
	}
	w.depth++
}

// createPart creates a new MIME part with the provided headers.
func (w *messageWriter) createPart(h map[string][]string) {
	w.partWriter, w.err = w.writers[w.depth-1].CreatePart(h)
}

// closeMultipart closes the current multipart section.
func (w *messageWriter) closeMultipart() {
	if w.depth > 0 {
		_ = w.writers[w.depth-1].Close()
		w.depth--
	}
}

// writePart writes a single part with the specified charset.
func (w *messageWriter) writePart(p *part, charset string) {
	w.writeHeaders(map[string][]string{
		"Content-Type":              {p.contentType + "; charset=" + charset},
		"Content-Transfer-Encoding": {string(p.encoding)},
	})
	w.writeBody(p.copier, p.encoding)
}

// addFiles adds the provided files as attachments or inline parts.
func (w *messageWriter) addFiles(files []*file, isAttachment bool) {
	for _, f := range files {
		if _, ok := f.Header["Content-Type"]; !ok {
			mediaType := mime.TypeByExtension(filepath.Ext(f.Name))
			if mediaType == "" {
				mediaType = "application/octet-stream"
			}
			f.setHeader("Content-Type", mediaType+`; name="`+f.Name+`"`)
		}

		if _, ok := f.Header["Content-Transfer-Encoding"]; !ok {
			f.setHeader("Content-Transfer-Encoding", string(Base64))
		}

		if _, ok := f.Header["Content-Disposition"]; !ok {
			var disp string
			if isAttachment {
				disp = "attachment"
			} else {
				disp = "inline"
			}
			f.setHeader("Content-Disposition", disp+`; filename="`+f.Name+`"`)
		}

		if !isAttachment {
			if _, ok := f.Header["Content-ID"]; !ok {
				f.setHeader("Content-ID", "<"+f.Name+">")
			}
		}

		w.writeHeaders(f.Header)
		w.writeBody(func(w io.Writer) error {
			r := bytes.NewReader(f.Data)
			if _, err := io.Copy(w, r); err != nil {
				return err
			}
			return nil
		}, Base64)
	}
}

// Write writes bytes to the underlying writer, tracking errors and byte count.
func (w *messageWriter) Write(p []byte) (int, error) {
	if w.err != nil {
		return 0, errors.New("gomail: cannot write as writer is in error")
	}

	var n int
	n, w.err = w.w.Write(p)
	w.n += int64(n)
	return n, w.err
}

// writeString writes a string to the underlying writer.
func (w *messageWriter) writeString(s string) {
	n, _ := io.WriteString(w.w, s)
	w.n += int64(n)
}

// writeHeader writes a single header field with the provided key and values.
func (w *messageWriter) writeHeader(k string, v ...string) {
	w.writeString(k)
	if len(v) == 0 {
		w.writeString(":\r\n")
		return
	}
	w.writeString(": ")

	// Max header line length is 78 characters in RFC 5322 and 76 characters
	// in RFC 2047. So for the sake of simplicity we use the 76 characters
	// limit.
	charsLeft := 76 - len(k) - len(": ")

	for i, s := range v {
		// If the line is already too long, insert a newline right away.
		if charsLeft < 1 {
			if i == 0 {
				w.writeString("\r\n ")
			} else {
				w.writeString(",\r\n ")
			}
			charsLeft = 75
		} else if i != 0 {
			w.writeString(", ")
			charsLeft -= 2
		}

		// While the header content is too long, fold it by inserting a newline.
		for len(s) > charsLeft {
			s = w.writeLine(s, charsLeft)
			charsLeft = 75
		}
		w.writeString(s)
		if i := lastIndexByte(s, '\n'); i != -1 {
			charsLeft = 75 - (len(s) - i - 1)
		} else {
			charsLeft -= len(s)
		}
	}
	w.writeString("\r\n")
}

func (w *messageWriter) writeLine(s string, charsLeft int) string {
	// If there is already a newline before the limit. Write the line.
	if i := strings.IndexByte(s, '\n'); i != -1 && i < charsLeft {
		w.writeString(s[:i+1])
		return s[i+1:]
	}

	for i := charsLeft - 1; i >= 0; i-- {
		if s[i] == ' ' {
			w.writeString(s[:i])
			w.writeString("\r\n ")
			return s[i+1:]
		}
	}

	// We could not insert a newline cleanly so look for a space or a newline
	// even if it is after the limit.
	for i := 75; i < len(s); i++ {
		if s[i] == ' ' {
			w.writeString(s[:i])
			w.writeString("\r\n ")
			return s[i+1:]
		}
		if s[i] == '\n' {
			w.writeString(s[:i+1])
			return s[i+1:]
		}
	}

	// Too bad, no space or newline in the whole string. Just write everything.
	w.writeString(s)
	return ""
}

func (w *messageWriter) writeHeaders(h map[string][]string) {
	if w.depth == 0 {
		for k, v := range h {
			if k != "Bcc" {
				w.writeHeader(k, v...)
			}
		}
	} else {
		w.createPart(h)
	}
}

func (w *messageWriter) writeBody(f func(io.Writer) error, enc Encoding) {
	var subWriter io.Writer
	if w.depth == 0 {
		w.writeString("\r\n")
		subWriter = w.w
	} else {
		subWriter = w.partWriter
	}

	switch enc {
	case Base64:
		wc := base64.NewEncoder(base64.StdEncoding, newBase64LineWriter(subWriter))
		w.err = f(wc)
		_ = wc.Close()
	case Unencoded:
		w.err = f(subWriter)
	default:
		wc := newQPWriter(subWriter)
		w.err = f(wc)
		_ = wc.Close()
	}
}

// As required by RFC 2045, 6.7. (page 21) for quoted-printable, and
// RFC 2045, 6.8. (page 25) for base64.
const maxLineLen = 76

// base64LineWriter limits text encoded in base64 to 76 characters per line
type base64LineWriter struct {
	w       io.Writer
	lineLen int
}

func newBase64LineWriter(w io.Writer) *base64LineWriter {
	return &base64LineWriter{w: w}
}

func (w *base64LineWriter) Write(p []byte) (int, error) {
	n := 0
	for len(p)+w.lineLen > maxLineLen {
		_, err := w.w.Write(p[:maxLineLen-w.lineLen])
		if err != nil {
			return n, err
		}
		_, err = w.w.Write([]byte("\r\n"))
		if err != nil {
			return n, err
		}
		p = p[maxLineLen-w.lineLen:]
		n += maxLineLen - w.lineLen
		w.lineLen = 0
	}

	_, err := w.w.Write(p)
	if err != nil {
		return n, err
	}
	w.lineLen += len(p)

	return n + len(p), nil
}

// Stubbed out for testing.
var now = time.Now
