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


## Rerefences
* [GoでGraphQLを話すサーバを作ってみた](https://qiita.com/ichikawa_0829/items/964682e3450d828968b9)
* [gqlgen](https://github.com/vektah/gqlgen)