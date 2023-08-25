package handle

import (
	"context"
	"encoding/json"
	"github.com/fatih/color"
	"io"
	"log"
	"log/slog"
)

type SlogHandlerOptions struct {
	Options *slog.HandlerOptions
}

type SlogHandler struct {
	slog.Handler
	options    SlogHandlerOptions
	attributes []slog.Attr
	log        *log.Logger
}

func (sho SlogHandlerOptions) NewSlogHandler(out io.Writer) *SlogHandler {
	return &SlogHandler{
		Handler: slog.NewJSONHandler(out, sho.Options),
		log:     log.New(out, "", 0),
	}
}

// Handle handles the Record. In this function we can add colors for level of debug.
// Also format a time and color text of message. Message we marshal in JSON.
func (sh SlogHandler) Handle(_ context.Context, record slog.Record) error {
	level := record.Level.String() + ":"
	switch record.Level {
	case slog.LevelDebug:
		level = color.BlueString(level)
	case slog.LevelInfo:
		level = color.GreenString(level)
	case slog.LevelWarn:
		level = color.CyanString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	time := record.Time.Format("[15:04:05.Mon]")
	time = color.YellowString(time)

	msg := color.WhiteString(record.Message)

	fields := make(map[string]interface{}, record.NumAttrs())
	record.Attrs(func(attr slog.Attr) bool {
		fields[attr.Key] = attr.Value.Any()
		return true
	})

	for _, attr := range sh.attributes {
		fields[attr.Key] = attr.Value.Any()
	}

	var b []byte
	var err error

	if len(fields) > 0 {
		b, err = json.MarshalIndent(fields, "", " ")
		if err != nil {
			return err
		}
	}

	sh.log.Println(level, time, msg, color.MagentaString(string(b)))
	return nil
}

func (sh SlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &SlogHandler{
		Handler:    sh.Handler,
		log:        sh.log,
		attributes: attrs,
	}
}
func (sh SlogHandler) WithGroup(name string) slog.Handler {
	return SlogHandler{
		Handler: sh.Handler.WithGroup(name),
		log:     sh.log,
	}
}
