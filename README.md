# qrpix

Lib para gerar qrcode estático de cobrança do pix em Go, além fazer encode e decode de BRCodes.

## TODOs

1. Criar consts com os IDs dos campos;
2. Criar decoder;
3. Completar a especificação (Unreserved Templates).

## Exemplos

```go
qr := NewStatic(
    "123e4567-e12b-12d1-a456-426655440000", // Chave Pix
    "Fulano de Tal", // Nome
    "BRASILIA", // Cidade
    "***", // ID Transação
)

// Gerando Imagem
if err := qr.SaveFile("example.png"); err != nil {
    return err
}

// Servindo via HTTP
http.HandleFunc("/", func(w http.ResponseWrite, r *http.Request) {
    if err := qr.Serve(w); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
})
```

- Campos Opcionais

```go
qr := NewStatic(
    "123e4567-e12b-12d1-a456-426655440000", // Chave Pix
    "Fulano de Tal", // Nome
    "BRASILIA", // Cidade
    "***", // ID Transação
    WithTransactionAmount(1000), // Valor da transação em centavos (10 reais)
)
```

- Decode

```go
code := "00020126580014br.gov.bcb.pix0136123e4567-e12b-12d1-a456-4266554400005204000053039865802BR5913Fulano de Tal6008BRASILIA62070503***63041D3D"
p := NewParser()
static, _ := p.ParseStatic(code)
fmt.Printf("%+v", static)
```
