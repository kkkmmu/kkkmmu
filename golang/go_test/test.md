import testing
go help test
go help testflag
go test -cover
go test -v -cover -coverprofile=c.out
go tool cover -html c.out
