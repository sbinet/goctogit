goctogit
========

``goctogit`` is a simple reimplementation of
[octogit](https://github.com/myusuf3/octogit) in ``go``.

Installation
------------

```sh
$ go get github.com/sbinet/goctogit
```

Using ``goctogit``
------------------

Available commands:

## ``login``
```sh

# store the github authentication data
$ goctogit login
```

## ``create``
```sh
# create a repository on github named 'reponame' with some description
$ goctogit create -descr 'some description' <reponame>
```

## ``dl-ls``
```sh
# list the available downloads on a github repo
$ goctogit dl-ls -org mana-fwk mana-release
github-dl-ls: listing downloads for [mana-fwk/mana-release]...
=== mana-20121116-000.tar.gz
    id=363354
    sz=4782080 bytes
    https://github.com/downloads/mana-fwk/mana-release/mana-20121116-000.tar.gz
=== mana-20121115-000.tar.gz
    id=362274
    sz=4780032 bytes
    https://github.com/downloads/mana-fwk/mana-release/mana-20121115-000.tar.gz
=== mana-20121108-001.tar.gz
    id=356891
    sz=4751360 bytes
    https://github.com/downloads/mana-fwk/mana-release/mana-20121108-001.tar.gz
=== mana-20121031.tar.gz
    id=350409
    sz=4731904 bytes
    https://github.com/downloads/mana-fwk/mana-release/mana-20121031.tar.gz
github-dl-ls: listing downloads for [mana-fwk/mana-release]... [done]
```

## ``dl-rm``
```sh
# delete the download with id=350409 from a github repository
$ goctogit dl-rm -org mana-fwk -repo mana-release 350409
github-dl-rm: deleting download id=350409 from [mana-fwk/mana-release]...
github-dl-rm: deleting download id=350409 from [mana-fwk/mana-release]... [done]
```

## ``dl-create``
```sh
# upload a file to the download page of a github repository
$ goctogit dl-create -descr "a new tarball" -f foo.tar.gz -repo=mana-release -org my-org
github-dl-create: uploading [foo.tar.gz] to [my-org/mana-release]...
github-dl-create: uploading [foo.tar.gz] to [my-org/mana-release]... [done]
```

TODO
----

- goctogit issues
- goctogit issues <number>
- goctogit issues <number> close
- others ?


