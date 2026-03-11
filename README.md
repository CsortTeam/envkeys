# envkeys

Извлекает env-ключи из struct-тегов [go-envconfig](https://github.com/sethvargo/go-envconfig).

```bash
go get github.com/CsortTeam/envkeys
```

```go
import "github.com/CsortTeam/envkeys"

keys := envkeys.FromStruct(&Config{})
m := envkeys.NamespaceToEnvKey(&Config{}, "Config")  // namespace → env key
```

API: `FromStruct`, `NamespaceToEnvKey`, `NamespaceToKey`. Теги: `env:"KEY"`, `env:", prefix=PREFIX"`.
