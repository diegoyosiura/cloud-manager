package utils

import (
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/emaildataplane"
	"log"
	"os"
	"regexp"
	"strings"
)

// GetEnvWithValidation retrieves an environment variable and ensures it is not empty.
func GetEnvWithValidation(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}
	return value
}

// GetOptionalEnv retrieves an environment variable or returns the provided default value if the variable is not set.
func GetOptionalEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func IsValidArn(arn string) bool {
	return strings.HasPrefix(arn, "arn:aws:")
}

func ConvertToOCIEmailList(l []string) []emaildataplane.EmailAddress {
	re := regexp.MustCompile(`(?P<email>[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})\s*(?:<(?P<name>[^>]+)>)?`)

	var r []emaildataplane.EmailAddress

	for _, e := range l {
		match := re.FindStringSubmatch(e)
		email := match[re.SubexpIndex("email")]
		name := match[re.SubexpIndex("name")]

		r = append(r, emaildataplane.EmailAddress{Email: common.String(email), Name: common.String(name)})
	}

	return r
}
