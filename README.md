## Setup

```shell
$ asdf plugin add golang https://github.com/kennyp/asdf-golang.git
$ asdf plugin add tinygo https://github.com/schmir/asdf-tinygo.git
$ asdf plugin add binaryen https://github.com/birros/asdf-binaryen.git
$ asdf install
```

## Usage

```shell
$ make
```

## Gomobile

```shell
$ go install golang.org/x/mobile/cmd/gomobile
$ go install golang.org/x/mobile/cmd/gobind
$ asdf reshim golang
$ gomobile init
```
