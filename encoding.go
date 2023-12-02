package abair

type Encoder interface {
	Encode(any) error
}

type Decoder interface {
	Decode(any) error
}
