package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// ANSI 颜色（开发环境终端）：debug 蓝 / info 白 / warn 黄 / error 红
const (
	ansiReset  = "\x1b[0m"
	ansiBlue   = "\x1b[34m"
	ansiYellow = "\x1b[33m"
	ansiRed    = "\x1b[31m"
	ansiWhite  = "\x1b[37m"
)

func init() {
	_ = zap.RegisterEncoder("dev_console", func(cfg zapcore.EncoderConfig) (zapcore.Encoder, error) {
		return newDevConsoleEncoder(cfg), nil
	})
}

type devConsoleEncoder struct {
	zapcore.Encoder
}

func newDevConsoleEncoder(cfg zapcore.EncoderConfig) *devConsoleEncoder {
	jsonCfg := devJSONEncoderConfig(cfg)
	return &devConsoleEncoder{Encoder: zapcore.NewJSONEncoder(jsonCfg)}
}

func devJSONEncoderConfig(cfg zapcore.EncoderConfig) zapcore.EncoderConfig {
	jsonCfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        zapcore.OmitKey,
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  zapcore.OmitKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	if cfg.EncodeTime != nil {
		jsonCfg.EncodeTime = cfg.EncodeTime
	}
	return jsonCfg
}

func (e *devConsoleEncoder) Clone() zapcore.Encoder {
	return &devConsoleEncoder{Encoder: e.Encoder.Clone()}
}

func (e *devConsoleEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	buf, err := e.Encoder.EncodeEntry(ent, fields)
	if err != nil {
		return nil, err
	}
	out := formatDevConsoleLine(buf.Bytes())
	buf.Free()
	result := bufPool.Get()
	result.Write(out)
	return result, nil
}

// formatDevConsoleLine 将 JSON 行转为：
// 2026-05-26 10:34:45.794 INFO [request_id=xxx] handler/foo.go:23 message {"k":"v"}
func formatDevConsoleLine(raw []byte) []byte {
	raw = bytes.TrimSpace(raw)
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return append(bytes.Clone(raw), '\n')
	}

	var b strings.Builder
	if ts, ok := fields["ts"]; ok {
		b.WriteString(strings.Trim(string(ts), `"`))
		delete(fields, "ts")
	}
	b.WriteByte(' ')

	if level, ok := fields["level"]; ok {
		b.WriteString(colorizeLevel(strings.Trim(string(level), `"`)))
		delete(fields, "level")
	}

	if rid, ok := fields["request_id"]; ok {
		ridStr := strings.Trim(string(rid), `"`)
		if ridStr != "" {
			b.WriteString(fmt.Sprintf(" [request_id=%s]", ridStr))
		}
		delete(fields, "request_id")
	}

	if caller, ok := fields["caller"]; ok {
		b.WriteByte(' ')
		b.WriteString(strings.Trim(string(caller), `"`))
		delete(fields, "caller")
	}

	if msg, ok := fields["msg"]; ok {
		b.WriteByte(' ')
		b.WriteString(strings.Trim(string(msg), `"`))
		delete(fields, "msg")
	}

	if method, ok := fields["http_method"]; ok {
		fields["method"] = method
		delete(fields, "http_method")
	}

	if len(fields) > 0 {
		b.WriteByte(' ')
		b.Write(marshalTailJSON(fields))
	}

	b.WriteByte('\n')
	return []byte(b.String())
}

func colorizeLevel(level string) string {
	switch level {
	case "DEBUG":
		return ansiBlue + level + ansiReset
	case "INFO":
		return ansiWhite + level + ansiReset
	case "WARN":
		return ansiYellow + level + ansiReset
	case "ERROR", "DPANIC", "PANIC", "FATAL":
		return ansiRed + level + ansiReset
	default:
		return level
	}
}

func marshalTailJSON(fields map[string]json.RawMessage) []byte {
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b bytes.Buffer
	b.WriteByte('{')
	for i, k := range keys {
		if i > 0 {
			b.WriteByte(',')
		}
		keyJSON, _ := json.Marshal(k)
		b.Write(keyJSON)
		b.WriteByte(':')
		b.Write(fields[k])
	}
	b.WriteByte('}')
	return b.Bytes()
}
