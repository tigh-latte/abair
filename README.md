# Overview

A wrapper around the [Chi](https://github.com/go-chi/chi) web framework, adding strongly typed routing and responses in JSON, and a global error handler.

All while not breaking the `http.Handler` interface.

The end goal is for this to be swapped out with the std library http router, once that becomes available.

## Usage

First, define your request, response, and path param data types.

```go
type PersonPostBody struct {
    Name string `json:"name"`
    Age  int `json:"age"`
}

type Person struct {
    ID   string `json:"id"`
    Name string `json:"data"`
    Age  int    `json:"age"`
}

type PersonGetPath struct {
    ID string `path:"id"`
}
```

Define a `abair.HandlerFunc` for this data.

```go
func handlerPersonPost(ctx context.Context, req abair.Request[PersonPostBody, struct{}]) (Person, error) {
    person := Person{
        ID:   uuid.New().String(),
        Name: req.Body.Name,
        Age:  req.Body.Age,
    }

    insertIntoStore(person)

    return person
}

func handlerPersonGet(ctx context.Context, req abair.Request[struct{}, PersonGetPath]) (Person, error) {
    person := fetchFromStore(req.PathParams.ID)
    return person
}
```

Create a server via `NewService()` and use the page level functions to build routes.

Golang should be able to infer all types leading to simple generic hookups.

Start the server with `http.ListenAndServe(...)`

```go
server := abair.NewService()
abair.Post(server, "/person", handlerPersonPost)
abair.Get(server, "/person/{id}", handlerPersonGet)

http.ListenAndServe(":8080", server)
```
