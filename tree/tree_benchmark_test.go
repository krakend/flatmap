package tree

import (
	"testing"
)

var result interface{}

func Benchmark_integration(b *testing.B) {
	var res interface{}

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

		flatten, _ := New(in)
		flatten.Move([]string{"a", "*", "b", "*", "c"}, []string{"a", "*", "b", "*", "x"})
		flatten.Move([]string{"b", "b"}, []string{"c"})
		flatten.Del([]string{"b"})
		res = flatten.Get([]string{})
	}
	result = res
}

func BenchmarkNew(b *testing.B) {
	var res interface{}

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

		res, _ = New(in)
	}
	result = res
}

func BenchmarkMove(b *testing.B) {
	var res *Tree

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

	res, _ = New(in)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		if n%2 == 0 {
			res.Move([]string{"a", "*", "b", "*", "c"}, []string{"a", "*", "b", "*", "x"})
			res.Move([]string{"b", "b"}, []string{"b", "c"})
		} else {
			res.Move([]string{"a", "*", "b", "*", "x"}, []string{"a", "*", "b", "*", "c"})
			res.Move([]string{"b", "c"}, []string{"b", "b"})
		}

	}
	result = res
}
