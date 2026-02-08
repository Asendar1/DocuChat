package main

import (
	"context"
	"time"

	pb "github.com/Asendar1/DocuChat/scrapper/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TestClient struct {
	client pb.TestClient
	conn   *grpc.ClientConn
}

func NewTestClient(addr string) (*TestClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &TestClient{
		client: pb.NewTestClient(conn),
		conn:   conn,
	}, nil
}

func (c *TestClient) CallTest(msg string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.client.Test(ctx, &pb.TestReq{Tm: msg})
	if err != nil {
		return "", err
	}

	return resp.GetTm(), nil
}

func (c *TestClient) Close() error {
	return c.conn.Close()
}
