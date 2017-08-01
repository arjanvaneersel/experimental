PACKAGE

package mime
import "github.com/arjanvaneersel/gockan/mime"

CKAN format tags to mime types.


VARIABLES

var MimeTypes map[string]MimeType


TYPES

type MimeType struct {
    // The human readable friendly label
    Label string
    // The actual mime type
    Value string
}

