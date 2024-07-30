package smartquery

import (
	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = Describe("Test mandatory SQL SELECT", func() {

	Context("Test mandatory SQL SELECT", func() {
		It("correctly generates the mandatory SQL SELECT", func() {
			squirrelizer := Squirrelizer{}
			_, err := squirrelizer.BuildSelect("playlist1", "user666", "title")
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			expectedSQL := "SELECT row_number() over (order by title) as id, 'playlist1' as playlist_id, media_file.id as media_file_id " +
				"FROM media_file " +
				"LEFT JOIN annotation on (annotation.item_id = media_file.id AND annotation.item_type = 'media_file' AND annotation.user_id = 'user666') WHERE "
			// fmt.Println("++++++++++++++++++++++++++++")
			// fmt.Printf("(1) %d %v", len([]byte(expectedSQL)), []byte(expectedSQL))
			// fmt.Println()
			// fmt.Println("++++++++++++++++++++++++++++")
			// fmt.Printf("(2) %d %v", len([]byte(squirrelizer.Sql)), []byte(squirrelizer.Sql))
			// fmt.Println()
			// fmt.Println("++++++++++++++++++++++++++++")
			gomega.Expect(squirrelizer.Sql).To(gomega.Equal(expectedSQL))
		})
	})
})

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

			squirrelizer := Squirrelizer{}
			_, err = squirrelizer.BuildRefreshSmartQueryPlaylistSQL("playlist1", "user666", smartQuery.Query, "title")
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

			expectedSQL := "INSERT INTO playlist_tracks (id,playlist_id,media_file_id) " +
				"SELECT row_number() over (order by title) as id, 'playlist1' as playlist_id, media_file.id as media_file_id " +
				"FROM media_file " +
				"LEFT JOIN annotation on (annotation.item_id = media_file.id AND annotation.item_type = 'media_file' AND annotation.user_id = 'user666') " +
				"WHERE title LIKE \"John*\""

			// fmt.Println("++++++++++++++++++++++++++++")
			// fmt.Printf("(1) [%v] %d %v", expectedSQL, len([]byte(expectedSQL)), []byte(expectedSQL))
			// fmt.Println()
			// fmt.Println("++++++++++++++++++++++++++++")
			// fmt.Printf("(2) [%v]  %d %v", squirrelizer.Sql, len([]byte(squirrelizer.Sql)), []byte(squirrelizer.Sql))
			// fmt.Println()
			// fmt.Println("++++++++++++++++++++++++++++")

			gomega.Expect(squirrelizer.Sql).To(gomega.Equal(expectedSQL))
		})
	})
})
