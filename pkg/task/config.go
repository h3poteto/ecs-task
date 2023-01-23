package task

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
)

// newConfig returns a new aws ConfigProvider
func newConfig(profile string, region string) (client.ConfigProvider, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: &region,
		},
		Profile: profile,
		SharedConfigState: session.SharedConfigEnable,
	})
	return sess, err
}

func getenv(value, key string) string {
	if len(value) == 0 {
		return os.Getenv(key)
	}
	return value
}
