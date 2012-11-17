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

```sh

# stores the github authentication data
$ goctogit login

# create a repository on github named 'reponame' with some description
$ goctogit create -descr 'some description' <reponame>

# lists the available downloads on a github repo
$ goctogit dl-ls -org mana-fwk mana-release
github-dl-ls: listing downloads for repository [mana-release] with account [mana-fwk]...
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
github-dl-ls: listing downloads for repository [mana-release] with account [mana-fwk]... [done]

# deletes the download with id=350409 from a github repository
$ goctogit dl-rm -org mana-fwk -repo mana-release 36354
github-dl-rm: deleting download id=350409 from repository [mana-release] with account [mana-fwk]...
github-dl-rm: deleting download id=350409 from repository [mana-release] with account [mana-fwk]...[done]

```

TODO
----

- goctogit issues
- goctogit issues <number>
- goctogit issues <number> close
- others ?


