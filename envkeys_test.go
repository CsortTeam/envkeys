package envkeys

import (
	"reflect"
	"testing"
)

type NestedConfig struct {
	Host string `env:"HOST, default=localhost"`
	Port int    `env:"PORT, default=5432"`
}

type TestConfig struct {
	Server   NestedConfig `env:", prefix=DB_"`
	Debug    bool         `env:"DEBUG, default=false"`
	Optional string       `env:", prefix="`
}

func TestFromStruct(t *testing.T) {
	keys := FromStruct(TestConfig{})
	got := make(map[string]struct{})
	for _, k := range keys {
		got[k] = struct{}{}
	}
	want := map[string]struct{}{
		"DB_HOST": {}, "DB_PORT": {}, "DEBUG": {},
	}
	for k := range want {
		if _, ok := got[k]; !ok {
			t.Errorf("FromStruct missing %q", k)
		}
	}
	for k := range got {
		if _, ok := want[k]; !ok {
			t.Errorf("FromStruct unexpected %q", k)
		}
	}
}

func TestNamespaceToEnvKey(t *testing.T) {
	type Config struct {
		DB   NestedConfig `env:", prefix=DB_"`
		Port int          `env:"PORT"`
	}
	m := NamespaceToEnvKey(Config{}, "Config")
	tests := []struct {
		ns   string
		want string
	}{
		{"Config.DB.Host", "DB_HOST"},
		{"Config.DB.Port", "DB_PORT"},
		{"Config.Port", "PORT"},
	}
	for _, tt := range tests {
		got, ok := m[tt.ns]
		if !ok || got != tt.want {
			t.Errorf("NamespaceToEnvKey[%q] = (%q, %v), want %q", tt.ns, got, ok, tt.want)
		}
	}
}

func TestNamespaceToKey(t *testing.T) {
	type C struct {
		X string `env:"MY_KEY"`
	}
	got := NamespaceToKey(C{}, "C", "C.X")
	if got != "MY_KEY" {
		t.Errorf("NamespaceToKey = %q, want MY_KEY", got)
	}
	got = NamespaceToKey(C{}, "C", "Unknown")
	if got != "Unknown" {
		t.Errorf("NamespaceToKey unknown = %q, want Unknown", got)
	}
}

func TestParseEnvTag(t *testing.T) {
	type S struct {
		A string `env:"KEY"`
		B int    `env:"PORT, default=80"`
		C struct {
			D string `env:"NESTED"`
		} `env:", prefix=PFX_"`
	}
	typ := reflect.TypeOf(S{})
	key, pre, has := parseEnvTag(typ.Field(0).Tag.Get("env"))
	if key != "KEY" || pre != "" || has {
		t.Errorf("KEY tag: got (%q, %q, %v)", key, pre, has)
	}
	key, pre, has = parseEnvTag(typ.Field(2).Tag.Get("env"))
	if key != "" || pre != "PFX_" || !has {
		t.Errorf("prefix tag: got (%q, %q, %v)", key, pre, has)
	}
}
