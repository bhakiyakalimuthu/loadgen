# loadgen
>*  loadgen is a CLI client which generates json rpc load/request and send batch message every interval to the configured endpoint/url(configurable).
>
>*   CLI contains a library that implements an HTTP client.A client is configured with a URL to which request are sent.It implements a function that takes json rpc tx inputs and load them by sending HTTP POST requests to the configured URL with the rpc request content in the request body.
client operations are non-blocking for the caller. It handles notification failures.

### Supported operation
```
Usage:
  notifier [flags]

Flags:
  -e, --env string          test or staging environment (default "dev")
  -h, --help                help for loadgen
  -i, --interval duration   Load interval (default 1m40s)
  -n, --numOfBatch int      Number of batch request / Each batch contains min 10 requests (default 1)
  -u, --url string          URL to which load to generated / request to be sent
  ```

### Architecture diagram
![plot](picture/Architecture_diagram.png)


# Pre requisites
- Ubuntu 20.04 (any linux based distros) 
- Vim or Goland 2021.3.1
- Go 1.17

# Build & Run
> loadgen library  is designed to run as a cli which means all business logics as single application.
* Application can be build and started by using Makefile.
* Make sure to cd to project folder.
* Run the below commands in the terminal shell.
* Make sure to run Pre-run and Go path is set properly
* Make sure to install go static check tool (go install honnef.co/go/tools/cmd/staticcheck@latest)
# Pre-run
    make mod
    make lint 

# How to run unit test
    make test

# How to run build
    make build

# How to start the cli
- **1.start using go build (Without go install)**
![plot](picture/go-build.png)
- **2.start using go install (go path is configured properly)**
![plot](picture/go-install.png)

# How to exit the cli
    press ctrl+c 
