# Golang Project0: URL shortner, a CRUD app with Concurrency and PostreSQL 

The app implements the following functionality:
- [Increment 14: Bulk Soft Delete with Transactions and Goroutines](https://github.com/allensuvorov/urlshortner/blob/main/README.md#Increment-14)
- Increment 13: Index and Duplicates handling
- Increment 12: Bulk create with JSON
- Increment 11: PostgreSQL, priority storage
- Increment 10: PostgreSQL, ping
- Increment 09: Authentification with Hash and Crypto via Middleware
- Increment 08: Compression with Gzip via Middleware
- Increment 07: App Config with Flags
- Increment 06: File, saving and restoring data
- Increment 05: Environment variables
- Increment 04: Encoding, API for JSON client, Unit-Tests for handlers
- Increment 03: HTTP libraries and frameworks - Chi router
- Increment 02: Unit-Tests, Table-driven
- Increment 01: Shortner - Server with POST and GET endpoints 

### Increment 14: Bulk Soft Delete with Transactions, SQL Statements and Goroutines
In the database table with short URLS, create an additional field with a flag, indicating that this URL should be considered deleted. Then add an asynchronous handler DELETE /api/user/urls, which accepts a list of short URL identifiers to be deleted in the format:
````
[ "a", "b", "c", "d", ...]
````
If the request is accepted successfully, the handler should return HTTP status 202 Accepted.The actual result of deletion may occur later - there is no need to notify the client of the operation's success or failure in any way.

Only the user who created the URL can successfully delete the URL. When requesting a deleted URL using the handler GET /{id}, a 410 Gone status should be returned.

Advice:
- Use batch update to effectively set the Deleted flag in the database.
- Use the fan-In pattern to maximize buffer load of update objects.

*Details on other increments to be provided upon request.*