package flatmap

import (
	"reflect"
	"testing"
)

func Benchmark(b *testing.B) {
	expectedRes := map[string]interface{}{
		"a": []interface{}{
			map[string]interface{}{
				"b": []interface{}{
					map[string]interface{}{
						"x": map[string]interface{}{
							"a": 1,
						},
						"aa": 1,
					},
					map[string]interface{}{
						"x": map[string]interface{}{
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
						"x": map[string]interface{}{
							"a": 1,
						},
						"aa": 1,
					},
				},
				"aa": 1,
			},
		},
		"tupu": false,
		"c":    map[string]interface{}{"a": 42},
	}

	for n := 0; n < b.N; n++ {
		in := map[string]interface{}{
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
			"b": map[string]interface{}{
				"a":  42,
				"bb": true,
				"b":  map[string]interface{}{"a": 42},
			},
		}

		flatten, _ := Flatten(in, DefaultTokenizer)
		flatten.Move("a.*.b.*.c", "a.*.b.*.x")
		flatten.Move("b.b", "c")
		flatten.Del("b")
		res := flatten.Expand()

		if !reflect.DeepEqual(res, expectedRes) {
			b.Errorf("unexpected result:\n%+v\n%+v", res, expectedRes)
		}
	}
}
