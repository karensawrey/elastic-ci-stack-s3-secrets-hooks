package main

import (
	"fmt"
	"log"
	"os"

	"github.com/buildkite/elastic-ci-stack-s3-secrets-hooks/s3secrets-helper/v2/env"
	"github.com/buildkite/elastic-ci-stack-s3-secrets-hooks/s3secrets-helper/v2/s3"
	"github.com/buildkite/elastic-ci-stack-s3-secrets-hooks/s3secrets-helper/v2/secrets"
	"github.com/buildkite/elastic-ci-stack-s3-secrets-hooks/s3secrets-helper/v2/sshagent"
)

func main() {
	log := log.New(os.Stderr, "", log.Lmsgprefix)
	if err := mainWithError(log); err != nil {
		log.Fatalf("fatal error: %v", err)
	}
}

func mainWithError(log *log.Logger) error {
	bucket := os.Getenv(env.EnvBucket)
	if bucket == "" {
		return nil
	}

	// May be empty string
	regionHint := os.Getenv(env.EnvRegion)

	prefix := os.Getenv(env.EnvPrefix)
	if prefix == "" {
		prefix = os.Getenv(env.EnvPipeline)
	}
	if prefix == "" {
		return fmt.Errorf("One of %s or %s environment variables required", env.EnvPrefix, env.EnvPipeline)
	}

	client, err := s3.New(log, bucket, regionHint)
	if err != nil {
		return err
	}

	agent := &sshagent.Agent{}

	credHelper := os.Getenv(env.EnvCredHelper)
	if credHelper == "" {
		return fmt.Errorf("%s environment variable required", env.EnvCredHelper)
	}

	return secrets.Run(secrets.Config{
		Repo:                os.Getenv(env.EnvRepo),
		Bucket:              bucket,
		Prefix:              prefix,
		Client:              client,
		Logger:              log,
		SSHAgent:            agent,
		EnvSink:             os.Stdout,
		GitCredentialHelper: credHelper,
	})
}
