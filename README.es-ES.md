

![Gopher](assets/gopher.png)
[![CircleCI](https://circleci.com/gh/alexfalkowski/go-service.svg?style=shield)](https://circleci.com/gh/alexfalkowski/go-service)
[![codecov](https://codecov.io/gh/alexfalkowski/go-service/graph/badge.svg?token=AGP01JOTM0)](https://codecov.io/gh/alexfalkowski/go-service)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexfalkowski/go-service/v2)](https://goreportcard.com/report/github.com/alexfalkowski/go-service/v2)
[![Go Reference](https://pkg.go.dev/badge/github.com/alexfalkowski/go-service/v2.svg)](https://pkg.go.dev/github.com/alexfalkowski/go-service/v2)
[![Stability: Active](https://masterminds.github.io/stability/active.svg)](https://masterminds.github.io/stability/active.html)

# 🧰 Go Service

`github.com/alexfalkowski/go-service/v2` es un marco de trabajo/biblioteca con enfoques propios para construir servicios en Go con una configuración consistente para inyección de dependencias, transportes, telemetría, criptografía, etc.

Este repositorio es principalmente una **biblioteca de paquetes** (sin binario de nivel superior `cmd/`). Los servicios construidos sobre él suelen definir su propio paquete `main` en otro lugar e importar este módulo.

Se espera que los servicios de ejecución prolongada comiencen desde [`go-service-template`](https://github.com/alexfalkowski/go-service-template), mientras que los comandos cliente de corta duración comienzan desde [`go-client-template`](https://github.com/alexfalkowski/go-client-template). Ambos componen los paquetes del módulo de alto nivel de este repositorio. Estos son los caminos principales soportados. La composición paquete por paquete a nivel inferior sigue disponible, pero es un modo avanzado y puede requerir registro manual adicional.

---

## 🚀 Instalación

Para un nuevo servicio de ejecución prolongada, comience desde `go-service-template` para que el `main` de la aplicación, la conexión del comando servidor, los fixtures de configuración y la composición estándar de módulos se generen juntos. Para un comando de control, migración o procesamiento por lotes de corta duración, comience desde `go-client-template`; demuestra `cli.Application.AddClient`, `module.Client` y el trabajo de inicio del ciclo de vida.

Para el uso directo de paquetes en un módulo existente, agregue la dependencia de la biblioteca con la ruta del módulo versionado:

```sh
go get github.com/alexfalkowski/go-service/v2
```

Utilice la versión de Go declarada en `go.mod` o una posterior al instalar o compilar este módulo.

---

## 🧩 Inyección de Dependencias (Fx)

El marco de trabajo está diseñado alrededor de la inyección de dependencias y utiliza [Uber Fx](https://github.com/uber-go/fx) (y Dig bajo el capó). La mayoría de los subsistemas exponen módulos Fx que compone en su servicio.

Si es nuevo en Fx, vale la pena leer primero su documentación y ejemplos.

### Paquetes de módulos (Bundles)

El paquete `module` expone tres paquetes de alto nivel:

- `module.Library` para fundaciones compartidas (entorno, compresión, codificación, criptografía, tiempo, configuración de búfer de sincronización, id)
- `module.Server` para procesos de servidor (Library + configuración, transportes, telemetría, depuración, salud, etc.)
- `module.Client` para procesos de corta duración/por lotes/cliente (Library + configuración, telemetría, sql, hooks, etc.)

Estos paquetes son el valor predeterminado destinado para servicios generados desde `go-service-template`. Manejan el registro interno esperado por el marco de trabajo, por lo que la mayoría de los servicios no necesitan conectar manualmente helpers de transporte o ciclo de vida a nivel inferior.

### Ejemplo mínimo de arranque CLI

Este repositorio es una biblioteca, por lo que su binario suele estar en otro módulo. Un `main` típico usa `cli.Application` y compone paquetes de módulos:

```go
package main

import (
    "github.com/alexfalkowski/go-service/v2/cli"
    "github.com/alexfalkowski/go-service/v2/context"
    "github.com/alexfalkowski/go-service/v2/module"
    "github.com/alexfalkowski/go-service/v2/os"
)

func main() {
    app := cli.NewApplication(func(commander cli.Commander) {
        serve := commander.AddServer("serve", "Run the service", module.Server)
        serve.AddConfig("file:./config.yaml") // agrega la bandera de configuración `-config` / `-c` con este valor predeterminado
    })

    os.Exit(app.RunCode(context.Background()))
}
```

El valor predeterminado `file:./config.yaml` anterior espera un archivo de configuración no vacío. Una configuración mínima de servidor puede comenzar con el entorno más un transporte habilitado:

```yaml
environment: development
transport:
  http:
    address: tcp://localhost:8000
    timeout: 10s
```

Use `app.RunCode(context.Background())` desde `main` al salir del proceso. Devuelve `os.ExitCodeSuccess` en caso de éxito, devuelve un código de salida de apagado no nulo solicitado como `os.ExitCodeServeFailure`, y devuelve `os.ExitCodeFailure` para otros errores. Use `app.Run(context.Background())` en pruebas o código de incrustación que necesite inspeccionar el error devuelto.

---

## 🖥️ CLI

Los servicios comúnmente exponen dos formas de comando:

- **Server**: proceso de daemon de ejecución prolongada
- **Client**: proceso de control/admin de corta duración

El marco de trabajo usa [acmd](https://github.com/cristalhq/acmd). El `main` de su servicio típicamente conecta módulos Fx + comandos.

> Este repositorio intencionalmente no distribuye un `main` listo para ejecutar: proporciona los bloques de construcción. En el uso normal, las aplicaciones de servidor las consumen a través de `go-service-template` más `module.Server`, mientras que los comandos de corta duración usan `go-client-template` más `module.Client`, en lugar de conectar cada subsistema manualmente.

---

## 🗂️ Estructura del repositorio

El repositorio está intencionalmente dividido entre composición de servicios de alto nivel y helpers reutilizables a nivel inferior:

- `module/` expone los paquetes Fx con enfoques propios (`Library`, `Server`, `Client`)
- `config/` define la forma estándar de configuración de nivel superior más las proyecciones utilizadas por la conexión de módulos
- paquetes de características como `cache/`, `crypto/`, `database/sql/`, `feature/`, `telemetry/`, `time/`, e `id/` proporcionan configuración, constructores y módulos Fx para un subsistema
- `net/...` contiene helpers de protocolo a nivel inferior y primitivas reutilizables (`net/http`, `net/grpc`, helpers de metadatos/encabezados, alias de protocolo de salud gRPC, y `net/server`)
- `transport/...` contiene la capa de transporte de servicio de nivel superior: pilas HTTP/gRPC compuestas, middleware de política, endpoints operativos y módulos específicos de transporte
- `internal/test/` contiene el entorno de prueba compartido y fixtures utilizados entre paquetes

Como regla general: si desea primitivas de protocolo o helpers compartidos, comience en `net/...`; si desea conexión de servicios y política de middleware, comience en `transport/...`. Los helpers compartidos de metadatos, encabezados y ciclo de vida viven bajo `net/...`, incluyendo `net/http/meta`, `net/grpc/meta`, `net/header`, y `net/server.Register`.

Para la mayoría de los autores de servicios, el punto de partida correcto sigue siendo los paquetes de módulos de alto nivel en lugar de estos paquetes a nivel inferior directamente.

---

## ⚙️ Configuración

### Formatos de configuración soportados

El decodificador de configuración soporta:

- JSON
- HJSON (`github.com/hjson/hjson-go/v4`)
- TOML (`github.com/BurntSushi/toml`)
- YAML (`go.yaml.in/yaml/v3`)

### Selección de la fuente de configuración (banderas `-config` / `-c`)

La entrada de configuración se enruta mediante banderas llamadas `-config` y `-c`:

- `file:<path>`
  Leer configuración desde un archivo en `<path>`; el parser se selecciona desde la extensión del archivo (`.json`, `.hjson`, `.yaml`, `.toml`).

- `env:<ENV_VAR>`
  Leer configuración desde la variable de entorno `<ENV_VAR>`. El valor de la variable de entorno debe estar formateado como:

  `"<extension>:<contenido-base64>"`

  Formato de ejemplo: `yaml:ZW52aXJvbm1lbnQ6IGRldmVsb3BtZW50Cg==`

  Comandos de ejemplo:

  ```sh
  # Linux (GNU base64)
  export SERVICE_CONFIG="yaml:$(base64 -w 0 < ./config.yaml)"
  ./your-service serve -config env:SERVICE_CONFIG
  ```

  ```sh
  # macOS/BSD base64
  export SERVICE_CONFIG="yaml:$(base64 < ./config.yaml | tr -d '\n')"
  ./your-service serve -c env:SERVICE_CONFIG
  ```

  HJSON funciona de la misma manera, por ejemplo `hjson:<contenido-base64>`.

  El helper del repositorio `make kind=configs/config encode-config` usa GNU `base64 -w 0`; en macOS/BSD, use `base64 | tr -d '\n'` para la carga útil de una sola línea equivalente.

- Los prefijos `kind:location` explícitos no soportados fallan al inicio en lugar de caer en otra fuente.

- Los valores sin prefijo, incluido un valor vacío, caen en **búsqueda predeterminada**, buscando:

  `<serviceName>.{yaml,hjson,toml,json}`

  La búsqueda predeterminada verifica extensiones primero (`.yaml`, `.hjson`, `.toml`, `.json`), y para cada extensión verifica:
  - directorio ejecutable
  - `$XDG_CONFIG_HOME/<serviceName>/` (vía `os.UserConfigDir()`)
  - `/etc/<serviceName>/`

> [!IMPORTANT]
> Debido a que el directorio de configuración del usuario es parte de esa búsqueda, se espera que los entornos de ejecución que usan la búsqueda predeterminada proporcionen `HOME` o `XDG_CONFIG_HOME`. Los servicios que no pueden depender de esas variables de entorno deben pasar una fuente explícita `-config file:<path>` o `-config env:<ENV_VAR>`.

### Decodificación tipada y validación

En tiempo de ejecución, los servicios típicamente decodifican en una estructura (a menudo incorporando `config.Config`) y la validan usando `go-playground/validator`.

La biblioteca proporciona un helper `config.NewConfig[T]` que:

- decodifica en `*T`
- rechaza un valor decodificado "vacío" (protege contra iniciar con una configuración de valor cero)
- valida la configuración decodificada

La detección de vacío usa semánticas de valor cero y soporta tipos de configuración que contienen mapas, slices u otros campos no comparables.

Ejemplo:

```go
type WorkerConfig struct {
    Queue string `yaml:"queue" json:"queue" toml:"queue" validate:"required"`
}

type AppConfig struct {
    Worker         *WorkerConfig `yaml:"worker" json:"worker" toml:"worker" validate:"required"`
    *config.Config `yaml:",inline" json:",inline" toml:",inline" validate:"required"`
}

func loadConfig(decoder config.Decoder, validator *config.Validator) (*AppConfig, error) {
    return config.NewConfig[AppConfig](decoder, validator)
}

func sharedConfig(cfg *AppConfig) *config.Config {
    return cfg.Config
}

func workerConfig(cfg *AppConfig) *WorkerConfig {
    return cfg.Worker
}

var AppConfigModule = di.Module(
    di.Constructor(config.NewConfig[AppConfig]),
    di.Decorate(sharedConfig),
    di.Constructor(workerConfig),
)
```

Composite `AppConfigModule` junto a `module.Server` o `module.Client`. El decorador proyecta el `*config.Config` compartido incorporado en el gráfico estándar, para que la configuración específica del servicio se decodifique una vez mientras las proyecciones existentes de transporte, SQL y telemetría continúan funcionando. Agregue constructores como `workerConfig` para sub-configuraciones poseídas por el servicio.

### La forma estándar de configuración de nivel superior

El tipo de configuración de nivel superior canónico es `config.Config` (en `config/config.go`). Contiene:

- `debug`, `cache`, `crypto`, `feature`, `hooks`, `id`, `sql`, `telemetry`, `time`, `transport`, `environment`

La mayoría de las sub-configuraciones son punteros opcionales. Por convención, `nil` significa **desactivado**.

---

## 🔐 Cadenas de origen (secretos, DSN, rutas)

Muchos campos aceptan una *cadena de origen* en lugar de solo una literal:

- `env:NAME` → leer desde la variable de entorno `NAME` (falla si `NAME` no está establecida; resuelve a un valor vacío si `NAME` está establecida explícitamente en `""`)
- `file:/ruta/a/archivo` → leer desde el sistema de archivos después de limpiar la ruta; los bytes devueltos se recortan de espacios en blanco al inicio y al final
- de lo contrario → tratar como cadena literal

Esto se usa para secretos y material de claves (claves TLS, claves HMAC, secretos de webhook, DSN de SQL, etc.).
Los valores `env:` y literales se devuelven exactamente como se proporcionaron; no se recortan.

Ejemplo:

```yaml
hooks:
  key: current
  secrets:
    current: env:WEBHOOK_SECRET
```

---

## 🌍 Entorno

El entorno de nivel superior es:

```yaml
environment: development
```

Este es un valor `env.Environment` utilizado para impulsar el comportamiento específico del entorno en los servicios.

---

## 🗜️ Compresión

Tipos de compresión utilizados por subsistemas que soportan compresión:

- `none`
- `zstd`
- `s2`
- `snappy`

---

## 🧾 Codificadores (Encoders)

Tipos de codificación utilizados por subsistemas que soportan codificación. `encoding.Map` registra cada codificador bajo exactamente un tipo canónico (sin alias):

- `json`
- `hjson`
- `toml`
- `yaml`
- `msgpack`
- `protobuf`
- `protojson`
- `prototext`
- `gob`
- `bytes`

> [!NOTE]
> - `bytes` es el codificador de paso a paso para cargas útiles `io.ReaderFrom`/`io.WriterTo`.
> - Los alias de tipo medio HTTP como `pb`, `proto`, `protobin`, `pbbin`, `pbtxt`, `prototxt`,
>   `pbjson`, `octet-stream`, `plain`, y `yml` se resuelven a los tipos canónicos anteriores mediante
>   `net/http/content/unary` antes de llegar a este registro. Consulte [Tipos de contenido HTTP](#http-content-types).
> - `encoding/stream.Map` es un registro separado para codificación de transmisión (multi-valor) — `json`, `msgpack`,
>   `gob`, `yaml` — utilizado por [Streaming HTTP (NDJSON)](#http-streaming-ndjson), no por este registro
>   de valor único.
> - No todos los tipos en este registro son intercambiables para decodificación de cuerpo de solicitud HTTP: `msgpack` y
>   `gob` permanecen como codecs de respuesta válidos pero son rechazados como `Content-Type` de solicitud. Consulte
>   [Tipos de contenido HTTP](#http-content-types).

---

## 💾 Caché

La configuración de caché se define en `cache/config.Config`:

```yaml
cache:
  kind: redis
  compressor: zstd
  encoder: json
  max_size: 4MB
  max_entries: 1024
  options:
    url: env:CACHE_URL
```

> [!NOTE]
> - Los tipos de controlador integrados en este repositorio son `redis` y `ttlcache`.
> - Los valores `kind` desconocidos devuelven `cache/driver/errors.ErrNotFound`.
> - Los valores `compressor` desconocidos o vacíos caen en `none`.
> - Para valores normales, los valores `encoder` desconocidos o vacíos caen en `json`.
> - Los valores `compressor` y `encoder` configurados son parte del espacio de nombres de claves del controlador de caché, por lo que cambiar cualquiera de estas configuraciones crea fallos de caché para valores escritos con el formato anterior.
> - Las operaciones de caché usan `bytes` para valores de flujo `io.WriterTo`/`io.ReaderFrom` y `protobuf` para mensajes protobuf, independientemente del `encoder` configurado.
> - `max_size` limita los valores de caché codificados antes de la compresión, después de la compresión y después de la descompresión. Un valor cero usa el predeterminado `4MB`.
> - `max_entries` limita las entradas retenidas por controladores de caché en memoria acotados. Un valor cero usa el predeterminado `1024`; los valores negativos son inválidos.
> - `options` es específico del backend y se decodifica como `map[string]any`.
> - `GetOrPersist` respaldado por Redis requiere [semánticas `SET` compatibles con Redis 7.0](https://redis.io/docs/latest/commands/set/) porque la publicación atómica usa `SET ... NX GET`.
> - Configure cada backend de caché para un servicio o propósito específico. Para Redis, use una base de datos dedicada, punto final o espacio de nombres de clave a nivel de implementación en la conexión/configuración en lugar de compartir una caché general para datos no relacionados.

> [!WARNING]
> `Cache.Flush` sigue las semánticas del backend; para Redis borra la base de datos seleccionada.

---

## 🚩 Banderas de características (OpenFeature)

`feature.Config` incorpora la configuración del lado del cliente (`config/client.Config`), por lo que soporta:

- `address`
- `timeout`
- `retry`
- `breaker`
- `limiter`
- `tls`
- `token`
- `options`

Ejemplo:

```yaml
feature:
  address: localhost:9000
  timeout: 10s
  breaker:
    max_requests: 2
    interval: 15s
    timeout: 5s
    consecutive_failures: 4
  retry:
    backoff: 100ms
    timeout: 1s
    attempts: 3
  tls:
    cert: file:test/certs/client-cert.pem
    key: file:test/certs/client-key.pem
    ca: file:test/certs/rootCA.pem
    server_name: localhost
```

> [!NOTE]
> - `feature.Config` incorpora la configuración del cliente; `IsEnabled` es verdadero solo cuando tanto la configuración de la característica como la configuración del cliente incorporada están presentes. Un bloque `feature:` vacío se trata como desactivado por los helpers de configuración de características.
> - Este repositorio no construye un proveedor OpenFeature integrado desde esta configuración.
> - Los servicios que necesitan un proveedor remoto o personalizado deben usar `feature.Config` en su propio constructor de proveedor y proporcionar el `openfeature.FeatureProvider` resultante en DI; `feature.Module` registra ese proveedor suministrado con el ciclo de vida del SDK OpenFeature.

---

## 🪝 Webhooks (Standard Webhooks)

Configurado vía `hooks.Config`:

```yaml
hooks:
  key: current
  secrets:
    current: env:WEBHOOK_SECRET_CURRENT
    previous: env:WEBHOOK_SECRET_PREVIOUS
  leeway: 30s
```

Cada valor `secrets` es una cadena de origen. El valor resuelto debe ser aceptado por la biblioteca Standard Webhooks, como un secreto generado por `hooks.Generator` con o sin el prefijo `whsec_`. Los secretos resueltos vacíos fallan al inicio.

La firma usa la `key` activa. La verificación acepta firmas de cada secreto configurado, probando primero el secreto activo. Standard Webhooks incluye un id de mensaje (`Webhook-Id`) pero no un id de clave de firma, por lo que go-service no extiende el protocolo con un encabezado selector personalizado.

`leeway` es una tolerancia opcional de desfase de reloj aplicada a la comprobación de frescura de `Webhook-Timestamp` durante la verificación. Un valor cero (el predeterminado) mantiene la ventana de frescura fija de 5 minutos de la biblioteca Standard Webhooks; un valor no nulo reemplaza esa ventana fija con esta tolerancia configurada, coincidiendo con el `leeway` de desfase de reloj ya expuesto por los verificadores de tokens JWT, PASETO y SSH. Como aquellos, `leeway` es una cadena de duración Go y debe ser una duración completa positiva.

La verificación entrante comprueba firmas y marcas de tiempo de Standard Webhooks, pero no almacena ni rechaza ids de webhook vistos anteriormente. Los receptores que realizan trabajo no idempotente deben desduplicar o procesar idempotentemente usando `Webhook-Id` o el id del evento, respaldado por almacenamiento compartido durable cuando se ejecutan más de una instancia de receptor.

> [!IMPORTANT]
> Los CloudEvents protegidos por webhook deben usar codificación HTTP estructurada. Los CloudEvents en modo binario con encabezados `ce-*` son rechazados antes de la verificación de firma.

---

## 🆔 Generación de IDs

Tipos de ID soportados:

- `uuid`
- `ksuid`
- `nanoid`
- `ulid`
- `xid`

Configuración:

```yaml
id:
  kind: uuid
```

> [!NOTE]
> Los generadores de ID producen identificadores operativos como ids de solicitud, ids de webhook y valores `jti` de token. No son una API de material secreto y no deben usarse como contraseñas, tokens de portador u otras credenciales. Omita `id` por completo para seleccionar el predeterminado `uuid`. Si `id` está presente, `kind` debe ser uno de los tipos registrados soportados. Los tipos ordenables como `ksuid`, `ulid` y `xid` exponen características de ordenación.

---

## 🚀 Mejoras de tiempo de ejecución

Los comandos de servidor creados a través de `cli.Application.AddServer` incluyen `runtime.Module`, que actualmente habilita:

- [automemlimit](https://github.com/KimMachineGun/automemlimit)

> [!NOTE]
> Este registro es de mejor esfuerzo y no falla al inicio si no se puede aplicar un límite de memoria. Las composiciones Fx directas y comandos estilo cliente deben incluir `runtime.Module` explícitamente cuando deseen este comportamiento.

---

## 🐘 SQL (Postgres)

La configuración raíz de SQL es `database/sql.Config`, con Postgres bajo `sql.pg`.

La configuración de Postgres incorpora configuración común de pool + DSN (`database/sql/config.Config`), incluyendo pools de escritor/lector. Cada pool de rol posee sus `dsns` y `settings`. Las configuraciones de pool SQL habilitadas deben establecer `max_open_conns` y `max_idle_conns` positivas; `max_idle_conns` no debe exceder `max_open_conns`.

Tanto `module.Server` como `module.Client` incluyen `sql.Module`, que actualmente conecta soporte PostgreSQL vía `database/sql/pg.Module`.

La habilitación es basada en presencia: un bloque `sql` nulo o un bloque `sql.pg` nulo deshabilita la conexión SQL. Cuando está habilitado, el controlador stdlib de pgx se registra bajo el nombre `pg`, y los DSN de lector/escritor se resuelven usando las reglas de cadena de origen descritas anteriormente. La configuración de PostgreSQL habilitada debe proporcionar al menos un `reader.dsns[].url` o `writer.dsns[].url` no vacío. La instrumentación del controlador se instala cuando el tracing o métricas están habilitados, las métricas estadísticas `database/sql` de OpenTelemetry se registran cuando las métricas están habilitadas, y los pools resultantes se cierran al detener el ciclo de vida.

La conexión SQL crea manejadores de pool `database/sql` y aplica configuraciones de pool, pero no hace ping a PostgreSQL durante la construcción. Llame a `DBs.Ping`, `DBs.PingWriter`, `DBs.PingReader`, o registre `health/checker.NewDBChecker` cuando el inicio o la disponibilidad deberían verificar la accesibilidad de la base de datos.

Ejemplo (con cadenas de origen para DSN):

```yaml
sql:
  pg:
    reader:
      dsns:
        - url: env:PG_READER_DSN
      settings:
        max_open_conns: 20
        max_idle_conns: 10
        conn_max_idle_time: 30m
        conn_max_lifetime: 1h
    writer:
      dsns:
        - url: env:PG_WRITER_DSN
      settings:
        max_open_conns: 3
        max_idle_conns: 2
        conn_max_idle_time: 10m
        conn_max_lifetime: 30m
```

Ejemplo (DSN literal; no recomendado para secretos de producción):

```yaml
sql:
  pg:
    writer:
      dsns:
        - url: postgres://user:pass@localhost:5432/dbname?sslmode=disable
      settings:
        max_open_conns: 10
        max_idle_conns: 5
```

### Dependencias

![Dependencies](./assets/database.png)

---

## 🩺 Estado de salud (Health)

Las verificaciones de salud se basan en [go-health](https://github.com/alexfalkowski/go-health).

El marco de trabajo proporciona endpoints estilo Kubernetes:

- `/<name>/healthz` — estado general de salud de servicio
- `/<name>/livez` — sonda de viabilidad (liveness)
- `/<name>/readyz` — sonda de disponibilidad (readiness)

Las respuestas de salud exitosas devuelven HTTP 200 con el cuerpo de texto plano `SERVING`.
Observadores faltantes o fallidos devuelven HTTP 503 con la respuesta de error estándar de go-service.
Durante el apagado del servidor, `/readyz` también devuelve HTTP 503 después de que el ciclo de vida comienza a drenar para que los orquestadores dejen de enviar tráfico nuevo antes de que el listener se detenga por completo.

Los helpers de verificador integrados bajo `health/checker` incluyen verificaciones de conectividad de DB y
verificaciones de conectividad de caché para controladores de caché pingables como Redis y ttlcache.

`module.Server` instala los transportes de salud HTTP/gRPC, pero los servicios poseen las
verificaciones y el mapeo de observadores. Cree valores `server.Registration` de go-health,
regístrelos bajo el nombre de servicio o gRPC service en `*server.Server`, y
mápelos a `healthz`, `livez`, `readyz`, o `grpc` con `Observe`. Consulte el
ejecutable [`Registrations` example](health/example_test.go) y el
[`go-service-template` health module](https://github.com/alexfalkowski/go-service-template/tree/master/internal/health)
para el patrón DI estándar. Un verificador no se expone por una sonda hasta que
exista ese registro y mapeo de observador.

Cuando el transporte gRPC está habilitado, `transport/grpc/health` registra el servicio
estándar `grpc.health.v1.Health` en el servidor gRPC. Las verificaciones nombradas usan el nombre de servicio
como la solicitud `service`; un servicio vacío verifica la salud general de gRPC:

```sh
grpcurl -plaintext -d '{"service":"<name>"}' localhost:9000 grpc.health.v1.Health/Check
```

`Check` devuelve `SERVING` o `NOT_SERVING` para servicios conocidos y `NotFound` para
servicios desconocidos. `List` devuelve los estados actuales para servicios registrados.
`Watch` transmite cambios de estado hasta que el cliente cancela; los servicios desconocidos transmiten
`SERVICE_UNKNOWN`. Los RPC de operación de salud omiten la verificación de token. El `Check` y `List` unarios también omiten el limitador del lado del servidor unario, mientras que `Watch` de salud
es un stream y aún usa limitación de stream.

Estos están modelados después de [endpoints de salud de la API de Kubernetes](https://kubernetes.io/docs/reference/using-api/health-checks/).

---

## 📡 Telemetría

La raíz de configuración de telemetría es `telemetry.Config`:

```yaml
telemetry:
  attributes:
    k8s.namespace.name: payments
  logger: ...
  metrics: ...
  propagation: ...
  tracer: ...
```

`attributes` son etiquetas de recurso OpenTelemetry planas adjuntas a registros, métricas,
y trazas. No son cadenas de origen. Los atributos de identidad fijos de go-service
como `host.id`, `service.instance.id`, `service.name`, `service.version`,
y `deployment.environment.name` tienen precedencia si se configura la misma clave.

### Propagación

La propagación de contexto OpenTelemetry por defecto es W3C Trace Context más W3C Baggage
para extracción e inyección:

```yaml
telemetry:
  propagation:
    formats:
      - tracecontext
      - baggage
```

Los ecosistemas de tracing mixtos pueden habilitar formatos adicionales:

```yaml
telemetry:
  propagation:
    formats:
      - tracecontext
      - baggage
      - b3
```

Los propagadores soportados son `tracecontext`, `baggage`, `b3`, `b3multi`, y
`none`. Use `none` solo como el único valor para `formats`.

B3 usa el propagador B3 upstream, que soporta tanto formatos de encabezado único como
multi-encabezado.

### Registros (Logging)

El logging usa `log/slog`.

Tipos de logger integrados soportados:

- `json`
- `text`
- `tint`
- `otlp`

Los niveles de logger soportados son `debug`, `info`, `warn`, y `error`. Cuando `level`
no está establecido, el logging por defecto es `info`; los valores desconocidos fallan en la construcción del logger.

#### Registro JSON

```yaml
telemetry:
  logger:
    kind: json
    level: info
```

#### Registro de texto

```yaml
telemetry:
  logger:
    kind: text
    level: info
```

#### Registro OTLP

```yaml
telemetry:
  logger:
    kind: otlp
    level: info
    protocol: http
    url: http://localhost:4318/v1/logs
    batch_timeout: 5s
    export_timeout: 30s
    max_queue_size: 2048
    max_export_batch_size: 512
    headers:
      Authorization: env:OTLP_LOGS_AUTH
```

> [!NOTE]
> - `batch_timeout`, `export_timeout`, `max_queue_size`, y `max_export_batch_size` ajustan la canalización de exportación por lotes OTLP y se aplican solo cuando `kind` es `otlp`. Cuando un valor no está establecido o es cero, se usa el predeterminado del SDK OpenTelemetry (cola `2048`, lote `512`).
> - Los valores `headers` son cadenas de origen.
> - Los mapas de encabezados de telemetría se resuelven durante la proyección de configuración; los valores `env:` no establecidos y los valores `file:` no legibles fallan rápidamente (pánico durante el inicio).

> [!WARNING]
> Los exportadores OTLP rechazan endpoints `http://` no de bucle local cuando se configuran encabezados. Use HTTPS para colectores remotos que requieran encabezados de autorización; el texto claro con encabezados se acepta solo para endpoints de bucle local.
>
> Los exportadores OTLP/gRPC usan `protocol: grpc` y un endpoint `host:port` como `localhost:4317`. Los endpoints remotos gRPC con encabezados requieren la configuración `tls` de la señal; los endpoints gRPC de bucle local aún pueden usar texto claro.
>
> Los endpoints de exportador OTLP deben establecerse en campos de configuración de go-service como `telemetry.logger.url`, `telemetry.metrics.url`, y `telemetry.tracer.url`. Las variables de entorno de endpoint OpenTelemetry estándar como `OTEL_EXPORTER_OTLP_ENDPOINT` no se usan como fuentes de respaldo.

Los exportadores remotos OTLP/gRPC pueden usar el mismo modelo de cadena de origen TLS que otros clientes de go-service:

```yaml
telemetry:
  tracer:
    kind: otlp
    protocol: grpc
    url: collector.example.com:4317
    tls:
      ca: file:/etc/otel/ca.pem
      cert: file:/etc/otel/client.crt
      key: file:/etc/otel/client.key
      server_name: collector.example.com
    headers:
      Authorization: env:OTLP_TRACES_AUTH
```

Use la misma forma `tls` bajo `telemetry.logger` o `telemetry.metrics` cuando esas señales exportan a través de OTLP/gRPC.

### Métricas

Tipos de métricas soportados:

- `prometheus`
- `otlp`

#### Prometheus

```yaml
telemetry:
  metrics:
    kind: prometheus
    prometheus:
      without_suffixes: true
      without_target_info: true
      without_scope_info: true
```

Cuando Prometheus está habilitado en el transporte HTTP, las métricas se exponen en `/<name>/metrics`.

El bloque opcional `prometheus` moldea la salida del exportador para compatibilidad con una
pila existente de Prometheus/Grafana/alertas. `without_suffixes` elimina sufijos de unidad
(por ejemplo `_seconds`, `_bytes`) y sufijos de contador `_total` de los nombres de métrica, `without_target_info` omite la métrica `target_info`, y
`without_scope_info` omite las etiquetas `otel_scope_name`/`otel_scope_version`. Cuando
se omite el bloque `prometheus`, el exportador mantiene su salida
convencional de OpenTelemetry por defecto.

#### Métricas OTLP

```yaml
telemetry:
  metrics:
    kind: otlp
    protocol: http
    url: http://localhost:9009/otlp/v1/metrics
    interval: 30s
    timeout: 5s
    headers:
      Authorization: env:OTLP_METRICS_AUTH
```

`interval` y `timeout` se aplican solo a métricas de empuje OTLP. Cuando cualquiera de los valores
no está establecido o es cero, se usa el predeterminado del SDK OpenTelemetry.

#### Buckets de histograma

Anule los límites de bucket de histograma predeterminados por instrumento con
`telemetry.metrics.views`, indexados por nombre de instrumento (coincidencia de nombre OpenTelemetry,
incluyendo comodines `*`):

```yaml
telemetry:
  metrics:
    views:
      http.server.request.duration: [0.005, 0.01, 0.05, 0.1, 0.5, 1, 5]
      "rpc.*.duration": [0.01, 0.1, 1]
```

Los límites están en la unidad del instrumento (segundos para histogramas de duración, bytes
para histogramas de tamaño) y deben listarse en orden ascendente. Las vistas se aplican a
instrumentos de histograma independientemente del tipo de métricas; un mapa no establecido o vacío mantiene los
buckets predeterminados del SDK OpenTelemetry.

### Seguimiento (Tracing)

El tracing soporta configuración de exportador OTLP:

```yaml
telemetry:
  tracer:
    kind: otlp
    protocol: http
    url: http://localhost:4318/v1/traces
    batch_timeout: 5s
    export_timeout: 30s
    max_queue_size: 2048
    max_export_batch_size: 512
    sampler:
      kind: ratio
      ratio: 0.25
    headers:
      Authorization: env:OTLP_TRACES_AUTH
```

> [!NOTE]
> `batch_timeout`, `export_timeout`, `max_queue_size`, y `max_export_batch_size` ajustan la canalización de exportación de spans por lotes OTLP. Cuando un valor no está establecido o es cero, se usa el predeterminado del SDK OpenTelemetry (cola `2048`, lote `512`).
>
> Los exportadores OTLP por defecto usan `protocol: http`. Establezca `protocol: grpc` y use una
> `url` `host:port`, como `localhost:4317`, para exportar a través de OTLP/gRPC.
>
> Tipos de sampler soportados:
>
> - `always_on`: registra cada traza.
> - `always_off`: descarta cada traza.
> - `ratio`: sigue la decisión de muestreo del span padre entrante cuando la solicitud
>   ya tiene contexto de traza; de lo contrario, registra la fracción configurada de nuevas
>   trazas raíz. Establezca `ratio` entre `0` y `1`, donde `0` descarta nuevas trazas
>   raíz y `1` registra todas las nuevas trazas raíz.
>
> Cuando se omite `sampler`, go-service preserva el
> sampler predeterminado del SDK OpenTelemetry y el manejo de entorno del sampler del SDK.

### Bibliotecas de telemetría utilizadas

- <https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/runtime>
- <https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/host>
- <https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp>
- <https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc>
- <https://github.com/redis/go-redis/tree/master/extra/redisotel>
- <https://github.com/XSAM/otelsql>

### Dependencias de Telemetría

![Dependencies](./assets/telemetry.png)

---

## 🎫 Tokens

La configuración de tokens está arraigada en `token.Config`, usualmente anidada bajo la configuración de transporte como `transport.http.token` y/o `transport.grpc.token` (vía la configuración compartida del lado del servidor de transporte).

Tipos de token `kind` soportados:

- `jwt`
- `paseto`
- `ssh`

### Control de acceso (Casbin)

El control de acceso se configura una vez a nivel de transporte y se comparte entre todos
los stacks de servidor HTTP y gRPC habilitados:

```yaml
transport:
  access:
    model: file:./config/rbac.conf
    policy: file:./config/rbac.csv
```

Cuando `access` está configurado, los stacks de servidor HTTP y gRPC estándar aplican
la política después de la autenticación de token y antes de que se ejecuten los manejadores de la aplicación. Omita
`access` para dejar la autorización de transporte deshabilitada.

El modelo se basa en RBAC de Casbin:
<https://github.com/casbin/casbin/blob/master/examples/rbac_model.conf>

Ejemplo `rbac.conf`:

```ini
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
```

Las políticas usan el id de usuario verificado como `sub`, `meta.TransportServiceMethod` como
`obj`, y `invoke` como `act`. Ejemplo `rbac.csv`:

```csv
p, reader, http:GET /users/{id}, invoke
p, writer, http:POST /users, invoke
p, greeter, grpc:/greet.v1.GreeterService/SayHello, invoke
g, frontend, reader
g, admin, reader
g, admin, writer
g, billing-service, greeter
```

Las filas `p` definen permisos y deben coincidir con la forma `p = sub, obj, act` del modelo, por lo que incluyen `invoke`. Las filas `g` definen membresía de rol y coinciden con `g = _, _`, por lo que solo contienen `subject, role`.

> [!WARNING]
> El adaptador de política de cadena de Casbin puede saltar filas de política malformadas sin fallar
> al inicio. Valide los archivos de política antes del despliegue; un inicio exitoso no
> prueba que cada fila configurada fue cargada.

Para servidores HTTP, el objeto usa el patrón de ruta coincidente cuando está disponible, como
`http:GET /users/{id}`. Los tokens HTTP se autentican contra el método y path de solicitud concretos, como `GET /users/123`; la aplicación de política de acceso
usa el patrón de ruta canónico. Los tokens gRPC se autentican contra el nombre completo
del método, como `/greet.v1.GreeterService/SayHello`; la aplicación de política de acceso
usa el objeto de método de servicio de transporte, como
`grpc:/greet.v1.GreeterService/SayHello`.

> [!NOTE]
> `access.model` y `access.policy` se resuelven a través de `os.FS.ReadSource`; use `file:` para archivos, `env:` para contenido proporcionado por el entorno, o contenido literal.
>
> La configuración de acceso construye un controller inyectable para verificaciones de autorización. Los stacks de servidor HTTP y gRPC integrados autentican tokens, almacenan el id de usuario verificado y aplican la política de Casbin configurada antes de que se ejecuten los manejadores de la aplicación.

### JWT

Configuración JWT:

```yaml
transport:
  http:
    token:
      kind: jwt
      jwt:
        iss: my-service
        exp: 1h
        leeway: 30s
        key: active
        keys:
          active:
            public: file:/keys/ed25519.pub
            private: file:/keys/ed25519
          old:
            public: file:/keys/ed25519-old.pub
```

Comportamiento importante:

- La generación JWT firma con `jwt.key`; la verificación requiere que el encabezado `kid` del token seleccione una entrada en `jwt.keys`.
- `exp` se parsea como una cadena de duración Go; los valores inválidos pueden fallar rápido.
- `leeway` es una tolerancia opcional de desfase de reloj para verificación; manténgala pequeña porque extiende la aceptación alrededor de `iat`/`nbf` y `exp`.

> [!IMPORTANT]
> La generación y verificación JWT usan material de clave Ed25519 de `jwt.keys`. Mantenga el material de clave privada solo en servicios que acuñan tokens; los verificadores solo necesitan claves públicas.

Todos los valores `exp` de token y `leeway` no nulos son cadenas de duración Go y deben ser duraciones completas positivas. Valores como
`1s`, `15m`, y `24h` validan; los valores sub-segundo como `500ms` no.

### Paseto

Configuración Paseto:

```yaml
transport:
  http:
    token:
      kind: paseto
      paseto:
        iss: my-service
        exp: 1h
        leeway: 30s
        key: active
        keys:
          active:
            public: file:/keys/ed25519.pub
            private: file:/keys/ed25519
          old:
            public: file:/keys/ed25519-old.pub
```

> [!NOTE]
> La implementación de PASETO emite tokens **v4 public**. La generación firma con `paseto.key`, escribe ese id como pie de página `kid`, y la verificación selecciona la clave pública desde `paseto.keys`. `paseto.leeway` es una tolerancia opcional de desfase de reloj para verificación.

### Tokens SSH

Las claves de verificación de token SSH son direccionables por id y soportan rotación.

Ejemplo solo verificación:

```yaml
transport:
  http:
    token:
      kind: ssh
      ssh:
        exp: 5m
        leeway: 30s
        keys:
          active:
            public: file:/keys/active.pub
```

Ejemplo firma + verificación:

```yaml
transport:
  http:
    token:
      kind: ssh
      ssh:
        exp: 5m
        leeway: 30s
        key: active
        keys:
          active:
            public: file:/keys/active.pub
            private: file:/keys/active
          old:
            public: file:/keys/old.pub
```

> [!NOTE]
> - `ssh.key` es el id de clave activa usado para acuñar tokens (la entrada `ssh.keys` correspondiente requiere material de clave privada).
> - `ssh.keys` es el mapa de claves de confianza usado para verificación (claves públicas).
> - `ssh.exp` establece la ventana de validez del token; las claves SSH permanecen de larga vida, mientras que los tokens generados son de corta vida.
> - `ssh.leeway` es una tolerancia opcional de desfase de reloj para verificación; manténgala pequeña porque extiende la aceptación alrededor de `iat` y `exp`.
> - Los tokens SSH llevan `sub` igual a `kid`, por lo que el sujeto verificado es el id de clave de par de confianza.

---

## 🚦 Limitador (Limiter)

La configuración del limitador es `transport/limiter.Config` y típicamente se aplica a nivel de transporte.

Tipos de clave soportados (integrados):

- `user-id`
- `transport-service-method`
- `service-method`
- `ip`
- `user-agent`

Ejemplo:

```yaml
transport:
  http:
    limiter:
      kind: user-agent
      tokens: 10
      interval: 1s
      max_keys: 4096
```

> [!NOTE]
> - `interval` se parsea como una cadena de duración Go. Los valores inválidos pueden fallar rápido.
> - `tokens` e `interval` usan los defaults de la tienda en memoria subyacente cuando se establecen en cero: `1` token por `1s`. Configure valores positivos para cuotas explícitas.
> - `max_keys` limita el número de claves derivadas del caller que reciben buckets en memoria independientes. Un valor cero usa el predeterminado `4096`; las claves distintas adicionales comparten un bucket de desbordamiento.
> - El limitador integrado es una salvaguarda en memoria, por proceso. Úselo como último recurso y prefiera un limitador de borde externo, gateway, ingress, balanceador de carga, o service-mesh para protección contra abuso en producción.
> - La clave `user-id` usa el principal verificado almacenado en metadatos. Para tokens JWT/PASETO esto es el claim de subject; para tokens SSH esto es el nombre de clave verificado. Prefiéralo cuando la identidad autenticada esté disponible.
> - La clave `transport-service-method` prefija el valor de método de servicio con el nombre del transporte, como `http:GET /users/{id}` o `grpc:/users.v1.Users/Get`, por lo que las operaciones HTTP y gRPC usan buckets separados.
> - La clave `service-method` usa metadatos de ruta/path HTTP o el nombre completo del método gRPC. Prefiera `transport-service-method` a menos que las operaciones entre transportes compartan intencionalmente cuota.
> - Los limitadores HTTP y gRPC del lado del servidor se ejecutan después de la extracción de metadatos y verificación de token, por lo que la autorización faltante, malformada o inválida se rechaza antes de llegar al limitador. Esto es intencional; aplique cuotas para esos intentos con un limitador de borde externo, gateway, ingress, balanceador de carga, o service-mesh.
> - Los limitadores HTTP del lado del servidor establecen encabezados `RateLimit` y `RateLimit-Policy`; las solicitudes HTTP denegadas también establecen `Retry-After` cuando el tiempo de reinicio está disponible. Los limitadores gRPC del lado del servidor establecen metadatos de respuesta `ratelimit` y `ratelimit-policy`; las solicitudes gRPC denegadas también adjuntan un detalle `google.rpc.RetryInfo` cuando el tiempo de reinicio está disponible.
> - Los limitadores de stream gRPC consumen un token cuando el stream se abre y un token para cada operación `RecvMsg` y `SendMsg`. Las solicitudes unarias HTTP y gRPC consumen un token por solicitud/RPC.

---

## 🕒 Tiempo (tiempo de red)

Configuración de tiempo:

```yaml
time:
  kind: nts
  address: time.cloudflare.com
  timeout: 2s
```

Tipos soportados:

- `ntp`
- `nts`

Omita el bloque `time` para deshabilitar el tiempo de red. Si el bloque está presente, `kind`
debe ser `ntp` o `nts`; los tipos vacíos o desconocidos fallan al inicio con el error de proveedor de tiempo no encontrado. `address` es específico del proveedor y se usa cuando el
proveedor de tiempo de red realiza E/S. `timeout` limita las operaciones de red para el
proveedor seleccionado; un valor cero usa el timeout predeterminado del cliente upstream, y
los valores negativos son inválidos.

---

## 🌐 Transporte

La capa de transporte proporciona conexión de nivel superior y política de middleware para comunicación de entrada/salida del servicio.

A un nivel alto:

- `transport/...` contiene la capa de transporte de servicio con enfoques propios: conexión Fx, pilas de servidor y cliente HTTP/gRPC compuestas, reintentos, circuit breakers, middleware de token, conexión de salud, y política relacionada.
- `net/...` contiene helpers de protocolo a nivel inferior y primitivas reutilizables como `net/http`, `net/grpc`, `net/http/meta`, `net/grpc/meta`, `net/grpc/health`, `net/header`, y `net/server`.

Los stacks soportados incluyen:

- gRPC (<https://grpc.io/>)
- Abstracción REST HTTP (`net/http/rest`) usando negociación de contenido
- Abstracción RPC HTTP (`net/http/rpc`) usando negociación de contenido
- Helpers MVC HTTP (`net/http/mvc`)
- CloudEvents (<https://github.com/cloudevents/sdk-go>)

La conexión HTTP de CloudEvents vive bajo `transport/http/events`: use
`NewReceiver(...).Register(...)` para recibir eventos en una ruta POST y
`NewSender(...).Send(...)` con `net/http/events.ContextWithTarget(...)` para
enviar eventos. El sender usa codificación HTTP estructurada por defecto; configure
`WithSenderEncoding(SenderEncodingBinary)` para integraciones salientes que
requieran CloudEvents en modo binario. Los receptores protegidos por webhook requieren codificación
estructurada y rechazan CloudEvents en modo binario con encabezados `ce-*` antes
de la verificación de firma. El registro de receptor marca la ruta del evento como
no autenticada para el middleware de token/acceso de transporte para que la verificación de webhook pueda
actuar como el límite de autenticación del evento.

### Tipos de contenido HTTP

Los helpers REST y RPC HTTP decodifican cuerpos de solicitud desde el `Content-Type` de la solicitud, cayendo en JSON cuando `Content-Type` está ausente. Un `Content-Type` sin parsear, no registrado o intencionalmente indecodificable se rechaza con HTTP 415 en lugar de caer en JSON. La codificación de respuesta usa el primer tipo medio `Accept` cuando está presente, cayendo en el `Content-Type` de la solicitud cuando `Accept` está ausente. Los helpers de cliente pueden establecer `ContentType` para el cuerpo de solicitud y `Accept` para un formato de respuesta independiente.

Los tipos medios de carga útil texto/objeto integrados incluyen:

- `application/json`
- `application/hjson`
- `application/yaml`
- `application/toml`
- `application/octet-stream`, `text/plain`

Los tipos medios de carga útil binaria interna incluyen:

- `application/vnd.msgpack`
- `application/gob`

Los alias de tipo medio orientados a protobuf integrados incluyen:

- `application/proto`, `application/pb`, `application/protobuf`, `application/protobin`, `application/pbbin`
- `application/protojson`, `application/pbjson`
- `application/prototext`, `application/prototxt`, `application/pbtxt`

> [!NOTE]
> - `application/hjson` mapea al tipo de codificador integrado `hjson`.
> - Los tipos medios desconocidos o inválidos caen en selección de JSON solo para salida (impulsada por `Accept`)
>   negociación. Un `Content-Type` de solicitud ausente aún predetermina a JSON, pero uno desconocido o inválido
>   se rechaza con HTTP 415 en lugar de decodificarse como un formato diferente al declarado por el caller.
> - `text/error` está reservado para respuestas de error y no debería enviarse por clientes como tipo de contenido de solicitud.
>
> `application/toml`, `application/vnd.msgpack`, y `application/gob` pueden resolverse como tipos medios y permanecen válidos
> como codecs de respuesta, pero la decodificación de cuerpo de solicitud REST/RPC — tanto para valores únicos como para streaming
> (NDJSON) — los rechaza con HTTP 415. Esto sigue la regla de límites de decodificador documentada en
> la documentación del paquete `net/http/content/unary`: un codec es admisible para decodificar entrada no confiable solo
> cuando está acotado por ratio y acotado por profundidad, lo que TOML, msgpack y gob no son.

### Streaming HTTP (NDJSON)

REST y RPC soportan rutas de streaming junto a los helpers de valor único anteriores, para respuestas (y,
sobre HTTP/2, solicitudes) que llegan como una secuencia de valores en lugar de una sola carga útil amortiguada:

| valor único | streaming | dirección | HTTP/2 requerido |
| --- | --- | --- | --- |
| `rest.Get`/`rest.Route` | `rest.StreamGet`/`rest.StreamRoute` | solo envío | no |
| `rest.Post`/`rest.Put`/`rest.Patch`/`rest.RouteRequest` | `rest.StreamPost`/`rest.StreamPut`/`rest.StreamPatch`/`rest.StreamRouteRequest` | bidireccional | sí |
| `rpc.Route` | `rpc.StreamRoute` | bidireccional | sí |

Un manejador de streaming solo de envío obtiene un `*stream.Stream[Res]` con `Send`; un manejador de streaming
bidireccional obtiene un `*stream.RequestStream[Req, Res]` con ambos `Send` y `Recv`. Las llamadas de cliente usan las
funciones coincidentes `client.Stream`/`client.RequestStream`, que toman el mismo tipo de callback.
Vea `ExampleClient_RequestStream` de `net/http/client` para una llamada completa de cliente bidireccional HTTP/2.

> [!IMPORTANT]
> Los helpers de valor único viven en `github.com/alexfalkowski/go-service/v2/net/http/content/unary`; los helpers de streaming
> viven en `github.com/alexfalkowski/go-service/v2/net/http/content/stream`. Importe `unary` para `Content`, `Media`,
> `NewContent`, `NewHandler`, y `NewRequestHandler`; importe `stream` para helpers de solicitud/respuesta incrementales.
> `stream.NewHandler` y `stream.NewRequestHandler` toman `*stream.Content`, mientras que los manejadores unarios toman
> `*unary.Content`. `rest.Register`, `rpc.Register`, y `client.NewClient` toman los propietarios de contenido unario y de streaming por separado.
> Migre las importaciones e identificadores raíz `content` a `content/unary` y `unary`, respectivamente; use
> `net/http.ContentTypeKey` y `net/http.AcceptKey` para los nombres de encabezado compartidos.

El formato de cable inicial es NDJSON (`application/x-ndjson`), valores JSON delimitados por nueva línea, resueltos
a través de un registro de codificador/decodificador de streaming separado (`encoding/stream.Map`) del de valor único
anterior — un tipo medio de streaming no registrado o sin parsear se rechaza directamente en lugar de caer
en JSON, a diferencia de la negociación de valor único.

> [!NOTE]
> - Las rutas de streaming bidireccionales requieren HTTP/2 (incluyendo h2c); una solicitud sobre HTTP/1.x se rechaza
>   con `505 HTTP Version Not Supported` antes de que se ejecute el manejador. Las rutas de streaming solo de envío no tienen
>   tal requisito y permanecen completamente soportadas en respuestas chunked HTTP/1.1.
> - Las respuestas de streaming no se comprimen con gzip, independientemente del `Accept-Encoding` del cliente.
> - `max_receive_size` se aplica por valor decodificado en un cuerpo de solicitud de streaming, no como un total acumulativo
>   a través de todo el stream; el volumen total del stream se controla mediante el limitador de tasa configurado,
>   que cobra un token por cada mensaje transmitido además del token cobrado cuando el stream se abre.
> - Un `Send` exitoso extiende el timeout de escritura configurado del servidor HTTP, y en un stream bidireccional
>   un `Recv` exitoso extiende ambos los timeouts de lectura y escritura (y `Send` también extiende ambos),
>   por lo que un stream lento pero activo no se corta por un límite de tiempo de todo el stream en ninguna dirección; limite
>   una llamada de streaming del lado del cliente con el contexto de solicitud en lugar del timeout de solicitud general del cliente.
> - Los timeouts de lectura/escritura por mensaje siguen la misma precedencia `options.read_timeout`/`options.write_timeout`
>   que los propios timeouts del servidor (vea [Configuración de transporte (servidores)](#transport-configuration-servers)),
>   cayendo en `timeout` cuando la opción correspondiente no está establecida.
> - Las solicitudes de streaming nunca se reintentan por el middleware de reintento del cliente.
> - Un fallo de stream después de que la respuesta se haya comprometido se registra como un error de traza y en el log de acceso, luego
>   aborta la respuesta para que los clientes no reciban un stream limpio pero truncado. Las métricas RED del servidor HTTP upstream
>   no registran streams abortados; use el log de acceso para investigar esa clase de fallo.
> - Durante el apagado estándar del servidor, los contextos del manejador de stream se cancelan. Los manejadores deben retornar después
>   de `ctx.Done()` cuando esperan en una fuente upstream; un `Recv` activo termina con la señal de drenaje. Un `Send` bloqueado
>   permanece sujeto al timeout de escritura configurado. Si expira el límite de apagado del ciclo de vida, el
>   servidor cierra por fuerza las conexiones HTTP restantes, por lo que los clientes observan un error de transporte. Un cliente HTTP/2 bidireccional
>   puede observar el cierre por fuerza del cuerpo de solicitud como un reinicio de stream y debería reconectarse a un
>   servidor no drenante.

### Rutas HTTP no encontradas

El transporte HTTP envuelve el mux con `net/http.NewNotFoundHandler` para que las respuestas 404 generadas puedan renderizarse consistentemente mientras se preservan otras respuestas del mux como 405 Method Not Allowed.

- Las rutas faltantes estilo REST/RPC usan `net/http/status.NotFoundHandler`, que escribe la respuesta estándar `status.WriteError`.
- Las rutas faltantes MVC pueden usar `net/http/mvc.NotFoundHandler` para renderizar la vista no encontrada MVC registrada cuando la solicitud acepta HTML (`Accept: text/html`) o es una solicitud HTMX (`Hx-Request: true`).
- Las rutas que coinciden y escriben su propio estado no son reemplazadas por este manejador de no encontrado a nivel de mux.

### Errores HTTP MVC

Cuando un controller MVC devuelve un error, `net/http/mvc.Route` renderiza la vista devuelta con un modelo `mvc.Error` seguro para el cliente. El modelo contiene el `Code` de estado HTTP y un `Message` visible seguro para el cliente.

La cadena de error cruda permanece disponible para templates como metadato `mvcModelError` por compatibilidad. Renderizar ese metadato puede exponer detalles de diagnóstico, por lo que prefiera `.Model.Message` para páginas de error visibles al cliente.

### Configuración de transporte (servidores)

La raíz de configuración de transporte es `transport.Config`:

- `transport.http` incorpora `config/server.Config`
- `transport.grpc` incorpora `config/server.Config`

Ejemplo mínimo:

```yaml
transport:
  http:
    address: tcp://localhost:8000
    timeout: 10s
  grpc:
    address: tcp://localhost:9000
    timeout: 10s
```

> [!NOTE]
> - La dirección puede usar `<network>://<address>` (por ejemplo `tcp://:8000`) o una dirección de escucha cruda como `:8000`, que predetermina a la red `tcp`.
> - Si la dirección se omite, los predeterminados son `tcp://:8080` (HTTP) y `tcp://:9090` (gRPC).
> - `transport.grpc.timeout` limita los manejadores RPC unarios y alimenta los predeterminados de keepalive/conexión del servidor gRPC; no limita la vida útil del stream. Los streams de larga vida permanecen abiertos hasta la cancelación del cliente o se apliquen controles específicos de stream.
> - `max_receive_size` limita el tamaño de carga útil entrante. Un valor cero usa el predeterminado `4MB`.
> - Para HTTP, `max_receive_size` se aplica por cuerpo de solicitud, excepto para rutas de streaming bidireccionales (vea [Streaming HTTP (NDJSON)](#http-streaming-ndjson)), donde se aplica por valor decodificado en su lugar, sin un total acumulativo. Para gRPC, se aplica por solicitud unaria entrante y por mensaje de stream entrante.
> - MVC no aplica sus propios límites de tamaño de cuerpo; la conexión de servidor HTTP soportada aplica `max_receive_size` antes de que se ejecuten los manejadores MVC, y los clientes HTTP de go-service aplican su límite de tamaño de respuesta configurado al leer respuestas.

Ejemplo de límite de recepción:

```yaml
transport:
  http:
    max_receive_size: 2MB
  grpc:
    max_receive_size: 3MB
```

Con opciones de servidor a nivel bajo:

```yaml
transport:
  http:
    address: tcp://localhost:8000
    timeout: 10s
    options:
      read_timeout: 10s
      write_timeout: 10s
      idle_timeout: 10s
      read_header_timeout: 10s
  grpc:
    address: tcp://localhost:9000
    timeout: 10s
    options:
      keepalive_enforcement_policy_ping_min_time: 10s
      keepalive_max_connection_idle: 10s
      keepalive_max_connection_age: 10s
      keepalive_max_connection_age_grace: 10s
      keepalive_ping_time: 10s
```

### TLS para transportes

La configuración TLS usa `crypto/tls/config.Config` y los campos son cadenas de origen:

```yaml
transport:
  http:
    tls:
      cert: file:test/certs/cert.pem
      key: file:test/certs/key.pem
      ca: file:test/certs/rootCA.pem
  grpc:
    tls:
      cert: file:test/certs/cert.pem
      key: file:test/certs/key.pem
      ca: file:test/certs/rootCA.pem
```

Establezca `ca` en la configuración TLS del servidor para requerir y verificar certificados de cliente para mTLS. Establezca `ca` en la configuración TLS del
cliente para verificar certificados de servidor emitidos por la misma CA local o privada. `server_name` solo es necesario
en clientes cuando la dirección de dial difiere del nombre DNS del certificado.

El TLS del lado del servidor requiere un par completo `cert` y `key` siempre que se configure material TLS. `ca` habilita
la verificación de certificado de cliente para mTLS, pero una configuración TLS de servidor solo con CA falla al inicio.

Los servidores en tiempo de ejecución requieren TLS 1.3 o posterior en handshakes entrantes; los clientes mantienen un piso TLS 1.2 para que las llamadas salientes permanezcan interoperables con endpoints solo TLS-1.2.

Los clientes gRPC usan credenciales de transporte inseguras cuando TLS no está configurado. Ese predeterminado está destinado para
tráfico local o protegido por plataforma; configure TLS del cliente para llamadas fuera de ese límite de confianza.

> [!IMPORTANT]
> Si está usando `go-service-template` o componiendo paquetes de transporte de servidor como `module.Server` o `transport.Module`, el registro de transporte requerido se maneja por usted mediante DI.
>
> `module.Client` no conecta transportes por defecto. Cuando un proceso de cliente construye configuración TLS HTTP o gRPC desde cadenas de origen como `file:`, llame a las funciones `Register(...)` relevantes a nivel de transporte, como `transport/http.Register(...)` o `transport/grpc.Register(...)`.
>
> Solo necesita llamar a las funciones `Register(...)` a nivel de transporte usted mismo cuando intencionalmente conecta transportes manualmente o compone paquetes a nivel inferior fuera del gráfico de módulos de transporte.
>
> Si está conectando el ciclo de vida del servidor manualmente, use `net/server.Register(...)`.

### IPs reenviadas y reflexión

> [!WARNING]
> La extracción de metadatos HTTP y gRPC intencionalmente confía en encabezados/metadatos de IP reenviada comunes como `X-Forwarded-For`, `X-Real-IP`, `CF-Connecting-IP`, y `True-Client-IP`. Los servicios que dependen de IPs extraídas para logging, política o limitación de tasa solo deberían recibir tráfico a través de infraestructura de borde confiable que elimine o sobrescriba los encabezados de reenvío suministrados por el cliente.

> [!WARNING]
> La reflexión del servidor gRPC está intencionalmente siempre registrada por `net/grpc.NewServer` para que las herramientas internas puedan descubrir servicios. Los servicios que no deberían exponer reflexión públicamente deberían restringir el acceso con direcciones de unión, TLS/autenticación de cliente, política de ingress, reglas de firewall, o autorización de service-mesh.

### Dependencias de Transporte

![Dependencies](./assets/transport.png)

### Circuit breakers (lado del cliente)

Los wrappers de cliente de transporte incluyen circuit breakers opcionales:

- Breaker HTTP (`transport/http/breaker`):
  - El alcance es por `"<METHOD> <HOST>"`.
  - Los estados de fallo predeterminados son `>=500` y `429`.
  - Los errores de transporte se cuentan como fallos.
  - Las respuestas de estado de fallo aún se devuelven a los llamadores (mientras la contabilidad del breaker registra un fallo).

- Breaker gRPC (`transport/grpc/breaker`):
  - El alcance es por `fullMethod`.
  - Los códigos de fallo predeterminados son `Unavailable`, `DeadlineExceeded`, `ResourceExhausted`, e `Internal`.
  - Los errores con otros códigos gRPC se tratan como exitosos para la contabilidad del breaker.

La configuración del cliente usa la forma compartida `transport/breaker.Config` para la mecánica del breaker. Cualquier tipo de configuración que
incorpore `config/client.Config` tiene su propio bloque `breaker` bajo esa configuración de cliente. Este ejemplo usa
`feature.Config` solo porque es una configuración de cliente de ese tipo:

```yaml
feature:
  address: localhost:9000
  breaker:
    max_requests: 2
    interval: 15s
    timeout: 5s
    consecutive_failures: 4
```

Al construir manualmente clientes HTTP o gRPC, pase una configuración de breaker específica de transporte a
`transport/http.WithClientBreaker(...)` o `transport/grpc.WithClientBreaker(...)`. Estas configuraciones
incorporan la mecánica compartida del breaker y agregan clasificación de fallo específica de protocolo:

```go
httpBreaker := httpbreaker.NewConfig(sharedBreaker, 429, 502, 503)
grpcBreaker := grpcbreaker.NewConfig(sharedBreaker, codes.Unavailable, codes.ResourceExhausted)
```

`NewConfig` devuelve `nil` cuando la configuración compartida del breaker es `nil`, preservando la conexión de opciones de cliente que
desactiva breakers omitiendo la configuración del breaker.

`max_requests` controla la concurrencia de probes semi-abierto. `interval` controla la
ventana de reinicio de conteo en estado cerrado. `timeout` controla cuánto tiempo permanece el breaker
abierto antes de permitir probes semi-abierto. `consecutive_failures` controla cuándo se
abre el breaker. Los valores cero mantienen los predeterminados del paquete.

En lugar de (o junto con) `consecutive_failures`, `failure_ratio` y
`min_requests` abren el breaker en una tasa de error sostenida en lugar de una ejecución ininterrumpida de fallos:

```yaml
feature:
  address: localhost:9000
  breaker:
    failure_ratio: 0.5
    min_requests: 10
```

`failure_ratio` es la fracción de solicitudes fallidas (0 < r <= 1) dentro del
`interval` actual que abre el breaker, evaluado solo una vez que se han observado `min_requests`
solicitudes. Cuando se establece `failure_ratio`, toma precedencia
sobre `consecutive_failures`.

HTTP `StatusCodes` y gRPC `Codes` son listas de reemplazo opcionales para clasificación de fallo. Cuando se omiten, se aplican las listas predeterminadas anteriores. Cuando se establecen, solo los
valores configurados cuentan como fallos de breaker, por lo que incluya los predeterminados también
al extender en lugar de reemplazar el comportamiento predeterminado.

### Reintentos del cliente

La configuración del cliente usa la forma compartida `transport/retry.Config` para la mecánica de reintento. Cualquier tipo de configuración que incorpore
`config/client.Config` tiene su propio bloque `retry` bajo esa configuración de cliente. Este ejemplo usa `feature.Config`
solo porque es una configuración de cliente de ese tipo:

```yaml
feature:
  address: localhost:9000
  retry:
    timeout: 1s
    backoff: 100ms
    attempts: 3
    strategy: exponential
```

Al construir manualmente clientes HTTP o gRPC, pase una configuración de reintento específica de transporte a
`transport/http.WithClientRetry(...)` o `transport/grpc.WithClientRetry(...)`. Estas configuraciones incorporan la
mecánica compartida de reintento y agregan clasificación de fallo específica de protocolo:

```go
httpRetry := httpretry.NewConfig(sharedRetry, 429, 502, 503)
grpcRetry := grpcretry.NewConfig(sharedRetry, codes.Unavailable, codes.ResourceExhausted)
```

`NewConfig` devuelve `nil` cuando la configuración de reintento compartida es `nil`, preservando la conexión de opciones de cliente que
desactiva reintentos omitiendo la configuración de reintento.

`attempts` es el número total de intentos, incluyendo la llamada inicial. Un valor
de `0` o `1` significa ningún reintento más allá del primer intento; los valores por encima de `10` son
rechazados durante la validación de configuración. `backoff` es el retraso base entre intentos de reintento.

`strategy` selecciona cómo crece `backoff` entre intentos: `constant` (el
predeterminado) reutiliza el retraso base para cada espera, `exponential` lo duplica en cada
intento, y `fibonacci` lo hace crecer a lo largo de la secuencia de Fibonacci. Un valor no establecido
aplica `constant`, se aplica jitter sobre la estrategia elegida, y cualquier
otro valor se rechaza durante la validación de configuración.

`timeout` es específico del transporte. Los reintentos unarios gRPC lo aplican por intento, por lo que
el tiempo transcurrido total puede incluir múltiples timeouts de intento más backoff a menos que el
contexto del llamador termine primero. Los reintentos HTTP no crean un timeout por intento poseído por el reintento; limite las llamadas HTTP salientes con el contexto de solicitud o
`http.Client.Timeout`.

`max_backoff` limita la duración de backoff por intento, aplicado antes del jitter. Es
más útil con el crecimiento `exponential` y `fibonacci`, que de otro modo crecerían
sin límites a través de los intentos. Un valor cero (el predeterminado) deja el backoff
sin límite:

```yaml
feature:
  address: localhost:9000
  retry:
    backoff: 1s
    strategy: exponential
    attempts: 10
    max_backoff: 30s
```

HTTP `StatusCodes` y gRPC `Codes` son listas de reemplazo opcionales para clasificación de fallo. Cuando se omiten, se aplican las listas predeterminadas a continuación. Cuando se establecen, solo los
valores configurados son reintentables, por lo que incluya los predeterminados también al extender
en lugar de reemplazar el comportamiento predeterminado. Los valores HTTP deben ser códigos de estado 4xx o 5xx. Los valores gRPC deben ser valores `codes.Code` no OK.

La política de reintento predeterminada es intencionalmente conservadora:

- HTTP reintenta métodos seguros de efecto secundario (`GET`, `HEAD`, `OPTIONS`) o solicitudes con un `Request-Id`.
- HTTP reintenta fallos de respuesta/estado solo para `429 Too Many Requests` y `503 Service Unavailable`, más errores de transporte seleccionados clasificados por `retryablehttp.DefaultRetryPolicy`.
- gRPC reintenta métodos de lectura estilo AIP nombrados `Get*` o `List*`, o llamadas con un `Request-Id`.
- gRPC reintenta solo `Unavailable` por defecto.

Las respuestas reintentables HTTP con un retraso `Retry-After` válido mayor que el
backoff jitter mínimo suprimen otro intento y devuelven la respuesta actual. Los errores de estado reintentables gRPC con `google.rpc RetryInfo.retry_delay`
usan la misma política de supresión.

`Request-Id` identifica la solicitud lógica, no un intento de cable individual.
Los servicios que permiten escrituras reintentadas deberían tratarlo como la clave de idempotencia y
desduplicar intentos repetidos cuando el procesamiento duplicado sería inseguro.

---

## 🔑 Criptografía

La configuración raíz de criptografía es `crypto.Config` y soporta múltiples tipos de clave. La mayoría de los campos son cadenas de origen.

Ejemplo:

```yaml
crypto:
  aes:
    key: file:test/secrets/aes
  ed25519:
    public: file:test/secrets/ed25519_public
    private: file:test/secrets/ed25519_private
  hmac:
    key: file:test/secrets/hmac
  rsa:
    public: file:test/secrets/rsa_public
    private: file:test/secrets/rsa_private
  ssh:
    public: file:test/secrets/ssh_public
    private: file:test/secrets/ssh_private
```

> [!NOTE]
> - Las claves AES deben ser de 16/24/32 bytes después de resolver la cadena de origen.
> - Las claves HMAC deberían ser secretos de alta entropía y deben permanecer privadas.
> - Las claves RSA esperan bloques PEM PKCS#1 (`RSA PUBLIC KEY` / `RSA PRIVATE KEY`) y deben ser de al menos 4096 bits.
> - Ed25519 espera bloques PEM PKIX `PUBLIC KEY` y PKCS#8 `PRIVATE KEY`.
> - Las claves SSH deben ser claves SSH Ed25519: las claves públicas usan formato `authorized_keys` y las claves privadas usan formato de clave privada SSH.

Las APIs de encriptación AES y RSA aceptan `crypto.Message`. `Data` se cifra o
descifra, mientras que `Meta` es contexto autenticado que debe coincidir durante
el descifrado. AES-GCM usa `Meta` como datos asociados; RSA-OAEP lo usa como la
etiqueta OAEP.

### Dependencias de Criptografía

![Dependencies](./assets/crypto.png)

---

## 🛠️ Endpoints de depuración

Configuración de servidor de depuración:

```yaml
debug:
  address: tcp://localhost:6060
  timeout: 10s
```

Habilitar TLS:

```yaml
debug:
  tls:
    cert: file:test/certs/cert.pem
    key: file:test/certs/key.pem
    ca: file:test/certs/rootCA.pem
```

El TLS de depuración usa el mismo contrato TLS del lado del servidor que los transportes: `cert` y `key`
son requeridos siempre que se configure material TLS, y `ca` agrega verificación de certificado de cliente
para mTLS.

Todos los endpoints de depuración están namespeados por nombre de servicio: `/<name>/debug/...`.

> [!WARNING]
> Si `debug.address` se omite mientras la depuración está habilitada, el servidor de depuración se une a `tcp://:6060`. Establezca una dirección explícita, TLS/mTLS, y controles de red o política apropiados para el despliegue.

### statsviz

```http
GET http://localhost:6060/<name>/debug/statsviz
```

<https://github.com/arl/statsviz>

### pprof

```http
GET http://localhost:6060/<name>/debug/pprof/
GET http://localhost:6060/<name>/debug/pprof/cmdline
GET http://localhost:6060/<name>/debug/pprof/profile
GET http://localhost:6060/<name>/debug/pprof/symbol
GET http://localhost:6060/<name>/debug/pprof/trace
```

<https://pkg.go.dev/net/http/pprof>

### fgprof

```http
GET http://localhost:6060/<name>/debug/fgprof?seconds=10
```

<https://pkg.go.dev/github.com/felixge/fgprof>

---

## 🧑‍💻 Desarrollo

### Estilo

Este repositorio generalmente sigue la [Guía de Estilo Go de Uber](https://github.com/uber-go/guide/blob/master/style.md).

Los identificadores Go exportados deberían tener comentarios GoDoc, y cada comentario debería comenzar con el nombre del identificador o `Deprecated:`.

### Dependencias de Desarrollo

Los objetivos comunes del repositorio esperan estas herramientas en `PATH`:

- `make`
- `gotestsum` para `make specs`
- `fieldalignment` para `make lint`
- `golangci-lint` para cobertura completa de `make lint` (el wrapper es no-op cuando falta)
- `govulncheck` y `trivy` para `make sec`
- `mkcert` para fixtures TLS locales y `make create-certs`
- `buf` para `make generate`
- `goda` y Graphviz `dot` para `make diagrams`

### Configuración (repo)

Este repositorio usa un submódulo git `bin/` para targets de `make`.

```sh
git submodule sync
git submodule update --init

mkcert -install
make create-certs

make dep
```

Si la obtención del submódulo falla, asegúrese de que el acceso SSH de GitHub esté configurado (`.gitmodules` usa URLs `git@github.com:...`).

### Descubrir targets

```sh
make help
```

### Dependencias (flujo de `vendor/`)

```sh
make dep
```

`make dep` ejecuta:

- `go mod download`
- `go mod tidy`
- `go mod vendor`

Las pruebas se ejecutan con `-mod vendor`, por lo que después de cambios de dependencia ejecute `make dep` antes de `make specs`.

### Dependencias de integración local

`make start` usa el entorno Docker compartido desde el hermano
`../docker` repo. Requiere Docker y puede requerir acceso SSH de GitHub si ese
repositorio hermano debe ser obtenido.

Iniciar servicios requeridos:

```sh
make start
```

Detenerlos:

```sh
make stop
```

### Pruebas

Ejecutar pruebas unitarias con race + cobertura:

```sh
make specs
```

Artefactos:

- JUnit XML: `test/reports/specs.xml`
- Perfil de cobertura: `test/reports/profile.cov`

### Lint y formato

```sh
make lint
make fix-lint
make format
```

### Verificaciones de seguridad

```sh
make sec
```

### Benchmarks

```sh
make benchmarks
make http-benchmarks
make grpc-benchmarks
make limiter-benchmarks
make sql-benchmarks
make cache-benchmarks
make bytes-benchmarks
make strings-benchmarks
make id-benchmarks
make net-http-benchmarks
make http-content-benchmarks
```

### Pruebas Fuzz

```sh
make fuzzes
make bytes-fuzz
make time-fuzz
make encoding-fuzz
make compress-fuzz
make net-fuzz
make package=encoding/json name=FuzzUnmarshal fuzztime=10s fuzz
```

### Informes de cobertura

```sh
make coverage
make html-coverage
make func-coverage
```

### Generación de código (Buf)

Los targets de generación raíz son para los fixtures protobuf de `internal/test`. Después
de cambiar esos fixtures, regénerelos. Para coincidir con la verificación de output estancado de CI,
ejecute `make generate-stale` desde un worktree limpio, o después de staged los cambios de fixture y archivo generados intencionados:

```sh
make generate
make generate-stale
```

### Diagramas de arquitectura

```sh
make diagrams
make crypto-diagram
make database-diagram
make telemetry-diagram
make transport-diagram
```
