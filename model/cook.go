package model

import (
	"encoding/json"
	"fmt"
	"time"
)

type User struct {
	UserID   uint64			`json:"user_id"`
	Image    string			`json:"image"`
	Email    string			`json:"email"`
	Password string			`json:"password"`
	NickName string			`json:"nickname"`
	BornDate *time.Time		`json:"borndate"`
}

type Recipe struct {
	RecipeID uint64			`json:"recipe_id"`
	UserID   uint64			`json:"user_id"`
	Name    string			`json:"name"`
	Image    string			`json:"image"`
	Description    []Step	`json:"description"`
	Products []Product		`json:"products"`
	Category string			`json:"category"`
	Tag string				`json:"tag"`
}

type Product struct{
	ProductID uint `json:"product_id"`
 	Name string `json:"name"`
	Weight uint  `json:"weight"`
}
type Step struct{
	StepNumber uint `json:"step_number"`
	Step string `json:"step"`
}


func (o *Recipe) MarshalDescription() string {
    js, err := json.Marshal(o.Description[0])
    if err != nil {
        panic(err)
    }
    items := string(js)
    for _, item := range o.Description[1:] {
        js, err := json.Marshal(item)
        if err != nil {
            panic(err)
        }
        items += fmt.Sprintf(", %v", string(js))
    }
    return fmt.Sprintf(`{"Description":[%v]}`, items)
}

func (o *Recipe) MarshalProducts() string {
    js, err := json.Marshal(o.Products[0])
    if err != nil {
        panic(err)
    }
    items := string(js)
    for _, item := range o.Products[1:] {
        js, err := json.Marshal(item)
        if err != nil {
            panic(err)
        }
        items += fmt.Sprintf(", %v", string(js))
    }
    return fmt.Sprintf(`{"Products":[%v]}`, items)
}

func (o *Recipe) UnmarshalProducts(data string) {
    json.Unmarshal([]byte(data), o)
}

func (o *Recipe) UnmarshalSteps(data string) {
    json.Unmarshal([]byte(data), o)
}
