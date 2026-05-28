package log

import (
	"bytes"
	"encoding/json"
	"sort"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var bufPool = buffer.NewPool()

const logTimeLayout = "2006-01-02 15:04:05.000"

// 固定前缀字段顺序：level → ts → request_id → caller → msg → 其余按字母序
var logFieldOrder = []string{"level", "ts", "request_id", "caller", "msg"}

func init() {
	_ = zap.RegisterEncoder("ordered_json", func(cfg zapcore.EncoderConfig) (zapcore.Encoder, error) {
		return newOrderedJSONEncoder(cfg), nil
	})
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(logTimeLayout))
}

type orderedJSONEncoder struct {
	zapcore.Encoder
}

func newOrderedJSONEncoder(cfg zapcore.EncoderConfig) *orderedJSONEncoder {
	return &orderedJSONEncoder{Encoder: zapcore.NewJSONEncoder(cfg)}
}

func (e *orderedJSONEncoder) Clone() zapcore.Encoder {
	return &orderedJSONEncoder{Encoder: e.Encoder.Clone()}
}

func (e *orderedJSONEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	buf, err := e.Encoder.EncodeEntry(ent, fields)
	if err != nil {
		return nil, err
	}
	out := reorderLogJSON(buf.Bytes())
	buf.Free()
	result := bufPool.Get()
	result.Write(out)
	return result, nil
}

func reorderLogJSON(raw []byte) []byte {
	raw = bytes.TrimSpace(raw)
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return append(bytes.Clone(raw), '\n')
	}

	ordered := make([]string, 0, len(fields))
	seen := make(map[string]struct{}, len(fields))
	for _, key := range logFieldOrder {
		if _, ok := fields[key]; ok {
			ordered = append(ordered, key)
			seen[key] = struct{}{}
		}
	}

	rest := make([]string, 0, len(fields)-len(seen))
	for key := range fields {
		if _, ok := seen[key]; !ok {
			rest = append(rest, key)
		}
	}
	sort.Strings(rest)
	ordered = append(ordered, rest...)

	var b bytes.Buffer
	b.WriteByte('{')
	for i, key := range ordered {
		if i > 0 {
			b.WriteByte(',')
		}
		keyJSON, _ := json.Marshal(key)
		b.Write(keyJSON)
		b.WriteByte(':')
		b.Write(fields[key])
	}
	b.WriteByte('}')
	b.WriteByte('\n')
	return b.Bytes()
}
