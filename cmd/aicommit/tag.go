package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aicommit/aicommit/internal/config"
	"github.com/aicommit/aicommit/internal/git"
	"github.com/aicommit/aicommit/internal/model"
	"github.com/aicommit/aicommit/pkg/editor"
	"github.com/aicommit/aicommit/pkg/prompt"
	"github.com/aicommit/aicommit/pkg/validator"
	"github.com/spf13/cobra"
)

const (
	defaultTagCommitLimit      = 50
	defaultTagDiffStatMaxLen   = 4000
	defaultTagNameStatusMaxLen = 4000
)

func newTagCmd() *cobra.Command {
	var versionFlag string

	cmd := &cobra.Command{
		Use:   "tag [version]",
		Short: "Generate an annotated git tag message with AI and create the tag",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			version := ""
			if len(args) > 0 {
				version = args[0]
			}
			if strings.TrimSpace(version) == "" {
				version = versionFlag
			}
			if strings.TrimSpace(version) == "" {
				var err error
				version, err = promptForVersion(cmd)
				if err != nil {
					return err
				}
			}
			version = strings.TrimSpace(version)
			if version == "" {
				return fmt.Errorf("tag version cannot be empty")
			}

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			if err := validator.ValidateRepository("."); err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}

			gitClient := git.New(".")
			if !gitClient.IsRepository() {
				return fmt.Errorf("not a git repository")
			}

			exists, err := gitClient.TagExists(version)
			if err != nil {
				return fmt.Errorf("failed to check tag existence: %w", err)
			}
			if exists {
				return fmt.Errorf("tag already exists: %s", version)
			}

			previousTag, hasPreviousTag, err := gitClient.LatestTag()
			if err != nil {
				return err
			}

			rangeSpec := ""
			if hasPreviousTag {
				rangeSpec = fmt.Sprintf("%s..HEAD", previousTag)
			}

			commitRange := rangeSpec
			if commitRange == "" {
				commitRange = ""
			}

			commitSubjects, truncated, err := gitClient.CommitSubjects(commitRange, defaultTagCommitLimit)
			if err != nil {
				return fmt.Errorf("failed to get commit subjects: %w", err)
			}

			diffStat := "unavailable (no previous tag found)"
			nameStatus := "unavailable (no previous tag found)"
			if hasPreviousTag {
				if s, err := gitClient.DiffStat(rangeSpec); err == nil {
					diffStat = truncateText(strings.TrimSpace(s), defaultTagDiffStatMaxLen)
				} else {
					diffStat = fmt.Sprintf("unavailable (%v)", err)
				}
				if s, err := gitClient.DiffNameStatus(rangeSpec); err == nil {
					nameStatus = truncateText(strings.TrimSpace(s), defaultTagNameStatusMaxLen)
				} else {
					nameStatus = fmt.Sprintf("unavailable (%v)", err)
				}
			}

			infoBlock := buildTagInfoBlock(version, previousTag, hasPreviousTag, rangeSpec, commitSubjects, truncated, diffStat, nameStatus)

			provider, err := model.NewProvider(cfg)
			if err != nil {
				return fmt.Errorf("failed to create provider: %w", err)
			}
			provider.SetTemplate(prompt.NewTagTemplate())

			ctx := context.Background()

			modelName := cfg.Model
			if cfg.Provider == "custom" && cfg.Custom.Model != "" {
				modelName = cfg.Custom.Model
			}

			fmt.Printf("Generating tag message using %s with model %s...\n", provider.Name(), modelName)

			tagMessage, err := provider.GenerateMessage(ctx, infoBlock)
			if err != nil {
				return fmt.Errorf("failed to generate tag message: %w", err)
			}

			tagMessage = prompt.CleanAIText(tagMessage)
			if strings.TrimSpace(tagMessage) == "" {
				return fmt.Errorf("generated tag message is empty")
			}

			fmt.Printf("\nGenerated tag message:\n%s\n", tagMessage)

			fmt.Println("\nOpening editor to review/edit tag message...")
			edited, err := editor.Open(tagMessage, cfg.Editor)
			if err != nil {
				return fmt.Errorf("failed to open editor: %w", err)
			}

			edited = strings.TrimSpace(edited)
			if edited == "" {
				fmt.Println("\nTag message is empty, aborting tag creation.")
				return nil
			}

			if dryRun {
				fmt.Println("\nDry run mode - no tag was created")
				return nil
			}

			if err := gitClient.CreateAnnotatedTag(version, edited); err != nil {
				return fmt.Errorf("failed to create tag: %w", err)
			}

			fmt.Printf("\nTag created: %s\n", version)
			if !hasPreviousTag {
				fmt.Println("Note: no previous tag was found; consider creating an initial baseline tag for better release notes.")
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&versionFlag, "version", "", "tag version (if not provided as an argument)")
	return cmd
}

func promptForVersion(cmd *cobra.Command) (string, error) {
	fmt.Fprint(cmd.OutOrStdout(), "Tag version: ")
	reader := bufio.NewReader(cmd.InOrStdin())
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read tag version: %w", err)
	}
	return strings.TrimSpace(line), nil
}

func buildTagInfoBlock(
	version string,
	previousTag string,
	hasPreviousTag bool,
	rangeSpec string,
	commitSubjects []string,
	truncated bool,
	diffStat string,
	nameStatus string,
) string {
	var b strings.Builder
	b.WriteString("Release version: ")
	b.WriteString(version)
	b.WriteString("\n")

	if hasPreviousTag {
		b.WriteString("Previous tag: ")
		b.WriteString(previousTag)
		b.WriteString("\n")
		b.WriteString("Range: ")
		b.WriteString(rangeSpec)
		b.WriteString("\n")
	} else {
		b.WriteString("Previous tag: (none)\n")
		b.WriteString("Range: (no previous tag; summary is for repository history up to HEAD)\n")
	}

	b.WriteString("\nCommit subjects:\n")
	if len(commitSubjects) == 0 {
		b.WriteString("(none)\n")
	} else {
		for _, s := range commitSubjects {
			s = strings.TrimSpace(s)
			if s == "" {
				continue
			}
			b.WriteString("- ")
			b.WriteString(s)
			b.WriteString("\n")
		}
		if truncated {
			b.WriteString("(commit list truncated)\n")
		}
	}

	b.WriteString("\nDiffstat:\n")
	b.WriteString(diffStat)
	b.WriteString("\n")

	b.WriteString("\nChanged files (name-status):\n")
	b.WriteString(nameStatus)
	b.WriteString("\n")

	return b.String()
}

func truncateText(s string, max int) string {
	if max <= 0 || len(s) <= max {
		return s
	}
	return s[:max] + "\n... (truncated)"
}
