package errors

import (
	"errors"
	"testing"
)

func TestResourceError_Error(t *testing.T) {
	tests := []struct {
		name     string
		resource string
		message  string
		err      error
		want     string
	}{
		{
			name:     "error with internal error",
			resource: "deployment",
			message:  "not found",
			err:      ErrResourceNotFound,
			want:     "deployment: not found (resource not found)",
		},
		{
			name:     "error without internal error",
			resource: "pod",
			message:  "invalid configuration",
			err:      nil,
			want:     "pod: invalid configuration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ResourceError{
				Resource: tt.resource,
				Message:  tt.message,
				Err:      tt.err,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("ResourceError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceError_Unwrap(t *testing.T) {
	err := ErrResourceNotFound
	re := &ResourceError{
		Resource: "deployment",
		Message:  "not found",
		Err:      err,
	}

	if got := re.Unwrap(); got != err {
		t.Errorf("ResourceError.Unwrap() = %v, want %v", got, err)
	}
}

func TestNewResourceNotFoundError(t *testing.T) {
	got := NewResourceNotFoundError("deployment", "not found")
	if !errors.Is(got, ErrResourceNotFound) {
		t.Error("NewResourceNotFoundError() did not return an ErrResourceNotFound")
	}
}

func TestNewInvalidMetricsError(t *testing.T) {
	got := NewInvalidMetricsError("cpu", "invalid value")
	if !errors.Is(got, ErrInvalidMetrics) {
		t.Error("NewInvalidMetricsError() did not return an ErrInvalidMetrics")
	}
}

func TestNewInvalidConfigurationError(t *testing.T) {
	got := NewInvalidConfigurationError("mimir", "invalid url")
	if !errors.Is(got, ErrInvalidConfiguration) {
		t.Error("NewInvalidConfigurationError() did not return an ErrInvalidConfiguration")
	}
}

func TestErrorChecks(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		checkFn  func(error) bool
		wantBool bool
	}{
		{
			name:     "IsResourceNotFound with correct error",
			err:      NewResourceNotFoundError("deployment", "not found"),
			checkFn:  IsResourceNotFound,
			wantBool: true,
		},
		{
			name:     "IsResourceNotFound with different error",
			err:      NewInvalidMetricsError("cpu", "invalid value"),
			checkFn:  IsResourceNotFound,
			wantBool: false,
		},
		{
			name:     "IsInvalidMetrics with correct error",
			err:      NewInvalidMetricsError("memory", "invalid value"),
			checkFn:  IsInvalidMetrics,
			wantBool: true,
		},
		{
			name:     "IsInvalidMetrics with different error",
			err:      NewResourceNotFoundError("pod", "not found"),
			checkFn:  IsInvalidMetrics,
			wantBool: false,
		},
		{
			name:     "IsInvalidConfiguration with correct error",
			err:      NewInvalidConfigurationError("mimir", "invalid url"),
			checkFn:  IsInvalidConfiguration,
			wantBool: true,
		},
		{
			name:     "IsInvalidConfiguration with different error",
			err:      NewResourceNotFoundError("service", "not found"),
			checkFn:  IsInvalidConfiguration,
			wantBool: false,
		},
		{
			name:     "IsUnavailableMetrics with correct error",
			err:      ErrUnavailableMetrics,
			checkFn:  IsUnavailableMetrics,
			wantBool: true,
		},
		{
			name:     "IsUnavailableMetrics with different error",
			err:      NewResourceNotFoundError("deployment", "not found"),
			checkFn:  IsUnavailableMetrics,
			wantBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.checkFn(tt.err); got != tt.wantBool {
				t.Errorf("%s = %v, want %v", tt.name, got, tt.wantBool)
			}
		})
	}
}
