package agent

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHTTP_AgentSelf(t *testing.T) {
	httpTest(t, nil, func(s *TestServer) {
		// Make the HTTP request
		req, err := http.NewRequest("GET", "/v1/agent/self", nil)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		respW := httptest.NewRecorder()

		// Make the request
		obj, err := s.Server.AgentSelfRequest(respW, req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		// Check the job
		self := obj.(agentSelf)
		if self.Config == nil {
			t.Fatalf("bad: %#v", self)
		}
		if len(self.Stats) == 0 {
			t.Fatalf("bad: %#v", self)
		}
	})
}

func TestHTTP_AgentJoin(t *testing.T) {
	httpTest(t, nil, func(s *TestServer) {
		// Determine the join address
		member := s.Agent.Server().LocalMember()
		addr := fmt.Sprintf("%s:%d", member.Addr, member.Port)

		// Make the HTTP request
		req, err := http.NewRequest("PUT",
			fmt.Sprintf("/v1/agent/join?address=%s&address=%s", addr, addr), nil)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		respW := httptest.NewRecorder()

		// Make the request
		obj, err := s.Server.AgentJoinRequest(respW, req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		// Check the job
		join := obj.(joinResult)
		if join.NumJoined != 2 {
			t.Fatalf("bad: %#v", join)
		}
		if join.Error != "" {
			t.Fatalf("bad: %#v", join)
		}
	})
}

func TestHTTP_AgentMembers(t *testing.T) {
	httpTest(t, nil, func(s *TestServer) {
		// Make the HTTP request
		req, err := http.NewRequest("GET", "/v1/agent/members", nil)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		respW := httptest.NewRecorder()

		// Make the request
		obj, err := s.Server.AgentMembersRequest(respW, req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		// Check the job
		members := obj.([]Member)
		if len(members) != 1 {
			t.Fatalf("bad: %#v", members)
		}
	})
}

func TestHTTP_AgentForceLeave(t *testing.T) {
	httpTest(t, nil, func(s *TestServer) {
		// Make the HTTP request
		req, err := http.NewRequest("PUT", "/v1/agent/force-leave?node=foo", nil)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		respW := httptest.NewRecorder()

		// Make the request
		_, err = s.Server.AgentForceLeaveRequest(respW, req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
	})
}

func TestHTTP_AgentSetServers(t *testing.T) {
	httpTest(t, nil, func(s *TestServer) {
		// Establish a baseline number of servers
		req, err := http.NewRequest("GET", "/v1/agent/servers", nil)
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		respW := httptest.NewRecorder()

		// Create the request
		req, err = http.NewRequest("PUT", "/v1/agent/servers", nil)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		// Send the request
		respW = httptest.NewRecorder()
		_, err = s.Server.AgentServersRequest(respW, req)
		if err == nil || !strings.Contains(err.Error(), "missing server address") {
			t.Fatalf("expected missing servers error, got: %#v", err)
		}

		// Create a valid request
		req, err = http.NewRequest("PUT", "/v1/agent/servers?address=127.0.0.1%3A4647&address=127.0.0.2%3A4647&address=127.0.0.3%3A4647", nil)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		// Send the request
		respW = httptest.NewRecorder()
		_, err = s.Server.AgentServersRequest(respW, req)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		// Retrieve the servers again
		req, err = http.NewRequest("GET", "/v1/agent/servers", nil)
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		respW = httptest.NewRecorder()

		// Make the request and check the result
		expected := map[string]bool{
			"127.0.0.1:4647": true,
			"127.0.0.2:4647": true,
			"127.0.0.3:4647": true,
		}
		out, err := s.Server.AgentServersRequest(respW, req)
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		servers := out.([]string)
		if n := len(servers); n != len(expected) {
			t.Fatalf("expected %d servers, got: %d: %v", len(expected), n, servers)
		}
		received := make(map[string]bool, len(servers))
		for _, server := range servers {
			received[server] = true
		}
		foundCount := 0
		for k, _ := range received {
			_, found := expected[k]
			if found {
				foundCount++
			}
		}
		if foundCount != len(expected) {
			t.Fatalf("bad servers result")
		}
	})
}
