package xfer

// SchemaCard declares a synced extension table.
// Wire: schema TABLE VERSION HASH MTIME SIZE\nJSON
type SchemaCard struct {
	Table   string // table name (without x_ prefix)
	Version int
	Hash    string // SHA1 of canonical schema definition
	MTime   int64
	Content []byte // raw JSON schema definition
}

func (c *SchemaCard) Type() CardType { return CardSchema }

// XIGotCard announces possession of a table sync row.
// Wire: xigot TABLE PK_HASH MTIME
type XIGotCard struct {
	Table  string
	PKHash string
	MTime  int64
}

func (c *XIGotCard) Type() CardType { return CardXIGot }

// XGimmeCard requests a table sync row.
// Wire: xgimme TABLE PK_HASH
type XGimmeCard struct {
	Table  string
	PKHash string
}

func (c *XGimmeCard) Type() CardType { return CardXGimme }

// XRowCard carries a table sync row payload.
// Wire: xrow TABLE PK_HASH MTIME SIZE\nJSON
type XRowCard struct {
	Table   string
	PKHash  string
	MTime   int64
	Content []byte // raw JSON row data
}

func (c *XRowCard) Type() CardType { return CardXRow }

// XDeleteCard marks a table sync row as deleted (tombstone).
// Wire: xdelete TABLE PK_HASH MTIME SIZE\nJSON_PK_DATA
type XDeleteCard struct {
	Table  string
	PKHash string
	MTime  int64
	PKData []byte // JSON-encoded PK column values
}

func (c *XDeleteCard) Type() CardType { return CardXDelete }
