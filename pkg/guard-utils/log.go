/*
Copyright 2022 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package guardutils

import (
	"encoding/json"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel

	}
}

func CreateLogger(logLevel string) *zap.SugaredLogger {
	rawJSON := []byte(`{
		"level": "info",
		"encoding": "console",
		"outputPaths": ["stdout"],
		"development": false,
		"errorOutputPaths": ["stderr"],
		"encoderConfig": {
		  "messageKey": "message",
		  "levelKey": "level",
		  "levelEncoder": "lowercase"
		}
	  }`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}

	cfg.Level = zap.NewAtomicLevelAt(getLogLevel(logLevel))
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return logger.Sugar()
}
