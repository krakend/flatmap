package tree

import (
	"encoding/json"
	"fmt"
)

func ExampleNew() {
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

	tree, err := New(sample)
	tree.Sort()
	fmt.Println("error:", err)
	// fmt.Println(tree)

	tree.Del([]string{"a", "b"})
	tree.Del([]string{"a", "collection", "0"})
	tree.Del([]string{"a", "collection2", "*", "d", "*", "d"})
	fmt.Println(tree)
	b, _ := json.MarshalIndent(tree.Get([]string{}), "", "\t")
	fmt.Println(string(b))

	// output:
	// error: <nil>
	//
	// ├── a
	// │   ├── a
	// │   │   ├── a
	// │   │   │   ├── b	true
	// │   │   │   ├── c	42
	// │   │   │   └── d	tupu
	// │   │   ├── b	true
	// │   │   ├── c	42
	// │   │   └── d	tupu
	// │   ├── aa	23
	// │   ├── ab
	// │   │   ├── b	true
	// │   │   ├── c	42
	// │   │   └── d	tupu
	// │   ├── abc	23
	// │   ├── ac	23
	// │   ├── c	42
	// │   ├── collection []
	// │   │   └── 1
	// │   │       ├── b	true
	// │   │       ├── c	42
	// │   │       └── d	tupu
	// │   ├── collection2 []
	// │   │   ├── 0
	// │   │   │   └── d []
	// │   │   │       ├── 0
	// │   │   │       │   └── a	1
	// │   │   │       ├── 1	<nil>
	// │   │   │       └── 2	<nil>
	// │   │   ├── 1
	// │   │   │   └── d	[1 2 3 4 5]
	// │   │   └── 2
	// │   │       └── d	[1 2 3 4 5]
	// │   └── d	tupu
	// ├── aa	23
	// ├── foo	bar
	// ├── supu	42
	// └── tupu	false
	//
	// {
	// 	"a": {
	// 		"a": {
	// 			"a": {
	// 				"b": true,
	// 				"c": 42,
	// 				"d": "tupu"
	// 			},
	// 			"b": true,
	// 			"c": 42,
	// 			"d": "tupu"
	// 		},
	// 		"aa": 23,
	// 		"ab": {
	// 			"b": true,
	// 			"c": 42,
	// 			"d": "tupu"
	// 		},
	// 		"abc": 23,
	// 		"ac": 23,
	// 		"c": 42,
	// 		"collection": [
	// 			{
	// 				"b": true,
	// 				"c": 42,
	// 				"d": "tupu"
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
	// 	"foo": "bar",
	// 	"supu": 42,
	// 	"tupu": false
	// }
}
