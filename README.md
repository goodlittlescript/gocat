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

Clone repo, build images, get a shell for development.

```bash
./Projectfile images shell
# go build
# ./gocat <<<"success"
# ./test/suite
```

See `./Projectfile -h` for additional help.
