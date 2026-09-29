package main

//logic for set,get,store

type store struct {
	data map[string]string
}

//initialize 
//constructor 
func newStore() *store {
	return &store{
     data: make(map[string]string),
	}
}


func (s *store) Set(key, value string) {
	s.data[key] = value
}

func (s *store) Get(key string) (string, bool) {
	value, exist := s.data[key]
	return value, exist
}

func (s *store) Delete(key string) {
	delete(s.data, key)
}
