package util

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql/driver"
	"fmt"
	"io/ioutil"
	"strings"

	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/constant"
	"github.com/go-sql-driver/mysql"
)

func AppendSSLConnectionStringParamForMysqlIfNecessary(dbKey, databaseConnectionString, params, sqlCACertFile string, clientCaFiles ...string) (string, *errors.Error) {
	if sqlCACertFile != "" {
		certBytes, err := ioutil.ReadFile(sqlCACertFile)
		if err != nil {
			return "", errors.Errorf(constant.SystemInternalError, "failed to read sql ca file[path:=%s], error:%++v", sqlCACertFile, err)
		}

		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(certBytes); !ok {
			return "", errors.Errorf(constant.SystemInternalError, "failed to parse sql ca, error:%++v", err)
		}
		var clientCert []tls.Certificate
		if len(clientCaFiles) == 2 && len(clientCaFiles[0]) > 0 && len(clientCaFiles[1]) > 0 {
			clientCert = make([]tls.Certificate, 0, 1)
			certs, err := tls.LoadX509KeyPair(clientCaFiles[0], clientCaFiles[1])
			if err != nil {
				return "", errors.Errorf(constant.SystemInternalError, "failed to load client ca file[path=%s] and key[path=%s], error:%++v", clientCaFiles[0], clientCaFiles[1], err)
			}
			clientCert = append(clientCert, certs)
		}

		tlsConfig := &tls.Config{
			InsecureSkipVerify: false,
			RootCAs:            caCertPool,
			Certificates:       clientCert,
		}

		mysql.RegisterTLSConfig(dbKey, tlsConfig)
		databaseConnectionString += fmt.Sprintf("?tls=%s", dbKey)
		if len(params) > 0 {
			if strings.HasPrefix(params, "?") {
				params = strings.Replace(params, "?", "&", 1)
			}
			databaseConnectionString += params
		}
	} else {
		if len(params) > 0 {
			databaseConnectionString += params
		}
	}

	return databaseConnectionString, nil
}

func AppendSSLConnectionStringParamForPostgresIfNecessary(databaseConnectionString, params, sqlCACertFile string, clientCaFiles ...string) (string, *errors.Error) {
	if sqlCACertFile != "" {
		databaseConnectionString += fmt.Sprintf(" sslrootcert=%s", sqlCACertFile)
		sslmode := "require"
		if len(clientCaFiles) == 2 && len(clientCaFiles[0]) > 0 && len(clientCaFiles[1]) > 0 {
			databaseConnectionString += fmt.Sprintf(" sslcert=%s sslkey=%s", clientCaFiles[0], clientCaFiles[1])
		}

		databaseConnectionString += fmt.Sprintf(" sslmode=%s", sslmode)

	} else {
		databaseConnectionString += " sslmode=disable"
	}

	if len(params) > 0 {
		databaseConnectionString += fmt.Sprintf(" %s", params)
	}

	return databaseConnectionString, nil
}

func NamedValueToValue(nvs []driver.NamedValue) []driver.Value {
	vs := make([]driver.Value, 0, len(nvs))
	for _, nv := range nvs {
		vs = append(vs, nv.Value)
	}
	return vs
}

func ValueToNamedValue(vs []driver.Value) []driver.NamedValue {
	nvs := make([]driver.NamedValue, 0, len(vs))
	for i, v := range vs {
		nvs = append(nvs, driver.NamedValue{
			Value:   v,
			Ordinal: i,
		})
	}
	return nvs
}
