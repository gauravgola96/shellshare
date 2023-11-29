
# ShellShare


### Run SSH Server

```go run cmd/shellshare/main.go ssh ```


### Run HTTP Server

```go run cmd/shellshare/main.go http ```

### Run SSH & HTTP Server combined as single Service - Preffered

```go run cmd/shellshare/main.go combined ```


## How to use locally

> run
> : ssh localhost -p 2222 < `<filepath>`

click on the link to download

> Your download link: http://localhost:8000/download/0188c94d-91d4-73ad-b901-b936c2678458

### Example Commands
>
> ssh localhost -p 2222 < file.yaml
> 
> ssh localhost -p 2222 filename="file.yaml" < file.yaml