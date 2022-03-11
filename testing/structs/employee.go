package structs

type Employee struct {
	Name   interface{} `json:"name"`
	Age    int         `json:"age"`
	Salary int         `json:"salary"`
}
