package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/superfly/flyctl/internal/version"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "version",
		Short: "Tool for working with flyctl version numbers",
	}

	rootCmd.AddCommand(newLatestVersionCmd())
	rootCmd.AddCommand(newNextVersionCmd())

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func newLatestVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "latest",
		Short: "Prints the latest version for the current track",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := refreshTags(); err != nil {
				return err
			}
			cmd.PrintErrln("refreshed tags")

			ref, err := gitRef()
			if err != nil {
				return err
			}
			cmd.PrintErrln("ref:", ref)

			time, err := gitCommitTime(ref)
			if err != nil {
				return err
			}
			cmd.PrintErrln("commit time:", time)

			track, err := trackFromRef(ref)
			if err != nil {
				return err
			}
			cmd.PrintErrln("track:", track)

			currentVersion, err := latestVersion(track)
			if err != nil {
				return err
			}
			cmd.PrintErrln("current version:", currentVersion)

			cmd.Print(currentVersion)

			return nil
		},
	}

	return cmd
}

func newNextVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "next",
		Short: "Prints the next version for the current track",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := refreshTags(); err != nil {
				return err
			}
			cmd.PrintErrln("refreshed tags")

			ref, err := gitRef()
			if err != nil {
				return err
			}
			cmd.PrintErrln("ref:", ref)

			time, err := gitCommitTime(ref)
			if err != nil {
				return err
			}
			cmd.PrintErrln("commit time:", time)

			track, err := trackFromRef(ref)
			if err != nil {
				return err
			}
			cmd.PrintErrln("track:", track)

			currentVersion, err := latestVersion(track)
			if err != nil {
				return err
			}
			cmd.PrintErrln("current version:", currentVersion)

			buildNumber, err := nextBuildNumber(track, time)
			if err != nil {
				return err
			}
			cmd.PrintErrln("build number:", buildNumber)

			nextVersion := version.Version{
				Major: time.Year(),
				Minor: int(time.Month()),
				Patch: time.Day(),
				Track: track,
				Build: buildNumber,
			}

			cmd.Print(nextVersion)

			return nil
		},
	}

	return cmd
}
