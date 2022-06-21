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

- Run PostgreSQL psql from docker:

  ```bash
  docker exec -it postgres12 psql -U root -d <db_name>
  ```
