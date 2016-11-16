# Tumblr API Go Client

[![Build Status](https://travis-ci.org/tumblr/tumblr.go.svg?branch=master)](https://travis-ci.org/tumblr/tumblr.go) [![GoDoc](https://godoc.org/github.com/tumblr/tumblr.go?status.svg)](https://godoc.org/github.com/tumblr/tumblr.go)

This is the Tumblr API Golang client

## Installation

Run `go get github.com/tumblr/tumblr.go`

## Usage

The mechanics of this library send HTTP requests through a `ClientInterface` object. There is intentionally no concrete client defined in this library to allow for maximum flexibility. There is [a separate repository](https://github.com/tumblr/tumblrclient.go) with a client implementation and convenience methods if you do not require a custom client behavior.

## Support/Questions
You can post a question in the [Google Group](https://groups.google.com/forum/#!forum/tumblr-api) or contact the Tumblr API Team at [api@tumblr.com](mailto:api@tumblr.com)

## License

Copyright 2016 Tumblr, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
