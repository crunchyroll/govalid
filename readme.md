TODO
----
 - See if there's a way to include the original comments.
 - Fuzz test with randomly-generated, but linguistically valid input go
   source files.  Auto-generate based on grammar from go/\* packages?
 - Godoc functions, etc.

Future Work
-----------
 - Add more types to validate beyond strconv "low-hanging fruit".
 - Handle weirder pre-existing imports in `.v` files, such as strconv
   imported as some other name.
