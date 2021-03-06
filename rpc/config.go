package rpc

import (
	"fmt"
	"time"
)

import (
	config "github.com/koding/multiconfig"
)

type (
	GettySessionParam struct {
		CompressEncoding bool   `default:"false" yaml:"compress_encoding" json:"compress_encoding,omitempty"`
		TcpNoDelay       bool   `default:"true" yaml:"tcp_no_delay" json:"tcp_no_delay,omitempty"`
		TcpKeepAlive     bool   `default:"true" yaml:"tcp_keep_alive" json:"tcp_keep_alive,omitempty"`
		KeepAlivePeriod  string `default:"180s" yaml:"keep_alive_period" json:"keep_alive_period,omitempty"`
		keepAlivePeriod  time.Duration
		TcpRBufSize      int    `default:"262144" yaml:"tcp_r_buf_size" json:"tcp_r_buf_size,omitempty"`
		TcpWBufSize      int    `default:"65536" yaml:"tcp_w_buf_size" json:"tcp_w_buf_size,omitempty"`
		PkgRQSize        int    `default:"1024" yaml:"pkg_rq_size" json:"pkg_rq_size,omitempty"`
		PkgWQSize        int    `default:"1024" yaml:"pkg_wq_size" json:"pkg_wq_size,omitempty"`
		TcpReadTimeout   string `default:"1s" yaml:"tcp_read_timeout" json:"tcp_read_timeout,omitempty"`
		tcpReadTimeout   time.Duration
		TcpWriteTimeout  string `default:"5s" yaml:"tcp_write_timeout" json:"tcp_write_timeout,omitempty"`
		tcpWriteTimeout  time.Duration
		WaitTimeout      string `default:"7s" yaml:"wait_timeout" json:"wait_timeout,omitempty"`
		waitTimeout      time.Duration
		MaxMsgLen        int    `default:"1024" yaml:"max_msg_len" json:"max_msg_len,omitempty"`
		SessionName      string `default:"rpc" yaml:"session_name" json:"session_name,omitempty"`
	}

	RegistryConfig struct {
		Type             string `default:"etcd" yaml:"type" json:"type,omitempty"`
		Addr             string `default:"127.0.0.1:2379" yaml:"addr" json:"addr,omitempty"`
		KeepaliveTimeout int    `default:"5" yaml:"keepalive_time" json:"keepalive_timeout,omitempty"`
		Root             string `default:"getty" yaml:"keepalive_time" json:"keepalive_timeout,omitempty"`
		IDC              string `default:"idc-bj" yaml:"idc" json:"idc,omitempty"`
		NodeID           string `default:"node0" yaml:"node_id" json:"node_id,omitempty"`
	}

	// Config holds supported types by the multiconfig package
	ServerConfig struct {
		// local address
		AppName     string   `default:"rcp-server" yaml:"app_name" json:"app_name,omitempty"`
		Host        string   `default:"127.0.0.1" yaml:"host" json:"host,omitempty"`
		Ports       []string `yaml:"ports" json:"ports,omitempty"` // `default:["10000"]`
		ProfilePort int      `default:"10086" yaml:"profile_port" json:"profile_port,omitempty"`
		CodecType   string   `default:"json" yaml:"codec_type" json:"codec_type,omitempty"`
		codecType   gettyCodecType

		// session
		SessionTimeout string `default:"60s" yaml:"session_timeout" json:"session_timeout,omitempty"`
		sessionTimeout time.Duration
		SessionNumber  int `default:"1000" yaml:"session_number" json:"session_number,omitempty"`

		// app
		FailFastTimeout string `default:"5s" yaml:"fail_fast_timeout" json:"fail_fast_timeout,omitempty"`
		failFastTimeout time.Duration

		// session tcp parameters
		GettySessionParam GettySessionParam `required:"true" yaml:"getty_session_param" json:"getty_session_param,omitempty"`

		// registry center
		Registry RegistryConfig `required:"true" yaml:"registry_config" json:"registry_config,omitempty"`
	}

	// Config holds supported types by the multiconfig package
	ClientConfig struct {
		// local address
		AppName     string   `default:"rcp-client" yaml:"app_name" json:"app_name,omitempty"`
		Host        string   `default:"127.0.0.1" yaml:"host" json:"host,omitempty"`
		Ports       []string `yaml:"ports" json:"ports,omitempty"` // `default:["10000"]`
		ProfilePort int      `default:"10086" yaml:"profile_port" json:"profile_port,omitempty"`

		// server
		ServerHost string `default:"127.0.0.1" yaml:"server_host" json:"server_host,omitempty"`
		ServerPort int    `default:"10000" yaml:"server_port" json:"server_port,omitempty"`

		// session pool
		ConnectionNum int `default:"16" yaml:"connection_num" json:"connection_num,omitempty"`

		// heartbeat
		HeartbeatPeriod string `default:"15s" yaml:"heartbeat_period" json:"heartbeat_period,omitempty"`
		heartbeatPeriod time.Duration

		// session
		SessionTimeout string `default:"60s" yaml:"session_timeout" json:"session_timeout,omitempty"`
		sessionTimeout time.Duration

		// app
		FailFastTimeout string `default:"5s" yaml:"fail_fast_timeout" json:"fail_fast_timeout,omitempty"`
		failFastTimeout time.Duration

		// session tcp parameters
		GettySessionParam GettySessionParam `required:"true" yaml:"getty_session_param" json:"getty_session_param,omitempty"`

		// registry center
		Registry RegistryConfig `required:"true" yaml:"registry_config" json:"registry_config,omitempty"`
	}
)

func loadClientConf(confFile string) *ClientConfig {
	var err error
	conf := new(ClientConfig)
	config.MustLoadWithPath(confFile, conf)
	conf.heartbeatPeriod, err = time.ParseDuration(conf.HeartbeatPeriod)
	if err != nil {
		panic(fmt.Sprintf("time.ParseDuration(HeartbeatPeroid{%#v}) = error{%v}", conf.HeartbeatPeriod, err))
	}
	conf.sessionTimeout, err = time.ParseDuration(conf.SessionTimeout)
	if err != nil {
		panic(fmt.Sprintf("time.ParseDuration(SessionTimeout{%#v}) = error{%v}", conf.SessionTimeout, err))
	}
	conf.failFastTimeout, err = time.ParseDuration(conf.FailFastTimeout)
	if err != nil {
		panic(fmt.Sprintf("time.ParseDuration(FailFastTimeout{%#v}) = error{%v}", conf.FailFastTimeout, err))
	}
	conf.GettySessionParam.keepAlivePeriod, err = time.ParseDuration(conf.GettySessionParam.KeepAlivePeriod)
	if err != nil {
		panic(fmt.Sprintf("time.ParseDuration(KeepAlivePeriod{%#v}) = error{%v}", conf.GettySessionParam.KeepAlivePeriod, err))
	}
	conf.GettySessionParam.tcpReadTimeout, err = time.ParseDuration(conf.GettySessionParam.TcpReadTimeout)
	if err != nil {
		panic(fmt.Sprintf("time.ParseDuration(TcpReadTimeout{%#v}) = error{%v}", conf.GettySessionParam.TcpReadTimeout, err))
	}
	conf.GettySessionParam.tcpWriteTimeout, err = time.ParseDuration(conf.GettySessionParam.TcpWriteTimeout)
	if err != nil {
		panic(fmt.Sprintf("time.ParseDuration(TcpWriteTimeout{%#v}) = error{%v}", conf.GettySessionParam.TcpWriteTimeout, err))
	}
	conf.GettySessionParam.waitTimeout, err = time.ParseDuration(conf.GettySessionParam.WaitTimeout)
	if err != nil {
		panic(fmt.Sprintf("time.ParseDuration(WaitTimeout{%#v}) = error{%v}", conf.GettySessionParam.WaitTimeout, err))
	}
	return conf
}

func loadServerConf(confFile string) *ServerConfig {
	var err error
	conf := new(ServerConfig)
	config.MustLoadWithPath(confFile, conf)

	conf.sessionTimeout, err = time.ParseDuration(conf.SessionTimeout)
	if err != nil {
		panic(fmt.Sprintf("time.ParseDuration(SessionTimeout{%#v}) = error{%v}", conf.SessionTimeout, err))
	}
	conf.failFastTimeout, err = time.ParseDuration(conf.FailFastTimeout)
	if err != nil {
		panic(fmt.Sprintf("time.ParseDuration(FailFastTimeout{%#v}) = error{%v}", conf.FailFastTimeout, err))
	}
	conf.GettySessionParam.keepAlivePeriod, err = time.ParseDuration(conf.GettySessionParam.KeepAlivePeriod)
	if err != nil {
		panic(fmt.Sprintf("time.ParseDuration(KeepAlivePeriod{%#v}) = error{%v}", conf.GettySessionParam.KeepAlivePeriod, err))
	}
	conf.GettySessionParam.tcpReadTimeout, err = time.ParseDuration(conf.GettySessionParam.TcpReadTimeout)
	if err != nil {
		panic(fmt.Sprintf("time.ParseDuration(TcpReadTimeout{%#v}) = error{%v}", conf.GettySessionParam.TcpReadTimeout, err))
	}
	conf.GettySessionParam.tcpWriteTimeout, err = time.ParseDuration(conf.GettySessionParam.TcpWriteTimeout)
	if err != nil {
		panic(fmt.Sprintf("time.ParseDuration(TcpWriteTimeout{%#v}) = error{%v}", conf.GettySessionParam.TcpWriteTimeout, err))
	}
	conf.GettySessionParam.waitTimeout, err = time.ParseDuration(conf.GettySessionParam.WaitTimeout)
	if err != nil {
		panic(fmt.Sprintf("time.ParseDuration(WaitTimeout{%#v}) = error{%v}", conf.GettySessionParam.WaitTimeout, err))
	}
	return conf
}
