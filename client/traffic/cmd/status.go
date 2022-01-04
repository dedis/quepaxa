package cmd

import (
	"bufio"
	"context"
	"fmt"
	"github.com/dedis/backsos/proto/control"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"os"
	"sync"
	"time"
)

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.PersistentFlags().IntVar(&operationType, "operationType", 6, "operation type")
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Fetch status of a quorum of instances",
	Run: func(cmd *cobra.Command, args []string) {
		type report struct {
			name      string
			resp      *control.StatusResponse
			timestamp time.Time
		}

		wg := sync.WaitGroup{}

		reports := make(chan *report, 3)
		for _, in := range cfgQuorum.Instances {
			wg.Add(1)
			go func(name, address string) {
				defer wg.Done()

				// connect to instance
				conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBackoffMaxDelay(60*time.Second))
				if err != nil {
					return
				}
				defer conn.Close()
				client := control.NewControlClient(conn)

				ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
				resp, err := client.Status(ctx, &control.StatusRequest{OperationType: int64(operationType)})
				if err != nil {
					resp = nil
				}
				reports <- &report{
					name:      name,
					resp:      resp,
					timestamp: time.Now(),
				}
				cancel()
				return

			}(in.Name, in.Address)
		}

		// close reports channel once we are done
		go func() {
			wg.Wait()
			close(reports)
		}()

		bw := bufio.NewWriter(os.Stdout)
		//fmt.Println("---Only the last 15 log entries are shown---")
		fmt.Print("NAME\tTime\n") //10 fields
		for r := range reports {
			if r.resp == nil {
				continue
			} else {
				fmt.Printf("%v\t %v\n", r.resp.Name, r.timestamp)
			}

		}
		bw.Flush()

	},
}
