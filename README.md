log-server
==========

[![Gobuild Download](http://gobuild.io/badge/github.com/jgoodall/log-server/download.png)](http://gobuild.io/github.com/jgoodall/log-server)

Simple server for saving and retrieving JSON logs via  HTTP interface. 

## Download

To download a binary, open a browser to [the gobuild page](http://gobuild.io/download/github.com/jgoodall/log-server).

## Usage 

To start the server:

    ./log-server -filepath="out.json" -logpath="log-server.log" -port=8000

To save a log:

    curl -i -XPOST http://localhost:8000/log -d '{"message": "test message one", "type": "curl", "tags": ["firsttag", "secondtag"]}' -H "content-type: application/json"

    curl -i -XPOST http://localhost:8000/log -d '{"message": "test message two", "type": "curl", "tags": ["secondtag"]}' -H "content-type: application/json"

To retrieve logs:

    curl -i -XGET http://localhost:8000/logs -H "accept: application/json"


## License