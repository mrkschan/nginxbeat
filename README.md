[![Build Status](https://travis-ci.org/mrkschan/nginxbeat.svg?branch=travisci-setup)](https://travis-ci.org/mrkschan/nginxbeat)

# Nginxbeat

Nginxbeat is the Beat used for Nginx monitoring. It is a lightweight agent that reads status from Nginx periodically. Nginx must either expose its status via stub module (http://nginx.org/en/docs/http/ngx_http_stub_status_module.html) or Nginx Plus status module (http://nginx.org/en/docs/http/ngx_http_status_module.html).


## Elasticsearch template

To apply Nginxbeat template for Nginx status:

```
curl -XPUT 'http://localhost:9200/_template/nginxbeat' -d@etc/nginxbeat.template.json
```


# Build, Test, Run

```
# Build
export GO15VENDOREXPERIMENT=1
GOPATH=<your go path> godep restore
GOPATH=<your go path> make

# Test
GOPATH=<your go path> make unit-tests

# Run whole testsuite
GOPATH=<your go path> make testsuite

# Run
./nginxbeat -c nginxbeat.yml

# Make binaries
GOPATH=<your go path> make crosscompile
```


## Exported fields

Nginxbeat only exports several types of document. The properties in the document varies according to the configured Nginx status page.

- `type: stub` holds Nginx stub status
- `type: plus` holds Nginx Plus status
- `type: zone` holds Nginx Plus status zone status
- `type: upstream` holds Nginx Plus upstream group status
- `type: cache` holds Nginx Plus cache zone status
- `type: tcpzone` holds Nginx Plus TCP zone status
- `type: tcpupstream` holds Nginx Plus TCP upstream status

**Sample of Nginx stub status document**

```
{
    "@timestamp": "2015-12-11T13:25:14.667Z",
    "beat": {
        "hostname": "vm-nginxbeat",
        "name": "vm-nginxbeat"
    },
    "count": 1,
    "source": "http://127.0.0.1:8080/status#stub",
    "type": "stub",
    "stub": {
        "accepts": 4,
        "active": 2,
        "current": 4,
        "dropped": 0,
        "handled": 4,
        "reading": 0,
        "requests": 4,
        "waiting": 1,
        "writing": 1
    }
}
```

**Sample of Nginx Plus status document**

```
{
    "@timestamp": "2015-12-11T13:25:14.978Z",
    "beat": {
        "hostname": "vm-nginxbeat",
        "name": "vm-nginxbeat"
    },
    "count": 1,
    "source": "http://demo.nginx.com/status#plus",
    "type": "plus",
    "plus": {
        "address": "206.251.255.64",
        "connections": {
            "accepted": 816383,
            "active": 4,
            "dropped": 0,
            "idle": 24
        },
        "generation": 1,
        "load_timestamp": 1449672299104,
        "nginx_version": "1.9.4",
        "pid": 22709,
        "processes": {
            "respawned": 0
        },
        "requests": {
            "current": 3,
            "total": 2558202
        },
        "ssl": {
            "handshakes": 22321,
            "handshakes_failed": 7985,
            "session_reuses": 4800
        },
        "timestamp": 1449840314663,
        "version": 6
    }
}
```

**Sample of Nginx Plus zone status document**

```
{
    "@timestamp": "2015-12-11T13:25:14.978Z",
    "beat": {
        "hostname": "vm-nginxbeat",
        "name": "vm-nginxbeat"
    },
    "count": 1,
    "source": "http://demo.nginx.com/status#plus",
    "type": "zone",
    "zone": {
        "discarded": 2280,
        "name": "trac.nginx.org",
        "nginx_version": "1.9.4",
        "processing": 0,
        "received": 47296049,
        "requests": 141060,
        "responses": {
            "1xx": 0,
            "2xx": 78270,
            "3xx": 57124,
            "4xx": 3149,
            "5xx": 237,
            "total": 138780
        },
        "sent": 2415440333,
        "version": 6
    }
}
```

**Sample of Nginx Plus upstream status document**

```
{
    "@timestamp": "2015-12-11T13:25:14.978Z",
    "beat": {
        "hostname": "vm-nginxbeat",
        "name": "vm-nginxbeat"
    },
    "count": 1,
    "source": "http://demo.nginx.com/status#plus",
    "type": "upstream",
    "upstream": {
        "keepalive": 0,
        "name": "demo-backend",
        "nginx_version": "1.9.4",
        "peers": [{
            "active": 0,
            "backup": false,
            "downstart": 0,
            "downtime": 0,
            "fails": 0,
            "health_checks": {
                "checks": 167454,
                "fails": 0,
                "last_passed": true,
                "unhealthy": 0
            },
            "id": 0,
            "received": 0,
            "requests": 0,
            "responses": {
                "1xx": 0,
                "2xx": 0,
                "3xx": 0,
                "4xx": 0,
                "5xx": 0,
                "total": 0
            },
            "selected": 0,
            "sent": 0,
            "server": "10.0.0.2:15431",
            "state": "up",
            "unavail": 0,
            "weight": 1
        }],
        "version": 6
    }
}
```

**Sample of Nginx Plus cache status document**

```
{
    "@timestamp": "2015-12-11T13:25:14.978Z",
    "beat": {
        "hostname": "vm-nginxbeat",
        "name": "vm-nginxbeat"
    },
    "count": 1,
    "source": "http://demo.nginx.com/status#plus",
    "type": "cache",
    "cache": {
        "bypass": {
            "bytes": 505700919,
            "bytes_written": 505700919,
            "responses": 10430,
            "responses_written": 10430
        },
        "cold": false,
        "expired": {
            "bytes": 313906093,
            "bytes_written": 309450949,
            "responses": 7591,
            "responses_written": 7308
        },
        "hit": {
            "bytes": 426450804,
            "responses": 39078
        },
        "max_size": 536870912,
        "miss": {
            "bytes": 3385601257,
            "bytes_written": 1344694835,
            "responses": 68405,
            "responses_written": 30029
        },
        "name": "http_cache",
        "nginx_version": "1.9.4",
        "revalidated": {
            "bytes": 0,
            "responses": 0
        },
        "size": 532258816,
        "stale": {
            "bytes": 0,
            "responses": 0
        },
        "updating": {
            "bytes": 0,
            "responses": 0
        },
        "version": 6
    }
}
```

**Sample of Nginx Plus TCP zone status document**

```
{
    "@timestamp": "2015-12-11T13:25:14.978Z",
    "beat": {
        "hostname": "vm-nginxbeat",
        "name": "vm-nginxbeat"
    },
    "count": 1,
    "source": "http://demo.nginx.com/status#plus",
    "type": "tcpzone",
    "tcpzone": {
        "connections": 167454,
        "name": "postgresql_loadbalancer",
        "nginx_version": "1.9.4",
        "processing": 0,
        "received": 17582670,
        "sent": 947252523,
        "version": 6
    }
}
```

**Sample of Nginx Plus TCP upstream status document**

```
{
    "@timestamp": "2015-12-11T13:25:14.978Z",
    "beat": {
        "hostname": "vm-nginxbeat",
        "name": "vm-nginxbeat"
    },
    "count": 1,
    "source": "http://demo.nginx.com/status#plus",
    "type": "tcpupstream",
    "tcpupstream": {
        "name": "postgresql_backends",
        "nginx_version": "1.9.4",
        "peers": [{
            "active": 0,
            "backup": false,
            "connect_time": 1,
            "connections": 55818,
            "downstart": 0,
            "downtime": 0,
            "fails": 0,
            "first_byte_time": 1,
            "health_checks": {
                "checks": 33582,
                "fails": 0,
                "last_passed": true,
                "unhealthy": 0
            },
            "id": 0,
            "max_conns": 42,
            "received": 315750682,
            "response_time": 1,
            "selected": 1449840312000,
            "sent": 5860890,
            "server": "10.0.0.2:15432",
            "state": "up",
            "unavail": 0,
            "weight": 1
        }],
        "version": 6
    }
}
```
