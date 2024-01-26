# Grace Material API

The Grace Material API makes available an HTTP web server with a REST API, serving materials to the [Grace](https://github.com/muzzarellimj/grace) client applications. This API is written in [Go](https://go.dev/) and utilises the [Gin](https://gin-gonic.com/) web framework.

## Examples

After this repository has been cloned, the dependencies fetched, and the web server started, a few prepackaged example requests can be made. For example...

```
curl --request GET \
  --url http://localhost:8080/api/ex/movies
```

... should garner response ...

```
[
	{
		"id": 568124,
		"title": "Encanto",
		"overview": "The tale of an extraordinary family, the Madrigals, who live hidden in the mountains of Colombia, in a magical house, in a vibrant town, in a wondrous, charmed place called an Encanto. The magic of the Encanto has blessed every child in the family—every child except one, Mirabel. But when she discovers that the magic surrounding the Encanto is in danger, Mirabel decides that she, the only ordinary Madrigal, might just be her exceptional family's last hope."
	},
	{
		"id": 38757,
		"title": "Tangled",
		"overview": "When the kingdom's most wanted-and most charming-bandit Flynn Rider hides out in a mysterious tower, he's taken hostage by Rapunzel, a beautiful and feisty tower-bound teen with 70 feet of magical, golden hair. Flynn's curious captor, who's looking for her ticket out of the tower where she's been locked away for years, strikes a deal with the handsome thief and the unlikely duo sets off on an action-packed escapade, complete with a super-cop horse, an over-protective chameleon and a gruff gang of pub thugs."
	}
]
```

Additional example requests can be found under the `/api/ex/...` endpoints and are handled by functions found in the [example/](/example/) subdirectory.