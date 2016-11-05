package reference

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSymbolicReferencePathWithErrors(t *testing.T) {
	datas := [][]string{
		[]string{".test", `branch name must not start with a dot`, "Must not start with a dot"},
		[]string{"te..st", `branch name must not contains any double dot`, "Must not contains double dot"},
		[]string{"test..", `branch name must not contains any double dot`, "Must not contains double dot"},
		[]string{"test/", `branch name must not end with a slash`, "Must not end with a slash"},
		[]string{"test/~", `branch name must not end with a slash`, "Must not end with a slash"},
		[]string{"test/^", `branch name must not end with a slash`, "Must not end with a slash"},
		[]string{"test.lock", `branch name cannot end with .lock`, "Must not end with .lock"},
		[]string{"test.lock~", `branch name cannot end with .lock`, "Must not end with .lock"},
		[]string{"test.lock^", `branch name cannot end with .lock`, "Must not end with .lock"},
		[]string{"test.LOCK", `branch name cannot end with .lock`, "Must not end with .lock"},
		[]string{"te st", `branch name contains a space character`, "Must not contains a space"},
		[]string{"master~^22", `level associated to a caret must be 1 or 2`, "A number must follow a tilde or circumflex accent"},
		[]string{"master~tg^", `must be a caret or a tilde and may be followed with a number`, "Must be compound of tilde or circumflex accent"},
	}

	for _, data := range datas {
		parser := newParser(bytes.NewBufferString(data[0]))
		_, err := parser.parseSymbolicReferencePath()

		assert.EqualError(t, err, data[1], data[2]+" : "+data[0])
	}
}

func TestParseSymbolicReferencePath(t *testing.T) {
	datas := []string{
		"master",
		"master~^",
		"master~3^2~2^1^~5^^~10",
	}

	expected := []symbolicRefPathStmt{
		symbolicRefPathStmt{
			"master",
			[]int{},
		},
		symbolicRefPathStmt{
			"master",
			[]int{1, 1},
		},
		symbolicRefPathStmt{
			"master",
			[]int{1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
	}

	for i := 0; i < len(datas); i++ {
		parser := newParser(bytes.NewBufferString(datas[i]))
		stmt, err := parser.parseSymbolicReferencePath()

		assert.NoError(t, err, "Must return no errors")
		assert.Equal(t, expected[i], *stmt, "Must return statement")
	}
}
