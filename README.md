# A Simple Banking App

golang_simplebank is a simple rest api application developed using Golang. It allows users to register an account, deposit to their account and transfer money across accounts.

### Project Requirement

1. Go enviroment - Go >= 1.18
2. SQLc
3. Gin Golang Framework
4. Gorm
5. Paseto/JWT
6. testify
7. gomock
8. Docker
9. Viper

### How to generate code

- Generate SQL CRUD using sqlc:

  ```bash
  make sqlc
  ```

- Generate DB mock using gomock:

  ```bash
  make mock
  ```

- Create a new db migration:

  ```bash
  migrate create -ext sql -dir db/migration -seq <migration_name>
  ```

### How to run

- Run server:

  ```bash
  make server
  ```

## Getting Setup

- Run PostgreSQL psql from docker:

  ```bash
  docker exec -it postgres12 psql -U root -d <db_name>
  ```

- Build app into docker image:

  ```bash
  docker build -t <image-tag> .
  ```

- Run our custom app docker image:

  ```bash
  docker run --name <image_name> -p 8080:8080 e GIN_MODE=release <image-name>:<image-tag>
  ```

  ### Installing useful tools

#### 1. [Postbird](https://github.com/paxa/postbird)

Postbird is a useful client GUI (graphical user interface) to interact with our provisioned Postgres database. We can establish a remote connection and complete actions like viewing data and changing schema (tables, columns, ect).

#### 2. [Postman](https://www.getpostman.com/downloads/)

Postman is a useful tool to issue and save requests. Postman can create GET, PUT, POST, etc. requests complete with bodies. It can also be used to test endpoints automatically.

---
