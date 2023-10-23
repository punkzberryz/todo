# todo

Todo list with Go

### deploy

- Create the todo-network
  If network hasn't been existed

  ```bash
  make network
  ```

- Compoase docker up with re-build
  ```bash
  make docker-up-rebuild
  ```

### Setup infrastructure

- Create the todo-network

  ```bash
  make network
  ```

- Start postgres container:
  ```bash
  make postgres
  ```

### How to generate code

- Generate SQL CRUD with sqlc:

  ```bash
  make sqlc
  ```

- Generate DB mock with gomock:

  ```bash
  make mock
  ```

- Create a new db migration:
  ```bash
  migrate create -ext sql -dir db/migration -seq  <migration_name>
  ```
