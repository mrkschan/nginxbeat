[![Build Status](https://travis-ci.org/mrkschan/nginxbeat.svg?branch=travisci-setup)](https://travis-ci.org/mrkschan/nginxbeat)

# nginxbeat

Nginxbeat is the Beat used for Nginx monitoring. It is a lightweight agent that reads status from Nginx periodically. Nginx must either expose its status via stub module (http://nginx.org/en/docs/http/ngx_http_stub_status_module.html) or Nginx Plus status module (http://nginx.org/en/docs/http/ngx_http_status_module.html).


## Elasticsearch template

To apply nginxbeat template for Nginx stub status:

```
curl -XPUT 'http://localhost:9200/_template/nginxbeat' -d@etc/nginxbeat-stub.template.json
```

To apply nginxbeat template for Nginx Plus status:

```
curl -XPUT 'http://localhost:9200/_template/nginxbeat' -d@etc/nginxbeat-plus.template.json
```


# Build, Test, Run

```
# Build
export GO15VENDOREXPERIMENT=1
GOPATH=<your go path> make

# Test
GOPATH=<your go path> make test

# Run
./nginxbeat -c etc/nginxbeat.yml
```


## Exported fields

Nginxbeat only exports a single type of document. Though, the properties in the document varies according to the configured Nginx status page.

- `type: nginx` holds either Nginx stub status or Nginx Plus status
- `type: zone` holds Nginx status zone status
- `type: upstream` holds Nginx upstream group status
- `type: cache` holds Nginx cache zone status
- `type: tcpzone` holds Nginx TCP zone status
- `type: tcpupstream` holds Nginx TCP upstream status

**Sample of Nginx stub status document**

```
{
    "count": 1,
    "nginx": {
        "format": "stub",
        "accepts": 10660,
        "active": 443,
        "current": 333,
        "dropped": 0,
        "handled": 10660,
        "reading": 212,
        "requests": 16882,
        "waiting": 110,
        "writing": 121
    },
    "shipper": "vm-nginxbeat",
    "@timestamp": "2015-11-01T14:25:42.776Z",
    "type": "nginx"
}
```

**Sample of Nginx Plus status document**

```
{
    "@timestamp": "2015-11-25T14:55:50.396Z",
    "count": 1,
    "nginx": {
        "format": "plus",
        "address": "206.251.255.64",
        "connections": {
            "accepted": 20293854,
            "active": 7,
            "dropped": 0,
            "idle": 32
        },
        "generation": 16,
        "load_timestamp": 1448100000394,
        "nginx_version": "1.9.4",
        "pid": 70469,
        "processes": {
            "respawned": 0
        },
        "requests": {
            "current": 7,
            "total": 45834498
        },
        "ssl": {
            "handshakes": 61215,
            "handshakes_failed": 7429,
            "session_reuses": 7799
        },
        "timestamp": 1448463350114,
        "version": 6
    },
    "shipper": "vm-nginxbeat",
    "type": "nginx"
}

{
    "@timestamp": "2015-11-25T14:55:50.396Z",
    "count": 1,
    "shipper": "vm-nginxbeat",
    "type": "zone",
    "zone": {
        "format": "plus",
        "discarded": 73,
        "name": "hg.nginx.org",
        "nginx_version": "1.9.4",
        "processing": 0,
        "received": 29203365,
        "requests": 94029,
        "responses": {
            "1xx": 0,
            "2xx": 87003,
            "3xx": 5204,
            "4xx": 1217,
            "5xx": 532,
            "total": 93956
        },
        "sent": 4108440100,
        "version": 6
    }
}

{
    "@timestamp": "2015-11-25T14:55:50.396Z",
    "count": 1,
    "shipper": "vm-nginxbeat",
    "type": "upstream",
    "upstream": {
        "format": "plus",
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
                "checks": 361666,
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

{
    "@timestamp": "2015-11-25T14:55:50.396Z",
    "cache": {
        "format": "plus",
        "bypass": {
            "bytes": 7630800604,
            "bytes_written": 7630793856,
            "responses": 172960,
            "responses_written": 172931
        },
        "cold": false,
        "expired": {
            "bytes": 9666071646,
            "bytes_written": 9572182461,
            "responses": 128922,
            "responses_written": 124066
        },
        "hit": {
            "bytes": 493360022371,
            "responses": 790197
        },
        "max_size": 536870912,
        "miss": {
            "bytes": 56916438362,
            "bytes_written": 26602168664,
            "responses": 1509079,
            "responses_written": 928035
        },
        "name": "http_cache",
        "nginx_version": "1.9.4",
        "revalidated": {
            "bytes": 0,
            "responses": 0
        },
        "size": 530309120,
        "stale": {
            "bytes": 0,
            "responses": 0
        },
        "updating": {
            "bytes": 0,
            "responses": 0
        },
        "version": 6
    },
    "count": 1,
    "shipper": "vm-nginxbeat",
    "type": "cache"
}

{
    "@timestamp": "2015-11-25T14:55:50.396Z",
    "count": 1,
    "shipper": "vm-nginxbeat",
    "tcpzone": {
        "format": "plus",
        "connections": 361666,
        "name": "postgresql_loadbalancer",
        "nginx_version": "1.9.4",
        "processing": 0,
        "received": 37974930,
        "sent": 2061658911,
        "version": 6
    },
    "type": "tcpzone"
}

{
    "@timestamp": "2015-11-25T14:55:50.396Z",
    "count": 1,
    "shipper": "vm-nginxbeat",
    "tcpupstream": {
        "format": "plus",
        "name": "postgresql_backends",
        "nginx_version": "1.9.4",
        "peers": [{
            "active": 0,
            "backup": false,
            "connect_time": 1,
            "connections": 120556,
            "downstart": 0,
            "downtime": 0,
            "fails": 0,
            "first_byte_time": 1,
            "health_checks": {
                "checks": 72607,
                "fails": 0,
                "last_passed": true,
                "unhealthy": 0
            },
            "id": 0,
            "max_conns": 42,
            "received": 687223188,
            "response_time": 1,
            "selected": 1448463349000,
            "sent": 12658380,
            "server": "10.0.0.2:15432",
            "state": "up",
            "unavail": 0,
            "weight": 1
        }, {
            "active": 0,
            "backup": false,
            "connect_time": 1,
            "connections": 120555,
            "downstart": 0,
            "downtime": 0,
            "fails": 0,
            "first_byte_time": 1,
            "health_checks": {
                "checks": 72607,
                "fails": 0,
                "last_passed": true,
                "unhealthy": 0
            },
            "id": 1,
            "received": 687217775,
            "response_time": 1,
            "selected": 1448463347000,
            "sent": 12658275,
            "server": "10.0.0.2:15433",
            "state": "up",
            "unavail": 0,
            "weight": 1
        }, {
            "active": 0,
            "backup": false,
            "connect_time": 1,
            "connections": 120555,
            "downstart": 0,
            "downtime": 0,
            "fails": 0,
            "first_byte_time": 1,
            "health_checks": {
                "checks": 72607,
                "fails": 0,
                "last_passed": true,
                "unhealthy": 0
            },
            "id": 2,
            "received": 687217948,
            "response_time": 1,
            "selected": 1448463348000,
            "sent": 12658275,
            "server": "10.0.0.2:15434",
            "state": "up",
            "unavail": 0,
            "weight": 1
        }, {
            "active": 0,
            "backup": false,
            "connections": 0,
            "downstart": 0,
            "downtime": 0,
            "fails": 0,
            "health_checks": {
                "checks": 0,
                "fails": 0,
                "unhealthy": 0
            },
            "id": 3,
            "received": 0,
            "selected": 0,
            "sent": 0,
            "server": "10.0.0.2:15435",
            "state": "down",
            "unavail": 0,
            "weight": 1
        }],
        "version": 6
    },
    "type": "tcpupstream"
}
```
