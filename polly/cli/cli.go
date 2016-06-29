package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/akutz/gofig"
	glog "github.com/akutz/golf/logrus"
	"github.com/akutz/gotil"
	// config "github.com/emccode/polly/core/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v1"

	"github.com/akutz/goof"
	// "github.com/emccode/polly/core"
	"github.com/emccode/polly/client"
	"github.com/emccode/polly/core"
	"github.com/emccode/polly/core/config"
	"github.com/emccode/polly/core/store"
	"github.com/emccode/polly/core/types"
	"github.com/emccode/polly/polly/cli/term"
	"github.com/emccode/polly/util"
)

func init() {
	log.SetFormatter(&glog.TextFormatter{log.TextFormatter{}})
}

type helpFlagPanic struct{}
type printedErrorPanic struct{}
type subCommandPanic struct{}

// CLI is the Polly command line interface.
type CLI struct {
	l  *log.Logger
	p  *types.Polly
	c  *cobra.Command
	pc client.Client

	serviceCmd           *cobra.Command
	versionCmd           *cobra.Command
	envCmd               *cobra.Command
	volumeCmd            *cobra.Command
	installCmd           *cobra.Command
	uninstallCmd         *cobra.Command
	serviceStartCmd      *cobra.Command
	serviceRestartCmd    *cobra.Command
	serviceStopCmd       *cobra.Command
	serviceStatusCmd     *cobra.Command
	serviceInitSysCmd    *cobra.Command
	volumeGetCmd         *cobra.Command
	volumeOfferCmd       *cobra.Command
	volumeOfferRevokeCmd *cobra.Command
	volumeLabelCmd       *cobra.Command
	volumeLabelRemoveCmd *cobra.Command
	volumeCreateCmd      *cobra.Command
	volumeRemoveCmd      *cobra.Command
	storeCmd             *cobra.Command
	storeEraseCmd        *cobra.Command
	storeGetCmd          *cobra.Command

	outputFormat     string
	client           string
	fg               bool
	force            bool
	cfgFile          string
	all              bool
	volumeID         string
	schedulers       []string
	labels           []string
	serviceName      string
	volumeType       string
	IOPS             int64
	size             int64
	name             string
	availabilityZone string
}

const (
	noColor     = 0
	black       = 30
	red         = 31
	redBg       = 41
	green       = 32
	yellow      = 33
	blue        = 34
	gray        = 37
	blueBg      = blue + 10
	white       = 97
	whiteBg     = white + 10
	darkGrayBg  = 100
	lightBlue   = 94
	lightBlueBg = lightBlue + 10
)

func validateConfig(path string) {
	if !gotil.FileExists(path) {
		return
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(
			os.Stderr, "polly: error reading config: %s\n%v\n", path, err)
		os.Exit(1)
	}

	s := string(buf)

	if _, err := gofig.ValidateYAMLString(s); err != nil {
		fmt.Fprintf(
			os.Stderr,
			"polly: invalid config: %s\n\n  %v\n\n", path, err)
		fmt.Fprint(
			os.Stderr,
			"paste the contents between ---BEGIN--- and ---END---\n")
		fmt.Fprint(
			os.Stderr,
			"into http://www.yamllint.com/ to discover the issue\n\n")
		fmt.Fprintln(os.Stderr, "---BEGIN---")
		fmt.Fprintln(os.Stderr, s)
		fmt.Fprintln(os.Stderr, "---END---")
		os.Exit(1)
	}
}

// New returns a new CLI using the current process's arguments.
func New() *CLI {
	return NewWithArgs(os.Args[1:]...)
}

// NewWithArgs returns a new CLI using the specified arguments.
func NewWithArgs(a ...string) *CLI {

	validateConfig(util.EtcFilePath("config.yml"))
	validateConfig(fmt.Sprintf("%s/.polly/config.yml", gotil.HomeDir()))

	s := "Polly:\n" +
		"  Polly-morphic storage scheduling"

	cfg, err := config.New()
	if err != nil {
		log.Error(goof.WithError("problem getting config", err))
		os.Exit(1)
	}

	pc, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	p := core.NewWithConfig(cfg)

	ps, err := store.NewWithConfig(p.Config.Scope("polly.store"))
	if err != nil {
		log.Error(goof.WithError("problem initialization store", err))
		os.Exit(1)
	}
	p.Store = ps

	c := &CLI{
		l:  log.New(),
		p:  p,
		pc: pc,
	}

	c.c = &cobra.Command{
		Use:              "polly",
		Short:            s,
		PersistentPreRun: c.preRun,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	c.c.SetArgs(a)

	c.initOtherCmdsAndFlags()
	c.initVolumeCmdsAndFlags()
	c.initStoreCmdsAndFlags()
	c.initServiceCmdsAndFlags()
	c.initUsageTemplates()

	return c
}

// Execute executes the CLI using the current process's arguments.
func Execute() {
	New().Execute()
}

// ExecuteWithArgs executes the CLI using the specified arguments.
func ExecuteWithArgs(a ...string) {
	NewWithArgs(a...).Execute()
}

// Execute executes the CLI.
func (c *CLI) Execute() {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		switch r := r.(type) {
		case int:
			log.Debugf("exiting with error code %d", r)
			os.Exit(r)
		case error:
			log.Panic(r)
		default:
			log.Debugf("exiting1 with default error code 1, r=%v", r)
			os.Exit(1)
		}
	}()
	c.execute()
}

func (c *CLI) execute() {
	defer func() {
		r := recover()
		if r != nil {
			switch r.(type) {
			case helpFlagPanic, subCommandPanic:
			// Do nothing
			case printedErrorPanic:
				log.Panic(r)
				os.Exit(1)
			default:
				log.Debugf("exiting2 with default error code 1, r=%v", r)
				panic(r)
			}
		}
	}()
	c.c.Execute()
}

func (c *CLI) marshalOutput(v interface{}) (string, error) {
	var err error
	var buf []byte
	if strings.ToUpper(c.outputFormat) == "JSON" {
		buf, err = marshalJSONOutput(v)
	} else {
		buf, err = marshalYamlOutput(v)
	}
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func marshalYamlOutput(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func marshalJSONOutput(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (c *CLI) addOutputFormatFlag(fs *pflag.FlagSet) {
	fs.StringVarP(
		&c.outputFormat, "format", "f", "yml", "The output format (yml, json)")
}

func (c *CLI) updateLogLevel() {
	switch c.logLevel() {
	case "panic":
		log.SetLevel(log.PanicLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	}

	log.WithField("logLevel", c.logLevel()).Debug("updated log level")
}

func (c *CLI) preRun(cmd *cobra.Command, args []string) {

	if c.cfgFile != "" && gotil.FileExists(c.cfgFile) {
		validateConfig(c.cfgFile)
		if err := c.p.Config.ReadConfigFile(c.cfgFile); err != nil {
			panic(err)
		}
		cmd.Flags().Parse(os.Args[1:])
	}

	c.updateLogLevel()

	if isHelpFlag(cmd) {
		cmd.Help()
		panic(&helpFlagPanic{})
	}

	if permErr := c.checkCmdPermRequirements(cmd); permErr != nil {
		if term.IsTerminal() {
			printColorizedError(permErr)
		} else {
			printNonColorizedError(permErr)
		}

		fmt.Println()
		cmd.Help()
		panic(&printedErrorPanic{})
	}

}

func isHelpFlags(cmd *cobra.Command) bool {
	help, _ := cmd.Flags().GetBool("help")
	verb, _ := cmd.Flags().GetBool("verbose")
	return help || verb
}

func (c *CLI) checkCmdPermRequirements(cmd *cobra.Command) error {
	if cmd == c.installCmd {
		return checkOpPerms("installed")
	}

	if cmd == c.uninstallCmd {
		return checkOpPerms("uninstalled")
	}

	if cmd == c.serviceStartCmd {
		return checkOpPerms("started")
	}

	if cmd == c.serviceStopCmd {
		return checkOpPerms("stopped")
	}

	if cmd == c.serviceRestartCmd {
		return checkOpPerms("restarted")
	}

	return nil
}

func printColorizedError(err error) {
	stderr := os.Stderr
	l := fmt.Sprintf("\x1b[%dm\xe2\x86\x93\x1b[0m", white)

	fmt.Fprintf(stderr, "Oops, an \x1b[%[1]dmerror\x1b[0m occured!\n\n", redBg)
	fmt.Fprintf(stderr, "  \x1b[%dm%s\n\n", red, err.Error())
	fmt.Fprintf(stderr, "\x1b[0m")
	fmt.Fprintf(stderr,
		"To correct the \x1b[%dmerror\x1b[0m please review:\n\n", redBg)
	fmt.Fprintf(
		stderr,
		"  - Debug output by using the flag \x1b[%dm-l debug\x1b[0m\n",
		lightBlue)
	fmt.Fprintf(stderr, "  - The Polly website at \x1b[%dm%s\x1b[0m\n",
		blueBg, "https://github.com/emccode/polly")
	fmt.Fprintf(stderr, "  - The on%[1]sine he%[1]sp be%[1]sow\n", l)
}

func printNonColorizedError(err error) {
	stderr := os.Stderr

	fmt.Fprintf(stderr, "Oops, an error occured!\n\n")
	fmt.Fprintf(stderr, "  %s\n", err.Error())
	fmt.Fprintf(stderr, "To correct the error please review:\n\n")
	fmt.Fprintf(stderr, "  - Debug output by using the flag \"-l debug\"\n")
	fmt.Fprintf(
		stderr,
		"  - The Polly website at https://github.com/emccode/polly\n")
	fmt.Fprintf(stderr, "  - The online help below\n")
}

func (c *CLI) logLevel() string {
	return c.p.Config.GetString("polly.logLevel")
}

func (c *CLI) host() string {
	return c.p.Config.GetString("polly.host")
}
