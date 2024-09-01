# structogqlgen
StructsToGqlGenTypes is a tool that helps to automatically converts Golang structs into GraphQL types that are readily usable with the popular GraphQL framework, gqlgen. It aims to reduce the boilerplate code required to define GraphQL schemas manually, thus accelerating the development of GraphQL APIs in Go projects.

# Features
- Automatic Conversion: Convert your existing Golang structs into GraphQL types with a single command, preserving struct relationships and data types.
- Compatibility with gqlgen: Generated GraphQL types are fully compatible with gqlgen, ensuring smooth integration into your existing GraphQL server setup.
- Customizable Mappings: Offers options to customize how your Go data types are mapped to GraphQL types, allowing for fine-tuned control over the schema generation process.
- CLI Tooling: Comes with a command-line interface that makes it easy to integrate into development workflows and CI/CD pipelines.

# Use Case
  This tool is perfect for developers using gqlgen who already have a collection of defined structs and wish to use them as templates for generating their GraphQL schemas. It removes the need to manually craft GraphQL type definitions that replicate existing Go data structures, thereby enhancing developer productivity and minimizing the likelihood of errors.

# Getting Started
To get started with structogqlgen, clone the repository, install the necessary dependencies, and follow the quick setup guide provided in the documentation to integrate it into your project.

## Installation

### Using the Interactive CLI:

Ensure that you have Go installed on your local machine. If not, you can download it from the [official Go website](https://golang.org/).

Run : 

```shell
go install github.com/VintageOps/structogqlgen@latest
```

This typically compiles and Install the binary in bin directory inside your Go workspace or the global Go installation directory ($GOPATH/bin or $GOBIN), if neither $GOPATH nor $GOBIN are set, then it will install under $HOME/go/bin. 

### Using the available packages Types and Functions

```shell
go get github.com/VintageOps/structogqlgen
```

```
import github.com/VintageOps/structogqlgen
```

## Usage

```shell
~/go/bin/structogqlgen -h
StructsToGqlGenTypes is a tool that helps automatically convert Golang structs into GraphQL types
that are readily usable with the popular GraphQL framework, gqlgen. It aims to reduce the boilerplate code
required to define GraphQL schemas manually, thus accelerating the development of GraphQL APIs in Go projects.

Usage:
  structogqlgen [path] [flags]

Examples:
structogqlgen pkg/examples --use-json-tags

Flags:
  -h, --help                        help for structogqlgen
  -r, --required-tags key=value     If there is a tag that make a field required, specified that tag using the format key=value. e.g. validate=required (default [])
  -i, --tags-value-ignored string   Specify a tag value that signals to ignore a field with this tag value.Automatically set to '-' for JSON tags if not specified. (default "-")
  -c, --use-custom-tags string      Specify a custom tag to use as field name. This takes precedence over JSON tags.
  -j, --use-json-tags               Use JSON Tag as field name when available. If not present, the field name will be used.
  -v, --version                     version for structogqlgen
```

Running structogqlgen prints the generated Schema Definition on standard output (stdout), the output is segmented into two sections:

- Custom Scalar Declaration
This section displays the Custom Scalars that needs to be defined.
It is to be highlighted that gqlgen ships with some built-in helpers for [common custom scalar use-cases](https://gqlgen.com/reference/scalars/#built-in-helpers) (Time, Any, Upload and Map). Adding any of these to a schema will automatically add the marshalling behaviour to Go types.
Any other Scalar highlighted in this section needs to be implemented, [Gqlgen documentation](https://gqlgen.com/reference/scalars/#custom-scalars-with-user-defined-types) describes how this should be achieved, and there are some packages, such as [gql-bingInt](https://github.com/xplorfin/gql-bigint), that can help in this endeavour.

- Graphql Type Definitions 

### Example:

Using the example in [pkg/examples](https://github.com/VintageOps/structogqlgen/blob/main/pkg/examples) with options to make use of json tags and to use the tag validate when set to "required" for finding the required fields.

```shell
~/go/bin/structogqlgen pkg/examples --use-json-tags --required-tags validate=required
```

```graphql
scalar error
scalar interfaceEmpty
scalar interfacevalues
scalar BigInt
scalar Time
scalar PublicationStatus

type AdditionalData {
    users: [User]
    articles: [Article]
    article_comments: ArticleCommentsMap
}

type ArticleCommentsMap {
    key: Int
    values: [Int]
}

type structWithAnonymous {
    name: String
    age: Int
    unnamed_type: UnNamedType
}

type UnNamedType {
    name: String
    created_at: Time
    deleted_at: Time
}

type Another {
    Name: String
}

type Metadata {
    created_at: Time
    updated_at: Time
}

type User {
    id: Int
    username: String
    email: String
    verified: Boolean
    created_at: Time
    updated_at: Time
}

type Article {
    id: Int!
    title: String
    content: String
    author: User
    tags: [String]
    comments: [Comment]
    published_at: Time
    status: PublicationStatus
    error: error
    anything: interfaceEmpty
    do_something: DoSomethingMap
    random_int: BigInt
    another_random_int64: BigInt
    created_at: Time
    updated_at: Time
}

type DoSomethingMap {
    key: String
    values: interfacevalues
}

type Comment {
    id: Int
    article_id: Int
    author: User
    content: String
    created_at: Time
    updated_at: Time
}

type CMSData {
    users: [User]
    articles: [Article]
    article_comments: ArticleCommentsMap
}

type ArticleCommentsMap {
    key: Int
    values: [Int]
}
```
