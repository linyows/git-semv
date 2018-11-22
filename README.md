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

```sh
$ git semv current
v1.2.3
$ git semv bump
v1.2.4
$ git semv --minor bump
v1.3.0
$ git semv --major bump
v2.0.0
$ git semv --major --pre bump
v2.0.0-rc.0
```

```sh
$ VER=$(git semv bump)
$ git tag $VER && git push origin $VER
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
