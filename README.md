# contentful-go
Contentful api client for GoLang (UNOFFICIAL)

## Install

`go get -u github.com/tolgaakyuz/contentful-go`

## Getting Started

Setting up SDK client is very easy and it only requires passing a simple settings to sdk. The following example is create a sdk instance and it is enough to have to get benefit of all the sdk functions.

```go
token := "your-cma-token"
cma := contentful.NewCMA(token)
```

### Setting Organization
When you set `organizationID` for contentful sdk client, every api request will have `X-Contentful-Organization: organizationID` header.

```go
token := "your-cma-token"
organizationID := "your-organization-id"

cma := contentful.NewCMA(token)
cma.SetOrganization(organizationID)
```

### Debug mode

When activated, sdk client starts to work in verbose mode and try to print as much informatin as possible. In debug mode, all outgoing http requests are printed nicely in the form of `curl` command so that you can easly drop into your command line to debug specific request.

```go
c.Debug = true
```

## Concepts

- [Querying](#querying) - Explains how to issue query requests
- [Collections](#collections) - General rules for multi resource endpoints, paginating,  type casting etc.

## Querying

### Geting/Querying entries

Please check [query.go](https://github.com/tolgaakyuz/contentful.go/blob/master/contentful/query.go) to see all the query options

```go
query := c.Space.Entries()
query.Query
    .ContentType("content_type_id")
    .Select([]string{"field1, field2"})
    .Equal("fields.field3", "cat")
    .NotEqual("fields.field3", "dog")
    .Exists("fields.field4")

entries, err := query.All()
if err != nil {
  fmt.Println(err)
}
```

### Collections

All the endpoints which return multiple entity objects are wrapped around `Collection` struct. The main goal of `Collection` is to give you the ability to cast api results into entity objects such as `Space`. `Collection` struct exposes the necessary converters such as `toSpace()`. The following example gets all spaces for the given account:

```go
col, err := c.GetSpaces().Next()
if err != nil {
  log.Fatal(err)
}
spaces := col.ToSpace()
for _, space := range spaces {
  fmt.Println(space.Name)
  fmt.Println(space.Sys.ID) // or space.ID()
}

// In order to access collection metadata
fmt.Println(col.Total)
fmt.Println(col.Skip)
fmt.Println(col.Limit)
```

You can also call `Collection.Next()` to paginate or nagivate through the collection:

```go
col := c.GetSpaces()

while col.HasMore() {
  _, err := col.Next()
  if err != nil {
    break
  }

  spaces := col.ToSpace() // next 100 space is populated
}
```

In order to change the pagination limit of collection, you can initialize `Collection` struct with a option parameter:

```go
col, err := c.GetSpaces(&CollectionOption{
  Limit; 60,
})

// now, everytime `col.Next()` is called, it will fetch the next 60 object from api
```