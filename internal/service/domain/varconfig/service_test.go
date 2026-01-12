package varconfig

import (
	"testing"
)

// TestServiceCreate testa a criação de um VarConfig
func TestServiceCreate(t *testing.T) {
	// Mock repository
	mockRepo := &mockRepository{
		configs: make(map[string]VarConfig),
	}

	svc := NewService(mockRepo)

	t.Run("Create com dados válidos", func(t *testing.T) {
		config := VarConfig{
			OrgID:       123,
			BenchmarkID: "test",
			Payload:     map[string]any{"key": "value"},
		}

		result, err := svc.Create(config)
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

		_, err := svc.Create(config)
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

		_, err := svc.Create(config)
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

		_, err := svc.Create(config)
		if err == nil {
			t.Error("Create deveria retornar erro para payload nulo")
		}
	})
}

// TestServiceGetByID testa a busca de um VarConfig
func TestServiceGetByID(t *testing.T) {
	mockRepo := &mockRepository{}
	svc := NewService(mockRepo)

	t.Run("GetByID com parâmetros inválidos", func(t *testing.T) {
		_, err := svc.GetByID(0, "test", 1)
		if err == nil {
			t.Error("GetByID deveria retornar erro para orgID = 0")
		}
	})
}

// mockRepository é um mock simples para testes
type mockRepository struct {
	configs map[string]VarConfig
}

func (m *mockRepository) FindAllByBenchmark(orgID int64, benchmarkID string) ([]VarConfig, error) {
	return nil, nil
}

func (m *mockRepository) Save(config VarConfig) (VarConfig, error) {
	return config, nil
}

func (m *mockRepository) FindByID(orgID int64, benchmarkID string, id int64) (VarConfig, error) {
	return VarConfig{}, nil
}

func (m *mockRepository) Update(config VarConfig) (VarConfig, error) {
	return config, nil
}

func (m *mockRepository) Delete(orgID int64, benchmarkID string, id int64) error {
	return nil
}
