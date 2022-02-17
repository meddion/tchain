package crypto

type HashFunc func([]byte) (HashValue, error)

func MerkleRoot(hashFunc HashFunc, values [][]byte) (HashValue, error) {
	switch len(values) {
	case 0:
		return hashFunc([]byte{})
	case 1:
		return hashFunc(values[0])
	}

	// TODO: simplify
	var hashes []HashValue
	if len(values)%2 != 0 {
		hashes = make([]HashValue, len(values)+1)
		values = append(values, values[len(values)-1])
	} else {
		hashes = make([]HashValue, len(values))
	}

	for i, v := range values {
		hash, err := hashFunc(v)
		if err != nil {
			return HashValue{}, nil
		}
		hashes[i] = hash
	}

	for len(hashes) > 1 {
		if len(hashes)%2 != 0 {
			hashes[len(hashes)] = hashes[len(hashes)-1]
			hashes = hashes[:len(hashes)+1]
		}

		for i, j := 0, 0; i < len(hashes); i, j = i+2, j+1 {
			buf := make([]byte, 0, HashLen*2)
			buf = append(buf, hashes[i][:]...)
			buf = append(buf, hashes[i+1][:]...)

			hash, err := hashFunc(buf)
			if err != nil {
				return HashValue{}, nil
			}
			hashes[j] = hash
		}

		hashes = hashes[:len(hashes)/2]
	}

	return hashes[0], nil
}
