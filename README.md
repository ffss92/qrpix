# qrpix

Lib para gerar qrcode estático de cobrança do pix em Go, além fazer encode e decode de BRCodes.

## TODOs

1. Criar consts com os IDs dos campos;
2. Criar decoder;
3. Completar a especificação (Unreserved Templates).

## Exemplos

- Salvando imagem

```go
qr := NewStatic(
    "123e4567-e12b-12d1-a456-426655440000",
    "Fulano de Tal",
    "BRASILIA",
    "***",
)
if err := qr.SaveFile("example.png"); err != nil {
    return err
}
```

- Servindo via HTTP

```go
qr := NewStatic(
    "123e4567-e12b-12d1-a456-426655440000",
    "Fulano de Tal",
    "BRASILIA",
    "***",
)
http.HandleFunc("/", func(w http.ResponseWrite, r *http.Request) {
    if err := qr.Serve(w); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
})
```
