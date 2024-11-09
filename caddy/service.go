package caddy

import (
	"encoding/json"
	caddy2 "github.com/caddyserver/caddy/v2"
	"paas/httpjson"
)

var RawConfig = `{
  "apps": {
    "http": {
      "servers": {
        "srv0": {
          "listen": [
            ":80"
          ],
          "routes": [
            {
              "handle": [
                {
                  "handler": "subroute",
                  "routes": [
                    {
                      "handle": [
                        {
                          "handler": "reverse_proxy",
                          "upstreams": [
                            {
                              "dial": ":3000"
                            }
                          ]
                        }
                      ]
                    }
                  ]
                }
              ],
              "match": [
                {
                  "host": [
                    "htmgo.dev"
                  ]
                }
              ],
              "terminal": true
            }
          ]
        }
      }
    }
  }
}`

func BuildConfig() *Config {
	var config Config
	err := json.Unmarshal([]byte(RawConfig), &config)
	if err != nil {
		panic(err)
	}
	return &config
}

func GetConfig() (*Config, error) {
	config, err := httpjson.Get[Config]("http://localhost:2019/config/")
	if err != nil {
		return nil, err
	}
	return config, nil
}

func ApplyConfig(config *Config) error {
	serialized, _ := json.Marshal(config)
	return caddy2.Load(serialized, true)
}
