package flatmap

import (
	"encoding/json"
	"fmt"
)

func ExampleFlatten() {
	sample := map[string]interface{}{
		"supu": 42,
		"tupu": false,
		"foo":  "bar",
		"a": map[string]interface{}{
			"b": true,
			"c": 42,
			"d": "tupu",
		},
		"collection": []interface{}{
			map[string]interface{}{
				"b": false,
				"d": "foobar",
			},
			map[string]interface{}{
				"b": true,
				"c": 42,
				"d": "tupu",
			},
		},
	}

	res, _ := Flatten(sample, DefaultTokenizer)
	b, _ := json.MarshalIndent(res.m, "", "\t")
	fmt.Println(string(b))

	// output:
	// {
	// 	"a.b": true,
	// 	"a.c": 42,
	// 	"a.d": "tupu",
	// 	"collection.#": 2,
	// 	"collection.0.b": false,
	// 	"collection.0.d": "foobar",
	// 	"collection.1.b": true,
	// 	"collection.1.c": 42,
	// 	"collection.1.d": "tupu",
	// 	"foo": "bar",
	// 	"supu": 42,
	// 	"tupu": false
	// }
}

func ExampleFlatten_collection() {
	sample := map[string]interface{}{
		"a": []interface{}{
			map[string]interface{}{
				"b": []interface{}{
					map[string]interface{}{
						"c": map[string]interface{}{
							"a": 1,
						},
						"aa": 1,
					},
					map[string]interface{}{
						"c": map[string]interface{}{
							"a": 2,
						},
						"aa": 1,
					},
				},
				"aa": 1,
			},
			map[string]interface{}{
				"b": []interface{}{
					map[string]interface{}{
						"c": map[string]interface{}{
							"a": 1,
						},
						"aa": 1,
					},
				},
				"aa": 1,
			},
		},
		"tupu": false,
	}
	res, _ := Flatten(sample, DefaultTokenizer)

	b, _ := json.MarshalIndent(res.m, "", "\t")
	fmt.Println(string(b))

	// output:
	// {
	// 	"a.#": 2,
	// 	"a.0.aa": 1,
	// 	"a.0.b.#": 2,
	// 	"a.0.b.0.aa": 1,
	// 	"a.0.b.0.c.a": 1,
	// 	"a.0.b.1.aa": 1,
	// 	"a.0.b.1.c.a": 2,
	// 	"a.1.aa": 1,
	// 	"a.1.b.#": 1,
	// 	"a.1.b.0.aa": 1,
	// 	"a.1.b.0.c.a": 1,
	// 	"tupu": false
	// }
}
