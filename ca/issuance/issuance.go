package issuance

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"math/big"
	"time"
)

type IssueType int

const (
	IssueTypeClient IssueType = iota + 1
	IssueTypeServer
	IssueTypeIntermediate // 中间证书
)

type Cert struct {
	caKey      []byte
	caCert     []byte
	csr        []byte
	issueType  IssueType
	days       int
	ocspServer []string
	crlDist    []string
}
type CertOption func(*Cert)

func NewCert(caKey, caCert, csr []byte, opts ...CertOption) *Cert {
	c := &Cert{
		caKey:     caKey,
		caCert:    caCert,
		csr:       csr,
		days:      365,
		issueType: IssueTypeClient,
	}
	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithDays(days int) CertOption {
	return func(c *Cert) {
		c.days = days
	}
}

func WithIssueType(issueType IssueType) CertOption {
	return func(c *Cert) {
		c.issueType = issueType
	}
}
func WithOcspServer(server []string) CertOption {
	return func(c *Cert) {
		c.ocspServer = server
	}
}

func WithCrlPoint(points []string) CertOption {
	return func(c *Cert) {
		c.crlDist = points
	}
}

// Sign signs the certificate
// return pemCert, error
func (c *Cert) Sign() ([]byte, error) {

	//解析csr
	block, _ := pem.Decode(c.csr)
	if block == nil || block.Type != "CERTIFICATE REQUEST" {
		return nil, errors.New("failed to decode PEM block containing CSR")
	}

	csr, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		return nil, err
	}

	//解析ca证书
	block, _ = pem.Decode(c.caCert)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing CA certificate")
	}

	caCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	//解析ca私钥
	block, _ = pem.Decode(c.caKey)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, err
	}
	caKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// 创建证书模板
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               csr.Subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Duration(c.days) * 24 * time.Hour),
		ExtKeyUsage:           make([]x509.ExtKeyUsage, 0),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}
	switch c.issueType {
	case IssueTypeClient:
		template.ExtKeyUsage = append(template.ExtKeyUsage, x509.ExtKeyUsageClientAuth)
	case IssueTypeServer:
		template.DNSNames = csr.DNSNames
		template.IPAddresses = csr.IPAddresses
		template.ExtKeyUsage = append(template.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
	case IssueTypeIntermediate:
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
		template.KeyUsage |= x509.KeyUsageDigitalSignature
		template.KeyUsage |= x509.KeyUsageCRLSign
		template.ExtKeyUsage = append(template.ExtKeyUsage, x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageOCSPSigning)
	}
	if c.ocspServer != nil {
		template.OCSPServer = c.ocspServer
	}
	if c.crlDist != nil {
		template.CRLDistributionPoints = c.crlDist
	}
	// 解析扩展配置文件并应用到证书模板
	// 注意：这是一个简化的示例。解析 OpenSSL 配置文件格式需要自定义代码。
	//if len(c.dns) > 0 {
	//	extConfigData := []byte(fmt.Sprintf("subjectAltName = DNS:%s", c.dns))
	//	if len(extConfigData) > 0 {
	//		// 假设扩展配置文件包含如下内容:
	//		// subjectAltName = @alt_names
	//		// [alt_names]
	//		// DNS.1 = example.com
	//		exts, err := parseExtensions(extConfigData)
	//		if err != nil {
	//			return nil, err
	//		}
	//		template.ExtraExtensions = exts
	//	}
	//}

	// 使用 CA 签署 CSR
	certDER, err := x509.CreateCertificate(rand.Reader, &template, caCert, csr.PublicKey, caKey)
	if err != nil {
		return nil, err
	}
	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	return pemCert, nil
}

// parseExtensions 解析扩展配置文件内容
// 注意：这是一个示例解析函数，您需要根据实际的扩展配置文件格式进行调整。
//func parseExtensions(data []byte) ([]pkix.Extension, error) {
//	// 示例解析逻辑，这里应根据具体的配置文件格式进行调整
//	// 假设配置文件是简单的 Key=Value 格式
//	var extensions []pkix.Extension
//	lines := string(data)
//	for _, line := range strings.Split(lines, "\n") {
//		parts := strings.SplitN(line, "=", 2)
//		if len(parts) == 2 {
//			key := strings.TrimSpace(parts[0])
//			value := strings.TrimSpace(parts[1])
//			ext := pkix.Extension{
//				Id:    oidFromString(key), // 需要实现 oidFromString 将字符串转换为 OID
//				Value: []byte(value),
//			}
//			extensions = append(extensions, ext)
//		}
//	}
//	return extensions, nil
//}
//
//// oidFromString 将字符串转换为 OID
//func oidFromString(s string) asn1.ObjectIdentifier {
//	// 这是一个示例实现，需要根据实际情况调整
//	// 这里假设输入格式为 "1.2.3.4"
//	parts := strings.Split(s, ".")
//	oid := make([]int, len(parts))
//	for i, part := range parts {
//		num, err := strconv.Atoi(part)
//		if err != nil {
//			return nil
//		}
//		oid[i] = num
//	}
//	return oid
//}
