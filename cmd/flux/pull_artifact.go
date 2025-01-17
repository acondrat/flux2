/*
Copyright 2022 The Flux authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	oci "github.com/fluxcd/pkg/oci/client"
)

var pullArtifactCmd = &cobra.Command{
	Use:   "artifact",
	Short: "Pull artifact",
	Long: `The pull artifact command downloads and extracts the OCI artifact content to the given path.
The pull command uses the credentials from '~/.docker/config.json'.`,
	Example: `  # Pull an OCI artifact created by flux from GHCR
  flux pull artifact oci://ghcr.io/org/manifests/app:v0.0.1 --output ./path/to/local/manifests
`,
	RunE: pullArtifactCmdRun,
}

type pullArtifactFlags struct {
	output string
}

var pullArtifactArgs pullArtifactFlags

func init() {
	pullArtifactCmd.Flags().StringVarP(&pullArtifactArgs.output, "output", "o", "", "path where the artifact content should be extracted.")
	pullCmd.AddCommand(pullArtifactCmd)
}

func pullArtifactCmdRun(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("artifact URL is required")
	}
	ociURL := args[0]

	if pullArtifactArgs.output == "" {
		return fmt.Errorf("invalid output path %s", pullArtifactArgs.output)
	}

	if fs, err := os.Stat(pullArtifactArgs.output); err != nil || !fs.IsDir() {
		return fmt.Errorf("invalid output path %s", pullArtifactArgs.output)
	}

	ociClient := oci.NewLocalClient()
	url, err := oci.ParseArtifactURL(ociURL)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), rootArgs.timeout)
	defer cancel()

	logger.Actionf("pulling artifact from %s", url)

	meta, err := ociClient.Pull(ctx, url, pullArtifactArgs.output)
	if err != nil {
		return err
	}

	logger.Successf("source %s", meta.Source)
	logger.Successf("revision %s", meta.Revision)
	logger.Successf("digest %s", meta.Digest)
	logger.Successf("artifact content extracted to %s", pullArtifactArgs.output)

	return nil
}
