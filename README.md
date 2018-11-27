git-semv
==

git-semv help semantic versioning for git tag.

[![travis](https://img.shields.io/travis/linyows/git-semv.svg?style=for-the-badge)][travis]
[![release](http://img.shields.io/github/release/linyows/git-semv.svg?style=for-the-badge)][release]
[![license](http://img.shields.io/badge/license-MIT-blue.svg?style=for-the-badge)][license]
[![godoc](http://img.shields.io/badge/go-documentation-blue.svg?style=for-the-badge)][godoc]
[![codecov](https://img.shields.io/codecov/c/github/linyows/git-semv.svg?style=for-the-badge)][codecov]

[travis]: https://travis-ci.org/linyows/git-semv
[release]: https://github.com/linyows/git-semv/releases
[license]: https://github.com/linyows/git-semv/blob/master/LICENSE
[godoc]: http://godoc.org/github.com/linyows/git-semv
[codecov]: https://codecov.io/gh/linyows/git-semv

Installation
--

```sh
$ go get github.com/linyows/git-semv
```

Usage
--

Show list:

```sh
$ git semv list
v1.2.1-alpha.0
v1.2.0
v1.2.0-rc.1
v1.2.0-rc.0
v1.2.0-beta.0+ba8a247.foobar
v1.2.0-alpha.0+a2a784b.anonymous
v1.1.0
v1.0.1
v1.0.0
```

Show latest version:

```sh
$ git semv now
v1.2.0
```

Show next version:

```sh
# with pre-release option
$ git semv patch --pre
v1.2.1-alpha.1
# specify pre-release name option
$ git semv patch --pre-name beta
v1.2.1-beta.0
# next minor version
$ git semv minor
v1.3.0
# with bump option
$ git semv minor --bump
git tag v1.3.0 && git push origin v1.3.0
# next major version with build option
$ git semv major --pre --build
v2.0.0-alpha.0+9125b23.linyows
```

Contribution
------------

1. Fork ([https://github.com/linyows/git-semv/fork](https://github.com/linyows/git-semv/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

Author
--

[linyows](https://github.com/linyows)
