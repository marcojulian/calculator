# calculator
Exercise project written in Go that I did on my own during the course "gRPC [Golang] Master Class: Build Modern API &amp; Microservices" taught by Stephane Maarek on Udemy. 

This project is a demostration of the things I learned while doing the course. Here I show the four types of comunication between server and client (unary, server
streaming, client streaming, and bidirectional streaming), alongside with some other topics like the errors in gRPC, deadlines, reflection, and ssl encryption.

### Generating gRPC files
To generate the files related to gRPC you just need to run the contents of `generate.sh` or just run the file with sh or bash.
```
$ sh generate.sh
```

### Reflection
If you want to explore the service run [evans](https://github.com/ktr0731/evans) with this command while running the server.
```
$ evans -t --cacert ssl/ca.crt --servername localhost -r -p 50051
```

### Running the code
You will need two terminals for both the server and the client and run this on each of them:
```
go run server/server.go
```
```
go run client/client.go
```
