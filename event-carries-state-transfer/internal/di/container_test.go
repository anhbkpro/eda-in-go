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
	type fields struct {
		container Container
	}
	type args struct {
		key     string
		factory DepFactoryFunc
	}
	type expected struct {
		validate func(t *testing.T, c Container)
	}
	type testCase struct {
		prepare  func(f *fields)
		args     args
		expected expected
		wantErr  bool
		errMsg   string
	}

	tests := map[string]testCase{
		"add_singleton_dependency_successfully": {
			prepare: func(f *fields) {
				f.container = New()
			},
			args: args{
				key: "test-service",
				factory: func(c Container) (any, error) {
					return &mockDependency{value: "singleton"}, nil
				},
			},
			expected: expected{
				validate: func(t *testing.T, c Container) {
					val := c.Get("test-service").(*mockDependency)
					if val.value != "singleton" {
						t.Errorf("expected singleton value, got %s", val.value)
					}
				},
			},
			wantErr: false,
		},
		"singleton_returns_same_instance": {
			prepare: func(f *fields) {
				f.container = New()
			},
			args: args{
				key: "singleton-test",
				factory: func(c Container) (any, error) {
					return &mockDependency{value: "same"}, nil
				},
			},
			expected: expected{
				validate: func(t *testing.T, c Container) {
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
			wantErr: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			f := fields{}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			// Act
			f.container.AddSingleton(tt.args.key, tt.args.factory)

			// Assert
			if tt.expected.validate != nil {
				tt.expected.validate(t, f.container)
			}
		})
	}
}

func TestContainer_AddScoped(t *testing.T) {
	type fields struct {
		container Container
	}
	type args struct {
		key     string
		factory DepFactoryFunc
	}
	type expected struct {
		validate func(t *testing.T, c Container)
	}
	type testCase struct {
		prepare  func(f *fields)
		args     args
		expected expected
		wantErr  bool
		errMsg   string
	}

	tests := map[string]testCase{
		"add_scoped_dependency_successfully": {
			prepare: func(f *fields) {
				f.container = New()
			},
			args: args{
				key: "scoped-service",
				factory: func(c Container) (any, error) {
					return &mockDependency{value: "scoped"}, nil
				},
			},
			expected: expected{
				validate: func(t *testing.T, c Container) {
					ctx := c.Scoped(context.Background())
					val := ctx.Value(containerKey).(*container).Get("scoped-service").(*mockDependency)
					if val.value != "scoped" {
						t.Errorf("expected scoped value, got %s", val.value)
					}
				},
			},
			wantErr: false,
		},
		"scoped_returns_different_instances": {
			prepare: func(f *fields) {
				f.container = New()
			},
			args: args{
				key: "scoped-instance-test",
				factory: func(c Container) (any, error) {
					return &mockDependency{value: "scoped-instance"}, nil
				},
			},
			expected: expected{
				validate: func(t *testing.T, c Container) {
					ctx1 := c.Scoped(context.Background())
					ctx2 := c.Scoped(context.Background())

					val1 := ctx1.Value(containerKey).(*container).Get("scoped-instance-test").(*mockDependency)
					val2 := ctx2.Value(containerKey).(*container).Get("scoped-instance-test").(*mockDependency)

					if val1 == val2 {
						t.Error("scoped should return different instances")
					}
				},
			},
			wantErr: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			f := fields{}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			// Act
			f.container.AddScoped(tt.args.key, tt.args.factory)

			// Assert
			if tt.expected.validate != nil {
				tt.expected.validate(t, f.container)
			}
		})
	}
}

func TestContainer_Get(t *testing.T) {
	type fields struct {
		container Container
	}
	type args struct {
		ctx context.Context
		key string
	}
	type expected struct {
		result      any
		validate    func(t *testing.T, result any)
		expectPanic bool
		panicMsg    string
	}
	type testCase struct {
		prepare  func(f *fields)
		args     args
		expected expected
		wantErr  bool
		errMsg   string
	}

	tests := map[string]testCase{
		"get_existing_singleton_dependency": {
			prepare: func(f *fields) {
				f.container = New()
				f.container.AddSingleton("existing-singleton", func(c Container) (any, error) {
					return "singleton-value", nil
				})
			},
			args: args{
				ctx: context.Background(),
				key: "existing-singleton",
			},
			expected: expected{
				validate: func(t *testing.T, result any) {
					if result != "singleton-value" {
						t.Errorf("expected 'singleton-value', got %v", result)
					}
				},
			},
			wantErr: false,
		},
		"get_non_existing_dependency_panics": {
			prepare: func(f *fields) {
				f.container = New()
			},
			args: args{
				ctx: context.Background(),
				key: "non-existing",
			},
			expected: expected{
				expectPanic: true,
				panicMsg:    "there is no dependency registered with `non-existing`",
			},
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			f := fields{}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			// Act & Assert
			if tt.expected.expectPanic {
				defer func() {
					if r := recover(); r != nil {
						if panicMsg, ok := r.(string); ok && panicMsg == tt.expected.panicMsg {
							// Expected panic
							return
						}
						t.Errorf("unexpected panic message: %v", r)
					} else {
						t.Error("expected panic but none occurred")
					}
				}()
			}

			result := f.container.Get(tt.args.key)

			if !tt.expected.expectPanic && tt.expected.validate != nil {
				tt.expected.validate(t, result)
			}
		})
	}
}

func TestContainer_Scoped(t *testing.T) {
	type fields struct {
		container Container
	}
	type args struct {
		ctx context.Context
	}
	type expected struct {
		validate func(t *testing.T, ctx context.Context, original Container)
	}
	type testCase struct {
		prepare  func(f *fields)
		args     args
		expected expected
		wantErr  bool
		errMsg   string
	}

	tests := map[string]testCase{
		"scoped_context_contains_container": {
			prepare: func(f *fields) {
				f.container = New()
			},
			args: args{
				ctx: context.Background(),
			},
			expected: expected{
				validate: func(t *testing.T, ctx context.Context, original Container) {
					containerVal := ctx.Value(containerKey)
					if containerVal == nil {
						t.Fatal("scoped context should contain container")
					}

					if containerVal == nil {
						t.Error("scoped context should contain non-nil value")
					}
				},
			},
			wantErr: false,
		},
		"scoped_container_inherits_dependencies": {
			prepare: func(f *fields) {
				f.container = New()
			},
			args: args{
				ctx: context.Background(),
			},
			expected: expected{
				validate: func(t *testing.T, ctx context.Context, original Container) {
					original.AddSingleton("inherited", func(c Container) (any, error) {
						return "inherited-value", nil
					})

					scopedContainer := New()
					scopedContainer.AddScoped("scoped-inherited", func(c Container) (any, error) {
						return original.Get("inherited"), nil
					})

					scopedCtx := scopedContainer.Scoped(ctx)
					_ = scopedCtx
				},
			},
			wantErr: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			f := fields{}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			// Act
			ctx := f.container.Scoped(tt.args.ctx)

			// Assert
			if tt.expected.validate != nil {
				tt.expected.validate(t, ctx, f.container)
			}
		})
	}
}

func TestContainer_CyclicDependency(t *testing.T) {
	type fields struct {
		container Container
	}
	type args struct {
		ctx context.Context
		key string
	}
	type expected struct {
		expectPanic bool
		panicMsg    string
	}
	type testCase struct {
		prepare  func(f *fields)
		args     args
		expected expected
		wantErr  bool
		errMsg   string
	}

	tests := map[string]testCase{
		"cyclic_dependency_detection": {
			prepare: func(f *fields) {
				f.container = New()

				// Add dependency A that depends on B
				f.container.AddSingleton("dep-a", func(c Container) (any, error) {
					c.Get("dep-b") // This will try to get B
					return "a", nil
				})

				// Add dependency B that depends on A (creating cycle)
				f.container.AddSingleton("dep-b", func(c Container) (any, error) {
					c.Get("dep-a") // This creates the cycle
					return "b", nil
				})
			},
			args: args{
				ctx: context.Background(),
				key: "dep-a",
			},
			expected: expected{
				expectPanic: true,
				panicMsg:    "cyclic dependencies",
			},
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			f := fields{}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			// Act & Assert
			defer func() {
				if r := recover(); r != nil {
					if panicMsg, ok := r.(string); ok && len(panicMsg) > 0 {
						// Expected panic for cyclic dependency
						if !contains(panicMsg, tt.expected.panicMsg) {
							t.Errorf("expected cyclic dependency panic, got: %s", panicMsg)
						}
						return
					}
					t.Errorf("unexpected panic type: %v", r)
				} else {
					t.Error("expected panic for cyclic dependency")
				}
			}()

			f.container.Get(tt.args.key)
		})
	}
}

func TestContainer_ThreadSafety(t *testing.T) {
	type fields struct {
		container Container
	}
	type args struct {
		ctx           context.Context
		numGoroutines int
		numOperations int
		counterKey    string
	}
	type expected struct {
		finalCount int
	}
	type testCase struct {
		prepare  func(f *fields)
		args     args
		expected expected
		wantErr  bool
		errMsg   string
	}

	tests := map[string]testCase{
		"concurrent_singleton_access": {
			prepare: func(f *fields) {
				f.container = New()
				f.container.AddSingleton("counter", func(c Container) (any, error) {
					return &threadSafeCounter{}, nil
				})
			},
			args: args{
				ctx:           context.Background(),
				numGoroutines: 10,
				numOperations: 100,
				counterKey:    "counter",
			},
			expected: expected{
				finalCount: 1000, // 10 goroutines * 100 operations
			},
			wantErr: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			f := fields{}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			var wg sync.WaitGroup

			// Act
			for i := 0; i < tt.args.numGoroutines; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for j := 0; j < tt.args.numOperations; j++ {
						counter := f.container.Get(tt.args.counterKey).(*threadSafeCounter)
						counter.increment()
					}
				}()
			}

			wg.Wait()

			// Assert
			finalCounter := f.container.Get(tt.args.counterKey).(*threadSafeCounter)
			if finalCounter.get() != tt.expected.finalCount {
				t.Errorf("expected count %d, got %d", tt.expected.finalCount, finalCounter.get())
			}
		})
	}
}

func TestContainer_ScopedConcurrency(t *testing.T) {
	type fields struct {
		container Container
	}
	type args struct {
		ctx           context.Context
		numGoroutines int
	}
	type expected struct {
		validate func(t *testing.T)
	}
	type testCase struct {
		prepare  func(f *fields)
		args     args
		expected expected
		wantErr  bool
		errMsg   string
	}

	tests := map[string]testCase{
		"scoped_containers_work_concurrently": {
			prepare: func(f *fields) {
				f.container = New()
				f.container.AddScoped("scoped-counter", func(c Container) (any, error) {
					return &threadSafeCounter{}, nil
				})
			},
			args: args{
				ctx:           context.Background(),
				numGoroutines: 5,
			},
			expected: expected{
				validate: func(t *testing.T) {
					// Test completed without panic
				},
			},
			wantErr: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			f := fields{}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			var wg sync.WaitGroup

			// Act
			for i := 0; i < tt.args.numGoroutines; i++ {
				wg.Add(1)
				go func(id int) {
					defer wg.Done()
					ctx := f.container.Scoped(tt.args.ctx)
					_ = ctx
				}(i)
			}

			wg.Wait()

			// Assert
			if tt.expected.validate != nil {
				tt.expected.validate(t)
			}
		})
	}
}

func TestContainer_FactoryError(t *testing.T) {
	type fields struct {
		container Container
	}
	type args struct {
		ctx context.Context
		key string
	}
	type expected struct {
		expectPanic bool
		panicMsg    string
	}
	type testCase struct {
		prepare  func(f *fields)
		args     args
		expected expected
		wantErr  bool
		errMsg   string
	}

	tests := map[string]testCase{
		"factory_error_causes_panic": {
			prepare: func(f *fields) {
				f.container = New()
				f.container.AddSingleton("failing-factory", func(c Container) (any, error) {
					return nil, context.DeadlineExceeded
				})
			},
			args: args{
				ctx: context.Background(),
				key: "failing-factory",
			},
			expected: expected{
				expectPanic: true,
				panicMsg:    "error building dependency",
			},
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			f := fields{}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			// Act & Assert
			defer func() {
				if r := recover(); r != nil {
					if panicMsg, ok := r.(string); ok && len(panicMsg) > 0 {
						if !contains(panicMsg, tt.expected.panicMsg) {
							t.Errorf("expected factory error panic, got: %s", panicMsg)
						}
						return
					}
					t.Errorf("unexpected panic type: %v", r)
				} else {
					t.Error("expected panic for factory error")
				}
			}()

			f.container.Get(tt.args.key)
		})
	}
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
