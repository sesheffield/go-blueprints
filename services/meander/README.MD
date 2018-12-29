* A simple flexible checklist for writing enumerators in Go is:
- Define a new type, based on a primitive integer type
- Use that type whenever you need users to specify one of the appropriate values
- Use the iota keyword to set the values in a const block, disregarding the first zero value
- Implement a map of sensible string representations to the values of your enumerator
- Implement a String method on the type that returns the appropriate string representation from the map
- Implement a ParseType function that converts from a string to your type using the map
