package task

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

// newConfig returns a new aws ConfigProvider
func newConfig(profile string, region string) (aws.Config, error) {
	return config.LoadDefaultConfig(context.Background(), config.WithRegion(region), config.WithSharedConfigProfile(profile))

}

func getenv(value, key string) string {
	if len(value) == 0 {
		return os.Getenv(key)
	}
	return value
}
