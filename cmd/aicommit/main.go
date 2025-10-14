package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aicommit/aicommit/internal/config"
	"github.com/aicommit/aicommit/internal/git"
	"github.com/aicommit/aicommit/internal/model"
	"github.com/aicommit/aicommit/pkg/prompt"
	"github.com/aicommit/aicommit/pkg/validator"
	"github.com/spf13/cobra"
)

var (
	version = "1.0.0"
	cfgFile string
	dryRun  bool
)

func main() {
	// 示例：你可以在这里自定义prompt
	// SetChinesePrompt()  // 使用中文prompt
	// SetDetailedPrompt() // 使用详细prompt
	// SetMinimalPrompt()  // 使用简洁prompt

	rootCmd := &cobra.Command{
		Use:   "aicommit",
		Short: "AI-powered git commit message generator",
		Long:  "aicommit uses AI models to generate meaningful commit messages based on your staged changes",
		RunE:  run,
	}

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.config/aicommit/aicommit.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "show the generated commit message without committing")

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("aicommit version %s\n", version)
		},
	}

	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
	}

	configInitCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize configuration file",
		RunE:  initConfig,
	}

	// Prompt customization commands
	promptCmd := &cobra.Command{
		Use:   "prompt",
		Short: "Manage prompt templates",
	}

	promptDefaultCmd := &cobra.Command{
		Use:   "default",
		Short: "Use default prompt template",
		Run: func(cmd *cobra.Command, args []string) {
			SetDefaultPrompt()
			fmt.Println("Switched to default prompt template")
		},
	}

	promptChineseCmd := &cobra.Command{
		Use:   "chinese",
		Short: "Use Chinese prompt template",
		Run: func(cmd *cobra.Command, args []string) {
			SetChinesePrompt()
			fmt.Println("Switched to Chinese prompt template")
		},
	}

	promptDetailedCmd := &cobra.Command{
		Use:   "detailed",
		Short: "Use detailed prompt template",
		Run: func(cmd *cobra.Command, args []string) {
			SetDetailedPrompt()
			fmt.Println("Switched to detailed prompt template")
		},
	}

	promptMinimalCmd := &cobra.Command{
		Use:   "minimal",
		Short: "Use minimal prompt template",
		Run: func(cmd *cobra.Command, args []string) {
			SetMinimalPrompt()
			fmt.Println("Switched to minimal prompt template")
		},
	}

	promptCmd.AddCommand(promptDefaultCmd, promptChineseCmd, promptDetailedCmd, promptMinimalCmd)

	configCmd.AddCommand(configInitCmd)
	rootCmd.AddCommand(versionCmd, configCmd, promptCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func run(cmd *cobra.Command, args []string) error {
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

	diff, err := gitClient.GetDiff()
	if err != nil {
		return fmt.Errorf("failed to get diff: %w", err)
	}

	provider, err := model.NewProvider(cfg)
	if err != nil {
		return fmt.Errorf("failed to create provider: %w", err)
	}

	ctx := context.Background()
	fmt.Printf("Generating commit message using %s with model %s...\n", provider.Name(), cfg.Model)

	commitMessage, err := provider.GenerateCommitMessage(ctx, diff)
	if err != nil {
		return fmt.Errorf("failed to generate commit message: %w", err)
	}

	commitMessage = prompt.CleanCommitMessage(commitMessage)

	if err := prompt.ValidateCommitMessage(commitMessage); err != nil {
		return fmt.Errorf("generated commit message is invalid: %w", err)
	}

	fmt.Printf("\nGenerated commit message:\n%s\n", commitMessage)

	if dryRun {
		fmt.Println("\nDry run mode - no commit was made")
		return nil
	}

	if err := gitClient.Commit(commitMessage); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	fmt.Println("\nCommit successful!")
	return nil
}

func initConfig(cmd *cobra.Command, args []string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := fmt.Sprintf("%s/.config/aicommit", homeDir)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configFile := fmt.Sprintf("%s/aicommit.yaml", configDir)
	if _, err := os.Stat(configFile); err == nil {
		return fmt.Errorf("config file already exists: %s", configFile)
	}

	defaultConfig := `# aicommit configuration file
model: claude-3-sonnet-20240229
provider: claude

# API keys - you can also use environment variables:
# AICOMMIT_CLAUDE_API_KEY, AICOMMIT_OPENAI_API_KEY, AICOMMIT_DEEPSEEK_API_KEY
api_keys:
  claude: ""    # Your Claude API key
  openai: ""    # Your OpenAI API key
  deepseek: ""  # Your DeepSeek API key
`

	if err := os.WriteFile(configFile, []byte(defaultConfig), 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("Configuration file created: %s\n", configFile)
	fmt.Println("Please edit the file to add your API keys")
	return nil
}
