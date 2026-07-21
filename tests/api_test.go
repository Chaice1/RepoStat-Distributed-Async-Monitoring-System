package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const address = "http://localhost:28080"

var client = http.Client{
	Timeout: 30 * time.Second,
}

type PingService struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type PingResponse struct {
	Status   string        `json:"status"`
	Services []PingService `json:"services"`
}

type Subscription struct {
	OwnerName string `json:"owner_name"`
	RepoName  string `json:"repo_name"`
}

type RepositoryInfoResponse struct {
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Stars       int64  `json:"stars"`
	Forks       int64  `json:"forks"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type GetSubscriptionsResponse struct {
	Subscriptions []*Subscription `json:"subscriptions"`
}

func waitForAPI(t *testing.T) {
	t.Helper()

	require.Eventually(t, func() bool {
		resp, err := client.Get(address + "/api/ping")
		if err != nil {
			return false
		}
		defer resp.Body.Close()

		return resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusServiceUnavailable
	}, 20*time.Second, 500*time.Millisecond, "api did not become ready")
}

func TestRepositoryInfo_Async(t *testing.T) {
	waitForAPI(t)
	targetURL := "https://github.com/golang/go"

	resp, err := client.Get(address + "/api/repositories/info?url=" + targetURL)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Contains(t, []int{http.StatusOK, http.StatusAccepted}, resp.StatusCode)

	var body RepositoryInfoResponse
	require.Eventually(t, func() bool {
		r, err := client.Get(address + "/api/repositories/info?url=" + targetURL)
		if err != nil || r.StatusCode != http.StatusOK {
			return false
		}
		defer r.Body.Close()
		json.NewDecoder(r.Body).Decode(&body)
		return body.FullName == "golang/go" && body.Stars > 0
	}, 10*time.Second, 1*time.Second, "data not  in processor cache")

	require.Equal(t, "golang/go", body.FullName)
}

func TestSubscriptions_FullFlow(t *testing.T) {
	waitForAPI(t)

	sub := Subscription{
		OwnerName: "google",
		RepoName:  "uuid",
	}

	bodyBytes, _ := json.Marshal(sub)

	resp, err := client.Post(address+"/subscriptions", "application/json", bytes.NewBuffer(bodyBytes))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	require.Eventually(t, func() bool {
		r, err := client.Get(address + "/subscriptions")
		if err != nil {
			return false
		}
		defer r.Body.Close()
		var list GetSubscriptionsResponse
		json.NewDecoder(r.Body).Decode(&list)

		for _, s := range list.Subscriptions {
			if s.RepoName == "uuid" {
				return true
			}
		}
		return false
	}, 5*time.Second, 500*time.Millisecond)

	req, _ := http.NewRequest(http.MethodDelete, address+"/subscriptions/google/uuid", nil)
	respDel, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, respDel.StatusCode)
}

func TestRepositoryInfo_NotFound(t *testing.T) {
	waitForAPI(t)

	badURL := "https://github.com/Chaice1/3721737173217"

	client.Get(address + "/api/repositories/info?url=" + badURL)

	require.Eventually(t, func() bool {
		r, err := client.Get(address + "/api/repositories/info?url=" + badURL)
		if err != nil {
			return false
		}
		defer r.Body.Close()
		return r.StatusCode == http.StatusNotFound
	}, 10*time.Second, 1*time.Second, "api didn't return 404 for unexisting repo")
}

func TestRateLimiter(t *testing.T) {
	waitForAPI(t)

	limitExceeded := false

	for i := 0; i < 15; i++ {
		resp, err := client.Get(address + "/api/ping")
		if err != nil {
			continue
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			limitExceeded = true
			break
		}
		_ = resp.Body.Close()
	}

	require.True(t, limitExceeded, "Api should return 429 for too many req")
}
