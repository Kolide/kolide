package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/kolide/fleet/server/kolide"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

type specMetadata struct {
	Kind    string      `json:"kind"`
	Version string      `json:"apiVersion"`
	Spec    interface{} `json:"spec"`
}

type specGroup struct {
	Queries []*kolide.QuerySpec
	Packs   []*kolide.PackSpec
	Labels  []*kolide.LabelSpec
	Options *kolide.OptionsSpec
}

func specGroupFromBytes(b []byte) (*specGroup, error) {
	specs := &specGroup{
		Queries: []*kolide.QuerySpec{},
		Packs:   []*kolide.PackSpec{},
		Labels:  []*kolide.LabelSpec{},
	}

	for _, spec := range strings.Split(string(b), "---") {
		if strings.TrimSpace(spec) == "" {
			continue
		}

		var s specMetadata
		if err := yaml.Unmarshal([]byte(spec), &s); err != nil {
			return nil, errors.Wrap(err, "error unmarshaling spec")
		}

		if s.Spec == nil {
			return nil, errors.Errorf("no spec field on %q document", s.Kind)
		}

		specBytes, err := yaml.Marshal(s.Spec)
		if err != nil {
			return nil, errors.Errorf("error marshaling spec for %q kind", s.Kind)
		}

		switch strings.ToLower(s.Kind) {
		case "query":
			var querySpec *kolide.QuerySpec
			if err := yaml.Unmarshal(specBytes, &querySpec); err != nil {
				return nil, errors.Wrap(err, "error unmarshaling query spec")
			}
			specs.Queries = append(specs.Queries, querySpec)

		case "pack":
			var packSpec *kolide.PackSpec
			if err := yaml.Unmarshal(specBytes, &packSpec); err != nil {
				return nil, errors.Wrap(err, "error unmarshaling pack spec")
			}
			specs.Packs = append(specs.Packs, packSpec)

		case "label":
			var labelSpec *kolide.LabelSpec
			if err := yaml.Unmarshal(specBytes, &labelSpec); err != nil {
				return nil, errors.Wrap(err, "error unmarshaling label spec")
			}
			specs.Labels = append(specs.Labels, labelSpec)

		case "options":
			if specs.Options != nil {
				return nil, errors.New("options defined twice in the same file")
			}

			var optionSpec *kolide.OptionsSpec
			if err := yaml.Unmarshal(specBytes, &optionSpec); err != nil {
				return nil, errors.Wrap(err, "error unmarshaling option spec")
			}
			specs.Options = optionSpec

		default:
			return nil, errors.Errorf("unknown kind %q", s.Kind)
		}
	}

	return specs, nil
}

func applyCommand() cli.Command {
	var (
		flFilename string
		flDebug    bool
	)
	return cli.Command{
		Name:      "apply",
		Usage:     "Apply files to declaratively manage osquery configurations",
		UsageText: `fleetctl apply [options]`,
		Flags: []cli.Flag{
			configFlag(),
			contextFlag(),
			cli.StringFlag{
				Name:        "f",
				EnvVar:      "FILENAME",
				Value:       "",
				Destination: &flFilename,
				Usage:       "A file to apply",
			},
			cli.BoolFlag{
				Name:        "debug",
				EnvVar:      "DEBUG",
				Destination: &flDebug,
				Usage:       "Whether or not to enable debug logging",
			},
		},
		Action: func(c *cli.Context) error {
			if flFilename == "" {
				return errors.New("-f must be specified")
			}

			b, err := ioutil.ReadFile(flFilename)
			if err != nil {
				return err
			}

			fleet, err := clientFromCLI(c)
			if err != nil {
				return err
			}

			specs, err := specGroupFromBytes(b)
			if err != nil {
				return err
			}

			if len(specs.Queries) > 0 {
				if err := fleet.ApplyQuerySpecs(specs.Queries); err != nil {
					return errors.Wrap(err, "error applying queries")
				}
				fmt.Printf("[+] applied %d queries\n", len(specs.Queries))
			}

			if len(specs.Labels) > 0 {
				if err := fleet.ApplyLabelSpecs(specs.Labels); err != nil {
					return errors.Wrap(err, "error applying labels")
				}
				fmt.Printf("[+] applied %d labels\n", len(specs.Labels))
			}

			if len(specs.Packs) > 0 {
				if err := fleet.ApplyPackSpecs(specs.Packs); err != nil {
					return errors.Wrap(err, "error applying packs")
				}
				fmt.Printf("[+] applied %d packs\n", len(specs.Packs))
			}

			return nil
		},
	}
}
