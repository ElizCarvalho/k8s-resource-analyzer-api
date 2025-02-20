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
			name:     "erro com erro interno",
			resource: "deployment",
			message:  "não encontrado",
			err:      ErrResourceNotFound,
			want:     "deployment: não encontrado (recurso não encontrado)",
		},
		{
			name:     "erro sem erro interno",
			resource: "pod",
			message:  "configuração inválida",
			err:      nil,
			want:     "pod: configuração inválida",
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
		Message:  "não encontrado",
		Err:      err,
	}

	if got := re.Unwrap(); got != err {
		t.Errorf("ResourceError.Unwrap() = %v, want %v", got, err)
	}
}

func TestNewResourceNotFoundError(t *testing.T) {
	got := NewResourceNotFoundError("deployment", "não encontrado")
	if !errors.Is(got, ErrResourceNotFound) {
		t.Error("NewResourceNotFoundError() não retornou um ErrResourceNotFound")
	}
}

func TestNewInvalidMetricsError(t *testing.T) {
	got := NewInvalidMetricsError("cpu", "valor inválido")
	if !errors.Is(got, ErrInvalidMetrics) {
		t.Error("NewInvalidMetricsError() não retornou um ErrInvalidMetrics")
	}
}

func TestNewInvalidConfigurationError(t *testing.T) {
	got := NewInvalidConfigurationError("mimir", "url inválida")
	if !errors.Is(got, ErrInvalidConfiguration) {
		t.Error("NewInvalidConfigurationError() não retornou um ErrInvalidConfiguration")
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
			name:     "IsResourceNotFound com erro correto",
			err:      NewResourceNotFoundError("deployment", "não encontrado"),
			checkFn:  IsResourceNotFound,
			wantBool: true,
		},
		{
			name:     "IsResourceNotFound com erro diferente",
			err:      NewInvalidMetricsError("cpu", "valor inválido"),
			checkFn:  IsResourceNotFound,
			wantBool: false,
		},
		{
			name:     "IsInvalidMetrics com erro correto",
			err:      NewInvalidMetricsError("memory", "valor inválido"),
			checkFn:  IsInvalidMetrics,
			wantBool: true,
		},
		{
			name:     "IsInvalidMetrics com erro diferente",
			err:      NewResourceNotFoundError("pod", "não encontrado"),
			checkFn:  IsInvalidMetrics,
			wantBool: false,
		},
		{
			name:     "IsInvalidConfiguration com erro correto",
			err:      NewInvalidConfigurationError("mimir", "url inválida"),
			checkFn:  IsInvalidConfiguration,
			wantBool: true,
		},
		{
			name:     "IsInvalidConfiguration com erro diferente",
			err:      NewResourceNotFoundError("service", "não encontrado"),
			checkFn:  IsInvalidConfiguration,
			wantBool: false,
		},
		{
			name:     "IsUnavailableMetrics com erro correto",
			err:      ErrUnavailableMetrics,
			checkFn:  IsUnavailableMetrics,
			wantBool: true,
		},
		{
			name:     "IsUnavailableMetrics com erro diferente",
			err:      NewResourceNotFoundError("deployment", "não encontrado"),
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
