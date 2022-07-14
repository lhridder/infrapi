# Infrapi - API gateway for infrared

utility for [lhridder/infrared](https://github.com/lhridder/infrared)

### Global config.json
```json
{
  "apiBind": ":5000",
  "redisHost": "localhost",
  "redisDB": 0,
  "redisPass": ""
}
```

## API
### Routes examples
GET `/` will return 200(OK)

GET `/proxies` will return
```json
[
"config",
"config2"
]
```

GET `/proxies/{name}` will return
```json
{
"domainNames": ["play.example.org"],
"proxyTo": "backend.example.org:25566"
}
```

POST `/proxies/{name}` with body
```json
{
"domainNames": ["play.example.org"],
"proxyTo": "backend.example.org:25566"
}
```
will return 200(ok)

DELETE `/proxies/{name}` will return 200(OK)

## Used sources
- [Redis library for golang](https://github.com/go-redis/redis/v8)
- [Chi router](https://github.com/go-chi/chi)