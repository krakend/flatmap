package flatmap

import (
	"testing"
)

var result interface{}

func BenchmarkIntegration(b *testing.B) {
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

		flatten, _ := Flatten(in, DefaultTokenizer)
		flatten.Move("a.*.b.*.c", "a.*.b.*.x")
		flatten.Move("b.b", "b.c")
		flatten.Del("b")
		res = flatten.Expand()
	}
	result = res
}

func BenchmarkExpand(b *testing.B) {
	var res interface{}
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
	flatten.Move("b.b", "b.c")
	flatten.Del("b")

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		res = flatten.Expand()
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

		res, _ = Flatten(in, DefaultTokenizer)
	}
	result = res
}

func BenchmarkMove(b *testing.B) {
	var res *Map

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

	res, _ = Flatten(in, DefaultTokenizer)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		if n%2 == 0 {
			res.Move("a.*.b.*.c", "a.*.b.*.x")
			res.Move("b.b", "b.c")
		} else {
			res.Move("a.*.b.*.x", "a.*.b.*.c")
			res.Move("b.c", "b.b")
		}

	}
	result = res
}
