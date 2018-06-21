package astconf_test

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/scjalliance/astconf"
)

type Name string

func (n Name) SectionName() string {
	return strings.ToLower(strings.Replace(string(n), " ", "_", -1))
}

type Elephant struct {
	Name Name `astconf:"elephant_name"`
	Age  int  `astconf:"age"`
}

func (elephant *Elephant) SectionName() string {
	return "elephant." + elephant.Name.SectionName()
}

func (elephant *Elephant) MarshalAsteriskPreamble(e *astconf.Encoder) error {
	var kind string
	switch {
	case elephant.Age > 50:
		kind = "old_elephant"
	case elephant.Age < 10:
		kind = "young_elephant"
	default:
		kind = "elephant"
	}
	return e.Printer().Setting("type", kind)
}

type Zookeeper struct {
	Name       Name `astconf:"zookeeper_name"`
	Experience int  `astconf:"experience_level"`
}

func (zk Zookeeper) SectionName() string {
	return "zookeeper." + zk.Name.SectionName()
}

func (*Zookeeper) MarshalAsteriskPreamble(e *astconf.Encoder) error {
	return e.Printer().Setting("type", "zookeeper")
}

type Zoo struct {
	Name        Name `astconf:"zoo_name"`
	Maintainers []*Zookeeper
	Elephants   []Elephant
}

func (zoo *Zoo) SectionName() string {
	return "zoo"
}

func (*Zoo) MarshalAsteriskPreamble(e *astconf.Encoder) error {
	return e.Printer().Setting("type", "zoo")
}

func ExampleEncoder() {
	var buf bytes.Buffer
	e := astconf.NewEncoder(&buf, astconf.AlignRight)
	err := e.Encode(&Zoo{
		Name: "Malarky McFee's Mighty Jungle",
		Elephants: []Elephant{
			Elephant{Name: "Matilda", Age: 47},
			Elephant{Name: "Franklin", Age: 52},
			Elephant{Name: "Georgey the Kid", Age: 5},
		},
		Maintainers: []*Zookeeper{
			&Zookeeper{
				Name:       "Gershwin McFee",
				Experience: 8000,
			},
		},
	})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Print(buf.String())
	}
	// Output:
	// [zoo]
	//     type = zoo
	// zoo_name = Malarky McFee's Mighty Jungle
	//
	// [zookeeper.gershwin_mcfee]
	//             type = zookeeper
	//   zookeeper_name = Gershwin McFee
	// experience_level = 8000
	//
	// [elephant.matilda]
	//          type = elephant
	// elephant_name = Matilda
	//           age = 47
	//
	// [elephant.franklin]
	//          type = old_elephant
	// elephant_name = Franklin
	//           age = 52
	//
	// [elephant.georgey_the_kid]
	//          type = young_elephant
	// elephant_name = Georgey the Kid
	//           age = 5
}
