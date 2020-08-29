package encryptor

// Config represent secrets configuration
type Config struct {
	AEADKey        string `mapstructure:"aeadkey"`
	PublicNonce    string `mapstructure:"public_nonce"`
	GenericHashKey string `mapstructure:"generic_hashkey"`
}
