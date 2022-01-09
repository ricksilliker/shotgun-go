# Shotgun Go API

A small, somewhat opinionated library to interact with a Shotgun server with Go. What makes it opinionated
is the return fields for each Shotgun entity. This API was written in a way to return useful info,
related to a previous place I worked.


## Basic Usage Idea

** Environment expectations **

The following env vars are expected.

- `SHOTGUN_URL` - The base url of your Shotgun server.

** Finding an entity from a known ID**
```go
req, _ := NewFindRequest("Shot", shotID, shotFields)
var resp ShotRecordResponse
if err = DoFindRequest(req, &resp); err != nil {
    logrus.Error("failed to make Shot find request")
    return nil, err
}
// Parse ShotRecordResponse struct for your needs.
```

** Creating a custom filter for a search request**
```go
filters := ShotgunFilters{
    Expressions: []ShotgunFilterExpression{
        {"sg_sequence.Sequence.id", "is", sequenceID},
    },
}

req, _ := NewSearchRequest("Shot", filters, shotFields, nil, sortBy)

var resp ShotMultiRecordResponse
if err = DoSearchRequest(req, &resp); err != nil {
    logrus.Error("failed to make Shot search request")
    return nil, err
}
// Same as above, use the response struct for your needs, it contains the Shotgun records.
```

