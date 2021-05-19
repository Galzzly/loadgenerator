package loadgenerator

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/user"
	"time"

	"github.com/colinmarc/hdfs/v2"
	"github.com/colinmarc/hdfs/v2/hadoopconf"
)

func ConnectToNamenode() (*hdfs.Client, error) {
	client, e := getClient()
	if e != nil {
		return nil, e
	}

	return client, nil
}

func getClient() (*hdfs.Client, error) {
	// Check for the namenode in env
	namenode := os.Getenv("HADOOP_NAMENODE")

	conf, e := hadoopconf.LoadFromEnvironment()
	if e != nil {
		return nil, e
	}

	// if namenode populated, set it in options.Addresses
	options := hdfs.ClientOptionsFromConf(conf)
	if namenode != "" {
		options.Addresses = []string{namenode}
	}

	// Otherwise, just die
	if options.Addresses == nil {
		return nil, errors.New("Cannot find Namenode to connect to")
	}

	/*
		I'm leaving this piece of code in here, it will just fail for clusters with Kerberos
		enabled for now. However, for clusters without kerberos enabled, it will work just fine.
		I will, however, be working towards enabling support for Kerberos in the future, once I
		have more time to work on it.
	*/

	if options.KerberosClient != nil {
		options.KerberosClient, e = getKerberosClient()
		if e != nil {
			return nil, fmt.Errorf("Problem with kerberos auth: %s", e)
		}
	} else {
		options.User = os.Getenv("HADOOP_USER_NAME")
		if options.User == "" {
			u, e := user.Current()
			if e != nil {
				return nil, fmt.Errorf("Unable to determine user: %s", e)
			}

			options.User = u.Username
		}
	}

	dialFunc := (&net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 5 * time.Second,
		DualStack: true,
	}).DialContext

	options.NamenodeDialFunc = dialFunc
	options.DatanodeDialFunc = dialFunc

	/*
		With Kerberos enabled, the hdfs.NewClient return fails due to being unable
		to connect to a namenode. This is likely down to the ciphers used, and so
		will need to have further testing performed.
	*/
	client, e := hdfs.NewClient(options)
	if e != nil {
		return nil, fmt.Errorf("Unable to connect to namenode: %s", e)
	}

	return client, nil
}
