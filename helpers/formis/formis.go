//
package formis

import (
	"html/template"
)

var (
	Functions = map[string]interface{}{
		"form_is": FormisCreate,
	}
)
