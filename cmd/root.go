package cmd

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/IShamraI/BitBlitz/internal/bitcoin"
	"github.com/IShamraI/BitBlitz/internal/sys"
	tgnotifyer "github.com/IShamraI/BitBlitz/internal/tg_notifyer"
	"github.com/spf13/cobra"
)

var (
	botToken            string = ""
	chatID              string = ""
	outputFile          string = "output.csv"
	rate                int    = 100
	workerPoolSize      int    = 10
	workerPoolSizeLimit int    = 100
	taskQueueThreshold         = 0.9
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "BitBlitz",
	Short: "Experimental BitBlitz CLI",
	Long:  `Experimental BitBlitz CLI.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		// var mutex sync.Mutex
		// var errLock *locker.Locker
		// var tasksCount int = 0

		// taskQueue := make(chan wallet.Wallet, 100)
		// errLock = locker.New()

		ip, err := sys.GetLocalIPv4()
		if err != nil {
			log.Fatalf("Failed to get local IP address: %s", err)
		}

		notifyer := tgnotifyer.New(tgnotifyer.WithToken(botToken), tgnotifyer.WithChatID(chatID))
		notifyer.Notify(fmt.Sprintf("Started BTC Finder on: %s", ip))

		bitcoin.Some()

		// // Start workers
		// for i := 0; i < workerPoolSize; i++ {
		// 	wg.Add(1)
		// 	log.Printf("Starting worker %d", i)
		// 	go func(workerID int) {
		// 		defer wg.Done()
		// 		bitcoin.NewWorker(workerID,
		// 			bitcoin.WithMutex(&mutex),
		// 			bitcoin.WithWaitGroup(&wg),
		// 			bitcoin.WithOutput(outputFile),
		// 			bitcoin.WithNotifyer(notifyer),
		// 			bitcoin.WithErrLock(errLock),
		// 		).Start(taskQueue)
		// 	}(i)
		// }

		// // Generate tasks with rate limiting
		// wg.Add(1)
		// log.Printf("Starting task generator")
		// go func() {
		// 	defer wg.Done()
		// 	ticker := time.NewTicker(time.Second / time.Duration(rate))
		// 	for _ = range ticker.C {
		// 		// Generate task
		// 		go func() {
		// 			task, err := wallet.GenWallet()
		// 			if err != nil {
		// 				log.Fatalf("Failed to generate wallet: %s", err)
		// 			}
		// 			tasksCount++
		// 			taskQueue <- *task
		// 		}()
		// 	}
		// }()

		// ticker := time.NewTicker(10 * time.Second)
		// log.Printf("Starting status checker")
		// for _ = range ticker.C {
		// 	log.Printf("Qlength: %d, Qcapacity: %d, Qutil: %.2f, Task count: %d, Worker pool size: %d",
		// 		len(taskQueue), cap(taskQueue), float64(len(taskQueue))/float64(cap(taskQueue)), tasksCount, workerPoolSize)
		// 	if float64(len(taskQueue))/float64(cap(taskQueue)) > taskQueueThreshold && workerPoolSize < workerPoolSizeLimit {
		// 		// Add additional worker if queue size exceeds threshold
		// 		log.Println("Adding additional worker")
		// 		wg.Add(1)
		// 		go func(workerID int) {
		// 			defer wg.Done()
		// 			log.Printf("Starting worker %d", workerID)
		// 			bitcoin.NewWorker(workerID,
		// 				bitcoin.WithMutex(&mutex),
		// 				bitcoin.WithWaitGroup(&wg),
		// 				bitcoin.WithOutput(outputFile),
		// 				bitcoin.WithNotifyer(notifyer),
		// 			).Start(taskQueue)
		// 			log.Printf("Worker %d finished", workerID)
		// 		}(workerPoolSize)
		// 		workerPoolSize++
		// 	}
		// }

		wg.Wait()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&botToken, "bot-token", "t", botToken, "Telegram bot token")
	rootCmd.Flags().StringVarP(&chatID, "chat-id", "c", chatID, "Telegram chat ID")
	rootCmd.Flags().StringVarP(&outputFile, "output-file", "o", outputFile, "Output file")
	rootCmd.Flags().IntVarP(&rate, "rate", "r", rate, "Rate of checks per second")
	// rootCmd.Flags().DurationVarP(&iterInterval, "interval", "i", iterInterval, "Interval between iterations")
	// Check if flags are set to default values and fetch values from environment variables if needed
	if botToken == "" {
		botToken = os.Getenv("BOT_TOKEN")
	}
	if chatID == "" {
		chatID = os.Getenv("CHAT_ID")
	}
	if outputFile == "output.csv" {
		outputFile = os.Getenv("OUTPUT_FILE")
	}
	if rate == 100 {
		rateEnv := os.Getenv("RATE")
		if rateEnv != "" {
			_, err := fmt.Sscanf(rateEnv, "%d", &rate)
			if err != nil {
				log.Fatalf("Failed to parse THREADS: %v", err)
			}
		}
	}
}
