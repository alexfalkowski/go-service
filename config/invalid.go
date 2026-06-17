package config

type invalidSourceDecoder struct{}

func (invalidSourceDecoder) Decode(any) error {
	return ErrInvalidSource
}
