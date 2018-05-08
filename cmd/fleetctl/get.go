package main

import (
	"fmt"
	"os"

	"github.com/ghodss/yaml"
	"github.com/kolide/fleet/server/kolide"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

type specGeneric struct {
	Kind    string      `json:"kind"`
	Version string      `json:"apiVersion"`
	Spec    interface{} `json:"spec"`
}

func defaultTable() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	return table
}

func getQueriesCommand() cli.Command {
	return cli.Command{
		Name:    "queries",
		Aliases: []string{"query", "q"},
		Usage:   "List information about one or more queries",
		Flags: []cli.Flag{
			configFlag(),
			contextFlag(),
		},
		Action: func(c *cli.Context) error {
			fleet, err := clientFromCLI(c)
			if err != nil {
				return err
			}

			name := c.Args().First()

			// if name wasn't provided, list all queries
			if name == "" {
				queries, err := fleet.GetQueries()
				if err != nil {
					return errors.Wrap(err, "could not list queries")
				}

				if len(queries) == 0 {
					fmt.Println("no queries found")
					return nil
				}

				data := [][]string{}

				for _, query := range queries {
					data = append(data, []string{
						query.Name,
						query.Description,
						query.Query,
					})
				}

				table := defaultTable()
				table.SetHeader([]string{"name", "description", "query"})
				table.AppendBulk(data)
				table.Render()

				return nil
			} else {
				query, err := fleet.GetQuery(name)
				if err != nil {
					return err
				}

				spec := specGeneric{
					Kind:    "query",
					Version: kolide.ApiVersion,
					Spec:    query,
				}

				b, err := yaml.Marshal(spec)
				if err != nil {
					return err
				}

				fmt.Printf(string(b))
				return nil
			}
		},
	}
}

func getPacksCommand() cli.Command {
	return cli.Command{
		Name:    "packs",
		Aliases: []string{"pack", "p"},
		Usage:   "List information about one or more packs",
		Flags: []cli.Flag{
			configFlag(),
			contextFlag(),
		},
		Action: func(c *cli.Context) error {
			fleet, err := clientFromCLI(c)
			if err != nil {
				return err
			}

			name := c.Args().First()

			// if name wasn't provided, list all packs
			if name == "" {
				packs, err := fleet.GetPacks()
				if err != nil {
					return errors.Wrap(err, "could not list packs")
				}

				if len(packs) == 0 {
					fmt.Println("no packs found")
					return nil
				}

				data := [][]string{}

				for _, pack := range packs {
					data = append(data, []string{
						pack.Name,
						pack.Platform,
						pack.Description,
					})
				}

				table := defaultTable()
				table.SetHeader([]string{"name", "platform", "description"})
				table.AppendBulk(data)
				table.Render()

				return nil
			} else {
				pack, err := fleet.GetPack(name)
				if err != nil {
					return err
				}

				spec := specGeneric{
					Kind:    "pack",
					Version: kolide.ApiVersion,
					Spec:    pack,
				}

				b, err := yaml.Marshal(spec)
				if err != nil {
					return err
				}

				fmt.Printf(string(b))
				return nil
			}
		},
	}
}

func getLabelsCommand() cli.Command {
	return cli.Command{
		Name:    "labels",
		Aliases: []string{"label", "l"},
		Usage:   "List information about one or more labels",
		Flags: []cli.Flag{
			configFlag(),
			contextFlag(),
		},
		Action: func(c *cli.Context) error {
			fleet, err := clientFromCLI(c)
			if err != nil {
				return err
			}

			name := c.Args().First()

			// if name wasn't provided, list all labels
			if name == "" {
				labels, err := fleet.GetLabels()
				if err != nil {
					return errors.Wrap(err, "could not list labels")
				}

				if len(labels) == 0 {
					fmt.Println("no labels found")
					return nil
				}

				data := [][]string{}

				for _, label := range labels {
					data = append(data, []string{
						label.Name,
						label.Platform,
						label.Description,
						label.Query,
					})
				}

				table := defaultTable()
				table.SetHeader([]string{"name", "platform", "description", "query"})
				table.AppendBulk(data)
				table.Render()

				return nil
			} else {
				label, err := fleet.GetLabel(name)
				if err != nil {
					return err
				}

				spec := specGeneric{
					Kind:    "label",
					Version: kolide.ApiVersion,
					Spec:    label,
				}

				b, err := yaml.Marshal(spec)
				if err != nil {
					return err
				}

				fmt.Printf(string(b))

				return nil
			}
		},
	}
}
