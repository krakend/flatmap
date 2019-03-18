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
				"a": map[string]interface{}{
					"b": true,
					"c": 42,
					"d": "tupu",
				},
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

	flatten, _ := Flatten(sample, DefaultTokenizer)

	flatten.Del("a.b")
	flatten.Del("a.*.a")
	flatten.Del("a.collection.1.b")
	flatten.Del("a.collection.*.d")
	flatten.Del("a.collection2.*.d.*.d")
	flatten.Move("tupu", "tuputupu")
	flatten.Move("a.c", "a.cb")
	flatten.Move("a.a", "acb")
	flatten.Del("a.ab")

	res := flatten.Expand()

	b, _ := json.MarshalIndent(res, "", "\t")
	fmt.Println(string(b))

	// output:
	// {
	// 	"a": {
	// 		"aa": 23,
	// 		"abc": 23,
	// 		"ac": 23,
	// 		"cb": 42,
	// 		"collection": [
	// 			{
	// 				"b": false,
	// 				"dx": "foobar"
	// 			},
	// 			{
	// 				"c": 42
	// 			}
	// 		],
	// 		"collection2": [
	// 			{
	// 				"d": [
	// 					{
	// 						"a": 1
	// 					},
	// 					null,
	// 					null
	// 				]
	// 			},
	// 			{
	// 				"d": [
	// 					1,
	// 					2,
	// 					3,
	// 					4,
	// 					5
	// 				]
	// 			},
	// 			{
	// 				"d": [
	// 					1,
	// 					2,
	// 					3,
	// 					4,
	// 					5
	// 				]
	// 			}
	// 		],
	// 		"d": "tupu"
	// 	},
	// 	"aa": 23,
	// 	"acb": {
	// 		"b": true,
	// 		"c": 42,
	// 		"d": "tupu"
	// 	},
	// 	"foo": "bar",
	// 	"supu": 42,
	// 	"tuputupu": false
	// }
}
