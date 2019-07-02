# shorten-go

A version of the original [shorten](https://github.com/amdavidson/shorten) re-written in Go and made Docker friendly.

## Run shorten-go in Docker

    docker run -d --restart always --v shorten-data:/data -p 80:8000 amdavidson/shorten-go
  
## Build shorten-go and run

    git clone https://github.com/amdavidson/shorten-go.git
    cd shorten-go
    go get -d 
    go build -o shorten-go .
    ./shorten-go
