package server

// Config contains the configuration information needed to start
// the user service.
type Config struct {
	DevMode            bool   `env:"DEV_MODE,default=false"`
	DBURL              string `env:"DATABASE_URL,required"`
	Port               int    `env:"PORT,default=8080"`
	CSRFSecret         string `env:"CSRF_SECRET"`
	CORSOrigins        string `env:"CORS_ORIGINS"`
	MJPublicKey        string `env:"MAILJET_API_KEY_PUBLIC"`
	MJPrivateKey       string `env:"MAILJET_API_KEY_PRIVATE"`
	SimultaneousEmails int    `env:"SIMULTANEOUS_EMAILS,default=10"`
}
