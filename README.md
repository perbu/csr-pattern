# CSR Pattern example code

This is a simple API, written in Go, that demonstrates the CSR pattern. The API is a simple CRUD API that allows you
to create, read, update, and delete key/values. The API is backed by a Sqlite database.

In particular this API demonstrates the steps needed to avoid abstraction leaks. Errors in the repo layer
must be handled in the service layer and new errors created and returned to the caller.

## Why prevent abstraction leaks?

Reduce complexity: Abstractions are meant to simplify complex systems by hiding unnecessary details. When an abstraction
leaks, it forces developers to understand those hidden details, making the code harder to maintain and modify. 
This can lead to increased development time and higher risk of introducing bugs.

Reusability: A well-designed abstraction can be reused in different parts of the application or even in different
projects. However, a leaky abstraction often ties the code to specific implementation details, limiting its reusability
and leading to code duplication.

Testability: Leaky abstractions can make it difficult to write unit tests because the code becomes dependent on the
underlying implementation. This can lead to complex and brittle tests that are hard to maintain.

## Mocking

I also use this to evaluate testify/mock.


