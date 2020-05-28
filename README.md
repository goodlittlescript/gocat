# GoCat

A package for building command line utilities in go that follow the `cat` and `cp` signatures. Also implements both per their respective specifications.

- [gocat](cmd/gocat/MANPAGE.md) - [POSIX specification](http://pubs.opengroup.org/onlinepubs/000095399/utilities/cat.html).
- [gocp](cmd/gocp/MANPAGE.md) - [POSIX specification](http://pubs.opengroup.org/onlinepubs/000095399/utilities/cp.html)

## Development

Clone repo, build images, get a shell for development:

```bash
./Projectfile images 
./Projectfile shell
```

Now in the shell:

```bash
go install ./cmd/gocat
gocat <<<"success" > afile
gocat afile
# success

go install ./cmd/gocp
gocp afile bfile
gocat afile bfile
# success
# success

./Projectfile manpages
man gocat
man gocp
./test/suite
```

See [./Projectfile](./Projectfile) for the workflow.

## Compiling from source

Same as usual, just note the command is in the cmd dir.

```shell
go install github.com/goodlittlescript/gocat/cmd/{gocat,gocp}
```

To add manpages, symlink as follows:

```shell
mkdir -p "$GOPATH/man/man1"
ln -s "$GOPATH"/github.com/goodlittlescript/gocat/man/man1/{gocat,gocp}.1 "$GOPATH/man/man1"
```
