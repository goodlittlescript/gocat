gocat(1) -- POSIX cat in Go
================================================================

## SYNOPSIS

`gocat` [options] [FILES...]

## DESCRIPTION

The `cat` utility in ruby, per the [POSIX specification](http://pubs.opengroup.org/onlinepubs/000095399/utilities/cat.html).

## OPTIONS

These options control how `gocat` operates.

* `-u`:
  Unbuffer output.

## EXAMPLES

```bash
gocat > file <<DOC
content
DOC

gocat file
# content
```

## INSTALLATION

Add `gocat` to your PATH (or execute it directly).

## DEVELOPMENT

Clone repo, build images.

```bash
make images
```

Run the utility, test, fix, and lint.

```bash
make run <<<"success"
make test fix lint
```

Get a shell for development.

```bash
make shell
# go build
# ./gocat <<<"success"
# ./test/suite
```

Package.

```bash
make artifacts
```
