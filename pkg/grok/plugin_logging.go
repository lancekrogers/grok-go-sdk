package grok

import (
	"context"
	"log"
	"strings"
)

type LoggingPlugin struct {
	SanitizeSecrets bool
	TruncateLength  int
	Logger          *log.Logger
}

func (p *LoggingPlugin) Name() string { return "logging" }

func (p *LoggingPlugin) logger() *log.Logger {
	if p.Logger != nil {
		return p.Logger
	}
	return log.Default()
}

func (p *LoggingPlugin) truncate(s string) string {
	if p.TruncateLength > 0 && len(s) > p.TruncateLength {
		return s[:p.TruncateLength] + "..."
	}
	return s
}

func (p *LoggingPlugin) sanitize(args []string) []string {
	if !p.SanitizeSecrets {
		return args
	}
	out := make([]string, len(args))
	skipNext := false
	for i, a := range args {
		if skipNext {
			out[i] = "***"
			skipNext = false
			continue
		}
		if strings.HasPrefix(a, "--secret") || strings.HasPrefix(a, "--xai-api-key") {
			out[i] = a
			if !strings.Contains(a, "=") {
				skipNext = true
			} else {
				eq := strings.IndexByte(a, '=')
				out[i] = a[:eq+1] + "***"
			}
			continue
		}
		out[i] = a
	}
	return out
}

func (p *LoggingPlugin) OnBeforeRun(_ context.Context, ev BeforeRunEvent) error {
	p.logger().Printf("grok run start prompt=%q args=%v", p.truncate(ev.Prompt), p.sanitize(ev.Args))
	return nil
}

func (p *LoggingPlugin) OnAfterRun(_ context.Context, ev AfterRunEvent) error {
	if ev.Err != nil {
		p.logger().Printf("grok run error: %v", ev.Err)
		return nil
	}
	if ev.Result != nil {
		p.logger().Printf("grok run ok text=%q", p.truncate(ev.Result.Text))
	}
	return nil
}
