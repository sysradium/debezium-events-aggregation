package debezium

import "encoding/json"

type DebeziumEvent struct {
	Schema  Schema  `json:"schema"`
	Payload Payload `json:"payload"`
}

type Schema struct {
	Type    string  `json:"type"`
	Fields  []Field `json:"fields"`
	Name    string  `json:"name"`
	Version int     `json:"version"`
}

type Field struct {
	Type       string     `json:"type"`
	Fields     []Field    `json:"fields,omitempty"`
	Optional   bool       `json:"optional"`
	Name       string     `json:"name,omitempty"`
	Field      string     `json:"field"`
	Version    int        `json:"version,omitempty"`
	Parameters Parameters `json:"parameters,omitempty"`
}

type Parameters struct {
	Allowed string `json:"allowed"`
}

type Operation string

const (
	OPERATION_SNAPSHOT Operation = "r"
	OPERATION_UPDATE             = "u"
	OPERATION_CREATE             = "c"
	OPERATION_DELETE             = "d"
)

type Payload struct {
	Before      json.RawMessage `json:"before"`
	After       json.RawMessage `json:"after"`
	Source      Source          `json:"source"`
	Op          Operation       `json:"op"`
	TsMs        int64           `json:"ts_ms"`
	Transaction Transaction     `json:"transaction"`
}

type Source struct {
	Version   string `json:"version"`
	Connector string `json:"connector"`
	Name      string `json:"name"`
	TsMs      int64  `json:"ts_ms"`
	Snapshot  string `json:"snapshot"`
	DB        string `json:"db"`
	Sequence  string `json:"sequence,omitempty"`
	Table     string `json:"table"`
	ServerID  int64  `json:"server_id"`
	Gtid      string `json:"gtid,omitempty"`
	File      string `json:"file"`
	Pos       int64  `json:"pos"`
	Row       int32  `json:"row"`
	Thread    int64  `json:"thread,omitempty"`
	Query     string `json:"query,omitempty"`
}
