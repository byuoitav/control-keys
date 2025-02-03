package opa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/byuoitav/control-keys/middleware"
	"github.com/gin-gonic/gin"
)

type Client struct {
	URL   string
	Token string
}

type opaResponse struct {
	DecisionID string    `json:"decision_id"`
	Result     opaResult `json:"result"`
}

type opaResult struct {
	Allow bool `json:"allow"`
}

type opaRequest struct {
	Input requestData `json:"input"`
}

type requestData struct {
	APIKey string `json:"api_key"`
	User   string `json:"user"`
	Path   string `json:"path"`
	Method string `json:"method"`
}

func (client *Client) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Initial data
		opaData := opaRequest{
			Input: requestData{
				Path:   c.FullPath(),
				Method: c.Request.Method,
			},
		}
		fmt.Printf("Context Output: %v\n", c.Request.Context().Value("user"))
		fmt.Printf("Context Output: %v\n", c.Request.Context().Value("userBYUID"))

		// use either the user netid for the authorization request or an
		// API key if one was used instead
		if user, ok := c.Request.Context().Value("user").(string); ok {
			opaData.Input.User = user
			fmt.Printf("User Found\n")
		} else if apiKey, ok := middleware.GetAVAPIKey(c.Request.Context()); ok {
			opaData.Input.APIKey = apiKey
		}

		// Prep the request
		oReq, err := json.Marshal(opaData)
		if err != nil {
			fmt.Printf("Error trying to create request to OPA: %s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while contacting authorization server"})
			return
		}

		req, err := http.NewRequest(
			"POST",
			fmt.Sprintf("%s/v1/data/viaalert", client.URL),
			bytes.NewReader(oReq),
		)

		req.Header.Set("authorization", fmt.Sprintf("Bearer %s", client.Token))

		fmt.Printf("Data: %s\n", opaData)
		fmt.Printf("URL: %s\n", client.URL)

		// Make the request
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("Error while making request to OPA: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while contacting authorization server"})
			return
		}
		if res.StatusCode != http.StatusOK {
			fmt.Printf("Got back non 200 status from OPA: %d", res.StatusCode)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while contacting authorization server"})
			return
		}

		// Read the body
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Unable to read body from OPA: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while contacting authorization server"})
			return
		}

		// Unmarshal the body
		oRes := opaResponse{}
		err = json.Unmarshal(body, &oRes)
		if err != nil {
			fmt.Printf("Unable to parse body from OPA: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while contacting authorization server"})
			return
		}
		fmt.Printf("Results: %v\n", oRes.Result)
		// If OPA approved then allow the request, else reject with a 403
		if oRes.Result.Allow {
			c.Next()
		} else {
			fmt.Printf("Unauthorized\n")
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
			c.Abort()
		}
	}
}
