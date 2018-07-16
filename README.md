# gqlgen

[GraphQLを試してみた](https://github.com/cipepser/graphql-sample)けど、[gqlgen](https://github.com/vektah/gqlgen)が互換性がないのかイマイチ理解してきれていないので、[GoでGraphQLを話すサーバを作ってみた](https://qiita.com/ichikawa_0829/items/964682e3450d828968b9)をもとに触ってみる。

## gqlgenのインストール

2018/07/16現在の最新版で進める。
開発が盛んみたいなので、適宜vgoとかでバージョン管理するのがよいと思う。

```sh
go get -u github.com/vektah/gqlgen
```

## スキーマの定義

`./schema.graphql`で定義する。

```graphql
type User {
  id: ID!
  name: String!
}

type Query {
  user(id: String!): User
  users: [User!]!
}

input NewUser {
  name: String!
}

type Mutation {
  createUser(input: NewUser!): User!
}
```

`graph/graph.go`を以下のようにする。

```go
package graph

type User struct {
	ID   string
	Name string
}
```

記事だと`type.json`を使っているが、jsonフォーマットは廃止されるらしいので、`gqlgen.yml`にする。

```yml
schema: ../schema.graphql
exec:
  filename: generated.go
model:
  filename: models_gen.go
models:
  User:
    model: github.com/cipepser/gqlgen/graph.User
```

以下を実行すると、上記`yml`で定義した通りにgenerateされる。

```sh
❯ cd graph
❯ gqlgen graph/gqlgen.yml
```

## Resolverの実装

生成された`generated.go`で定義される以下のmethodを実装する。

```go
type Resolvers interface {
	Mutation_createUser(ctx context.Context, input NewUser) (User, error)
	Query_user(ctx context.Context, id string) (*User, error)
	Query_users(ctx context.Context) ([]User, error)
}
```

実装するのは`graph.go`

```go
// Resolver implements Resolvers interface
// type Resolvers interface {
// 	Mutation_createUser(ctx context.Context, input NewUser) (User, error)
// 	Query_user(ctx context.Context, id string) (*User, error)
// 	Query_users(ctx context.Context) ([]User, error)
// }
type Resolver struct {
	users []User
}

// Mutation_createUser creates a new user and add user to Resolver.
func (r *Resolver) Mutation_createUser(ctx context.Context, input NewUser) (User, error) {
	user := User{
		ID:   fmt.Sprintf("%d", rand.Int()),
		Name: input.Name,
	}
	r.users = append(r.users, user)
	return user, nil
}

// Query_user returns a user specified by the id.
func (r *Resolver) Query_user(ctx context.Context, id string) (*User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return &user, nil
		}
	}

	return nil, errors.New("user not found")
}

// Query_users returns all the users who Resolver knows.
func (r *Resolver) Query_users(ctx context.Context) ([]User, error) {
	return r.users, nil
}
```

## mainの実装

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cipepser/gqlgen/graph"
	"github.com/vektah/gqlgen/handler"
)

func main() {
	resolvers := &graph.Resolver{}
	http.Handle("/", handler.Playground("Todo", "/query"))
	http.Handle("/query", handler.GraphQL(graph.MakeExecutableSchema(resolvers)))

	fmt.Println("Listening on: 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

`go run main.go`で実行して`http://localhost:8080`にアクセスするとPlayGroundで遊べる。

```graphql
query {
  users {
    id
    name
  }
}
```

や

```graphql
query {
  user(id: "1") {
    name
  }
}
```

などとするとレスポンスが返ってくる（この段階ではまだuserがいないのでnot foundしか返ってこない）。

ユーザを登録するには以下。

```graphql
mutation ($NewUser: NewUser!) {
  createUser(input: $NewUser) {
    id,
    name
  }
}
```

このとき一緒に`QEURY VARIABLES`を以下のように設定する。

```json
{
  "NewUser": {
    "name": "hoge"
  }
}
```

上記を実行する or `users`のクエリを投げると`id`がわかるので、以下のようにするとちゃんと`name`が返ってくる。

```graphql
query {
  user(id: "5577006791947779410") {
    name
  }
}
```

```json
{
  "data": {
    "user": {
      "name": "hoge"
    }
  }
}
```

記事ではTODOリストを追加しているがここでは割愛。

## Rerefences
* [GoでGraphQLを話すサーバを作ってみた](https://qiita.com/ichikawa_0829/items/964682e3450d828968b9)
* [gqlgen](https://github.com/vektah/gqlgen)