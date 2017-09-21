package license

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func RsaEncrypt(origData []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// 解密
func RsaDecrypt(ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

// 公钥和私钥

var privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIIEpgIBAAKCAQEAybYWWudxZvSOA3hjtkyqPNgaK+PkJTifrFJrCRmzW30/fi1D
5Rcv/TSSxjNX7FfGvdIai7z2NvGXZVGgWGISUkJoWjdNZHRAp74ieGdpR1GNUY58
DC8IK26br5PrHEYt5qhPSkDONi9A6ugwavvVMtmbIM1vNtchlTcKCY2UzZOua1X1
LfajmJfiBucyY/C+oEVHYhoJQvOdWVTLsxrjJreGoAABXFj+CpGMsMK26DUnmYhG
IMBBOaiyVuHNjwAlO//EuHdkt6s/c/NU7DJYKNXRh1XypksFtgdq2xMUBag+zsJo
RHtFVKVN21sP8ApZAWgI/5pj+YW99CVl9vKsbQIDAQABAoIBAQCpISNXb25XnECj
SkOZLGklgTFYkcnPJ85CaAzVYZZQ5RDi1EN0iF+1mGplA9H6GpNKwCt/9Z4g7as6
yhl/YiPok0P6ORqMWymHPNacTGEq6odq1eTPNnRMLn8d1hIx7+o21/M72GDPcAmn
ra5DVgsqiukWtQpGWlYPTkn9PpiBUY6DkKLiBL6/HEl5gJ+j9fh/vAr5WClJdniM
qKVPCi9MUyAGa2GlqkZEEgMN+B6PWpchbsiIvV47M6slu8ZxWZVj0A8BvUxMc9z3
LQ7itoHHIyZ3JC6Tb5AR1xKF6Ka+cIrbONa7EWYhpLnwRQOnH0EfVwrsj5p2rZon
hq8oAfBhAoGBAMvtJFucy1k19mOobmpZgdrgRHDAwfTVBG6ltIV5BV62IfgT8TYn
AJh9BLu2QyFTzQTCMsqw+7FwI6OSG06yM9gk/bJBSqdm/EIhbqt2w0jQk00uy+YP
IKuAjm1DO90c9azGU5gjkGY5EC0yGIoU5HbpYOCpa48LVvicQJ5o7TppAoGBAP04
JRgGOai5e5VFrFCzbKy0SLzG/+prR6ye9r6ZiKVoMSHhZzhdtpiJPU5HZN080WR1
ch+94D+ZW87XFGqEc06ppKNe3A9etJnwLXz27za2YAlMjKZpwHeZRUI9Nag9mwiI
Wk4+nIQXMdfpNvN3F/Y8PA0ZXtp917SNhBlc53llAoGBAIYTn9EAEQ7RlPLHHfmc
ae1PgJAlnCBuIeDh4APVovs3grQJ4JD7KcAYipEkb5Ss9WIIkL6FiKaMFVKJUJz4
O3AEOi0GEqBn7LlKI+pmRlVMlVysxPC4x4EbIAmkp+pkDP8Q0ot37ovXPtSIWGwI
+oFYRhLQMWC2UvNYQIAmerrpAoGBANAHXSoUizAbWfUzbTJLhZ+I6Sz1y+95PUQK
wKmrlwBUzfCIrnU3QDimnw/9RVcgjOIcPqrnIiD9y9ftCN/NerGIWzLATsftxS+m
iqhccgAf6lwZYH+i57UZ3TVi9f8GxaRH6PDZLVqzd4ZrGXMBR1VK/QbB7hRQSHMT
xA/WLHClAoGBAKtFLlKKLub7dNmi/zpHnihFl+gsU5/uDJsvGi16dpc3yH3Hlg6i
B2GniTkzuRkk/cS+yB+CGF+nFQXtM0KohsVIbgMSAztvYxGs2RY5C8GvU/cIaG8n
QtNfn2TCBv/Qs1/Ccya/V/csu+h6rBYVBGiocRN0D1k9c0BlrXTrcale
-----END RSA PRIVATE KEY-----
`)

var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAybYWWudxZvSOA3hjtkyq
PNgaK+PkJTifrFJrCRmzW30/fi1D5Rcv/TSSxjNX7FfGvdIai7z2NvGXZVGgWGIS
UkJoWjdNZHRAp74ieGdpR1GNUY58DC8IK26br5PrHEYt5qhPSkDONi9A6ugwavvV
MtmbIM1vNtchlTcKCY2UzZOua1X1LfajmJfiBucyY/C+oEVHYhoJQvOdWVTLsxrj
JreGoAABXFj+CpGMsMK26DUnmYhGIMBBOaiyVuHNjwAlO//EuHdkt6s/c/NU7DJY
KNXRh1XypksFtgdq2xMUBag+zsJoRHtFVKVN21sP8ApZAWgI/5pj+YW99CVl9vKs
bQIDAQAB
-----END PUBLIC KEY-----
`)
