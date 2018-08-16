package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/urfave/cli"
)

type resultOutput struct {
	HostIdentifier string              `json:"host"`
	Rows           []map[string]string `json:"rows"`
}

func queryCommand() cli.Command {
	var (
		flHosts, flLabels, flQuery string
		flDebug, flQuiet, flExit   bool
	)
	return cli.Command{
		Name:      "query",
		Usage:     "Run a live query",
		UsageText: `fleetctl query [options]`,
		Flags: []cli.Flag{
			configFlag(),
			contextFlag(),
			cli.StringFlag{
				Name:        "hosts",
				EnvVar:      "HOSTS",
				Value:       "",
				Destination: &flHosts,
				Usage:       "Comma separated hostnames to target",
			},
			cli.StringFlag{
				Name:        "labels",
				EnvVar:      "LABELS",
				Value:       "",
				Destination: &flLabels,
				Usage:       "Comma separated label names to target",
			},
			cli.BoolFlag{
				Name:        "quiet",
				EnvVar:      "QUIET",
				Destination: &flQuiet,
				Usage:       "Quiet and don't show anything other the results",
			},
			cli.BoolFlag{
				Name:        "exit",
				EnvVar:      "EXIT",
				Destination: &flExit,
				Usage:       "Only print results (no status information)",
			},
			cli.StringFlag{
				Name:        "query",
				EnvVar:      "QUERY",
				Value:       "",
				Destination: &flQuery,
				Usage:       "Query to run",
			},
			cli.BoolFlag{
				Name:        "debug",
				EnvVar:      "DEBUG",
				Destination: &flDebug,
				Usage:       "Whether or not to enable debug logging",
			},
		},
		Action: func(c *cli.Context) error {
			fleet, err := clientFromCLI(c)
			if err != nil {
				return err
			}

			if flHosts == "" && flLabels == "" {
				return errors.New("No hosts or labels targeted")
			}

			if flQuery == "" {
				return errors.New("No query specified")
			}

			hosts := strings.Split(flHosts, ",")
			labels := strings.Split(flLabels, ",")

			res, err := fleet.LiveQuery(flQuery, labels, hosts)
			if err != nil {
				return err
			}

			tick := time.NewTicker(100 * time.Millisecond)
			defer tick.Stop()

			// See charsets at
			// https://godoc.org/github.com/briandowns/spinner#pkg-variables
			s := spinner.New(spinner.CharSets[24], 200*time.Millisecond)
			s.Writer = os.Stderr
			if !flQuiet {
				s.Start()
			}

			for {
				select {
				case hostResult := <-res.Results():
					out := resultOutput{hostResult.Host.HostName, hostResult.Rows}
					if err := json.NewEncoder(os.Stdout).Encode(out); err != nil {
						fmt.Fprintf(os.Stderr, "Error writing output: %s\n", err)
					}

				case err := <-res.Errors():
					fmt.Fprintf(os.Stderr, "Error talking to server: %s\n", err.Error())

				case <-tick.C:
					// Print status message to stderr
					status := res.Status()
					totals := res.Totals()
					var percentTotal, percentOnline float64
					var responded, total, online uint
					if status != nil && totals != nil {
						total = totals.Total
						online = totals.Online
						responded = status.ActualResults
						if total > 0 {
							percentTotal = 100 * float64(responded) / float64(total)
						}
						if online > 0 {
							percentOnline = 100 * float64(responded) / float64(online)
						}
					}

					if responded >= online && flExit {
						return nil
					}

					if !flQuiet {
						s.Suffix = fmt.Sprintf(
							"  %.f%% responded (%.f%% online) | %d/%d targeted hosts (%d/%d online)",
							percentTotal, percentOnline,
							responded, total,
							responded, online,
						)
					}
				}
			}
		},
	}
}
