package flatmap

import (
	"reflect"
	"testing"
)

func TestMap_Del(t *testing.T) {
	for _, tc := range []struct {
		name    string
		pattern string
		in      map[string]interface{}
		out     map[string]interface{}
	}{
		{
			name:    "unknown",
			pattern: "abc",
			in: map[string]interface{}{
				"supu": 42,
				"tupu": false,
			},
			out: map[string]interface{}{
				"supu": 42,
				"tupu": false,
			},
		},
		{
			name:    "plain",
			pattern: "supu",
			in: map[string]interface{}{
				"supu": 42,
				"tupu": false,
			},
			out: map[string]interface{}{"tupu": false},
		},
		{
			name:    "element_in_struct",
			pattern: "internal.supu",
			in: map[string]interface{}{
				"internal": map[string]interface{}{
					"supu": 42,
					"tupu": false,
				},
				"tupu": false,
			},
			out: map[string]interface{}{
				"internal.tupu": false,
				"tupu":          false,
			},
		},
		{
			name:    "struct",
			pattern: "internal",
			in: map[string]interface{}{
				"internal": map[string]interface{}{
					"supu": 42,
					"tupu": false,
				},
				"tupu": false,
			},
			out: map[string]interface{}{
				"tupu": false,
			},
		},
		{
			name:    "element_in_substruct",
			pattern: "internal.internal.supu",
			in: map[string]interface{}{
				"internal": map[string]interface{}{
					"supu": 42,
					"tupu": false,
					"internal": map[string]interface{}{
						"supu": 42,
						"tupu": false,
					},
				},
				"tupu": false,
			},
			out: map[string]interface{}{
				"internal.supu":          42,
				"internal.tupu":          false,
				"internal.internal.tupu": false,
				"tupu":                   false,
			},
		},
		{
			name:    "similar_names",
			pattern: "a.a.a",
			in: map[string]interface{}{
				"a": map[string]interface{}{
					"a": map[string]interface{}{
						"a": map[string]interface{}{
							"a": 1,
						},
						"aa": 1,
					},
					"aa": 1,
				},
				"tupu": false,
			},
			out: map[string]interface{}{
				"a.aa":   1,
				"a.a.aa": 1,
				"tupu":   false,
			},
		},
		{
			name:    "collection_element_atributes",
			pattern: "a.*.a",
			in: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{
						"a": map[string]interface{}{
							"a": map[string]interface{}{
								"a": 1,
							},
							"aa": 1,
						},
						"aa": 1,
					},
					map[string]interface{}{
						"a":  42,
						"aa": 1,
					},
				},
				"tupu": false,
			},
			out: map[string]interface{}{
				"a.#":    2,
				"a.1.aa": 1,
				"a.0.aa": 1,
				"tupu":   false,
			},
		},
		{
			name:    "nested_collection_element_atributes",
			pattern: "a.*.b.*.c",
			in: map[string]interface{}{
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
			},
			out: map[string]interface{}{
				"a.#":        2,
				"a.1.aa":     1,
				"a.0.aa":     1,
				"a.0.b.#":    2,
				"a.0.b.0.aa": 1,
				"a.0.b.1.aa": 1,
				"a.1.b.#":    1,
				"a.1.b.0.aa": 1,
				"tupu":       false,
			},
		},
		{
			name:    "large_collection_element_atributes",
			pattern: "a.*.a",
			in: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{
						"a":  1,
						"aa": 1,
					},
					map[string]interface{}{
						"a":  2,
						"aa": 1,
					},
					map[string]interface{}{
						"a":  1,
						"aa": 1,
					},
					map[string]interface{}{
						"a":  2,
						"aa": 1,
					},
					map[string]interface{}{
						"a":  1,
						"aa": 1,
					},
					map[string]interface{}{
						"a":  2,
						"aa": 1,
					},
					map[string]interface{}{
						"a":  1,
						"aa": 1,
					},
					map[string]interface{}{
						"a":  2,
						"aa": 1,
					},
					map[string]interface{}{
						"a":  1,
						"aa": 1,
					},
					map[string]interface{}{
						"a":  2,
						"aa": 1,
					},
					map[string]interface{}{
						"a":  1,
						"aa": 1,
					},
					map[string]interface{}{
						"a":  2,
						"aa": 1,
					},
				},
				"tupu": false,
			},
			out: map[string]interface{}{
				"a.#":     12,
				"a.0.aa":  1,
				"a.1.aa":  1,
				"a.2.aa":  1,
				"a.3.aa":  1,
				"a.4.aa":  1,
				"a.5.aa":  1,
				"a.6.aa":  1,
				"a.7.aa":  1,
				"a.8.aa":  1,
				"a.9.aa":  1,
				"a.10.aa": 1,
				"a.11.aa": 1,
				"tupu":    false,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			res, _ := Flatten(tc.in, DefaultTokenizer)

			res.Del(tc.pattern)
			if !reflect.DeepEqual(res.m, tc.out) {
				t.Errorf("unexpected result (%s):\n%+v\n%+v", tc.pattern, res.m, tc.out)
			}
		})
	}
}

func TestMap_Move(t *testing.T) {
	for _, tc := range []struct {
		name string
		src  string
		dst  string
		in   map[string]interface{}
		out  map[string]interface{}
	}{
		{
			name: "plain",
			src:  "a",
			dst:  "b",
			in:   map[string]interface{}{"a": 42},
			out:  map[string]interface{}{"b": 42},
		},
		{
			name: "from_struct",
			src:  "b.a",
			dst:  "c",
			in: map[string]interface{}{
				"a": 42,
				"b": map[string]interface{}{"a": 42},
			},
			out: map[string]interface{}{"a": 42, "c": 42},
		},
		{
			name: "from_collection",
			src:  "b.*.c",
			dst:  "b.*.x",
			in: map[string]interface{}{
				"a": 42,
				"b": []interface{}{
					map[string]interface{}{"c": 42},
					map[string]interface{}{"c": map[string]interface{}{"d": 42}},
				},
			},
			out: map[string]interface{}{
				"a":       42,
				"b.#":     2,
				"b.0.x":   42,
				"b.1.x.d": 42,
			},
		},
		{
			name: "from_struct_nested",
			src:  "b.b",
			dst:  "c",
			in: map[string]interface{}{
				"a": 42,
				"b": map[string]interface{}{
					"a":  42,
					"bb": true,
					"b":  map[string]interface{}{"a": 42},
				},
			},
			out: map[string]interface{}{"a": 42, "b.a": 42, "b.bb": true, "c.a": 42},
		},
		{
			name: "collection",
			src:  "a.*.b",
			dst:  "a.*.c",
			in: map[string]interface{}{
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
			},
			out: map[string]interface{}{
				"a.#":         2,
				"a.0.aa":      1,
				"a.0.c.#":     2,
				"a.0.c.0.aa":  1,
				"a.0.c.0.c.a": 1,
				"a.0.c.1.aa":  1,
				"a.0.c.1.c.a": 2,
				"a.1.aa":      1,
				"a.1.c.#":     1,
				"a.1.c.0.aa":  1,
				"a.1.c.0.c.a": 1,
				"tupu":        false,
			},
		},
		{
			name: "recursive_collection",
			src:  "a.*.b.*.c",
			dst:  "a.*.b.*.x",
			in: map[string]interface{}{
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
			},
			out: map[string]interface{}{
				"a.#":         2,
				"a.0.aa":      1,
				"a.0.b.#":     2,
				"a.0.b.0.aa":  1,
				"a.0.b.0.x.a": 1,
				"a.0.b.1.aa":  1,
				"a.0.b.1.x.a": 2,
				"a.1.aa":      1,
				"a.1.b.#":     1,
				"a.1.b.0.aa":  1,
				"a.1.b.0.x.a": 1,
				"tupu":        false,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			res, _ := Flatten(tc.in, DefaultTokenizer)

			res.Move(tc.src, tc.dst)
			if !reflect.DeepEqual(res.m, tc.out) {
				t.Errorf("unexpected result (%s -> %s):\n%+v\n%+v", tc.src, tc.dst, res.m, tc.out)
			}
		})
	}
}

func TestMap_Expand(t *testing.T) {
	m, err := newMap(DefaultTokenizer)
	if err != nil {
		t.Error(err)
		return
	}

	m.m = map[string]interface{}{
		"a.#":         2,
		"a.0.aa":      1,
		"a.0.b.#":     2,
		"a.0.b.0.aa":  1,
		"a.0.b.0.c.a": 1,
		"a.0.b.1.aa":  1,
		"a.0.b.1.c.a": 2,
		"a.1.aa":      1,
		"a.1.b.#":     1,
		"a.1.b.0.aa":  1,
		"a.1.b.0.c.a": 1,
		"tupu":        false,
	}

	res := m.Expand()

	expectedRes := map[string]interface{}{
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

	if !reflect.DeepEqual(res, expectedRes) {
		t.Errorf("unexpected result:\n%+v\n%+v", res, expectedRes)
	}
}
