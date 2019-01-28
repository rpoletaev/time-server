# time-server
implementation of rfc868. TCP only

### build:

  ```GO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .```
  
### run:

  ```./time-server -p 11037```
