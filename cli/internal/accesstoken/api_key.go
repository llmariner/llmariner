package accesstoken

import "os"

// GetAPIKeyEnvVar returns the value of the LLMARINER_API_KEY environment variable.
func GetAPIKeyEnvVar() string {
	return os.Getenv("LLMARINER_API_KEY")
}
