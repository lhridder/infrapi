package infrapi

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-redis/redis/v8"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Event struct {
	Name   string          `json:"name"`
	Config json.RawMessage `json:"config"`
}

var (
	ctx context.Context
	rdb *redis.Client
)

func ConnectRedis() error {
	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr:     Config.RedisHost + ":6379",
		Password: Config.RedisPass,
		DB:       Config.RedisDB,
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return err
	}

	return nil
}

func ListenAndServe() error {
	apiBind := Config.ApiBind
	log.Printf("Starting API on %s", apiBind)

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", getIndex())
	router.Get("/proxies", getProxies())
	router.Get("/proxies/{name}", getProxy())
	router.Post("/proxies/{name}", postProxy())
	router.Delete("/proxies/{name}", deleteProxy())

	return http.ListenAndServe(apiBind, router)
}

// getIndex sends back empty 200(OK)
func getIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

// getProxies sends back list of proxies
func getProxies() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		configs, err := configList()
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
		}

		data, err := json.Marshal(configs)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
		}

		_, err = w.Write(data)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

// getProxy sends back the config of a requested proxy
func getProxy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		config, err := rdb.Get(ctx, "config:"+name).Result()
		if err != nil {
			if err == redis.Nil {
				http.Error(w, http.StatusText(404), 404)
			} else {
				log.Println(err)
				http.Error(w, http.StatusText(500), 500)
			}
		}

		_, err = w.Write([]byte(config))
		if err != nil {
			log.Println(err)
			return
		}
	}
}

// postProxy accepts a valid ProxyConfig and stores it in redis
func postProxy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")

		rawData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(400), 400)
		}

		var cfg ProxyConfig
		err = json.Unmarshal(rawData, &cfg)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(400), 400)
		}

		configs, err := configList()
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
		}

		event := Event{
			Name:   name,
			Config: rawData,
		}

		data, err := json.Marshal(event)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
		}

		if err = rdb.Set(ctx, "config:"+name, rawData, 0).Err(); err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
		}

		if contains(configs, name) {
			if err := rdb.Publish(ctx, "infrared-edit-config", data).Err(); err != nil {
				log.Println(err)
				http.Error(w, http.StatusText(500), 500)
			}
			return
		} else {
			if err := rdb.Publish(ctx, "infrared-add-config", data).Err(); err != nil {
				log.Println(err)
				http.Error(w, http.StatusText(500), 500)
			}
			return
		}
	}
}

// deleteProxy will delete the specified config
func deleteProxy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")

		_, err := rdb.Get(ctx, "config:"+name).Result()
		if err != nil {
			if err == redis.Nil {
				http.Error(w, http.StatusText(404), 404)
			} else {
				log.Println(err)
				http.Error(w, http.StatusText(500), 500)
			}
		}

		if err = rdb.Del(ctx, "config:"+name).Err(); err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
		}
		if err := rdb.Publish(ctx, "infrared-delete-config", name).Err(); err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
		}
		return
	}
}

func configList() ([]string, error) {
	configList, err := rdb.Keys(ctx, "config:*").Result()
	if err != nil {
		return nil, err
	}

	var configs []string
	for _, config := range configList {
		configs = append(configs, strings.Split(config, "config:")[1])
	}
	return configs, nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
