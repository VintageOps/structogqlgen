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
NAME:
   structogqlgen - Converts Golang structs into GraphQL types that are readily usable with the popular GraphQL framework, gqlgen

USAGE:
   structogqlgen [global options] 

DESCRIPTION:
   StructsToGqlGenTypes is a tool that helps to automatically converts Golang structs into GraphQL types that are readily usable with the popular GraphQL framework, gqlgen.
   It aims to reduce the boilerplate code required to define GraphQL schemas manually, thus accelerating the development of GraphQL APIs in Go projects.

AUTHOR:
   VintageOps

GLOBAL OPTIONS:
   --src SRC_PATH, -s SRC_PATH              SRC_PATH is the required path to the source file containing the structs to import (required)
   --use-json-tags, -j                      Use JSON Tag as field name when available. If this is selected and a field has no Json tag, then the field name will be used. (default: false)
   --use-custom-tags value, -c value        Specify a custom tag to use as field name. Specifying this takes precedence over JSON tags. If specifed and a field does not have this tag, the field name will be used
   --tags-value-ignored value, -i value     Specify a tag value that signal to ignore Field with tag having this value. When using json tags with use-json-tags option, if this not specified, it is automatically set to '-'
   --required-tags key=value, -r key=value  If there is a tag that make a field required, specified that tag using the format key=value. e.g. validate=required
   --help, -h                               show help
```

### Example:

Using the example in pkg/examples_test/examples_test.go with

```shell
~/go/bin/structogqlgen --src pkg/examples_test/examples_test.go --use-json-tags --required-tags validate=required
```

```graphql
scalar interfaceEmpty
scalar interfacevalues
scalar BigInt
scalar PublicationStatus
scalar error

type Another {
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
    Metadata: Metadata
}

type DoSomethingMap {
    key: String
    values: interfacevalues
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

type Comment {
    id: Int
    article_id: Int
    author: User
    content: String
    Metadata: Metadata
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
    Metadata: Metadata
}
```