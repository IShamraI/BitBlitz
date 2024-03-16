package bitcion

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/IShamraI/BitBlitz/internal/crypto"
	tgnotifyer "github.com/IShamraI/BitBlitz/internal/tg_notifyer"
)

type BlockCypherResponse struct {
	Balance int `json:"balance"`
}

type Worker struct {
	ID            int
	apiURL        string
	output        string
	notifyer      *tgnotifyer.TgNotifyer
	checkInterval time.Duration
	wg            *sync.WaitGroup
	mutex         *sync.Mutex
}

type Option func(*Worker)

func WithCheckInterval(interval time.Duration) Option {
	return func(w *Worker) {
		w.checkInterval = interval
	}
}

func WithNotifyer(notifyer *tgnotifyer.TgNotifyer) Option {
	return func(w *Worker) {
		w.notifyer = notifyer
	}
}

func WithOutput(output string) Option {
	return func(w *Worker) {
		w.output = output
	}
}

func WithWaitGroup(wg *sync.WaitGroup) Option {
	return func(w *Worker) {
		w.wg = wg
	}
}

func WithMutex(mutex *sync.Mutex) Option {
	return func(w *Worker) {
		w.mutex = mutex
	}
}

func NewWorker(id int, opts ...Option) *Worker {
	worker := &Worker{
		ID:            id,
		apiURL:        "https://api.blockcypher.com/v1/btc/main/addrs/%s/balance",
		notifyer:      nil,
		wg:            &sync.WaitGroup{},
		mutex:         &sync.Mutex{},
		output:        "",
		checkInterval: 10 * time.Second,
	}
	for _, opt := range opts {
		opt(worker)
	}
	return worker
}

// LogWallet appends a line of text to the output file. The output file is a CSV
// file with the format:
// <private key>,<public address>,<current balance in satoshis>
//
// If the output file is not specified, this function does nothing.
func (w *Worker) LogWallet(data string) {
	if w.output == "" {
		return
	}
	w.mutex.Lock() // Acquire the lock to make the log function thread-safe.
	defer w.mutex.Unlock()
	if file, err := os.OpenFile(w.output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		defer file.Close()
		if _, err := file.WriteString(data + "\n"); err != nil {
			log.Printf("Worker %d: Failed to write to file: %s", w.ID, err)
		}
	} else {
		log.Printf("Worker %d: Failed to open file: %s", w.ID, err)
	}
}

// Notify sends a notification to the worker's notifyer. If the notifyer is not
// set, this function does nothing.
//
// The message will be sent to the notifyer as-is, so the message should be a
// meaningful message to the user.
func (w *Worker) Notify(message string) {
	if w.notifyer != nil {
		w.notifyer.Notify(message) // Send the notification to the user.
	}
}

func (w *Worker) Start() {
	var balance int
	w.wg.Add(1)
	ticker := time.NewTicker(w.checkInterval)
	for _ = range ticker.C {
		privateKey, publicAddress, err := crypto.GenerateKeyAndAddress()
		if err != nil {
			log.Printf("Worker %6d: Failed to generate key and address: %s", w.ID, err)
			continue
		}
		balance, err = w.CheckBalance(publicAddress)
		if err != nil {
			log.Printf("Worker %6d: Failed to check balance: %s", w.ID, err)
			continue
		}
		if balance == 0 {
			log.Printf("Worker %6d: No balance found for %s", w.ID, publicAddress)
			continue
		}
		log.Printf("Worker %6d: Logging wallet data", w.ID)
		w.LogWallet(fmt.Sprintf("%s,%s,%d", privateKey, publicAddress, balance))
		log.Printf("Worker %6d: Sending notification", w.ID)
		w.Notify(fmt.Sprintf("Private key: %s\nPublic address: %s\nBalance: %d", privateKey, publicAddress, balance))
	}
}

func (w *Worker) CheckBalance(address string) (int, error) {
	return w.checkBalance(address)
}

func (w *Worker) checkBalance(address string) (int, error) {
	url := fmt.Sprintf(w.apiURL, address)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var response BlockCypherResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, err
	}

	return response.Balance, nil
}
