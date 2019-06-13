package github

import (
	"strings"

	"github.com/src-d/go-mysql-server/sql"
)

const labelComponentDelim = "::"

func toType(v string) (label string, t sql.Type) {
	p := strings.Split(v, labelComponentDelim)
	// label has no second component
	if len(p) == 0 {
		return v, sql.Boolean
	}
	// label has second component
	if len(p) > 1 && len(p[1]) > 0 {
		if p[1][0] == '{' {
			return p[0], sql.JSON
		}
	}
	// label is plain text content
	return p[0], sql.Text
}
