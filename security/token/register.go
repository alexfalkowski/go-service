package token

type generatorRegister map[string]Generator

type verifierRegister map[string]Verifier

var (
	genRegister = generatorRegister{}
	verRegister = verifierRegister{}
)

func init() {
	RegisterGenerator("none", NewGenerator())
	RegisterVerifier("none", NewVerifier())
}

// RegisterGenerator for token.
func RegisterGenerator(kind string, g Generator) {
	genRegister[kind] = g
}

// RegisterVerifier for token.
func RegisterVerifier(kind string, v Verifier) {
	verRegister[kind] = v
}
