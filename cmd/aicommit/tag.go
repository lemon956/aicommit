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
			return runTag(cmd, args, versionFlag)
		},
	}

	cmd.Flags().StringVar(&versionFlag, "version", "", "tag version (if not provided as an argument)")
	return cmd
}

func runTag(cmd *cobra.Command, args []string, versionFlag string) error {
	version, err := resolveTagVersion(cmd, args, versionFlag)
	if err != nil {
		return err
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	gitClient, err := mustOpenRepo()
	if err != nil {
		return err
	}

	if err := ensureTagDoesNotExist(gitClient, version); err != nil {
		return err
	}

	infoBlock, hasPreviousTag, err := buildTagContext(gitClient, version)
	if err != nil {
		return err
	}

	tagMessage, err := generateTagMessage(cfg, infoBlock)
	if err != nil {
		return err
	}

	edited, aborted, err := reviewTagMessage(cmd, cfg, tagMessage)
	if err != nil {
		return err
	}
	if aborted {
		return nil
	}

	if dryRun {
		fmt.Fprintln(cmd.OutOrStdout(), "\nDry run mode - no tag was created")
		return nil
	}

	if err := gitClient.CreateAnnotatedTag(version, edited); err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "\nTag created: %s\n", version)
	if !hasPreviousTag {
		fmt.Fprintln(cmd.OutOrStdout(), "Note: no previous tag was found; consider creating an initial baseline tag for better release notes.")
	}
	return nil
}

func resolveTagVersion(cmd *cobra.Command, args []string, versionFlag string) (string, error) {
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
			return "", err
		}
	}
	version = strings.TrimSpace(version)
	if version == "" {
		return "", fmt.Errorf("tag version cannot be empty")
	}
	return version, nil
}

func mustOpenRepo() (*git.Git, error) {
	if err := validator.ValidateRepository("."); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	gitClient := git.New(".")
	if !gitClient.IsRepository() {
		return nil, fmt.Errorf("not a git repository")
	}

	return gitClient, nil
}

func ensureTagDoesNotExist(gitClient *git.Git, tag string) error {
	exists, err := gitClient.TagExists(tag)
	if err != nil {
		return fmt.Errorf("failed to check tag existence: %w", err)
	}
	if exists {
		return fmt.Errorf("tag already exists: %s", tag)
	}
	return nil
}

func buildTagContext(gitClient *git.Git, version string) (infoBlock string, hasPreviousTag bool, err error) {
	previousTag, hasPreviousTag, err := gitClient.LatestTag()
	if err != nil {
		return "", false, err
	}

	rangeSpec := ""
	if hasPreviousTag {
		rangeSpec = fmt.Sprintf("%s..HEAD", previousTag)
	}

	commitSubjects, truncated, err := gitClient.CommitSubjects(rangeSpec, defaultTagCommitLimit)
	if err != nil {
		return "", false, fmt.Errorf("failed to get commit subjects: %w", err)
	}

	diffStat := "unavailable (no previous tag found)"
	nameStatus := "unavailable (no previous tag found)"
	if hasPreviousTag {
		diffStat = formatOrUnavailable(func() (string, error) { return gitClient.DiffStat(rangeSpec) }, defaultTagDiffStatMaxLen)
		nameStatus = formatOrUnavailable(func() (string, error) { return gitClient.DiffNameStatus(rangeSpec) }, defaultTagNameStatusMaxLen)
	}

	infoBlock = buildTagInfoBlock(version, previousTag, hasPreviousTag, rangeSpec, commitSubjects, truncated, diffStat, nameStatus)
	return infoBlock, hasPreviousTag, nil
}

func formatOrUnavailable(fetch func() (string, error), maxLen int) string {
	s, err := fetch()
	if err != nil {
		return fmt.Sprintf("unavailable (%v)", err)
	}
	return truncateText(strings.TrimSpace(s), maxLen)
}

func generateTagMessage(cfg *config.Config, infoBlock string) (string, error) {
	provider, err := model.NewProvider(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to create provider: %w", err)
	}
	provider.SetTemplate(prompt.NewTagTemplate())

	modelName := cfg.Model
	if cfg.Provider == "custom" && cfg.Custom.Model != "" {
		modelName = cfg.Custom.Model
	}

	fmt.Printf("Generating tag message using %s with model %s...\n", provider.Name(), modelName)

	tagMessage, err := provider.GenerateMessage(context.Background(), infoBlock)
	if err != nil {
		return "", fmt.Errorf("failed to generate tag message: %w", err)
	}

	tagMessage = prompt.CleanAIText(tagMessage)
	if strings.TrimSpace(tagMessage) == "" {
		return "", fmt.Errorf("generated tag message is empty")
	}

	return tagMessage, nil
}

func reviewTagMessage(cmd *cobra.Command, cfg *config.Config, tagMessage string) (edited string, aborted bool, err error) {
	fmt.Fprintf(cmd.OutOrStdout(), "\nGenerated tag message:\n%s\n", tagMessage)

	fmt.Fprintln(cmd.OutOrStdout(), "\nOpening editor to review/edit tag message...")
	edited, err = editor.Open(tagMessage, cfg.Editor)
	if err != nil {
		return "", false, fmt.Errorf("failed to open editor: %w", err)
	}

	edited = strings.TrimSpace(edited)
	if edited == "" {
		fmt.Fprintln(cmd.OutOrStdout(), "\nTag message is empty, aborting tag creation.")
		return "", true, nil
	}

	return edited, false, nil
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
