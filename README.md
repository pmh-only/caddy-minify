# File (html/css/js/etc..) minification module for Caddy v2
![](docs/draw.png)
powered by [tdewolff/minify](https://github.com/tdewolff/minify) & inspired by [hacdias/caddy-v1-minify](https://github.com/hacdias/caddy-v1-minify)

## Caddy module ID
```
http.handlers.minify
```

## Support formats
* html: `text/html`
* css: `text/css`
* svg: `image/svg+xml`
* js: `^(application|text)/(x-)?(java|ecma)script$`
* json: `[/+]json$`
* xml: `[/+]xml$`

## Performance
See https://github.com/tdewolff/minify#performance

## How to use
Make sure to order the handler in the correct place:

```Caddyfile
{
    order minify after encode
}
```

With file_server:

```
example.com {
    minify
    root ./test
    file_server     
}
```

or reverse_proxy:

```
example.com {
    minify
    reverse_proxy localhost:8080
}
```

limit formats (format all when not specified.):

```
example.com {
    minify {
        formats html css js
    }
    
    reverse_proxy
}
```

or simply:

```
example.com {
    minify html css js
    reverse_proxy localhost:8080
}
```

## Limitations
This module doesn't minify the original responses that have already been compressed, It just skips them.
 
To work around this, you may send the `Accept-Encoding: identity` request header to the upstream to tell it not to compress the response. For example:

```caddyfile
reverse_proxy localhost:8080 {
    header_up Accept-Encoding identity
}
```
