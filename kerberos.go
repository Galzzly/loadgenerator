package loadgenerator

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	krb "github.com/jcmturner/gokrb5/v8/client"
	"github.com/jcmturner/gokrb5/v8/config"
	"github.com/jcmturner/gokrb5/v8/credentials"
)

func getKerberosClient() (*krb.Client, error) {
	conf := os.Getenv("KRB5_CONFIG")
	if conf == "" {
		conf = "/etc/krb5.conf"
	}

	cfg, e := config.Load(conf)
	if e != nil {
		return nil, e
	}

	ccachePath := os.Getenv("KRB5CCNAME")
	if strings.Contains(ccachePath, ":") {
		if strings.HasPrefix(ccachePath, "FILE:") {
			ccachePath = strings.SplitN(ccachePath, ":", 2)[1]
		} else {
			return nil, fmt.Errorf("Unusable ccache: %s", ccachePath)
		}
	} else if ccachePath == "" {
		u, e := user.Current()
		if e != nil {
			return nil, e
		}

		ccachePath = fmt.Sprintf("/tmp/krb5cc_%s", u.Uid)
	}

	ccache, e := credentials.LoadCCache(ccachePath)
	if e != nil {
		return nil, e
	}

	//client, e := krb.NewClientFromCCache(ccache, cfg)
	client, e := krb.NewFromCCache(ccache, cfg)
	if e != nil {
		return nil, e
	}

	return client, nil
}
