package assistant

import (
	internal "github.com/invxp/go-layout-assistant/internal/http"
	"github.com/invxp/go-layout-assistant/internal/util/convert"
	"net/http"
	"testing"
	"time"
)

var testConfig = &Config{
	HTTP: struct {
		Enable  bool
		Address string
	}{Enable: true, Address: ":80"},

	Log: struct {
		Enable               bool
		Path                 string
		MaxAgeHours          uint
		MaxRotationMegabytes uint
	}{Enable: false, Path: "logs", MaxAgeHours: 3 * 24, MaxRotationMegabytes: 1 * 1024},

	MySQL: struct {
		Enable                bool
		Host                  string
		Username              string
		Password              string
		Database              string
		MaxConnIdleTimeMinute uint
		MaxConnLifeTimeMinute uint
		MaxOpenConnections    uint
		MaxIdleConnections    uint
	}{Enable: false, Host: "127.0.0.1:3306", Username: "user", Password: "pwd", Database: "test", MaxConnIdleTimeMinute: 10, MaxConnLifeTimeMinute: 10, MaxOpenConnections: 1000, MaxIdleConnections: 100},

	Redis: struct {
		Enable                bool
		Host                  string
		Password              string
		Database              uint
		MaxIdle               uint
		MaxActive             uint
		MaxConnTimeoutSecond  uint
		MaxConnIdleTimeMinute uint
	}{Enable: false, Host: "127.0.0.1:6379", Password: "pwd", Database: 1, MaxIdle: 100, MaxActive: 1000, MaxConnTimeoutSecond: 5, MaxConnIdleTimeMinute: 10},
}

func TestMain(m *testing.M) {
	srv, err := New(
		WithConfig(testConfig),
		WithMySQLConfig(map[string]string{"timeout": "5s"}))

	if err != nil {
		panic(err)
	}

	go func() {
		panic(srv.Serv())
	}()

	time.Sleep(time.Second)

	m.Run()
}

func TestHTTPGet(t *testing.T) {
	resp, err := internal.Request(http.MethodGet, "http://localhost/test?key=k&value=v", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(convert.ByteToString(resp.Data.Payload))

	resp, err = internal.Request(http.MethodGet, "http://localhost/test?key=k", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(convert.ByteToString(resp.Data.Payload))

	resp, err = internal.Request(http.MethodGet, "http://localhost/test?value=v", nil, nil)
	if err == nil {
		t.Fatal("bind query key was not found")
	}
	t.Log(err)
}

func TestHTTPPost(t *testing.T) {
	resp, err := internal.Request(http.MethodPost, "http://localhost/api/cron", internal.RequestPOST{Key: "*/1 * * * * *", Value: convert.StringToByte("...CRON...")}, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(convert.ByteToString(resp.Data.Payload))

	resp, err = internal.Request(http.MethodPost, "http://localhost/api/cron", internal.RequestPOST{Key: "*/1 * * * * *"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(convert.ByteToString(resp.Data.Payload))

	resp, err = internal.Request(http.MethodPost, "http://localhost/api/cron", internal.RequestPOST{}, nil)
	if err == nil {
		t.Fatal("bind post key was not found")
	}
	t.Log(err)

	time.Sleep(5 * time.Second)
}

func TestHTTPGetCustomHeaderForAuth(t *testing.T) {
	_, err := internal.Request(http.MethodGet, "http://localhost/test?key=k&value=v", nil, map[string]string{"Auth": "FALSE"})
	if err == nil {
		t.Fatal("auth success")
	}
	t.Log(err)
}
