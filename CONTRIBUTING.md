# Contributing to Nginxbeat

We try to follow the [Golang coding standards](https://github.com/golang/go/wiki/CodeReviewComments) as close as possible. And, we also follow [the seven rules of a great git commit message](http://chris.beams.io/posts/git-commit/).

Anyway, make sure to run `go fmt`, `go vet`, `goimports`, `golint` before you push your code.


## Dependency Management

Nginxbeats are using [godep](https://github.com/tools/godep) for dependency management.
For updating dependencies we have the following strategy:

* If possible use the most recent release
* If no release tag exist, try to stay as close as possible to master

NOTE, make sure `GO15VENDOREXPERIMENT=1` is set so that dependencies are living in the `vendor/` folder.


### Update Dependencies

Godep allows to update all dependencies at once. We DON'T do that. If a dependency
is updated, the newest dependency must be loaded into the `$GOPATH` through either
using

`go get your-go-package-path`

or by having the package already in the `$GOPATH`with the correct version / tag.
To then save the most recent packages into Godep, run

`godep update your-go-package-path`

Avoid using `godep save ./...` or `godep update ...` as this will update all packages at
once and in case of issues it will be hard to track which one cause the issue.

After you updated the package, open a pull request where you state which package
you updated.
