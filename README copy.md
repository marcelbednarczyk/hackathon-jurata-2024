# Go Client 4 Hackaton Rainbow 2024

## GIT Submodules
```sh
git submodule update --init --recursive
```

## Start

### Parametry: 

```sh
count -> liczba gier
new -> czy tworzymy nowy pokój
room -> nazwa pokoju
name -> nazwa gracza
addr -> adres serwera
debug -> wlacza debugowe logi

```

### Tworzymy nową grę: 

```sh
go run main.go --count=200 --new --name test --room 123

```

### Dołączamy do gry: 

```sh
go run main.go --name test --room 123

```