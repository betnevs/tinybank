# tinybank
A practice project implemented by Go.

## Install
1. Deploy mysql with [Docker Desktop](https://www.docker.com/products/docker-desktop/).
``` makefile
make mysql
```
2. Create database `tinybank`.
``` makefile
make createdb
```
3. Create tables.
``` makefile
make migrateup
```