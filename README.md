# network-access-checker

Application for checking the availability of hosts over the network, the result is issued as a metric for Prometheus. In addition to the application, busybox-extras, kurl, bash are installed in the container.

##### Docker run example:
```
docker run -d -p 8084:2112 -v $(pwd):/appconfig network-access-checker:latest --config.filepath=/appconfig/config.json
```
##### Config example:
```
{
    "resourceList":[
        {
            "name":"Google",
            "host":"www.google.com",
            "ports":["443","80"]
        },
        {
            "name":"Yandex",
            "host":"www.yandex.com",
            "ports":["443"]
        },
        {
            "name":"for testing",
            "host":"www.213y4y.com",
            "ports":["443","80"]
        }
    ],
    "metricsPort":"2112",
    "checkPeriodSeconds":3,
    "checkConnectionTimeout":2
}
```
Metrics endpoint <b>GET /metrics<b>. 
##### Metrics example:
```
monitoring_network_network_access_checker{resourceAddr="www.213y4y.com",resourceName="for testing"} 0
monitoring_network_network_access_checker{resourceAddr="www.google.com",resourceName="Google"} 1
monitoring_network_network_access_checker{resourceAddr="www.yandex.com",resourceName="Yandex"} 1
```

Docker - https://hub.docker.com/repository/docker/buneyev/network-access-checker