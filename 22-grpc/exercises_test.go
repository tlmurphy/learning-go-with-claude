package grpcmod

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"testing"
	"time"
)

// =============================================================================
// Exercise 1: UserServiceImpl Tests
// =============================================================================

func TestUserServiceImpl(t *testing.T) {
	t.Run("GetUser existing", func(t *testing.T) {
		store := NewUserStore()
		svc := NewUserServiceImpl(store)
		ctx := context.Background()

		user, err := svc.GetUser(ctx, &GetUserRequest{ID: "user-1"})
		if err != nil {
			t.Fatalf("GetUser returned unexpected error: %v", err)
		}
		if user == nil {
			t.Fatal("GetUser returned nil user for existing user-1")
		}
		if user.Name != "Alice" {
			t.Errorf("Expected name 'Alice', got %q. Look up the user in store.users by ID.", user.Name)
		}
	})

	t.Run("GetUser not found", func(t *testing.T) {
		store := NewUserStore()
		svc := NewUserServiceImpl(store)
		ctx := context.Background()

		_, err := svc.GetUser(ctx, &GetUserRequest{ID: "nonexistent"})
		if err == nil {
			t.Fatal("GetUser should return an error for nonexistent user. Use StatusError(NotFound, ...).")
		}
		st := FromError(err)
		if st.Code() != NotFound {
			t.Errorf("Expected NotFound status code, got %s. Return StatusError(NotFound, ...) when user doesn't exist.", st.Code())
		}
	})

	t.Run("CreateUser success", func(t *testing.T) {
		store := NewUserStore()
		svc := NewUserServiceImpl(store)
		ctx := context.Background()

		user, err := svc.CreateUser(ctx, &CreateUserRequest{Name: "Diana", Email: "diana@example.com"})
		if err != nil {
			t.Fatalf("CreateUser returned unexpected error: %v", err)
		}
		if user == nil {
			t.Fatal("CreateUser returned nil user")
		}
		if user.ID == "" {
			t.Error("CreateUser should generate an ID for the new user (e.g., 'user-4').")
		}
		if user.Name != "Diana" {
			t.Errorf("Expected name 'Diana', got %q", user.Name)
		}
		if user.Email != "diana@example.com" {
			t.Errorf("Expected email 'diana@example.com', got %q", user.Email)
		}

		// Verify it's actually stored
		stored, ok := store.users[user.ID]
		if !ok {
			t.Errorf("Created user not found in store. Make sure to add the user to store.users.")
		} else if stored.Name != "Diana" {
			t.Errorf("Stored user has wrong name: %q", stored.Name)
		}
	})

	t.Run("CreateUser invalid input", func(t *testing.T) {
		store := NewUserStore()
		svc := NewUserServiceImpl(store)
		ctx := context.Background()

		_, err := svc.CreateUser(ctx, &CreateUserRequest{Name: "", Email: "test@example.com"})
		if err == nil {
			t.Fatal("CreateUser should return InvalidArgument when Name is empty.")
		}
		st := FromError(err)
		if st.Code() != InvalidArgument {
			t.Errorf("Expected InvalidArgument for empty name, got %s", st.Code())
		}

		_, err = svc.CreateUser(ctx, &CreateUserRequest{Name: "Test", Email: ""})
		if err == nil {
			t.Fatal("CreateUser should return InvalidArgument when Email is empty.")
		}
		st = FromError(err)
		if st.Code() != InvalidArgument {
			t.Errorf("Expected InvalidArgument for empty email, got %s", st.Code())
		}
	})

	t.Run("UpdateUser success", func(t *testing.T) {
		store := NewUserStore()
		svc := NewUserServiceImpl(store)
		ctx := context.Background()

		user, err := svc.UpdateUser(ctx, &UpdateUserRequest{
			ID:    "user-1",
			Name:  "Alice Updated",
			Email: "alice.new@example.com",
		})
		if err != nil {
			t.Fatalf("UpdateUser returned unexpected error: %v", err)
		}
		if user == nil {
			t.Fatal("UpdateUser returned nil user")
		}
		if user.Name != "Alice Updated" {
			t.Errorf("Expected updated name 'Alice Updated', got %q", user.Name)
		}
		if user.Email != "alice.new@example.com" {
			t.Errorf("Expected updated email, got %q", user.Email)
		}
	})

	t.Run("UpdateUser not found", func(t *testing.T) {
		store := NewUserStore()
		svc := NewUserServiceImpl(store)
		ctx := context.Background()

		_, err := svc.UpdateUser(ctx, &UpdateUserRequest{ID: "nonexistent"})
		if err == nil {
			t.Fatal("UpdateUser should return NotFound for nonexistent user.")
		}
		st := FromError(err)
		if st.Code() != NotFound {
			t.Errorf("Expected NotFound, got %s", st.Code())
		}
	})

	t.Run("DeleteUser success", func(t *testing.T) {
		store := NewUserStore()
		svc := NewUserServiceImpl(store)
		ctx := context.Background()

		result, err := svc.DeleteUser(ctx, &DeleteUserRequest{ID: "user-1"})
		if err != nil {
			t.Fatalf("DeleteUser returned unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("DeleteUser should return &Empty{}, not nil")
		}
		if _, exists := store.users["user-1"]; exists {
			t.Error("DeleteUser should remove the user from the store.")
		}
	})

	t.Run("DeleteUser not found", func(t *testing.T) {
		store := NewUserStore()
		svc := NewUserServiceImpl(store)
		ctx := context.Background()

		_, err := svc.DeleteUser(ctx, &DeleteUserRequest{ID: "nonexistent"})
		if err == nil {
			t.Fatal("DeleteUser should return NotFound for nonexistent user.")
		}
		st := FromError(err)
		if st.Code() != NotFound {
			t.Errorf("Expected NotFound, got %s", st.Code())
		}
	})
}

// =============================================================================
// Exercise 2: GetUserHandler Tests
// =============================================================================

func TestGetUserHandler(t *testing.T) {
	t.Run("successful lookup", func(t *testing.T) {
		store := NewUserStore()
		ctx := context.Background()

		user, reqID, err := GetUserHandler(ctx, store, "user-1")
		if err != nil {
			t.Fatalf("GetUserHandler returned unexpected error: %v", err)
		}
		if user == nil {
			t.Fatal("Expected a user, got nil")
		}
		if user.Name != "Alice" {
			t.Errorf("Expected Alice, got %q", user.Name)
		}
		if reqID != "" {
			t.Errorf("Expected empty request ID when no metadata, got %q", reqID)
		}
	})

	t.Run("with metadata", func(t *testing.T) {
		store := NewUserStore()
		md := NewMetadata("x-request-id", "req-123")
		ctx := NewIncomingContext(context.Background(), md)

		user, reqID, err := GetUserHandler(ctx, store, "user-2")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if user == nil || user.Name != "Bob" {
			t.Error("Expected Bob")
		}
		if reqID != "req-123" {
			t.Errorf("Expected request ID 'req-123', got %q. "+
				"Extract it from incoming metadata using FromIncomingContext.", reqID)
		}
	})

	t.Run("empty user ID", func(t *testing.T) {
		store := NewUserStore()
		ctx := context.Background()

		_, _, err := GetUserHandler(ctx, store, "")
		if err == nil {
			t.Fatal("Should return InvalidArgument for empty user ID.")
		}
		st := FromError(err)
		if st.Code() != InvalidArgument {
			t.Errorf("Expected InvalidArgument, got %s", st.Code())
		}
	})

	t.Run("user not found", func(t *testing.T) {
		store := NewUserStore()
		ctx := context.Background()

		_, _, err := GetUserHandler(ctx, store, "nonexistent")
		if err == nil {
			t.Fatal("Should return NotFound for nonexistent user.")
		}
		st := FromError(err)
		if st.Code() != NotFound {
			t.Errorf("Expected NotFound, got %s", st.Code())
		}
	})

	t.Run("cancelled context", func(t *testing.T) {
		store := NewUserStore()
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		_, _, err := GetUserHandler(ctx, store, "user-1")
		if err == nil {
			t.Fatal("Should return DeadlineExceeded for cancelled context. "+
				"Check ctx.Err() before doing work.")
		}
		st := FromError(err)
		if st.Code() != DeadlineExceeded {
			t.Errorf("Expected DeadlineExceeded, got %s", st.Code())
		}
	})
}

// =============================================================================
// Exercise 3: Server Streaming Tests
// =============================================================================

// mockUserStream is a test double for UserStream.
type mockUserStream struct {
	users   []*User
	ctx     context.Context
	sendErr error // If set, Send returns this error
}

func newMockStream(ctx context.Context) *mockUserStream {
	return &mockUserStream{ctx: ctx}
}

func (m *mockUserStream) Send(user *User) error {
	if m.sendErr != nil {
		return m.sendErr
	}
	m.users = append(m.users, user)
	return nil
}

func (m *mockUserStream) Context() context.Context {
	return m.ctx
}

func TestStreamUsers(t *testing.T) {
	t.Run("stream all users", func(t *testing.T) {
		store := NewUserStore()
		stream := newMockStream(context.Background())

		err := StreamUsers(store, &ListUsersRequest{}, stream)
		if err != nil {
			t.Fatalf("StreamUsers returned error: %v", err)
		}
		if len(stream.users) != 3 {
			t.Errorf("Expected 3 users streamed, got %d. "+
				"Iterate over store.users and call stream.Send for each.", len(stream.users))
		}
	})

	t.Run("stream with page size", func(t *testing.T) {
		store := NewUserStore()
		stream := newMockStream(context.Background())

		err := StreamUsers(store, &ListUsersRequest{PageSize: 2}, stream)
		if err != nil {
			t.Fatalf("StreamUsers returned error: %v", err)
		}
		if len(stream.users) != 2 {
			t.Errorf("Expected 2 users (PageSize=2), got %d. "+
				"Honor the PageSize field in ListUsersRequest.", len(stream.users))
		}
	})

	t.Run("stream with cancelled context", func(t *testing.T) {
		store := NewUserStore()
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately
		stream := newMockStream(ctx)

		// With a cancelled context, the function should stop streaming.
		// It may return an error or may have sent 0 users — both are fine.
		_ = StreamUsers(store, &ListUsersRequest{}, stream)
		// The key check: it shouldn't send ALL users when context is cancelled.
		// (It might send 0 or some, depending on when cancellation is checked.)
		if len(stream.users) == 3 {
			t.Error("StreamUsers sent all users despite cancelled context. " +
				"Check stream.Context() for cancellation between sends.")
		}
	})

	t.Run("stream handles send error", func(t *testing.T) {
		store := NewUserStore()
		stream := newMockStream(context.Background())
		stream.sendErr = fmt.Errorf("connection closed")

		err := StreamUsers(store, &ListUsersRequest{}, stream)
		if err == nil {
			t.Error("StreamUsers should return the error from Send when it fails.")
		}
	})
}

// =============================================================================
// Exercise 4: Domain Error to Status Conversion Tests
// =============================================================================

func TestDomainToStatus(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode Code
		expectedMsg  string
	}{
		{
			name:         "nil error",
			err:          nil,
			expectedCode: OK,
			expectedMsg:  "",
		},
		{
			name:         "not found",
			err:          &DomainError{Kind: ErrNotFound, Message: "user not found"},
			expectedCode: NotFound,
			expectedMsg:  "user not found",
		},
		{
			name:         "already exists",
			err:          &DomainError{Kind: ErrAlreadyExists, Message: "email taken"},
			expectedCode: AlreadyExists,
			expectedMsg:  "email taken",
		},
		{
			name:         "invalid input",
			err:          &DomainError{Kind: ErrInvalidInput, Message: "name required"},
			expectedCode: InvalidArgument,
			expectedMsg:  "name required",
		},
		{
			name:         "unauthorized",
			err:          &DomainError{Kind: ErrUnauthorized, Message: "bad token"},
			expectedCode: Unauthenticated,
			expectedMsg:  "bad token",
		},
		{
			name:         "forbidden",
			err:          &DomainError{Kind: ErrForbidden, Message: "admin only"},
			expectedCode: PermissionDenied,
			expectedMsg:  "admin only",
		},
		{
			name:         "internal",
			err:          &DomainError{Kind: ErrInternal, Message: "db crashed"},
			expectedCode: Internal,
			expectedMsg:  "db crashed",
		},
		{
			name:         "unknown error type",
			err:          fmt.Errorf("something weird happened"),
			expectedCode: Unknown,
			expectedMsg:  "something weird happened",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DomainToStatus(tt.err)

			if tt.err == nil {
				if result != nil {
					t.Error("DomainToStatus(nil) should return nil. " +
						"Check for nil input before converting.")
				}
				return
			}

			if result == nil {
				t.Fatalf("DomainToStatus returned nil for non-nil error %v. "+
					"Should return a Status error.", tt.err)
			}

			st := FromError(result)
			if st.Code() != tt.expectedCode {
				t.Errorf("Expected code %s, got %s. "+
					"Map DomainErrorKind to the appropriate gRPC status code.",
					tt.expectedCode, st.Code())
			}
			if st.Message() != tt.expectedMsg {
				t.Errorf("Expected message %q, got %q. "+
					"Preserve the error message from the domain error.",
					tt.expectedMsg, st.Message())
			}
		})
	}
}

// =============================================================================
// Exercise 5: Logging Interceptor Tests
// =============================================================================

func TestLoggingInterceptor(t *testing.T) {
	t.Run("logs successful call", func(t *testing.T) {
		log := NewCallLog()
		interceptor := NewLoggingInterceptor(log)

		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			time.Sleep(5 * time.Millisecond) // Simulate some work
			return "result", nil
		}
		info := &UnaryServerInfo{FullMethod: "/test.Service/Method"}

		resp, err := interceptor(context.Background(), "request", info, handler)
		if err != nil {
			t.Fatalf("Interceptor returned unexpected error: %v", err)
		}
		if resp != "result" {
			t.Errorf("Expected response 'result', got %v. "+
				"The interceptor should pass through the handler's return values.", resp)
		}

		records := log.Records()
		if len(records) != 1 {
			t.Fatalf("Expected 1 call record, got %d. "+
				"Append a CallRecord to the log after calling the handler.", len(records))
		}

		record := records[0]
		if record.Method != "/test.Service/Method" {
			t.Errorf("Expected method '/test.Service/Method', got %q", record.Method)
		}
		if record.Duration < 5*time.Millisecond {
			t.Errorf("Duration should be >= 5ms, got %v. "+
				"Record time.Since(start) after calling the handler.", record.Duration)
		}
		if record.Code != OK {
			t.Errorf("Expected code OK for successful call, got %s", record.Code)
		}
	})

	t.Run("logs failed call", func(t *testing.T) {
		log := NewCallLog()
		interceptor := NewLoggingInterceptor(log)

		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, StatusError(NotFound, "not found")
		}
		info := &UnaryServerInfo{FullMethod: "/test.Service/GetUser"}

		resp, err := interceptor(context.Background(), "request", info, handler)
		if err == nil {
			t.Fatal("Interceptor should pass through the handler's error.")
		}
		if resp != nil {
			t.Errorf("Expected nil response on error, got %v", resp)
		}

		records := log.Records()
		if len(records) != 1 {
			t.Fatalf("Expected 1 call record, got %d", len(records))
		}
		if records[0].Code != NotFound {
			t.Errorf("Expected NotFound code, got %s. "+
				"Use FromError to extract the status code from the handler's error.",
				records[0].Code)
		}
		if records[0].Error != "not found" {
			t.Errorf("Expected error message 'not found', got %q", records[0].Error)
		}
	})

	t.Run("multiple calls are logged", func(t *testing.T) {
		log := NewCallLog()
		interceptor := NewLoggingInterceptor(log)
		info := &UnaryServerInfo{FullMethod: "/test.Service/Method"}

		okHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return "ok", nil
		}
		errHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, StatusError(Internal, "oops")
		}

		_, _ = interceptor(context.Background(), nil, info, okHandler)
		_, _ = interceptor(context.Background(), nil, info, errHandler)
		_, _ = interceptor(context.Background(), nil, info, okHandler)

		records := log.Records()
		if len(records) != 3 {
			t.Errorf("Expected 3 records, got %d", len(records))
		}
	})
}

// =============================================================================
// Exercise 6: Deadline Handler Tests
// =============================================================================

func TestDoWorkWithDeadline(t *testing.T) {
	t.Run("completes within deadline", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		result, err := DoWorkWithDeadline(ctx, 10*time.Millisecond)
		if err != nil {
			t.Fatalf("Expected success, got error: %v. "+
				"The work (10ms) should complete within the deadline (1s).", err)
		}
		if result != "completed" {
			t.Errorf("Expected 'completed', got %q", result)
		}
	})

	t.Run("deadline exceeded", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		result, err := DoWorkWithDeadline(ctx, 1*time.Second)
		if err == nil {
			t.Fatal("Expected DeadlineExceeded error. " +
				"Use select with ctx.Done() and time.After(workDuration).")
		}
		st := FromError(err)
		if st.Code() != DeadlineExceeded {
			t.Errorf("Expected DeadlineExceeded, got %s", st.Code())
		}
		if result != "" {
			t.Errorf("Expected empty result on deadline, got %q", result)
		}
	})

	t.Run("already cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Already cancelled

		_, err := DoWorkWithDeadline(ctx, 10*time.Millisecond)
		if err == nil {
			t.Fatal("Should return error for already-cancelled context. " +
				"Check ctx.Err() before starting work.")
		}
		st := FromError(err)
		if st.Code() != DeadlineExceeded {
			t.Errorf("Expected DeadlineExceeded, got %s", st.Code())
		}
	})
}

// =============================================================================
// Exercise 7: Health Check Service Tests
// =============================================================================

func TestServiceHealth(t *testing.T) {
	t.Run("all services initially serving", func(t *testing.T) {
		sh := NewServiceHealth("db", "cache", "queue")

		for _, name := range []string{"db", "cache", "queue"} {
			status := sh.CheckHealth(name)
			if status != StatusServing {
				t.Errorf("Service %q should initially be StatusServing, got %s. "+
					"Initialize all services as StatusServing in NewServiceHealth.",
					name, status)
			}
		}
	})

	t.Run("overall health when all serving", func(t *testing.T) {
		sh := NewServiceHealth("db", "cache")

		overall := sh.OverallHealth()
		if overall != StatusServing {
			t.Errorf("OverallHealth should be StatusServing when all services are serving, got %s",
				overall)
		}
	})

	t.Run("overall health with one not serving", func(t *testing.T) {
		sh := NewServiceHealth("db", "cache")
		_ = sh.SetStatus("cache", StatusNotServing)

		overall := sh.OverallHealth()
		if overall != StatusNotServing {
			t.Errorf("OverallHealth should be StatusNotServing when any service is not serving, got %s",
				overall)
		}
	})

	t.Run("unknown service check", func(t *testing.T) {
		sh := NewServiceHealth("db")

		status := sh.CheckHealth("nonexistent")
		if status != StatusServiceUnknown {
			t.Errorf("CheckHealth for unknown service should return StatusServiceUnknown, got %s",
				status)
		}
	})

	t.Run("set status for unknown service returns error", func(t *testing.T) {
		sh := NewServiceHealth("db")

		err := sh.SetStatus("nonexistent", StatusServing)
		if err == nil {
			t.Error("SetStatus should return an error for unknown service names.")
		}
	})

	t.Run("update and check status", func(t *testing.T) {
		sh := NewServiceHealth("db", "cache")

		_ = sh.SetStatus("db", StatusNotServing)
		if sh.CheckHealth("db") != StatusNotServing {
			t.Error("CheckHealth should reflect the updated status after SetStatus.")
		}

		_ = sh.SetStatus("db", StatusServing)
		if sh.CheckHealth("db") != StatusServing {
			t.Error("CheckHealth should reflect the restored status.")
		}
	})

	t.Run("concurrent access is safe", func(t *testing.T) {
		sh := NewServiceHealth("db", "cache", "queue")
		var wg sync.WaitGroup

		// Hammer it from multiple goroutines
		for i := 0; i < 100; i++ {
			wg.Add(3)
			go func() {
				defer wg.Done()
				_ = sh.SetStatus("db", StatusNotServing)
			}()
			go func() {
				defer wg.Done()
				sh.CheckHealth("cache")
			}()
			go func() {
				defer wg.Done()
				sh.OverallHealth()
			}()
		}
		wg.Wait()
		// If we get here without a race detector panic, we're good.
	})
}

// =============================================================================
// Exercise 8: Retry Wrapper Tests
// =============================================================================

func TestRetryCall(t *testing.T) {
	t.Run("succeeds on first try", func(t *testing.T) {
		calls := 0
		result, err := RetryCall(
			context.Background(),
			DefaultRetryConfig(),
			func(ctx context.Context) (interface{}, error) {
				calls++
				return "success", nil
			},
		)
		if err != nil {
			t.Fatalf("Expected success, got error: %v", err)
		}
		if result != "success" {
			t.Errorf("Expected 'success', got %v", result)
		}
		if calls != 1 {
			t.Errorf("Should only call once on success, called %d times", calls)
		}
	})

	t.Run("retries on retryable error then succeeds", func(t *testing.T) {
		calls := 0
		config := RetryConfig{
			MaxRetries:  3,
			InitialWait: 1 * time.Millisecond,
			MaxWait:     10 * time.Millisecond,
			Multiplier:  2.0,
		}

		result, err := RetryCall(
			context.Background(),
			config,
			func(ctx context.Context) (interface{}, error) {
				calls++
				if calls < 3 {
					return nil, StatusError(Unavailable, "server busy")
				}
				return "finally", nil
			},
		)
		if err != nil {
			t.Fatalf("Expected success after retries, got error: %v", err)
		}
		if result != "finally" {
			t.Errorf("Expected 'finally', got %v", result)
		}
		if calls != 3 {
			t.Errorf("Expected 3 calls (2 failures + 1 success), got %d. "+
				"Retry retryable errors up to MaxRetries times.", calls)
		}
	})

	t.Run("gives up after max retries", func(t *testing.T) {
		calls := 0
		config := RetryConfig{
			MaxRetries:  2,
			InitialWait: 1 * time.Millisecond,
			MaxWait:     10 * time.Millisecond,
			Multiplier:  2.0,
		}

		_, err := RetryCall(
			context.Background(),
			config,
			func(ctx context.Context) (interface{}, error) {
				calls++
				return nil, StatusError(Unavailable, "still busy")
			},
		)
		if err == nil {
			t.Fatal("Should return error after exhausting retries.")
		}
		// 1 initial call + 2 retries = 3 total calls
		if calls != 3 {
			t.Errorf("Expected 3 total calls (1 initial + 2 retries), got %d", calls)
		}
	})

	t.Run("no retry on non-retryable error", func(t *testing.T) {
		calls := 0
		result, err := RetryCall(
			context.Background(),
			DefaultRetryConfig(),
			func(ctx context.Context) (interface{}, error) {
				calls++
				return nil, StatusError(NotFound, "user not found")
			},
		)
		if err == nil {
			t.Fatal("Should return the non-retryable error immediately.")
		}
		if result != nil {
			t.Errorf("Expected nil result on error, got %v", result)
		}
		if calls != 1 {
			t.Errorf("Should not retry NotFound errors. Expected 1 call, got %d. "+
				"Use IsRetryable() to check if the error is worth retrying.", calls)
		}
	})

	t.Run("respects context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		config := RetryConfig{
			MaxRetries:  100,
			InitialWait: 200 * time.Millisecond, // Longer than the context timeout
			MaxWait:     1 * time.Second,
			Multiplier:  2.0,
		}

		_, err := RetryCall(
			ctx,
			config,
			func(ctx context.Context) (interface{}, error) {
				return nil, StatusError(Unavailable, "busy")
			},
		)
		if err == nil {
			t.Fatal("Should return error when context is cancelled during backoff wait.")
		}
	})
}

// =============================================================================
// Helper: sort user names for deterministic comparison
// =============================================================================

func sortedUserNames(users []*User) []string {
	names := make([]string, len(users))
	for i, u := range users {
		names[i] = u.Name
	}
	sort.Strings(names)
	return names
}
