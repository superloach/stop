# stop it

proof-of-concept system for automatically managing a function context in go.

## todo
- [x] make Never private
- [x] remove Handle type
- [x] add GoNothing
- [x] gofumpt
- [x] automagically pass context down to children
- [x] free context after function exit (refcount for children)
- [x] look upwards into call stack
- [x] implement a trimap?
- [ ] add Yield function which fetches the Handle for the current fn
- [ ] graceful handling instead of panic
- [ ] remove debug logs
- [ ] add warning message to any program that uses this
- [ ] add proper tests
- [ ] integrate with stdlib context?
- [ ] add docs

## license
this software is provided without warranty, primarily for educational purposes. you are free to use this code for any purpose, although it is not recommended in any critical use cases.
