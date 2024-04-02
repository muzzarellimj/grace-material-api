package helper

import (
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	api "github.com/muzzarellimj/grace-material-api/internal/api/third_party/openlibrary.org"
	"github.com/muzzarellimj/grace-material-api/internal/database"
	"github.com/muzzarellimj/grace-material-api/internal/database/service"
	model "github.com/muzzarellimj/grace-material-api/internal/model/book"
	OLModel "github.com/muzzarellimj/grace-material-api/internal/model/third_party/openlibrary.org"
	"github.com/muzzarellimj/grace-material-api/internal/util"
)

func ProcessBookStorage(edition OLModel.OLEditionResponse, work OLModel.OLWorkResponse) (int, error) {
	bookId, err := storeBookFragment(edition, work)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to store book '%s' fragment: %v\n", ExtractResourceId(edition.ID), err)

		return 0, err
	}

	authorIdSlice := processAuthorFragmentSliceStorage(edition.Authors)
	publisherIdSlice := processPublisherFragmentSliceStorage(edition.Publishers)
	topicIdSlice := processTopicFragmentSliceStorage(work.Subjects)

	service.StoreRelationshipSlice(database.Connection, database.TableBookAuthorRelationships, database.PropertiesBookAuthorRelationships, service.RelationshipSliceArgument{
		SourceName:          "book",
		SourceArgument:      bookId,
		DestinationName:     "author",
		DestinationArgument: authorIdSlice,
	})

	service.StoreRelationshipSlice(database.Connection, database.TableBookPublisherRelationships, database.PropertiesBookPublisherRelationships, service.RelationshipSliceArgument{
		SourceName:          "book",
		SourceArgument:      bookId,
		DestinationName:     "publisher",
		DestinationArgument: publisherIdSlice,
	})

	service.StoreRelationshipSlice(database.Connection, database.TableBookTopicRelationships, database.PropertiesBookTopicRelationships, service.RelationshipSliceArgument{
		SourceName:          "book",
		SourceArgument:      bookId,
		DestinationName:     "topic",
		DestinationArgument: topicIdSlice,
	})

	return bookId, nil
}

func storeBookFragment(edition OLModel.OLEditionResponse, work OLModel.OLWorkResponse) (int, error) {
	bookId, err := service.StoreFragment(database.Connection, database.TableBookFragments, database.PropertiesBookFragments, pgx.NamedArgs{
		"title":             edition.Title,
		"subtitle":          edition.Subtitle,
		"description":       ExtractDescription(work.Description),
		"publish_date":      util.ParseDateTime(edition.PublishDate),
		"pages":             edition.Pages,
		"isbn10":            ExtractISBN(edition.ISBN10),
		"isbn13":            ExtractISBN(edition.ISBN13),
		"image":             fmt.Sprintf("https://covers.openlibrary.org/b/olid/%s-L.jpg", ExtractResourceId(edition.ID)),
		"edition_reference": ExtractResourceId(edition.ID),
		"work_reference":    ExtractResourceId(work.ID),
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to store book fragment '%s': %v\n", ExtractResourceId(edition.ID), err)

		return 0, err
	}

	return bookId, nil
}

func processAuthorFragmentSliceStorage(authors []OLModel.OLResourceReference) []int {
	var authorIdSlice []int

	for _, resource := range authors {
		existingAuthorFragment, err := service.FetchFragment[model.BookAuthorFragment](database.Connection, database.TableBookAuthorFragments, fmt.Sprintf("reference='%s'", ExtractResourceId(resource.ID)))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch existing author '%s' fragment: %v\n", ExtractResourceId(resource.ID), err)

			continue
		}

		if existingAuthorFragment.ID != 0 {
			authorIdSlice = append(authorIdSlice, existingAuthorFragment.ID)

			continue
		}

		author, err := api.OLGetAuthor(ExtractResourceId(resource.ID))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch author '%s' OL record: %v\n", ExtractResourceId(resource.ID), err)

			continue
		}

		firstName, middleName, lastName := ExtractName(author.Name)

		authorId, err := service.StoreFragment(database.Connection, database.TableBookAuthorFragments, database.PropertiesBookAuthorFragments, pgx.NamedArgs{
			"first_name":  firstName,
			"middle_name": middleName,
			"last_name":   lastName,
			"biography":   author.Biography,
			"image":       fmt.Sprintf("https://covers.openlibrary.org/a/olid/%s-L.jpg", ExtractResourceId(author.ID)),
			"reference":   ExtractResourceId(author.ID),
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to store new author '%s' fragment: %v\n", ExtractResourceId(author.ID), err)

			continue
		}

		if authorId != 0 {
			authorIdSlice = append(authorIdSlice, authorId)
		}
	}

	return authorIdSlice
}

func processPublisherFragmentSliceStorage(publishers []string) []int {
	var publisherIdSlice []int

	for _, publisher := range publishers {
		existingPublisherFragment, err := service.FetchFragment[model.BookPublisherFragment](database.Connection, database.TableBookPublisherFragments, fmt.Sprintf("name='%s'", util.FormatPSQLString(publisher)))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch existing publisher '%s' fragment: %v\n", publisher, err)

			continue
		}

		if existingPublisherFragment.ID != 0 {
			publisherIdSlice = append(publisherIdSlice, existingPublisherFragment.ID)

			continue
		}

		publisherId, err := service.StoreFragment(database.Connection, database.TableBookPublisherFragments, database.PropertiesBookPublisherFragments, pgx.NamedArgs{
			"name": publisher,
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to store new publisher '%s' fragment: %v\n", publisher, err)

			continue
		}

		if publisherId != 0 {
			publisherIdSlice = append(publisherIdSlice, publisherId)
		}
	}

	return publisherIdSlice
}

func processTopicFragmentSliceStorage(topics []string) []int {
	var topicIdSlice []int

	for _, topic := range topics {
		existingTopicFragment, err := service.FetchFragment[model.BookTopicFragment](database.Connection, database.TableBookTopicFragments, fmt.Sprintf("name='%s'", util.FormatPSQLString(topic)))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch existing topic '%s' fragment: %v\n", topic, err)

			continue
		}

		if existingTopicFragment.ID != 0 {
			topicIdSlice = append(topicIdSlice, existingTopicFragment.ID)

			continue
		}

		topicId, err := service.StoreFragment(database.Connection, database.TableBookTopicFragments, database.PropertiesBookTopicFragments, pgx.NamedArgs{
			"name": topic,
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to store new topic '%s' fragment: %v\n", topic, err)

			continue
		}

		if topicId != 0 {
			topicIdSlice = append(topicIdSlice, topicId)
		}
	}

	return topicIdSlice
}
