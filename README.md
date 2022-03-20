## Setup

```shell
$ asdf plugin add golang https://github.com/kennyp/asdf-golang.git
$ asdf plugin add tinygo https://github.com/schmir/asdf-tinygo.git 
$ make
```

## Gomobile

```shell
$ go install golang.org/x/mobile/cmd/gomobile
$ go install golang.org/x/mobile/cmd/gobind
$ asdf reshim golang
$ gomobile init
```
