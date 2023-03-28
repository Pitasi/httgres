# httgres

How to run the HTTP server:

- Edit your postgres credentials in `main.go`
- `go run .`
- profit:

```
$ curl http://localhost:8080 -d'{"sql":"select * from \"Article\""}' | jq
{
  "columns": [
    {
      "name": "id",
      "type": "TEXT"
    },
    {
      "name": "createdAt",
      "type": "TIMESTAMP"
    },
    {
      "name": "updatedAt",
      "type": "TIMESTAMP"
    },
    {
      "name": "title",
      "type": "TEXT"
    },
    {
      "name": "content",
      "type": "TEXT"
    },
    {
      "name": "published",
      "type": "BOOL"
    },
    {
      "name": "authorId",
      "type": "TEXT"
    },
    {
      "name": "slug",
      "type": "TEXT"
    }
  ],
  "data": [
    [
      "clfl546fk0000yys3z7k7grt8",
      "2023-03-23T13:19:23.793Z",
      "2023-03-23T13:18:52.225Z",
      "Sample Article",
      "The body of the article",
      true,
      "clfjv00mh0000yycj17dzz2fn",
      "test-article"
    ]
  ]
}
```


# Rationale

There is some buzzing around edge computing (eg. Cloudflare workers). Code
running in such a runtime cannot directly access a database, it can however make
HTTP requests using `fetch()`.

I've done the tinyest proof of concept I could for the server-side that can talk
to any Postgres instance. Writing a Fetch-API client for such a server should be
somehow trivial.

Inspired by [database-js](https://github.com/planetscale/database-js) from
Planetscale.