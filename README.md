git-semv
==

This is a [Git][git] plugin for [Semantic Versioning][semver].

The usefulness of Semantic Versioning has been accepted by OSS in advance.
And, with the appearance of Go modules, Semantic Versioning became indispensable for the development of the Go library.
However, `git tag` used for versioning can not support pre-releases and sorts including build information for managing Semantic Versioning.
Also, it is not easy to increment the version.
This `git-semv` is a CLI tool for solving these problems and functions as git subcommand.

[![travis](https://img.shields.io/travis/linyows/git-semv.svg?style=for-the-badge)][travis]
[![codecov](https://img.shields.io/codecov/c/github/linyows/git-semv.svg?style=for-the-badge)][codecov]
[![release](http://img.shields.io/github/release/linyows/git-semv.svg?style=for-the-badge)][release]
[![godoc](http://img.shields.io/badge/go-documentation-blue.svg?style=for-the-badge)][godoc]
[![license](http://img.shields.io/badge/license-MIT-blue.svg?style=for-the-badge)][license]

[travis]: https://travis-ci.org/linyows/git-semv
[release]: https://github.com/linyows/git-semv/releases
[license]: https://github.com/linyows/git-semv/blob/master/LICENSE
[godoc]: http://godoc.org/github.com/linyows/git-semv
[codecov]: https://codecov.io/gh/linyows/git-semv
[semver]: https://semver.org/
[git]: https://git-scm.com/

Japanese: https://tomohisaoda.com/posts/2018/do-semantic-versioning-for-app.html

Installation
--

Download the binary in [Github Releases][release] and place it in the directory where `$PATH` passed.
Or, you can download using `go get` depending on the version of Go1.11 or higher.

```sh
$ go get -u github.com/linyows/git-semv/cmd/git-semv
```

### Homebrew

```sh
$ brew tap linyows/git-semv
$ brew install git-semv
```

Usage
--

Show list:

```sh
# Only release versions
$ git semv
v0.0.1
v0.0.2
v1.0.0
v1.1.0
v1.1.1

# All versions including pre-release
$ git semv -a
v0.0.1
v0.0.2
v1.0.0-alpha.0+a2a784b.linyows
v1.0.0-beta.0+ba8a247.foobar
v1.0.0-rc.0
v1.0.0-rc.1
v1.0.0
v1.1.0
v1.1.1
v2.0.0-alpha.0
```

Show latest version:

```sh
$ git semv now
v1.1.1
```

Show next version(major|minor|patch):

```sh
# Next patch version
$ git semv patch
v1.1.2

# Next minor version
$ git semv minor
v1.2.0

# Next major version
$ git semv major
v2.0.0
```

Use options(pre|pre-name|build|build-name|bump):

```sh
# Next pre-release as major
$ git semv major --pre
v2.0.0-alpha.1

# Specify pre-release name as major
$ git semv major --pre-name beta
v2.0.0-beta.0

# Next minor version with build info
$ git semv minor --build
v1.2.0+9125b23.linyows

# Specify build name
$ git semv minor --build-name superproject
v1.2.0+superproject

# Create tag and Push origin
$ git semv patch --bump
Bumped version to v1.1.2
#==> git tag v1.1.2 && git push origin v1.1.2
```

VS.
--

### motemen/gobump

[gobump][gobump] will increment the version according to semver in version in the source code.
On the other hand, `git-semv` does not do anything to the source code.
Even if you do, you just create a tag and push it remotely.
When focusing on `Go`, `Go` can add version and other information to the build, so there is no need to manage version in code.
Also, in other languages, you can easily replace them in code by combining with commands such as `sed`.

and `git-semv` supports versioning of pre-release and build information.

[gobump]: https://github.com/motemen/gobump

Development flow
--

The assumed development flow is...

1. Development
1. Remote push
1. Pull-request create
1. Continuous Integration
1. Master branch merge
1. Tag create and push(git-semv)
1. Continuous Integration
1. Release create([goreleaser][goreleaser])

Generally, development in Go will upload the product binary to github releases and release the product.
There is a great tool called [goreleaser][goreleaser] which makes that work easier.
By running this tool on the CI, we will automatically release the binary after pushing the created tag.
And `git-semv` solves troublesome versioning and tag creation problem which is the next bottleneck.

[goreleaser]: https://github.com/goreleaser/goreleaser

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
