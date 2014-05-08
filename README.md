log-server
==========

[![Gobuild Download](http://gobuild.io/badge/github.com/jgoodall/log-server/download.png)](http://gobuild.io/github.com/jgoodall/log-server)

Simple server for saving and retrieving JSON logs via  HTTP interface. 

## Download

To download a binary, open a browser to [the gobuild page](http://gobuild.io/download/github.com/jgoodall/log-server).

## Usage 

To start the server:

    ./log-server -filepath="out.json" -logpath="log-server.log" -port=8000

There is no `OriginValidator`, so by default any host will be able to log to the server.

To save a log:

    curl -i -XPOST http://localhost:8000/log -d '{"message": "test message one", "type": "curl", "tags": ["firsttag", "secondtag"]}' -H "content-type: application/json"

    curl -i -XPOST http://localhost:8000/log -d '{"message": "test message two", "type": "curl", "tags": ["secondtag"]}' -H "content-type: application/json"

To retrieve logs:

    curl -i -XGET http://localhost:8000/logs -H "accept: application/json"


## License

This software is freely distributable under the terms of the MIT License.

Copyright (c) UT-Battelle, LLC (the "Original Author")

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS, THE U.S. GOVERNMENT, OR UT-BATTELLE BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.