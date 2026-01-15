package varconfig

import (
	"context"
	"projeto-crud-credencials/internal/common/logger"
	"projeto-crud-credencials/internal/common/metrics"
	"testing"
)

// TestServiceCreate testa a criação de um VarConfig
func TestServiceCreate(t *testing.T) {
	// Mock repository
	mockRepo := &mockRepository{
		configs: make(map[string]VarConfig),
	}

	svc := NewService(mockRepo, &logger.DefaultLogger{}, &metrics.DefaultCollector{})

	t.Run("Create com dados válidos", func(t *testing.T) {
		config := VarConfig{
			OrgID:       123,
			BenchmarkID: "test",
			Payload:     map[string]any{"key": "value"},
		}

		result, err := svc.Create(context.Background(), config)
		if err != nil {
			t.Fatalf("Create falhou: %v", err)
		}

		if result.CreatedAt.IsZero() {
			t.Error("CreatedAt não foi definido")
		}
	})

	t.Run("Create com orgID inválido", func(t *testing.T) {
		config := VarConfig{
			OrgID:       0,
			BenchmarkID: "test",
			Payload:     map[string]any{},
		}

		_, err := svc.Create(context.Background(), config)
		if err == nil {
			t.Error("Create deveria retornar erro para orgID = 0")
		}
	})

	t.Run("Create com benchmarkID vazio", func(t *testing.T) {
		config := VarConfig{
			OrgID:       123,
			BenchmarkID: "",
			Payload:     map[string]any{},
		}

		_, err := svc.Create(context.Background(), config)
		if err == nil {
			t.Error("Create deveria retornar erro para benchmarkID vazio")
		}
	})

	t.Run("Create com payload nulo", func(t *testing.T) {
		config := VarConfig{
			OrgID:       123,
			BenchmarkID: "test",
			Payload:     nil,
		}

		_, err := svc.Create(context.Background(), config)
		if err == nil {
			t.Error("Create deveria retornar erro para payload nulo")
		}
	})
}

// TestServiceGetByID testa a busca de um VarConfig
func TestServiceGetByID(t *testing.T) {
	mockRepo := &mockRepository{}
	svc := NewService(mockRepo, &logger.DefaultLogger{}, &metrics.DefaultCollector{})

	t.Run("GetByID com parâmetros inválidos", func(t *testing.T) {
		_, err := svc.GetByID(context.Background(), 0, "test", 1)
		if err == nil {
			t.Error("GetByID deveria retornar erro para orgID = 0")
		}
	})
}

// mockRepository é um mock simples para testes
type mockRepository struct {
	configs map[string]VarConfig
}

func (m *mockRepository) FindAllByBenchmark(ctx context.Context, orgID int64, benchmarkID string) ([]VarConfig, error) {
	return nil, nil
}

func (m *mockRepository) Save(ctx context.Context, config VarConfig) (VarConfig, error) {
	return config, nil
}

func (m *mockRepository) FindByID(ctx context.Context, orgID int64, benchmarkID string, id int64) (VarConfig, error) {
	return VarConfig{}, nil
}

func (m *mockRepository) Update(ctx context.Context, config VarConfig) (VarConfig, error) {
	return config, nil
}

func (m *mockRepository) Delete(ctx context.Context, orgID int64, benchmarkID string, id int64) error {
	return nil
}
