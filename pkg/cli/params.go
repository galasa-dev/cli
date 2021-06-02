package cli

type GalasaParams struct {
	bootstrap string
}

var _ Params = (*GalasaParams)(nil)

func (p *GalasaParams) SetBootstrap(url string) {
	p.bootstrap = url
}
