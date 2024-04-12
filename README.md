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
