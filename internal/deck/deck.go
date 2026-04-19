package deck

import "time"

type ArtifactType int

const (
	Checkin    ArtifactType = iota
	Wiki
	Ticket
	Event
	Cluster
	ForumPost
	Attachment
	Control
)

type Deck struct {
	Type ArtifactType
	A    *AttachmentCard
	B    string
	C    string
	D    time.Time
	E    *EventCard
	F    []FileCard
	G    string
	H    string
	I    string
	J    []TicketField
	K    string
	L    string
	M    []string
	N    string
	P    []string
	Q    []CherryPick
	R    string
	T    []TagCard
	U    string
	W    []byte
	Z    string
}

type FileCard struct {
	Name    string
	UUID    string
	Perm    string
	OldName string
}

type TagCard struct {
	Type  TagType
	Name  string
	UUID  string
	Value string
}

type TagType byte

const (
	TagSingleton   TagType = '+'
	TagPropagating TagType = '*'
	TagCancel      TagType = '-'
)

type CherryPick struct {
	IsBackout bool
	Target    string
	Baseline  string
}

type AttachmentCard struct {
	Filename string
	Target   string
	Source   string
}

type EventCard struct {
	Date time.Time
	UUID string
}

type TicketField struct {
	Name  string
	Value string
}
