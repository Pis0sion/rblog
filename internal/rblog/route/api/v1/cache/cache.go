package cache

import (
	"context"
	pb "github.com/Pis0sion/rblog/internal/proto/v1"
)

type Cache struct {
	pb.UnimplementedCacheServer
}

func NewCache() pb.CacheServer {
	return &Cache{}
}

func (c *Cache) ListSecrets(context.Context, *pb.ListSecretsRequest) (*pb.ListSecretsResponse, error) {
	return nil, nil
}

func (c *Cache) ListPolicies(context.Context, *pb.ListPoliciesRequest) (*pb.ListPoliciesResponse, error) {
	return nil, nil
}
