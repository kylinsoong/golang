package as3

type AS3Manager struct {
        as3Validation             bool
        as3Validated              map[string]bool
        sslInsecure               bool
        enableTLS                 string
        tls13CipherGroupReference string
        ciphers                   string
}
