package libfossil

import "github.com/danmestas/libfossil/internal/fsltype"

// FslID is a row-id in the blob table (content-addressed artifacts).
type FslID = fsltype.FslID

// FslSize represents a blob size; negative values indicate phantom blobs.
type FslSize = fsltype.FslSize

const (
	PhantomSize         FslSize = fsltype.PhantomSize
	FossilApplicationID int32   = fsltype.FossilApplicationID
)
