package main


import(
	"fmt"
)



func main(){
      s:= newStore()

	  s.Set("name","anshit")
	  s.Set("learning","GO")

	  value,exists :=s.Get("name")

	  if exists {
		fmt.Println("value : ",value)
	  }else{
		fmt.Println("key not found")
	  }

	 value,exists =s.Get("age")

	  if exists {
		fmt.Println("value : ",value)
	  }else{
		fmt.Println("key not found")
	  }

	  s.Delete("name")

	  _,exists=s.Get("name")
	  fmt.Println("Name exists :", exists)

}