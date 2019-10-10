package main

import (
	"context"
	"database/sql"
	"fmt"
	pb "github.com/Calvinn097/pmb-grpc-server/grpc/account/account"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	port             = ":50051"
	database_setting = "root:root@tcp(user_mysql:3306)/user_db"
)

type accountServer struct {
}

type User struct {
	Id      int64
	Name    string
	Address string
	Age     int32
}

func (s *accountServer) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	name := request.GetName()
	address := request.GetAddress()
	age := request.GetAge()
	db, err := sql.Open("mysql", database_setting)

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO users(name, address, age) VALUES(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(name, address, age)
	if err != nil {
		log.Fatal(err)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	return &pb.CreateUserResponse{Id: lastId, Address: address, Name: name, Age: age}, nil
}

func (s *accountServer) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	id := request.GetId()
	/**name := request.GetName()
	address := request.GetAddress()
	age := request.GetAge()**/

	return &pb.UpdateUserResponse{Id: id, Updated: true}, nil
}

func (s *accountServer) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	id := request.GetId()

	db, err := sql.Open("mysql", database_setting)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var user User

	err = db.QueryRow("SELECT * FROM users where id = ?", id).Scan(&user.Id, &user.Name, &user.Address, &user.Age)
	if err != nil {
		log.Fatal(err)
	}

	return &pb.GetUserResponse{Id: user.Id, Name: user.Name, Address: user.Address, Age: user.Age}, nil
}

func (s *accountServer) DeleteUser(ctx context.Context, request *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	id := request.GetId()

	db, err := sql.Open("mysql", database_setting)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	del, err := db.Prepare("DELETE FROM users where id=?")
	if err != nil {
		log.Fatal(err)
	}
	del.Exec(id)

	return &pb.DeleteUserResponse{Deleted: true}, nil
}

func main() {

	fmt.Println("running")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAccountServer(s, &accountServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
