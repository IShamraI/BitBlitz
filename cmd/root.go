package cmd

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/IShamraI/BitBlitz/internal/bitcion"
	"github.com/IShamraI/BitBlitz/internal/sys"
	tgnotifyer "github.com/IShamraI/BitBlitz/internal/tg_notifyer"
	"github.com/spf13/cobra"
)

var (
	botToken     string = ""
	chatID       string = ""
	outputFile   string = "output.csv"
	threadsNum   int    = 1
	iterInterval        = 5 * time.Second
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
		var mutex sync.Mutex

		ip, err := sys.GetLocalIPv4()
		if err != nil {
			log.Fatalf("Failed to get local IP address: %s", err)
		}

		notifyer := tgnotifyer.New(tgnotifyer.WithToken(botToken), tgnotifyer.WithChatID(chatID))
		notifyer.Notify(fmt.Sprintf("Started BTC Finder on: %s", ip))

		for i := 0; i < threadsNum; i++ {
			wg.Add(1)
			go bitcion.NewWorker(i,
				bitcion.WithMutex(&mutex),
				bitcion.WithWaitGroup(&wg),
				bitcion.WithOutput(outputFile),
				bitcion.WithNotifyer(notifyer),
				bitcion.WithCheckInterval(iterInterval),
			).Start()
			time.Sleep(iterInterval / time.Duration(threadsNum))
		}

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
	rootCmd.Flags().IntVarP(&threadsNum, "threads", "n", threadsNum, "Number of threads")
	rootCmd.Flags().DurationVarP(&iterInterval, "interval", "i", iterInterval, "Interval between iterations")
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
	if threadsNum == 1 {
		threadsEnv := os.Getenv("THREADS")
		if threadsEnv != "" {
			_, err := fmt.Sscanf(threadsEnv, "%d", &threadsNum)
			if err != nil {
				log.Fatalf("Failed to parse THREADS: %v", err)
			}
		}
	}
	if iterInterval == 5*time.Second {
		iterIntervalEnv := os.Getenv("INTERVAL")
		if iterIntervalEnv != "" {
			parsedInterval, err := time.ParseDuration(iterIntervalEnv)
			if err != nil {
				log.Fatalf("Failed to parse INTERVAL: %v", err)
			}
			iterInterval = parsedInterval
		}
	}
}
