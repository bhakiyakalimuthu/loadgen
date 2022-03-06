package pkg

//type Generator_ struct {
//	log           *zap.Logger
//	client        internal.HttpClient
//	url           string
//	waitPeriod    time.Duration // default 2 second
//	numberOfRound int           // each round 100 request
//}
//
//func NewGenerator_(log *zap.Logger, client internal.HttpClient, url string, waitPeriod time.Duration, numberOfRound int) *Generator_ {
//	return &Generator_{
//		log:           log,
//		client:        client,
//		url:           url,
//		waitPeriod:    waitPeriod,
//		numberOfRound: numberOfRound,
//	}
//}
//
//func (g *Generator_) Generate() {
//	wg := new(sync.WaitGroup)
//	for i := 0; i < 10; i++ {
//		wg.Add(1)
//		go func() {
//			defer wg.Done()
//			g.Load()
//		}()
//	}
//	wg.Wait()
//}
//
//func (g *Generator_) Load() {
//	req := createRequest()
//	wg := new(sync.WaitGroup)
//	for i := 0; i < 10; i++ {
//		wg.Add(1)
//		go func() {
//			defer wg.Done()
//			if err := g.client.GenerateLoad(req); err != nil {
//				g.log.Error("http client make request failed", zap.Error(err))
//			}
//		}()
//	}
//	wg.Wait()
//}
//
//func createRequest() []*internal.JsonRpcRequest {
//	var batch []*internal.JsonRpcRequest
//	txList := internal.NewTxList()
//	for _, tx := range txList {
//		rawTxReq := internal.NewJsonRpcRequest(1, "eth_sendRawTransaction", []interface{}{tx})
//		batch = append(batch, rawTxReq)
//	}
//	return batch
//}

//const (
//	Dev   = "dev"
//	Stage = "stage"
//)
//
//type config struct {
//	Env string `env:"FB_ENV" envDefault:"dev"`
//}
//
//func main() {
//	cfg := &config{}
//	logger := loggerSetup(cfg)
//	client := internal.NewHttpClient(logger)
//	generator := internal.NewGenerator(logger, client, newURL(cfg), time.Second*2, 10)
//	generator.Generate()
//}
//
//func newURL(c *config) string {
//	switch c.Env {
//	case Dev:
//		return "http://localhost:9000"
//	case Stage:
//		return "https://rpc-staging.flashbots.net"
//	default:
//		return "http://localhost:9000"
//	}
//	return "https://rpc-staging.flashbots.net"
//}
//
//// loggerSetup setup zap logger
//func loggerSetup(c *config) *zap.Logger {
//	// setup dev logger to show different colors
//	cfg := zap.NewDevelopmentEncoderConfig()
//	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
//	log := zap.New(zapcore.NewCore(
//		zapcore.NewConsoleEncoder(cfg),
//		zapcore.AddSync(colorable.NewColorableStdout()),
//		zapcore.InfoLevel,
//	))
//	log.Info("logger setup done")
//	return log
//}
