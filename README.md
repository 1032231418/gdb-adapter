gdb adapter [![Build Status](https://travis-ci.org/vance-liu/gdb-adapter.svg?branch=master)](https://travis-ci.org/vance-liu/gdb-adapter) [![Coverage Status](https://coveralls.io/repos/github/vance-liu/gdb-adapter/badge.svg?branch=master)](https://coveralls.io/github/vance-liu/gdb-adapter?branch=master) [![Godoc](https://godoc.org/github.com/vance-liu/gdb-adapter?status.svg)](https://godoc.org/github.com/vance-liu/gdb-adapter)
====

[GF ORM](https://github.com/gogf/gf) adapter for [Casbin](https://github.com/casbin/casbin). 

Based on [GF ORM](https://github.com/gogf/gf), and tested in:
- MySQL
- PostgreSQL

基于 [vance-liu/gdb-adapter](https://github.com/vance-liu/gdb-adapter)

## Installation

    go get github.com/flipped-aurora/gdb-adapter

## Usage example

```go
db, err = gdb.New("casbin")
if err != nil{
	panic(err)
}
a := NewAdapterByDB(db)
e := casbin.NewEnforcer("examples/rbac_model.conf", a)
```

## Notice

- you should create the database on your own.

-  Not tested, please do not use in production environment.

## Getting Help

- [Casbin](https://github.com/casbin/casbin)

## License

This project is under Apache 2.0 License. See the [LICENSE](LICENSE) file for the full license text.

