package image

import (
	"github.com/bthuilot/dockerleaks/internal/util"
	"github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"strings"
)

// runLineRegex is the regular expression to match a RUN line in
// a Dockerfile
var runLineRegex = regexp.MustCompile(`^(?:RUN )?\|(\d)+`)

// shellLineRegex is the regular expression to match a SHELL line in
// a Dockerfile
var shellLineRegex = regexp.MustCompile(`SHELL\s+\[\s*"?([^]"])"?`)

// argLineRegex is the regular expression to match a ARG line in
// a Dockerfile
var argLineRegex = regexp.MustCompile(`ARG\s+([A-Za-z0-9_\-]+)`)

// ParseBuildArguments will parse out each build argument by inspecting
// each response item in the docker images history, collecting the current set shell
// and build arguments, to be able to parse out each build arguments value
func (i image) ParseBuildArguments() ([]BuildArg, error) {
	history, err := i.cli.ImageHistory(i.ctx, i.ref.String())
	var (
		// buildArgs is the list of build arguments discovered
		buildArgs []BuildArg
		// definedArgs is the list of build arguments defined in the current context,
		// this is needed to be able to parse out the build arguments value from a run command
		definedArgs []string
		// shell is the current shell being used in the Dockerfile
		// this is needed to be able to parse out the build arguments value from a run command
		shell = "/bin/sh"
	)
	for _, h := range util.Reverse(history) {
		logrus.Debugf("parsing history line %s", h.CreatedBy)
		// If this is a shell line
		if matches := shellLineRegex.FindStringSubmatch(h.CreatedBy); matches != nil {
			// set the shell to the new shell
			logrus.Debugf("found match for SHELL line: %s", matches[1])
			shell = matches[1]
		}

		if matches := argLineRegex.FindStringSubmatch(h.CreatedBy); matches != nil {
			// If this is a build arg line
			// add the build arg to the current list of args
			logrus.Debugf("found match for ARG line: %s", matches[1])
			definedArgs = append(definedArgs, matches[1])
		}

		if matches := runLineRegex.FindStringSubmatch(h.CreatedBy); matches != nil {
			amt := 0
			if amt, err = strconv.Atoi(matches[1]); err != nil {
				logrus.Warnf("invalid build arg amount %s, skipping", matches[1])
			} else if len(definedArgs) != amt {
				logrus.Warnf("amount of counted args %d, differs from the amount of build args %d\n", len(definedArgs), amt)
			}
			regxp, err := regexp.Compile(strings.Join(append(definedArgs, shell), `=(.*)\s`))
			if err != nil {
				logrus.Errorf("unable to compile regex %s\n", err)
			}
			if matches = regxp.FindStringSubmatch(h.CreatedBy); matches != nil {
				// TODO(i dont love this)
				buildArgs = append(buildArgs, util.ZipApply(func(name string, value string) BuildArg {
					return BuildArg{
						Name:     name,
						Value:    value,
						Location: h.CreatedBy,
					}
				}, definedArgs, matches[1:])...)
			} else {
				logrus.Debugf("RUN regexp %s failed to match", regxp.String())
			}
		}
	}

	// remove duplicates
	buildArgs = uniqueBuildArgs(buildArgs)

	return buildArgs, err
}

func uniqueBuildArgs(buildArgs []BuildArg) (unique []BuildArg) {
	seen := make(map[string][]string)
	for _, arg := range buildArgs {
		if _, ok := seen[arg.Name]; !ok {
			unique = append(unique, arg)
			seen[arg.Name] = []string{arg.Value}
		} else if !util.Any(seen[arg.Name], func(s string) bool {
			return s == arg.Value
		}) {
			seen[arg.Name] = append(seen[arg.Name], arg.Value)
			unique = append(unique, arg)
		}
	}
	return
}
