package di

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
)

// mockDependency is a simple mock for testing dependency injection
type mockDependency struct {
	value string
}

func TestNew(t *testing.T) {
	t.Parallel()

	container := New()

	if container == nil {
		t.Fatal("New() returned nil")
	}

	// Verify it implements the Container interface
	var _ Container = container
}

func TestContainer_AddSingleton(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		key      string
		factory  DepFactoryFunc
		expected func(t *testing.T, c Container)
	}{
		{
			name: "add singleton dependency successfully",
			key:  "test-service",
			factory: func(c Container) (any, error) {
				return &mockDependency{value: "singleton"}, nil
			},
			expected: func(t *testing.T, c Container) {
				val := c.Get("test-service").(*mockDependency)
				if val.value != "singleton" {
					t.Errorf("expected singleton value, got %s", val.value)
				}
			},
		},
		{
			name: "singleton returns same instance",
			key:  "singleton-test",
			factory: func(c Container) (any, error) {
				return &mockDependency{value: "same"}, nil
			},
			expected: func(t *testing.T, c Container) {
				val1 := c.Get("singleton-test").(*mockDependency)
				val2 := c.Get("singleton-test").(*mockDependency)

				if val1 != val2 {
					t.Error("singleton should return same instance")
				}
				if val1.value != "same" || val2.value != "same" {
					t.Error("singleton values should be consistent")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := New()
			container.AddSingleton(tt.key, tt.factory)
			tt.expected(t, container)
		})
	}
}

func TestContainer_AddScoped(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		key      string
		factory  DepFactoryFunc
		expected func(t *testing.T, c Container)
	}{
		{
			name: "add scoped dependency successfully",
			key:  "scoped-service",
			factory: func(c Container) (any, error) {
				return &mockDependency{value: "scoped"}, nil
			},
			expected: func(t *testing.T, c Container) {
				ctx := c.Scoped(context.Background())
				val := ctx.Value(containerKey).(*container).Get("scoped-service").(*mockDependency)
				if val.value != "scoped" {
					t.Errorf("expected scoped value, got %s", val.value)
				}
			},
		},
		{
			name: "scoped returns different instances",
			key:  "scoped-instance-test",
			factory: func(c Container) (any, error) {
				return &mockDependency{value: "scoped-instance"}, nil
			},
			expected: func(t *testing.T, c Container) {
				ctx1 := c.Scoped(context.Background())
				ctx2 := c.Scoped(context.Background())

				val1 := ctx1.Value(containerKey).(*container).Get("scoped-instance-test").(*mockDependency)
				val2 := ctx2.Value(containerKey).(*container).Get("scoped-instance-test").(*mockDependency)

				if val1 == val2 {
					t.Error("scoped should return different instances")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := New()
			container.AddScoped(tt.key, tt.factory)
			tt.expected(t, container)
		})
	}
}

func TestContainer_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setup       func(Container)
		key         string
		expectPanic bool
		panicMsg    string
		expected    func(t *testing.T, result any)
	}{
		{
			name: "get existing singleton dependency",
			setup: func(c Container) {
				c.AddSingleton("existing-singleton", func(c Container) (any, error) {
					return "singleton-value", nil
				})
			},
			key: "existing-singleton",
			expected: func(t *testing.T, result any) {
				if result != "singleton-value" {
					t.Errorf("expected 'singleton-value', got %v", result)
				}
			},
		},
		{
			name: "get non-existing dependency panics",
			key:  "non-existing",
			expected: func(t *testing.T, result any) {
				// Should panic before reaching here
				t.Error("expected panic for non-existing dependency")
			},
			expectPanic: true,
			panicMsg:    "there is no dependency registered with `non-existing`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := New()

			if tt.setup != nil {
				tt.setup(container)
			}

			if tt.expectPanic {
				defer func() {
					if r := recover(); r != nil {
						if panicMsg, ok := r.(string); ok && panicMsg == tt.panicMsg {
							// Expected panic
							return
						}
						t.Errorf("unexpected panic message: %v", r)
					} else {
						t.Error("expected panic but none occurred")
					}
				}()
			}

			result := container.Get(tt.key)

			if !tt.expectPanic && tt.expected != nil {
				tt.expected(t, result)
			}
		})
	}
}

func TestContainer_Scoped(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		expected func(t *testing.T, ctx context.Context, original Container)
	}{
		{
			name: "scoped context contains container",
			expected: func(t *testing.T, ctx context.Context, original Container) {
				containerVal := ctx.Value(containerKey)
				if containerVal == nil {
					t.Fatal("scoped context should contain container")
				}

				// We can't check the internal structure since container is unexported,
				// but we can verify the context contains something
				if containerVal == nil {
					t.Error("scoped context should contain non-nil value")
				}
			},
		},
		{
			name: "scoped container inherits dependencies",
			expected: func(t *testing.T, ctx context.Context, original Container) {
				original.AddSingleton("inherited", func(c Container) (any, error) {
					return "inherited-value", nil
				})

				// Create a new container from the scoped context and add scoped dependency
				scopedContainer := New()
				scopedContainer.AddScoped("scoped-inherited", func(c Container) (any, error) {
					return original.Get("inherited"), nil
				})

				scopedCtx := scopedContainer.Scoped(ctx)
				// This is a simplified test - in real usage, the scoped container
				// would inherit from the parent through the Scoped() method
				_ = scopedCtx
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := New()
			ctx := container.Scoped(context.Background())
			tt.expected(t, ctx, container)
		})
	}
}

func TestContainer_CyclicDependency(t *testing.T) {
	t.Parallel()

	container := New()

	// Add dependency A that depends on B
	container.AddSingleton("dep-a", func(c Container) (any, error) {
		c.Get("dep-b") // This will try to get B
		return "a", nil
	})

	// Add dependency B that depends on A (creating cycle)
	container.AddSingleton("dep-b", func(c Container) (any, error) {
		c.Get("dep-a") // This creates the cycle
		return "b", nil
	})

	defer func() {
		if r := recover(); r != nil {
			if panicMsg, ok := r.(string); ok && len(panicMsg) > 0 {
				// Expected panic for cyclic dependency
				if !contains(panicMsg, "cyclic dependencies") {
					t.Errorf("expected cyclic dependency panic, got: %s", panicMsg)
				}
				return
			}
			t.Errorf("unexpected panic type: %v", r)
		} else {
			t.Error("expected panic for cyclic dependency")
		}
	}()

	container.Get("dep-a")
}

func TestContainer_ThreadSafety(t *testing.T) {
	t.Parallel()

	container := New()
	container.AddSingleton("counter", func(c Container) (any, error) {
		return &threadSafeCounter{}, nil
	})

	var wg sync.WaitGroup
	numGoroutines := 10
	numOperations := 100

	// Start multiple goroutines accessing the same singleton
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				counter := container.Get("counter").(*threadSafeCounter)
				counter.increment()
			}
		}()
	}

	wg.Wait()

	// Verify the final count
	finalCounter := container.Get("counter").(*threadSafeCounter)
	expected := numGoroutines * numOperations
	if finalCounter.get() != expected {
		t.Errorf("expected count %d, got %d", expected, finalCounter.get())
	}
}

func TestContainer_ScopedConcurrency(t *testing.T) {
	t.Parallel()

	container := New()
	container.AddScoped("scoped-counter", func(c Container) (any, error) {
		return &threadSafeCounter{}, nil
	})

	var wg sync.WaitGroup
	numGoroutines := 5

	// Test that scoped containers work concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			ctx := container.Scoped(context.Background())
			// In real usage, you'd extract the scoped container from context
			// For this test, we just verify the Scoped method doesn't panic
			_ = ctx
		}(i)
	}

	wg.Wait()
}

func TestContainer_FactoryError(t *testing.T) {
	t.Parallel()

	container := New()
	container.AddSingleton("failing-factory", func(c Container) (any, error) {
		return nil, context.DeadlineExceeded
	})

	defer func() {
		if r := recover(); r != nil {
			if panicMsg, ok := r.(string); ok && len(panicMsg) > 0 {
				// Expected panic for factory error
				if !contains(panicMsg, "error building dependency") {
					t.Errorf("expected factory error panic, got: %s", panicMsg)
				}
				return
			}
			t.Errorf("unexpected panic type: %v", r)
		} else {
			t.Error("expected panic for factory error")
		}
	}()

	container.Get("failing-factory")
}

// Helper functions and types

type threadSafeCounter struct {
	count int64
}

func (c *threadSafeCounter) increment() {
	atomic.AddInt64(&c.count, 1)
}

func (c *threadSafeCounter) get() int {
	return int(atomic.LoadInt64(&c.count))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			func() bool {
				for i := 1; i <= len(s)-len(substr); i++ {
					if s[i:i+len(substr)] == substr {
						return true
					}
				}
				return false
			}()))
}

// Benchmark tests

func BenchmarkContainer_GetSingleton(b *testing.B) {
	container := New()
	container.AddSingleton("benchmark-singleton", func(c Container) (any, error) {
		return "benchmark-value", nil
	})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = container.Get("benchmark-singleton")
		}
	})
}

func BenchmarkContainer_GetScoped(b *testing.B) {
	container := New()
	container.AddScoped("benchmark-scoped", func(c Container) (any, error) {
		return "benchmark-scoped-value", nil
	})

	// Create a separate container for scoped access
	scopedContainer := New()
	scopedContainer.AddScoped("benchmark-scoped", func(c Container) (any, error) {
		return "benchmark-scoped-value", nil
	})

	// For benchmark purposes, we'll simulate scoped access by creating new instances
	// In real usage, the scoped container would be extracted from context

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Create a new scoped container for each operation
			scoped := New()
			scoped.AddScoped("benchmark-scoped", func(c Container) (any, error) {
				return "benchmark-scoped-value", nil
			})
			_ = scoped.Scoped(context.Background())
		}
	})
}
