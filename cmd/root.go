/*
Copyright © 2022
Author Bhakiyaraj Kalimuthu
Email bhakiya.kalimuthu@gmail.com
*/

package cmd

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"loadgen/internal"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	workerPoolSize = 5     // default worker pool size
	env            = "dev" // development env
)

var (
	//go:embed payload.json
	payloadBytes []byte
)

// rootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		Use:   "loadgen",
		Short: "load generator",
		Long:  `load generator can generate rpc http load to the configured URL`,
		Run:   runRootCmd,
	}
	rootArgs struct {
		env        string        // environment dev /stage
		url        string        // url where notification to be sent
		interval   time.Duration // interval in which notification to be sent
		numOfBatch int           // number of batch request
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	root := rootCmd.Flags()
	root.StringVarP(&rootArgs.env, "env", "e", "dev", "test or staging environment")
	root.StringVarP(&rootArgs.url, "url", "u", "", "URL to which load to generated / request to be sent")
	root.DurationVarP(&rootArgs.interval, "interval", "i", 100*time.Second, "Load interval")
	root.IntVarP(&rootArgs.numOfBatch, "numOfBatch", "n", 1, "Number of batch request / Each batch contains min 10 requests")

	//cobra.MarkFlagRequired(root, "url")
}

func runRootCmd(cmd *cobra.Command, args []string) {
	// logger setup
	l := loggerSetup()

	// init clock
	clock := internal.NewClock()
	defer func() {
		l.Info("Time taken to complete", zap.Duration("time_taken", <-clock.Since()))
	}()

	if rootArgs.env == "" {
		rootArgs.env = env
	}
	if !isValidURL(rootArgs.url) {

		rootArgs.url = defaultURL(rootArgs.env)
		l.Warn(fmt.Sprintf("url not provided, loadgen pointing to default url: %s", rootArgs.url))
	}

	// producer channel
	pChan := make(chan []*internal.JsonRpcRequest, 1)
	// consumer channel
	cChan := make(chan []*internal.JsonRpcRequest, workerPoolSize)

	// create http client
	httpClient := internal.NewHttpClient(l, rootArgs.url)

	// create notifier
	notifier := internal.NewGenerator(l, httpClient, rootArgs.interval, pChan, cChan)

	// setup cancellation context and wait group
	// root background with cancellation support
	ctx, cancel := context.WithCancel(context.Background())
	wg := new(sync.WaitGroup)

	// start notifier and pass the cancellation ctx
	go notifier.Start(ctx)

	// start workers and add worker pool
	wg.Add(workerPoolSize)
	for i := 1; i <= workerPoolSize; i++ {
		go notifier.Process(wg, i)
	}

	doneCh := make(chan os.Signal, 1)

	// send in request
	go func() {
		for i := 0; i < rootArgs.numOfBatch; i++ {
			req := createRequest()
			pChan <- req // send in data to producer channel
		}
	}()

	// handle manual interruption
	signal.Notify(doneCh, syscall.SIGINT, syscall.SIGTERM)

	switch <-doneCh { // blocks here until interrupted
	case syscall.SIGINT, syscall.SIGTERM:
		l.Warn("CTRL-C received.Terminating......")
	default:
		l.Warn("file read is completed,exiting......")
	}
	signal.Stop(doneCh)

	// handle shut down
	cancel() // cancel context
	// even if cancellation received, current running job will be not be interrupted until it completes
	wg.Wait() // wait for the workers to be completed
	l.Warn("All jobs are done, shutting down")

}

// loggerSetup setup zap logger
func loggerSetup() *zap.Logger {
	if env == "prod" {
		logger, err := zap.NewProduction()
		if err != nil {
			log.Fatalf("failed to create zap logger : %v", err)
		}
		logger.Info("logger setup done")
		return logger
	}

	// setup dev logger to show different colors
	cfg := zap.NewDevelopmentEncoderConfig()
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	log := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	))
	log.Info("logger setup done")
	return log
}

func isValidURL(URL string) bool {
	if rootArgs.url == "" {
		return false
	}

	// parse url if valid
	_, err := url.ParseRequestURI(URL)
	if err != nil {
		fmt.Printf("Error: invalid url %v", err)
		return false
	}
	return true
}

func createRequest() []*internal.JsonRpcRequest {
	var batch []*internal.JsonRpcRequest
	if err := json.Unmarshal(payloadBytes, &batch); err != nil {
		panic(err)
	}
	//txList := internal.NewTxList()
	//for _, tx := range txList {
	//	rawTxReq := internal.NewJsonRpcRequest(1, "eth_sendRawTransaction", []interface{}{tx})
	//	batch = append(batch, rawTxReq)
	//}
	return batch
}

// TODO: Refactor, too many redundant and hardcoded values
func defaultURL(env string) string {
	switch env {
	case "dev":
		return "http://localhost:9000"
	case "stage":
		return "https://rpc-staging.flashbots.net"
	}
	return "http://localhost:9000"
}
