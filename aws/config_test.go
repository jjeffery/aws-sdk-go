package aws

import (
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
)

var testCredentials = credentials.NewChainCredentials([]credentials.Provider{
	&credentials.EnvProvider{},
	&credentials.SharedCredentialsProvider{
		Filename: "TestFilename",
		Profile:  "TestProfile"},
	&credentials.EC2RoleProvider{ExpiryWindow: 5 * time.Minute},
})

var copyTestConfig = Config{
	Credentials:             testCredentials,
	Endpoint:                String("CopyTestEndpoint"),
	Region:                  String("COPY_TEST_AWS_REGION"),
	DisableSSL:              Bool(true),
	HTTPClient:              http.DefaultClient,
	LogHTTPBody:             Bool(true),
	LogLevel:                Int(2),
	Logger:                  os.Stdout,
	MaxRetries:              Int(DefaultRetries),
	DisableParamValidation:  Bool(true),
	DisableComputeChecksums: Bool(true),
	S3ForcePathStyle:        Bool(true),
}

func TestCopy(t *testing.T) {
	want := copyTestConfig
	got := copyTestConfig.Copy()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Copy() = %+v", got)
		t.Errorf("    want %+v", want)
	}
}

func TestCopyReturnsNewInstance(t *testing.T) {
	want := copyTestConfig
	got := copyTestConfig.Copy()
	if &got == &want {
		t.Errorf("Copy() = %p; want different instance as source %p", &got, &want)
	}
}

var mergeTestZeroValueConfig = Config{}

var mergeTestConfig = Config{
	Credentials:             testCredentials,
	Endpoint:                String("MergeTestEndpoint"),
	Region:                  String("MERGE_TEST_AWS_REGION"),
	DisableSSL:              Bool(true),
	HTTPClient:              http.DefaultClient,
	LogHTTPBody:             Bool(true),
	LogLevel:                Int(2),
	Logger:                  os.Stdout,
	MaxRetries:              Int(10),
	DisableParamValidation:  Bool(true),
	DisableComputeChecksums: Bool(true),
	S3ForcePathStyle:        Bool(true),
}

var mergeTests = []struct {
	cfg  *Config
	in   *Config
	want *Config
}{
	{&Config{}, nil, &Config{}},
	{&Config{}, &mergeTestZeroValueConfig, &Config{}},
	{&Config{}, &mergeTestConfig, &mergeTestConfig},
}

var mergeErrorFmt = `
Merge Failed [%d]
  Config  %+v
    Merge(%+v)
      got %+v
     want %+v
`

func TestMerge(t *testing.T) {
	for i, tt := range mergeTests {
		got := tt.cfg.Merge(tt.in)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf(mergeErrorFmt, i, tt.cfg, tt.in, got, tt.want)
		}
	}
}
