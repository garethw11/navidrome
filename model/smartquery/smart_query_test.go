package smartquery

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = Describe("Smart SQL Query", func() {

	Context("Simple SELECT", func() {
		It("correctly parses the SQL", func() {
			smartSQL := "title LIKE \"John*\""
			smartQuery := SmartQuery{"Songs about John", "Any song starting with John", smartSQL, ""}
			err := smartQuery.ValidateQuery()
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(smartQuery.Name).To(gomega.Equal("Songs about John"))
			gomega.Expect(smartQuery.Comment).To(gomega.Equal("Any song starting with John"))
			gomega.Expect(smartQuery.Query).To(gomega.Equal(smartSQL))
			gomega.Expect(smartQuery.OrderBy).To(gomega.Equal(""))
		})
	})

	Context("SELECT with subquery", func() {
		It("correctly parses the SQL", func() {
			smartSQL := strings.Join([]string{
				"-- comment ",
				"album_id = ",
				"(SELECT id FROM album ORDER BY random() /* inline comment */ LIMIT 1)",
				"-- another comment ",
				"ORDER BY disc_number, track_number ASC"}, "\n")
			smartQuery := SmartQuery{"Random Album", "Pick a Random Album", smartSQL, ""}
			err := smartQuery.ValidateQuery()
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(smartQuery.Name).To(gomega.Equal("Random Album"))
			gomega.Expect(smartQuery.Comment).To(gomega.Equal("Pick a Random Album"))
			gomega.Expect(smartQuery.Query).To(gomega.Equal(smartSQL))
			gomega.Expect(smartQuery.OrderBy).To(gomega.Equal("disc_number, track_number ASC"))
		})
	})

	Context("SELECT with comment at end of line", func() {
		It("correctly parses the SQL", func() {
			smartSQL := strings.Join([]string{
				"album_id = 'foo' -- comment",
				"AND disc_number > 1"}, "\n")
			smartQuery := SmartQuery{"foo", "bar", smartSQL, ""}
			err := smartQuery.ValidateQuery()
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(smartQuery.Query).To(gomega.Equal(smartSQL))
		})
	})

	Context("SELECT with comment line", func() {
		It("correctly parses the SQL", func() {
			smartSQL := strings.Join([]string{
				"album_id = 'foo'",
				"      -- comment",
				"AND disc_number > 1"}, "\n")
			smartQuery := SmartQuery{"foo", "bar", smartSQL, ""}
			err := smartQuery.ValidateQuery()
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(smartQuery.Query).To(gomega.Equal(smartSQL))
		})
	})

	Context("SELECT with inline comments", func() {
		It("correctly parses the SQL", func() {
			smartSQL := "album_id = /* inline comment */ 'foo'"
			smartQuery := SmartQuery{"foo", "bar", smartSQL, ""}
			err := smartQuery.ValidateQuery()
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(smartQuery.Query).To(gomega.Equal(smartSQL))
		})
	})

	Context("INSERT", func() {
		It("rejects the SQL", func() {
			smartSQL := "INSERT foo INTO bar"
			smartQuery := SmartQuery{"foo", "bar", smartSQL, ""}
			err := smartQuery.ValidateQuery()
			gomega.Expect(err).Should(gomega.HaveOccurred())
			// (todo) I can't figure this out gomega.Expect(err).To(gomega.Equal("sql parse failed INSERT foo INTO"))

		})
	})

	Context("DELETE", func() {
		It("rejects the SQL", func() {
			smartSQL := "DELETE FROM bar"
			smartQuery := SmartQuery{"foo", "bar", smartSQL, ""}
			err := smartQuery.ValidateQuery()
			gomega.Expect(err).Should(gomega.HaveOccurred())
		})
	})

	Context("UPDATE", func() {
		It("rejects the SQL", func() {
			smartSQL := "UPDATE foo SET bar VALUES (\"asd\")"
			smartQuery := SmartQuery{"foo", "bar", smartSQL, ""}
			err := smartQuery.ValidateQuery()
			gomega.Expect(err).Should(gomega.HaveOccurred())
		})
	})

	Context("EXPLAIN", func() {
		It("rejects the SQL", func() {
			smartSQL := strings.Join([]string{
				"album_id = ",
				"(EXPLAIN SELECT id FROM album ORDER BY random() LIMIT 10)"}, "\n")
			smartQuery := SmartQuery{"foo", "bar", smartSQL, ""}
			err := smartQuery.ValidateQuery()
			gomega.Expect(err).Should(gomega.HaveOccurred())
		})
	})

	Context("EXPLAIN QUERY PLAN", func() {
		It("rejects the SQL", func() {
			smartSQL := strings.Join([]string{
				"album_id = ",
				"(EXPLAIN QUERY PLAN SELECT id FROM album ORDER BY random() LIMIT 10)"}, "\n")
			smartQuery := SmartQuery{"foo", "bar", smartSQL, ""}
			err := smartQuery.ValidateQuery()
			gomega.Expect(err).Should(gomega.HaveOccurred())
		})
	})

	Context("Marshall JSON", func() {
		It("turns into JSON", func() {
			smartSQL := strings.Join([]string{
				"album_id = (SELECT id FROM album ORDER BY random() LIMIT 1)",
				"ORDER BY disc_number, track_number ASC"}, "\n")
			smartQuery := SmartQuery{"Random Album", "Pick a Random Album", smartSQL, ""}
			err := smartQuery.ValidateQuery()
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(smartQuery.Name).To(gomega.Equal("Random Album"))
			gomega.Expect(smartQuery.Comment).To(gomega.Equal("Pick a Random Album"))
			gomega.Expect(smartQuery.Query).To(gomega.Equal(smartSQL))
			gomega.Expect(smartQuery.OrderBy).To(gomega.Equal("disc_number, track_number ASC"))
			json, err := smartQuery.MarshalJSON()
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			expected := "{\"name\":\"Random Album\",\"comment\":\"Pick a Random Album\"," +
				"\"query\":\"album_id = (SELECT id FROM album ORDER BY random() LIMIT 1)\\n" +
				"ORDER BY disc_number, track_number ASC\"," +
				"\"orderby\":\"disc_number, track_number ASC\"}"
			got := string(json)
			gomega.Expect(got).To(gomega.Equal(expected))
		})
	})

	Context("Unmarshall JSON", func() {
		It("eats JSON", func() {
			json := "{\"name\":\"Random Album\",\"comment\":\"Pick a Random Album\"," +
				"\"query\":\"album_id = (SELECT id FROM album ORDER BY random() LIMIT 1)\\n" +
				"ORDER BY disc_number, track_number ASC\"," +
				"\"orderby\":\"disc_number, track_number ASC\"}"

			squery := SmartQuery{}
			smartQuery, err := squery.UnmarshalJSON([]byte(json))
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			validationErr := smartQuery.ValidateQuery()
			gomega.Expect(validationErr).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(smartQuery.Name).To(gomega.Equal("Random Album"))
			gomega.Expect(smartQuery.Comment).To(gomega.Equal("Pick a Random Album"))
			expectedSQL := strings.Join([]string{"album_id = (SELECT id FROM album ORDER BY random() LIMIT 1)",
				"ORDER BY disc_number, track_number ASC"}, "\n")
			gomega.Expect(smartQuery.Query).To(gomega.Equal(expectedSQL))
			gomega.Expect(smartQuery.OrderBy).To(gomega.Equal("disc_number, track_number ASC"))
		})
	})
})
