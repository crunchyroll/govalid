TODO
----
 - Finish writing validators.
 - Fuzz test with randomly-generated, but linguistically valid input go
   source files.  Auto-generate based on grammar from go/\* packages?
 - Godoc functions, etc.

Future Work
-----------
 - Include comments from original source `.v`.
 - Add more types to validate beyond strconv "low-hanging fruit".
 - Handle weirder pre-existing imports in `.v` files, such as strconv
   imported as some other name.
