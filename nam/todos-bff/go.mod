module github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff

go 1.25.0

tool github.com/99designs/gqlgen

require (
	github.com/99designs/gqlgen v0.17.89
	github.com/go-chi/chi/v5 v5.2.5
	github.com/google/wire v0.7.0
	github.com/graph-gophers/dataloader/v7 v7.1.3
	github.com/vektah/gqlparser/v2 v2.5.32
	google.golang.org/grpc v1.80.0
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	go.opentelemetry.io/otel v1.43.0 // indirect
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260120221211-b8f7ae30c516 // indirect
)

require (
	github.com/agnivade/levenshtein v1.2.1 // indirect
	github.com/goccy/go-yaml v1.19.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/sosodev/duration v1.4.0 // indirect
	github.com/tuannguyenandpadcojp/fresher26/nam/todos v0.0.0
	github.com/tuannguyenandpadcojp/fresher26/nam/users v0.0.0
	github.com/urfave/cli/v3 v3.7.0 // indirect
	golang.org/x/mod v0.35.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	golang.org/x/tools v0.44.0 // indirect
)

replace github.com/tuannguyenandpadcojp/fresher26/nam/todos => ../todos

replace github.com/tuannguyenandpadcojp/fresher26/nam/users => ../users
