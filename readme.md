# Test task Music library API

Implementation of online song library 🎶

The following needs to be implemented

1. Set rest methods
* Getting a data library with filtering by all fields and pagination
* Getting song lyrics with pagination by verses
* Deleting a song
* Changing song data
* Adding new songs in the format

JSON
```json
{
"group": "Muse",
"song": "Supermassive black hole"
}
```

2. When adding a request to the API, described by the swagger. The API, described by the swagger, will be raised against the background of the test task. No need to implement it separately

```yaml
openapi: 3.0.3
info:
title: Music Information
version: 0.0.1
paths:
/info:
get:
parameters:
  - name: group
  in: query
  required: true
  schema:
  type: string
  - name: song
  in : query
  required: true
  schema:
  type: string
  responses:
    '200':
      description: Ok
      content:
      application/json:
      schema:
      $ref: '#/components/schemas/SongDetail'
    '400':
      description: Bad request
    '500':
      description: Internal server error
      components:
      schemas:
      SongDetail:
      required:
      - releaseDate
      - text
      - link
      type: object
      properties:
      releaseDate:
      type: string
      example: 07/16/2006
      text:
      type: string
      example: Oh baby, don't you know I'm suffering?\nOh baby, can you hear me moaning?\nYou caught me under false pretenses\nHow long did it take before you let me go?\n\nOh\nYou set my soul on fire\nOh\ nYou set my soul on fire
      link:
      type: string
      example: https://www.youtube.com/watch?v=Xsp3_a-PMTw
```

3. Put the enriched information into the postgres DB (the DB structure should be created via migration when the service is started)
4. Cover the code with debug and info logs.
5. Move the configuration data to the .env file.
6. Generate a swagger on the implemented API.


To run app

`go run main.go`

To run db

`docker compose -f containers/docker-compose.postgres.yaml up -d`

Docs will be served on http://localhost:PORT/v1/swagger/index.html
