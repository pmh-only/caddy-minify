# File (html/css/js/etc..) minification plugin for Caddy v2
powered by [tdewolff/minify](https://github.com/tdewolff/minify) & inspired by [hacdias/caddy-v1-minify](https://github.com/hacdias/caddy-v1-minify)

## Caddy module ID
```
http.handlers.minify
```

## How to use
JSON config:
```json
{
  "handler": "minify"
}
```

or Caddyfile (Make sure to order the handler in the correct place.):

```Caddyfile
{
    order minify after encode
}
```

```
example.com {
    minify
    root ./test
    file_server     
}
```

```
example.com {
    minify
    reverse_proxy example.net
}
```
