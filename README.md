# wuzu - Build tool with Docker

wuzu is a command-line tool that builds in an isolated (docker) environment.

## `.wuzu.yml` Example

```yml
# golang build

build:
  from: 'golang:1.6'
  src: $PWD
  dest: /go/src/github.com/subicura/wuzu
  run: go build -v
```

```yml
# vert.x build

build:
  from: 'vertx/vertx3'
  src: $PWD
  dest: /app
  run: ./gradlew build
```

## Build

- Go (>= 1.6)
- [godep](https://github.com/tools/godep)
