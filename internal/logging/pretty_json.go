package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
)

// PrettyJSONHandler wraps slog.JSONHandler to pretty-print JSON output
type PrettyJSONHandler struct {
	handler slog.Handler
	writer  io.Writer
}

// NewPrettyJSONHandler creates a new pretty-printing JSON handler
func NewPrettyJSONHandler(w io.Writer, opts *slog.HandlerOptions) *PrettyJSONHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}

	// Create a buffer to capture JSON output
	buf := &bytes.Buffer{}

	return &PrettyJSONHandler{
		handler: slog.NewJSONHandler(buf, opts),
		writer:  w,
	}
}

func (h *PrettyJSONHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *PrettyJSONHandler) Handle(ctx context.Context, r slog.Record) error {
	// Get the underlying JSONHandler's buffer
	buf := &bytes.Buffer{}
	tempHandler := slog.NewJSONHandler(buf, nil)

	if err := tempHandler.Handle(ctx, r); err != nil {
		return err
	}

	// Parse and re-format the JSON with indentation
	var data map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
		// If parsing fails, just write the original
		_, err := h.writer.Write(buf.Bytes())
		return err
	}

	// Pretty-print with indentation
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		_, err := h.writer.Write(buf.Bytes())
		return err
	}

	// Write the pretty JSON with a newline
	_, err = h.writer.Write(append(prettyJSON, '\n'))
	return err
}

func (h *PrettyJSONHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyJSONHandler{
		handler: h.handler.WithAttrs(attrs),
		writer:  h.writer,
	}
}

func (h *PrettyJSONHandler) WithGroup(name string) slog.Handler {
	return &PrettyJSONHandler{
		handler: h.handler.WithGroup(name),
		writer:  h.writer,
	}
}
