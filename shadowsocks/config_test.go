// Copyright 2021 The Outline Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shadowsocks

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"reflect"
	"testing"
	"time"
)

const (
	redirectURL    = "https://127.0.0.1/200/"
	examplePemCert = `-----BEGIN CERTIFICATE-----
MIIG1TCCBb2gAwIBAgIQD74IsIVNBXOKsMzhya/uyTANBgkqhkiG9w0BAQsFADBP
MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMSkwJwYDVQQDEyBE
aWdpQ2VydCBUTFMgUlNBIFNIQTI1NiAyMDIwIENBMTAeFw0yMDExMjQwMDAwMDBa
Fw0yMTEyMjUyMzU5NTlaMIGQMQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZv
cm5pYTEUMBIGA1UEBxMLTG9zIEFuZ2VsZXMxPDA6BgNVBAoTM0ludGVybmV0IENv
cnBvcmF0aW9uIGZvciBBc3NpZ25lZCBOYW1lcyBhbmQgTnVtYmVyczEYMBYGA1UE
AxMPd3d3LmV4YW1wbGUub3JnMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC
AQEAuvzuzMoKCP8Okx2zvgucA5YinrFPEK5RQP1TX7PEYUAoBO6i5hIAsIKFmFxt
W2sghERilU5rdnxQcF3fEx3sY4OtY6VSBPLPhLrbKozHLrQ8ZN/rYTb+hgNUeT7N
A1mP78IEkxAj4qG5tli4Jq41aCbUlCt7equGXokImhC+UY5IpQEZS0tKD4vu2ksZ
04Qetp0k8jWdAvMA27W3EwgHHNeVGWbJPC0Dn7RqPw13r7hFyS5TpleywjdY1nB7
ad6kcZXZbEcaFZ7ZuerA6RkPGE+PsnZRb1oFJkYoXimsuvkVFhWeHQXCGC1cuDWS
rM3cpQvOzKH2vS7d15+zGls4IwIDAQABo4IDaTCCA2UwHwYDVR0jBBgwFoAUt2ui
6qiqhIx56rTaD5iyxZV2ufQwHQYDVR0OBBYEFCYa+OSxsHKEztqBBtInmPvtOj0X
MIGBBgNVHREEejB4gg93d3cuZXhhbXBsZS5vcmeCC2V4YW1wbGUuY29tggtleGFt
cGxlLmVkdYILZXhhbXBsZS5uZXSCC2V4YW1wbGUub3Jngg93d3cuZXhhbXBsZS5j
b22CD3d3dy5leGFtcGxlLmVkdYIPd3d3LmV4YW1wbGUubmV0MA4GA1UdDwEB/wQE
AwIFoDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwgYsGA1UdHwSBgzCB
gDA+oDygOoY4aHR0cDovL2NybDMuZGlnaWNlcnQuY29tL0RpZ2lDZXJ0VExTUlNB
U0hBMjU2MjAyMENBMS5jcmwwPqA8oDqGOGh0dHA6Ly9jcmw0LmRpZ2ljZXJ0LmNv
bS9EaWdpQ2VydFRMU1JTQVNIQTI1NjIwMjBDQTEuY3JsMEwGA1UdIARFMEMwNwYJ
YIZIAYb9bAEBMCowKAYIKwYBBQUHAgEWHGh0dHBzOi8vd3d3LmRpZ2ljZXJ0LmNv
bS9DUFMwCAYGZ4EMAQICMH0GCCsGAQUFBwEBBHEwbzAkBggrBgEFBQcwAYYYaHR0
cDovL29jc3AuZGlnaWNlcnQuY29tMEcGCCsGAQUFBzAChjtodHRwOi8vY2FjZXJ0
cy5kaWdpY2VydC5jb20vRGlnaUNlcnRUTFNSU0FTSEEyNTYyMDIwQ0ExLmNydDAM
BgNVHRMBAf8EAjAAMIIBBQYKKwYBBAHWeQIEAgSB9gSB8wDxAHcA9lyUL9F3MCIU
VBgIMJRWjuNNExkzv98MLyALzE7xZOMAAAF1+73YbgAABAMASDBGAiEApGuo0EOk
8QcyLe2cOX136HPBn+0iSgDFvprJtbYS3LECIQCN6F+Kx1LNDaEj1bW729tiE4gi
1nDsg14/yayUTIxYOgB2AFzcQ5L+5qtFRLFemtRW5hA3+9X6R9yhc5SyXub2xw7K
AAABdfu92M0AAAQDAEcwRQIgaqwR+gUJEv+bjokw3w4FbsqOWczttcIKPDM0qLAz
2qwCIQDa2FxRbWQKpqo9izUgEzpql092uWfLvvzMpFdntD8bvTANBgkqhkiG9w0B
AQsFAAOCAQEApyoQMFy4a3ob+GY49umgCtUTgoL4ZYlXpbjrEykdhGzs++MFEdce
MV4O4sAA5W0GSL49VW+6txE1turEz4TxMEy7M54RFyvJ0hlLLNCtXxcjhOHfF6I7
qH9pKXxIpmFfJj914jtbozazHM3jBFcwH/zJ+kuOSIBYJ5yix8Mm3BcC+uZs6oEB
XJKP0xgIF3B6wqNLbDr648/2/n7JVuWlThsUT6mYnXmxHsOrsQ0VhalGtuXCWOha
/sgUKGiQxrjIlH/hD4n6p9YJN6FitwAntb7xsV5FKAazVBXmw8isggHOhuIr4Xrk
vUzLnF7QYsJhvYtaYrZ2MLxGD+NFI8BkXw==
-----END CERTIFICATE-----`
	exampleCertFingerprint = "IA3K+nZ8hFDs5kSHnAYqDN9SJA/gW7frKEYRw67z7C4="
)

var proxies = []ProxyConfig{
	{"ssconf.test", 123, "passw0rd", "chacha20-ietf-poly1305", "ssconf-test-1", "", ""},
	{"ssconf-ii.test", 456, "dr0wssap", "chacha20-ietf-poly1305", "ssconf-test-2", "", ""},
}

func TestFetchConfig(t *testing.T) {
	serverAddr := "127.0.0.1:9999"
	cert, err := makeTLSCertificate()
	if err != nil {
		t.Fatalf("Failed to generate TLS certificate: %v", err)
	}
	certFingerprint := computeCertificateFingerprint(cert.Certificate[0])
	server := makeOnlineConfigServer(serverAddr, cert)
	go server.ListenAndServeTLS("", "")
	defer server.Close()

	t.Run("Success", func(t *testing.T) {
		req := FetchConfigRequest{
			fmt.Sprintf("https://%s/200", serverAddr), "GET", certFingerprint}
		res, err := FetchConfig(req)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if res.HTTPStatusCode != 200 {
			t.Errorf("Expected 200 HTTP status code, got %d", res.HTTPStatusCode)
		}
		if res.RedirectURL != "" {
			t.Errorf("Unexpected redirect URL: %s", res.RedirectURL)
		}
		if !reflect.DeepEqual(proxies, res.Proxies) {
			t.Errorf("Proxy configurations don't match. Want %v, have %v", proxies, res.Proxies)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		req := FetchConfigRequest{
			fmt.Sprintf("https://%s/404", serverAddr), "GET", certFingerprint}
		res, err := FetchConfig(req)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if res.HTTPStatusCode != 404 {
			t.Errorf("Expected 404 HTTP status code, got %d", res.HTTPStatusCode)
		}
		if res.RedirectURL != "" {
			t.Errorf("Unexpected redirect URL: %s", res.RedirectURL)
		}
		if len(res.Proxies) > 0 {
			t.Errorf("Expected empty proxy configurations, got: %v", res.Proxies)
		}
	})

	t.Run("Redirect", func(t *testing.T) {
		req := FetchConfigRequest{
			fmt.Sprintf("https://%s/301", serverAddr), "GET", certFingerprint}
		res, err := FetchConfig(req)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if res.HTTPStatusCode != 301 {
			t.Errorf("Expected 301 HTTP status code, got %d", res.HTTPStatusCode)
		}
		if res.RedirectURL != redirectURL {
			t.Errorf("Expected redirect URL %s , got %s", redirectURL, res.RedirectURL)
		}
		if len(res.Proxies) > 0 {
			t.Errorf("Expected empty proxy configurations, got: %v", res.Proxies)
		}
	})

	t.Run("CertificateFingerprint", func(t *testing.T) {
		req := FetchConfigRequest{
			fmt.Sprintf("https://%s/success", serverAddr), "GET", "wrongcertfp"}
		_, err := FetchConfig(req)
		if err == nil {
			t.Fatalf("Expected TLS certificate validation error")
		}
	})

	t.Run("NonHTTPSURL", func(t *testing.T) {
		req := FetchConfigRequest{
			fmt.Sprintf("http://%s/success", serverAddr), "GET", certFingerprint}
		_, err := FetchConfig(req)
		if err == nil {
			t.Fatalf("Expected error for non-HTTPs URL")
		}
	})
}

// HTTP handler for a fake online config server.
type onlineConfigHandler struct{}

func (h onlineConfigHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/200" {
		res := sip008Response{proxies, 1}
		data, _ := json.Marshal(res)
		h.sendResponse(w, 200, data)
	} else if req.URL.Path == "/404" {
		h.sendResponse(w, 404, []byte("Not Found"))
	} else if req.URL.Path == "/301" {
		w.Header().Add("Location", redirectURL)
		h.sendResponse(w, 301, []byte{})
	}
}

func (onlineConfigHandler) sendResponse(w http.ResponseWriter, code int, data []byte) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

// Returns a SIP008 online config HTTPs server with TLS certificate cert.
func makeOnlineConfigServer(addr string, cert tls.Certificate) http.Server {
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	return http.Server{
		Addr:      addr,
		TLSConfig: tlsConfig,
		Handler:   onlineConfigHandler{},
	}
}

// Generates a self-signed TLS certificate for localhost.
func makeTLSCertificate() (tls.Certificate, error) {
	now := time.Now()
	template := &x509.Certificate{
		SerialNumber: big.NewInt(now.Unix()),
		Subject: pkix.Name{
			Organization: []string{"online config"},
		},
		IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1)}, // Valid for localhost
		NotBefore:             now,
		NotAfter:              now.AddDate(0, 0, 1), // Valid for one day
		SubjectKeyId:          []byte{55, 43, 04, 45, 87, 65},
		BasicConstraintsValid: true,
		IsCA:                  true, // Self-signed
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature |
			x509.KeyUsageCertSign,
	}

	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return tls.Certificate{}, err
	}

	derCert, err := x509.CreateCertificate(rand.Reader, template, template,
		key.Public(), key)
	if err != nil {
		return tls.Certificate{}, err
	}

	var cert tls.Certificate
	cert.Certificate = append(cert.Certificate, derCert)
	cert.PrivateKey = key
	return cert, nil
}

func TestComputeCertificateFingerprint(t *testing.T) {
	pemCertData := []byte(examplePemCert)
	block, _ := pem.Decode(pemCertData)
	if block == nil || block.Type != "CERTIFICATE" {
		t.Fatalf("Failed to decode certificate PEM block")
	}

	certFingerprint := computeCertificateFingerprint(block.Bytes)
	if certFingerprint != exampleCertFingerprint {
		t.Errorf("Certificate fingerprints don't match. Want %s, got %s",
			exampleCertFingerprint, certFingerprint)
	}
}