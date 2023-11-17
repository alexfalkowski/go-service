package token

type generatorRegister map[string]Generator

type verifierRegister map[string]Verifier

var (
	genRegister = generatorRegister{}
	verRegister = verifierRegister{}
)

// RegisterGenerator for token.
func RegisterGenerator(kind string, g Generator) {
	genRegister[kind] = g
}

// RegisterVerifier for token.
func RegisterVerifier(kind string, v Verifier) {
	verRegister[kind] = v
}
