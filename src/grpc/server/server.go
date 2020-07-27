package main

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "grpc/pb"
	"os"
)

type server struct {
	jsonFile string
	fr *os.File
	dec *json.Decoder
}

type customer struct {
	CustomerId int `json:"customer_id"`
	CustomerName string `json:"customer_name"`
}

func (serve server) GetAllCustomers(ctx context.Context, apiType *pb.ApiType) (*pb.ListCustomers, error) {
	if apiType.Type != 1 {
		panic("Wrong Type received in grpc call")
	}

	fr, err := os.Open(serve.jsonFile)
	checkError(err)
	defer fr.Close()

	dec := json.NewDecoder(fr)

	_, err = dec.Token()
	checkError(err)

	customers := make([] *pb.Customer, 0, 5)
	for dec.More() {
		var tmp customer
		err := dec.Decode(&tmp)
		checkError(err)
		tt := pb.Customer{CustomerId: int64(tmp.CustomerId), CustomName: tmp.CustomerName}
		customers = append(customers, &tt)
	}
	//fmt.Println("CUSTOMERS ---> ", customers)
	return &pb.ListCustomers{Customers: customers}, nil
}

func (serve server) PutCustomer(ctx context.Context, apiType *pb.ApiType) (*pb.ApiResponse, error) {
	if apiType.Type != 2 {
		panic("Wrong Type received in grpc call")
	}

	cust := customer{CustomerId: int(apiType.CustomerData.CustomerId), CustomerName: apiType.CustomerData.CustomName}

	f, err := os.OpenFile(serve.jsonFile, os.O_RDWR, os.ModePerm)
	defer f.Close()
	checkError(err)

	custJson, err := json.Marshal(cust)
	checkError(err)

	custString := string(custJson)
	custString = "," + custString

	off := int64(1)
	stat, err := os.Stat(serve.jsonFile)
	fmt.Println("Size : ", stat.Size())
	start := stat.Size() - off

	tmp := []byte(custString)
	_, err = f.WriteAt(tmp, start)
	checkError(err)

	str := []byte("]")
	_, err = f.WriteAt(str, start + int64(len(custString)))
	checkError(err)

	return &pb.ApiResponse{StatusCode: 200}, nil

}

func (serve server) GetCustomer(ctx context.Context, apiType *pb.ApiType) (*pb.Customer, error) {
	if apiType.Type != 3 {
		panic("Wrong Type received in grpc call")
	}
	return &pb.Customer{CustomerId: 0, CustomName: "Testing"}, nil
}

func (serve *server) init(customerFilePath string) {
	serve.jsonFile = customerFilePath

	var err error
	serve.fr, err = os.Open(customerFilePath)
	checkError(err)

	serve.dec = json.NewDecoder(serve.fr)

	_, err = serve.dec.Token()
	checkError(err)
}

func checkError(err error)  {
	if err != nil {
		log.Fatal(err.Error())
	}
}

var serve server

func main() {
	serve.init("customer.json")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterApiServiceServer(s, &serve)

	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen, err is %v", err)
	}

}
