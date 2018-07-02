package servo

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
)

/**
* TODO:
* . finish documentation
**/

var (
	// from https://golang.org/pkg/crypto/tls/#pkg-constants
	ciphers map[string]uint16 = map[string]uint16{
		"TLS_RSA_WITH_RC4_128_SHA":                tls.TLS_RSA_WITH_RC4_128_SHA,
		"TLS_RSA_WITH_3DES_EDE_CBC_SHA":           tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
		"TLS_RSA_WITH_AES_128_CBC_SHA":            tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		"TLS_RSA_WITH_AES_256_CBC_SHA":            tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		"TLS_RSA_WITH_AES_128_CBC_SHA256":         tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
		"TLS_RSA_WITH_AES_128_GCM_SHA256":         tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		"TLS_RSA_WITH_AES_256_GCM_SHA384":         tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		"TLS_ECDHE_ECDSA_WITH_RC4_128_SHA":        tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
		"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA":    tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
		"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA":    tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		"TLS_ECDHE_RSA_WITH_RC4_128_SHA":          tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
		"TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA":     tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
		"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA":      tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA":      tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256": tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
		"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256":   tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
		"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256":   tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256": tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384":   tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384": tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305":    tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305":  tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
	}

	// from https://golang.org/pkg/crypto/tls/#CurveID
	curves map[string]tls.CurveID = map[string]tls.CurveID{
		"CurveP256": tls.CurveP256,
		"CurveP384": tls.CurveP384,
		"CurveP521": tls.CurveP521,
		"X25519":    tls.X25519,
	}

	// TLS dfault preferences
	defaultCurves []tls.CurveID = []tls.CurveID{
		tls.CurveP521,
		tls.CurveP384,
		tls.X25519,
		tls.CurveP256,
	}

	defaultCiphers []uint16 = []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	}
)

// TLSOpts is for configuring TLS connection.
type TLSOpts struct {
	Curves           []tls.CurveID
	Ciphers          []uint16
	Cert             string
	Key              string
	Ca               string
	HandshakeTimeout int
	ShouldVerify     bool
}

// genTLSConfig is for generating TLS configuration
// and performs setup according to supplied options.
func genTLSConfig(opts TLSOpts) (*tls.Config, error) {
	var (
		config   *tls.Config
		cert     tls.Certificate
		certpool *x509.CertPool
		leaf     *x509.Certificate
		err      error
	)
	cert, err = tls.LoadX509KeyPair(opts.Cert, opts.Key)
	if err != nil {
		log.Println("[server(tls)]: no certificate is given.", err)
		return nil, err
	}
	leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		log.Println("[server(tls)]: failed to parse certificate.")
		return nil, err
	}
	cert.Leaf = leaf
	if opts.Ciphers == nil {
		opts.Ciphers = make([]uint16, len(defaultCiphers))
		opts.Ciphers = append(opts.Ciphers, defaultCiphers...)
	}
	if opts.Curves == nil {
		opts.Curves = make([]tls.CurveID, len(defaultCurves))
		opts.Curves = append(opts.Curves, defaultCurves...)
	}
	config = &tls.Config{
		Certificates:             []tls.Certificate{cert},
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         opts.Curves,
		CipherSuites:             opts.Ciphers,
	}
	if opts.ShouldVerify {
		config.ClientAuth = tls.RequireAndVerifyClientCert
	}
	if opts.Ca != _EMPTY_ {
		certpool = x509.NewCertPool()
		var rpem []byte
		rpem, err = ioutil.ReadFile(opts.Ca)
		if rpem == nil || err != nil {
			return nil, err
		}
		if !certpool.AppendCertsFromPEM([]byte(rpem)) {
			return nil, EINVAL
		}
		config.ClientCAs = certpool
	}

	return config, nil
}
