# outersystems/http-exporter

Reverse proxy with exporter for prometheus. Project to learn golang. Please be gentle but please criticize.

# Options
    Usage of ./http-exporter:
      -metrics string
        	Port on which the metrics will be available (default ":9696")
      -port string
        	Listening port (default ":8080")
      -target string
        	Target (default "http://127.0.0.1:8081")

# Build
    go get github.com/prometheus/client_golang/prometheus
    go build
    docker build .

The steps are handled by the (over-engeenired) script ```./make.sh```.

By default a call to make.sh compiles and creates the docker image. It's possible to ask for just the compilation with ```./make.sh compile```.

The compilation is done inside a container and I use the image of my IDE (vim) because the default golang images are not really maintained.

# Test

The ```docker-compose.yml``` and ```docker-compose.override.yml``` files are starting the http-exporter and a webserver.
