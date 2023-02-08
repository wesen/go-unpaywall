package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/help"
	"github.com/spf13/cobra"
	"github.com/wesen/go-unpaywall/pkg"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "unpaywall",
	Short: "unpaywall is a tool to manage bibliography data",
}

var lookupCmd = &cobra.Command{
	Use:   "lookup",
	Short: "Lookup a bibliography entry",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		email, _ := cmd.Flags().GetString("email")
		if email == "" {
			// lookup email in environment under UNPAYWALL_EMAIL
			email = os.Getenv("UNPAYWALL_EMAIL")
		}
		if email == "" {
			_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "email is required")
			os.Exit(1)
		}
		baseUrl, _ := cmd.Flags().GetString("base-url")
		if baseUrl == "" {
			_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "base-url is required")
			os.Exit(1)
		}

		c := pkg.NewClient(pkg.WithEmail(email), pkg.WithBaseURL(baseUrl))

		ctx := context.Background()

		gp, of, err := cli.SetupProcessor(cmd)
		cobra.CheckErr(err)

		for _, arg := range args {
			entry, err := c.GetDOI(ctx, arg)
			cobra.CheckErr(err)

			// serialize to json string
			s, err := json.Marshal(entry)
			cobra.CheckErr(err)
			// deserialize to map
			var m map[string]interface{}
			err = json.Unmarshal(s, &m)
			cobra.CheckErr(err)

			delete(m, "oa_locations")
			delete(m, "first_oa_location")
			if _, ok := m["best_oa_location"]; ok {
				m["location"] = m["best_oa_location"]
				delete(m, "best_oa_location")
			}

			err = gp.ProcessInputObject(m)
			cobra.CheckErr(err)
		}

		s, err := of.Output()
		cobra.CheckErr(err)

		fmt.Print(s)
	},
}

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for bibliography entries",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		email, _ := cmd.Flags().GetString("email")
		if email == "" {
			// lookup email in environment under UNPAYWALL_EMAIL
			email = os.Getenv("UNPAYWALL_EMAIL")
		}
		if email == "" {
			_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "email is required")
			os.Exit(1)
		}
		baseUrl, _ := cmd.Flags().GetString("base-url")
		if baseUrl == "" {
			_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "base-url is required")
			os.Exit(1)
		}
		req := pkg.SearchRequest{
			Query: args[0],
		}
		onlyOpenAccess, _ := cmd.Flags().GetBool("only-open-access")
		if onlyOpenAccess {
			req.IsOA = &onlyOpenAccess
		}
		onlyNonOpenAccess, _ := cmd.Flags().GetBool("only-non-open-access")
		if onlyNonOpenAccess {
			onlyOpenAccess := false
			req.IsOA = &onlyOpenAccess
		}

		page, _ := cmd.Flags().GetInt("page")
		if page > 0 {
			req.Page = &page
		}

		c := pkg.NewClient(pkg.WithEmail(email), pkg.WithBaseURL(baseUrl))

		ctx := context.Background()

		gp, of, err := cli.SetupProcessor(cmd)
		cobra.CheckErr(err)

		results, err := c.Search(ctx, req)
		cobra.CheckErr(err)

		for _, entry := range results {
			// serialize to json string
			s, err := json.Marshal(entry.Response)
			cobra.CheckErr(err)
			// deserialize to map
			var m map[string]interface{}
			err = json.Unmarshal(s, &m)
			cobra.CheckErr(err)

			m["score"] = entry.Score
			m["snippet"] = entry.Snippet

			delete(m, "oa_locations")
			delete(m, "first_oa_location")
			if _, ok := m["best_oa_location"]; ok {
				m["location"] = m["best_oa_location"]
				delete(m, "best_oa_location")
			}

			err = gp.ProcessInputObject(m)
			cobra.CheckErr(err)

		}

		s, err := of.Output()
		cobra.CheckErr(err)

		fmt.Print(s)
	},
}

func main() {
	_ = rootCmd.Execute()
}

//go:embed doc/*
var docFS embed.FS

func init() {
	helpSystem := help.NewHelpSystem()
	err := helpSystem.LoadSectionsFromFS(docFS, ".")
	if err != nil {
		panic(err)
	}

	helpFunc, usageFunc := help.GetCobraHelpUsageFuncs(helpSystem)
	helpTemplate, usageTemplate := help.GetCobraHelpUsageTemplates(helpSystem)

	_ = usageFunc
	_ = usageTemplate

	rootCmd.SetHelpFunc(helpFunc)
	rootCmd.SetUsageFunc(usageFunc)
	rootCmd.SetHelpTemplate(helpTemplate)
	rootCmd.SetUsageTemplate(usageTemplate)

	helpCmd := help.NewCobraHelpCommand(helpSystem)
	rootCmd.SetHelpCommand(helpCmd)

	rootCmd.PersistentFlags().String("email", "", "email address to use for API requests (required)")
	rootCmd.PersistentFlags().String("base-url", "https://api.unpaywall.org", "base URL for API requests")

	// TODO(manuel, 2023-02-02): Provide better defaults for the output fields
	cli.AddFlags(lookupCmd, cli.NewFlagsDefaults())
	rootCmd.AddCommand(lookupCmd)

	// TODO(manuel, 2023-02-02): Provide better defaults for the output fields
	cli.AddFlags(searchCmd, cli.NewFlagsDefaults())
	searchCmd.Flags().Int("page", 1, "page number to return")
	searchCmd.Flags().Bool("only-open-access", false, "only return open access entries")
	searchCmd.Flags().Bool("only-non-open-access", false, "only return non-open access entries")
	rootCmd.AddCommand(searchCmd)
}
