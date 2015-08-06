TODO
----
 - Fuzz test with randomly-generated, but linguistically valid input go
   source files.  Auto-generate based on grammar from go/\* packages?
 - Generate GoDoc comment for generated function.
 - Document how if the struct has no validate-able fields, there will be
   a compiler error because of the unused err variable.
 - Godoc functions, etc.
 - Document how ints, uints can be specified in "any" base.

Future Work
-----------
 - Include comments from original source `.v`.
 - Add more types to validate beyond strconv "low-hanging fruit".
    - Bounded strings
    - Bounded numbers
    - `*url.URL`
    - `*mail.Address`
    - Enums?
 - Handle weirder pre-existing imports in `.v` files, such as strconv
   imported as some other name.
 - Nicer error reporting.  Ideally, the validation library would give
   you an error object from which you could easily generate a
   human-readable string indicating all of the bad fields passed in and
   why they were bad.
