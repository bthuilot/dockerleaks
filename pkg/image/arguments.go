package image

import (
	"context"
	"github.com/bthuilot/dockerleaks/internal/util"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"strings"
)

var runLineRegex = regexp.MustCompile(`^(?:RUN )?\|(\d)+`)
var shellLineRegex = regexp.MustCompile(`SHELL\s+\[\s*"?([^]"])"?`)
var argLineRegex = regexp.MustCompile(`ARG\s+([A-Za-z0-9_\-]+)`)

func parseBuildArgs(cli *client.Client, ctx context.Context, ref string) ([]BuildArg, error) {
	history, err := cli.ImageHistory(ctx, ref)
	var (
		envVars []BuildArg
		args    []string
		shell   = "/bin/sh"
	)
	for _, h := range util.Reverse(history) {
		logrus.Debugf("parsing history line %s", h.CreatedBy)
		if matches := shellLineRegex.FindStringSubmatch(h.CreatedBy); matches != nil {
			logrus.Debugf("found match for SHELL line: %s", matches[1])
			shell = matches[1]
		}

		if matches := argLineRegex.FindStringSubmatch(h.CreatedBy); matches != nil {
			logrus.Debugf("found match for ARG line: %s", matches[1])
			args = append(args, matches[1])
		}

		if matches := runLineRegex.FindStringSubmatch(h.CreatedBy); matches != nil {
			amt := 0
			if amt, err = strconv.Atoi(matches[1]); err != nil {
				logrus.Warnf("invalid build arg amount %s, skipping", matches[1])
			} else if len(args) != amt {
				logrus.Warnf("amount of counted args %d, differs from the amount of build args %d\n", len(args), amt)
			}
			regxp, err := regexp.Compile(strings.Join(append(args, shell), `=(.*)\s`))
			if err != nil {
				logrus.Errorf("unable to compile regex %s\n", err)
			}
			if matches = regxp.FindStringSubmatch(h.CreatedBy); matches != nil {
				// TODO(i dont love this)
				envVars = append(envVars, util.ZipApply(func(name string, value string) BuildArg {
					return BuildArg{
						Name:     name,
						Value:    value,
						Location: h.CreatedBy,
					}
				}, args, matches[1:])...)

			} else {
				logrus.Debugf("RUN regexp %s failed to match", regxp.String())
			}
		}
	}
	return envVars, err
}
