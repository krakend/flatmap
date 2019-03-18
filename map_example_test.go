package flatmap

import (
	"encoding/json"
	"fmt"
)

func ExampleMap() {
	sample := map[string]interface{}{
		"supu": 42,
		"tupu": false,
		"foo":  "bar",
		"aa":   23,
		"a": map[string]interface{}{
			"aa":  23,
			"ac":  23,
			"abc": 23,
			"a": map[string]interface{}{
				"b": true,
				"c": 42,
				"d": "tupu",
			},
			"ab": map[string]interface{}{
				"b": true,
				"c": 42,
				"d": "tupu",
			},
			"b": true,
			"c": 42,
			"d": "tupu",
			"collection": []interface{}{
				map[string]interface{}{
					"b":  false,
					"d":  "aaaa",
					"dx": "foobar",
				},
				map[string]interface{}{
					"b": true,
					"c": 42,
					"d": "tupu",
				},
			},
			"collection2": []interface{}{
				map[string]interface{}{
					"d": []interface{}{
						map[string]interface{}{
							"a": 1,
							"d": []int{1, 2, 3, 4, 5},
						},
						map[string]interface{}{
							"d": []int{1, 2, 3, 4, 5},
						},
						map[string]interface{}{
							"d": []int{1, 2, 3, 4, 5},
						},
					},
				},
				map[string]interface{}{
					"d": []int{1, 2, 3, 4, 5},
				},
				map[string]interface{}{
					"d": []int{1, 2, 3, 4, 5},
				},
			},
		},
	}

	res, _ := Flatten(sample, DefaultTokenizer)

	res.Del("a.b")
	res.Del("a.collection.1.b")
	res.Del("a.collection.*.d")
	res.Del("a.collection2.*.d.*.d")
	res.Move("tupu", "tuputupu")
	res.Move("a.c", "a.cb")
	res.Move("a.a", "acb")
	res.Del("a.ab")

	b, _ := json.MarshalIndent(res.m, "", "\t")
	fmt.Println(string(b))

	// fmt.Println(res.Expand())

	// output:
	// {
	// 	"a.aa": 23,
	// 	"a.abc": 23,
	// 	"a.ac": 23,
	// 	"a.cb": 42,
	// 	"a.collection.#": 2,
	// 	"a.collection.0.b": false,
	// 	"a.collection.0.dx": "foobar",
	// 	"a.collection.1.c": 42,
	// 	"a.collection2.#": 3,
	// 	"a.collection2.0.d.#": 3,
	// 	"a.collection2.0.d.0.a": 1,
	// 	"a.collection2.1.d": [
	// 		1,
	// 		2,
	// 		3,
	// 		4,
	// 		5
	// 	],
	// 	"a.collection2.2.d": [
	// 		1,
	// 		2,
	// 		3,
	// 		4,
	// 		5
	// 	],
	// 	"a.d": "tupu",
	// 	"aa": 23,
	// 	"acb.b": true,
	// 	"acb.c": 42,
	// 	"acb.d": "tupu",
	// 	"foo": "bar",
	// 	"supu": 42,
	// 	"tuputupu": false
	// }
}
