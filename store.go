package kairogo

//logic for set,get,store

type store struct{
	data map[string]string
}


func (s *store) Set(key , value string){
	s.data[key]=value
}

func (s *store) Get(key string){
	
}