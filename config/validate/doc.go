// Package validate provides the shared go-playground/validator wiring used to validate decoded
// configuration.
//
// [Validator] is built by [NewValidator], which starts from a base validator and registers every [Validation]
// contributed to the "validations" group. Packages that own a validation rule (such as `bytes`,
// `time`, or `telemetry/otlp`) contribute their own [Validation] through their own [go.uber.org/fx]
// module, tagging the result with [ValidationResult], so this package never needs to know about
// their rules by name.
package validate
