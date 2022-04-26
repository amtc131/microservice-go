# microservice-go

Run the command: go run main.go and open navegador in the link,
[http://localhost:9090/docs](http://localhost:9090/docs)
for the opend the documentation of API


For generating the documentation run the command *make swagger* or *GO111MODULE=off swagger generate spec -o ./swagger.yaml --scan-models*
in terminal<br\>


For generate client run the command <br\>
swagger generate client -f ../swagger.yaml -A product-apil


> Update rates, change *JPY*

curl localhost:9090/products?currency=JPY


Testing

To test the system install grpccurl which is a command line tool which can interact with gRPC API's

List Services

grpcurl --plaintext localhost:9092 list
Currency
grpc.reflection.v1alpha.ServerReflection


List Methods

grpcurl --plaintext localhost:9092 list Currency        
Currency.GetRate
Currency.SubscribeRates

Method detail for GetRate

grpcurl --plaintext localhost:9092 describe Currency.GetRate

Currency.GetRate is a method:
rpc GetRate ( .RateRequest ) returns ( .RateResponse );


RateRequest detail

grpcurl --plaintext --msg-template localhost:9092 describe .RateRequest    
RateRequest is a message:
message RateRequest {
  string Base = 1 [json_name = "base"];
  string Destination = 2 [json_name = "destination"];
}

Message template:
{
  "Base": "EUR",
  "Destination": "EUR"
}


Execute a request for GetRate

➜ grpcurl --plaintext -d '{"Base": "GBP", "Destination": "USD"}' localhost:9092 Currency/GetRate
{
  "rate": 1.2229967868538965
}


Execute a request for SubscribeRates
The parameter -d @ means that gRPCurl will read the messages from StdIn.

grpcurl --plaintext --msg-template -d @ localhost:9092 Currency/SubscribeRates 
You can send a message to the server using the following payload

{
  "Base": "EUR",
  "Destination": "GBP"
}






