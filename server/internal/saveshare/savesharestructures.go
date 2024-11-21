package saveshare

type json_data struct { //make a structure json_data containing hash and data, why do we use this?
	Hash string `json:"hash"` //we are supposed to use this because we encrypt this later using random generated key
	Data string `json:"data"` //and if we decrypt this into json and hash inside of json is not same as hash user send then key was wrong
} //and data inside of decrypted data is not valid => key was wrong => user unauthorized
