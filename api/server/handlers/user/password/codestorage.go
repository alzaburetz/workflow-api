package password

var Storage map[string]string

func CodeStorageInit() {
	Storage = make(map[string]string)
}

