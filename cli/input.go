package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/danielgtaylor/shorthand"
	yaml "gopkg.in/yaml.v2"
)

// GetBody returns the request body if one was passed either as shorthand
// arguments or via stdin.
func GetBody(mediaType string, args []string) (string, error) {
	var body string

	info, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}
	if (info.Mode() & os.ModeCharDevice) == 0 {
		// Data is available on stdin
		input, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}

		body = string(input)
		LogDebug("Body from stdin is: %s", body)
	}

	if len(args) > 0 {
		bodyInput := strings.Join(args, " ")
		result, err := shorthand.ParseAndBuild("stdin", bodyInput)
		if err != nil {
			return "", err
		}

		if strings.Contains(mediaType, "json") {
			if body != "" {
				// Have a body from stdin. Should be JSON, so let's merge.
				var curBody map[string]interface{}
				if err := json.Unmarshal([]byte(body), &curBody); err != nil {
					return "", err
				}

				shorthand.DeepAssign(curBody, result)
				result = curBody
			}

			marshalled, err := json.Marshal(result)
			if err != nil {
				return "", err
			}

			body = string(marshalled)
		} else if strings.Contains(mediaType, "yaml") {
			if body != "" {
				// Have a body from stdin. Should be YAML, so let's merge.
				var curBody map[string]interface{}
				if err := yaml.Unmarshal([]byte(body), &curBody); err != nil {
					return "", err
				}

				shorthand.DeepAssign(curBody, result)
				result = curBody
			}

			marshalled, err := yaml.Marshal(result)
			if err != nil {
				return "", err
			}

			body = string(marshalled)
		} else {
			return "", fmt.Errorf("Not sure how to marshal %s", mediaType)
		}
	}

	return body, nil
}
