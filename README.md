# microservice-go

Run the command: go run main.go and open navegador in the link,
[http://localhost:9090/docs](http://localhost:9090/docs)
for the opend the documentation of API


For generating the documentation run the command *make swagger* or *GO111MODULE=off swagger generate spec -o ./swagger.yaml --scan-models*
in terminal<br\>


For generate client run the command <br\>
swagger generate client -f ../swagger.yaml -A product-apil
