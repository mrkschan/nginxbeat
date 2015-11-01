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
