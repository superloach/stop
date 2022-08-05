# stop it

proof-of-concept system for automatically managing a function context in go.

todo:
- [x] make Never private
- [x] remove Handle type
- [x] add GoNothing
- [x] gofumpt
- [ ] automagically pass context down to children
- [ ] free context after function exit (refcount for children)
- [ ] look upwards into call stack
- [ ] graceful handling instead of panic
- [ ] remove debug logs
- [ ] add warning message to any program that uses this
- [ ] add proper tests
- [ ] integrate with stdlib context?
- [ ] add docs
