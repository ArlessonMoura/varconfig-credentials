# RESTWrapper - Guia Técnico Completo

## Visão Geral e Propósito

O **restwrapper** é um módulo Go especializado em fornecer uma camada de abstração unificada para tratamento de requisições e respostas HTTP, suportando múltiplos gateways web incluindo AWS API Gateway Lambda Proxy Integration e Gin Framework. 

### Propósito Principal
- **Wrapper HTTP**: Abstrai a complexidade do tratamento de requisições/respostas HTTP
- **Gerenciador de Respostas Padronizadas**: Fornece métodos consistentes para diferentes tipos de resposta HTTP
- **Middleware de Validação**: Integra validação automática de dados usando struct tags
- **Parser Universal**: Converte automaticamente entre diferentes tipos de dados (query params, path params, headers, body)

### Casos de Uso
- Aplicações Lambda com API Gateway
- Aplicações web com Gin Framework
- APIs REST que necessitam de validação estruturada
- Sistemas que requerem tratamento padronizado de erros
- Aplicações com autenticação baseada em JWT/claims
- Projetos que precisam suportar múltiplos gateways com o mesmo código

## Estrutura de Exportação

### Interfaces Principais

```go
// Validatable define a interface para estruturas que podem ser validadas
type Validatable interface {
    Validate() error
}

// IRequestWrapper define a interface para wrapping de requisições HTTP
type IRequestWrapper interface {
    BindHeaders(Validatable) error
    BindPathParams(Validatable) error
    BindQueryParams(Validatable) error
    BindBody(any) error
    RequesterEmail() (string, error)
    Method() string
    GetPathParam(string) (*string, bool)
}

// IResponseWrapper define a interface para wrapping de respostas HTTP
type IResponseWrapper interface {
    WriteSuccessResponse(int, any) error
    WriteClientErrorResponse(int, string) error
    WriteApplicationErrorResponse(apperror.IAppError) error
    WriteServerErrorResponse() error
    WriteServiceUnavailableErrorResponse(string) error
}
```

### Struct Principal

```go
// Wrapper unifica request e response wrappers
type Wrapper struct {
    RequestWrapper  IRequestWrapper
    ResponseWrapper IResponseWrapper
}
```

### Funções de Factory

```go
// Factory functions para API Gateway
func NewAPIGatewayWrapper(req eve.APIGatewayProxyRequest, res *eve.APIGatewayProxyResponse) *restwrapper.Wrapper
func NewAPIGatewayRequestWrapper(req eve.APIGatewayProxyRequest) *APIGatewayRequestWrapper
func NewAPIGatewayResponseWrapper(res *eve.APIGatewayProxyResponse) *APIGatewayResponseWrapper

// Factory functions para Gin Framework
func NewGinWrapper(req *gin.Context, res *eve.APIGatewayProxyResponse) *restwrapper.Wrapper
func NewGinRequestWrapper(req *gin.Context) *GinRequestWrapper
func NewGinResponseWrapper(res *eve.APIGatewayProxyResponse) *GinResponseWrapper

// Factory functions para componentes internos
func NewParser() *Parser

// Factory function para Handler Universal
func NewHandler(handler IRestHandler) *Handler
```

## Padrões de Design

### 1. Interface Segregation Principle (ISP)
- Separa responsabilidades entre `IRequestWrapper` e `IResponseWrapper`
- Interface `Validatable` para validação independente

### 2. Factory Pattern
- `NewAPIGatewayWrapper()` e `NewGinWrapper()` funcionam como factories principais
- Cria instâncias configuradas com implementações específicas
- `NewHandler()` cria handlers universais compatíveis com múltiplos gateways

### 3. Strategy Pattern
- Diferentes implementações de wrappers para diferentes gateways (API Gateway vs Gin)
- Parser implementa estratégias diferentes para map[string]string e map[string][]string
- Handler universal implementa estratégia para múltiplos gateways

### 4. Template Method Pattern
- Métodos de response seguem template: criar body → serializar → definir headers → escrever response
- Handlers universais seguem template: criar wrapper → executar lógica → retornar resposta

### 5. Decorator Pattern
- Wrapper adiciona funcionalidades (validação, parsing, CORS) às requisições/respostas originais
- Handler universal adiciona camada de abstração sobre implementações específicas

## Exemplos de Implementação

### 1. Inicialização do Módulo

#### API Gateway (Serverless)
```go
package main

import (
    "github.com/Wizzi-Cloud/restwrapper/apigw"
    eve "github.com/aws/aws-lambda-go/events"
)

func HandleRequest(req eve.APIGatewayProxyRequest) (eve.APIGatewayProxyResponse, error) {
    var res eve.APIGatewayProxyResponse
    
    // Inicialização do wrapper
    wrapper := apigw.NewAPIGatewayWrapper(req, &res)
    
    // Processamento...
    
    return res, nil
}
```

#### Gin Framework (Web Applications)
```go
package main

import (
    "github.com/Wizzi-Cloud/restwrapper/gin"
    "github.com/gin-gonic/gin"
    eve "github.com/aws/aws-lambda-go/events"
)

func HandleGinRequest(c *gin.Context) {
    var res eve.APIGatewayProxyResponse
    
    // Inicialização do wrapper
    wrapper := gin.NewGinWrapper(c, &res)
    
    // Processamento...
    
    // Aplicar resposta ao contexto Gin
    // (ver seção de Handler Universal)
}
```

### 2. Estruturas de Dados com Validação

#### Tags Suportadas por Gateway

**API Gateway (JSON tags):**
```go
type QueryParameters struct {
    Page              *uint     `json:"page" validate:"min=1"`
    Limit             *uint     `json:"limit" validate:"min=1,max=100"`
    CollaboratorEmail *string   `json:"collaboratorEmail" validate:"required,email"`
    TeamID            *uint     `json:"teamID"`
    FromNow           *bool     `json:"fromNow"`
    Filters           *[]string `json:"filter"`
}
```

**Gin Framework (Form/URI/Header tags):**
```go
type QueryParameters struct {
    Page              *uint     `form:"page" validate:"min=1"`
    Limit             *uint     `form:"limit" validate:"min=1,max=100"`
    CollaboratorEmail *string   `form:"collaboratorEmail" validate:"required,email"`
    TeamID            *uint     `form:"teamID"`
    FromNow           *bool     `form:"fromNow"`
    Filters           *[]string `form:"filter"`
}

type PathParameters struct {
    OrgID  string `uri:"org_id" validate:"required"`
    UserID string `uri:"user_id" validate:"required"`
}

type CustomHeaders struct {
    Authorization string `header:"authorization" validate:"required"`
    XRequestID    string `header:"x-request-id"`
    ContentType   string `header:"content-type"`
}
```

// Estrutura para path parameters (API Gateway)
type PathParameters struct {
    OrgID  string `json:"org_id" validate:"required"`
    UserID string `json:"user_id" validate:"required"`
}

func (pp *PathParameters) Validate() error {
    if pp.OrgID == "" {
        return errors.New("org_id is required")
    }
    if pp.UserID == "" {
        return errors.New("user_id is required")
    }
    return nil
}

// Estrutura para body (ambos gateways)
type CreateUserRequest struct {
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"min=18,max=120"`
}

func (cur *CreateUserRequest) Validate() error {
    // Validação customizada se necessário
    return nil
}
```

### 3. Handler Universal - Multi-Gateway Pattern

#### Interface do Handler Universal
```go
type IRestHandler interface {
    Handle(context *restwrapper.Wrapper)
}
```

#### Implementação do Handler Universal
```go
type UserHandler struct {
    userService UserService
}

func (h *UserHandler) Handle(wrapper *restwrapper.Wrapper) {
    // Extrair email do requester
    email, err := wrapper.RequestWrapper.RequesterEmail()
    if err != nil {
        wrapper.ResponseWrapper.WriteClientErrorResponse(401, "unauthorized")
        return
    }
    
    // Binding de path parameters
    var pathParams PathParameters
    if err := wrapper.RequestWrapper.BindPathParams(&pathParams); err != nil {
        wrapper.ResponseWrapper.WriteClientErrorResponse(400, err.Error())
        return
    }
    
    // Verificar método HTTP
    switch wrapper.RequestWrapper.Method() {
    case "GET":
        h.handleGetUsers(wrapper, pathParams)
    case "POST":
        h.handleCreateUser(wrapper, pathParams)
    case "PUT":
        h.handleUpdateUser(wrapper, pathParams)
    case "DELETE":
        h.handleDeleteUser(wrapper, pathParams)
    default:
        wrapper.ResponseWrapper.WriteClientErrorResponse(405, "method not allowed")
    }
}
```

#### Uso em API Gateway
```go
func LambdaHandler(req eve.APIGatewayProxyRequest) (eve.APIGatewayProxyResponse, error) {
    userHandler := &UserHandler{userService: NewUserService()}
    handler := resthandler.NewHandler(userHandler)
    
    return handler.HandleApiGateway(req)
}
```

#### Uso em Gin Framework
```go
func GinUserHandler(c *gin.Context) {
    userHandler := &UserHandler{userService: NewUserService()}
    handler := resthandler.NewHandler(userHandler)
    
    handler.HandleGin(c)
}
```

### 4. Métodos Avançados dos Wrappers

#### Acesso Direto a Path Parameters
```go
func handleWithDirectAccess(wrapper *restwrapper.Wrapper) {
    // Acesso direto sem binding completo
    if orgID, exists := wrapper.RequestWrapper.GetPathParam("org_id"); exists {
        fmt.Printf("Org ID: %s\n", *orgID)
    } else {
        wrapper.ResponseWrapper.WriteClientErrorResponse(400, "org_id required")
        return
    }
    
    // Verificar método HTTP
    switch wrapper.RequestWrapper.Method() {
    case "GET":
        // Lógica para GET
    case "POST":
        // Lógica para POST
    }
}
```

#### Headers Customizados (Gin vs API Gateway)
```go
// API Gateway - usa JSON tags
type APIGatewayHeaders struct {
    Authorization string `json:"authorization" validate:"required"`
    XRequestID    string `json:"x-request-id"`
    ContentType   string `json:"content-type"`
}

// Gin - usa Header tags
type GinHeaders struct {
    Authorization string `header:"authorization" validate:"required"`
    XRequestID    string `header:"x-request-id"`
    ContentType   string `header:"content-type"`
}

func handleWithHeaders(wrapper *restwrapper.Wrapper) {
    var headers APIGatewayHeaders  // ou GinHeaders dependendo do gateway
    if err := wrapper.RequestWrapper.BindHeaders(&headers); err != nil {
        wrapper.ResponseWrapper.WriteClientErrorResponse(400, err.Error())
        return
    }
    
    fmt.Printf("Request ID: %s\n", headers.XRequestID)
}
```

### 5. Tratamento de Erros Específico

```go
// Erros customizados usando apperror
func handleBusinessLogic() error {
    return apperror.NewAppError(
        "BUSINESS_RULE_VIOLATION",
        "User cannot be created with this email domain",
        errors.New("business rule violation"),
    )
}

// Tratamento no handler
func handleWithError(wrapper *restwrapper.Wrapper) {
    err := someBusinessOperation()
    if err != nil {
        // Erro de aplicação com logging automático
        wrapper.ResponseWrapper.WriteApplicationErrorResponse(err)
        return
    }
    
    // Erro genérico de servidor
    wrapper.ResponseWrapper.WriteServerErrorResponse()
    
    // Serviço indisponível
    wrapper.ResponseWrapper.WriteServiceUnavailableErrorResponse("Database connection failed")
    
    // Erro de cliente
    wrapper.ResponseWrapper.WriteClientErrorResponse(400, "Invalid input data")
}
```

### 6. Testes e Validação

#### Estrutura de Testes
```go
// Teste unitário completo
func TestWrapperIntegration(t *testing.T) {
    // Setup da requisição API Gateway
    req := eve.APIGatewayProxyRequest{
        HTTPMethod: "GET",
        Path:       "org/{org_id}/users/{user_id}",
        RequestContext: eve.APIGatewayProxyRequestContext{
            Authorizer: map[string]interface{}{
                "claims": map[string]interface{}{
                    "email": "test@example.com",
                },
            },
        },
        PathParameters: map[string]string{
            "org_id": "123",
            "user_id": "456",
        },
        QueryStringParameters: map[string]string{
            "page":  "1",
            "limit": "10",
            "filter": "active",
        },
        MultiValueQueryStringParameters: map[string][]string{
            "filter": {"active", "verified"},
        },
    }
    
    var res eve.APIGatewayProxyResponse
    wrapper := apigw.NewAPIGatewayWrapper(req, &res)
    
    // Testar email extraction
    email, err := wrapper.RequestWrapper.RequesterEmail()
    assert.NoError(t, err)
    assert.Equal(t, "test@example.com", email)
    
    // Testar path params binding
    var pathParams PathParameters
    err = wrapper.RequestWrapper.BindPathParams(&pathParams)
    assert.NoError(t, err)
    assert.Equal(t, "123", pathParams.OrgID)
    assert.Equal(t, "456", pathParams.UserID)
    
    // Testar query params binding
    var queryParams QueryParameters
    err = wrapper.RequestWrapper.BindQueryParams(&queryParams)
    assert.NoError(t, err)
    assert.Equal(t, uint(1), *queryParams.Page)
    assert.Equal(t, uint(10), *queryParams.Limit)
    assert.Contains(t, *queryParams.Filters, "active")
    assert.Contains(t, *queryParams.Filters, "verified")
}
```

#### Testes com Gin Framework
```go
func TestGinWrapper(t *testing.T) {
    // Setup do Gin context
    gin.SetMode(gin.TestMode)
    router := gin.New()
    
    // Simular requisição
    req, _ := http.NewRequest("GET", "/users/123?page=1&limit=10", nil)
    req.Header.Set("Authorization", "Bearer token")
    req.Header.Set("X-Request-ID", "test-123")
    
    // Criar contexto manualmente
    c, _ := gin.CreateTestContext(httptest.NewRecorder())
    c.Request = req
    c.Params = gin.Params{{Key: "id", Value: "123"}}
    c.Set("email", "test@example.com")
    
    var res eve.APIGatewayProxyResponse
    wrapper := ginWrapper.NewGinWrapper(c, &res)
    
    // Testar métodos
    assert.Equal(t, "GET", wrapper.RequestWrapper.Method())
    
    id, exists := wrapper.RequestWrapper.GetPathParam("id")
    assert.True(t, exists)
    assert.Equal(t, "123", *id)
    
    email, err := wrapper.RequestWrapper.RequesterEmail()
    assert.NoError(t, err)
    assert.Equal(t, "test@example.com", email)
}
```

### 7. Guia para IAs (IA Implementation Context)

### Regras Obrigatórias
- **Sempre** use a struct `Wrapper` como ponto de entrada principal
- **Nunca** ignore o erro retornado pelos métodos `Bind*` - sempre trate com `WriteClientErrorResponse`
- **Sempre** implemente o método `Validate()` em structs que usam métodos `Bind*`
- **Sempre** use ponteiros (`*string`, `*uint`, `*bool`) para campos opcionais em query parameters
- **Nunca** acesse diretamente os campos internos dos wrappers - use apenas os métodos da interface
- **Escolha** as tags corretas: `json:"..."` para API Gateway, `form:"..."`/`uri:"..."`/`header:"..."` para Gin

### Padrões de Código
- **Validação**: Implementar `Validate()` com early returns para múltiplos erros
- **Binding**: Usar `BindQueryParams()` para query strings, `BindPathParams()` para path params, `BindBody()` para JSON body
- **Respostas**: Usar `WriteSuccessResponse()` para sucesso, `WriteApplicationErrorResponse()` para erros de negócio
- **Autenticação**: Sempre chamar `RequesterEmail()` primeiro para validar autorização
- **Multi-Gateway**: Use Handler Universal (`IRestHandler`) para código compatível com ambos gateways

### Tratamento de Erros
- **Erros de Cliente (4xx)**: Use `WriteClientErrorResponse(statusCode, message)`
- **Erros de Servidor (5xx)**: Use `WriteApplicationErrorResponse(appError)` para erros com contexto
- **Erro Genérico**: Use `WriteServerErrorResponse()` para erros inesperados
- **Serviço Indisponível**: Use `WriteServiceUnavailableErrorResponse(message)`

### Performance e Boas Práticas
- **Validação**: Prefira validação no método `Validate()` em vez de nos handlers
- **Pointers**: Use ponteiros para valores opcionais para distinguir "não informado" de "valor zero"
- **Slices**: Para arrays em query params, use `*[]string` com tag `json:"fieldName"` (API Gateway) ou `form:"fieldName"` (Gin)
- **Logging**: Erros passados para `WriteApplicationErrorResponse()` são logados automaticamente
- **Handler Universal**: Prefira `IRestHandler` para máxima portabilidade entre gateways

### Integração com Ecossistema
- **apperror**: Use `apperror.IAppError` para erros de negócio com mensagens públicas/privadas
- **validator**: Use tags `validate` nos structs para validação automática
- **AWS Lambda**: O módulo é otimizado para `APIGatewayProxyRequest/Response`
- **Gin Framework**: Integração nativa com `*gin.Context` e tags de binding do Gin

### Anti-Patterns a Evitar
- **Nunca** faça parsing manual do JSON body
- **Nunca** ignore validação de estruturas
- **Nunca** retorne erros diretamente do handler Lambda - sempre use os métodos de response
- **Nunca** acesse `req.Body` diretamente - use `BindBody()`
- **Nunca** misture tags de diferentes gateways na mesma struct

## Dependências e Compatibilidade

### Dependências Principais
- `github.com/Wizzi-Cloud/apperror` - Tratamento de erros de aplicação
- `github.com/Wizzi-Cloud/core_interfaces` - Interfaces de validação
- `github.com/aws/aws-lambda-go` - AWS Lambda events
- `github.com/go-playground/validator/v10` - Validação de structs
- `github.com/gin-gonic/gin` - Web framework para aplicações HTTP

### Versão Go
- **Mínima**: Go 1.23.7
- **Recomendada**: Go 1.21+

### Compatibilidade
- **AWS API Gateway**: Full suporte para Lambda Proxy Integration
- **Gin Framework**: Full suporte com integração nativa
- **Outros Gateways**: Extensível via implementação das interfaces
- **Testing**: Inclui testes unitários completos

### Instalação
```bash
go get github.com/Wizzi-Cloud/restwrapper
```

### Configuração Mínima
```go
// Para API Gateway
import "github.com/Wizzi-Cloud/restwrapper/apigw"

// Para Gin Framework
import "github.com/Wizzi-Cloud/restwrapper/gin"

// Para Handler Universal
import "github.com/Wizzi-Cloud/restwrapper/handler"
```

## Guia de Decisão e Migração

### Quando Usar API Gateway
- **Aplicações Serverless**: Lambda functions com API Gateway
- **Infraestrutura AWS**: Integração com outros serviços AWS
- **Pay-per-use**: Custos baseados em requisições
- **Escalabilidade Serverless**: Escalonamento automático

### Quando Usar Gin Framework
- **Aplicações Web Traditional**: Servidores web convencionais
- **Controle Total**: Sobre infraestrutura e deployment
- **Performance Local**: Menos overhead que serverless
- **Microserviços**: Com containers ou VMs

### Migração entre Gateways
#### De API Gateway para Gin
```go
// Antes (API Gateway)
type QueryParams struct {
    Page *uint `json:"page" validate:"min=1"`
}

// Depois (Gin)
type QueryParams struct {
    Page *uint `form:"page" validate:"min=1"`
}

// Handler permanece o mesmo com IRestHandler
type UserHandler struct { /* ... */ }

func (h *UserHandler) Handle(wrapper *restwrapper.Wrapper) {
    // Lógica idêntica para ambos gateways
}
```

#### De Gin para API Gateway
```go
// Antes (Gin)
type PathParams struct {
    ID string `uri:"id" validate:"required"`
}

// Depois (API Gateway)
type PathParams struct {
    ID string `json:"id" validate:"required"`
}
```

### Handler Universal - Máxima Portabilidade
```go
// Estruturas separadas por gateway
type APIGatewayQueryParams struct { /* json tags */ }
type GinQueryParams struct { /* form tags */ }

// Handler único que detecta o gateway
type UniversalHandler struct{}

func (h *UniversalHandler) Handle(wrapper *restwrapper.Wrapper) {
    // Detectar tipo de wrapper e usar structs apropriadas
    switch wrapper.RequestWrapper.(type) {
    case *apigw.APIGatewayRequestWrapper:
        var params APIGatewayQueryParams
        wrapper.RequestWrapper.BindQueryParams(&params)
    case *ginWrapper.GinRequestWrapper:
        var params GinQueryParams
        wrapper.RequestWrapper.BindQueryParams(&params)
    }
    
    // Lógica de negócio comum
}
```

## Considerações de Performance

### Memory Usage
- Parser usa reflection otimizado para minimizações de alocações
- Reuse de instâncias de Parser quando possível
- Ponteiros reduzem cópias desnecessárias

### CPU Optimization
- Validação acontece durante binding (single pass)
- CORS headers pré-definidos para evitar alocações
- JSON marshaling otimizado para respostas

### Scalability
- Stateless design permite scaling horizontal
- Sem estado global compartilhado
- Thread-safe por design (immutabilidade)
- Handler Universal facilita scaling entre diferentes arquiteturas

## Arquitetura e Estrutura do Projeto

### Organização dos Pacotes
```
restwrapper/
├── wrapper.go              # Interfaces principais e struct Wrapper
├── wrapper_test.go         # Testes integrados
├── apigw/                  # API Gateway implementation
│   ├── wrapper.go          # Factory function
│   ├── request/            # Request wrapper para API Gateway
│   │   └── request.go
│   ├── response/           # Response wrapper para API Gateway
│   │   └── response.go
│   └── utils/              # Utilitários compartilhados
│       ├── parser.go       # Parser universal
│       └── parser_test.go  # Testes do parser
├── gin/                    # Gin Framework implementation
│   ├── wrapper.go          # Factory function
│   ├── request/            # Request wrapper para Gin
│   │   └── request.go
│   └── response/           # Response wrapper para Gin
│       └── response.go
└── handler/                # Handler Universal
    └── handler.go          # Implementação universal multi-gateway
```

### Fluxo de Arquitetura
1. **Request** → Gateway específico (API Gateway/Gin)
2. **Wrapper Factory** → Cria Wrapper apropriado
3. **Handler Universal** → Processa requisição de forma agnóstica
4. **Business Logic** → Executa lógica de negócio
5. **Response** → Retorna resposta formatada

### Extensibilidade
- **Novos Gateways**: Implementar interfaces `IRequestWrapper` e `IResponseWrapper`
- **Custom Parsers**: Estender `Parser` para novos formatos
- **Middleware**: Adicionar camadas de validação customizadas
- **Response Formats**: Suportar diferentes formatos de resposta

## Conclusão

O **restwrapper** é uma solução completa e madura para tratamento de requisições HTTP em Go, oferecendo:

### ✅ **Vantagens Principais**
- **Multi-Gateway**: Suporte nativo para API Gateway e Gin Framework
- **Handler Universal**: Código portável entre diferentes arquiteturas
- **Validação Integrada**: Sistema robusto de validação com struct tags
- **Tratamento de Erros**: Padronização e logging automático
- **Performance**: Otimizado para alta performance e baixo overhead
- **Testabilidade**: Arquitetura testável com testes completos

### 🎯 **Casos de Uso Ideais**
- Projetos que precisam suportar múltiplas arquiteturas
- Migração gradual entre serverless e traditional
- APIs que requerem validação robusta
- Equipes que buscam consistência no tratamento de HTTP

### 🚀 **Próximos Passos**
1. Escolher o gateway adequado para seu projeto
2. Implementar structs com tags corretas
3. Usar Handler Universal para máxima portabilidade
4. Escrever testes para sua implementação
5. Considerar migração futura entre gateways

O módulo está pronto para produção em ambientes críticos, com cobertura completa de testes e documentação detalhada para facilitar adoção e manutenção.
