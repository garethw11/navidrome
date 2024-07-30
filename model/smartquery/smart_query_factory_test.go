package smartquery

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = Describe("test Smart Query Factory", func() {

	Context("Simple SELECT", func() {
		It("correctly parses the smart playlist definition", func() {
			smartplaylistDefinition := []string{"PLAYLIST name: Random Album, description: smart playlist",
				"-- comment",
				"album_id=(SELECT id FROM album ORDER BY random() /* inline comment */ LIMIT 1)",
				"-- another comment",
				"ORDER BY disc_number, track_number ASC"}

			factory := SmartQueryFactory{smartplaylistDefinition}
			smartQuery, err := factory.CreateSmartQuery()
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			validationErr := smartQuery.ValidateQuery()
			gomega.Expect(validationErr).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(smartQuery.Name).To(gomega.Equal("Random Album"))
			gomega.Expect(smartQuery.Comment).To(gomega.Equal("smart playlist"))
			expectedSQL := strings.Join([]string{
				"-- comment",
				"album_id=(SELECT id FROM album ORDER BY random() /* inline comment */ LIMIT 1)",
				"-- another comment",
				"ORDER BY disc_number, track_number ASC"}, "\n")
			gomega.Expect(smartQuery.Query).To(gomega.Equal(expectedSQL))
			gomega.Expect(smartQuery.OrderBy).To(gomega.Equal("disc_number, track_number ASC"))
		})
	})
})
