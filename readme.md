TODO
----
 - Fuzz test with randomly-generated, but linguistically valid input go
   source files.  Auto-generate based on grammar from go/\* packages?
 - Document how if the struct has no validate-able fields, there will be
   a compiler error because of the unused err variable.
 - Godoc functions, etc.

Integer Bases
-------------
The generated validation code for `int`s and `uint`s uses a base of 0,
so input strings may be in any base represented by the
`strconv.ParseInt` or `strconv.ParseUint` functions.  For example, a
hexadecimal values would be parsed by passing in `0xbeef`.

Future Work
-----------
 - Add more types to validate:
    - `*url.URL`
    - `*mail.Address`
 - Add more types to validate:
    - Bounded strings
    - Bounded numbers
 - Optional fields.
 - Nicer error reporting.  Ideally, the validation library would give
   you an error object from which you could easily generate a
   human-readable string indicating all of the bad fields passed in and
   why they were bad.
 - Add more types to validate:
    - Enums?
 - Include comments from original source.
 - Handle weirder pre-existing imports source files, such as strconv
   imported as some other name.
