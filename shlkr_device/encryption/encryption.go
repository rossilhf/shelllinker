package encryption

func Encrypt(key int, s string) string {
	b := []byte(s)
	n := len(b)
	var c []byte = make([]byte, 2*n, 4*n)
	j := 0

	for i := 0; i < n; i++ {
		b1 := b[i]
		b2 := int(b1) ^ key
		c1 := b2 % 19
		c2 := b2 / 19
		c1 = c1 + 46
		c2 = c2 + 46
		c[j] = byte(c1)
		c[j + 1] = byte(c2)
		j = j + 2
	}

	str := string(c)
	return str
}

func Decrypt(key int, s string) string {
	c := []byte(s)
	n := len(c)
	if n % 2 != 0 {
		return ""
	}

	n = n / 2
	var b []byte = make([]byte, n, 2*n)
	j := 0
	
	for i := 0; i < n; i++ {
		c1 := int(c[j])
		c2 := int(c[j + 1])
		j = j + 2
		c1 = c1 - 46
		c2 = c2 - 46
		b2 := c2 * 19 + c1
		b1 := b2 ^ key
		b[i] = byte(b1)
	}

	str := string(b)
	return str
}
