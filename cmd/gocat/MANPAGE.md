NAME
====

**gocat(1)** -- posix cat in go

SYNOPSIS
========

| **gocat** `[options]` `[file ...]`

DESCRIPTION
===========

The **cat** utility in go, per the [POSIX specification](http://pubs.opengroup.org/onlinepubs/000095399/utilities/cat.html).

OPTIONS
=======

These options control how `gocat` operates.

`-u`
  ~ Unbuffer output.

EXAMPLE
========

Use **gocat** to create a file:

```sh
gocat > file <<DOC
content
DOC
```

Use **gocat** to print the contents of the file:

```sh
gocat file
# content
```
