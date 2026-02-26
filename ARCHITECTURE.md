# Guia de Arquitetura — Services, Handlers & Storage

Este documento descreve uma arquitetura, e organização de arquivos e diretórios, baseados no módulo
`github.com/Wizzi-Cloud/backend-design`.

O objetivo é consolidar o fluxo entre camadas, esclarecer responsabilidades,
registrar decisões ratificadas e servir como **referência única** para criação de
novos domínios e serviços **sem alterar conceitos ou padrões existentes**.

---

## 1. Visão Geral

A plataforma adota um fluxo **handler → service → storage**, com contratos explícitos
e governança rigorosa de dependências.

```
HTTP Request
    ↓
Handler
    ↓ (Service port — definido pelo handler)
Service
    ↓ (Repository port — definido pelo service)
Storage (driver)
    ↓
Database
```

Características:

- Fluxo síncrono e direto
- Dependências sempre apontam para dentro
- **Quem consome define o contrato**
- Wiring explícito no handler
- Domínios representam **contextos de negócio**, não entidades

---

## 2. Conceito Fundamental: Domínio ≠ Entidade

> **Domínio é um contexto de negócio (DDD).**
> **Entidades vivem dentro de domínios.**

Exemplo adotado no projeto:

- `core` → domínio central do sistema
- `user`, `customer` → agregados do domínio `core`

Essa decisão é **organizacional e intencional**, baseada em:

- regras que evoluem juntas
- ownership compartilhado
- convenção do time

---

## 3. Panorama End-to-End

```
Entrada externa (REST, Lambda, Scheduler)
        ↓
Handler em pkg/handler/domain/<domínio>/<agregado>
        # valida entrada, define contratos e executa wiring
        ↓
Service em internal/service/domain/<domínio>/<agregado>
        # aplica regras de negócio e depende apenas de ports
        ↓
Storage em internal/storage/<driver>/<domínio>/<agregado>
        # implementa ports usando o driver selecionado
```

Notas importantes:

- Handler **define o port do service** que consome
- Service **define o port do repository** que consome
- Service não conhece HTTP, handlers ou drivers
- Storage não conhece handlers nem regras de negócio
- Domínio é sempre a **primeira camada de organização**

---

## 4. Estrutura do Projeto

```
backend-design/
├─ dto/                                  # DTOs de entrada/saída dos handlers
├─ internal/
│  ├─ common/                            # infraestrutura compartilhada
│  ├─ service/
│  │  └─ domain/
│  │     └─ <domínio>/                  # ex: core
│  │        └─ <agregado>/              # ex: user, customer
│  │           ├─ ports.go              # port do repository (definido pelo service)
│  │           ├─ service.go            # implementação do service
│  │           └─ helpers/              # utilitários opcionais
│  │
│  └─ storage/
│     └─ <driver>/                      # ex: sqlite, postgres
│        ├─ common/                     # utilidades do driver
│        └─ <domínio>/
│           └─ <agregado>/
│              └─ repository.go         # implementação do repository port
│
├─ pkg/
│  ├─ models/                           # modelos de persistência
│  └─ handler/
│     └─ domain/
│        └─ <domínio>/                  # ex: core
│           └─ <agregado>/              # ex: user, customer
│              ├─ handler.go            # handlers HTTP
│              ├─ ports.go              # port do service (definido pelo handler)
│              ├─ path_parameter.go     # validação de path params
│              ├─ query_parameter.go    # validação de query params
│              └─ init.go               # wiring explícito
│
├─ routes/
├─ ARCHITECTURE.md
└─ go.mod
```

---

## 5. Contratos e Implementações (Exemplo Canônico)

### 5.1 Port do Service (definido pelo Handler)

**`pkg/handler/domain/core/user/ports.go`**

```go
type IUserService interface {
    List(ctx context.Context, role UserRole) ([]dto.UserDTO, error)
    GetByID(ctx context.Context, id, requesterID string, role UserRole) (dto.UserDTO, error)
}
```

> Este contrato representa **exatamente o que o handler precisa** para atender HTTP.

---

### 5.2 Port do Repository (definido pelo Service)

**`internal/service/domain/core/user/ports.go`**

```go
type IUserRepository interface {
    List(ctx context.Context) ([]models.User, error)
    GetByID(ctx context.Context, id string) (models.User, error)
}
```

> O domínio define este contrato sem qualquer dependência de HTTP ou drivers.

---

### 5.3 Implementação do Service

**`internal/service/domain/core/user/service.go`**

```go
type Service struct {
    repo IUserRepository
}

func NewService(repo IUserRepository) *Service {
    return &Service{repo: repo}
}
```

---

### 5.4 Implementação do Storage (SQLite)

**`internal/storage/sqlite/core/user/repository.go`**

```go
type Repository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
    return &Repository{db: db}
}
```

---

### 5.5 Wiring no Handler

**`pkg/handler/domain/core/user/init.go`**

```go
repo := userStorage.NewRepository(db)
svc  := userService.NewService(repo)
return NewHandler(svc)
```

> O wiring conecta **implementações concretas** aos **ports definidos pelo consumidor**.

---

## 6. Decisões Arquiteturais Ratificadas

- Domínios representam **contextos de negócio**, não entidades
- `core` é o domínio central atual
- User e Customer são agregados do domínio `core`
- **Handler define o port do service que consome**
- **Service define o port do repository que consome**
- Ports vivem próximos ao agregado que os define
- Wiring ocorre exclusivamente no `init.go` do handler
- Import cycles são evitados por **ports unidirecionais**

---

## 7. Governança de Dependências (depguard)

```yaml
linters-settings:
  depguard:
    rules:
      default:
        deny:
          - pkg: 'internal/storage/.*'
            desc: 'Storage só pode ser importado no wiring do handler'
      handler:
        allow:
          - 'internal/service/.*'
          - 'internal/storage/.*'
      service:
        allow:
          - 'internal/common/.*'
```

---

## 8. Qualidade, Testes e Observabilidade

- Services possuem testes de regra de negócio
- Repositórios implementam contract tests
- Métricas padronizadas:
  - `service_<domínio>_<agregado>_<caso>_latency_ms`
  - `storage_<driver>_<domínio>_<agregado>_errors_total`

---

## 9. Glossário

- **Domínio:** Contexto de negócio (DDD)
- **Agregado:** Conjunto coeso de entidades dentro de um domínio
- **Port:** Interface definida por quem consome
- **Wiring:** Montagem explícita das dependências

---

## Histórico

- **Out/2025** — Consolidação domain-first
