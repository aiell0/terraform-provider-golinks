

# Proveedor de Terraform Golinks

## Requisitos

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24

## Construir el Proveedor

1. Clona el repositorio
1. Ingresa al directorio del repositorio
1. Compila el proveedor utilizando el comando `install` de Go:

```shell
go install
```

## Añadir Dependencias

Este proveedor utiliza [módulos de Go](https://github.com/golang/go/wiki/Modules).
Consulta la documentación de Go para obtener la información más actualizada sobre el uso de módulos de Go.

Para agregar una nueva dependencia `github.com/author/dependency` a tu proveedor de Terraform:

```shell
go get github.com/author/dependency
go mod tidy
```

Luego, confirma los cambios en `go.mod` y `go.sum`.

## Uso del proveedor

* El ejemplo completo está en examples/full.

## Desarrollo del Proveedor

Si deseas trabajar en el proveedor, primero necesitarás tener [Go](http://www.golang.org) instalado en tu máquina (ver [Requisitos](#requirements) arriba).

Para compilar el proveedor, ejecuta `go install`. Esto compilará el proveedor y colocará el binario del proveedor en el directorio `$GOPATH/bin`.

Para generar o actualizar la documentación, ejecuta `make generate`.

Para ejecutar el conjunto completo de pruebas de aceptación, ejecuta `make testacc`.

```shell
make testacc
```
