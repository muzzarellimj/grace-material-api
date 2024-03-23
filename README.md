# Grace Material API

The Grace Material API makes available an HTTP web server with a REST API, serving materials to the [Grace](https://github.com/muzzarellimj/grace) client applications. This API is written in [Go](https://go.dev/) and utilises the [Gin](https://gin-gonic.com/) web framework.

## Usage

### Practical

A user on the Grace client application uses the search bar to search for a game, "Super Smash Bros.", which does not yet exist in the Grace material database. As a search is made, results are displayed in a list with the cover image, title, and year published, and clicking one of these results will add it to the user's collection. The first option is clicked and the user receives a notification that the game is being added to their collection. Moments later, the user is rerouted to their collection overview page where the new game has been added successfully.

### Technical

After this repository has been cloned, the dependencies fetched, the .env properties added, and the web server started, an example request flow can begin with a search:

```
curl --request GET \
  --url 'http://localhost:8080/api/game/search?query=super%20smash%20bros'
```

... will garner response ...

```
[
	"status": 200,
	"data": [
		{
			"id": 1626,
			"name": "Super Smash Bros.",
			"first_release_date": 916876800,
			"cover": {
				"image_id": "co2tso"
			}
		},
		...
	]
]
```

Once a search has been made, a request can be made to either (a) store a new resource in the local database or (b) fetch an existing resource in the local database, all based on the source numeric identifier (e.g., IGDB ID):

```
curl --request POST \
  --url 'http://localhost:8080/api/game?id=1626'
```

... will garner response ...

```
{
	"message": "Game stored with numeric identifier '1'.",
	"status": 201
}
```

Now that the resource exists in the local database, local fetch requests can be made:

```
curl --request GET \
  --url 'http://localhost:8080/api/game?id=1'
```

... will garner response ...

```
{
	"data": {
		"id": 3,
		"title": "Super Smash Bros.",
		"summary": "Super Smash Bros. is a crossover fighting video game between several different Nintendo franchises, and the first installment in the Super Smash Bros. series. Players must defeat their opponents multiple times in a fighting frenzy of items and power-ups. Super Smash Bros. is a departure from the general genre of fighting games: instead of depleting an opponent's life bar, the players seek to knock opposing characters off a stage. Each player has a damage total, represented by a percentage, which rises as the damage is taken.",
		"storyline": "",
		"franchises": [
			{
				"id": 3,
				"name": "Pokémon",
				"reference": 60
			},
			...
			{
				"id": 13,
				"name": "Super Smash Bros.",
				"reference": 1787
			}
		],
		"genres": [
			{
				"id": 5,
				"name": "Fighting",
				"reference": 4
			},
			{
				"id": 6,
				"name": "Platform",
				"reference": 8
			}
		],
		"platforms": [
			{
				"id": 7,
				"name": "Nintendo 64",
				"reference": 4
			},
			{
				"id": 8,
				"name": "Wii",
				"reference": 5
			}
		],
		"studios": [
			{
				"id": 3,
				"name": "HAL Laboratory",
				"description": "HAL Laboratory is a Japanese game developer partnered with Nintendo. It was founded by a small group of friends who shared the desire to create games, and among said group was Satoru Iwata. The group developed games for numerous platforms at first before developing exclusively for Nintendo. The name HAL originates from the fictional computer HAL 9000 from 2001: A Space Odyssey, as well as each letter being one letter before IBM. HAL Laboratory is most notable for developing the Kirby series created by ex-employee Masahiro Sakurai, they are also responsible for the Super Smash Bros. and the EarthBound series.",
				"reference": 762
			}
		],
		"release_date": 916876800,
		"image": "https://images.igdb.com/igdb/image/upload/t_cover_big/co2tso.jpg",
		"reference": 1626
	},
	"status": 200
}
```