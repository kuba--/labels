### Github labels color/description updater

The tool lets you update label's color and/or description for your repo.

```sh
GO111MODULE=on go build .
usage: ./labels <list | update> -r "<org/repo>" [options]
update options:
	-f string
		JSON file with labels (default "default.json")
	-a bool
		Add a new label if doesn't exist (default false)
```

##### List
```sh
GITHUB_TOKEN="my github token" ./labels list -r "<org/repo>"

# response
[
	{
		"name": "Please! ♥",
		"color": "d4c5f9",
		"description": "Particularly useful features that everyone would love!"
	},
	{
		"name": "bug",
		"color": "d73a4a"
	},
...
]
```

##### Update (add if doesn't exist)
```sh
GITHUB_TOKEN="my github token" ./labels update -r "<org/repo>" -f "mylabels.json" -a

# response
[invalid]: 200 OK
[duplicate]: 200 OK
[question]: 200 OK
[help wanted]: 200 OK
[wontfix]: 200 OK
[enhancement]: 200 OK
[bug]: 200 OK
[good first issue]: 201 Created
[meta]: 201 Created
[under discussion]: 201 Created
[change request]: 201 Created
[has PR]: 201 Created
[Please! ♥]: 201 Created
[enterprise]: 201 Created
[cleanup]: 201 Created
[do not merge yet]: 201 Created
[keyboard shortcuts]: 201 Created
[needs readme update]: 201 Created
```
