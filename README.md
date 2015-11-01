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


## Exported fields

Nginxbeat only exports a single type of document. Though, the properties in the document varies according to the configured Nginx status page.

- `type: nginx` holds either Nginx stub status or Nginx Plus status

**Sample of Nginx stub status document**

```
{
    "count": 1,
    "nginx": {
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
    "timestamp": "2015-11-01T14:25:42.776Z",
    "type": "nginx"
}
```

**Sample of Nginx Plus status document**

```
{
    "count": 1,
    "nginx": {
        "address": "206.251.255.64",
        "caches": [{
            "bypass": {
                "bytes": 2683460584,
                "bytes_written": 2683456136,
                "responses": 60927,
                "responses_written": 60911
            },
            "cold": false,
            "expired": {
                "bytes": 2399255970,
                "bytes_written": 2372163959,
                "responses": 44397,
                "responses_written": 43052
            },
            "hit": {
                "bytes": 19273543372,
                "responses": 222198
            },
            "max_size": 536870912,
            "miss": {
                "bytes": 19611112157,
                "bytes_written": 9580189053,
                "responses": 551535,
                "responses_written": 336504
            },
            "name": "http_cache",
            "revalidated": {
                "bytes": 0,
                "responses": 0
            },
            "size": 536276992,
            "stale": {
                "bytes": 0,
                "responses": 0
            },
            "updating": {
                "bytes": 0,
                "responses": 0
            }
        }],
        "connections": {
            "accepted": 9601253,
            "active": 2,
            "dropped": 0,
            "idle": 18
        },
        "generation": 12,
        "load_timestamp": 1446285600278,
        "nginx_version": "1.9.4",
        "pid": 92677,
        "processes": {
            "respawned": 0
        },
        "requests": {
            "current": 2,
            "total": 20097683
        },
        "server_zones": [{
            "discarded": 1732,
            "name": "hg.nginx.org",
            "processing": 0,
            "received": 9404252,
            "requests": 33082,
            "responses": {
                "1xx": 0,
                "2xx": 29893,
                "3xx": 857,
                "4xx": 453,
                "5xx": 147,
                "total": 31350
            },
            "sent": 943707757
        }, {
            "discarded": 396,
            "name": "trac.nginx.org",
            "processing": 1,
            "received": 16824065,
            "requests": 47339,
            "responses": {
                "1xx": 0,
                "2xx": 22694,
                "3xx": 21060,
                "4xx": 3061,
                "5xx": 127,
                "total": 46942
            },
            "sent": 771492649
        }, {
            "discarded": 112,
            "name": "lxr.nginx.org",
            "processing": 0,
            "received": 888931,
            "requests": 3684,
            "responses": {
                "1xx": 0,
                "2xx": 3361,
                "3xx": 125,
                "4xx": 80,
                "5xx": 6,
                "total": 3572
            },
            "sent": 82231555
        }],
        "ssl": {
            "handshakes": 2763,
            "handshakes_failed": 348,
            "session_reuses": 590
        },
        "stream": {
            "server_zones": [{
                "connections": 101119,
                "name": "postgresql_loadbalancer",
                "processing": 0,
                "received": 10617495,
                "sent": 571476342
            }],
            "upstreams": [{
                "name": "postgresql_backends",
                "peers": [{
                    "active": 0,
                    "backup": false,
                    "connect_time": 1,
                    "connections": 33707,
                    "downstart": 0,
                    "downtime": 0,
                    "fails": 0,
                    "first_byte_time": 1,
                    "health_checks": {
                        "checks": 20320,
                        "fails": 0,
                        "last_passed": true,
                        "unhealthy": 0
                    },
                    "id": 0,
                    "max_conns": 42,
                    "received": 190495640,
                    "response_time": 1,
                    "selected": 1446387304000,
                    "sent": 3539235,
                    "server": "10.0.0.2:15432",
                    "state": "up",
                    "unavail": 0,
                    "weight": 1
                }, {
                    "active": 0,
                    "backup": false,
                    "connect_time": 1,
                    "connections": 33706,
                    "downstart": 0,
                    "downtime": 0,
                    "fails": 0,
                    "first_byte_time": 1,
                    "health_checks": {
                        "checks": 20320,
                        "fails": 0,
                        "last_passed": true,
                        "unhealthy": 0
                    },
                    "id": 1,
                    "received": 190490305,
                    "response_time": 1,
                    "selected": 1446387302000,
                    "sent": 3539130,
                    "server": "10.0.0.2:15433",
                    "state": "up",
                    "unavail": 0,
                    "weight": 1
                }, {
                    "active": 0,
                    "backup": false,
                    "connect_time": 1,
                    "connections": 33706,
                    "downstart": 0,
                    "downtime": 0,
                    "fails": 0,
                    "first_byte_time": 1,
                    "health_checks": {
                        "checks": 20320,
                        "fails": 0,
                        "last_passed": true,
                        "unhealthy": 0
                    },
                    "id": 2,
                    "received": 190490397,
                    "response_time": 1,
                    "selected": 1446387303000,
                    "sent": 3539130,
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
                }]
            }]
        },
        "timestamp": 1446387304962,
        "upstreams": [{
            "keepalive": 0,
            "name": "trac-backend",
            "peers": [{
                "active": 0,
                "backup": false,
                "downstart": 0,
                "downtime": 0,
                "fails": 0,
                "health_checks": {
                    "checks": 10152,
                    "fails": 0,
                    "last_passed": true,
                    "unhealthy": 0
                },
                "id": 0,
                "received": 692128923,
                "requests": 17301,
                "responses": {
                    "1xx": 0,
                    "2xx": 14968,
                    "3xx": 630,
                    "4xx": 1702,
                    "5xx": 1,
                    "total": 17301
                },
                "selected": 1446387293000,
                "sent": 7430064,
                "server": "10.0.0.1:8080",
                "state": "up",
                "unavail": 0,
                "weight": 1
            }, {
                "active": 0,
                "backup": true,
                "downstart": 1446285601039,
                "downtime": 101703923,
                "fails": 0,
                "health_checks": {
                    "checks": 10166,
                    "fails": 10166,
                    "last_passed": false,
                    "unhealthy": 1
                },
                "id": 1,
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
                "server": "10.0.0.1:8081",
                "state": "unhealthy",
                "unavail": 0,
                "weight": 1
            }]
        }, {
            "keepalive": 0,
            "name": "hg-backend",
            "peers": [{
                "active": 0,
                "backup": false,
                "downstart": 0,
                "downtime": 0,
                "fails": 0,
                "health_checks": {
                    "checks": 10139,
                    "fails": 0,
                    "last_passed": true,
                    "unhealthy": 0
                },
                "id": 0,
                "received": 961850742,
                "requests": 29623,
                "responses": {
                    "1xx": 0,
                    "2xx": 29168,
                    "3xx": 0,
                    "4xx": 453,
                    "5xx": 2,
                    "total": 29623
                },
                "selected": 1446387287000,
                "sent": 8740856,
                "server": "10.0.0.1:8088",
                "state": "up",
                "unavail": 0,
                "weight": 5
            }, {
                "active": 0,
                "backup": true,
                "downstart": 1446285601018,
                "downtime": 101703944,
                "fails": 0,
                "health_checks": {
                    "checks": 10166,
                    "fails": 10166,
                    "last_passed": false,
                    "unhealthy": 1
                },
                "id": 1,
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
                "server": "10.0.0.1:8089",
                "state": "unhealthy",
                "unavail": 0,
                "weight": 1
            }]
        }, {
            "keepalive": 0,
            "name": "lxr-backend",
            "peers": [{
                "active": 0,
                "backup": false,
                "downstart": 0,
                "downtime": 0,
                "fails": 0,
                "health_checks": {
                    "checks": 0,
                    "fails": 0,
                    "unhealthy": 0
                },
                "id": 0,
                "received": 81183944,
                "requests": 2980,
                "responses": {
                    "1xx": 0,
                    "2xx": 2980,
                    "3xx": 0,
                    "4xx": 0,
                    "5xx": 0,
                    "total": 2980
                },
                "selected": 1446387286000,
                "sent": 2087704,
                "server": "unix:/tmp/cgi.sock",
                "state": "up",
                "unavail": 0,
                "weight": 1
            }, {
                "active": 0,
                "backup": true,
                "downstart": 0,
                "downtime": 0,
                "fails": 0,
                "health_checks": {
                    "checks": 0,
                    "fails": 0,
                    "unhealthy": 0
                },
                "id": 1,
                "max_conns": 42,
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
                "server": "unix:/tmp/cgib.sock",
                "state": "up",
                "unavail": 0,
                "weight": 1
            }]
        }, {
            "keepalive": 0,
            "name": "demo-backend",
            "peers": [{
                "active": 0,
                "backup": false,
                "downstart": 0,
                "downtime": 0,
                "fails": 0,
                "health_checks": {
                    "checks": 101119,
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
            }]
        }],
        "version": 6
    },
    "shipper": "vm-nginxbeat",
    "timestamp": "2015-11-01T14:15:06.166Z",
    "type": "nginx"
}
```

