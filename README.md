# Go Url Shortener

This is a url shortener to reduce a long link. Made with Go just as a POC to undestand a little bit more about Go web development

## Autores

- [@pedrofhasler](https://www.github.com/pedrofhasler)

## Instalação

```bash
  go run main.go
  //or
  go build

```

## Documentação da API

#### Retorna todos os itens

```http
  POST /api/shortener
```

```curl
  "localhost:port/api/shortener" \ -d "your_url"
```

#### Returns a URL

```http
  GET /r/<short_id>
```

#### Returns your original url

```http
  GET /api/stats/<short_id>
```

```curl
  "localhost:port/api/stats/<short_id>"
```

#### Returns a JSON with your Url clicks
