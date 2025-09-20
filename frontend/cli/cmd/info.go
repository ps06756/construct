package cmd

import (
	"github.com/spf13/cobra"
)

type infoOptions struct {
	RenderOptions RenderOptions
}

type VersionInfo struct {
	Version   string `detail:"default"`
	Commit    string `detail:"default"`
	BuildDate string `detail:"default"`
}

func NewInfoCmd() *cobra.Command {
	options := infoOptions{}
	cmd := &cobra.Command{
		Use:     "info",
		Short:   "Print information about the Construct CLI",
		GroupID: "system",
		RunE: func(cmd *cobra.Command, args []string) error {
			renderer := getRenderer(cmd.Context())
			versionInfo := &VersionInfo{
				Version:   Version,
				Commit:    GitCommit,
				BuildDate: BuildDate,
			}

			return renderer.Render(versionInfo, &options.RenderOptions)
		},
	}

	addRenderOptions(cmd, WithCardFormat(&options.RenderOptions))
	return cmd
}
