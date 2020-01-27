[![Go Report Card](https://goreportcard.com/badge/github.com/coby9241/frontend-service)](https://goreportcard.com/report/github.com/coby9241/frontend-service)
[![Build Status](https://travis-ci.org/coby9241/frontend-service.svg?branch=master)](https://travis-ci.org/coby9241/frontend-service)

# Frontend Service

This is a frontend service consisting of an admin page to manage certain resources and is built with QOR Admin.

## Development

To compile QOR Admin assets and login template using bindatafs:
```
 go run cmd/compile/main.go
```

To run locally:
```bash
go run -tags 'bindatafs' main.go
```
