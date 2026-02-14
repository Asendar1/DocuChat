package main

import (
	"context"
	"time"

	pb "github.com/Asendar1/DocuChat/scrapper/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ConCloser interface {
	Close() error
}

type DocClient struct {
	client	pb.DocumentProcessorClient
	conn	*grpc.ClientConn
}

func NewDocClient(addr string) (*DocClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &DocClient{
		client: pb.NewDocumentProcessorClient(conn),
		conn:   conn,
	}, nil
}

func (c *DocClient) CallTest(msg string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.client.Test(ctx, &pb.TestReq{Tm: msg})
	if err != nil {
		return "", err
	}

	return resp.GetTm(), nil
}

func (c *DocClient) Close() error {
	return c.conn.Close()
}

type VectorSearchClient struct {
	client	pb.VectorSearchClient
	conn	*grpc.ClientConn
}

func NewVectorSearchClient(addr string) (*VectorSearchClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &VectorSearchClient{
		client: pb.NewVectorSearchClient(conn),
		conn:   conn,
	}, nil
}

// TODO implement vector search client methods

func (c *VectorSearchClient) Close() error {
	return c.conn.Close()
}
