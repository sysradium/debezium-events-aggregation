package models

type Transaction struct {
	ID                  string `json:"id"`
	TotalOrder          int64  `json:"total_order"`
	DataCollectionOrder int64  `json:"data_collection_order"`
}

type TransactionMetadataEvent struct {
    Schema  TransactionSchema `json:"schema"`
    Payload TransactionPayload `json:"payload"`
}

type TransactionSchema struct {
    Type   string                `json:"type"`
    Fields []TransactionField    `json:"fields"`
    Name   string                `json:"name"`
    Version int                  `json:"version"`
}

type TransactionField struct {
    Type     string              `json:"type"`
    Optional bool                `json:"optional"`
    Field    string              `json:"field"`
    Items    *TransactionItems   `json:"items,omitempty"`
}

type TransactionItems struct {
    Type   string                  `json:"type"`
    Fields []TransactionInnerField `json:"fields"`
    Name   string                  `json:"name"`
    Version int                    `json:"version"`
}

type TransactionInnerField struct {
    Type     string `json:"type"`
    Optional bool   `json:"optional"`
    Field    string `json:"field"`
}

type TransactionPayload struct {
    Status          string                `json:"status"`
    ID              string                `json:"id"`
    EventCount      int64                 `json:"event_count,omitempty"`
    DataCollections []DataCollection      `json:"data_collections,omitempty"`
    TsMs            int64                 `json:"ts_ms"`
}

type DataCollection struct {
    DataCollection string `json:"data_collection"`
    EventCount     int64  `json:"event_count"`
}
