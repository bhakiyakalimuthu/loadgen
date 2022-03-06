/*
Copyright © 2022
Author Bhakiyaraj Kalimuthu
Email bhakiya.kalimuthu@gmail.com
*/

package internal

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Generator is the interface that groups the Start and Process methods
type Generator interface {
	Process(wg *sync.WaitGroup, workerID int)
	Start(ctx context.Context)
}

// generator type
type generator struct {
	logger       *zap.Logger            // logger
	httpClient   HttpClient             // http client for generating load notification
	interval     time.Duration          // interval in which request to be sent
	producerChan chan []*JsonRpcRequest // channel to receive from stdio
	consumerChan chan []*JsonRpcRequest // chanel to consume the data

}

// NewGenerator constructor
func NewGenerator(logger *zap.Logger, httpClient HttpClient, interval time.Duration, producerChan, consumerChan chan []*JsonRpcRequest) Generator {
	return &generator{
		logger:       logger,
		interval:     interval,
		producerChan: producerChan,
		consumerChan: consumerChan,
		httpClient:   httpClient,
	}
}

// Process starts the worker process based on the number items in the consumer channel until it closes
func (n *generator) Process(wg *sync.WaitGroup, workerID int) {
	defer wg.Done()
	for job := range n.consumerChan {
		<-time.After(n.interval) // wait for the provided interval
		n.logger.Debug("starting job", zap.Int("workerID", workerID))
		n.httpClient.GenerateLoad(job) // call http client to make notification
		n.logger.Warn("worker finishing job", zap.Int("workerID", workerID))
	}
	n.logger.Warn("gracefully finishing job", zap.Int("workerID", workerID))
}

// Start acts as a proxy between producer and consumer channel,also supports the graceful cancellation
func (n *generator) Start(ctx context.Context) {
	for {
		select {
		case job := <-n.producerChan: // fetch job from producer
			n.logger.Debug("received msg from consumerChan")
			n.consumerChan <- job // pass job to consumer
		case <-ctx.Done():
			n.logger.Warn("received context cancellation......")
			close(n.consumerChan) // when context is done, close the consumer channel
			return
		}
	}
}
