package searchengine

import (
	"context"
	"fmt"

	"github.com/knowledge-search-system/document-processor/config"
	"github.com/knowledge-search-system/document-processor/internal/model"
	searchenginev1 "github.com/knowledge-search-system/search-engine/proto/searchengine/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	grpcClient searchenginev1.SearchServiceClient
	conn       *grpc.ClientConn
}

func NewClient(cfg *config.Config) (*Client, error) {
	conn, err := grpc.NewClient(cfg.SearchEngine.GRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("dial search-engine grpc at %q: %w", cfg.SearchEngine.GRPCAddr, err)
	}

	return &Client{
		grpcClient: searchenginev1.NewSearchServiceClient(conn),
		conn:       conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) IndexChunks(ctx context.Context, chunks []model.Chunk) (int, error) {
	req := &searchenginev1.IndexChunksRequest{
		Chunks: make([]*searchenginev1.Chunk, 0, len(chunks)),
	}

	for _, chunk := range chunks {
		req.Chunks = append(req.Chunks, &searchenginev1.Chunk{
			ChunkId:    chunk.ChunkID,
			DocumentId: chunk.DocumentID,
			FileName:   chunk.FileName,
			PageNumber: chunk.PageNumber,
			Text:       chunk.Text,
		})
	}

	resp, err := c.grpcClient.IndexChunks(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("call search-engine IndexChunks: %w", err)
	}

	return int(resp.GetIndexedCount()), nil
}
