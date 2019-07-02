# Face Detection Processor

Serviço responsável por ler os frames do kafka, detectar todas as faces presentes no frame e publicar novamente no kafka.

## Requisitos

A aplicação foi testada com as ferramentas e versões abaixo.

* go - 1.10.4

**Observação**: O `docker-compose.yml` possui alguns parâmetros que só estão presentes no docker-compose versão utilize uma versão acima da 1.18, pois utilizamos algumas configurações

### Bibliotecas

* [pigo](https://github.com/esimov/pigo)
* [gg](https://github.com/fogleman/gg)

## Instalação

Verifique se você possui o `go` instalado em sua máquina. A saída do comando deve ser parecida com a mostrada abaixo.

```bash
$ go version
go version go1.10.4 linux/amd64
```

Caso você tenha o `go` instalado em sua máquina, utilize o comando abaixo para instalar as bibliotecas necessárias.

```bash
go get -u github.com/esimov/pigo/cmd/pigo
go get -u github.com/fogleman/gg
go get -u github.com/disintegration/imaging
go get -u github.com/lovoo/goka
```

## Como utilizar

Primeiro, verifique se o comando abaixo mostra uma saída parecida.

```bash
$ docker-compose --version
docker-compose version 1.24.0, build 0aa59064
```

**Importante**: Caso sua versão seja abaixo da 1.18, é necessário que você atualize-o antes de prosseguir.

Agora podemos "levantar" os containers do kafka usando o comando abaixo.

```bash
docker-compose up
```

Depois disso, compile a aplicação e execute-a.

```bash
go build
./face-detection-processor
```
