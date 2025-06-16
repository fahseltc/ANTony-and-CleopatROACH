package environment

import (
	"bytes"
	"encoding/json"
	"gamejam/data"
	"io"
	"sync"
)

type Config struct {
	values map[string]any
	sync.Mutex
}

func newConfig() *Config {
	config := &Config{
		values: make(map[string]any),
	}
	jsonFile, _ := data.Files.Open("config.json")
	jsonBytes, _ := io.ReadAll(jsonFile)

	decoder := json.NewDecoder(
		// Use a bytes.Reader to treat embeddedConfig as an io.Reader
		bytes.NewReader(jsonBytes),
	)
	decoder.UseNumber()

	var data map[string]interface{}

	for decoder.More() {
		decoder.Decode(&data)
	}
	var process func(prefix string, m map[string]interface{})
	process = func(prefix string, m map[string]interface{}) {
		for key, val := range m {
			fullKey := key
			if prefix != "" {
				fullKey = prefix + "." + key
			}
			switch v := val.(type) {
			case map[string]interface{}:
				process(fullKey, v)
			case json.Number:
				if n, err := v.Int64(); err == nil {
					config.Add(fullKey, int(n))
				} else if f, err := v.Float64(); err == nil {
					config.Add(fullKey, f)
				} else {
					config.Add(fullKey, v)
				}
			case string:
				config.Add(fullKey, v)
			default:
				config.Add(fullKey, v)
				//log.Fatalf("Unknown JSON type in config, key: %v, val:%v, type:%v", fullKey, v, reflect.TypeOf(v))
			}
		}
	}

	process("", data)

	return config
}

func (r *Config) Add(k string, v any) {
	r.Lock()
	defer r.Unlock()
	r.values[k] = v
}

func (r *Config) Get(k string) any {
	r.Lock()
	defer r.Unlock()
	return r.values[k]
}
