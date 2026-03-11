// Package envkeys extracts environment variable keys from struct tags compatible with
// github.com/sethvargo/go-envconfig. Use for contract-check drift detection and
// mapping validator namespace errors to env key names.
package envkeys

import (
	"reflect"
	"sort"
	"strings"
)

// FromStruct extracts env keys from a struct type (struct or *struct).
func FromStruct(v interface{}) []string {
	t := reflect.TypeOf(v)
	m := extractEnvKeys(t, "")
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// NamespaceToEnvKey maps validator namespace to env key. rootName is the top-level struct name.
func NamespaceToEnvKey(v interface{}, rootName string) map[string]string {
	t := reflect.TypeOf(v)
	return buildNamespaceToEnvMap(t, rootName, "")
}

// NamespaceToKey returns the env key for namespace ns, or ns if unknown.
func NamespaceToKey(v interface{}, rootName, ns string) string {
	m := NamespaceToEnvKey(v, rootName)
	if k, ok := m[ns]; ok {
		return k
	}
	return ns
}

func extractEnvKeys(t reflect.Type, prefix string) map[string]struct{} {
	out := make(map[string]struct{})
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return out
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("env")
		if tag == "" {
			continue
		}
		key, prefixVal, hasPrefix := parseEnvTag(tag)
		if hasPrefix {
			ft := f.Type
			if ft.Kind() == reflect.Ptr {
				ft = ft.Elem()
			}
			if ft.Kind() == reflect.Struct {
				for k := range extractEnvKeys(ft, prefix+prefixVal) {
					out[k] = struct{}{}
				}
			}
			continue
		}
		if key != "" {
			out[prefix+key] = struct{}{}
		}
	}
	return out
}

func buildNamespaceToEnvMap(t reflect.Type, nsPrefix, envPrefix string) map[string]string {
	out := make(map[string]string)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return out
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name := f.Name
		fullNs := nsPrefix + "." + name
		tag := f.Tag.Get("env")
		if tag == "" {
			continue
		}
		key, prefixVal, hasPrefix := parseEnvTag(tag)
		if hasPrefix {
			ft := f.Type
			if ft.Kind() == reflect.Ptr {
				ft = ft.Elem()
			}
			if ft.Kind() == reflect.Struct {
				for k, v := range buildNamespaceToEnvMap(ft, fullNs, envPrefix+prefixVal) {
					out[k] = v
				}
			}
			continue
		}
		if key != "" {
			out[fullNs] = envPrefix + key
		}
	}
	return out
}

func parseEnvTag(tag string) (key, prefix string, hasPrefix bool) {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return "", "", false
	}
	if strings.HasPrefix(tag, ",") {
		rest := strings.TrimSpace(strings.TrimPrefix(tag, ","))
		for _, part := range strings.Split(rest, ",") {
			part = strings.TrimSpace(part)
			if strings.HasPrefix(part, "prefix=") {
				prefix = strings.TrimPrefix(part, "prefix=")
				prefix = strings.Trim(prefix, "\"")
				return "", prefix, true
			}
		}
		return "", "", false
	}
	key = strings.SplitN(tag, ",", 2)[0]
	key = strings.TrimSpace(key)
	return key, "", false
}
