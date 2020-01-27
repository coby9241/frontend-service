[![Go Report Card](https://goreportcard.com/badge/github.com/coby9241/frontend-service)](https://goreportcard.com/report/github.com/coby9241/frontend-service)
[![Build Status](https://travis-ci.org/coby9241/frontend-service.svg?branch=master)](https://travis-ci.org/coby9241/frontend-service)
[![Maintainability](https://api.codeclimate.com/v1/badges/15063695ed48e8287dc6/maintainability)](https://codeclimate.com/github/coby9241/frontend-service/maintainability)

# Frontend Service

This is a frontend service consisting of an admin page to manage certain resources and is built with QOR Admin.

## Requirements

- [Docker](https://github.com/docker/docker-ce)
- [Docker Compose](https://github.com/docker/compose)

## Development

To run locally:
```bash
docker-compose build && docker-compose up
```

Then open the Web UI by entering in the url: `http://localhost:8082/admin`.
