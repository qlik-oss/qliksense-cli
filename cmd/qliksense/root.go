package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	ansi "github.com/mattn/go-colorable"
	"github.com/mitchellh/go-homedir"
	"github.com/qlik-oss/sense-installer/pkg"
	"github.com/qlik-oss/sense-installer/pkg/api"
	"github.com/qlik-oss/sense-installer/pkg/qliksense"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ttacon/chalk"
)

// To run this project in debug mode, run:
// export QLIKSENSE_DEBUG=true
// qliksense <command>

const (
	qlikSenseHomeVar        = "QLIKSENSE_HOME"
	qlikSenseDirVar         = ".qliksense"
	keepPatchFilesFlagName  = "keep-config-repo-patches"
	keepPatchFilesFlagUsage = "Keep config repo patch files (for debugging)"
)

func initAndExecute() error {
	var (
		qlikSenseHome string
		err           error
	)
	qlikSenseHome, err = setUpPaths()
	if err != nil {
		log.Fatal(err)
	}
	// create dirs and appropriate files for setting up contexts
	api.LogDebugMessage("QliksenseHomeDir: %s", qlikSenseHome)

	qliksenseClient := qliksense.New(qlikSenseHome)
	cmd := rootCmd(qliksenseClient)
	if err := cmd.Execute(); err != nil {
		//levenstein checks (auto-suggestions)
		levenstein(cmd)
		return err
	}

	return nil
}

func setUpPaths() (string, error) {
	var (
		homeDir, qlikSenseHome string
		err                    error
	)

	if qlikSenseHome = os.Getenv(qlikSenseHomeVar); qlikSenseHome == "" {
		if homeDir, err = homedir.Dir(); err != nil {
			return "", err
		}
		if homeDir, err = homedir.Expand(homeDir); err != nil {
			return "", err
		}
		qlikSenseHome = filepath.Join(homeDir, qlikSenseDirVar)
	}

	if err := os.MkdirAll(qlikSenseHome, os.ModePerm); err != nil {
		return "", err
	}

	return qlikSenseHome, nil
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of qliksense cli",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s (%s, %s)\n", pkg.Version, pkg.Commit, pkg.CommitDate)
	},
}

func commandUsesContext(command string) bool {
	return command != "" && command != "help" && command != "version"
}

func globalPreRun(cmd *cobra.Command, p *qliksense.Qliksense) {
	if command := cmd.CalledAs(); commandUsesContext(command) {
		if isEulaEnforced() {
			enforceEula(p)
		}

		if err := p.SetUpQliksenseDefaultContext(); err != nil {
			panic(err)
		}

		if isEulaEnforced() {
			if err := p.SetEulaAccepted(); err != nil {
				panic(err)
			}
		}
	}
}

func rootCmd(p *qliksense.Qliksense) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "qliksense",
		Short: "Qliksense cli tool",
		Long:  `qliksense cli tool provides functionality to perform operations on qliksense-k8s, qliksense operator, and kubernetes cluster`,
		Args:  cobra.ArbitraryArgs,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			globalPreRun(cmd, p)
		},
	}

	cmd.Flags().SetInterspersed(false)

	cobra.OnInitialize(initConfig)

	// For qliksense overrides/commands

	cmd.AddCommand(getInstallableVersionsCmd(p))
	cmd.AddCommand(pullQliksenseImages(p))
	cmd.AddCommand(pushQliksenseImages(p))
	cmd.AddCommand(about(p))
	// add version command
	cmd.AddCommand(versionCmd)

	// add operator command
	cmd.AddCommand(operatorCmd)
	//operatorCmd.AddCommand(operatorViewCmd(p))
	operatorCmd.AddCommand(operatorCrdCmd(p))
	operatorCmd.AddCommand(operatorControllerCmd(p))

	//add fetch command
	cmd.AddCommand(fetchCmd(p))

	// add install command
	cmd.AddCommand(installCmd(p))

	// add config command
	configCmd := configCmd(p)
	cmd.AddCommand(configCmd)
	configCmd.AddCommand(configApplyCmd(p))
	configCmd.AddCommand(configViewCmd(p))

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	//add upgrade command
	cmd.AddCommand(upgradeCmd(p))

	// add the set-context config command as a sub-command to the app config command
	configCmd.AddCommand(setContextConfigCmd(p))

	// add the set profile/namespace/storageClassName/git-repository config command as a sub-command to the app config command
	configCmd.AddCommand(setOtherConfigsCmd(p))

	// add the set ### config command as a sub-command to the app config sub-command
	configCmd.AddCommand(setConfigsCmd(p))

	// add the set ### config command as a sub-command to the app config sub-command
	configCmd.AddCommand(setSecretsCmd(p))

	// add the list config command as a sub-command to the app config sub-command
	configCmd.AddCommand(listContextConfigCmd(p))

	// add the delete-context config command as a sub-command to the app config command
	configCmd.AddCommand(deleteContextConfigCmd(p))

	// add set-image-registry command as a sub-command to the app config sub-command
	configCmd.AddCommand(setImageRegistryCmd(p))

	// add clean-config-repo-patches command as a sub-command to the app config sub-command
	configCmd.AddCommand(cleanConfigRepoPatchesCmd(p))

	// add uninstall command
	cmd.AddCommand(uninstallCmd(p))

	// add crds
	cmd.AddCommand(crdsCmd)
	crdsCmd.AddCommand(crdsViewCmd(p))
	crdsCmd.AddCommand(crdsInstallCmd(p))

	// add preflight command
	preflightCmd := preflightCmd(p)
	preflightCmd.AddCommand(preflightCheckDnsCmd(p))
	preflightCmd.AddCommand(preflightCheckK8sVersionCmd(p))
	preflightCmd.AddCommand(preflightAllChecksCmd(p))
	//preflightCmd.AddCommand(preflightCheckMongoCmd(p))
	//preflightCmd.AddCommand(preflightCheckAllCmd(p))

	cmd.AddCommand(preflightCmd)
	cmd.AddCommand(loadCrFile(p))
	cmd.AddCommand((applyCmd(p)))
	return cmd
}

func initConfig() {
	viper.SetEnvPrefix("QLIKSENSE")
	viper.AutomaticEnv()
}

func copy(src, dst string) (int64, error) {
	var (
		source, destination *os.File
		sourceFileStat      os.FileInfo
		err                 error
		nBytes              int64
	)
	if sourceFileStat, err = os.Stat(src); err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	if source, err = os.Open(src); err != nil {
		return 0, err
	}
	defer source.Close()

	if destination, err = os.Create(dst); err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err = io.Copy(destination, source)
	return nBytes, err
}

func levenstein(cmd *cobra.Command) {
	cmd.SuggestionsMinimumDistance = 2
	if len(os.Args) > 1 {
		args := os.Args[1]
		suggest := cmd.SuggestionsFor(args)
		if len(suggest) > 0 {
			arg := []string{}
			for _, cm := range os.Args {
				arg = append(arg, cm)
			}
			if !strings.EqualFold(arg[1], suggest[0]) {
				arg[1] = suggest[0]
				out := ansi.NewColorableStdout()
				fmt.Fprintln(out, chalk.Green.Color("Did you mean: "), chalk.Bold.TextStyle(strings.Join(arg, " ")), "?")
			}
		}
	}
}
