package api

import (
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/muzzarellimj/grace-material-api/pkg/api/book/helper"
	api "github.com/muzzarellimj/grace-material-api/pkg/api/third_party/openlibrary.org"
	"github.com/muzzarellimj/grace-material-api/pkg/database"
	"github.com/muzzarellimj/grace-material-api/pkg/database/connection"
	"github.com/muzzarellimj/grace-material-api/pkg/database/service"
	model "github.com/muzzarellimj/grace-material-api/pkg/model/book"
	tmodel "github.com/muzzarellimj/grace-material-api/pkg/model/third_party/openlibrary.org"
	"github.com/muzzarellimj/grace-material-api/pkg/util"
)

func fetchBook(constraint string) (model.Book, error) {
	var connection connection.PgxPool = connection.Book

	bookFragment, err := service.FetchFragment[model.BookFragment](connection, database.TableBookFragments, constraint)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch book with constraint '%s': %v\n", constraint, err)

		return model.Book{}, err
	}

	authorFragmentSlice, err := fetchAuthorSlice(bookFragment)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch authors related to book '%d': %v\n", bookFragment.ID, err)
	}

	publisherFragmentSlice, err := fetchPublisherSlice(bookFragment)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch publishers related to book '%d': %v\n", bookFragment.ID, err)
	}

	topicFragmentSlice, err := fetchTopicSlice(bookFragment)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch topics related to book '%d': %v\n", bookFragment.ID, err)
	}

	return mapBook(bookFragment, authorFragmentSlice, publisherFragmentSlice, topicFragmentSlice), nil
}

func fetchAuthorSlice(bookFragment model.BookFragment) ([]model.BookAuthorFragment, error) {
	bookAuthorRelationshipSlice, err := service.FetchRelationshipSlice[model.BookAuthorRelationship](connection.Book, database.TableBookAuthorRelationships, fmt.Sprintf("book=%d", bookFragment.ID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relationships between book '%d' and authors: %v\n", bookFragment.ID, err)

		return []model.BookAuthorFragment{}, err
	}

	var authorFragmentSlice []model.BookAuthorFragment

	for _, relationship := range bookAuthorRelationshipSlice {
		authorFragment, err := service.FetchFragment[model.BookAuthorFragment](connection.Book, database.TableBookAuthorFragments, fmt.Sprintf("id=%d", relationship.Author))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch author '%d': %v\n", relationship.Author, err)
		}

		if authorFragment.ID != 0 {
			authorFragmentSlice = append(authorFragmentSlice, authorFragment)
		}
	}

	return authorFragmentSlice, nil
}

func fetchPublisherSlice(bookFragment model.BookFragment) ([]model.BookPublisherFragment, error) {
	bookPublisherRelationshipSlice, err := service.FetchRelationshipSlice[model.BookPublisherRelationship](connection.Book, database.TableBookPublisherRelationships, fmt.Sprintf("book=%d", bookFragment.ID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relationships between book '%d' and publishers: %v\n", bookFragment.ID, err)

		return []model.BookPublisherFragment{}, err
	}

	var publisherFragmentSlice []model.BookPublisherFragment

	for _, relationship := range bookPublisherRelationshipSlice {
		publisherFragment, err := service.FetchFragment[model.BookPublisherFragment](connection.Book, database.TableBookPublisherFragments, fmt.Sprintf("id=%d", relationship.Publisher))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch publisher '%d': %v\n", relationship.Publisher, err)
		}

		if publisherFragment.ID != 0 {
			publisherFragmentSlice = append(publisherFragmentSlice, publisherFragment)
		}
	}

	return publisherFragmentSlice, nil
}

func fetchTopicSlice(bookFragment model.BookFragment) ([]model.BookTopicFragment, error) {
	bookTopicRelationshipSlice, err := service.FetchRelationshipSlice[model.BookTopicRelationship](connection.Book, database.TableBookTopicRelationships, fmt.Sprintf("book=%d", bookFragment.ID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relationships between book '%d' and topics: %v\n", bookFragment.ID, err)

		return []model.BookTopicFragment{}, err
	}

	var topicFragmentSlice []model.BookTopicFragment

	for _, relationship := range bookTopicRelationshipSlice {
		topicFragment, err := service.FetchFragment[model.BookTopicFragment](connection.Book, database.TableBookTopicFragments, fmt.Sprintf("id=%d", relationship.Topic))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch topic '%d': %v\n", relationship.Topic, err)
		}

		if topicFragment.ID != 0 {
			topicFragmentSlice = append(topicFragmentSlice, topicFragment)
		}
	}

	return topicFragmentSlice, nil
}

func mapBook(bookFragment model.BookFragment, authorFragmentSlice []model.BookAuthorFragment, publisherFragmentSlice []model.BookPublisherFragment, topicFragmentSlice []model.BookTopicFragment) model.Book {
	return model.Book{
		ID:               bookFragment.ID,
		Title:            bookFragment.Title,
		Subtitle:         bookFragment.Subtitle,
		Description:      bookFragment.Description,
		Authors:          authorFragmentSlice,
		Publishers:       publisherFragmentSlice,
		Topics:           topicFragmentSlice,
		PublishDate:      bookFragment.PublishDate,
		Pages:            bookFragment.Pages,
		ISBN10:           bookFragment.ISBN10,
		ISBN13:           bookFragment.ISBN13,
		Image:            bookFragment.Image,
		EditionReference: bookFragment.EditionReference,
		WorkReference:    bookFragment.WorkReference,
	}
}

func storeBook(edition tmodel.OLEditionResponse, work tmodel.OLWorkResponse) (int, error) {
	bookId, err := storeBookFragment(edition, work)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to store book fragment: %v\n", err)

		return 0, err
	}

	authorIdSlice := storeAuthorFragments(edition.Authors)
	publisherIdSlice := storePublisherFragments(edition.Publishers)
	topicIdSlice := storeTopicFragments(work.Subjects)

	service.StoreRelationshipSlice(connection.Book, database.TableBookAuthorRelationships, database.PropertiesBookAuthorRelationships, service.RelationshipSliceArgument{
		SourceName:          "book",
		SourceArgument:      bookId,
		DestinationName:     "author",
		DestinationArgument: authorIdSlice,
	})

	service.StoreRelationshipSlice(connection.Book, database.TableBookPublisherRelationships, database.PropertiesBookPublisherRelationships, service.RelationshipSliceArgument{
		SourceName:          "book",
		SourceArgument:      bookId,
		DestinationName:     "publisher",
		DestinationArgument: publisherIdSlice,
	})

	service.StoreRelationshipSlice(connection.Book, database.TableBookTopicRelationships, database.PropertiesBookTopicRelationships, service.RelationshipSliceArgument{
		SourceName:          "book",
		SourceArgument:      bookId,
		DestinationName:     "topic",
		DestinationArgument: topicIdSlice,
	})

	return bookId, nil
}

func storeBookFragment(olEdition tmodel.OLEditionResponse, olWork tmodel.OLWorkResponse) (int, error) {
	var imageId int

	if len(olEdition.Images) > 0 {
		imageId = olEdition.Images[0]
	}

	bookId, err := service.StoreFragment(connection.Book, database.TableBookFragments, database.PropertiesBookFragments, pgx.NamedArgs{
		"title":             olEdition.Title,
		"subtitle":          olEdition.Subtitle,
		"description":       olWork.Description,
		"publish_date":      olEdition.PublishDate,
		"pages":             olEdition.Pages,
		"isbn10":            helper.ExtractISBN(olEdition.ISBN10),
		"isbn13":            helper.ExtractISBN(olEdition.ISBN13),
		"image":             fmt.Sprint(imageId),
		"edition_reference": olEdition.ID,
		"work_reference":    olWork.ID,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to store book fragment: %v\n", err)

		return 0, err
	}

	return bookId, nil
}

func storeAuthorFragments(authors []tmodel.OLResourceReference) []int {
	var authorIdSlice []int

	for _, author := range authors {
		existingAuthorFragment, err := service.FetchFragment[model.BookAuthorFragment](connection.Book, database.TableBookAuthorFragments, fmt.Sprintf("reference='%s'", author.ID))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch existing author fragment: %v\n", err)

			continue
		}

		if existingAuthorFragment.ID != 0 {
			authorIdSlice = append(authorIdSlice, existingAuthorFragment.ID)

			continue
		}

		olAuthor, err := api.OLGetAuthor(strings.ReplaceAll(author.ID, "/authors/", ""))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch OL author resource: %v\n", err)

			continue
		}

		first, middle, last := formatAuthorName(olAuthor.Name)

		var imageId int

		if len(olAuthor.Images) > 0 {
			imageId = olAuthor.Images[0]
		}

		storedAuthorId, err := service.StoreFragment(connection.Book, database.TableBookAuthorFragments, database.PropertiesBookAuthorFragments, pgx.NamedArgs{
			"first_name":  first,
			"middle_name": middle,
			"last_name":   last,
			"biography":   olAuthor.Biography,
			"image":       fmt.Sprint(imageId),
			"reference":   olAuthor.ID,
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to store new author fragment: %v\n", err)

			continue
		}

		if storedAuthorId != 0 {
			authorIdSlice = append(authorIdSlice, storedAuthorId)
		}
	}

	return authorIdSlice
}

func formatAuthorName(name string) (string, string, string) {
	nameSlice := strings.Split(name, " ")

	if len(nameSlice) == 1 {
		return nameSlice[0], "", ""
	}

	if len(nameSlice) == 2 {
		return nameSlice[0], "", nameSlice[1]
	}

	if len(nameSlice) == 3 {
		return nameSlice[0], nameSlice[1], nameSlice[2]
	}

	return name, "", ""
}

func storePublisherFragments(publishers []string) []int {
	var publisherIdSlice []int

	for _, publisher := range publishers {
		existingPublisherFragment, err := service.FetchFragment[model.BookPublisherFragment](connection.Book, database.TableBookPublisherFragments, fmt.Sprintf("name='%s'", publisher))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch existing publisher fragment: %v\n", err)

			continue
		}

		if existingPublisherFragment.ID != 0 {
			publisherIdSlice = append(publisherIdSlice, existingPublisherFragment.ID)

			continue
		}

		storedPublisherId, err := service.StoreFragment(connection.Book, database.TableBookPublisherFragments, database.PropertiesBookPublisherFragments, pgx.NamedArgs{
			"name": publisher,
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to store new publisher fragment: %v\n", err)

			continue
		}

		if storedPublisherId != 0 {
			publisherIdSlice = append(publisherIdSlice, storedPublisherId)
		}
	}

	return publisherIdSlice
}

func storeTopicFragments(topics []string) []int {
	var topicIdSlice []int

	for _, topic := range topics {
		existingTopicFragment, err := service.FetchFragment[model.BookTopicFragment](connection.Book, database.TableBookTopicFragments, fmt.Sprintf("name='%s'", util.FormatPSQLString(topic)))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch existing topic fragment: %v\n", err)

			continue
		}

		if existingTopicFragment.ID != 0 {
			topicIdSlice = append(topicIdSlice, existingTopicFragment.ID)

			continue
		}

		storedTopicId, err := service.StoreFragment(connection.Book, database.TableBookTopicFragments, database.PropertiesBookTopicFragments, pgx.NamedArgs{
			"name": topic,
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to store new topic fragment: %v\n", err)

			continue
		}

		if storedTopicId != 0 {
			topicIdSlice = append(topicIdSlice, storedTopicId)
		}
	}

	return topicIdSlice
}
