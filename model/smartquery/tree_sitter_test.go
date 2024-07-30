package smartquery

import (
	"strings"
	"testing"
)

func TestParseSql(t *testing.T) {
	smartSQL := "title LIKE \"John*\""
	smartQuery := SmartQuery{"Songs about John", "Any song starting with John", smartSQL, ""}
	err := smartQuery.ValidateQuery()
	if err != nil {
		t.Fatalf("got error %v", err)
	}
	expectedOrderBy := ""
	if expectedOrderBy != smartQuery.OrderBy {
		t.Fatalf("expected %q , got %#q", expectedOrderBy, smartQuery.OrderBy)
	}
}

func TestParseSqlWithSpecificOrderBy(t *testing.T) {
	// EXTRA_DEBUG = true
	smartSQL := strings.Join([]string{
		"-- comment ",
		"album_id = ",
		"(SELECT id FROM album ORDER BY random() /* inline comment */ LIMIT 1)",
		"-- another comment ",
		"ORDER BY disc_number, track_number ASC"}, "\n")
	smartQuery := SmartQuery{"Random Album", "Pick a Random Album", smartSQL, ""}
	err := smartQuery.ValidateQuery()
	if err != nil {
		t.Fatalf("got error %v", err)
	}
	expected := "disc_number, track_number ASC"
	if expected != smartQuery.OrderBy {
		t.Fatalf("expected %q , got %#q", expected, smartQuery.OrderBy)
	}
}
