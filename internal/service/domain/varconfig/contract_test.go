package varconfig

import (
	"testing"
	"time"
)

// VarConfigRepositoryContract define testes de contrato que podem ser chamados por testes de implementação
func VarConfigRepositoryContract(t *testing.T, repo VarConfigRepository) {
	// Setup
	orgID := int64(123)
	benchmarkID := "benchmark_test"
	
	t.Run("Save e FindByID", func(t *testing.T) {
		// Arrange
		config := VarConfig{
			OrgID:       orgID,
			BenchmarkID: benchmarkID,
			Payload: map[string]any{
				"max_mem":      1024,
				"allowed_types": []string{"a", "b"},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Act
		saved, err := repo.Save(config)
		if err != nil {
			t.Fatalf("Save falhou: %v", err)
		}

		// Assert
		found, err := repo.FindByID(orgID, benchmarkID, saved.ID)
		if err != nil {
			t.Fatalf("FindByID falhou: %v", err)
		}

		if found.Payload == nil {
			t.Error("Payload não foi persistido corretamente")
		}
	})

	t.Run("FindAllByBenchmark", func(t *testing.T) {
		// Arrange
		config1 := VarConfig{
			OrgID:       orgID,
			BenchmarkID: benchmarkID,
			Payload:     map[string]any{"version": 1},
		}
		config2 := VarConfig{
			OrgID:       orgID,
			BenchmarkID: benchmarkID,
			Payload:     map[string]any{"version": 2},
		}

		// Act
		repo.Save(config1)
		repo.Save(config2)

		configs, err := repo.FindAllByBenchmark(orgID, benchmarkID)
		if err != nil {
			t.Fatalf("FindAllByBenchmark falhou: %v", err)
		}

		// Assert
		if len(configs) < 2 {
			t.Errorf("Esperava no mínimo 2 configs, obteve %d", len(configs))
		}
	})

	t.Run("Update", func(t *testing.T) {
		// Arrange
		config := VarConfig{
			OrgID:       orgID,
			BenchmarkID: benchmarkID,
			Payload:     map[string]any{"version": 1},
		}
		saved, _ := repo.Save(config)

		// Act
		saved.Payload = map[string]any{"version": 2}
		updated, err := repo.Update(saved)
		if err != nil {
			t.Fatalf("Update falhou: %v", err)
		}

		// Assert
		version := updated.Payload["version"].(float64)
		if int(version) != 2 {
			t.Errorf("Payload não foi atualizado: esperava 2, obteve %v", version)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		// Arrange
		config := VarConfig{
			OrgID:       orgID,
			BenchmarkID: benchmarkID,
			Payload:     map[string]any{"test": true},
		}
		saved, _ := repo.Save(config)

		// Act
		err := repo.Delete(orgID, benchmarkID, saved.ID)
		if err != nil {
			t.Fatalf("Delete falhou: %v", err)
		}

		// Assert
		_, err = repo.FindByID(orgID, benchmarkID, saved.ID)
		if err == nil {
			t.Error("Config não foi deletado")
		}
	})

	t.Run("Delete com ID inválido", func(t *testing.T) {
		// Act
		err := repo.Delete(orgID, benchmarkID, 99999)

		// Assert
		if err == nil {
			t.Error("Delete deveria retornar erro para ID inválido")
		}
	})

	t.Run("FindByID com ID inválido", func(t *testing.T) {
		// Act
		_, err := repo.FindByID(orgID, benchmarkID, 99999)

		// Assert
		if err == nil {
			t.Error("FindByID deveria retornar erro para ID inválido")
		}
	})
}
