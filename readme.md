`govalid` generates validation code for maps of strings to strings to
marshal the data into well-typed structures.

Documentation
-------------
 - [GoDoc Documentation](https://godoc.org/chrispennello.com/go/govalid)

Testing
-------
Run `go test` to test.  You may want to run `go test -short` to avoid
[fuzz testing][1] of random input programs based off of the
`test/struct.ebnf` grammar.

Future Work
-----------
 - Real functional testing.
 - Nicer error reporting.  Ideally, the validation library would give
   you an error object from which you could easily generate a
   human-readable string indicating all of the bad fields passed in and
   why they were bad.
 - Add more types to validate:
    - Enums?
 - Include comments from original source.
 - Handle weirder pre-existing imports source files, such as strconv
   imported as some other name.

Bugs
----
 - Users cannot add more types to validate as they develop without
   modifying this library.

[1]: https://en.wikipedia.org/wiki/Fuzz_testing
