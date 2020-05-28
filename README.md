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

## Install pre-compiled binary

Pick your binary from the available [releases](https://github.com/goodlittlescript/gocat/releases), download, add to PATH.

```shell
RELEASE_URL="https://github.com/goodlittlescript/gocat/releases/download/...tar.gz"
INSTALL_DIR="/usr/local/gocat"

mkdir -p "$INSTALL_DIR"
curl -L "$RELEASE_URL" | tar zxvf - -C "$INSTALL_DIR"
export PATH="$INSTALL_DIR/bin:$PATH"
```

## Install from source

Same as usual for go; note the commands are in the cmd dir.

```shell
go install github.com/goodlittlescript/gocat/cmd/{gocat,gocp}
```

One way to install manpages is to put them under GOPATH (if `$GOPATH/bin` is on your PATH, then they will be autodiscovered by `man`).

```shell
mkdir -p "$GOPATH/man/man1"
ln -s "$GOPATH"/github.com/goodlittlescript/gocat/man/man1/{gocat,gocp}.1 "$GOPATH/man/man1"
```
