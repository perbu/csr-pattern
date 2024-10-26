# CSR Pattern example code

This is a simple API, written in Go, that demonstrates the CSR pattern. The API is a simple CRUD API that allows you 
to create, read, update, and delete key/values. The API is backed by a Sqlite database.

In particular this API demonstrates the steps needed to avoid abstraction leaks. Errors in the repo layer
must be handled in the service layer and new errors created and returned to the caller. 

## Mocking

I also use this to evaluate testify/mock.


